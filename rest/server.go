package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
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
	Listen(int)
}

type _HandlerFunc func(http.ResponseWriter, *http.Request)

type _Server struct {
	Schema  store.SchemaHolder
	Store   store.Store
	Context context.Context
	Router  *mux.Router
}

func (d *_Server) Listen(port int) {
	log.Println(http.ListenAndServe(
		fmt.Sprintf(":%d", port), d.Router))
}

func Server(schema store.SchemaHolder, store store.Store) Endpoint {
	server := &_Server{
		Schema:  schema,
		Store:   store,
		Context: context.Background(),
		Router:  mux.NewRouter(),
	}

	addHandler(server.Router, "/id/{id}", makeIdHandler(server))
	for k, v := range schema.ObjectMethods() {
		addHandler(server.Router,
			fmt.Sprintf("/%s/{pkey}", strings.ToLower(k)),
			makeObjectHandler(server, k, v))
		addHandler(server.Router,
			fmt.Sprintf("/%s", strings.ToLower(k)),
			makeTypeHandler(server, k, v))
		addHandler(server.Router,
			fmt.Sprintf("/%s/", strings.ToLower(k)),
			makeTypeHandler(server, k, v))
	}

	return server
}

func addHandler(router *mux.Router, pattern string, handler _HandlerFunc) {
	log.Printf("SERVER serving %s", pattern)
	router.HandleFunc(pattern, handler)
}

func makeIdHandler(server *_Server) _HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prepResponse(w, r)
		id := store.ObjectIdentity(mux.Vars(r)["id"])
		existing, _ := server.Store.Get(server.Context, id)
		var robject store.Object = nil
		data, err := utils.ReadStream(r.Body)
		if err != nil {
			robject, _ = utils.UnmarshalObject(data, server.Schema, utils.ObjeectKind(data))
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
		prepResponse(w, r)
		var robject store.Object = nil
		id := store.ObjectIdentity(strings.ToLower(t) + "/" + mux.Vars(r)["pkey"])
		data, err := utils.ReadStream(r.Body)
		if err == nil {
			robject, _ = utils.UnmarshalObject(data, server.Schema, t)
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

func makeTypeHandler(server *_Server, t string, methods []string) _HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prepResponse(w, r)
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
				ps, _ := strconv.Atoi(pageSize[0])
				opts = append(opts, options.PageSize(ps))
			}

			pageOffset, ok := vals[PageOffsetArg]
			if ok {
				ps, _ := strconv.Atoi(pageOffset[0])
				opts = append(opts, options.PageOffset(ps))
			}

			orderBy, ok := vals[OrderByArg]
			if ok {
				ob := orderBy[0]
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
				store.ObjectIdentity(
					fmt.Sprintf("%s/", strings.ToLower(t))),
				opts...)

			// log.Printf("size of list %d", len(ret))

			if err != nil {
				reportError(w, err, http.StatusBadRequest)
			} else if ret != nil {
				resp, _ := json.Marshal(ret)
				w.Write(resp)
			}
		case http.MethodPost:
			// method validation
			if !slices.Contains(methods, r.Method) {
				reportError(w,
					fmt.Errorf("method not allowed"),
					http.StatusMethodNotAllowed)
				return
			}

			data, err := utils.ReadStream(r.Body)
			if err != nil {
				reportError(w,
					err,
					http.StatusBadRequest)
				return
			}
			robject, err := utils.UnmarshalObject(data, server.Schema, utils.ObjeectKind(data))
			if err != nil {
				reportError(w,
					err,
					http.StatusBadRequest)
				return
			}

			server.handlePath(w, r, store.ObjectIdentity(""), robject)
		default:
			reportError(w,
				fmt.Errorf("method not allowed"),
				http.StatusMethodNotAllowed)
		}
	}
}

func prepResponse(w http.ResponseWriter, r *http.Request) {
	log.Printf("SERVER %s %s", strings.ToLower(r.Method), r.URL)

	w.Header().Add("Content-Type", "application/json")
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
		object.Metadata().SetIdentity(store.ObjectIdentity(uuid.New().String()))
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
	// log.Panicf(err.Error())
	http.Error(w, err.Error(), code)
}
