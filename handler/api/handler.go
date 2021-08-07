package api

import (
	"context"
	"net/http"

	"gae-go-recruiting-server/adapter"
	"gae-go-recruiting-server/handler"

	"github.com/twitchtv/twirp"

	pb "gae-go-recruiting-server/proto/go/pb"
)

type Handler func(mux *http.ServeMux)

func NewHandler(
	applyClientService pb.ApplyClientService,
	projectService pb.ProjectService,
	companyService pb.CompanyService,
	userService pb.UserService,
	clientService pb.ClientService,
	customerService pb.CustomerService,
	entryService pb.EntryService,
	messageRoomService pb.MessageRoomService,
	messageService pb.MessageService,
	prefectureService pb.PrefectureService,
	contractService pb.ContractService,
	checkMaintenance adapter.CheckMaintenance,
	captureHTTP adapter.CaptureHTTP,
	cros adapter.CROS,
	userAuth adapter.UserAuthenticate) Handler {

	public := func(server pb.TwirpServer) http.Handler {
		return handler.ApplyMiddleware(
			server,
			checkMaintenance,
			captureHTTP,
			cros)
	}

	auth := func(server pb.TwirpServer) http.Handler {
		return handler.ApplyMiddleware(
			server,
			checkMaintenance,
			userAuth,
			captureHTTP,
			cros)
	}

	return func(mux *http.ServeMux) {
		hooks := &twirp.ServerHooks{
			RequestReceived: func(ctx context.Context) (context.Context, error) {
				return ctx, nil
			},
			ResponseSent: func(ctx context.Context) {

			},
			Error: func(ctx context.Context, err twirp.Error) context.Context {
				return ctx
			},
		}

		mux.Handle(pb.ApplyClientServicePathPrefix, public(pb.NewApplyClientServiceServer(applyClientService, hooks)))

		mux.Handle(pb.ProjectServicePathPrefix, auth(pb.NewProjectServiceServer(projectService, hooks)))
		mux.Handle(pb.CompanyServicePathPrefix, auth(pb.NewCompanyServiceServer(companyService, hooks)))
		mux.Handle(pb.UserServicePathPrefix, auth(pb.NewUserServiceServer(userService, hooks)))
		mux.Handle(pb.ClientServicePathPrefix, auth(pb.NewClientServiceServer(clientService, hooks)))
		mux.Handle(pb.CustomerServicePathPrefix, auth(pb.NewCustomerServiceServer(customerService, hooks)))
		mux.Handle(pb.EntryServicePathPrefix, auth(pb.NewEntryServiceServer(entryService, hooks)))
		mux.Handle(pb.MessageRoomServicePathPrefix, auth(pb.NewMessageRoomServiceServer(messageRoomService, hooks)))
		mux.Handle(pb.MessageServicePathPrefix, auth(pb.NewMessageServiceServer(messageService, hooks)))
		mux.Handle(pb.PrefectureServicePathPrefix, auth(pb.NewPrefectureServiceServer(prefectureService, hooks)))
		mux.Handle(pb.ContractServicePathPrefix, auth(pb.NewContractServiceServer(contractService, hooks)))
	}
}
