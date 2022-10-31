package gateway

import (
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/protobuf/encoding/protojson"
)

func NewGatewayServer() *runtime.ServeMux {
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
		// runtime.WithForwardResponseOption(CustomForwardResponse),
		// runtime.WithErrorHandler(CustomErrorResponse),
	)

	return gwMux
}
