package discovery

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Registry interface {
	Register(ctx context.Context, serviceId, serviceName, serviceHost, servicePort string) error
	Deregister(ctx context.Context, serviceId, serviceName string) error
	Discover(ctx context.Context, serviceName string) ([]string, error)
	HealthCheck(serviceId, serviceName string) error
}

func GenerateInstanceId(serviceName string) string {
	return fmt.Sprintf("%s-%d", serviceName, rand.New(rand.NewSource(time.Now().UnixNano())).Int())
}

func ConnectService(ctx context.Context, serviceName string, registry Registry) (*grpc.ClientConn, error) {
	addrs, err := registry.Discover(ctx, serviceName)
	if err != nil {
		return nil, err
	}

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	addr := addrs[rand.Intn(len(addrs))]
	return grpc.NewClient(addr, opts...)
}
