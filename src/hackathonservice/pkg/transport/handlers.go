package transport

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	log "github.com/sirupsen/logrus"
	hackathon "go-hackathon/api/hackathonservice"
	"go-hackathon/src/common/cmd"
	"go-hackathon/src/hackathonservice/pkg/hackathon/api"
	"net/http"
)

type server struct {
	api api.Api
}

func Router(ctx context.Context, api api.Api) http.Handler {
	srv := &server{api: api}

	router := runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{EmitDefaults: true, OrigName: true}))
	err := hackathon.RegisterHackathonServiceHandlerServer(ctx, router, srv)
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()
	r.Use(cmd.LogMiddleware)
	r.PathPrefix("/api").Handler(router)
	r.HandleFunc("/hackathon/{ID:[0-9a-zA-Z-]+}/participants", srv.getHackathonParticipantsPage).Methods(http.MethodGet)

	return r
}
