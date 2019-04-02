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
	streams   map[string]*Stream
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
		streams:   make(map[string]*Stream),
	}
	return &tc
}
