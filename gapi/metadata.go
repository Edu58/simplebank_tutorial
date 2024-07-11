package gapi

import (
	"context"
	"log"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	GRPCGatewayUserAgent     = "grpcgateway-user-agent"
	GRPCGatewayXForwardedFor = "x-forwarded-for"
	GRPCUserAgent            = "user-agent"
)

type Metadata struct {
	userAgent string
	clientIp  string
}

func (server *Server) extractMetadata(ctx context.Context) *Metadata {
	mtdata := &Metadata{}

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		log.Printf("%v", md)

		if userAgents := md.Get(GRPCGatewayUserAgent); len(userAgents) > 0 {
			mtdata.userAgent = userAgents[0]
		}

		if clientIP := md.Get(GRPCGatewayXForwardedFor); len(clientIP) > 0 {
			mtdata.clientIp = clientIP[0]
		}

		if userAgents := md.Get(GRPCUserAgent); len(userAgents) > 0 {
			mtdata.userAgent = userAgents[0]
		}
	}

	if peer, ok := peer.FromContext(ctx); ok {
		mtdata.clientIp = peer.Addr.String()
	}

	return mtdata
}
