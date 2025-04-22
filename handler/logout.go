package handler

import (
	"log"
	"net/http"
)

func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	session, err := h.session.Get(r, sessionName)
	if err != nil {
		log.Fatal(err)
	}

	authUserID := session.Values["authUserId"]
	if authUserID != "" {
		session.Values["authUserId"] = ""
	}

	if err := session.Save(r, w); err != nil {
		log.Fatal(err)
	}

	http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
}
