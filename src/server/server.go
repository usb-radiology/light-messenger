package server

import (
	"bytes"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/gorilla/mux"
	"github.com/usb-radiology/light-messenger/src/configuration"
	"github.com/usb-radiology/light-messenger/src/lmdatabase"
)

// InitServer ...
func InitServer(initConfig *configuration.Configuration) *http.Server {
	port := strconv.Itoa(initConfig.Server.HTTPPort)
	r := getRouter(initConfig)
	server := &http.Server{Addr: ":" + port, Handler: r}
	return server
}

func arduinoStatusHandler(config *configuration.Configuration, db *sql.DB, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	department := vars["department"]
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	status := lmdatabase.ArduinoStatus{
		DepartmentID: department,
		StatusAt:     time.Now().Unix(),
	}

	errInsert := lmdatabase.InsertStatus(db, status)
	if errInsert != nil {
		log.Fatal(errInsert)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Something bad happened!"))
		return errInsert
	}

	w.Write([]byte(fmt.Sprintf("%+v", status)))
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

func priorityHandler(config *configuration.Configuration, db *sql.DB, w http.ResponseWriter, r *http.Request) error {
	cardTemplateHTML, _ := ioutil.ReadFile("templates/card.html")
	cardTpl := template.Must(template.New("card_view").Parse(string(cardTemplateHTML)))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	vars := mux.Vars(r)
	modality := vars["modality"]
	department := vars["department"]
	priority := vars["priority"]

	priorityMap := map[string]string{
		"1": "is-danger",
		"2": "is-warning",
		"3": "is-info",
	}
	priorityNumber, _ := strconv.Atoi(priority)

	arduinoStatus, errInsert := lmdatabase.IsAlive(db, department, time.Now().Unix())
	if errInsert != nil {
		log.Fatal(errInsert)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Something bad happened!"))
		return errInsert
	}

	data := map[string]interface{}{
		"Modality":       modality,
		"Department":     department,
		"Priority":       priority,
		"PriorityName":   priorityMap[priority],
		"PriorityNumber": priorityNumber,
		"ArduinoStatus":  arduinoStatus,
	}

	if err := cardTpl.ExecuteTemplate(w, "card_view", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return nil
}

func cancelHandler(config *configuration.Configuration, db *sql.DB, w http.ResponseWriter, r *http.Request) error {
	cardTemplateHTML, _ := ioutil.ReadFile("templates/card.html")
	cardTpl := template.Must(template.New("card_view").Parse(string(cardTemplateHTML)))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	vars := mux.Vars(r)
	modality := vars["modality"]
	department := vars["department"]

	data := map[string]interface{}{
		"Modality":       modality,
		"Department":     department,
		"PriorityNumber": 99, // needed because of le comparison in template
	}

	if err := cardTpl.ExecuteTemplate(w, "card_view", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return nil
}

func getRouter(initConfig *configuration.Configuration) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", mainHandler)
	r.Handle("/nce-rest/arduino-status/{department}-status", handler{initConfig, arduinoStatusHandler})
	r.Handle("/modality/{modality}/department/{department}/prio/{priority}", handler{initConfig, priorityHandler})
	r.Handle("/modality/{modality}/department/{department}/cancel", handler{initConfig, cancelHandler})
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
