package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/exp/slices"

	"github.com/wazofski/store"
	"github.com/wazofski/store/options"
	"github.com/wazofski/store/utils"
)

const (
	FilterArg      = "filter"
	IncrementalArg = "inc"
	PageSizeArg    = "pageSize"
	PageOffsetArg  = "pageOffset"
	OrderByArg     = "orderBy"
)

type Endpoint interface {
}

type _HandlerFunc func(http.ResponseWriter, *http.Request)

type _Server struct {
	Schema  store.SchemaHolder
	Store   store.Store
	Context context.Context
}

func Server(schema store.SchemaHolder, store store.Store) Endpoint {
	server := &_Server{
		Schema:  schema,
		Store:   store,
		Context: context.Background(),
	}

	http.HandleFunc("/id/{id}", makeIdHandler(server))
	for k, v := range schema.ObjectMethods() {
		http.HandleFunc(
			fmt.Sprintf("/%s/{pkey}", k),
			makeObjectHandler(server, k, v))
		http.HandleFunc(
			fmt.Sprintf("/%s", k),
			makeTypeHandler(server, k))
	}

	log.Println(http.ListenAndServe(":8000", nil))
	return server
}

func makeIdHandler(server *_Server) _HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := store.ObjectIdentity(mux.Vars(r)["id"])
		existing, _ := server.Store.Get(server.Context, id)
		var robject store.Object = nil
		data, err := utils.ReadStream(r.Body)
		if err != nil {
			robject, _ = utils.UnmarshalObject(data, server.Schema)
		}

		kind := ""
		if existing != nil {
			kind = existing.Metadata().Kind()
		} else if robject != nil {
			kind = robject.Metadata().Kind()
		}

		// method validation
		objMethods := server.Schema.ObjectMethods()[kind]
		if objMethods == nil || !slices.Contains(objMethods, r.Method) {
			reportError(w,
				fmt.Errorf("method not allowed"),
				http.StatusMethodNotAllowed)
		}

		server.handlePath(w, r, id, robject)
	}
}

func makeObjectHandler(server *_Server, t string, methods []string) _HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := store.ObjectIdentity(mux.Vars(r)["pkey"])
		var robject store.Object = nil
		data, err := utils.ReadStream(r.Body)
		if err == nil {
			robject, _ = utils.UnmarshalObject(data, server.Schema)
		}

		// method validation
		if !slices.Contains(methods, r.Method) {
			reportError(w,
				fmt.Errorf("method not allowed"),
				http.StatusMethodNotAllowed)
		}

		server.handlePath(w, r, id, robject)
	}
}

func makeTypeHandler(server *_Server, t string) _HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			opts := []store.ListOption{}

			vals := r.URL.Query()
			filter, ok := vals[FilterArg]
			if ok {
				flt := store.Filter{}
				err := json.Unmarshal([]byte(filter[0]), &flt)
				if err != nil {
					reportError(w, err, http.StatusBadRequest)
				}
				opts = append(opts, options.PropFilter(flt.Key, flt.Value))
			}

			pageSize, ok := vals[PageSizeArg]
			if ok {
				ps := 0
				err := json.Unmarshal([]byte(pageSize[0]), &ps)
				if err != nil {
					reportError(w, err, http.StatusBadRequest)
				}
				opts = append(opts, options.PageSize(ps))
			}

			pageOffset, ok := vals[PageOffsetArg]
			if ok {
				ps := 0
				err := json.Unmarshal([]byte(pageOffset[0]), &ps)
				if err != nil {
					reportError(w, err, http.StatusBadRequest)
				}
				opts = append(opts, options.PageOffset(ps))
			}

			orderBy, ok := vals[OrderByArg]
			if ok {
				ob := ""
				err := json.Unmarshal([]byte(orderBy[0]), &ob)
				if err != nil {
					reportError(w, err, http.StatusBadRequest)
				}
				opts = append(opts, options.OrderBy(ob))
			}

			orderInc, ok := vals[IncrementalArg]
			if ok {
				ob := false
				err := json.Unmarshal([]byte(orderInc[0]), &ob)
				if err != nil {
					reportError(w, err, http.StatusBadRequest)
				}
				opts = append(opts, options.OrderIncremental(ob))
			}

			ret, err := server.Store.List(
				server.Context,
				store.ObjectIdentity(t),
				opts...)

			if err != nil {
				reportError(w, err, http.StatusBadRequest)
			} else if ret != nil {
				resp, _ := json.Marshal(ret)
				w.Write(resp)
			}
		default:
			reportError(w,
				fmt.Errorf("method not allowed"),
				http.StatusMethodNotAllowed)
		}
	}
}

func (d *_Server) handlePath(
	w http.ResponseWriter,
	r *http.Request,
	identity store.ObjectIdentity,
	object store.Object) {

	var ret store.Object = nil
	var err error = nil
	switch r.Method {
	case http.MethodGet:
		ret, err = d.Store.Get(d.Context, identity)
		if err != nil {
			reportError(w, err, http.StatusNotFound)
		}
	case http.MethodPost:
		ret, err = d.Store.Create(d.Context, object)
		if err != nil {
			reportError(w, err, http.StatusNotAcceptable)
		}
	case http.MethodPut:
		ret, err = d.Store.Update(d.Context, identity, object)
		if err != nil {
			reportError(w, err, http.StatusNotAcceptable)
		}
	case http.MethodDelete:
		err = d.Store.Delete(d.Context, identity)
		if err != nil {
			reportError(w, err, http.StatusNotFound)
		}
	}

	if err == nil && ret != nil {
		resp, _ := json.Marshal(ret)
		w.Write(resp)
	}
}

func reportError(w http.ResponseWriter, err error, code int) {
	str := fmt.Sprintf("{ \"error\": \"%s\" }", err)
	http.Error(w, str, code)
}
