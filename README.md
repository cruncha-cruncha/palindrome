# Palindrome

This was intended to be a simple REST API. As implemented, data is not persisted: stopping the server will erase all user data. The design brief was:

> Create an application which manages messages and provides details about those messages, specifically whether or not a message is a palindrome. Your application should support the
> following operations:
>
> - Create, retrieve, update, and delete a message
> - List messages

But my implementation quickly got out of hand. If actually given this brief in a work setting, I would write much simpler code. Conversely, I would spend much more time reviewing code of this complexity if it was actually intended for release. My goal was to have some fun and demonstrate my capabilities, while still satisfying the task requirements. Please keep this in mind when reviewing the repo.

## Endpoints

| Request               | Response Status    | Handler           |
| --------------------- | ------------------ | ----------------- |
| POST /messages        | 201, 400, 500      | CreateMessage     |
| GET /messages         | 200, 500           | GetAllMessages    |
| DELETE /messages      | 204, 500           | DeleteAllMessages |
| GET /messages/{id}    | 200, 400, 404, 500 | GetMessage        |
| PUT /messages/{id}    | 200, 400, 404, 500 | UpdateMessage     |
| DELETE /messages/{id} | 204, 400, 404, 500 | DeleteMessage     |

All Handlers are methods on the SharedState struct, detailed later.

_Design Notes_

- The DELETE /messages endpoint was not required, but was convenient during testing
- PUT was chosen over PATCH, as the endpoint effectively replaces a message in it's entirety
- Favour using a plural noun (messages not message), especially if there's a "get all" endpoint

### Request/Response Payloads

All request/response payloads are JSON. See [network_types.go](./network_types.go) for exact definitions.

Request payloads:

```js
// POST /messages
{
    "text": "some message text"
}

// PUT /messages/{id}
{
    "id": 123,
    "text": "updated message text"
}
```

All other endpoints do not require a request payload. Response payloads:

```js
// POST /messages
{
    "id": 123
}

// GET /messages
{
    "messages": [{
        "id": 123,
        "text": "some message text",
        "is_palindrome": false // null / true / false
    }]
}

// GET /message/{id}
{
    "text": "the text"
    "is_palindrome": true // can be true / false / null
}
```

All other endpoints do not return a response payload.

