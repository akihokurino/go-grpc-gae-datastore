package apply_client_app

import (
	"gae-go-sample/adapter"
	"gae-go-sample/domain"
)

type application struct {
	logger                  adapter.CompositeLogger
	applyClientRepository   adapter.ApplyClientRepository
	fireUserRepository      adapter.FireUserRepository
	userRepository          adapter.UserRepository
	clientRepository        adapter.ClientRepository
	companyRepository       adapter.CompanyRepository
	transaction             adapter.Transaction
	idFactory               domain.IDFactory
	rollbackFireUserService adapter.RollbackFireUserService
}

func NewApplication(
	logger adapter.CompositeLogger,
	applyClientRepository adapter.ApplyClientRepository,
	fireUserRepository adapter.FireUserRepository,
	userRepository adapter.UserRepository,
	clientRepository adapter.ClientRepository,
	companyRepository adapter.CompanyRepository,
	transaction adapter.Transaction,
	idFactory domain.IDFactory,
	rollbackFireUserService adapter.RollbackFireUserService) adapter.ApplyClientApplication {
	return &application{
		logger:                  logger,
		applyClientRepository:   applyClientRepository,
		fireUserRepository:      fireUserRepository,
		userRepository:          userRepository,
		clientRepository:        clientRepository,
		companyRepository:       companyRepository,
		transaction:             transaction,
		idFactory:               idFactory,
		rollbackFireUserService: rollbackFireUserService,
	}
}

func (a *application) BuildAsAdmin(id domain.AdminUserID) adapter.ApplyClientApplicationForAdmin {
	return &adminApplication{
		executorID:  id,
		application: a,
	}
}

func (a *application) BuildAsPublic() adapter.ApplyClientApplicationForPublic {
	return &publicApplication{
		application: a,
	}
}
