package middlewares

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"time"
)

// Middleware to check for location id
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

// Middleware to check for culture id
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

// Middleware to check for microclimate id
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

//func Top(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		fmt.Println("top")
//		next.ServeHTTP(w, r)
//	})
//}
//
//func Kek(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		fmt.Println("kek")
//		next.ServeHTTP(w, r)
//	})
//}
//func Zozzle(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		urlParams := r.URL.Query()
//		//parse direct into string YYYY-MM-DD
//		fromDate, _ := time.Parse("20060102", urlParams.Get("from"))
//
//		toDate, _ := time.Parse("20060102", urlParams.Get("to"))
//
//		fmt.Println(fromDate, toDate)
//
//		next.ServeHTTP(w, r)
//	})
//}

// Middleware to check for from and to query params
func DatesMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		urlParams := r.URL.Query()
		//parse direct into string YYYY-MM-DD
		fromDate, err := time.Parse("20060102", urlParams.Get("from"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error in parsing: fromDate is not in YYYYMMDD format")
			return
		}
		toDate, err := time.Parse("20060102", urlParams.Get("to"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error in parsing: toDate is not in YYYYMMDD format")
			return
		}
		difference := toDate.Sub(fromDate)
		if int(difference.Hours()/24) < 1 {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error in parsing: toDate is lower than fromDate")
			return
		}
		ctx := context.WithValue(r.Context(), "from", fromDate)
		ctx = context.WithValue(ctx, "to", toDate)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
