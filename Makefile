clean-book-v1:
	rm proto/book/v1/*.go

clean-book-v2:
	rm proto/book/v2/*.go

book-v1:
	buf generate proto/book/v1/*.proto

book-v2:
	buf generate proto/book/v2/*.proto
	
book-3:
	buf generate proto/book/v3/*.proto

podcast-v1:
	buf generate proto/podcast/v1/*.proto

podcast-v2:
	buf generate proto/podcast/v2/*.proto

echo-v1:
	buf generate proto/echo/v1/*.proto
echo-v2:
	buf generate proto/echo/v2/*.proto

echo-v3:
	buf generate proto/echo/v3/*.
	
echo-v4:
	buf generate proto/echo/v4/*.proto

echo-v5:
	buf generate proto/echo/v5/*.proto

	
movies-v1:
	buf generate proto/movie/v1/*.

movies-v2:
	buf generate proto/movie/v2/*.proto

movies-v3:
	buf generate proto/movie/v3/*.proto

google-v1:
	buf generate proto/google/v1/*.proto