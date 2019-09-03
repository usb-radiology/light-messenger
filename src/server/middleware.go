package server

import (
	"database/sql"
	"net/http"

	"github.com/usb-radiology/light-messenger/src/configuration"
	"github.com/usb-radiology/light-messenger/src/lmdatabase"
)

type handler struct {
	initConfig   *configuration.Configuration
	routeHandler func(config *configuration.Configuration, db *sql.DB, w http.ResponseWriter, r *http.Request) error
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	db, errDb := lmdatabase.GetDB(h.initConfig)
	if errDb != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	err := h.routeHandler(h.initConfig, db, w, r)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)

	}
}
