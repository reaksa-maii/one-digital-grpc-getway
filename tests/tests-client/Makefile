test-local:
	grpcurl -plaintext -d '{"name": "your name"}' localhost:8080 book.v1.GreeterService/SayHello

test-debug:
	grpcdebug localhost:8080 health