package gateway

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"github.com/xdorro/golang-grpc-base-project/internal/service"
)

func NewGatewayServer(service *service.Service) *runtime.ServeMux {
	jsonOption := &runtime.JSONPb{
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

	gwMux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, jsonOption),
		runtime.WithForwardResponseOption(CustomForwardResponse),
		// runtime.WithErrorHandler(CustomErrorResponse),
	)

	// register Gateway Server handler
	service.RegisterGatewayServerHandler(gwMux)

	return gwMux
}

// CustomForwardResponse forwards the response from the backend to the client.
func CustomForwardResponse(_ context.Context, w http.ResponseWriter, _ proto.Message) error {
	headers := w.Header()
	if location, ok := headers["Grpc-Metadata-Location"]; ok {
		w.Header().Set("Location", location[0])
		w.WriteHeader(http.StatusFound)
	}

	return nil
}
