package routes

import (
	"github.com/gorilla/mux"
	"go-backend-technical-test/controllers"
	"go-backend-technical-test/middlewares"
	"net/http"
)

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000") // Specify allowed origin
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept, X-Requested-With")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func AppRoutes() *mux.Router {
	r := mux.NewRouter()

	r.Use(CORS)
	
	r.HandleFunc("/register", controllers.Register).Methods(http.MethodPost)
	r.HandleFunc("/login", controllers.Login).Methods(http.MethodPost)

	protected := r.PathPrefix("/api").Subrouter()
	protected.Use(middlewares.Authenticate)
	protected.HandleFunc("/events", controllers.GetEvents).Methods(http.MethodGet)
	protected.HandleFunc("/events/{id}", controllers.GetEventByID).Methods(http.MethodGet)
	protected.HandleFunc("/events", controllers.CreateEvent).Methods(http.MethodPost)
	protected.HandleFunc("/events/{id}", controllers.EditEvent).Methods(http.MethodPut)
	protected.HandleFunc("/events/{id}", controllers.DeleteEvent).Methods(http.MethodDelete)

	return r
}

