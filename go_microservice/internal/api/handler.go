package api

import (
	"log"
	"net/http"
)

type HttpHandler struct {
	emailUseCase MessageSender
}

func NewHttpHandler(emailUseCase MessageSender) *HttpHandler {
	return &HttpHandler{emailUseCase: emailUseCase}
}

func (h *HttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/verify":
		h.verifyHandler(w, r)
		log.Printf("'/verify' url triggered")
	default:
		http.NotFound(w, r)
	}
}

func (h *HttpHandler) verifyHandler(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	email, err := h.emailUseCase.SendVerifyMessageToRMQ(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Email " + email + " verified successfully"))
}
