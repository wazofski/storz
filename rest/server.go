package rest

import (
	"fmt"
	"log"
	"net/http"

	"github.com/wazofski/store"
)

type Endpoint interface {
}

type _Server struct {
}

func Server(schema store.SchemaHolder, store store.Store) Endpoint {
	http.HandleFunc("/id/{id}", idHandler)
	for k, v := range schema.ObjectMethods() {
		http.HandleFunc(
			fmt.Sprintf("/%s/{pkey}", k),
			makeObjectHandler(k, v))
		http.HandleFunc(
			fmt.Sprintf("/%s", k),
			makeTypeHandler(k))
	}

	log.Println(http.ListenAndServe(":8000", nil))

	return &_Server{}
}

type _HandlerFunc func(http.ResponseWriter, *http.Request)

func idHandler(w http.ResponseWriter, r *http.Request) {
	// vars := mux.Vars(r)
	// id := vars["id"]
}

func makeObjectHandler(t string, methods []string) _HandlerFunc {
	return pathHandler
}

func makeTypeHandler(t string) _HandlerFunc {
	return pathHandler
}

func pathHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// Serve the resource.
	case http.MethodPost:
		// Create a new record.
	case http.MethodPut:
		// Update an existing record.
	case http.MethodDelete:
		// Remove the record.
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}
