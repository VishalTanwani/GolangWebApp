package handler

import (
	"net/url"
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
	{"search-availability","/search-availability","POST",[]postData{
		{key:"start", value:"2020-02-02"},
		{key:"end", value:"2020-02-02"},
	},http.StatusOK},
	{"search-availability-json","/search-availability-json","POST",[]postData{
		{key:"start", value:"2020-02-02"},
		{key:"end", value:"2020-02-02"},
	},http.StatusOK},
	{"make-reservation","/make-reservation","POST",[]postData{
		{key:"first_name", value:"vishal"},
		{key:"last_name", value:"tanwani"},
		{key:"email", value:"vishal@gamil.com"},
		{key:"phone", value:"444-4444-44"},
	},http.StatusOK},
}

func TestHandlers(t *testing.T){
	routes := getRoutes()
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()
	for _,e := range theTests {
		if e.method == "GET"{
			resp,err := ts.Client().Get(ts.URL+e.url)
			if err !=nil {
				t.Log(err)
				t.Fatal(err)
			}
			if resp.StatusCode!=e.expectedStatusCode{
				t.Errorf("for %s expected %d but got %d",e.name,e.expectedStatusCode,resp.StatusCode)
			}
		} else {
			values := url.Values{}
			for _,x:=range e.params{
				values.Add(x.key,x.value)
			}
			resp,err := ts.Client().PostForm(ts.URL+e.url,values)
			if err !=nil {
				t.Log(err)
				t.Fatal(err)
			}
			if resp.StatusCode!=e.expectedStatusCode{
				t.Errorf("for %s expected %d but got %d",e.name,e.expectedStatusCode,resp.StatusCode)
			}
		}
		
	}
}