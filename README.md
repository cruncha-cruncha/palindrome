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

Aside: the brief says "Provide some REST API documentation (either in Readme or in-line code)", so I should make sure this is very visible.

I like using plurals (messages not message) as it's more flexible. I used a PUT instead of a PATCH, as we're effectively replacing the entire data. I think these endpoints are fairly self-descriptive. The DELETE /messages was not required by the brief, but helped with testing.

### Request / Response

Their types. Note that although isPalindrome is used internally, is_palindrome is returned to the user (all JSON uses snake_case). See `network_types.go` (TODO: link directly?) for exact definitions.

All request/response payloads are JSON.

Request payloads
```js
// POST /messages request payload
{
    "text": "some string, which will be check if it's a palindrome"
}

// PUT /messages/{id} request payload
{
    "id": 1234 // some integer
    "text": "this text will completely replace the previous text of the message with id 1234"
}

// all other endpoints do not require a payload
```

Response payloads
```js
// POST /messages response payload
{
    "id": 1 // an integer, greater than zero. guaranteed to be unique
}

// GET /messages response payload
{
    // this array is sorted in ascending order by message id. Can be empty.
    "messages": [{
        "id": 1
        "text": "some text"
        "is_palindrome": false // null / true / false
    }]
}

// DELETE /messages does not have a response payload

// GET /message/{id} response payload
{
    "text": "the text"
    "is_palindrome": true // can be true / false / null
}

// PUT /messages/{id} does not have a response payload

// DELETE /messages/{id} does not have a response payload. It returns 204, 404, or 500

```

Could have used the same response type for all messages (id, text, is_palindrome). Decided not to to de-duplicate data (keep it as minimal as possible). This is probably not the right decision; it doesn't make a lot of sense for GetMessage and GetMessages to return slightly different message data.

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

I thought to myself "this is too simple, how could I make it more challenging?". If I'm going to make people wait 4 days for this, it better be impressive.

What if calculating whether or not some text was a palindrome took a long time? As if the REST API was more like an RPC server for long-running tasks.

There are several consequences to this decision:
1. The API still needs to be responsive. It can't wait for work to finish. So we need explicit goroutines
2. Goroutines means more thought needs to go into synchronization, and concurrent / parallel safe code
3. We can't know the result of some work until later, so null values are required, or explicit in_progress status fields
4. Some things can be done quickly and immediately. Separating the immediate and the eventual is crucial

I wanted the long-running tasks to be somewhat efficient. So if two messages come in with the same text, only one calculation is required. If a message is created then immediately deleted, cancel the work being done in the background.

## Files / Architecture

- main: registered handlers to routes, starts the server
- handlers: defines all the handlers
- messages: defines Messages, which implements MessageOrchestrator
- palindromes: defines Palindromes, which implements WorkOrchestrator
- palindrome_calculation: functions for determining if a string is a palindrome, as well as Palindromes.doWork which can be artifically slowed down (via S_DELAY)
- helpers: small, self-contained functions which could be useful in several places and don't belong anywhere else. 
- shared_state: defines SharedState (which all handlers have access too), and provides the actual definition for some important interfaces and structs. I think it would be more typical to define the Message struct (for instance) in the messages.go file, but I chose to define it in here so one could get a quick overview of how it all comes together, instead of having to look in different files. See UML diagram for more details.

I'm not sure where 'orchestrator' came from, and MessageOrchestrator shares nothing in common (in terms of inheritance / composition / implementation) with WorkOrchestrator. I just like that word.

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

## Development Approach and Decsisions / Reasoning / Design

Make this a separate document?

Read the entire brief several times. First question: what language? Decided on golang as I'm good at it and know the company uses it. It's also easy to get up and running.

Next, endpoints. WHat endpoints do we need and what will they do? Need to follow REST best practices. Came up with the list as they are above, but without the DeleteAll endpoint.

Next, what's my timeline? I wanted to get this done quickly, it's a pretty simple service. But it was a Thursday, and already mostly gone. Let's aim for Monday. That gave me two extra days (Saturday and Sunday) with which to hang myself (increase complexity beyong what's needed).

