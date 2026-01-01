package http

import (
	"encoding/json"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/raft"
	"io"
	"main/internal/database"
	"net/http"
	"strconv"
	"time"
)

type HttpServer struct {
	R *raft.Raft
	L hclog.Logger
	T *database.PostgresAccessor
}

func (hs HttpServer) JoinNode(w http.ResponseWriter, r *http.Request) {
	followerId := r.URL.Query().Get("followerId")
	followerAddr := r.URL.Query().Get("followerAddr")

	if hs.R.State() != raft.Leader {
		err := json.NewEncoder(w).Encode(struct {
			Error string `json:"error"`
		}{
			"Not the leader",
		})
		if err != nil {
			return
		}
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err := hs.R.AddVoter(raft.ServerID(followerId), raft.ServerAddress(followerAddr), 0, 0).Error()
	if err != nil {
		hs.L.Error("Failed to add follower: %s", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	hs.L.Info("Peer joined raft: %s", followerAddr)
	w.WriteHeader(http.StatusOK)
}

func (hs HttpServer) SetValue(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	bs, err := io.ReadAll(r.Body)
	if err != nil {
		hs.L.Error("Could not read key-value in http request: %s", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	future := hs.R.Apply(bs, 500*time.Millisecond)
	if err := future.Error(); err != nil {
		hs.L.Error("Could not write key-value: %s", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	e := future.Response()
	if e != nil {
		hs.L.Error("Could not write key-value, application: %s", e)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (hs HttpServer) GetValue(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	number, err := strconv.ParseUint(key, 10, 64)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	value, err := hs.T.GetValue(number)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	rsp := struct {
		Data string `json:"data"`
	}{value[:]}
	err = json.NewEncoder(w).Encode(rsp)
	if err != nil {
		hs.L.Error("Could not encode key-value in http response: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (hs HttpServer) DeleteValue(w http.ResponseWriter, r *http.Request) {
	/*key := r.URL.Query().Get("key")
	deleted := hs.T.Delete([]byte(key))
	if !deleted {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
	w.WriteHeader(http.StatusNoContent)*/
	panic("implement me")
}
