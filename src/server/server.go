package server

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"text/template"

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

	indexTemplateHTML, readFileErr := ioutil.ReadFile("templates/index.html")
	if readFileErr != nil {
		log.Fatal(readFileErr)
		return nil
	}

	indexTpl := template.Must(template.New("index_view").Parse(string(indexTemplateHTML)))

	// routes
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		data := map[string]interface{}{
			"foo": "bar",
		}

		render(w, r, indexTpl, "index_view", data)
	})

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// log.Printf("%+v", routerAPI)

	return r
}

func render(w http.ResponseWriter, r *http.Request, tpl *template.Template, name string, data interface{}) {
	buf := new(bytes.Buffer)
	if err := tpl.ExecuteTemplate(buf, name, data); err != nil {
		log.Fatalf("\nRender Error: %v\n", err)
		return
	}
	w.Write(buf.Bytes())
}

// Start ...
func Start(server *http.Server, port int) {
	log.Println("Server listening on " + strconv.Itoa(port))

	// returns ErrServerClosed on graceful close
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("ListenAndServe(): %s", err)
	}
}
