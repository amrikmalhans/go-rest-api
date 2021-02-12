package routes  
 
import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"restapi/controllers"
	"restapi/helpers"
)

// Routes ... routes func for all the routes
func Routes() {
	r := mux.NewRouter()
	corsWrapper := cors.New(cors.Options{
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Origin", "Accept", "*"},
		AllowCredentials: true,
		AllowedOrigins:   []string{"http://localhost:4000"},
	})
	r.HandleFunc("/api/books", controllers.GetBooks).Methods("GET")
	r.HandleFunc("/api/books/{id}", controllers.GetBook).Methods("GET")
	r.Handle("/api/books", helpers.Middleware(controllers.CreateBook)).Methods("POST")
	r.HandleFunc("/api/books/{id}", controllers.UpdateBook).Methods("PUT")
	r.HandleFunc("/api/books/{id}", controllers.DeleteBook).Methods("DELETE")
	r.HandleFunc("/signup", controllers.Signup).Methods("POST")
	r.HandleFunc("/login", controllers.Login).Methods("POST", "OPTIONS")
	r.HandleFunc("/tokens/refresh", controllers.Refresh).Methods("POST", "OPTIONS")
	r.HandleFunc("/collection", controllers.CreateCollection).Methods("POST", "OPTIONS")
	log.Fatal(http.ListenAndServe(":9000", corsWrapper.Handler(r)))
}
