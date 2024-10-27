package service

import (
	"example.com/oms/common/discovery/consul"
	"google.golang.org/grpc"
)

type Service struct {
	Name string
}

type GRPCService struct {
	Service
	Registry   *consul.Registry
	Connection *grpc.ClientConn
}

func (s *GRPCService) Close() {
	s.Connection.Close()
}
