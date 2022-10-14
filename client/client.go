package client

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/wazofski/store"
	"github.com/wazofski/store/rest"
	"github.com/wazofski/store/utils"
)

type restStore struct {
	BaseURL     *url.URL
	Schema      store.SchemaHolder
	MakeRequest requestMaker
	Headers     []headerOption
}

type requestMaker func(path *url.URL, content []byte, method string, headers map[string]string) ([]byte, error)

type restOptions struct {
	store.CommonOptionHolder
	Headers map[string]string
}

func newRestOptions(d *restStore) restOptions {
	res := restOptions{
		CommonOptionHolder: store.CommonOptionHolder{},
		Headers:            make(map[string]string),
	}

	for _, h := range d.Headers {
		h.ApplyFunction()(&res)
	}

	return res
}

func (d *restOptions) CommonOptions() *store.CommonOptionHolder {
	return &d.CommonOptionHolder
}

func Factory(serviceUrl string, headers ...headerOption) store.Factory {
	return func(schema store.SchemaHolder) (store.Store, error) {
		URL, err := url.Parse(serviceUrl)
		if err != nil {
			return nil, fmt.Errorf("%s; expected format: http(s)://address:port/argo/api", err)
		}

		client := &restStore{
			BaseURL:     URL,
			Schema:      schema,
			MakeRequest: makeHttpRequest,
			Headers:     headers,
		}

		log.Printf("REST client initialized: %s", serviceUrl)
		return client, nil
	}
}

func makeHttpRequest(path *url.URL, content []byte, requestType string, headers map[string]string) ([]byte, error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	req, err := http.NewRequest(requestType, path.String(), strings.NewReader(string(content)))
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	req.Close = true

	// req.ContentLength = contentLength
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	rd, err := utils.ReadStream(resp.Body)

	defer resp.Body.Close()

	if err != nil {
		return rd, err
	}

	if resp.StatusCode < 200 || 300 <= resp.StatusCode {
		return rd, fmt.Errorf("http %d", resp.StatusCode)
	}

	return rd, nil
}

func processRequest(
	client *restStore,
	requestUrl *url.URL,
	content []byte,
	method string,
	headers map[string]string) ([]byte, error) {

	reqId := uuid.New().String()
	requestUrl.Path = strings.ReplaceAll(requestUrl.Path, "//", "/")
	origin := strings.ReplaceAll(requestUrl.String(), requestUrl.Path, "")
	headers["Origin"] = strings.ReplaceAll(origin, requestUrl.RawQuery, "")
	headers["X-Request-ID"] = reqId
	headers["Content-Type"] = "application/json"
	headers["X-Requested-With"] = "XMLHttpRequest"

	log.Printf("CLIENT %s %s", strings.ToLower(method), requestUrl)
	// log.Printf("X-Request-ID %s", reqId)

	data, err := client.MakeRequest(requestUrl, content, method, headers)
	cerr := errorCheck(data)
	if err == nil {
		err = cerr
	} else if cerr != nil {
		err = fmt.Errorf("%s %s", err, cerr)
	}

	if err != nil {
		// log.Println(err)
		if len(content) > 0 {
			var js interface{}
			if json.Unmarshal([]byte(content), &js) == nil {
				r, _ := json.MarshalIndent(js, "", "    ")
				log.Printf("CLIENT request content: %s", r)
			} else {
				log.Printf("CLIENT request content: %s", content)
			}
		}
		if len(data) > 0 {
			log.Printf("CLIENT response content: %s", string(data))
		}
		return nil, err
	}

	// mol := store.ObjectList{}
	// if err := json.Unmarshal(data, &mol); err != nil {
	// 	// ignore errors
	// 	// log.Printf("Unable to Unmarshal")
	// 	err = nil
	// }

	// if mol.Items != nil {
	// 	// This is a response of a GET on a collection, which is a list.
	// 	data = *mol.Items
	// }

	return data, err
}

func errorCheck(response []byte) error {
	str := string(response)
	if len(str) == 0 {
		return nil
	}

	um := make(map[string]interface{})

	err := json.Unmarshal(response, &um)
	if err == nil {
		if v, found := um["errors"]; found {
			return errors.New(v.([]interface{})[0].(string))
		}
		if v, found := um["error"]; found {
			m := v.(map[string]interface{})

			return fmt.Errorf("%v %s",
				m["status_code"],
				m["status"])
		}
	}

	return nil
}

func serialize(mo store.Object) ([]byte, error) {
	if mo == nil {
		return nil, errors.New("cannot serialize nil store.Object")
	}

	return json.Marshal(mo)
}

func makePathForType(baseUrl *url.URL, obj store.Object) *url.URL {
	u, _ := url.Parse(fmt.Sprintf("%s/%s", baseUrl, strings.ToLower(obj.Metadata().Kind())))
	return u
}

func removeTrailingSlash(val string) string {
	if strings.HasSuffix(val, "/") {
		return val[:len(val)-1]
	}
	return val
}

func makePathForIdentity(baseUrl *url.URL, identity store.ObjectIdentity, params string) *url.URL {
	if len(params) > 0 {
		path := fmt.Sprintf("%s/%s?%s",
			baseUrl,
			removeTrailingSlash(identity.Path()),
			params)

		// log.Printf(`made path %s # %s # %s`, path, identity.Path(), string(identity))

		u, _ := url.ParseRequestURI(path)
		return u
	}

	u, _ := url.Parse(fmt.Sprintf("%s/%s", baseUrl, identity.Path()))
	return u
}

