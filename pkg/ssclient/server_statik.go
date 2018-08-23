// +build statik

package ssclient

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rakyll/statik/fs"

	"fmt"

	_ "github.com/rubensayshi/xlspaceship/statik" // registers our static files to serve
)

func ServeAddStaticHandler(r *mux.Router) error {
	fmt.Printf("using statik to serve statics \n")

	// open statik bundle of our static files
	statikFS, err := fs.New()
	if err != nil {
		return err
	}

	// add file server for statik
	r.PathPrefix("/gui/").Handler(http.StripPrefix("/gui/", http.FileServer(statikFS)))

	return nil
}
