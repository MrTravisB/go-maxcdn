// Package maxcdn is the golang bindings for MaxCDN's REST API.
//
// At this time it should be considered very alpha.
package maxcdn

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/garyburd/go-oauth/oauth"
)

const (
	ApiPath     = "https://rws.netdna.com"
	UserAgent   = "Go MaxCDN API Client"
	ContentType = "application/x-www-form-urlencoded"
)

// MaxCDN is the core struct for interacting with MaxCDN.
//
// HttpClient can be overridden as needed, but will be set to
// http.DefaultClient by default.
type MaxCDN struct {
	Alias      string
	client     oauth.Client
	HttpClient *http.Client
}

// NewMaxCDN sets up a new MaxCDN instance.
func NewMaxCDN(alias, token, secret string) *MaxCDN {
	return &MaxCDN{
		HttpClient: http.DefaultClient,
		Alias:      alias,
		client: oauth.Client{
			Credentials: oauth.Credentials{
				Token:  token,
				Secret: secret,
			},
			TemporaryCredentialRequestURI: ApiPath + "oauth/request_token",
			TokenRequestURI:               ApiPath + "oauth/access_token",
		},
	}
}

func (max *MaxCDN) Get(endpoint string, form url.Values) (*GenericResponse, error) {
	return max.do("GET", endpoint, form)
}

func (max *MaxCDN) Post(endpoint string, form url.Values) (*GenericResponse, error) {
	return max.do("POST", endpoint, form)
}

func (max *MaxCDN) Put(endpoint string, form url.Values) (*GenericResponse, error) {
	return max.do("PUT", endpoint, form)
}

func (max *MaxCDN) Delete(endpoint string) (*GenericResponse, error) {
	return max.do("DELETE", endpoint, nil)
}

// PurgeZone purges a specified zones cache.
func (max *MaxCDN) PurgeZone(zone int) (*GenericResponse, error) {
	return max.Delete(fmt.Sprintf("/zones/pull.json/%d/cache", zone))
}

// PurgeZone purges a multiple zones caches.
func (max *MaxCDN) PurgeZones(zones []int) (responses []GenericResponse, last error) {
	var rc chan *GenericResponse
	var ec chan error

	waiter := sync.WaitGroup{}
	mutex := sync.Mutex{}

	done := func() {
		waiter.Done()
	}

	send := func(zone int) {
		defer done()
		r, e := max.PurgeZone(zone)

		rc <- r
		ec <- e
	}

	collect := func() {
		defer done()
		r := <-rc
		e := <-ec

		mutex.Lock()
		responses = append(responses, *r)
		last = e
		mutex.Unlock()
	}

	for _, zone := range zones {
		waiter.Add(2)
		go send(zone)
		go collect()
	}

	waiter.Wait()
	return
}

// TODO:
//
// func (max *MaxCDN) PurgeFiles(files []string) (responses []Genericresponse, last error)
//
// max.Delete will have to support form params first

func (max *MaxCDN) url(endpoint string) string {
	endpoint = strings.TrimPrefix(endpoint, "/")
	return fmt.Sprintf("%s/%s/%s", ApiPath, max.Alias, endpoint)
}

func (max *MaxCDN) do(method, endpoint string, form url.Values) (response *GenericResponse, err error) {
	var req *http.Request

	req, err = http.NewRequest(method, max.url(endpoint), nil)
	if err != nil {
		return
	}

	if method == "GET" && req.URL.RawQuery != "" {
		return nil, errors.New("oauth: url must not contain a query string")
	}

	if form != nil {
		if method == "GET" {
			req.URL.RawQuery = form.Encode()
		} else {
			req.Body = ioutil.NopCloser(strings.NewReader(form.Encode()))
		}

		// Only post needs a signed form.
		if method != "POST" {
			form = nil
		}
	}

	req.Header.Set("Authorization", max.client.AuthorizationHeader(nil, method, req.URL, form))
	req.Header.Set("Content-Type", ContentType)
	req.Header.Set("User-Agent", UserAgent)

	resp, err := max.HttpClient.Do(req)
	defer resp.Body.Close()

	if err != nil {
		return
	}

	parser := &genericParser{}
	r, e := parser.parse(resp)
	rr := r.(GenericResponse)

	if e == nil && (rr.Error.Message != "" || rr.Error.Type != "") {
		e = errors.New(fmt.Sprintf("%s, %s", rr.Error.Type, rr.Error.Message))
	}

	return &rr, e
}
