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

### Request / Response

Their types. Note that although isPalindrome is used internally, is_palindrome is returned to the user (all JSON uses snake_case).

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
- handlers do a lot of work. That work could be encapsulated into other functions or even another orchestrator (to make sure Messages and Palindromes stay in sync). This approach was not taken due to time constraints.

## Diagrams

![Data Flow](./diagrams/DataFlow.drawio.png)

This is the general data flow through the server. All requests come in and hit ListenAndServe. gorilla/mux routes the request to a handler (the middle column of rectangles). Each handler gets their own goroutine. All handlers have access to SharedState, through which they can access Messages and Palindromes. Calling the Palindromes.Add method may spawn a new doWork goroutine. 

I don't know exactly when / how many goroutines / threads are spawned / used by net/http (especially how responses are returned).

![UpdateMessage sequence diagram](./diagrams/UpdateMessage_Sequence.drawio.png)

The above diagram shows how the UpdateMessage handler interacts with Messages (a MessageOrchestrator) and Palindromes (a WorkOrchestrator). `msg_1` and `msg_2` are different variables but have the same message id. `msg_1` is the original, `msg_2` has updated text and hash.

Inbetween the `update (id, text)` call to Messages and the `add (msg_2)` call to Palindromes, a message exists without any corresponding palindrome work. This is handled by the code by returning a status of P_UNKNOWN, aka it is unknown if the message text is a palindrome or not. On the client side, is_palindrome will be null. This race condition exists even when `S_DELAY=0`. It could be eliminated by splitting the MessageOrchestrator 'add' into two functions: one to create a message and another to save it. The flow would then be:

1. check if the message exists (`get (id)`)
2. create a new message, using the new `create (id, text)`.
3. add palindrome work (`add (msg_2)`)
4. insert the new message, using the new `insert (msg_2)` function
5. remove old palindrome work (`remove (msg_1 key)`)

On insert, Messages would have to verify that the id of the message to be inserted does not already exist. I think this is an overall better approach in that it's more flexible and prevents the race condition, but it puts more work whatever code is calling Messages. 

![Messages and Palindromes UML](./diagrams/MP_UML.drawio.png)

If Messages and Palindromes were to be persistent and store data to a database, I see two possible approaches:

1. Pass a db pool/connection into the constructor. This approach is simple, and requires little modification to existing code. However the db connection could not be modified after instantiation, and I'm not sure how transactions across multiple methods could be implemented.
2. Modify the interface so all methods require a db pool/connection/tx. This approach requires modifying a lot of existing code and puts more work on the calling code (has to manage the db connection). It is more flexible, keeps the db connection in shared state, and could support transactions across methods. I would prefer this approach.
