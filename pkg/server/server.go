package server

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/sevenreup/mem-value/pkg/mem"
)

var memValue *mem.Cache

func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got / request\n")
	io.WriteString(w, "This is my website!\n")
}

func get(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")

	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	value, found := memValue.Get(key)

	if !found {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	io.WriteString(w, value.(string))
}

func add(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")

	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	memValue.Set(key, string(body), mem.DefaultExpiration)
}

func set(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")

	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	memValue.Set(key, string(body), mem.NoExpiration)
}

func delete(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")

	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	memValue.Delete(key)
}

func ServerRun() {
	memValue = mem.New(mem.DefaultExpiration, 0)
	mux := http.NewServeMux()
	mux.HandleFunc("/", getRoot)
	mux.HandleFunc("/get", get)
	mux.HandleFunc("/add", add)
	mux.HandleFunc("/set", set)
	mux.HandleFunc("/delete", delete)

	err := http.ListenAndServe(":3333", mux)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	} else {
		fmt.Printf("server is running on http://localhost:3333 \n")
	}
}
