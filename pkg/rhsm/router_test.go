package rhsm

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"
)

func TestRouter_Call(t *testing.T) {
	router := NewRouter()
	router.AddFunc(200, func(resp *http.Response) error {
		fmt.Printf("%d from %s", resp.StatusCode, resp.Request.URL.Path)
		return nil
	})
	resp := &http.Response{
		StatusCode: 200,
		Request:    &http.Request{URL: &url.URL{Path: "/subscriptions"}},
	}
	if err := router.Call(resp); err != nil {
		t.Fatal(err)
	}
}
