# Chirpy
HTTP Server

---

Chirpy is a HTTP server which handles the back-end of a social media network. The functionality is basic, with the ability to create user accounts and users to create "chirps" (posts). 

Here is a rough outline of the HTTP inputs (to be updated):

- Health: `GET /api/healthz`: returns OK in plain text.

- Chirps:
  - `GET /api/chirps`: Returns a list of chirps with optional query parameters `author_id` (which takes a user ID and returns chirps made by that user) and `sort` (which takes `asc` or `desc`, with default `asc`, sorting by date posted).
  - `POST /api/chirps`: Takes a JSON object with "body" and "user_id", with HTTP header "Authorization" of the shape "Bearer <token>" to create a chirp under the user `user_id` with the contents of `body`.
  - `GET /api/chirps/{chirpID}`: Returns chirp with ID `chirpID`.
  - `DELETE /api/chirps/{chirpID}`: Deletes the chirp with ID `chirpID`, given that the user has authorization to, using the header "Authorization" of the shape "Bearer <token>".

- Users: 
  - `POST /api/users`: Takes a JSON object with "email" and "password" to create a user.
  - `PUT /api/users`: Takes a JSON object with optional strings "email" and "password" to update the email and/or password of the user who bears the token in the header "Authorization: Bearer <token>".
  - `POST /api/login`: Takes a JSON object with "email" and "password" to login to a user: returns the user, including the JWT in "token" and refresh token "refresh_token". 

- Webhook: `POST /api/polka/webhooks`: Takes a JSON object of the form: 
```json
{
  "event": "user.upgraded",
  "data": {
    "user_id": "3311741c-680c-4546-99f3-fc9efac2036c"
  }
}
```
to upgrade the user to "Chirpy Red", the premium version. To authorize, the API Key must be passed in the "Authorization" header, like "ApiKey f271c81ff7084ee5b99a5091b42d486e".

- Refresh Tokens;
  - `POST /api/refresh`: Returns a new JWT for the user who bears the refresh token in the header "Authorization: Bearer <token>".
  - `POST /api/revoke`: Revokes the refresh token in the header "Authorization: Bearer <token>".

- Admin:
  - `GET /admin/metrics`: Returns the number of times clients have made requests to pages with metrics included (currently only `/app`).
  - `POST /admin/reset`: Resets the databases and metrics.

---

## Chapter 1, Servers

| Component | Function |
|:-------|:-------|
|`ServeMux`| Routing requests to handlers - "\[METHOD \]\[HOST\]/\[PATH\]" |
|`http.Server`| Listens on a port and accepts connections |
|`ListenAndServe`| Starts the server |

`mux.Handle` and `mux.HandleFunc` are used to define a map from "\[METHOD \]\[HOST\]/\[PATH\]" to an object (that implements the interface `Handler`) and a function respectively. The handlers define how the server responds.


## Chapter 2, Routing

### Middleware

Middleware is a way to wrap a handler. `middlewareMetrics` has been defined to count the number of requests made and return the next handler.

### Stateful handlers

Stateful refers to the fact that the state of the server can be accessed. This is as simple as defining a config struct which will hold stateful, in-memory data and writing methods on it.

## Chapter 3, Architecture

|Monolith | Decoupled |
|:--|:--|
|A large program containing all of the functionality of the front and back-end|The frunt and back-end are separated into two different codebases|
|Somtimes REST APIs for raw data are hosted, like https://chirpy.com/api|Sometimes the back-end might be hosted on a subdomain, like https://api.chirpy.com/|
|Tightly coupled monoliths may inject dynamic data directly into HTML|Embedding is still possible, but it is more complicated|
|Simpler to get started with|Easier and cheaper to scale|
|Everything is always in sync|Can practice good separation of concerns as the codebase grows|
|Embedded data can improve SEO and UX|Can be hosted on separate servers using separate technologies|

Generally: start building monolithic but with logically decoupled front and back-end to make it easier to migrate as and when needed.

## Chapter 4, JSON

- Decoder created with `decoder := json.NewDecoder(r.Body)`; data decoded into struct with `decoder.Decode(&struct)`. Preferred when working with files, connections, or I/O operations where data is processed incrementally.
- Working with in-memory JSON, marshalling like `dat, _ := json.Marshal(struct)` is suitable.

## Chapter 5, Storage

This database uses _Postgres_, _Goose_, and _SQLC_. Goose migrations in `/sql/schema`, and queries are generated with SQLC in `/sql/queries`. Handlers merely reference the queries generated in `/internal/database`.

### REST API good practice

It's the convention to name the end points after resources, which will typically be plural. At this stage we have the following paths:

- `/api/healthz` (GET)
- `/api/chirps` (GET and POST)
- `/api/chirps/{chirpID}` (GET)
- `/api/users` (POST)
- `/admin/metrics` (GET)
- `/admin/reset` (POST)

The method is what defines how the server should respond. For example, a client can send a request to `/api/chirps`. Depending on whether this is a POST or a GET request, the server retrieves the chirps or creates a new chirp (singular), to provide consistent CRUD endpoints.


## Chapter 6, Authentication

### Hashing passwords

To hash passwords, [argon2id](https://github.com/alexedwards/argon2id) is used here: `argon2id.CreateHash(password, argon2id.DefaultParams)` and `argon2id.ComparePasswordAndHash(password, hash)` are sufficient. 

### Types of Authentication

1. Password + ID
2. 3rd Party Authentication
3. Magic Links
4. API Keys

All have their self-explanatory pros and cons.

### JWTs

A JWT is a JSON Web Token. This is a cryptographically signed JSON object containing information about the user.

The process is as follows:

1. The user logs in. (ex. `POST api/login`)
2. The server does `CheckPaswordHash` to log user with `userID` in. 
3. The server creates a JWT (`MakeJWT`) and sends it to the client
4. The user makes an API request.
5. The user's token is sent with any request it makes (e.g. in a header "Authorization" with body like "Bearer <token>").
6. The server validates the JWT (`ValidateJWT`) to ensure that who is claiming to send the message is sending the message.

So the token generation is unique to the server, a `TOKEN_SECRET` is defined in `.env`, which is just a random string.

JWTs are short-lived, stateless, and irrevocable, meaning the server doesn't need to keep track of them. They are short-lived because, since they are irrevocable, if a JWT is stolen, they can be used by anyone. To overcome the issue of them being short-lived, so users don't have to login in every time they make a new request each hour, refresh tokens can be used.

### Refresh Tokens

Refresh tokens are stateful, last longer, and can be revoked. All they do are make new JWTs. 

In our case, refresh tokens are made with user login, lasting 60 days. Whenever a JWT expires, a `POST /api/refresh/` request can be made using the refresh token in the headers, given that the refresh token hasn't expired or been revoked.


## Chapter 7, Authorization

Authorization simply refers to what the user is allowed to do. For example, a user is authorized to post a chirp with a user ID which is their user ID, but not another user's. This is done by getting the user ID of the JWT bearer to check that it matches the user ID in the request.

## Chapter 8, Webhooks

Webhooks are just idempotent HTTP requests where the client is an automated system.

For example, if we use an external service to accept payments, when a user sends a payment, the service makes a request to an exposed HTTP input, like `/api/polka/webhooks` (where `polka` is the service which accepts payments). The main difference is that the client defines the API contract.

The webhook can be made secure by defining an API Key in `.env`.
