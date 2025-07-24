# Run get client by json

## Local

- grpcurl -plaintext -d '{"name": "your name"}' localhost:8080 book.v1.GreeterService/SayHello

## Production

- grpcurl -plaintext -d '{"name": "your name"}' getway-grpc.itedev.online:8080 book.v1.GreeterService/SayHello
