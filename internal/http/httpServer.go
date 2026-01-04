package http

import (
	"encoding/json"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/raft"
	"io"
	"main/internal/database"
	"main/internal/http/responses"
	"net/http"
	"time"
)

type HttpServer struct {
	raft     *raft.Raft
	logger   hclog.Logger
	database *database.PostgresAccessor
}

func NewHttpServer(r *raft.Raft, logger hclog.Logger, db *database.PostgresAccessor) HttpServer {
	return HttpServer{r, logger, db}
}

func (hs HttpServer) JoinNode(w http.ResponseWriter, r *http.Request) {
	followerId := r.URL.Query().Get("followerId")
	followerAddress := r.URL.Query().Get("followerAddress")

	if hs.raft.State() != raft.Leader {
		err := json.NewEncoder(w).Encode(responses.NewErrorResponse("Not the leader"))
		if err != nil {
			return
		}
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err := hs.raft.AddVoter(raft.ServerID(followerId), raft.ServerAddress(followerAddress), 0, 0).Error()
	if err != nil {
		hs.logger.Error("Failed to add follower: %s", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	hs.logger.Info("Peer joined raft: %s", followerAddress)
	w.WriteHeader(http.StatusOK)
}

func (hs HttpServer) SetValue(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	bs, err := io.ReadAll(r.Body)
	if err != nil {
		hs.logger.Error("Could not read key-value in http request: %s", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	promise := hs.raft.Apply(bs, 500*time.Millisecond)
	if err := promise.Error(); err != nil {
		hs.logger.Error("Could not write key-value: %s", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	e := promise.Response()
	if e != nil {
		hs.logger.Error("Could not write key-value, application: %s", e)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (hs HttpServer) GetValue(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	value, err := hs.database.GetValue(key)
	if err != nil {
		hs.logger.Error("Could not retrieve key-value in http response: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	rsp := responses.NewDataResponse(value[:])
	err = json.NewEncoder(w).Encode(rsp)
	if err != nil {
		hs.logger.Error("Could not encode key-value in http response: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (hs HttpServer) DeleteValue(w http.ResponseWriter, r *http.Request) {
	/*key := r.URL.Query().Get("key")
	deleted := hs.database.Delete([]byte(key))
	if !deleted {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
	w.WriteHeader(http.StatusNoContent)*/
	panic("implement me")
}
