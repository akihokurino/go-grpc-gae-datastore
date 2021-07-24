package algolia

import (
	"context"

	"gae-go-sample/adapter"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/compression"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
)

const MaxQueryLength = 100

func NewClientFactory(appID string, apiKey string) adapter.AlgoliaClientFactory {
	return func(ctx context.Context) *search.Client {
		return search.NewClientWithConfig(search.Configuration{
			AppID:       appID,
			APIKey:      apiKey,
			Compression: compression.None,
		})
	}
}
