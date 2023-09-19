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

	host, err := m.NewHost("{db}.example.com")

	if err != nil {
		fmt.Printf("Error creating host: %s", err)
		os.Exit(1)
	}

	host.Handle(http.MethodGet, "/docs/{id}", http.HandlerFunc(GetDoc))
	host.Handle(http.MethodDelete, "/docs/{id}", http.HandlerFunc(DeleteDoc))
	host.Handle(http.MethodPut, "/docs/{id}", http.HandlerFunc(PutDoc))

	err = http.ListenAndServe("127.0.0.1:8080", m)
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
	db := muxcontext.HostParams(r.Context())["db"]
	id := muxcontext.PathParams(r.Context())["id"]

	key := db + "/" + id

	if doc, ok := docs[key]; ok {
		w.Write([]byte(doc))
	} else {
		notFound(w, r)
	}
}

func DeleteDoc(w http.ResponseWriter, r *http.Request) {
	db := muxcontext.HostParams(r.Context())["db"]
	id := muxcontext.PathParams(r.Context())["id"]

	key := db + "/" + id

	if _, ok := docs[key]; ok {
		delete(docs, id)
		w.Write([]byte("ok"))
	} else {
		notFound(w, r)
	}
}

func PutDoc(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Error: %s", err)
	} else {
		db := muxcontext.HostParams(r.Context())["db"]
		id := muxcontext.PathParams(r.Context())["id"]

		key := db + "/" + id

		if _, ok := docs[key]; ok {
			w.WriteHeader(http.StatusCreated)
		} else {
			w.WriteHeader(http.StatusOK)
		}

		docs[key] = string(body)
		w.Write([]byte("ok"))
	}
}
