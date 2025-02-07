# Palindrome

This is a small REST API. Defualt port is 8090.

`go test -v`

`go run .` or `P_DELAY=10 PORT=3000 go run .`

Or
```
docker build -t liam/palindrome-demo .
docker run -p 3000:8090 liam/palindrome-demo
```

## Architecture
- gorilla mux, net/http -> every request gets a new goroutine
- handlers are all methods on a shared state object
- the shared state object has a message orchestrator
- handlers all follow a pattern of: get payload data / url variables, call message orchestrator, then return
- all request and response payloads have their own type