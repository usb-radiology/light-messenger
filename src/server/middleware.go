package server

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/usb-radiology/light-messenger/src/configuration"
)

type handler struct {
	db           *sql.DB
	initConfig   *configuration.Configuration
	routeHandler func(config *configuration.Configuration, db *sql.DB, w http.ResponseWriter, r *http.Request) error
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h.routeHandler(h.initConfig, h.db, w, r)
	if err != nil {
		log.Printf("%+v", err)

		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
	}
}
