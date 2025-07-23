clean:
	rm proto/helloworld/*.go

generate-v1:
	buf generate proto/book/v1/*.proto

generate-v2:
	buf generate proto/v2/book/*.proto

book-v1:
	buf generate proto/v1/book/*.proto

podcast-v1:
	buf generate proto/v1/podcast/*.proto

podcast-v2:
	buf generate proto/v2/podcast/*.proto