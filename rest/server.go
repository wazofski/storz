package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"golang.org/x/exp/slices"

	"github.com/wazofski/store"
	"github.com/wazofski/store/options"
	"github.com/wazofski/store/utils"
)

const (
	PropFilterArg  = "pf"
	KeyFilterArg   = "kf"
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
		if err == nil {
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
			filter, ok := vals[PropFilterArg]
			if ok {
				flt := store.PropFilter{}
				err := json.Unmarshal([]byte(filter[0]), &flt)
				if err != nil {
					reportError(w, err, http.StatusBadRequest)
				}
				opts = append(opts, options.PropFilter(flt.Key, flt.Value))
			}

			keyFilter, ok := vals[KeyFilterArg]
			if ok {
				flt := store.KeyFilter{}
				err := json.Unmarshal([]byte(keyFilter[0]), &flt)
				if err != nil {
					reportError(w, err, http.StatusBadRequest)
				}
				opts = append(opts, options.KeyFilter(flt...))
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
				writeResponse(w, resp)
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

			robject, err := utils.UnmarshalObject(data, server.Schema, t)
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
		writeResponse(w, resp)
	}
}

var writtenResp map[*http.ResponseWriter]bool = make(map[*http.ResponseWriter]bool)

func reportError(w http.ResponseWriter, err error, code int) {
	// log.Panicf(err.Error())
	if !isWritten(w) {
		http.Error(w, err.Error(), code)
	}
}

func isWritten(w http.ResponseWriter) bool {
	ret, ok := writtenResp[&w]
	if !ok {
		return false
	}
	return ret
}

func writeResponse(w http.ResponseWriter, data []byte) {
	if !isWritten(w) {
		w.Write(data)
		writtenResp[&w] = true
	}
}

func prepResponse(w http.ResponseWriter, r *http.Request) {
	log.Printf("SERVER %s %s", strings.ToLower(r.Method), r.URL)

	if !isWritten(w) {
		w.Header().Add("Content-Type", "application/json")
	}
}
