package gcs

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"

	"gae-go-recruiting-server/adapter"

	"cloud.google.com/go/storage"
)

func NewSignature(projectID string, encodedPrivateKey string) adapter.GCSSignature {
	return func(ctx context.Context, gsURL *url.URL) (*url.URL, error) {
		if gsURL == nil {
			return nil, nil
		}

		if gsURL.String() == "" {
			return nil, nil
		}

		paths := strings.Split(gsURL.Path, "/")
		bucketID := gsURL.Host
		objectID := strings.Join(paths[1:], "/")

		expires := time.Now().Add(time.Hour * 1)

		privateKey, _ := base64.StdEncoding.DecodeString(encodedPrivateKey)

		urlStringWithSignature, err := storage.SignedURL(bucketID, objectID, &storage.SignedURLOptions{
			GoogleAccessID: fmt.Sprintf("%s@appspot.gserviceaccount.com", projectID),
			PrivateKey:     privateKey,
			Method:         "GET",
			Expires:        expires,
		})
		if err != nil {
			return nil, errors.WithStack(err)
		}

		urlWithSignature, err := gsURL.Parse(urlStringWithSignature)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		return urlWithSignature, nil
	}

}
