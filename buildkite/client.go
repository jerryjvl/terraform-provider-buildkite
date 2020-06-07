package buildkite

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"golang.org/x/oauth2"
)

// Client provides a connection to both the Buildkite API and the
// Buildkite GraphQL interface.
type Client struct {
	slug    string
	httpAPI *http.Client
}

func expectToken(decoder *json.Decoder, expected string) error {
	token, err := decoder.Token()
	if err != nil {
		return err
	}
	if fmt.Sprintf("%v", token) != expected {
		return fmt.Errorf("GraphQL expected token '%s', but received '%s' instead", expected, token)
	}
	return nil
}

func (client *Client) Query(result interface{}, query string, args ...interface{}) error {
	requestBody := fmt.Sprintf(`{ "query": "{ `+query+` }" }`, args...)
	res, err := client.httpAPI.Post(
		"https://graphql.buildkite.com/v1",
		"application/json",
		strings.NewReader(requestBody))
	if err != nil {
		return err
	}

	responseBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	decoder := json.NewDecoder(bytes.NewReader(responseBytes))
	if err := expectToken(decoder, "{"); err != nil {
		return err
	}
	if err := expectToken(decoder, "data"); err != nil {
		return err
	}
	if err := expectToken(decoder, "{"); err != nil {
		return err
	}
	// Skip over the name of the top-level object; we only want the contents
	if _, err = decoder.Token(); err != nil {
		return err
	}

	return decoder.Decode(result)
}

// NewClient creates a connection to Buildkite.
func NewClient(apiToken string, organizationSlug string) *Client {
	oauth2Token := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: apiToken})
	httpClient := oauth2.NewClient(context.Background(), oauth2Token)

	return &Client{
		slug:    organizationSlug,
		httpAPI: httpClient,
	}
}
