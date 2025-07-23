# Run get client by json

- grpcurl -plaintext -d '{"name": "your name"}' localhost:8080 helloworld.v1.GreeterService/SayHello
