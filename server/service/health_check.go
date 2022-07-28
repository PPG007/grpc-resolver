package service

import (
	"context"
	"nameresolver/proto"
)

func (NameResolverService) HealthCheck(ctx context.Context, req *proto.EmptyRequest) (*proto.EmptyResponse, error) {
	return &proto.EmptyResponse{}, nil
}
