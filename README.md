# get chucked
This project is a small backend service written in Go to integrate with the Chuck Norris API (https://api.chucknorris.io/). It supports both unauthenticated users (retrieve random jokes) and authenticated users (access personalized jokes and submit search queries).

## To Run
**Mac/Linux:**

(without logs)
```sh
make run
```

(with logs)
```sh
make dev
```

**Windows (or without Make)**
```sh
docker compose down -v
docker compose build
docker compose up -d
```
View logs with either of the following:
```
make logs
```

```
docker compose logs -f
```

Stop with either of the following:
```
make down
```

```
docker compose down -v
```
## API Endpoints

All endpoints return JSON.  

Base API path: `/api/v1`

Server listens on HTTPS only.

For reviewer convenience, this repo includes a self-signed TLS certificate to be used for local development only. Because it is self-signed, TLS verification will fail by default.

When testing the following curl commands, make sure to include **-k** to make curl skip the verification step and proceed without checking.

---

## Health

### GET /health
Service health check.

**Auth:** Not required

**Example:**
```sh
curl -k https://localhost:8080/health
```

## Jokes

### GET /api/v1/jokes/random

Returns a random joke.

**Auth:** Not required

**Example:**
```sh
curl -k https://localhost:8080/api/v1/jokes/random
```

### GET /api/v1/jokes/search

Returns a random joke based on submitted query.

**Auth:** Required

**Query Parameter:** query
  * string
  * required
  * length between 3 and 120

**Example:**
```sh
curl -k -H "Authorization: Bearer <token>" \
  "https://localhost:8080/api/v1/jokes/search?query=shark"
```

### GET /api/v1/jokes/personalized

Returns a random joke with submitted name for Chuck Norris.

**Auth:** Required

**Query Parameter:** name 
  * string
  * required
  * length between 1 and 100

**Example:**
```sh
curl -k -H "Authorization: Bearer <token>" \
  "https://localhost:8080/api/v1/jokes/personalized?name=Dave"
```

## Users and Auth

### POST /api/v1/users

Create a new user.

**Auth:** Not required

**Request Body** 
1) email 
    * string
    * required
    * length between 1 and 254 
2) password
    * string
    * required
    * length between 8 and 30

**Example:**
```sh
curl -k -X POST https://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password"}'
```

### POST /api/v1/auth/login

Authenticate a user and return an access token.

**Auth:** Not required

**Request Body** 
1) email 
    * string
    * required
    * length between 1 and 254 
2) password
    * string
    * required
    * length between 8 and 30

**Example:**
```sh
curl -k -X POST https://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password"}'
```

# Reflections
The following section is in no way meant to be a comprehensive overview of the decisions made and the rationales behind them. Rather, it's a series of observations, possible conversation starters, invitations for further discussions, suggestions, elaborations, etc.

## Approaching the project
The prompt encouraged dialogue throughout the design and implementation process, and in spite of the potent blend of craziness that inevitably accompanies the end of the year, I did my best to honor that via email. I'll replicate some of those thoughts here to give a sense of how I approach problems, my decision-making process, communication, and so on.

