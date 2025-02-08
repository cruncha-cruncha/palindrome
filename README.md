# Palindrome

This was intended to be a simple REST API. The design brief was:

> Create an application which manages messages and provides details about those messages, specifically whether or not a message is a palindrome. Your application should support the
following operations:
> - Create, retrieve, update, and delete a message
> - List messages

But it quickly got out of hand. I thought to myself "this is too simple, how could I make it more challenging?". What if calculating whether or not some text was a palindrome took a long time? As if the REST API was more like an RPC server for long-running tasks.

There are several consequences to this decision:
1. The API still needs to be responsive. It can't wait for work to finish. So we need explicit goroutines
2. Goroutines means more thought needs to go into synchronization, and concurrent / parallel safe code
3. We can't know the result of some work until later, so null values are required, or explicit in_progress status fields
4. Some things can be done quickly and immediately. Separating the immediate and the eventual is crucial

Like every REST API though, there needs to be some handlers. The built-in net/http doesn't support url variables, so let's use gorilla/mux. Now we can easily handle calls like GET /messages/17.

I decided not to make the data persistent. This eases set up (so other people can run the code more easily). It also eases testing (just re-start the server, and we get a blank slate).

## Endpoints

- POST /messages
- GET /messages
- DELETE /messages
- GET /messages/{id}
- PUT /messages/{id}
- DELETE /messages/{id}

I like using plurals (messages not message) as it's more flexible. I used a PUT instead of a PATCH, as we're effectively replacing the entire data. I think these endpoints are fairly self-descriptive. The DELETE /messages was not required by the brief, but helped with testing.

## Files

main, handler, messages, palindromes, palindrome_calculation, helpers, shared_state

## Setup

Default port is 8090.

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