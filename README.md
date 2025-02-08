# Palindrome

This was intended to be a simple REST API. The design brief was:

> Create an application which manages messages and provides details about those messages, specifically whether or not a message is a palindrome. Your application should support the
following operations:
> - Create, retrieve, update, and delete a message
> - List messages

But it quickly got out of hand. Please keep this in mind when reviewing the repo. I was bored. If I was actually given this brief, I would write much simpler code. Conversely, if code of this complexity was actually going into production, I would spend much more time on it. I'm looking to have some fun and demonstrate what I can do, while still satisfying the brief. More later.

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

## Setup

Default port is 8090.

Run unit tests:
```shell
go test -v
```
There are more tests in the end-to-end-testing folder, but they use Python. See the README in there for more information.

Run 
```shell
go run .
```
Run on port 3000 (default is 8090)
```shell
PORT=3000 go run .
```
Run with an artificial delay of 10 seconds (calculating palindrome will take ten seconds, but responses are regular speed)
```shell
S_DELAY=10 go run .
```
All together now
```shell
PORT=3000 S_DELAY=10 go run .
```
And of course it can be compiled
```shell
go build . && ./palindrome
```
Alternatively, run it in a docker container, on port 3000
```shell
docker build -t liam/palindrome-demo .
docker run -p 3000:8090 liam/palindrome-demo
```

## What Actually Happened

I thought to myself "this is too simple, how could I make it more challenging?".

What if calculating whether or not some text was a palindrome took a long time? As if the REST API was more like an RPC server for long-running tasks.

There are several consequences to this decision:
1. The API still needs to be responsive. It can't wait for work to finish. So we need explicit goroutines
2. Goroutines means more thought needs to go into synchronization, and concurrent / parallel safe code
3. We can't know the result of some work until later, so null values are required, or explicit in_progress status fields
4. Some things can be done quickly and immediately. Separating the immediate and the eventual is crucial

I wanted the long-running tasks to be somewhat efficient. So if two messages come in with the same text, only one calculation is required. If a message is created then immediately deleted, cancel the work being done in the background.

## Files

main, handler, messages, palindromes, palindrome_calculation, helpers, shared_state


## Architecture

- gorilla mux, net/http -> every request gets a new goroutine
- handlers are all methods on a shared state object
- the shared state object has a message orchestrator, and palindrome work orchestrator
- handlers all follow a pattern of: get payload data / url variables, call orchestrators, then return
- generally update messages before updating palindrome work, but the code will still handle cases where a message exists but it's work does not
- all request and response payloads have their own type