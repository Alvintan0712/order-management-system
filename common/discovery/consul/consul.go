package consul

import (
	"context"
	"log"
	"strconv"

	capi "github.com/hashicorp/consul/api"
)

type Registry struct {
	client *capi.Client
}

func NewRegistry(consulAddr, serviceHost, servicePort, serviceName string) (*Registry, error) {
	config := capi.DefaultConfig()
	config.Address = consulAddr

	client, err := capi.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &Registry{client}, nil
}

func (r *Registry) Register(ctx context.Context, serviceId, serviceName, serviceHost, servicePort string) error {
	portValue, err := strconv.Atoi(servicePort)
	if err != nil {
		return err
	}

	agentServiceRegistration := &capi.AgentServiceRegistration{
		ID:      serviceId,
		Name:    serviceName,
		Address: serviceHost,
		Port:    portValue,
		Check: &capi.AgentServiceCheck{
			CheckID:                        serviceId,
			TLSSkipVerify:                  true,
			TTL:                            "5s",
			Timeout:                        "1s",
			DeregisterCriticalServiceAfter: "10s",
		},
	}

	return r.client.Agent().ServiceRegister(agentServiceRegistration)
}

func (r *Registry) Deregister(ctx context.Context, serviceId, serviceName string) error {
	log.Printf("Deregistering service %s", serviceId)
	return r.client.Agent().CheckDeregister(serviceId)
}

func (r *Registry) HealthCheck(serviceId, serviceName string) error {
	return r.client.Agent().UpdateTTL(serviceId, "online", capi.HealthPassing)
}
