# Should get 200 + "Hello, world!"
curl -v http://localhost:8080/hello

# Should echo the body back
curl -v -X POST http://localhost:8080/echo \
  -H "Content-Type: text/plain" \
  -d "this is my body"

# Should get 404
curl -v http://localhost:8080/missing

# Bash loop for load testing
for i in {1..20}; do curl -s http://localhost:8080/hello & done; wait

# Semaphore pressure test — 15 concurrent requests against a cap of 10
# Watch the terminal: goroutine count should plateau at ~10, not spike to 15
for i in {1..15}; do curl -s http://localhost:8080/hello & done; wait

# POST with a larger body — verifies Content-Length parsing still works after refactor
curl -v -X POST http://localhost:8080/echo \
  -H "Content-Type: text/plain" \
  -d "the quick brown fox jumps over the lazy dog"