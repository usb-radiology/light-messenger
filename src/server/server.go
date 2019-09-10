package server

import (
	"bytes"
	"log"
	"net/http"
	"strconv"
	"text/template"
	"time"

	rice "github.com/GeertJohan/go.rice"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/usb-radiology/light-messenger/src/configuration"
	"github.com/usb-radiology/light-messenger/src/lmdatabase"
)

// globals ...
var (
	templates            = make(map[string]*template.Template)
	templateIndexID      = "index"
	templateCardID       = "card"
	templateRadiologieID = "radiologie"
	templateVisierungID  = "visierung"
	box                  = rice.MustFindBox("../../static")
	priorityMap          = map[int]string{
		1: "is-danger",
		2: "is-warning",
		3: "is-info",
	}
)

// InitServer ...
func InitServer(initConfig *configuration.Configuration) *http.Server {
	port := strconv.Itoa(initConfig.Server.HTTPPort)
	r := getRouter(initConfig)
	server := &http.Server{Addr: ":" + port, Handler: r}
	return server
}

// Start ...
func Start(server *http.Server, port int) {
	log.Println("Server listening on " + strconv.Itoa(port))

	// returns ErrServerClosed on graceful close
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("ListenAndServe(): %s", err)
	}
}

func getRouter(initConfig *configuration.Configuration) *mux.Router {

	db, errDb := lmdatabase.GetDB(initConfig)
	if errDb != nil {
		log.Fatalf("%+v", errors.WithStack(errDb))
		return nil
	}

	errCompileTemplates := compileTemplates()
	if errCompileTemplates != nil {
		log.Fatalf("%+v", errCompileTemplates)
		return nil
	}

	r := mux.NewRouter()
	r.Handle("/", handler{db, initConfig, mainHandler})
	r.Handle("/mtra/{modality}", handler{db, initConfig, visierungHandler})
	r.Handle("/radiologie/{department}", handler{db, initConfig, radiologieHandler})
	r.Handle("/nce-rest/arduino-status/{department}-status", handler{db, initConfig, arduinoStatusHandler})
	r.Handle("/nce-rest/arduino-status/{department}-open-notifications", handler{db, initConfig, openStatusHandler})
	r.Handle("/notification/{department}/{id}", handler{db, initConfig, confirmHandler})
	r.Handle("/modality/{modality}/department/{department}/prio/{priority}", handler{db, initConfig, priorityHandler})
	r.Handle("/modality/{modality}/department/{department}/cancel", handler{db, initConfig, cancelHandler})

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(box.HTTPBox())))
	return r
}

func compileTemplates() error {
	{
		templateString, err := box.String("templates/index.html")
		if err != nil {
			return err
		}
		indexTpl := template.Must(template.New("index").Parse(templateString))
		templates[templateIndexID] = indexTpl
	}

	{
		cardTemplateHTML, _ := box.String("templates/card.html")
		cardTpl := template.Must(template.New("card_view").Parse(string(cardTemplateHTML)))
		templates[templateCardID] = cardTpl
	}

	{
		templateString, err := box.String("templates/radiologie.html")
		if err != nil {
			return err
		}
		radiologieTpl, _ := template.New("radiologie").Parse(templateString)
		templates[templateRadiologieID] = radiologieTpl
	}

	{
		funcMap := template.FuncMap{
			"priorityMap": func(prio int) string {
				priorityMap := map[int]string{
					1: "is-danger",
					2: "is-warning",
					3: "is-info",
				}
				return priorityMap[prio]
			},
			"priorityName": func(prio int) string {
				priorityMap := map[int]string{
					1: "Hoch",
					2: "Mittel",
					3: "Tief",
				}
				return priorityMap[prio]
			},
			"toTime": func(now int64) string {
				if now == -1 {
					return ""
				}
				return time.Unix(now, 0).Format("2006-01-02 15:04:05")
			},
		}

		templateString, err := box.String("templates/visierung.html")
		if err != nil {
			return err
		}

		visierungTpl := template.Must(template.New("visierung.html").Funcs(funcMap).Parse(templateString))
		templates[templateVisierungID] = visierungTpl
	}
	return nil
}

func renderTemplateName(w http.ResponseWriter, r *http.Request, tpl *template.Template, name string, data interface{}) error {
	buf := new(bytes.Buffer)
	err := tpl.ExecuteTemplate(buf, name, data)
	if err != nil {
		return err
	}

	errWrite := writeBytes(w, buf.Bytes())
	if errWrite != nil {
		return errWrite
	}

	return nil
}

func renderTemplate(w http.ResponseWriter, r *http.Request, tpl *template.Template, data interface{}) error {
	err := tpl.Execute(w, data)
	if writeInternalServerError(err, w) != nil {
		return err
	}
	return nil
}

func writeBytes(w http.ResponseWriter, bytes []byte) error {
	_, errWrite := w.Write(bytes)
	if errWrite != nil {
		return errWrite
	}

	return nil
}

func writeInternalServerError(err error, w http.ResponseWriter) error {
	if err != nil {
		// we don't necessarily want to kill the application on 500
		log.Printf("%+v", errors.WithStack(err))
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	return nil
}
