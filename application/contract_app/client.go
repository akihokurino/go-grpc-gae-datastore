package contract_app

import (
	"context"
	"net/url"
	"sync"
	"time"

	"go.mercari.io/datastore/boom"

	"gae-go-recruiting-server/adapter"
	"gae-go-recruiting-server/domain"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

type clientApplication struct {
	executorID domain.ClientID
	*application
}

func (a *clientApplication) GetAllWithPager(
	ctx context.Context,
	page int32,
	offset int32) ([]*domain.Contract, error) {
	me, err := a.clientRepository.Get(ctx, a.executorID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	contracts, err := a.contractRepository.GetAllByCompanyWithPager(ctx, me.CompanyID, domain.NewPager(page, offset))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	for i := range contracts {
		urlWithSignature, err := a.publishResourceService(ctx, contracts[i].GSFileURL)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		contracts[i].SignedFileURL = urlWithSignature
	}

	return contracts, nil
}

func (a *clientApplication) GetAllNewestByIDs(
	ctx context.Context,
	idParams []adapter.ContractIDWithoutCompanyIDParams) ([]*domain.Contract, error) {
	me, err := a.clientRepository.Get(ctx, a.executorID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	contractMap := make(map[domain.ContractID]*domain.Contract, 0)

	mutex := sync.Mutex{}
	eg := errgroup.Group{}

	for i := range idParams {
		idParams := idParams[i]

		eg.Go(func() error {
			id := domain.NewContractID(idParams.ProjectID, me.CompanyID, idParams.CustomerID)
			contract, err := a.contractRepository.Get(ctx, id)
			if err != nil && !domain.IsNoSuchEntityErr(err) {
				return err
			}
			if err != nil && domain.IsNoSuchEntityErr(err) {
				return nil
			}

			mutex.Lock()
			contractMap[id] = contract
			mutex.Unlock()

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, errors.WithStack(err)
	}

	contracts := make([]*domain.Contract, 0, len(contractMap))
	for _, contract := range contractMap {
		contracts = append(contracts, contract)
	}

	for i := range contracts {
		urlWithSignature, err := a.publishResourceService(ctx, contracts[i].GSFileURL)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		contracts[i].SignedFileURL = urlWithSignature
	}

	return contracts, nil
}

func (a *clientApplication) Get(
	ctx context.Context,
	idParams adapter.ContractIDWithoutCompanyIDParams) (*domain.Contract, error) {
	me, err := a.clientRepository.Get(ctx, a.executorID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	company, err := a.companyRepository.Get(ctx, me.CompanyID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	contractID := domain.NewContractID(idParams.ProjectID, company.ID, idParams.CustomerID)

	contract, err := a.contractRepository.Get(ctx, contractID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	urlWithSignature, err := a.publishResourceService(ctx, contract.GSFileURL)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	contract.SignedFileURL = urlWithSignature

	return contract, nil
}

func (a *clientApplication) GetCount(ctx context.Context) (int64, error) {
	me, err := a.clientRepository.Get(ctx, a.executorID)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	count, err := a.contractRepository.GetCountByCompany(ctx, me.CompanyID)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return count, nil
}

func (a *clientApplication) Create(
	ctx context.Context,
	idParams adapter.ContractIDWithoutCompanyIDParams,
	fileURL *url.URL,
	now time.Time) (*domain.Contract, error) {
	me, err := a.clientRepository.Get(ctx, a.executorID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	company, err := a.companyRepository.Get(ctx, me.CompanyID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	project, err := a.projectRepository.Get(ctx, idParams.ProjectID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	customer, err := a.customerRepository.Get(ctx, idParams.CustomerID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	isExists, err := a.entryRepository.Exists(ctx, domain.NewEntryID(customer.ID, project.ID))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if !isExists {
		return nil, errors.WithStack(domain.ErrCustomerDidNotEntry)
	}

	isExists, err = a.contractRepository.Exists(ctx, domain.NewContractID(project.ID, company.ID, customer.ID))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if isExists {
		return nil, errors.WithStack(domain.ErrContractAlreadyExists)
	}

	contract, err := company.Contract(project, customer, fileURL, now)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if err := a.transaction(ctx, func(tx *boom.Transaction) error {
		if err := a.contractRepository.Put(tx, contract); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, errors.WithStack(err)
	}

	urlWithSignature, err := a.publishResourceService(ctx, contract.GSFileURL)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	contract.SignedFileURL = urlWithSignature

	return contract, nil
}

func (a *clientApplication) Update(
	ctx context.Context,
	idParams adapter.ContractIDWithoutCompanyIDParams,
	fileURL *url.URL,
	now time.Time) (*domain.Contract, error) {
	me, err := a.clientRepository.Get(ctx, a.executorID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	contractID := domain.NewContractID(idParams.ProjectID, me.CompanyID, idParams.CustomerID)

	contract, err := a.contractRepository.Get(ctx, contractID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if err := contract.Update(fileURL, now); err != nil {
		return nil, errors.WithStack(err)
	}

	if err := a.transaction(ctx, func(tx *boom.Transaction) error {
		if err := a.contractRepository.Put(tx, contract); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, errors.WithStack(err)
	}

	urlWithSignature, err := a.publishResourceService(ctx, contract.GSFileURL)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	contract.SignedFileURL = urlWithSignature

	return contract, nil
}

func (a *clientApplication) Delete(
	ctx context.Context,
	idParams adapter.ContractIDWithoutCompanyIDParams) error {
	me, err := a.clientRepository.Get(ctx, a.executorID)
	if err != nil {
		return errors.WithStack(err)
	}

	contractID := domain.NewContractID(idParams.ProjectID, me.CompanyID, idParams.CustomerID)

	if err := a.transaction(ctx, func(tx *boom.Transaction) error {
		if err := a.contractRepository.Delete(tx, contractID); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
