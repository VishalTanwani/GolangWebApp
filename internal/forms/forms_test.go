package forms

import (
	"net/url"
	"net/http/httptest"
	"testing"
)

func TestValid(t *testing.T){
	r := httptest.NewRequest("POST","/whatever",nil)
	form := New(r.PostForm)
	if !form.Valid() {
		t.Error("got invalid form")
	}
}

func TestRequired(t *testing.T){
	r := httptest.NewRequest("POST","/whatever",nil)
	form := New(r.PostForm)

	form.Required("a","b","c")
	if form.Valid() {
		t.Error("got valid form where it should be invalid")
	}

	postData := url.Values{}

	postData.Add("a","a")
	postData.Add("b","b")
	postData.Add("c","c")
	
	r = httptest.NewRequest("POST","/whatever",nil)

	r.PostForm = postData
	form = New(r.PostForm)

	if !form.Valid() {
		t.Error("got invalid form")
	}
}

func TestHas(t *testing.T){
	r := httptest.NewRequest("POST","/whatever",nil)
	form := New(r.PostForm)

	if form.Has("a") {
		t.Error("it has a value but it should not be")
	}

	postData := url.Values{}
	postData.Add("a","a")

	form = New(postData)

	if !form.Has("a") {
		t.Error("does not have a value filed when it should")
	}
}

func TestIsEmail(t *testing.T){
	r := httptest.NewRequest("POST","/whatever",nil)
	form := New(r.PostForm)
	form.IsEmail("email")

	if form.Valid() {
		t.Error("got valid form where it should be invalid because it does not have a email")
	}

	postData := url.Values{}
	postData.Add("email","vishal@gmail.com")

	form = New(postData)

	form.IsEmail("email")

	if !form.Valid() {
		t.Error("got invalid email where is should be valid")
	}

	postData = url.Values{}
	postData.Add("email","vishalgmail.com")

	form = New(postData)
	
	form.IsEmail("email")

	if form.Valid() {
		t.Error("got valid email where is should be invalid")
	}
}

func TestMinLength(t *testing.T){
	r := httptest.NewRequest("POST","/whatever",nil)
	form := New(r.PostForm)

	form.MinLength("first_name",3)

	if form.Valid() {
		t.Error("it shows length of non-exist filed")
	}

	isError := form.Errors.Get("first_name")
	if isError == "" {
		t.Error("it shows no error but it should be")
	}

	postData := url.Values{}
	postData.Add("first_name","vishal")

	form = New(postData)

	form.MinLength("first_name",100)

	if form.Valid(){
		t.Error("shows minlength of 100 met when data is shorter")
	}

	postData = url.Values{}
	postData.Add("first_name","vishal")

	form = New(postData)

	isError = form.Errors.Get("first_name")
	if isError != "" {
		t.Error("it shows error but it should not be")
	}
	form.MinLength("first_name",1)

	if !form.Valid() {
		t.Error("shows minlength of 1 not met when data is correct")
	}
}