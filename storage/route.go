package storage

import "github.com/gorilla/mux"

func AddRoutes(r *mux.Router, h *Handler)  {
	r.HandleFunc("/", h.ListAll).Methods("GET")
	r.HandleFunc("/share", h.ShareForm).Methods("GET")
	r.HandleFunc("/share", h.Share).Methods("POST")
}
