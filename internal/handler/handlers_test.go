package handler

import (
	"net/http/httptest"
	"testing"
	"net/http"
)

type postData struct {
	key string
	value string
}

var theTests = []struct {
	name string
	url string
	method string
	params []postData
	expectedStatusCode int
}{
	{"home","/","GET",[]postData{},http.StatusOK},
	{"about","/about","GET",[]postData{},http.StatusOK},
	{"general-quarters","/general-quarters","GET",[]postData{},http.StatusOK},
	{"major-suite","/major-suite","GET",[]postData{},http.StatusOK},
	{"search-availability","/search-availability","GET",[]postData{},http.StatusOK},
	{"make-reservation","/make-reservation","GET",[]postData{},http.StatusOK},
	{"contact","/contact","GET",[]postData{},http.StatusOK},
}

func TestHandlers(t *testing.T){
	routes := getRoutes()
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()
	for _,e := range theTests {
		resp,err := ts.Client().Get(ts.URL+e.url)
		if err !=nil {
			t.Log(err)
			t.Fatal(err)
		}
		if resp.StatusCode!=e.expectedStatusCode{
			t.Errorf("for %s expected %d but got %d",e.name,e.expectedStatusCode,resp.StatusCode)
		}
	}
}