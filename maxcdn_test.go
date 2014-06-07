package maxcdn

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

	. "github.com/jmervine/GoT"
)

//var (
//alias  = os.Getenv("ALIAS")
//token  = os.Getenv("TOKEN")
//secret = os.Getenv("SECRET")
//)

func Test(T *testing.T) {
	max := NewMaxCDN("alias", "token", "secret")
	Go(T).AssertEqual(max.Alias, "alias")
	Go(T).AssertEqual(max.client.Credentials.Token, "token")
	Go(T).AssertEqual(max.client.Credentials.Secret, "secret")
}

func TestMaxCDN_Get(T *testing.T) {
	max := NewMaxCDN("alias", "token", "secret")

	var recorder http.Response
	max.HTTPClient = stubHTTPOkRecorded(&recorder)

	payload, err := max.Get("/account.json", nil)
	Go(T).AssertNil(err)
	Go(T).RefuteNil(payload)
	Go(T).AssertNil(payload.Error)
	Go(T).RefuteNil(payload.Raw)
	Go(T).RefuteNil(payload.Body)

	Go(T).AssertEqual(recorder.Request.Method, "GET")
	Go(T).AssertEqual(recorder.Request.URL.Path, "/alias/account.json")
	Go(T).AssertEqual(recorder.Request.URL.Query().Encode(), "")
	Go(T).AssertEqual(recorder.Request.Header.Get("Content-Type"), contentType)
	Go(T).RefuteEqual(recorder.Request.Header.Get("Authorization"), "")

	// check body
	Go(T).AssertNil(recorder.Request.Body)
}

func TestMaxCDN_Put(T *testing.T) {
	max := NewMaxCDN("alias", "token", "secret")

	var recorder http.Response
	max.HTTPClient = stubHTTPOkRecorded(&recorder)

	form := url.Values{}
	form.Add("name", "foo")

	payload, err := max.Put("/account.json", form)
	Go(T).AssertNil(err)
	Go(T).RefuteNil(payload)
	Go(T).AssertNil(payload.Error)
	Go(T).RefuteNil(payload.Raw)
	Go(T).RefuteNil(payload.Body)

	Go(T).AssertEqual(recorder.Request.Method, "PUT")
	Go(T).AssertEqual(recorder.Request.URL.Path, "/alias/account.json")
	Go(T).AssertEqual(recorder.Request.URL.Query().Encode(), "")
	Go(T).AssertEqual(recorder.Request.Header.Get("Content-Type"), contentType)
	Go(T).RefuteEqual(recorder.Request.Header.Get("Authorization"), "")

	// check body
	body, err := ioutil.ReadAll(recorder.Request.Body)
	Go(T).AssertNil(err)
	Go(T).AssertEqual(string(body), "name=foo")
}

func TestMaxCDN_Post(T *testing.T) {
	max := NewMaxCDN("alias", "token", "secret")

	var recorder http.Response
	max.HTTPClient = stubHTTPOkRecorded(&recorder)

	form := url.Values{}
	form.Add("name", "foo")

	payload, err := max.Post("/zones/pull.json", form)
	Go(T).AssertNil(err)
	Go(T).RefuteNil(payload)
	Go(T).AssertNil(payload.Error)
	Go(T).RefuteNil(payload.Raw)
	Go(T).RefuteNil(payload.Body)

	Go(T).AssertEqual(recorder.Request.Method, "POST")
	Go(T).AssertEqual(recorder.Request.URL.Path, "/alias/zones/pull.json")
	Go(T).AssertEqual(recorder.Request.URL.Query().Encode(), "")
	Go(T).AssertEqual(recorder.Request.Header.Get("Content-Type"), contentType)
	Go(T).RefuteEqual(recorder.Request.Header.Get("Authorization"), "")

	// check body
	body, err := ioutil.ReadAll(recorder.Request.Body)
	Go(T).AssertNil(err)
	Go(T).AssertEqual(string(body), "name=foo")
}

