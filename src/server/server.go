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
	"github.com/pkg/errors"
	"github.com/usb-radiology/light-messenger/src/configuration"
	"github.com/usb-radiology/light-messenger/src/lmdatabase"
	"github.com/usb-radiology/light-messenger/src/version"
)

var templates = make(map[string]*template.Template)
var templateIndexID = "index"
var templateCardID = "card"
var templateRadiologieID = "radiologie"
var templateVisierungID = "visierung"

func compileTemplates() {
	{
		indexTpl := template.Must(template.ParseFiles("templates/index.html"))
		templates[templateIndexID] = indexTpl
	}

	{
		cardTemplateHTML, _ := ioutil.ReadFile("templates/card.html")
		cardTpl := template.Must(template.New("card_view").Parse(string(cardTemplateHTML)))
		templates[templateCardID] = cardTpl
	}

	{
		radiologieTpl, _ := template.New("radiologie.html").ParseFiles("templates/radiologie.html")
		templates[templateRadiologieID] = radiologieTpl
	}

	{
		visierungTpl := template.Must(template.ParseFiles("templates/visierung.html"))
		templates[templateVisierungID] = visierungTpl
	}
}

var priorityMap = map[int]string{
	1: "is-danger",
	2: "is-warning",
	3: "is-info",
}

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

	status := lmdatabase.ArduinoStatus{
		DepartmentID: department,
		StatusAt:     time.Now().Unix(),
	}

	errInsert := lmdatabase.ArduinoStatusInsert(db, status)
	if writeInternalServerError(errInsert, w) != nil {
		return errInsert
	}

	w.Write([]byte(fmt.Sprintf("%+v", status)))
	return nil
}

func mainHandler(config *configuration.Configuration, db *sql.DB, w http.ResponseWriter, r *http.Request) error {

	data := map[string]interface{}{
		"Version":   version.Version,
		"BuildTime": version.BuildTime,
	}

	err := renderTemplate(w, r, templates[templateIndexID], data)
	if writeInternalServerError(err, w) != nil {
		return err
	}

	return nil
}

func visierungHandler(config *configuration.Configuration, db *sql.DB, w http.ResponseWriter, r *http.Request) error {

	vars := mux.Vars(r)
	modality := vars["modality"]

	data := map[string]interface{}{
		"Modality":  modality,
		"AOD":       getCardHTML(db, modality, "aod"),
		"CTD":       getCardHTML(db, modality, "ctd"),
		"MSK":       getCardHTML(db, modality, "msk"),
		"NR":        getCardHTML(db, modality, "nr"),
		"NUK_NUK":   getCardHTML(db, "nuk", "NUK"),
		"Version":   version.Version,
		"BuildTime": version.BuildTime,
	}

	err := renderTemplate(w, r, templates[templateVisierungID], data)
	if writeInternalServerError(err, w) != nil {
		return err
	}

	return nil
}

func confirmHandler(config *configuration.Configuration, db *sql.DB, w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("X-IC-Remove", "true")
	vars := mux.Vars(r)
	department := vars["department"]
	notificationID := vars["id"]
	log.Print("confirmHandler ", department, ", ", notificationID)
	lmdatabase.NotificationConfirm(db, notificationID, time.Now().Unix())
	return nil
}

func radiologieHandler(config *configuration.Configuration, db *sql.DB, w http.ResponseWriter, r *http.Request) error {

	vars := mux.Vars(r)
	department := vars["department"]

	log.Print("radiologieHandler ", department)

	data := map[string]interface{}{
		"Department":    department,
		"Notifications": createNotificationTmpl(db, department),
		"Version":       version.Version,
		"BuildTime":     version.BuildTime,
	}

	err := renderTemplate(w, r, templates[templateRadiologieID], data)
	if writeInternalServerError(err, w) != nil {
		return err
	}

	return nil
}

func priorityHandler(config *configuration.Configuration, db *sql.DB, w http.ResponseWriter, r *http.Request) error {

	vars := mux.Vars(r)
	modality := vars["modality"]
	department := vars["department"]
	priority := vars["priority"]
	log.Print("priorityHandler ", modality, ", ", department, ", ", priority)

	priorityNumber, errPriorityType := strconv.Atoi(priority)
	if writeInternalServerError(errPriorityType, w) != nil {
		return errPriorityType
	}

	notification, _ := lmdatabase.NotificationGetByDepartmentAndModality(db, department, modality)
	now := time.Now().Unix()
	if notification.NotificationID == "" {
		errInsertNotification := lmdatabase.NotificationInsert(db, department, priorityNumber, modality, now)
		if writeInternalServerError(errInsertNotification, w) != nil {
			return errPriorityType
		}
	} else {
		lmdatabase.NotificationUpdatePriority(db, notification.NotificationID, priorityNumber)
	}

	arduinoStatus, errInsert := lmdatabase.ArduinoStatusQueryWithin5MinutesFromNow(db, department, now)
	if writeInternalServerError(errInsert, w) != nil {
		return errPriorityType
	}

	data := map[string]interface{}{
		"Modality":       modality,
		"Department":     department,
		"Priority":       priority,
		"PriorityName":   priorityMap[priorityNumber],
		"PriorityNumber": priorityNumber,
		"ArduinoStatus":  arduinoStatus,
		"CreatedAt":      time.Unix(now, 0).Format("15:04:05"),
	}

	err := renderTemplateName(w, r, templates[templateCardID], "card_view", data)
	if writeInternalServerError(err, w) != nil {
		return err
	}

	return nil
}

func cancelHandler(config *configuration.Configuration, db *sql.DB, w http.ResponseWriter, r *http.Request) error {

	vars := mux.Vars(r)
	modality := vars["modality"]
	department := vars["department"]

	data := map[string]interface{}{
		"Modality":       modality,
		"Department":     department,
		"PriorityNumber": 99, // needed because of le comparison in template
	}

	lmdatabase.NotificationCancel(db, modality, department, time.Now().Unix())

	err := renderTemplateName(w, r, templates[templateCardID], "card_view", data)
	if writeInternalServerError(err, w) != nil {
		return err
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

// Start ...
func Start(server *http.Server, port int) {
	compileTemplates()

	log.Println("Server listening on " + strconv.Itoa(port))

	// returns ErrServerClosed on graceful close
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("ListenAndServe(): %s", err)
	}
}

func renderTemplateName(w http.ResponseWriter, r *http.Request, tpl *template.Template, name string, data interface{}) error {
	buf := new(bytes.Buffer)
	err := tpl.ExecuteTemplate(buf, name, data)
	if writeInternalServerError(err, w) != nil {
		return err
	}
	w.Write(buf.Bytes())
	return nil
}

func renderTemplate(w http.ResponseWriter, r *http.Request, tpl *template.Template, data interface{}) error {
	err := tpl.Execute(w, data)
	if writeInternalServerError(err, w) != nil {
		return err
	}
	return nil
}

func writeInternalServerError(err error, w http.ResponseWriter) error {
	if err != nil {
		log.Fatalf("%+v", errors.WithStack(err))
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Something bad happened!"))
		return err
	}

	return nil
}
