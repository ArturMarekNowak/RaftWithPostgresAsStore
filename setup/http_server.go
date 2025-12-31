package setup

import (
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/raft"
	"main/internal/database"
	httptutil "main/internal/http"
	"net/http"
	"os"
)

func HttpServer(r *raft.Raft, logger hclog.Logger, db *database.PostgresAccessor) {
	httpServer := &httptutil.HttpServer{r, logger, db}

	router := http.NewServeMux()
	router.HandleFunc("POST /join", httpServer.JoinNode)
	router.HandleFunc("POST /value", httpServer.SetValue)
	router.HandleFunc("GET /value", httpServer.GetValue)
	router.HandleFunc("DELETE /value", httpServer.DeleteValue)
	err := http.ListenAndServe(os.Getenv("HTTP_ADDRESS"), router)
	if err != nil {
		panic("Couldn't start HTTP server")
	}
}
