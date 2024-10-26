package consul

import (
	"context"
	"fmt"
	"log"
	"strconv"

	capi "github.com/hashicorp/consul/api"
)

type Registry struct {
	client *capi.Client
}

// Create new consul client
func NewRegistry(consulAddr, serviceHost, servicePort, serviceName string) (*Registry, error) {
	config := capi.DefaultConfig()
	config.Address = consulAddr

	client, err := capi.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &Registry{client}, nil
}

// Register address to consul
// https://developer.hashicorp.com/consul/api-docs/agent/service#register-service
func (r *Registry) Register(ctx context.Context, serviceId, serviceName, serviceHost, servicePort string) error {
	portValue, err := strconv.Atoi(servicePort)
	if err != nil {
		return err
	}

	// This part will load in the http body
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

	// This method will call PUT /v1/agent/service/register to register your service to consul
	return r.client.Agent().ServiceRegister(agentServiceRegistration)
}

// https://developer.hashicorp.com/consul/api-docs/agent/check#deregister-check
func (r *Registry) Deregister(ctx context.Context, serviceId, serviceName string) error {
	log.Printf("Deregistering service %s", serviceId)
	// This method will call PUT /v1/agent/check/deregister/{checkID} to deregister your service to consul
	return r.client.Agent().CheckDeregister(serviceId)
}

// https://developer.hashicorp.com/consul/api-docs/agent/check#ttl-check-update
func (r *Registry) HealthCheck(serviceId, serviceName string) error {
	return r.client.Agent().UpdateTTL(serviceId, "online", capi.HealthPassing)
}

// https://developer.hashicorp.com/consul/api-docs/health#list-checks-for-service
func (r *Registry) Discover(ctx context.Context, serviceName string) ([]string, error) {
	entries, _, err := r.client.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return nil, err
	}

	var instances []string
	for _, entry := range entries {
		instances = append(instances, fmt.Sprintf("%s:%d", entry.Service.Address, entry.Service.Port))
	}

	return instances, nil
}
