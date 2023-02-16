package grpcapi

import (
	context "context"

	empty "github.com/golang/protobuf/ptypes/empty"
)

type server struct {
	UnimplementedStateProviderServer
}

func (s *server) GetState(ctx context.Context, e *empty.Empty) (*StateT, error) {
	return &StateT{}, nil
}
