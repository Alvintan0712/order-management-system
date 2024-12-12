package broker

import "context"

type Broker interface {
	Start(ctx context.Context) error
	Close(ctx context.Context) error
}
