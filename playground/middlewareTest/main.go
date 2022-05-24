package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"time"
)

func main() {
	router := mux.NewRouter()
	//my example
	locationRouter := router.PathPrefix("/location/{locationId}").Subrouter()
	locationRouter.Use(locationMiddleware)

	cultureRouter := locationRouter.PathPrefix("/culture/{cultureId}").Subrouter()
	cultureRouter.Use(cultureMiddleware)

	cultureRouter.Handle("/yield", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Context().Value("locationId"))
		fmt.Println("handling yield")
	}))
	cultureRouter.Handle("/gdd", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { fmt.Println("handling gdd") }))

	microclimateRouter := locationRouter.PathPrefix("/microclimate/{microclimateId}").Subrouter()
	microclimateRouter.Use(microclimateMiddleware)
	microclimateRouter.Handle("", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { fmt.Println("handling microclimate") }))

	router.Handle("/location", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { fmt.Println("handling location") }))

	server := &http.Server{
		Addr: "0.0.0.0:8080",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router, // Pass our instance of gorilla/mux in.
	}

	fmt.Println("starting server")
	if err := server.ListenAndServe(); err != nil {
		fmt.Println(err)
	}

	//github example
	//router.Use(middleware1)
	//
	//wsRouter := router.PathPrefix("/ws").Subrouter()
	//wsRouter.Use(middleware2)
	//wsRouter.Use(middleware3)
	//
	//wsRouter.Handle("/sub/{subId}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//	fmt.Println("handling ws /sub")
	//	w.Write([]byte("/sub (ws)"))
	//}))
	//
	//chainRouter := router.PathPrefix("/chain").Subrouter()
	//chainRouter.HandleFunc("/sub1", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//	fmt.Println("handling chain /sub1")
	//	w.Write([]byte("/sub1"))
	//}))
	//chainRouter.HandleFunc("/sub2", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//	fmt.Println("handling chain /sub2")
	//	w.Write([]byte("/sub2"))
	//}))
	//
	//restRouter := router.PathPrefix("/").Subrouter()
	//restRouter.Use(middleware3)
	//
	//restRouter.HandleFunc("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//	fmt.Println("handling rest /")
	//	w.Write([]byte("/ (rest)"))
	//}))
	//
	//server := &http.Server{
	//	Addr: "0.0.0.0:8080",
	//	// Good practice to set timeouts to avoid Slowloris attacks.
	//	WriteTimeout: time.Second * 15,
	//	ReadTimeout:  time.Second * 15,
	//	IdleTimeout:  time.Second * 60,
	//	Handler:      router, // Pass our instance of gorilla/mux in.
	//}
	//
	//fmt.Println("starting server")
	//if err := server.ListenAndServe(); err != nil {
	//	fmt.Println(err)
	//}
}

//func middleware1(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		fmt.Println("middleware1")
//		next.ServeHTTP(w, r)
//	})
//}
//
//func middleware2(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		fmt.Println("middleware2")
//		next.ServeHTTP(w, r)
//	})
//}
//
//func middleware3(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		fmt.Println("middleware3")
//		next.ServeHTTP(w, r)
//	})
//}
func locationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		cultureId, _ := strconv.ParseUint(params["locationId"], 10, 32)
		ctx := context.WithValue(r.Context(), "locationId", cultureId)
		r = r.WithContext(ctx)
		fmt.Println("location")
		next.ServeHTTP(w, r)
	})
}

func cultureMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("culture")
		next.ServeHTTP(w, r)
	})
}

func microclimateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("microclimate")
		next.ServeHTTP(w, r)
	})
}
