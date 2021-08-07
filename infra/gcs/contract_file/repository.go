package contract_file

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"gae-go-recruiting-server/adapter"
	"gae-go-recruiting-server/domain"

	"cloud.google.com/go/storage"
	"github.com/pkg/errors"
)

type repository struct {
	bucketName string
}

func NewRepository(bucketName string) adapter.ContractFileRepository {
	return &repository{
		bucketName: bucketName,
	}
}

func (r *repository) Get(ctx context.Context, contract *domain.Contract) (*domain.File, error) {
	s, err := storage.NewClient(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	defer func() {
		_ = s.Close()
	}()

	if contract.GSFileURL == nil {
		return nil, errors.WithStack(domain.ErrContractFileIsNotExists)
	}

	paths := strings.Split(contract.GSFileURL.Path, "/")
	path := strings.Join(paths[1:], "/")

	reader, err := s.Bucket(r.bucketName).Object(path).NewReader(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer func() {
		_ = reader.Close()
	}()

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(reader); err != nil {
		return nil, err
	}

	name := fmt.Sprintf("contract-%s", contract.ID())

	return domain.NewFile(name, buf.Bytes(), reader.ContentType()), nil
}
