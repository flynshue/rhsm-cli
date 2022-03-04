package rhsm

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
)

func TestAPI_Call(t *testing.T) {
	api := NewAPI("https://api.access.redhat.com/management/v1", os.Getenv("RHSM_TOKEN"))
	router := NewRouter()
	router.AddFunc(200, func(resp *http.Response) error {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		fmt.Println(string(b))
		return nil
	})
	resource := NewRestResource("GET", "/systems?limit=10&filter={{ .filter }}", router)
	api.AddResource("systems", resource)

	params := map[string]string{"filter": "ocp"}
	if err := api.Call("systems", params, nil); err != nil {
		t.Fatal(err)
	}
}
