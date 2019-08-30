package server

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/usb-radiology/light-messenger/src/configuration"
)

// InitServer ...
func InitServer(initConfig *configuration.Configuration) *http.Server {

	port := strconv.Itoa(initConfig.Server.HTTPPort)

	r := getRouter(initConfig)

	server := &http.Server{Addr: ":" + port, Handler: r}

	// log.Println("Server listening on " + port)
	// log.Fatal(http.ListenAndServe(":"+port, nil))
	// log.Fatal(http.ListenAndServe(":"+port, r))

	return server
}

// getRouter ...
func getRouter(initConfig *configuration.Configuration) *mux.Router {

	r := mux.NewRouter()

	// r.Use(loggingMiddleware)

	// api
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		response := "hello world"
		w.Write([]byte(response))
	})

	// fs := http.FileServer(http.Dir("public"))
	// r.PathPrefix("/").Handler(fs)
	// r.Handle("/", fs)

	// log.Printf("%+v", routerAPI)

	return r
}

// Start ...
func Start(server *http.Server, port int) {
	log.Println("Server listening on " + strconv.Itoa(port))

	// returns ErrServerClosed on graceful close
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("ListenAndServe(): %s", err)
	}
}
