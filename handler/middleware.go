package handler

import (
	"log"
	"net/http"
)

func (h *Handler) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := h.session.Get(r, sessionName)
		if err != nil {
			log.Fatal(err)
		}
		authUserID := session.Values["authUserId"]
		if authUserID != "" {
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "Unauthorized Access", http.StatusUnauthorized)
		}
	})
}

func (h *Handler) restrictMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := h.session.Get(r, sessionName)
		if err != nil {
			log.Fatal(err)
		}
		authUserID := session.Values["authUserId"]
		if authUserID != "" {
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