With the timeline sorted out, let's hammer out the scope. How to serve? Use a framework / package, write my own, or built-in net/http? I decided on net/http, but would later add gorilla/mux for url variable support.

Next, do I want messages to be persisted to disk somehow? I decided know. It's doable, but I wanted to spend my time in other areas. The brief was unclear. If this was a real assignment, I would definitely reach out for clarification.

Next, the overall architecture of the service. How to break it into manageable pieces with nice abstractions and encapsulations? My original idea was to have a Messages struct, which would handle everything related to messages. 

Next, how to schedule my time? I like writing code first (for better or worse). I wanted to have the code finalized by noon Friday, so I could spend a lot of time on documentation. 

What are the unknowns for the code? What questions do I have? Do I need to reach out to answer them? I decided to:

1. Get a minimal API service running with all the packages I'd need.
2. Establish how to share state
3. Add non-function endpoints. It was during this step I realized the necessity for gorilla/mux.
4. Write the Messages struct
5. Make the endpoints functional
6. Add unit tests
7. Add end-to-end tests
8. Documentation

Right off the bat I knew I wanted end-to-end tests, but didn't know exactly how. I'd also never done unit testing in production go code but knew it was possible. After reading some documentation on go test, quickly integrated a bunch of those.

Earlier, branching strategy? Git commit message format? It was a given from the start that I would use Github; it's what I'm familiar with but also was requested in the brief. Could just commit everything to main. Decided to use some sort of 'feature branch' strategy, which was more like a what-part-of-the-project-i'm-on strategy. So there ended up being five branches (as of right now): main, comments, scratch, separate-palindromes, and tests.

- scratch was for figuring out all the packages, getting things running, getting an initial version out (all the way up to step 5)
- tests came along afterwards, and was for step 6 and 7
- then separate-palindromes came in, when I realized I didn't like the architecture of the program (more later)
- then finally comments, for all documentation

The separate-palindromes branch became necessary when I realized I wanted to do a big re-write on Friday. We know the plan, but how did it actually work out in practice?

The service was working by Thursday night. And the code supported an artifical delay for calculating palindromes. But as usual, in my time spent not coding, I was thinking about the code.

I realized that the Messages struct was doing too much, and I should have a Palindromes struct for long-running work. This separation of concerns would support extensibility and clarity. This came to me Thursday night, and I wanted to do a big re-write on Friday. Ended up missing the deadline of code-freeze Friday noon, ended up being Saturday noon.

Now onto documentation. A couple of things were obvious to me: I knew other languages have some commenting standards or tools which can automatically understand comments, and I wanted that. I quickly re-discovered the godoc specification and ran with that.

I took a brief detour at looking into OpenAPI docs (I've used Swaggerhub before), but decided it wasn't worth it.

Comments are good, also need a README. Left that until last so I could write with confidence. 

Then diagrams. I was dreading this. I'd used Lucid Chart and Mermaid before, but was not impressed with either (found both cumbersome). Then I realized I should figure out what diagrams I wanted? Looked around at what was possible, some I'd heard of and some I hadn't (C4, sequence, service architecture diagram, UML, flowchart with swimlanes). Decided I needed a sequence diagram for the UpdateMessage handler, as I had looked at that code several times to make sure it was correct, and a diagram would help me understand. Knew I wanted an overall architecture diagram but wasn't sure how (still not quite sure).

Now I went looking for software to help me, found draw io, which reminded me alittle bit of Lucid Chart but was open source so felt better. And they have a desktop app! Figured out how to export higher-quality pngs. Figured out how to embed a png in markdown. And we're off to the races.

Came up with a basic architecture diagram, not following any sort of specific template but I like it.

Decided a UML diagram would be helpful for the Messages and Palindromes structs, and their associates. Three different diagrams, all helpful I think. 

Now I'm finishing the README. Will go through several edits, clean up my verbiage, re-organize (while practicing a presentation). Will definitely keep tweaking it. And will go back and read the brief again to figure out what I missed.

But back to my overall approach. Questions: how can I accomplish this task simply? What questions do I have? How long will that take? How can I accomplish this task better?
