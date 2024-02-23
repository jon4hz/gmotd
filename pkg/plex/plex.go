package plex

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"runtime"
	"time"

	"github.com/jon4hz/gmotd/version"
)

type Client struct {
	http *http.Client

	server string
	token  string
}

func New(server, token string, timeout int, tlsVerify bool) *Client {
	return &Client{
		http: &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: !tlsVerify},
			},
		},
		server: server,
		token:  token,
	}
}

func (c *Client) CountSessions(ctx context.Context) (int, error) {
	u, err := url.Parse(c.server + "/status/sessions")
	if err != nil {
		return 0, err
	}
	v := u.Query()
	v.Set("X-Plex-Token", c.token)
	v.Set("X-Plex-Platform", runtime.GOOS)
	v.Set("X-Plex-Platform-Version", "0.0.0")
	v.Set("X-Plex-Client-Identifier", "gmotd-v"+version.Version)
	v.Set("X-Plex-Product", "gmotd")
	v.Set("X-Plex-Version", version.Version)
	v.Set("X-Plex-Device", runtime.GOOS+" "+runtime.GOARCH)
	u.RawQuery = v.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var cs CurrentSessions
	if err := json.Unmarshal(bodyBytes, &cs); err != nil {
		return 0, err
	}

	return cs.MediaContainer.Size, nil
}

type CurrentSessions struct {
	MediaContainer struct {
		Size int `json:"size"`
	} `json:"MediaContainer"`
}
