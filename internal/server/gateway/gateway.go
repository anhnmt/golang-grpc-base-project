package gateway

import (
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"github.com/xdorro/golang-grpc-base-project/internal/service"
)

func NewGatewayServer(service *service.Service) *runtime.ServeMux {
	jsonpb := &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			Multiline:       false,
			Indent:          "",
			AllowPartial:    false,
			UseProtoNames:   true,
			UseEnumNumbers:  false,
			EmitUnpopulated: false,
			Resolver:        nil,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	}

	gatewayServer := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, jsonpb),
		runtime.WithForwardResponseOption(customForwardResponse),
		runtime.WithErrorHandler(customErrorHandler),
	)

	// register Gateway Server handler
	err := service.RegisterGatewayServerHandler(gatewayServer)
	if err != nil {
		log.Fatal().Err(err).Msg("Register Gateway server failed")
	}

	return gatewayServer
}

// customForwardResponse forwards the response from the backend to the client.
func customForwardResponse(_ context.Context, w http.ResponseWriter, _ proto.Message) error {
	headers := w.Header()
	if location, ok := headers["Grpc-Metadata-Location"]; ok {
		w.Header().Set("Location", location[0])
		w.WriteHeader(http.StatusFound)
	}

	return nil
}

// customErrorHandler handles the error from the backend to the client.
func customErrorHandler(
	ctx context.Context, _ *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, _ *http.Request,
	err error,
) {
	logger := log.Err(err)

	val, ok := runtime.RPCMethod(ctx)
	if ok {
		logger.Str("method", val)
	}

	// return Internal when Marshal failed
	const fallback = `{"error": true, "message": "failed to marshal error message"}`

	var customStatus *runtime.HTTPStatusError
	if errors.As(err, &customStatus) {
		err = customStatus.Err
	}

	s := status.Convert(err)
	pb := s.Proto()

	w.Header().Del("Trailer")
	w.Header().Del("Transfer-Encoding")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	customResponse := map[string]interface{}{
		"error":   true,
		"message": pb.GetMessage(),
	}

	responseBody, merr := marshaler.Marshal(customResponse)
	if merr != nil {
		grpclog.Infof("Failed to marshal error message %q: %v", s, merr)
		if _, err = io.WriteString(w, fallback); err != nil {
			grpclog.Infof("Failed to write response: %v", err)
		}
		return
	}

	if _, err = w.Write(responseBody); err != nil {
		grpclog.Infof("Failed to write response: %v", err)
	}

	logger.Msg("Logger custom error handler")
}
