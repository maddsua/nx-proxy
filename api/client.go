package api

import (
	"net/http"
	"net/url"

	nxproxy "github.com/maddsua/nx-proxy"
)

type Client struct {
	URL   *url.URL
	Token *nxproxy.ServerToken
}

func (client *Client) PostMetrics(metrics *ModelMetrics) error {
	return beacon(client, http.MethodPost, "/metrics", metrics)
}

func (client *Client) GetProxyTable() (*ModelTable, error) {
	return fetch[ModelTable](client, http.MethodPost, "/table", nil)
}
