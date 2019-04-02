package tibber

import (
	"net/http"
	"time"

	"github.com/machinebox/graphql"
)

const graphQlEndpoint = "https://api.tibber.com/v1-beta/gql"

// Client for requests and streams
type Client struct {
	Token     string
	gqlClient *graphql.Client
}

// NewClient init tibber client
func NewClient(token string) *Client {
	c := &http.Client{
		Timeout: time.Second * 10,
	}
	gql := graphql.NewClient(graphQlEndpoint, graphql.WithHTTPClient(c))
	tc := Client{
		Token:     token,
		gqlClient: gql,
	}
	return &tc
}
