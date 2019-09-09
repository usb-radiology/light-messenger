package server

import (
	"database/sql"
	"net/http"

	"github.com/usb-radiology/light-messenger/src/configuration"
)

type handler struct {
	db           *sql.DB
	initConfig   *configuration.Configuration
	routeHandler func(config *configuration.Configuration, db *sql.DB, w http.ResponseWriter, r *http.Request) error
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := h.routeHandler(h.initConfig, h.db, w, r)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)

	}
}
