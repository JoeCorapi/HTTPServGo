# GET - no body
curl http://localhost:8080/hello

# POST - with a body
curl -X POST http://localhost:8080/submit \
  -H "Content-Type: application/json" \
  -d '{"name":"alice"}'