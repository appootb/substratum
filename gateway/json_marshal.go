package gateway

import (
	"io"

	"github.com/appootb/substratum/v2/util/jsonpb"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

const (
	MIMEJSON = "application/json"
)

type JSONMarshal struct{}

// Marshal marshals "v" into byte sequence.
func (j *JSONMarshal) Marshal(v interface{}) ([]byte, error) {
	data, err := jsonpb.Marshal(v)
	if err != nil {
		return nil, err
	}
	return []byte(`{"code":0,"message":"","data":` + string(data) + "}"), nil
}

// Unmarshal unmarshals "data" into "v".
// "v" must be a pointer value.
func (j *JSONMarshal) Unmarshal(data []byte, v interface{}) error {
	return jsonpb.Unmarshal(data, v)
}

// NewDecoder returns a Decoder which reads byte sequence from "r".
func (j *JSONMarshal) NewDecoder(r io.Reader) runtime.Decoder {
	return jsonpb.Marshaler.NewDecoder(r)
}

// NewEncoder returns an Encoder which writes bytes sequence into "w".
func (j *JSONMarshal) NewEncoder(w io.Writer) runtime.Encoder {
	return jsonpb.Marshaler.NewEncoder(w)
}

// ContentType returns the Content-Type which this marshaler is responsible for.
// The parameter describes the type which is being marshalled, which can sometimes
// affect the content type returned.
func (j *JSONMarshal) ContentType(_ interface{}) string {
	return MIMEJSON
}
