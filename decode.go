package main

import (
	"encoding/json"
	"net/http"
)

func decodeJSON[T any](r *http.Request) (T, error) {
    var v T
    err := json.NewDecoder(r.Body).Decode(&v)
    return v, err
}
