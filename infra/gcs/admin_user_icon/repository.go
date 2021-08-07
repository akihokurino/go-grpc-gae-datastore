package admin_user_icon

import (
	"context"
	"fmt"
	"net/url"

	"gae-go-recruiting-server/adapter"

	"cloud.google.com/go/storage"
	"github.com/pkg/errors"
)

type repository struct {
	bucketName string
}

func NewRepository(bucketName string) adapter.AdminUserFileRepository {
	return &repository{
		bucketName: bucketName,
	}
}

func (r *repository) Save(ctx context.Context, username string, data []byte) (*url.URL, error) {
	s, err := storage.NewClient(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	path := fmt.Sprintf("admin-user/%s", username)

	writer := s.Bucket(r.bucketName).Object(path).NewWriter(ctx)
	defer func() {
		_ = writer.Close()
	}()

	writer.ContentType = "image/png"

	if _, err := writer.Write(data); err != nil {
		return nil, errors.WithStack(err)
	}

	u, err := url.Parse(fmt.Sprintf("gs://%s/%s", r.bucketName, path))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return u, err
}
