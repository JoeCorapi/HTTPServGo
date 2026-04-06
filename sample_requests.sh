# Should get 200 + "Hello, world!"
curl -v http://localhost:8080/hello

# Should echo the body back
curl -v -X POST http://localhost:8080/echo \
  -H "Content-Type: text/plain" \
  -d "this is my body"

# Should get 404
curl -v http://localhost:8080/missing