package tibber

import (
	"net/http"
	"time"

	"github.com/machinebox/graphql"
)

const graphQlEndpoint = "https://api.tibber.com/v1-beta/gql"

type TibberClient struct {
	Token     string
	gqlClient *graphql.Client
}

func NewTibberClient(token string) *TibberClient {
	c := &http.Client{
		Timeout: time.Second * 10,
	}
	gql := graphql.NewClient(graphQlEndpoint, graphql.WithHTTPClient(c))
	tc := TibberClient{
		Token:     token,
		gqlClient: gql,
	}
	return &tc
}
