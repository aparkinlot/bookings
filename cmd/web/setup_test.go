package main

import (
	"net/http"
	"os"
	"testing"
)

func TestMain(m *testing.M) {

	// After everything is done, now run the tests
	os.Exit(m.Run())
}

type myHandler struct {
}

func (mh *myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
