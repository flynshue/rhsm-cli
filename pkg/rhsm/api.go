package rhsm

import "fmt"

type API struct {
	BaseURL       string
	Client        *Client
	Resources     map[string]*RestResource
	DefaultRouter *CBRouter
}

func NewAPI(baseurl, token string) *API {
	return &API{
		BaseURL:       baseurl,
		Client:        NewClient(token),
		Resources:     make(map[string]*RestResource),
		DefaultRouter: NewRouter(),
	}
}

func (a *API) AddResource(name string, resource *RestResource) {
	a.Resources[name] = resource
}

func (a *API) ListResources() []string {
	resources := make([]string, 0, len(a.Resources))
	for k := range a.Resources {
		resources = append(resources, k)
	}
	return resources
}

func (a *API) Call(name string, params map[string]string, body interface{}) error {
	resource, ok := a.Resources[name]
	if !ok {
		return fmt.Errorf("resource %s not found", name)
	}
	return a.Client.ProcessRequest(a.BaseURL, resource, params, body)
}
