package project_index

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/pkg/errors"

	"gae-go-recruiting-server/adapter"
	"gae-go-recruiting-server/domain"
	"gae-go-recruiting-server/infra/algolia"
	pb "gae-go-recruiting-server/proto/go/pb"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/opt"
)

type repository struct {
	cf     adapter.AlgoliaClientFactory
	logger adapter.CompositeLogger
}

func NewRepository(cf adapter.AlgoliaClientFactory, logger adapter.CompositeLogger) adapter.ProjectIndexRepository {
	return &repository{
		cf:     cf,
		logger: logger,
	}
}

func buildFilter(
	status pb.Project_Status) string {
	str := fmt.Sprintf("status=%d", int32(status))
	if status == pb.Project_Status_Unknown {
		str = ""
	}

	return strings.TrimSpace(str)
}

func (r *repository) SearchByStatusWithConditionWithPager(
	ctx context.Context,
	query string,
	status pb.Project_Status,
	pager *domain.Pager) ([]domain.ProjectIDWithHighlight, error) {
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

	res, err := r.cf(ctx).InitIndex(indexName).Search(query, params...)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var indexes []index

	if err = res.UnmarshalHits(&indexes); err != nil {
		return nil, err
	}

	ids := make([]domain.ProjectIDWithHighlight, 0, len(res.Hits))
	for _, index := range indexes {
		highlights := make([]*domain.SearchHighlight, 0)

		for key, val := range index.Highlighted {
			body, err := json.Marshal(val)
			if err != nil {
				continue
			}

			switch key {
			case "occupations":
				fallthrough
			case "companyBusinesses":
				var highlightedList []search.HighlightedResult
				if err := json.Unmarshal(body, &highlightedList); err != nil {
					continue
				}

				for _, highlighted := range highlightedList {
					if len(highlighted.MatchedWords) > 0 {
						highlights = append(highlights, &domain.SearchHighlight{
							Key:         key,
							Val:         highlighted.Value,
							MatchedWord: strings.Join(highlighted.MatchedWords, ""),
						})
					}
				}
			default:
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
		}

		ids = append(ids, domain.ProjectIDWithHighlight{
			ID:         domain.ProjectID(index.ObjectID),
			Highlights: highlights,
		})
	}

	return ids, nil
}

func (r *repository) SearchByConditionWithPager(
	ctx context.Context,
	query string,
	pager *domain.Pager) ([]domain.ProjectIDWithHighlight, error) {
	queryLength := utf8.RuneCountInString(query)
	if algolia.MaxQueryLength < queryLength {
		return nil, domain.ErrQueryIsTooLong
	}

	filter := buildFilter(pb.Project_Status_Unknown)
	params := []interface{}{
		opt.HitsPerPage(pager.Limit()),
		opt.Page(pager.Page() - 1),
		opt.Filters(filter),
	}

	res, err := r.cf(ctx).InitIndex(indexName).Search(query, params...)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var indexes []index

	if err = res.UnmarshalHits(&indexes); err != nil {
		return nil, err
	}

	ids := make([]domain.ProjectIDWithHighlight, 0, len(res.Hits))
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

		ids = append(ids, domain.ProjectIDWithHighlight{
			ID:         domain.ProjectID(index.ObjectID),
			Highlights: highlights,
		})
	}

	return ids, nil
}

func (r *repository) SearchCountByStatusWithCondition(
	ctx context.Context,
	query string,
	status pb.Project_Status) (int64, error) {
	queryLength := utf8.RuneCountInString(query)
	if algolia.MaxQueryLength < queryLength {
		return 0, domain.ErrQueryIsTooLong
	}

	filter := buildFilter(status)
	params := []interface{}{
		opt.Filters(filter),
	}

	res, err := r.cf(ctx).InitIndex(indexName).Search(query, params...)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return int64(res.NbHits), nil
}

func (r *repository) SearchCountByCondition(
	ctx context.Context,
	query string) (int64, error) {
	queryLength := utf8.RuneCountInString(query)
	if algolia.MaxQueryLength < queryLength {
		return 0, domain.ErrQueryIsTooLong
	}

	filter := buildFilter(pb.Project_Status_Unknown)
	params := []interface{}{
		opt.Filters(filter),
	}

	res, err := r.cf(ctx).InitIndex(indexName).Search(query, params...)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return int64(res.NbHits), nil
}

func (r *repository) Put(ctx context.Context, project *domain.Project, company *domain.Company) error {
	index := r.cf(ctx).InitIndex(indexName)

	if _, err := index.SaveObject(newIndex(project, company)); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *repository) Delete(ctx context.Context, projectID domain.ProjectID) error {
	index := r.cf(ctx).InitIndex(indexName)

	if _, err := index.DeleteObject(projectID.String()); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *repository) DeleteMulti(ctx context.Context, projectIDs []domain.ProjectID) error {
	index := r.cf(ctx).InitIndex(indexName)

	ids := make([]string, 0, len(projectIDs))
	for _, id := range projectIDs {
		ids = append(ids, id.String())
	}

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

	return nil
}
