package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"proto.zip/studio/mux/pkg/mux"
	"proto.zip/studio/mux/pkg/muxcontext"
)

var docs map[string]string

func main() {
	docs = make(map[string]string)
	m := mux.NewHTTP()

	m.HandleFunc(http.MethodGet, "/docs/{id}", GetDoc)
	m.HandleFunc(http.MethodDelete, "/docs/{id}", DeleteDoc)
	m.HandleFunc(http.MethodPut, "/docs/{id}", PutDoc)

	err := http.ListenAndServe("127.0.0.1:8080", m)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error starting server: %s\n", err)
		os.Exit(1)
	}
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("Not Found"))
}

func GetDoc(w http.ResponseWriter, r *http.Request) {
	id := muxcontext.PathParams(r.Context())["id"]

	if doc, ok := docs[id]; ok {
		w.Write([]byte(doc))
	} else {
		notFound(w, r)
	}
}

func DeleteDoc(w http.ResponseWriter, r *http.Request) {
	id := muxcontext.PathParams(r.Context())["id"]

	if _, ok := docs[id]; ok {
		delete(docs, id)
		w.Write([]byte("ok"))
	} else {
		notFound(w, r)
	}
}

func PutDoc(w http.ResponseWriter, r *http.Request) {
	id := muxcontext.PathParams(r.Context())["id"]
	body, err := io.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Error: %s", err)
	} else {

		if _, ok := docs[id]; ok {
			w.WriteHeader(http.StatusCreated)
		} else {
			w.WriteHeader(http.StatusOK)
		}

		docs[id] = string(body)
		w.Write([]byte("ok"))
	}
}
