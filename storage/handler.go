package storage

import (
	"github.com/sirupsen/logrus"
	"html/template"
	"net/http"
)

type Handler struct {
	storage *Storage
	logger  *logrus.Logger
}

var homeTmpl = template.Must(template.ParseFiles("./views/home.html"))
var shareTmpl = template.Must(template.ParseFiles("./views/share.html"))

func NewHandler(storage *Storage, logger *logrus.Logger) *Handler {
	return &Handler{
		storage: storage,
		logger:  logger,
	}
}

func (h *Handler) ListAll(w http.ResponseWriter, r *http.Request) {
	items, err := h.storage.ListAll(r.Context())
	if err != nil {
		h.logger.WithError(err).Error("Can't read data from redis")
		return
	}

	data := map[string]interface{}{
		"Items": items,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	err = homeTmpl.Execute(w, data)
	if err != nil {
		h.logger.WithError(err).Error(err)
	}
}

func (h *Handler) ShareForm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/html; charset=utf-8")
	err := shareTmpl.Execute(w, nil)
	if err != nil {
		h.logger.WithError(err).Error(err)
	}
}

func (h *Handler) Share(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		h.logger.WithError(err).Error("can't parse form")
		return
	}

	txt := r.Form.Get("data")

	item := &sharedItem{
		Type:  itemText,
		Value: txt,
	}

	err = h.storage.Add(r.Context(), item)
	if err != nil {
		h.logger.WithError(err).Error("can't insert data to redis")
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
