// This application demonstrates how to use the new routing capabilities of ServeMux in Go v1.22 and
// how to use a stdlib compatible third-party package, such as Chi, to provide middleware services.
// A basic comparison of processing times is made between these alternatives.

// COMPARISON OF THESE TWO OPTIONS

// STDLIB
//2024/05/14 21:45:41 "GET http://localhost:8090/task/f0cd2e/ HTTP/1.1" from [::1]:65380 - 200 34B in 39.187µs
//2024/05/14 21:45:41 "GET http://localhost:8090/task/f0cd2e/ HTTP/1.1" from [::1]:65381 - 200 34B in 19.451µs
//2024/05/14 21:45:42 "GET http://localhost:8090/task/f0cd2e/ HTTP/1.1" from [::1]:65382 - 200 34B in 14.573µs
// CHI
//2024/05/14 21:45:44 "GET http://localhost:8091/task/f0cd2e/ HTTP/1.1" from 127.0.0.1:65388 - 200 34B in 3.626µs
//2024/05/14 21:45:45 "GET http://localhost:8091/task/f0cd2e/ HTTP/1.1" from 127.0.0.1:65390 - 200 34B in 2.813µs
//2024/05/14 21:45:45 "GET http://localhost:8091/task/f0cd2e/ HTTP/1.1" from 127.0.0.1:65392 - 200 34B in 4.059µs

// STDLIB
//2024/05/14 21:45:50 "GET http://localhost:8090/path/ HTTP/1.1" from [::1]:65393 - 200 14B in 4.848µs
//2024/05/14 21:45:50 "GET http://localhost:8090/path/ HTTP/1.1" from [::1]:65396 - 200 14B in 4.697µs
//2024/05/14 21:45:51 "GET http://localhost:8090/path/ HTTP/1.1" from [::1]:65399 - 200 14B in 4.952µs
// CHI
//2024/05/14 21:45:54 "GET http://localhost:8091/path/ HTTP/1.1" from 127.0.0.1:65402 - 200 14B in 3.274µs
//2024/05/14 21:45:54 "GET http://localhost:8091/path/ HTTP/1.1" from 127.0.0.1:65404 - 200 14B in 2.581µs
//2024/05/14 21:45:55 "GET http://localhost:8091/path/ HTTP/1.1" from 127.0.0.1:65406 - 200 14B in 2.488µs

package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	go useChi(&wg)
	go useStdLibWithChiMiddleware(&wg)
	wg.Wait()

}

// Use Chi as the HTTP router and Chi's logging middleware.
func useChi(wg *sync.WaitGroup) {
	defer wg.Done()

	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Get("/path/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, "Chi: Hit GET path endpoint\n")
	})

	r.Get("/task/{id}/", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		_, _ = fmt.Fprintf(w, "Chi: Hit GET task by id={%v} endpoint\n", id)
	})

	_ = http.ListenAndServe(":8091", r)
}

// Use Go v1.22 ServerMux as the HTTP router and Chi's logging middleware.
func useStdLibWithChiMiddleware(wg *sync.WaitGroup) {
	defer wg.Done()

	r := http.NewServeMux()

	pathHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, "StdLib: Hit GET path endpoint\n")
	})
	r.Handle("GET /path/", middleware.Logger(pathHandler))

	taskHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		_, _ = fmt.Fprintf(w, "StdLib: Hit GET task by id={%v} endpoint\n", id)
	})
	r.Handle("GET /task/{id}/", middleware.Logger(taskHandler))

	_ = http.ListenAndServe("localhost:8090", r)
}
