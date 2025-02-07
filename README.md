# Palidrome

This is a small REST API. Defualt port is 8090.

`go test -v`

`go run .` or `P_DELAY=10 PORT=3000 go run .`

Or
```
docker build -t liam/palindrome-demo .
docker run -p 3000:8090 liam/palindrome-demo
```