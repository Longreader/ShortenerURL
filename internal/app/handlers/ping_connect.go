package handlers

import "net/http"

func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) {

	ok, err := h.st.Ping(r.Context())

	if err != nil {
		http.Error(w, "Connection error", http.StatusInternalServerError)
		return
	}

	if !ok {
		http.Error(w, "Connection error", http.StatusInternalServerError)
		return
	} else {
		w.WriteHeader(http.StatusOK)
		return
	}
}
