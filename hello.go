package hello

import (
    "fmt"
    "net/http"
)

func init() {
    http.HandleFunc("/posts", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "Mapping Sexual Violence Project")
}
