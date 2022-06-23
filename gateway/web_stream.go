package gateway

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/net/websocket"
	"google.golang.org/grpc/metadata"
)

// WebStream implements grpc.ServerStream for websocket connection.
type WebStream struct {
	*websocket.Conn

	ctx      context.Context
	inbound  runtime.Marshaler
	outbound runtime.Marshaler
}

func NewWebsocketStream(ctx context.Context, c *websocket.Conn, in, out runtime.Marshaler) *WebStream {
	return &WebStream{
		Conn:     c,
		ctx:      ctx,
		inbound:  in,
		outbound: out,
	}
}

// SetHeader sets the header metadata. It may be called multiple times.
// When call multiple times, all the provided metadata will be merged.
// All the metadata will be sent out when one of the following happens:
//  - ServerStream.SendHeader() is called;
//  - The first response is sent out;
//  - An RPC status is sent out (error or success).
func (ws *WebStream) SetHeader(metadata.MD) error {
	// TODO
	return nil
}

// SendHeader sends the header metadata.
// The provided md and headers set by SetHeader() will be sent.
// It fails if called multiple times.
func (ws *WebStream) SendHeader(metadata.MD) error {
	// TODO
	return nil
}

// SetTrailer sets the trailer metadata which will be sent with the RPC status.
// When called more than once, all the provided metadata will be merged.
func (ws *WebStream) SetTrailer(metadata.MD) {
	// TODO
}

// Context returns the context for this stream.
func (ws *WebStream) Context() context.Context {
	return ws.ctx
}

// SendMsg sends a message. On error, SendMsg aborts the stream and the
// error is returned directly.
//
// SendMsg blocks until:
//   - There is sufficient flow control to schedule m with the transport, or
//   - The stream is done, or
//   - The stream breaks.
//
// SendMsg does not wait until the message is received by the client. An
// untimely stream closure may result in lost messages.
//
// It is safe to have a goroutine calling SendMsg and another goroutine
// calling RecvMsg on the same stream at the same time, but it is not safe
// to call SendMsg on the same stream in different goroutines.
func (ws *WebStream) SendMsg(m interface{}) error {
	writer, err := ws.Conn.NewFrameWriter(websocket.TextFrame)
	if err != nil {
		return err
	}
	return ws.outbound.NewEncoder(writer).Encode(m)
}

// RecvMsg blocks until it receives a message into m or the stream is
// done. It returns io.EOF when the client has performed a CloseSend. On
// any non-EOF error, the stream is aborted and the error contains the
// RPC status.
//
// It is safe to have a goroutine calling SendMsg and another goroutine
// calling RecvMsg on the same stream at the same time, but it is not
// safe to call RecvMsg on the same stream in different goroutines.
func (ws *WebStream) RecvMsg(m interface{}) error {
	reader, err := ws.Conn.NewFrameReader()
	if err != nil {
		return err
	}
	return ws.inbound.NewDecoder(reader).Decode(m)
}
