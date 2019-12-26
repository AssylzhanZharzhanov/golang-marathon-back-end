package main

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"io/ioutil"
// 	"log"
// 	// "marathon/controller"
// 	"net/http"
// 	"os"
// )
import (
	"log"
	"marathon/controller"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	STATIC_DIR = "/static/"
	PORT       = "8080"
)

func main() {
	// httpServer := echo.New()
	// httpServer.GET("/image/:id", streamImage)

	r := mux.NewRouter()

	r.HandleFunc("api/v1/signup", controller.RegistrationHandler).
		Methods("POST")
	r.HandleFunc("api/v1/login", controller.LoginHandler).
		Methods("POST")
	r.HandleFunc("/api/v1/search", controller.SearchImageHandler).
		Methods("GET")
	r.HandleFunc("/api/v1/images", controller.GetImagesHandler).
		Methods("GET")
	r.HandleFunc("/api/v1/marathons", controller.GetMarathonsHandler).
		Methods("GET")
	r.HandleFunc("/api/v1/marathons", controller.AddMarathonHandler).
		Methods("POST")
	r.HandleFunc("/api/v1/image/{id}", controller.SearchImageHandler).
		Methods("GET")
	r.HandleFunc("/api/v1/images", controller.UploadManyFiles).
		Methods("POST")

	log.Println("Server is started...")
	http.ListenAndServe(":5000", r)
}

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	// Server CSS, JS & Images Statically.
	router.
		PathPrefix(STATIC_DIR).
		Handler(http.StripPrefix(STATIC_DIR, http.FileServer(http.Dir("."+STATIC_DIR))))

	return router
}
