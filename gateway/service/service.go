package service

import (
	"example.com/oms/common"
	"example.com/oms/common/discovery/consul"
	"google.golang.org/grpc"
)

var debug = common.EnvString("DEBUG", "false") == "true"

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
