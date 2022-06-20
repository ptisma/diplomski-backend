package middlewares

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func LocationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		locationId, err := strconv.ParseUint(params["locationId"], 10, 32)
		if err != nil {
			w.WriteHeader(http.StatusPreconditionFailed)
			fmt.Fprintf(w, "Error in parsing: locationId is not uint")
			return
		}
		ctx := context.WithValue(r.Context(), "locationId", locationId)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func CultureMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		cultureId, err := strconv.ParseUint(params["cultureId"], 10, 32)
		if err != nil {
			w.WriteHeader(http.StatusPreconditionFailed)
			fmt.Fprintf(w, "Error in parsing: cultureId is not uint")
			return
		}
		ctx := context.WithValue(r.Context(), "cultureId", cultureId)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func MicroclimateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		microclimateId, err := strconv.ParseUint(params["microclimateId"], 10, 32)
		if err != nil {
			w.WriteHeader(http.StatusPreconditionFailed)
			fmt.Fprintf(w, "Error in parsing: microclimateId is not uint")
			return
		}
		ctx := context.WithValue(r.Context(), "microclimateId", microclimateId)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
