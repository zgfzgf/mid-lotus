package jsonrpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
)

const (
	rpcParseError     = -32700
	rpcMethodNotFound = -32601
	rpcInvalidParams  = -32602
)

type rpcHandler struct {
	paramReceivers []reflect.Type
	nParams        int

	receiver    reflect.Value
	handlerFunc reflect.Value

	hasCtx int

	errOut int
	valOut int
}

// RPCServer provides a jsonrpc 2.0 http server handler
type RPCServer struct {
	methods map[string]rpcHandler
}

// NewServer creates new RPCServer instance
func NewServer() *RPCServer {
	return &RPCServer{
		methods: map[string]rpcHandler{},
	}
}

type param struct {
	data []byte // from unmarshal

	v reflect.Value // to marshal
}

func (p *param) UnmarshalJSON(raw []byte) error {
	p.data = make([]byte, len(raw))
	copy(p.data, raw)
	return nil
}

func (p *param) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.v.Interface())
}

type request struct {
	Jsonrpc string  `json:"jsonrpc"`
	ID      *int64  `json:"id,omitempty"`
	Method  string  `json:"method"`
	Params  []param `json:"params"`
}

type respError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *respError) Error() string {
	if e.Code >= -32768 && e.Code <= -32000 {
		return fmt.Sprintf("RPC error (%d): %s", e.Code, e.Message)
	}
	return e.Message
}

type response struct {
	Jsonrpc string      `json:"jsonrpc"`
	Result  interface{} `json:"result,omitempty"`
	ID      int64       `json:"id"`
	Error   *respError  `json:"error,omitempty"`
}

// TODO: return errors to clients per spec
func (s *RPCServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.rpcError(w, &req, rpcParseError, err)
		return
	}

	handler, ok := s.methods[req.Method]
	if !ok {
		s.rpcError(w, &req, rpcMethodNotFound, fmt.Errorf("method '%s' not found", req.Method))
		return
	}

	if len(req.Params) != handler.nParams {
		s.rpcError(w, &req, rpcInvalidParams, fmt.Errorf("wrong param count"))
		return
	}

	callParams := make([]reflect.Value, 1+handler.hasCtx+handler.nParams)
	callParams[0] = handler.receiver
	if handler.hasCtx == 1 {
		callParams[1] = reflect.ValueOf(r.Context())
	}

	for i := 0; i < handler.nParams; i++ {
		rp := reflect.New(handler.paramReceivers[i])
		if err := json.NewDecoder(bytes.NewReader(req.Params[i].data)).Decode(rp.Interface()); err != nil {
			s.rpcError(w, &req, rpcParseError, err)
			return
		}

		callParams[i+1+handler.hasCtx] = reflect.ValueOf(rp.Elem().Interface())
	}

	callResult := handler.handlerFunc.Call(callParams)
	if req.ID == nil {
		return // notification
	}

	resp := response{
		Jsonrpc: "2.0",
		ID:      *req.ID,
	}

	if handler.errOut != -1 {
		err := callResult[handler.errOut].Interface()
		if err != nil {
			resp.Error = &respError{
				Code:    1,
				Message: err.(error).Error(),
			}
		}
	}
	if handler.valOut != -1 {
		resp.Result = callResult[handler.valOut].Interface()
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		fmt.Println(err)
		return
	}
}

func (s *RPCServer) rpcError(w http.ResponseWriter, req *request, code int, err error) {
	w.WriteHeader(500)
	if req.ID == nil { // notification
		return
	}

	resp := response{
		Jsonrpc: "2.0",
		ID:      *req.ID,
		Error: &respError{
			Code:    code,
			Message: err.Error(),
		},
	}

	_ = json.NewEncoder(w).Encode(resp)
}

// Register registers new RPC handler
//
// Handler is any value with methods defined
func (s *RPCServer) Register(namespace string, r interface{}) {
	val := reflect.ValueOf(r)
	//TODO: expect ptr

	for i := 0; i < val.NumMethod(); i++ {
		method := val.Type().Method(i)

		funcType := method.Func.Type()
		hasCtx := 0
		if funcType.NumIn() >= 2 && funcType.In(1) == contextType {
			hasCtx = 1
		}

		ins := funcType.NumIn() - 1 - hasCtx
		recvs := make([]reflect.Type, ins)
		for i := 0; i < ins; i++ {
			recvs[i] = method.Type.In(i + 1 + hasCtx)
		}

		valOut, errOut, _ := processFuncOut(funcType)

		fmt.Println(namespace + "." + method.Name)

		s.methods[namespace+"."+method.Name] = rpcHandler{
			paramReceivers: recvs,
			nParams:        ins,

			handlerFunc: method.Func,
			receiver:    val,

			hasCtx: hasCtx,

			errOut: errOut,
			valOut: valOut,
		}
	}
}

func processFuncOut(funcType reflect.Type) (valOut int, errOut int, n int) {
	errOut = -1
	valOut = -1
	n = funcType.NumOut()

	switch n {
	case 0:
	case 1:
		if funcType.Out(0) == errorType {
			errOut = 0
		} else {
			valOut = 0
		}
	case 2:
		valOut = 0
		errOut = 1
		if funcType.Out(1) != errorType {
			panic("expected error as second return value")
		}
	default:
		panic("too many error values")
	}

	return
}

var _ error = &respError{}
