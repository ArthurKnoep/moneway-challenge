# Moneway Challenge

## Configration
The configuration is set from environment variables

| env  |  default  |  description  | 
|---|---|---|
| PORT  |  8080  |  Port used by the gRPC server  |
| DB_KEYSPACE  | moneway |  The keyspace used by the ScyllaDB  |
| DB_HOST | 127.0.0.1 |  ScyllaDB host  |


## Run
```bash
dep ensure
./script/gen-protoc.sh
./script/build.sh
```