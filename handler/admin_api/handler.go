package admin_api

import (
	"context"
	"net/http"

	"gae-go-sample/adapter"
	"gae-go-sample/handler"

	"github.com/twitchtv/twirp"

	pb "gae-go-sample/proto/go/pb"
)

type Handler func(mux *http.ServeMux)

func NewHandler(
	applyClientService pb.AdminApplyClientService,
	companyService pb.AdminCompanyService,
	projectService pb.AdminProjectService,
	clientService pb.AdminClientService,
	customerService pb.AdminCustomerService,
	entryService pb.AdminEntryService,
	contractService pb.AdminContractService,
	messageRoomService pb.AdminMessageRoomService,
	messageService pb.AdminMessageService,
	captureHTTP adapter.CaptureHTTP,
	cros adapter.CROS,
	checkMaintenance adapter.CheckMaintenance,
	adminAuth adapter.AdminAuthenticate) Handler {

	auth := func(server pb.TwirpServer) http.Handler {
		return handler.ApplyMiddleware(
			server,
			checkMaintenance,
			adminAuth,
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

		mux.Handle(pb.AdminApplyClientServicePathPrefix, auth(pb.NewAdminApplyClientServiceServer(applyClientService, hooks)))
		mux.Handle(pb.AdminCompanyServicePathPrefix, auth(pb.NewAdminCompanyServiceServer(companyService, hooks)))
		mux.Handle(pb.AdminProjectServicePathPrefix, auth(pb.NewAdminProjectServiceServer(projectService, hooks)))
		mux.Handle(pb.AdminClientServicePathPrefix, auth(pb.NewAdminClientServiceServer(clientService, hooks)))
		mux.Handle(pb.AdminCustomerServicePathPrefix, auth(pb.NewAdminCustomerServiceServer(customerService, hooks)))
		mux.Handle(pb.AdminEntryServicePathPrefix, auth(pb.NewAdminEntryServiceServer(entryService, hooks)))
		mux.Handle(pb.AdminContractServicePathPrefix, auth(pb.NewAdminContractServiceServer(contractService, hooks)))
		mux.Handle(pb.AdminMessageRoomServicePathPrefix, auth(pb.NewAdminMessageRoomServiceServer(messageRoomService, hooks)))
		mux.Handle(pb.AdminMessageServicePathPrefix, auth(pb.NewAdminMessageServiceServer(messageService, hooks)))
	}
}
