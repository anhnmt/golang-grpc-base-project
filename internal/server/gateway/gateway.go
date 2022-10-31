package gateway

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog/log"
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
		runtime.WithForwardResponseOption(CustomForwardResponse),
		// runtime.WithErrorHandler(CustomErrorResponse),
	)

	// register Gateway Server handler
	err := service.RegisterGatewayServerHandler(gatewayServer)
	if err != nil {
		log.Fatal().Err(err).Msg("Register Gateway server failed")
	}

	return gatewayServer
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
