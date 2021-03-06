// +build wireinject

package di

import (
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"gae-go-recruiting-server/infra/algolia"
	"gae-go-recruiting-server/infra/gcs"

	"gae-go-recruiting-server/application/company_app"
	"gae-go-recruiting-server/application/contract_app"
	"gae-go-recruiting-server/application/customer_app"
	"gae-go-recruiting-server/application/entry_app"
	"gae-go-recruiting-server/application/message_app"
	"gae-go-recruiting-server/application/message_room_app"
	"gae-go-recruiting-server/application/prefecture_app"
	"gae-go-recruiting-server/application/project_app"
	"gae-go-recruiting-server/application/user_app"

	"gae-go-recruiting-server/application/apply_client_app"
	"gae-go-recruiting-server/application/client_app"
	"gae-go-recruiting-server/domain/service"

	"gae-go-recruiting-server/adapter"
	"gae-go-recruiting-server/domain"
	"gae-go-recruiting-server/handler"
	"gae-go-recruiting-server/handler/admin_api"
	"gae-go-recruiting-server/handler/api"
	"gae-go-recruiting-server/handler/batch"
	"gae-go-recruiting-server/handler/middleware"
	"gae-go-recruiting-server/handler/subscriber"
	"gae-go-recruiting-server/infra/algolia/customer_index"
	"gae-go-recruiting-server/infra/algolia/project_index"
	"gae-go-recruiting-server/infra/datastore"
	"gae-go-recruiting-server/infra/datastore/apply_client_table"
	"gae-go-recruiting-server/infra/datastore/client_table"
	"gae-go-recruiting-server/infra/datastore/company_table"
	"gae-go-recruiting-server/infra/datastore/contract_table"
	"gae-go-recruiting-server/infra/datastore/customer_table"
	"gae-go-recruiting-server/infra/datastore/entry_table"
	"gae-go-recruiting-server/infra/datastore/message_room_table"
	"gae-go-recruiting-server/infra/datastore/message_table"
	"gae-go-recruiting-server/infra/datastore/no_entry_support_table"
	"gae-go-recruiting-server/infra/datastore/no_message_support_table"
	"gae-go-recruiting-server/infra/datastore/project_table"
	"gae-go-recruiting-server/infra/datastore/user_table"
	"gae-go-recruiting-server/infra/firebase"
	"gae-go-recruiting-server/infra/firebase/fire_user"
	"gae-go-recruiting-server/infra/firebase/rt_message_room"
	"gae-go-recruiting-server/infra/gcs/admin_user_icon"
	"gae-go-recruiting-server/infra/gcs/contract_file"
	"gae-go-recruiting-server/infra/government"
	"gae-go-recruiting-server/infra/government/prefecture"
	"gae-go-recruiting-server/infra/logger"

	"github.com/google/wire"
)

var providerSet = wire.NewSet(
	provideLogger,
	provideSwitchProvider,
	provideThresholdProvider,
	provideAlgoliaClientFactory,
	provideGovernmentClient,
	provideFireClient,
	provideContractFileRepository,
	provideAdminUserFileRepository,
	provideFireUserRepository,
	provideGCSSignature,
	provideDSFactory,
	datastore.NewClient,
	datastore.NewTransaction,
	apply_client_table.NewRepository,
	company_table.NewRepository,
	project_table.NewRepository,
	user_table.NewRepository,
	client_table.NewRepository,
	customer_table.NewRepository,
	customer_index.NewRepository,
	entry_table.NewRepository,
	message_room_table.NewRepository,
	message_table.NewRepository,
	rt_message_room.NewRtMessageRoomRepository,
	rt_message_room.NewRtMemberRepository,
	project_index.NewRepository,
	prefecture.NewRepository,
	contract_table.NewRepository,
	no_entry_support_table.NewRepository,
	no_message_support_table.NewRepository,

	service.NewBindCustomerEmailService,
	service.NewBindClientEmailService,
	service.NewCustomerAlreadyEntryToAnyoneService,
	service.NewCustomerAlreadyEntryToThisService,
	service.NewRollbackFireUserService,
	service.NewCloseNoEntrySupportService,
	service.NewCloseNoMessageSupportService,
	service.NewOpenNoEntrySupportService,
	service.NewOpenNoMessageSupportService,
	service.NewValidCompanyService,
	service.NewValidProjectService,
	service.NewPublishResourceService,

	domain.NewIDFactory,

	apply_client_app.NewApplication,
	client_app.NewApplication,
	company_app.NewApplication,
	contract_app.NewApplication,
	customer_app.NewApplication,
	entry_app.NewApplication,
	message_room_app.NewApplication,
	message_app.NewApplication,
	prefecture_app.NewApplication,
	project_app.NewApplication,
	user_app.NewApplication,

	handler.NewContextProvider,
	handler.NewErrorConverter,
	admin_api.NewApplyClientHandler,
	admin_api.NewClientHandler,
	admin_api.NewCompanyHandler,
	admin_api.NewContractHandler,
	admin_api.NewCustomerHandler,
	admin_api.NewEntryHandler,
	admin_api.NewMessageHandler,
	admin_api.NewMessageRoomHandler,
	admin_api.NewProjectHandler,
	api.NewApplyClientHandler,
	api.NewClientHandler,
	api.NewCompanyHandler,
	api.NewContractHandler,
	api.NewCustomerHandler,
	api.NewEntryHandler,
	api.NewMessageHandler,
	api.NewMessageRoomHandler,
	api.NewPrefectureHandler,
	api.NewProjectHandler,
	api.NewUserHandler,
	batch.NewProjectHandler,
	subscriber.NewMessageHandler,
	middleware.NewUserAuthenticate,
	middleware.NewAdminAuthenticate,
	middleware.NewBatchAuthenticate,
	middleware.NewCaptureHTTP,
	middleware.NewCROS,
	middleware.NewCheckMaintenance,
	admin_api.NewHandler,
	api.NewHandler,
	batch.NewHandler,
	subscriber.NewHandler,
)

