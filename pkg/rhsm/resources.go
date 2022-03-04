package rhsm

import (
	"html/template"
	"log"
	"strings"
)

type RestResource struct {
	Method   string
	Endpoint string
	Router   *CBRouter
}

func NewRestResource(method, endpoint string, router *CBRouter) *RestResource {
	return &RestResource{
		Method:   method,
		Endpoint: endpoint,
		Router:   router,
	}
}

func (r *RestResource) RenderEndpoint(params map[string]string) string {
	if params == nil {
		return r.Endpoint
	}
	t, err := template.New("").Parse(r.Endpoint)
	if err != nil {
		log.Fatal(err)
	}
	str := &strings.Builder{}
	if err := t.Execute(str, params); err != nil {
		log.Fatal(err)
	}
	return str.String()
}
