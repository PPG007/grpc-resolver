package service

import (
	"context"
	"nameresolver/proto"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	id primitive.ObjectID
)

func init() {
	id = primitive.NewObjectID()
}

func (NameResolverService) Hello(ctx context.Context, req *proto.EmptyRequest) (*proto.StringResponse, error) {
	return &proto.StringResponse{
		Message: id.Hex(),
	}, nil
}