func TestMaxCDN_Delete(T *testing.T) {
	max := NewMaxCDN("alias", "token", "secret")

	var recorder http.Response
	max.HTTPClient = stubHTTPOkRecorded(&recorder)

	payload, err := max.Delete("/zones/pull.json/123456", nil)
	Go(T).AssertNil(err)
	Go(T).RefuteNil(payload)
	Go(T).AssertNil(payload.Error)
	Go(T).RefuteNil(payload.Raw)
	Go(T).RefuteNil(payload.Body)

	Go(T).AssertEqual(recorder.Request.Method, "DELETE")
	Go(T).AssertEqual(recorder.Request.URL.Path, "/alias/zones/pull.json/123456")
	Go(T).AssertEqual(recorder.Request.URL.Query().Encode(), "")
	Go(T).AssertEqual(recorder.Request.Header.Get("Content-Type"), contentType)
	Go(T).RefuteEqual(recorder.Request.Header.Get("Authorization"), "")

	// check body
	Go(T).AssertNil(recorder.Request.Body)
}

func TestMaxCDN_PurgeZone(T *testing.T) {
	max := NewMaxCDN("alias", "token", "secret")

	var recorder http.Response
	max.HTTPClient = stubHTTPOkRecorded(&recorder)

	payload, err := max.PurgeZone(123456)
	Go(T).AssertNil(err)
	Go(T).RefuteNil(payload)
	Go(T).AssertNil(payload.Error)
	Go(T).RefuteNil(payload.Raw)
	Go(T).RefuteNil(payload.Body)

	Go(T).AssertEqual(recorder.Request.Method, "DELETE")
	Go(T).AssertEqual(recorder.Request.URL.Path, "/alias/zones/pull.json/123456/cache")
	Go(T).AssertEqual(recorder.Request.URL.Query().Encode(), "")
	Go(T).AssertEqual(recorder.Request.Header.Get("Content-Type"), contentType)
	Go(T).RefuteEqual(recorder.Request.Header.Get("Authorization"), "")

	// check body
	Go(T).AssertNil(recorder.Request.Body)
}

func TestMaxCDN_PurgeZones(T *testing.T) {
	max := NewMaxCDN("alias", "token", "secret")

	var recorder http.Response
	max.HTTPClient = stubHTTPOkRecorded(&recorder)

	payload, err := max.PurgeZones([]int{12345, 23456, 34567})
	Go(T).AssertNil(err)
	Go(T).RefuteNil(payload)

	Go(T).AssertEqual(recorder.Request.Method, "DELETE")
	Go(T).AssertEqual(recorder.Request.URL.Query().Encode(), "")
	Go(T).AssertEqual(recorder.Request.Header.Get("Content-Type"), contentType)
	Go(T).RefuteEqual(recorder.Request.Header.Get("Authorization"), "")

	// check body
	Go(T).AssertNil(recorder.Request.Body)
}

func TestMaxCDN_PurgeFile(T *testing.T) {
	max := NewMaxCDN("alias", "token", "secret")

	var recorder http.Response
	max.HTTPClient = stubHTTPOkRecorded(&recorder)

	payload, err := max.PurgeFile(123456, "/master.css")
	Go(T).AssertNil(err)
	Go(T).RefuteNil(payload)
	Go(T).AssertNil(payload.Error)
	Go(T).RefuteNil(payload.Raw)
	Go(T).RefuteNil(payload.Body)

	Go(T).AssertEqual(recorder.Request.Method, "DELETE")
	Go(T).AssertEqual(recorder.Request.URL.Path, "/alias/zones/pull.json/123456/cache")
	Go(T).AssertEqual(recorder.Request.URL.Query().Encode(), "")
	Go(T).AssertEqual(recorder.Request.Header.Get("Content-Type"), contentType)
	Go(T).RefuteEqual(recorder.Request.Header.Get("Authorization"), "")

	// check body
	body, err := ioutil.ReadAll(recorder.Request.Body)
	Go(T).AssertNil(err)
	Go(T).AssertEqual(string(body), "file=%2Fmaster.css")
}

func TestMaxCDN_PurgeFiles(T *testing.T) {
	max := NewMaxCDN("alias", "token", "secret")

	var recorder http.Response
	max.HTTPClient = stubHTTPOkRecorded(&recorder)

	files := []string{"/master.css", "/master.js", "/index.html"}
	payload, err := max.PurgeFiles(123456, files)
	Go(T).AssertNil(err)
	Go(T).RefuteNil(payload)

	Go(T).AssertEqual(recorder.Request.Method, "DELETE")
	Go(T).AssertEqual(recorder.Request.URL.Query().Encode(), "")
	Go(T).AssertEqual(recorder.Request.Header.Get("Content-Type"), contentType)
	Go(T).RefuteEqual(recorder.Request.Header.Get("Authorization"), "")

	// check body
	Go(T).RefuteNil(recorder.Request.Body)
	Go(T).AssertNil(err)
}
