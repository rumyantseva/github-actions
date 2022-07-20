// Copyright 2021 FerretDB Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package graphql

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/sethvargo/go-githubactions"
	"github.com/shurcooL/githubv4"
	"golang.org/x/net/context/ctxhttp"
	"golang.org/x/oauth2"
)

// GraphQLClient returns GitHub GraphQL client instance with an access token provided from GitHub Actions.
// The token to access API must be provided in the environment variable named `tokenVar`.
func GraphQLClient(ctx context.Context, action *githubactions.Action, tokenVar string) (*Client, error) {
	token := action.Getenv(tokenVar)
	if token == "" {
		return nil, fmt.Errorf("env %s is not set", tokenVar)
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	httpClient := oauth2.NewClient(ctx, ts)
	qlClient := NewClient(httpClient)

	// check that the client is able to make queries,
	// for that we call a simple rate limit query
	var q = `{rateLimit {cost limit remaining resetAt}}`
	var rl struct {
		RateLimit struct {
			Cost      githubv4.Int
			Limit     githubv4.Int
			Remaining githubv4.Int
			ResetAt   githubv4.DateTime
		}
	}
	res, err := qlClient.do(ctx, q, nil)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(res.Data, &rl)

	action.Debugf(
		"Rate limit remaining: %d, reset at: %s",
		rl.RateLimit.Remaining, rl.RateLimit.ResetAt.Format(time.RFC822),
	)

	return qlClient, nil
}

type Client struct {
	httpClient *http.Client
	url        string
}

func NewClient(httpClient *http.Client) *Client {
	return &Client{
		httpClient: httpClient,
		url:        "https://api.github.com/graphql",
	}
}

type out struct {
	Data   json.RawMessage
	Errors json.RawMessage
}

// do runs a GraphQL query.
// The implementation is based on github.com/shurcooL/graphql, but the response is returned a raw json.
func (c *Client) do(ctx context.Context, query string, variables map[string]any) (*out, error) {
	query = fmt.Sprintf("%s", query)
	in := struct {
		Query     string                 `json:"query"`
		Variables map[string]interface{} `json:"variables,omitempty"`
	}{
		Query:     query,
		Variables: variables,
	}
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(in)
	if err != nil {
		return nil, err
	}
	resp, err := ctxhttp.Post(ctx, c.httpClient, c.url, "application/json", &buf)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("non-200 OK status code: %v body: %q", resp.Status, body)
	}
	o := new(out)
	err = json.NewDecoder(resp.Body).Decode(o)
	return o, err
}
