package gapi

import (
	"context"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	grpcgatewayUserAgent = "grpcgateway-user-agent"
	xForwardedFor        = "x-forwarded-for"

	grpcClient = "grpc-client"
	userAgent  = "user-agent"
)

type Metadata struct {
	UserAgent string
	ClientIp  string
}

func getMetadata(ctx context.Context) *Metadata {
	data := &Metadata{}
	md, has := metadata.FromIncomingContext(ctx)
	if has {
		if agent, has := md[grpcgatewayUserAgent]; has {
			data.UserAgent = agent[0]
		}
		if clientIP, has := md[xForwardedFor]; has {
			data.ClientIp = clientIP[0]
		}

		if agent, has := md[userAgent]; has {
			data.UserAgent = agent[0]
		}
	}

	if p, has := peer.FromContext(ctx); has {
		data.ClientIp = p.Addr.String() // grpc-client
	}
	return data
}
