package scratchauth

import (
	"io"
	"fmt"
	"context"
	"net/http"
	"encoding/json"

	"github.com/nekto-sns/nekto-server/app/model"
)

type client struct{
	httpClient     *http.Client
	scratchAuthURL string
	redirectURIs   []string
}

type scratchAuth struct{
	Valid    bool   `json:"valid"`
	Username string `json:"username"`
	Redirect string `json:"redirect"`
}

func New(http *http.Client, scratchAuthURL string, redirectURIs []string) *client {
	return &client{
		httpClient: http,
		scratchAuthURL: scratchAuthURL,
		redirectURIs: redirectURIs,
	}
}

func (c *client) isAllowedRedirectURL(url string) bool {
	for _, v := range c.redirectURIs {
		if v == url {
			return true
		}
	}
	return false
}

func (c *client) Verify(ctx context.Context, privateCode string) (string, bool, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.scratchAuthURL + "/api/auth/verifyToken?privateCode=" + privateCode, nil)
	if err != nil {
		return "", false, fmt.Errorf("Failed to create HTTP request (%v): %w", err, model.ErrInternal)
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return "", false, fmt.Errorf("Failed to send HTTP request (%v): %w", err, model.ErrInternal)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", false, fmt.Errorf("Failed to read response body (%v): %w", err, model.ErrInternal)
	}

	var auth scratchAuth
	err = json.Unmarshal(body, &auth)
	if err != nil {
		return "", false, fmt.Errorf("Failed to unmarshal JSON into a struct (%v): %w", err, model.ErrInternal)
	}

	if auth.Valid == false {
		return "", false, fmt.Errorf("Invalid authentication: %w", model.ErrForbidden)
	}

	if !c.isAllowedRedirectURL(auth.Redirect) {
		return "", false, fmt.Errorf("Invalid redirect URL (%s): %w", auth.Redirect, model.ErrBadRequest)
	}

	return auth.Username, auth.Valid, nil
}