func toBytes(obj interface{}) []byte {
	if obj == nil {
		return []byte{}
	}

	jsn, _ := json.Marshal(obj)

	return []byte(string(jsn))
}

func listParameters(ropt restOptions) string {
	opt := ropt.CommonOptions()

	q := url.Values{}
	if len(opt.OrderBy) > 0 {
		q.Add(rest.OrderByArg, opt.OrderBy)
		q.Add(rest.IncrementalArg, strconv.FormatBool(opt.OrderIncremental))
	}

	if opt.PageOffset > 0 {
		q.Add(rest.PageOffsetArg, fmt.Sprintf("%d", opt.PageOffset))
	}

	if opt.PageSize > 0 {
		q.Add(rest.PageSizeArg, fmt.Sprintf("%d", opt.PageSize))
	}

	if opt.PropFilter != nil {
		content, err := json.Marshal(opt.PropFilter)
		if err != nil {
			log.Fatalln(err)
		}

		if len(content) > 0 {
			q.Add(rest.PropFilterArg, string(content))
		}
	}

	if opt.KeyFilter != nil {
		content, err := json.Marshal(opt.KeyFilter)
		if err != nil {
			log.Fatalln(err)
		}

		if len(content) > 0 {
			q.Add(rest.KeyFilterArg, string(content))
		}
	}

	return q.Encode()
}

func (d *restStore) Create(
	ctx context.Context,
	obj store.Object,
	opt ...store.CreateOption) (store.Object, error) {

	copt := newRestOptions(d)
	var err error
	for _, o := range opt {
		err = o.ApplyFunction()(&copt)
		if err != nil {
			return nil, err
		}
	}

	data, err := serialize(obj)
	if err != nil {
		return nil, err
	}

	data, err = processRequest(d,
		makePathForType(d.BaseURL, obj),
		data,
		http.MethodPost,
		copt.Headers)

	if err != nil {
		return nil, err
	}

	clone := obj.Clone()
	err = json.Unmarshal(data, &clone)
	if err != nil {
		log.Println(string(data))
		clone = nil
	}

	return clone, err
}

func (d *restStore) Update(
	ctx context.Context,
	identity store.ObjectIdentity,
	obj store.Object,
	opt ...store.UpdateOption) (store.Object, error) {

	copt := newRestOptions(d)
	var err error
	for _, o := range opt {
		err = o.ApplyFunction()(&copt)
		if err != nil {
			return nil, err
		}
	}

	data, err := serialize(obj)
	if err != nil {
		return nil, err
	}

	data, err = processRequest(d,
		makePathForIdentity(d.BaseURL, identity, ""),
		data,
		http.MethodPut,
		copt.Headers)

	if err != nil {
		return nil, err
	}

	clone := obj.Clone()
	err = json.Unmarshal(data, &clone)

	return clone, err
}

func (d *restStore) Delete(
	ctx context.Context,
	identity store.ObjectIdentity,
	opt ...store.DeleteOption) error {

	var err error
	copt := newRestOptions(d)
	for _, o := range opt {
		err = o.ApplyFunction()(&copt)
		if err != nil {
			return err
		}
	}

	_, err = processRequest(d,
		makePathForIdentity(d.BaseURL, identity, ""),
		[]byte{},
		http.MethodDelete,
		copt.Headers)

	return err
}

func (d *restStore) Get(
	ctx context.Context,
	identity store.ObjectIdentity,
	opt ...store.GetOption) (store.Object, error) {

	var err error
	copt := newRestOptions(d)
	for _, o := range opt {
		err = o.ApplyFunction()(&copt)
		if err != nil {
			return nil, err
		}
	}

	resp, err := processRequest(d,
		makePathForIdentity(d.BaseURL, identity, ""),
		[]byte{},
		http.MethodGet,
		copt.Headers)

	if err != nil {
		return nil, err
	}

	tp := identity.Type()
	if tp == "id" {
		tp = utils.ObjeectKind(resp)
	}

	return utils.UnmarshalObject(resp, d.Schema, tp)
}

func (d *restStore) List(
	ctx context.Context,
	identity store.ObjectIdentity,
	opt ...store.ListOption) (store.ObjectList, error) {

	var err error
	copt := newRestOptions(d)
	for _, o := range opt {
		err = o.ApplyFunction()(&copt)
		if err != nil {
			return nil, err
		}
	}

	params := listParameters(copt)
	path := makePathForIdentity(d.BaseURL, identity, params)
	res, err := processRequest(
		d,
		path,
		[]byte{},
		http.MethodGet,
		copt.Headers)

	if err != nil {
		return nil, err
	}

	parsed := []*json.RawMessage{}
	err = json.Unmarshal(res, &parsed)
	if err != nil {
		return nil, err
	}

	marshalledResult := store.ObjectList{}
	if len(parsed) == 0 {
		return marshalledResult, nil
	}

	resource := d.Schema.ObjectForKind(utils.ObjeectKind(*parsed[0]))

	for _, r := range parsed {
		clone := resource.Clone()
		clone.UnmarshalJSON(toBytes(r))

		marshalledResult = append(marshalledResult, clone)
	}

	return marshalledResult, nil
}
