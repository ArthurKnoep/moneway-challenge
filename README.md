# Moneway Challenge

## General
### File structure
```
./
 |-> api/proto             # Protobuf file
 |-> cmd/                  # Main file
 |-> lib
 |   `-> database          # Common database function
 |       `-> models        # Database model
 |-> pkg
 |   |-> api               # Protogen file
 |   |-> cmd/*             # Init service server
 |   |   |-> config        # Config for the service
 |   |   |-> database      # Database initialisation for the service
 |   |   `-> grpc          # Initialize grpc client
 |   |-> protocol/grpc/*   # Run service server
 |   `-> service           # Implementation of the service
 |-> proto_extern          # Extern protobuf file given by google
 `-> script                # Shell script     
```
### Build
```bash
dep ensure
./script/gen-protoc.sh
./script/build.sh
```

### Docker
A docker compose file is available
```bash
docker-compose build
docker-compose up
```

```bash
docker build -f cmd/account/Dockerfile .      # Acount service
docker build -f cmd/balance/Dockerfile .      # Balance service
docker build -f cmd/transaction/Dockerfile .  # Transaction service
```

## Services
### Account
This service allow the creation or the deletion of bank account

#### Configration
The configuration is set from environment variables

| env  |  default  |  description  | 
|---|---|---|
| ACCOUNT_PORT  |  8080  |  Port used by the gRPC server  |
| DB_KEYSPACE  | moneway |  The keyspace used by the ScyllaDB  |
| DB_HOST | 127.0.0.1 |  ScyllaDB host  |


#### Build
```bash
dep ensure
./script/gen-protoc.sh
go build -o account ./cmd/account/main.go
```

#### Proto
See `./api/proto/v1/account.proto` file

### Balance
This service allow operation on the balance of an account

#### Configuration
The configuration is set from environment variables

| env  |  default  |  description  | 
|---|---|---|
| BALANCE_PORT  |  8081  |  Port used by the gRPC server  |
| DB_KEYSPACE  | moneway |  The keyspace used by the ScyllaDB  |
| DB_HOST | 127.0.0.1 |  ScyllaDB host  |

#### Build
```bash
dep ensure
./script/gen-prootoc.sh
go build -o balance ./cmd/balance/main.go 
```

#### Proto
See `./api/proto/v1/balance.proto` file

### Transaction
This service allow transaction between two account

#### Configuration
The configuration is set from environment variables

| env  |  default  |  description  | 
|---|---|---|
| BALANCE_PORT  |  8081  |  Port used by the balance gRPC server  |
| BALANCE_HOST  |  127.0.0.1  |  Host used by the balance gRPC server  |
| TRANSACTION_PORT  |  8082  |  Host used by the gRPC server  |
| DB_KEYSPACE  | moneway |  The keyspace used by the ScyllaDB  |
| DB_HOST | 127.0.0.1 |  ScyllaDB host  |

#### Build
```bash
dep ensure
./script/gen-prootoc.sh
go build -o balance ./cmd/transaction/main.go 
```

#### Proto
See `./api/proto/v1/transaction.proto` file