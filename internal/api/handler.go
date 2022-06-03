package api

import (
	"encoding/xml"
	"net/http"

	"github.com/Kichiyaki/lubimyczytacrss/internal/lubimyczytac"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	client *lubimyczytac.Client
}

func NewHandler(client *lubimyczytac.Client) *Handler {
	return &Handler{client: client}
}

func (h *Handler) Register(r chi.Router) {
	r.Get("/api/v1/rss/authors/{authorID}", h.getRSSAuthor)
}

func (h *Handler) getRSSAuthor(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	author, err := h.client.GetAuthor(ctx, chi.URLParamFromCtx(ctx, "authorID"))
	if err == lubimyczytac.ErrAuthorNotFound {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`author not found`))
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`something went wrong while getting author info: ` + err.Error()))
		return
	}

	w.Header().Set("Content-Type", "text/xml; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_ = xml.NewEncoder(w).Encode(rssMainFromAuthor(author))
}
