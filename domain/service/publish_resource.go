package service

import (
	"context"
	"net/url"
	"strings"

	"github.com/pkg/errors"

	"gae-go-recruiting-server/adapter"
)

func NewPublishResourceService(signature adapter.GCSSignature) adapter.PublishResourceService {
	return func(ctx context.Context, resourceURL *url.URL) (*url.URL, error) {
		if resourceURL == nil || resourceURL.String() == "" {
			return nil, nil
		}

		if strings.HasPrefix(resourceURL.String(), "gs://") {
			sigURL, err := signature(ctx, resourceURL)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			return sigURL, nil
		}

		return resourceURL, nil
	}
}
