package storage

import (
	"github.com/sirupsen/logrus"
	"net/http"
)

type Handler struct {
	storage *Storage
	logger *logrus.Logger
}

func NewHandler(storage *Storage, logger *logrus.Logger) *Handler {
	return &Handler{
		storage: storage,
		logger: logger,
	}
}

func (h *Handler) ListAll(w http.ResponseWriter, r *http.Request) {
	items, err := h.storage.ListAll(r.Context())
	if err != nil {
		h.logger.WithError(err).Error("Can't read data from redis")
		return
	}

	h.logger.Debugf("Data: %v", items)

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) ShareForm(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) Share(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		h.logger.WithError(err).Error("can't parse form")
		return
	}

	txt := r.PostFormValue("data")

	item := &sharedItem{
		Type: itemText,
		Value: txt,
	}

	err = h.storage.Add(r.Context(), item)
	if err != nil {
		h.logger.WithError(err).Error("can't insert data to redis")
		return
	}

	w.WriteHeader(http.StatusOK)
}
