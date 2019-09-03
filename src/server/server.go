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
	lmdatabase "github.com/usb-radiology/light-messenger/src/db"
)

type handler struct {
	*configuration.Configuration
	H func(config *configuration.Configuration, w http.ResponseWriter, r *http.Request) error
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h.H(h.Configuration, w, r)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
	}
}

// InitServer ...
func InitServer(initConfig *configuration.Configuration) *http.Server {
	port := strconv.Itoa(initConfig.Server.HTTPPort)
	r := getRouter(initConfig)
	server := &http.Server{Addr: ":" + port, Handler: r}
	return server
}

func arduinoStatusHandler(config *configuration.Configuration, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	department := vars["department"]
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	db, errDB := lmdatabase.GetDB(config)
	if errDB != nil {
		log.Fatal(errDB)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Something bad happened!"))
		return errDB
	}

	status := lmdatabase.ArduinoStatus{
		DepartmentID: department,
	}

	errInsert := lmdatabase.InsertStatus(db, status)
	if errInsert != nil {
		log.Fatal(errInsert)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Something bad happened!"))
		return errInsert
	}

	w.Write([]byte(department))
	return nil
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	indexTemplateHTML, _ := ioutil.ReadFile("templates/index.html")

	indexTpl := template.Must(template.New("index_view").Parse(string(indexTemplateHTML)))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	data := map[string]interface{}{
		"foo": "bar",
	}

	render(w, r, indexTpl, "index_view", data)
}

func getRouter(initConfig *configuration.Configuration) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", mainHandler)
	r.Handle("/nce-rest/arduino-status/{department}-status", handler{initConfig, arduinoStatusHandler})
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
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
