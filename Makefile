clean:
	rm proto/helloworld/*.go

generate:
	buf generate proto/helloworld/*.proto