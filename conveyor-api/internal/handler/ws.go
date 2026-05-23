package handler

import "net/http"

func (h *Handler) ServeWS(w http.ResponseWriter, r *http.Request) {
	h.Hub.ServeWS(w, r)
}