Given the openness of the prompt (basically, build a wrapper around someone's API), the biggest question I had initially was about the expected scope. Was I looking at a 4-hour project or a 4-day one? Would I be summarily bounced if I didn't land in the ballpark of what the team was thinking, or would I have the opportunity for a back-and-forth to explain the whats and the whys? I anticipated that our communication at this point in the year would be more asynchronous than normal and wanted to come up with a solid starting point so I could better understand the expected scope against something fairly concrete. Under ordinary circumstances, I prioritize making sure all parties are on the same page from the start, but adapted that approach to this situation.

For instance, a project on the shorter side of the time spectrum would probably forego the concept of users, as having users would open questions like how we'd want to do auth, how we'd communicate with users (we'd probably want an email service and supporting infrastructure, like 2fa, more complicated middleware, blocking spam signups), and so on. Conversely, no auth whatsoever, at least not without explanation, would probably raise red flags about my awareness of important issues around security, industy standards, and so on. I knew I'd have to balance the tradeoffs between meeting production best practices and normal expectations around a take-home assignment, erring where possible on the side of best practices.

I decided to prioritize hitting the minimum viable product as defined in the project checklist, drawing clear lines for myself about what this project would do and what it wouldn't, and making notes for myself to track the decisions. I told myself I could always expand functionality based on whether I got any feedback on my proposal. 

To prepare, I spent some time making note of things I wanted to include, thinking through ways they might intersect, identifying possible friction points and areas where I'd need to invest some additional time to compensate for less familiarity (hello, devops). Because my work for the past several years has revolved around gRPC, a mix of large monoliths, and some microservices, I returned to resources by Alex Edwards, Kat Zien-Mendes, and Ben Johnson to ground myself in the realm where this project might reasonably live. Using gRPC for this seemed very much like overkill, even though some packages make also supporting HTTP easy, so I went with the more idiomatic choice of HTTP. Choosing a flatter structure is often the gopher way on smaller projects, but, not knowing the team or company culture, I worried this approach might not be seen as appropriate for a senior-level assessment, so I used a slightly more hybrid architecture that respects the limited scope of the project while setting us up to support additional growth (I elaborate on this further down).

I decided to start by developing the project scaffolding, so that if this were a production project with a team, it would be easier for others to join and start working on more detail-oriented implementation. I laid out what I wanted the git history to look like, though in practice I diverged somewhat from my plan, drilling down when I found (or made) pockets of time to work on this during the end-of-year holiday season. This project landed roughly in the "two days of fairly active work" category -- I wanted to put my best foot forward and not short myself on the opportunity to revisit and relearn things I hadn't touched in a while and still be present in a house with two young kids on school vacation who were determined to be bouncing around the room.

## What it does, what it doesn't, and why
I made the decision to keep the project scope small so I would have sufficient time to develop the infrastructure as well as the core functionality specified in the prompt. The app provides an endpoint to get a random Chuck Norris joke (open to all), an endpoint to get a personalized Chuck Norris joke (requires auth), and an endpoint to get a random Chuck Norris joke based on an input query (requires auth). 

Currently, jokes are not associated with users, though that could certainly change in response to future requirements (implement the ability to favorite a joke, for instance). Each of the current services could benefit from more management APIs (e.g. RevokeToken) and supporting services (email service, 2fa, etc). We could develop the concept of user to include active/inactive/banned, different levels of user-permissions, or change the way we authenticate, and so on.

I chose the current token-based auth approach as a nice middle-ground between a basic auth approach and something like JWT. Having registered (and authenticated) users lets us have more meaningful metrics about who is doing what (opening up possible monetization prospects).

If we think of this project as essentially a wrapper of the Chuck Norris API, we can consider the `jokes` database table to act as a cache, speeding up response times and helping to prevent abuse of the Chuck Norris API. When an authenticated user submits a search query, we first try to find a match in the database. If we have no matches, we next call the Chuck Norris API, retrieving up to 100 records, inserting them into our database, and returning a joke to the user. Although we won't call the API again in response to subsequent requests with the same query, this seems like a large-enough sample size to satisfy users. The other two tables, `users` and `tokens`, follow a more traditional relational approach, where tokens are associated with users.

When planning the architecture, I sought to respect the limited scope of the prompt while staying aware of (and allowing for) future growth. Certain packages, like `domain` or `mock`, use a flat structure appropriate for smaller projects (having different files in these directories, each containing a few lines of code, felt like unhelpful clutter). On the other hand, I broke apart services into their own packages because that made sense structurally.

## Seeding the database
Because this project depends entirely on a 3rd-party API, it made sense to seed the database with a few jokes to ensure there'd be something to return to a user in the event of an outage (though if any API were to be unbreakable, it would be the Chuck Norris API). If this were a production project, I'd have included many more records in the initial seed, but given that part of the prompt required API integration, I used a small seed to help ensure we'd hit the API quickly during a demo.

I thought about also inserting a user and including the token in the readme, but ultimately decided that copy-pasting the curl examples isn't too much of an ask of anyone looking at the project.

## Sticking to std lib
I decided to stick mainly to the standard library, with a few exceptions like choosing zap over slog for logging and testify over `if got != want` assertions for testing. The former was to save a little time since I'm more familiar with it, the latter to save some keystrokes. Both are standard choices.

I hadn't worked with HTTP in quite a while before starting this project. As a result, I considered reaching for something to hide a little of the api boilerplate. Ultimately though, I chose to use this as an opportunity to reaquaint myself with the standard library and to see how it has evolved since I last used it for HTTP. 

## Approaches to testing
I used a range of testing techniques in this project. While I'd normally reach for `moq` to generate mocks, it felt a bit like overkill here, so I instead handrolled some simple mocks that did what I needed. When I wanted to test against an actual database, I tried out `testcontainers` for the first time. For `TestPersonalize` (joke_test.go), I used table-driven tests to ensure that no matter what variant of Chuck we got, we'd be able to replace it with the user-submitted name. With more time, I'd add some end-to-end tests for additional peace of mind, but given the scope and timeline of this project, I think the current coverage is appropriate.

## TLS Note
I chose to include a self-signed certificate used only for local development to minimize setup time for reviewers. In a production setup, certificates would be issued and managed by a trusted CA and certificate verification would not be skipped.