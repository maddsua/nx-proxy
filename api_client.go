package nxproxy

import (
	"net/http"
	"net/url"

	"github.com/maddsua/nx-proxy/api_models"
)

type Client struct {
	URL   *url.URL
	Token *ServerToken
}

func (client *Client) PostMetrics(metrics *api_models.Metrics) error {
	return beacon(client, http.MethodPost, "/metrics", metrics)
}

func (client *Client) PullTable() (*api_models.Table, error) {
	return fetch[api_models.Table](client, http.MethodGet, "/table", nil)
}
