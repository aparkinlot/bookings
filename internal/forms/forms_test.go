package forms

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestValidForm(t *testing.T) {
	r := httptest.NewRequest("POST", "/blank", nil)
	form := New(r.PostForm)

	isValid := form.Valid()
	if !isValid {
		t.Error("Form is invalid when it should have been valid")
	}
}

func TestRequiredForm(t *testing.T) {
	r := httptest.NewRequest("POST", "/blank", nil)
	form := New(r.PostForm)

	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("Form is valid when required fields are missing")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "a")
	postedData.Add("c", "a")

	r, _ = http.NewRequest("POST", "/blank", nil)

	r.PostForm = postedData
	form = New(r.PostForm)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("Form doesn't have fields when it does")
	}
}

func TestHasForm(t *testing.T) {
	postedData := url.Values{}
	form := New(url.Values{})

	has := form.Has("something")
	if has {
		t.Error("Form shows has field param when it doesn't")
	}

	postedData = url.Values{}
	postedData.Add("a", "a")

	form = New(postedData)
	has = form.Has("a")
	if !has {
		t.Error("Form shows no field params when it should")
	}
}

func TestMinLengthForm(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)

	form.MinLength("x", 10)
	if form.Valid() {
		t.Error("Form shows min lenth err for blank field")
	}

	isErr := form.Errors.Get("x")
	if isErr == "" {
		t.Error("Should have err but does not")
	}

	postedData = url.Values{}
	postedData.Add("some_field", "some_value")
	form = New(postedData)

	form.MinLength("some_field", 100)
	if form.Valid() {
		t.Error("Form has a field which the min length of 100 is not met")
	}

	postedData = url.Values{}
	postedData.Add("another_field", "another_value")
	form = New(postedData)

	form.MinLength("another_field", 1)
	if !form.Valid() {
		t.Error("Form has field larger than min length of 1")
	}

	isErr = form.Errors.Get("another_field")
	if isErr != "" {
		t.Error("Got an err when not expected")
	}

}

func TestValidEmailForm(t *testing.T) {
	postedValues := url.Values{}
	form := New(postedValues)

	form.IsEmail("x")
	if form.Valid() {
		t.Error("form shows valid email for non-existent field")
	}

	postedValues = url.Values{}
	postedValues.Add("email", "me@here.com")
	form = New(postedValues)

	form.IsEmail("email")
	if !form.Valid() {
		t.Error("got an invalid email when we should not have")
	}

	postedValues = url.Values{}
	postedValues.Add("email", "x")
	form = New(postedValues)

	form.IsEmail("email")
	if form.Valid() {
		t.Error("got valid for invalid email address")
	}
}
