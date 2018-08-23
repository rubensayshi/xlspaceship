// +build !statik

package ssclient

import (
	"net/http"

	"fmt"

	"github.com/gorilla/mux"
)

func ServeAddStaticHandler(r *mux.Router) error {
	fmt.Printf("using fs to serve statics \n")

	// add file server for statik
	r.PathPrefix("/gui/").Handler(http.StripPrefix("/gui/", http.FileServer(http.Dir("./gui/web/www"))))

	return nil
}
