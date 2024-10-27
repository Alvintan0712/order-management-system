# Order Management System

## Getting Started

### Consul

1. [Install Consul](https://developer.hashicorp.com/consul/docs/install)
   1. Compiling from source
      ```sh
      git clone https://github.com/hashicorp/consul.git
      cd consul
      ```
   2. Build Consul
      ```sh
      make dev
      ```
   3. Verifying the installation
      ```sh
      consul version
      ```
2. Run Consul
   ```sh
   consul agent -dev
   ```

### Application

## Architecture Design

![order management system](images/oms.svg)

## Tech Stack

- Go
- Air (hot reload)
- Consul (Service Discovery)
- gRPC

# References

- [GitHub](https://github.com/sikozonpc/oms-repo/blob/main/gateway/main.go)
- [YouTube](https://youtu.be/KdnxzgSNLTU?si=sBJAgPfzgljNM8kH)
