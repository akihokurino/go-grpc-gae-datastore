package customer_index

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/pkg/errors"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/opt"

	"gae-go-sample/adapter"
	"gae-go-sample/domain"
	"gae-go-sample/infra/algolia"
	pb "gae-go-sample/proto/go/pb"
)

type repository struct {
	cf     adapter.AlgoliaClientFactory
	logger adapter.CompositeLogger
}

func NewRepository(cf adapter.AlgoliaClientFactory, logger adapter.CompositeLogger) adapter.CustomerIndexRepository {
	return &repository{
		cf:     cf,
		logger: logger,
	}
}

func buildFilter(
	status pb.Customer_Status) string {
	str := fmt.Sprintf("status=%d", int32(status))
	if status == pb.Customer_Status_Unknown {
		str = ""
	}

	return strings.TrimSpace(str)
}

func (r *repository) SearchByStatusWithConditionWithPager(
	ctx context.Context,
	query string,
	status pb.Customer_Status,
	pager *domain.Pager,
	orderBy pb.SearchCustomerRequest_OrderBy) ([]domain.CustomerIDWithHighlight, error) {
	queryLength := utf8.RuneCountInString(query)
	if algolia.MaxQueryLength < queryLength {
		return nil, domain.ErrQueryIsTooLong
	}

	filter := buildFilter(status)
	params := []interface{}{
		opt.HitsPerPage(pager.Limit()),
		opt.Page(pager.Page() - 1),
		opt.Filters(filter),
	}

	targetIndex := indexName
	switch orderBy {
	case pb.SearchCustomerRequest_OrderBy_CreatedAt_DESC:
		targetIndex = indexNameForClientOrderByCreatedAt
	case pb.SearchCustomerRequest_OrderBy_Bookmark_DESC:
		targetIndex = indexNameForClientOrderByBookmark
	default:
		targetIndex = indexName
	}

	res, err := r.cf(ctx).InitIndex(targetIndex).Search(query, params...)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var indexes []index

	if err = res.UnmarshalHits(&indexes); err != nil {
		return nil, err
	}

	ids := make([]domain.CustomerIDWithHighlight, 0, len(res.Hits))
	for _, index := range indexes {
		highlights := make([]*domain.SearchHighlight, 0)

		for key, val := range index.Highlighted {
			body, err := json.Marshal(val)
			if err != nil {
				continue
			}

			var highlighted search.HighlightedResult
			if err := json.Unmarshal(body, &highlighted); err != nil {
				continue
			}

			if len(highlighted.MatchedWords) > 0 {
				highlights = append(highlights, &domain.SearchHighlight{
					Key:         key,
					Val:         highlighted.Value,
					MatchedWord: strings.Join(highlighted.MatchedWords, ""),
				})
			}
		}

		ids = append(ids, domain.CustomerIDWithHighlight{
			ID:         domain.CustomerID(index.ObjectID),
			Highlights: highlights,
		})
	}

	return ids, nil
}

func (r *repository) SearchCountByStatusWithCondition(
	ctx context.Context,
	query string,
	status pb.Customer_Status) (int64, error) {
	queryLength := utf8.RuneCountInString(query)
	if algolia.MaxQueryLength < queryLength {
		return 0, domain.ErrQueryIsTooLong
	}

	filter := buildFilter(status)
	params := []interface{}{
		opt.Filters(filter),
	}

	targetIndex := indexName

	res, err := r.cf(ctx).InitIndex(targetIndex).Search(query, params...)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return int64(res.NbHits), nil
}

func (r *repository) Put(
	ctx context.Context,
	customer *domain.Customer) error {
	index := r.cf(ctx).InitIndex(indexName)
	if _, err := index.SaveObject(newIndex(customer)); err != nil {
		return errors.WithStack(err)
	}

	index = r.cf(ctx).InitIndex(indexNameForClientOrderByCreatedAt)
	if _, err := index.SaveObject(newIndex(customer)); err != nil {
		return errors.WithStack(err)
	}

	index = r.cf(ctx).InitIndex(indexNameForClientOrderByBookmark)
	if _, err := index.SaveObject(newIndex(customer)); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *repository) Delete(ctx context.Context, customerID domain.CustomerID) error {
	index := r.cf(ctx).InitIndex(indexName)
	if _, err := index.DeleteObject(customerID.String()); err != nil {
		return errors.WithStack(err)
	}

	index = r.cf(ctx).InitIndex(indexNameForClientOrderByCreatedAt)
	if _, err := index.DeleteObject(customerID.String()); err != nil {
		return errors.WithStack(err)
	}

	index = r.cf(ctx).InitIndex(indexNameForClientOrderByBookmark)
	if _, err := index.DeleteObject(customerID.String()); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *repository) DeleteMulti(ctx context.Context, customerIDs []domain.CustomerID) error {
	ids := make([]string, 0, len(customerIDs))
	for _, id := range customerIDs {
		ids = append(ids, id.String())
	}

	index := r.cf(ctx).InitIndex(indexName)
	if _, err := index.DeleteObjects(ids); err != nil {
		return errors.WithStack(err)
	}

	index = r.cf(ctx).InitIndex(indexNameForClientOrderByCreatedAt)
	if _, err := index.DeleteObjects(ids); err != nil {
		return errors.WithStack(err)
	}

	index = r.cf(ctx).InitIndex(indexNameForClientOrderByBookmark)
	if _, err := index.DeleteObjects(ids); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *repository) DeleteAll(ctx context.Context) error {
	index := r.cf(ctx).InitIndex(indexName)
	if _, err := index.ClearObjects(); err != nil {
		return errors.WithStack(err)
	}

	index = r.cf(ctx).InitIndex(indexNameForClientOrderByCreatedAt)
	if _, err := index.ClearObjects(); err != nil {
		return errors.WithStack(err)
	}

	index = r.cf(ctx).InitIndex(indexNameForClientOrderByBookmark)
	if _, err := index.ClearObjects(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
