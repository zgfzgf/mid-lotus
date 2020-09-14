package cborrpc

import (
	"io"
	"io/ioutil"

	cbor "github.com/ipfs/go-ipld-cbor"
)

const MessageSizeLimit = 1 << 20

func WriteCborRPC(w io.Writer, obj interface{}) error {
	data, err := cbor.DumpObject(obj)
	if err != nil {
		return err
	}

	_, err = w.Write(data)
	return err
}

func ReadCborRPC(r io.Reader, out interface{}) error {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	return cbor.DecodeInto(b, out)
}
