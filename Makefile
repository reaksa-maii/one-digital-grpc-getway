clean:
	rm proto/helloworld/*.go

generate-v1:
	buf generate proto/v1/*.proto