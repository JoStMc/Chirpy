# Chiy
HTTP Server

## Chapter 1

| Component | Function |
|:-------|:-------|
|`ServeMux`| Routing requests to handlers - "\[METHOD \]\[HOST\]/\[PATH\]" |
|`http.Server`| Listens on a port and accepts connections |
|`ListenAndServe`| Starts the server |

`mux.Handle` and `mux.HandleFunc` are used to define a map from "\[METHOD \]\[HOST\]/\[PATH\]" to an object (that implements the interface `Handler`) and a function respectively. The handlers define how the server responds.


## Chapter 2

### Middleware

Middleware is a way to wrap a handler. `middlewareMetrics` has been defined to count the number of requests made and return the next handler.

### Stateful handlers

Stateful refers to the fact that the state of the server can be accessed. This is as simple as defining a config struct which will hold stateful, in-memory data and writing methods on it.
