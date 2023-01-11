package handlers

import (
	"context"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"sf-api-gateway/internal/constants"
)

func (h *Handler) Middlware(handle httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		ctx := r.Context()
		requestId := uuid.New().String()
		ctx = context.WithValue(ctx, constants.RequestIdKey, requestId)
		r = r.WithContext(ctx)

		w.Header().Add(constants.RequestIdKey, requestId)
		w.Header().Add("Content-Type", "application/json")

		handle(w, r, ps)
	}
}

func (h *Handler) Cors(handle httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Credentials", "true")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

		if r.Method == "OPTIONS" {
			http.Error(w, "No Content", http.StatusNoContent)
			return
		}

		handle(w, r, ps)
	}
}
