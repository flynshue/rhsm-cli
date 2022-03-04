package rhsm

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
)

func TestClient_ProcessRequest(t *testing.T) {
	client := NewClient(os.Getenv("RHSM_TOKEN"))
	router := NewRouter()
	router.AddFunc(200, func(resp *http.Response) error {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		fmt.Println(string(b))
		return nil
	})
	resource := NewRestResource("GET", "/systems?limit=10", router)
	err := client.ProcessRequest("https://api.access.redhat.com/management/v1", resource, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
}
