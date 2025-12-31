package http

import (
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/raft"
	"main/internal/database"
	"net/http"
)

type HttpServer struct {
	R *raft.Raft
	L hclog.Logger
	T *database.PostgresAccessor
}

func (hs HttpServer) JoinNode(w http.ResponseWriter, r *http.Request) {
	/*followerId := r.URL.Query().Get("followerId")
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
		log.Printf("Failed to add follower: %s", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusOK)*/
}

func (hs HttpServer) SetValue(w http.ResponseWriter, r *http.Request) {
	/*defer r.Body.Close()
	bs, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Could not read key-value in http request: %s", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	future := hs.R.Apply(bs, 500*time.Millisecond)
	if err := future.Error(); err != nil {
		log.Printf("Could not write key-value: %s", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	e := future.Response()
	if e != nil {
		log.Printf("Could not write key-value, application: %s", e)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)*/
}

func (hs HttpServer) GetValue(w http.ResponseWriter, r *http.Request) {
	/*key := r.URL.Query().Get("key")
	value, _ := hs.T.Find([]byte(key))
	if value == nil {
		value = []byte("")
	}

	rsp := struct {
		Data string `json:"data"`
	}{string(value[:])}
	err := json.NewEncoder(w).Encode(rsp)
	if err != nil {
		log.Printf("Could not encode key-value in http response: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}*/
}

func (hs HttpServer) DeleteValue(w http.ResponseWriter, r *http.Request) {
	/*key := r.URL.Query().Get("key")
	deleted := hs.T.Delete([]byte(key))
	if !deleted {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
	w.WriteHeader(http.StatusNoContent)*/
}