func provideDSFactory() adapter.DSFactory {
	return datastore.NewDSFactory(os.Getenv("PROJECT_ID"))
}

func provideGCSSignature() adapter.GCSSignature {
	return gcs.NewSignature(os.Getenv("PROJECT_ID"), os.Getenv("SERVICE_ACCOUNT_PEM"))
}

func provideLogger() adapter.CompositeLogger {
	return logger.NewLoggerWithMinLevel(adapter.LogLevelDebug)
}

func provideSwitchProvider() adapter.SwitchProvider {
	isMaintenance := os.Getenv("MAINTENANCE") == "True" || os.Getenv("MAINTENANCE") == "true"

	return adapter.SwitchProvider{
		IsMaintenance: isMaintenance,
	}
}

func provideThresholdProvider() adapter.ThresholdProvider {
	noMessageDurationThreshold, err := strconv.Atoi(os.Getenv("NO_MESSAGE_DURATION_THRESHOLD"))
	if err != nil {
		panic("failed inject project application")
	}

	noEntryDurationThreshold, err := strconv.Atoi(os.Getenv("NO_ENTRY_DURATION_THRESHOLD"))
	if err != nil {
		panic("failed inject project application")
	}

	return adapter.ThresholdProvider{
		NoMessageDuration: time.Duration(noMessageDurationThreshold),
		NoEntryDuration:   time.Duration(noEntryDurationThreshold),
	}
}

func provideAlgoliaClientFactory() adapter.AlgoliaClientFactory {
	appID := os.Getenv("ALGOLIA_APP_ID")
	apiKey := os.Getenv("ALGOLIA_API_KEY")
	if appID == "" || apiKey == "" {
		panic("failed inject algolia")
	}

	return algolia.NewClientFactory(appID, apiKey)
}

func provideGovernmentClient() adapter.GovernmentClient {
	endpoint, _ := url.Parse(os.Getenv("GOVERNMENT_ENDPOINT"))
	key := os.Getenv("GOVERNMENT_KEY")
	if endpoint == nil || key == "" {
		panic("failed inject government")
	}

	return government.NewClient(key, endpoint)
}

func provideFireClient() adapter.FirebaseClient {
	rtdbURL, _ := url.Parse(os.Getenv("RTDB_URL"))
	adminKey := os.Getenv("FIREBASE_ADMIN_KEY")
	if rtdbURL == nil || adminKey == "" {
		panic("failed inject firebase")
	}

	return firebase.NewClient(rtdbURL, adminKey)
}

func provideFireUserRepository(
	fireClient adapter.FirebaseClient,
	logger adapter.CompositeLogger) adapter.FireUserRepository {
	emailsString := os.Getenv("ADMIN_EMAILS")
	return fire_user.NewRepository(fireClient, logger, strings.Split(emailsString, "/"))
}

func provideContractFileRepository() adapter.ContractFileRepository {
	bucket := os.Getenv("STORAGE_BUCKET")
	if bucket == "" {
		panic("failed inject admin user file repository")
	}

	return contract_file.NewRepository(bucket)
}

func provideAdminUserFileRepository() adapter.AdminUserFileRepository {
	bucket := os.Getenv("STORAGE_BUCKET")
	if bucket == "" {
		panic("failed inject admin user file repository")
	}

	return admin_user_icon.NewRepository(bucket)
}

func ResolveAdminAPIHandler() admin_api.Handler {
	panic(wire.Build(providerSet))
}

func ResolveAPIHandler() api.Handler {
	panic(wire.Build(providerSet))
}

func ResolveBatchHandler() batch.Handler {
	panic(wire.Build(providerSet))
}

func ResolveSubscriberHandler() subscriber.Handler {
	panic(wire.Build(providerSet))
}
