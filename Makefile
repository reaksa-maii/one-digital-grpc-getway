clean:
	rm proto/helloworld/*.go

generate-v1:
	buf generate proto/v1/*.proto

generate-v2:
	buf generate proto/v2/*.proto