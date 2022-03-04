package rhsm

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/oauth2"
	oauthClient "golang.org/x/oauth2/clientcredentials"
)

type Client struct {
	Client *http.Client
}

func NewClient(token string) *Client {
	return &Client{Client: newClient(token)}
}

func newClient(token string) *http.Client {
	values := url.Values{}
	values.Set("grant_type", "refresh_token")
	values.Set("refresh_token", token)

	config := oauthClient.Config{
		ClientID:       "rhsm-api",
		TokenURL:       "https://sso.redhat.com/auth/realms/redhat-external/protocol/openid-connect/token",
		EndpointParams: values,
		AuthStyle:      oauth2.AuthStyleInParams,
	}
	return config.Client(context.Background())
}

func (c *Client) ProcessRequest(baseurl string, resource *RestResource, params map[string]string, body interface{}) error {
	trimmedEndpoint := strings.TrimLeft(resource.RenderEndpoint(params), "/")
	trimmedBaseURL := strings.TrimRight(baseurl, "/")
	trimmedURL := trimmedBaseURL + "/" + trimmedEndpoint
	req, err := buildRequest(resource.Method, trimmedURL, body)
	if err != nil {
		return err
	}
	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	return resource.Router.Call(resp)
}

func buildRequest(method, url string, body interface{}) (*http.Request, error) {
	if body == nil {
		return http.NewRequest(method, url, nil)
	}
	b, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(b)
	return http.NewRequest(method, url, buf)
}