_Design Note_
Messages retrieved via GET /messages have fields [id, text, is_palindrome] while a message retrieved via GET /messages/{id} has only [text, is_palindrome]. At the time of writing, I wanted to remove redundant fields (this is also the reason why PUT doesn't respond with payload). In retrospect this was probably not a good decision: downstream (future) code would be simpler to write if messages had a consistent type with no optional fields.

### Details

Message ids are positive integers starting at 1. They're guaranteed to be unique and are not re-used. Messages in the 'messages' payload response array field of a GET /messages request are sorted in ascending order by id.

Note that all is_palindrome fields are parsed into a Golang struct field named isPalindrome, and vice-versa (in this project, JSON uses snake_case while Golang uses camelCase).

## Setup

This code was written and tested on MacOS using go1.23.6 darwin/arm64. It will not compile on go versions below 1.23.0, since [sync.Map.Clear function](https://pkg.go.dev/sync#Map.Clear) is used.

The server listens on port 8090 by default, but this is configurable (see below). To run unit tests:

```shell
go test -v
```

There are more tests in the end-to-end-testing folder, but they use Python. See [the README in there](./end-to-end-testing/README.md) for more information.

Run the server:

```shell
go run .
```

Run on port 3000 (default is 8090):

```shell
PORT=3000 go run .
```

Run with an artificial delay of 10 seconds (responses are normal speed, see [next section](#purposefully-overcomplicating-the-implementation)):

```shell
S_DELAY=10 go run .
```

Or run the server in a docker container, on port 4000:

```shell
docker build -t liam/palindrome-demo .
docker run -p 4000:8090 liam/palindrome-demo
```

## Purposefully Overcomplicating the Implementation

_AKA: things get out of hand_

How can I challenge myself? If I'm going to make people wait 4 days for a simple REST API, the results better be impressive. What if palindrome calculation took a long time? Like determing whether "racecar" was a palindrome or not took ten seconds. Imagine that the palindrome calculation is a stand-in for any time-consuming workload, and now I have to build a server which is responsive but also capable of processing long-running tasks, almost like an RPC server.

Implementing this seemed fun and challenging to me, while also being vaguely applicable to the real world. I had thought about persisting data to disk (using a plain text file, sqlite, or even postgres), but this concept was more exciting. With these new design goals in mind, let's continue.

The `S_DELAY` environment variable is used to artificially slow down the `Palindromes.doWork(msg)` method, which determines whether a string is a palindrome and then stores the result.

## Architecture

![Data Flow](./diagrams/DataFlow.drawio.png)
_Fig. 1_

Figure 1 depicts general data flow. All incoming requests hit [ListenAndServe](https://pkg.go.dev/net/http#ListenAndServe), which has been configured to use [gorilla/mux](https://github.com/gorilla/mux) (so we can use url variables). Matched requests are routed to a handler (the middle column of rectangles in Figure 1), running in a per-request goroutine. All handlers have access to a SharedState struct, through which they can access a Messages and a Palindromes struct.

Messages and Palindromes are two separate structs because they're responsible for different things. Messages methods return immediately, whereas Palindromes kicks off work that could take awhile. Currently, each handler is responsible for ensuring consistency between Messages and Palindromes, a not-ideal situation discussed in more detail later on (see Figure 2 in [Handlers](#handlers)).

The `doWork` method determines if some text is a palindrome. It may take time to calculate, so is always invoked in a new goroutine. If this code was running in production and doing real work, spawning a heavy goroutine without first checking how many are already running is *not ideal*.

### Files

- [main.go](./main.go): registered handlers to routes and starts the server (calls `ListenAndServe`)
- [handlers.go](./handlers.go): defines all the handlers
- [messages.go](./messages.go): defines `Messages`, which implements `MessageOrchestrator`
- [palindromes.go](./palindromes.go): defines `Palindromes`, which implements `WorkOrchestrator`
- [palindrome_calculation.go](./palindrome_calculation.go): defines functions for determining if text is a palindrome.
- [helpers.go](./helpers.go): small, self-contained functions which could be useful in several places and don't belong anywhere else
- [shared_state.go](./shared_state.go): defines `SharedState` and provides the actual definition for some important interfaces (like `MessageOrchestrator` and `WorkOrchestrator`) and structs (like `Message`). I think it would be more typical to define the `Message` struct (for instance) in the `messages.go` file, but I chose to define it in `shared_state.go` so we can get a quick overview of how important structs come together, instead of having to look in different files. See [Figure 3](#shared-state) for more details.

Some files have an associated x_test.go file for unit testing.

## Handlers

Every handler follows three basic steps:

1. get request data
2. do something with Messages and Palindromes
3. return response data

The first and third steps are fairly standard, it's the second step that can get tricky. The second step could be encapsulated into other functions or even another orchestrator (to make sure Messages and Palindromes stay in sync). Let's look at the `UpdateMessage` handler.

![UpdateMessage sequence diagram](./diagrams/UpdateMessage_Sequence.drawio.png)
_Fig. 2_

The above diagram details step 2 of UpdateMessage: how it interacts with Messages and Palindromes. `msg_1` and `msg_2` are the same message at two different points in time; different variables having the same message id. `msg_1` is the original, `msg_2` has updated text and hash.

Inbetween the `update (id, text)` call to `Messages` and the `add (msg_2)` call to `Palindromes`, it is possible for a message to exist without any corresponding palindrome work. This race condition is not affected by `S_DELAY`. This situation is handled by simply returning `null` for `is_palindrome`, aka `P_UNKNOWN`. It could be eliminated by replacing the `MessageOrchestrator` 'add' method in with two others: one to create a message and another to save it. The flow would then be:

1. check if the message exists (`get (id)`)
2. create a new message, using the new `create (id, text)`
3. add palindrome work (`add (msg_2)`)
4. insert the new message, using the new `insert (msg_2)`
5. remove old palindrome work (`remove (msg_1 key)`)

On insert, `Messages` would have to verify that the id of the message to be inserted does not already exist. I think this is an overall better approach because it's more flexible and prevents the race condition, but it does put more work on whatever code is calling `Messages`.

## Shared State

`SharedState` consists of `Messages` which implements `MessageOrchestrator` and `Palindromes` which implements `WorkOrchestrator`. The two interfaces share nothing in common (in terms of inheritance / composition / implementation), I just like the word 'orchestrator'. Let's look at the specifications for `Messages` and `Palindromes`.

![Messages and Palindromes UML](./diagrams/MP_UML.drawio.png)
_Fig. 3_

Both `Messages` and `Palindromes` are thread-safe. `Messages` can get by with using pre-defined data structures provided by the sync package, but `Palindromes` uses an explicit mutex as it's operations are more complex.

A `PWKey`'s `messageId` and `hash` are identical to some `Message`'s `id` and `hash`. Each `Message` has a corresponding onChange channel (stored in a `PalindromeWork`'s `listeners`) which will receive all updates on the message's palindrome work.

To determine if some message text is a palindrome, we must call `Palindromes.Add(msg)`. If the message text is already known to be a palindrome or if a calculation is already running for some identical message text (based on the presence of a matching `PalindromeWork`), no new work is started. Otherwise, `Palindromes.Add(msg)` will call `Palindromes.doWork(msg)` in a new goroutine. `doWork` may be cancelled if `Palindromes.Remove(key)` is called.

## Persistence

If `Messages` or `Palindromes` were to store data to disk, I see two possible approaches:

1. Pass a db pool/connection into the constructor. This approach is simple, and requires little modification to existing code. However the db connection could not be modified after instantiation, and I'm not sure how transactions across multiple methods could be implemented.
2. Modify the interface so all methods require a db pool/connection/tx. This approach requires modifying a lot of existing code and puts more work on the calling code (has to manage the db connection). It is more flexible, keeps the db connection in shared state, and could support transactions across methods. I would prefer this approach.

TODO: above has been edited, below has not

## Closing Thoughts

Strengths:

- Splitting fast and slow task processing into Messages and Palindromes illustrates a clean separation of concerns and is extensible.
- Two channels (onChange and cancel) can be used to safely and successfully interact with a long-running goroutine, as long as the cancel channel isn't closed prematurely. Writes are also asynchronous (they don't wait for a read) and buffered, de-coupling logic and improving overall speed.
- Despite over-complicating the implementation, development and delivery was on-schedule. I identified unknowns early on and managed scope well. I followed a simple three-step plan: coding (and exploration), then testing, and finally documentation. Each step was time-boxed to stay on track.

Learnings:

- [Gorilla/mux](https://github.com/gorilla/mux) is an http multiplexer and works with Golang's net/http instead of replacing it entirely like [gin](https://github.com/gin-gonic/gin) or [fiber](https://github.com/gofiber/fiber).
- [Draw.io](https://app.diagrams.net/) is a free, open-source, less-polished version of [Lucid Chart](https://www.lucidchart.com/pages).
- I've previously been left unsatisfied after writing a Python script for simple end-to-end testing of a REST API, but not this time. The tests are understandable, useful, and easy to write.
- By leveraging channels such that there are clear producers and consumers (writers and readers), and dis-allowing consumers to close a channel, it means that writes will never panic.
