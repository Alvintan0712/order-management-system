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

### Database

1. Install Sqlite3
2. cd database directory in each services.
3. Execute `sqlite3 <service name>.db < script.sql`

### Protobuf API

1. `cd common`
2. `make` to generate protobuf and avro go files.

### Work Directory

```sh
go work init ./common ./coordinator ./gateway ./kitchen ./menu ./order ./payment ./stock
```

### Application

## Architecture Design

![order management system](images/oms.svg)

## Tech Stack

- Go
- Air (hot reload)
- Consul (Service Discovery)
- gRPC
- Kafka
- Schema Registry (Avro format)

# References

- [GitHub](https://github.com/sikozonpc/oms-repo/blob/main/gateway/main.go)
- [YouTube](https://youtu.be/KdnxzgSNLTU?si=sBJAgPfzgljNM8kH)
