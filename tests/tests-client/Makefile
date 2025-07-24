test-local:
	grpcurl -plaintext -d '{"name": "your name"}' localhost:8080 book.v1.GreeterService/SayHello

health:
	grpcdebug localhost:8080 health

health-service:
	grpcdebug localhost:50051 health helloworld.Greeter

channelz:
	grpcdebug localhost:50051 channelz servers --json
channelz-client:
	grpcdebug localhost:50051 channelz channels

channelz-server:
	grpcdebug localhost:50051 channelz servers

channels-7:
	grpcdebug localhost:50051 channelz channel 7

xds:
	grpcdebug localhost:50051 xds status