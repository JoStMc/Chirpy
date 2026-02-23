# Chirpy
HTTP Server

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
