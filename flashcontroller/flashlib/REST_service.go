package flashlib

// Description for this file
// 1. Starts REST service
// 2. Registers REST handlers
// 3. Handlers are defined as well

import (
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

// register routes and their respective Handlers
func muxRouter(r *mux.Router) {
	// Routes consist of a path and a handler function.
	r.HandleFunc("/container/create", Create).Methods("POST")
	r.HandleFunc("/container/list", ContList).Methods("POST")
	r.HandleFunc("/container/start", ContStart).Methods("POST")
	r.HandleFunc("/container/stop", ContStop).Methods("POST")
	r.HandleFunc("/container/stat", ContStat).Methods("POST")
	r.HandleFunc("/container/destroy", ContDestroy).Methods("POST")

	r.HandleFunc("/agent/register", Register).Methods("POST")
	r.HandleFunc("/agent/unregister", UnRegister).Methods("POST")
	r.HandleFunc("/agent/status", AgentStatus).Methods("POST")
	r.HandleFunc("/agent/update", AgentUpdate).Methods("POST")

	r.HandleFunc("/container/snapshot", ContSnap).Methods("POST")
	r.HandleFunc("/container/snapshot/list", SnapList).Methods("POST")
	r.HandleFunc("/container/snapshot/destroy", SnapDestroy).Methods("POST")

	r.PathPrefix("/").Handler(http.FileServer(http.Dir("/home/user/experiment/web-stuff/")))
}

// StartRESTServer starts the REST service
// Listens to input via REST calls
func StartRESTServer(wg *sync.WaitGroup) {
	defer wg.Done()
	r := mux.NewRouter()
	muxRouter(r)
	http.ListenAndServe(":8989", r)
}
