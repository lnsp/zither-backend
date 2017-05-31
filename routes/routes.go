package routes

import (
	"encoding/json"
	"net/http"

	"github.com/zither-oss/zither-backend/player"
)

type responseStatus struct {
	Status string `json:"status"`
}

type Router struct {
	internal http.ServeMux
	backend  player.Player
}

func mustSendJSON(w http.ResponseWriter, obj interface{}) {
	json, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}
	if _, err := w.Write(json); err != nil {
		panic(err)
	}
}

func (router *Router) init() {
	// Setup routes here
	router.internal.HandleFunc("/play", func(w http.ResponseWriter, r *http.Request) {
		router.backend.Play()
		mustSendJSON(w, responseStatus{"OK"})
	})
	router.internal.HandleFunc("/stop", func(w http.ResponseWriter, r *http.Request) {
		router.backend.Stop()
		mustSendJSON(w, responseStatus{"OK"})
	})
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router.internal.ServeHTTP(w, r)
}

func New(backend player.Player) *Router {
	r := &Router{backend: backend}
	r.init()
	return r
}
