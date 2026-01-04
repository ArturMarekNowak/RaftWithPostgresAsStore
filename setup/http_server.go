package setup

import (
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/raft"
	"main/internal/database"
	httptutil "main/internal/http"
	"net/http"
)

func HttpServer(httpPort string, raft *raft.Raft, logger hclog.Logger, database *database.PostgresAccessor) {
	httpServer := httptutil.NewHttpServer(raft, logger, database)

	router := http.NewServeMux()
	router.HandleFunc("POST /join", httpServer.JoinNode)
	router.HandleFunc("POST /value", httpServer.SetValue)
	router.HandleFunc("GET /value", httpServer.GetValue)
	router.HandleFunc("DELETE /value", httpServer.DeleteValue)
	err := http.ListenAndServe("127.0.0.1:"+httpPort, router)
	if err != nil {
		panic("Couldn't start HTTP server")
	}
}
