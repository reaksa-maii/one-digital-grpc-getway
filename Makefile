clean-book-v1:
	rm proto/book/v1/*.go

clean-book-v2:
	rm proto/book/v2/*.go

book-v1:
	buf generate proto/book/v1/*.proto

book-v2:
	buf generate proto/book/v2/*.proto

podcast-v1:
	buf generate proto/podcast/v1/*.proto

podcast-v2:
	buf generate proto/podcast/v2/*.proto

echo-v1:
	buf generate proto/echo/v1/*.proto
echo-v2:
	buf generate proto/echo/v2/*.proto

echo-v3:
	buf generate proto/echo/v3/*.proto