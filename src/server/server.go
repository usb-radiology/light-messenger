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

	errInsert := lmdatabase.ArduinoStatusInsert(db, status)
	if errInsert != nil {
		log.Fatal(errInsert)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Something bad happened!"))
		return errInsert
	}
	w.Write([]byte(fmt.Sprintf("%+v", status)))
	return nil
}

func mainHandler(config *configuration.Configuration, db *sql.DB, w http.ResponseWriter, r *http.Request) error {

	indexTpl := template.Must(template.ParseFiles("templates/index.html"))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	data := map[string]interface{}{}
	err := indexTpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return nil
}

func visierungHandler(config *configuration.Configuration, db *sql.DB, w http.ResponseWriter, r *http.Request) error {
	indexTpl := template.Must(template.ParseFiles("templates/visierung.html"))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	vars := mux.Vars(r)
	modality := vars["modality"]

	data := map[string]interface{}{
		"Modality": modality,
		"AOD":      create(db, modality, "aod"),
		"CTD":      create(db, modality, "ctd"),
		"MSK":      create(db, modality, "msk"),
		"NR":       create(db, modality, "nr"),
		"NUK_NUK":  create(db, "nuk", "NUK"),
	}
	err := indexTpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return nil
}

func confirmHandler(config *configuration.Configuration, db *sql.DB, w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("X-IC-Remove", "true")
	vars := mux.Vars(r)
	department := vars["department"]
	notificationID := vars["id"]
	log.Print("confirmHandler ", department, ", ", notificationID)
	lmdatabase.NotificationConfirm(db, notificationID, time.Now().Unix())
	return nil
}

func radiologieHandler(config *configuration.Configuration, db *sql.DB, w http.ResponseWriter, r *http.Request) error {

	indexTpl, _ := template.New("radiologie.html").ParseFiles("templates/radiologie.html")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	vars := mux.Vars(r)
	department := vars["department"]

	log.Print("radiologieHandler ", department)

	data := map[string]interface{}{
		"Department":    department,
		"Notifications": createNotificationTmpl(db, department),
	}
	err := indexTpl.Execute(w, data)
	if err != nil {
		log.Fatal(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return nil
}

func priorityHandler(config *configuration.Configuration, db *sql.DB, w http.ResponseWriter, r *http.Request) error {
	cardTemplateHTML, _ := ioutil.ReadFile("templates/card.html")
	cardTpl := template.Must(template.New("card_view").Parse(string(cardTemplateHTML)))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	vars := mux.Vars(r)
	modality := vars["modality"]
	department := vars["department"]
	priority := vars["priority"]
	log.Print("priorityHandler ", modality, ", ", department, ", ", priority)
	priorityMap := map[string]string{
		"1": "is-danger",
		"2": "is-warning",
		"3": "is-info",
	}
	priorityNumber, _ := strconv.Atoi(priority)

	notification, _ := lmdatabase.NotificationGetByDepartmentAndModality(db, department, modality)
	now := time.Now().Unix()
	if notification.NotificationID == "" {
		errInsertNotification := lmdatabase.NotificationInsert(db, department, priorityNumber, modality, now)
		if errInsertNotification != nil {
			log.Fatal(errInsertNotification)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Something bad happened!"))
			return errInsertNotification
		}
	} else {
		lmdatabase.NotificationUpdatePriority(db, notification.NotificationID, priorityNumber)
	}

	arduinoStatus, errInsert := lmdatabase.ArduinoStatusQueryWithin5MinutesFromNow(db, department, now)
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
		"CreatedAt":      time.Unix(now, 0).Format("15:04:05"),
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

	lmdatabase.NotificationCancel(db, modality, department, time.Now().Unix())

	if err := cardTpl.ExecuteTemplate(w, "card_view", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return nil
}

func getRouter(initConfig *configuration.Configuration) *mux.Router {
	r := mux.NewRouter()
	r.Handle("/", handler{initConfig, mainHandler})
	r.Handle("/mta/{modality}", handler{initConfig, visierungHandler})
	r.Handle("/radiologie/{department}", handler{initConfig, radiologieHandler})
	r.Handle("/nce-rest/arduino-status/{department}-status", handler{initConfig, arduinoStatusHandler})
	r.Handle("/notification/{department}/{id}", handler{initConfig, confirmHandler})
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
