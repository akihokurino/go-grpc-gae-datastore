package contract_app

import (
	"context"

	"go.mercari.io/datastore/boom"

	"gae-go-recruiting-server/adapter"
	"gae-go-recruiting-server/domain"

	"github.com/pkg/errors"
)

type adminApplication struct {
	executorID domain.AdminUserID
	*application
}

func (a *adminApplication) GetAllWithPager(ctx context.Context, page int32, offset int32) ([]*domain.Contract, error) {
	contracts, err := a.contractRepository.GetAllWithPager(ctx, domain.NewPager(page, offset))
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

func (a *adminApplication) Get(ctx context.Context, idParams adapter.ContractIDParams) (*domain.Contract, error) {
	contractID := domain.NewContractID(idParams.ProjectID, idParams.CompanyID, idParams.CustomerID)

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

func (a *adminApplication) GetTotalCount(ctx context.Context) (int64, error) {
	count, err := a.contractRepository.GetTotalCount(ctx)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return count, nil
}

func (a *adminApplication) Accept(ctx context.Context, idParams adapter.ContractIDParams) error {
	contractID := domain.NewContractID(idParams.ProjectID, idParams.CompanyID, idParams.CustomerID)

	contract, err := a.contractRepository.Get(ctx, contractID)
	if err != nil {
		return errors.WithStack(err)
	}

	if err := contract.Accept(); err != nil {
		return errors.WithStack(err)
	}

	if err := a.transaction(ctx, func(tx *boom.Transaction) error {
		if err := a.contractRepository.Put(tx, contract); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (a *adminApplication) Cancel(ctx context.Context, idParams adapter.ContractIDParams) error {
	contractID := domain.NewContractID(idParams.ProjectID, idParams.CompanyID, idParams.CustomerID)

	contract, err := a.contractRepository.Get(ctx, contractID)
	if err != nil {
		return errors.WithStack(err)
	}

	if err := contract.Cancel(); err != nil {
		return errors.WithStack(err)
	}

	if err := a.transaction(ctx, func(tx *boom.Transaction) error {
		if err := a.contractRepository.Put(tx, contract); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (a *adminApplication) DownloadFile(
	ctx context.Context,
	idParams adapter.ContractIDParams) (*domain.File, error) {
	contractID := domain.NewContractID(idParams.ProjectID, idParams.CompanyID, idParams.CustomerID)

	contract, err := a.contractRepository.Get(ctx, contractID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	file, err := a.contractFileRepository.Get(ctx, contract)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return file, err
}
