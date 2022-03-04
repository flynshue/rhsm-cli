package rhsm

import "testing"

func TestRenderEndpoint(t *testing.T) {
	resource := NewRestResource("GET", "/systems/{{ .systemID }}", NewRouter())
	params := map[string]string{"systemID": "fakeSystemID"}
	got := resource.RenderEndpoint(params)
	want := "/systems/fakeSystemID"
	if got != want {
		t.Fatalf("got %s, wanted %s", got, want)
	}
}
