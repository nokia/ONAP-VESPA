package main

import (
	"net/http"
	"net/http/httputil"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func loggerMiddleware(hdl http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		log.Debugf("HTTP Request : %s %s", req.Method, req.URL.String())
		if *debug > 1 {
			dumped, err := httputil.DumpRequest(req, true)
			if err != nil {
				log.Error(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			log.Tracef("HTTP dump:\n%s\n\n", string(dumped))
			log.Trace("********* End of HTTP Request **********")
		}
		hdl.ServeHTTP(w, req)
	})
}

func initRoutes() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	if *debug > 0 {
		router.Use(loggerMiddleware)
	}

	router.Methods(http.MethodGet).
		Path("/").
		Handler(http.RedirectHandler("./doc/", http.StatusMovedPermanently))

	router.Methods(http.MethodGet).
		PathPrefix("/doc/").
		Handler(http.StripPrefix("/doc", http.FileServer(assets)))

	router.Methods(http.MethodGet).
		Path("/doc").
		Handler(http.RedirectHandler("./doc/", http.StatusMovedPermanently))

	router.Methods(http.MethodPost).
		Path("/api/eventListener/v5"+*topic).
		Headers("Content-Type", "application/json").
		Handler(errorWrapper(handlePostEvent))

	router.Methods(http.MethodPost).
		Path("/api/eventListener/v5/eventBatch").
		Headers("Content-Type", "application/json").
		Handler(errorWrapper(handlePostBatch))

	router.Methods(http.MethodGet).
		Path("/testControl/v5/events").
		Handler(errorWrapper(handleGetEvents))

	router.Methods(http.MethodGet).
		Path("/testControl/v5/stats").
		Handler(errorWrapper(handleGetStats))

	router.Methods(http.MethodDelete).
		Path("/testControl/v5/events").
		HandlerFunc(handleClearEvents)

	router.Methods(http.MethodPost).
		Path("/testControl/v5/commandList").
		Headers("Content-Type", "application/json").
		Handler(errorWrapper(handleSetCommandList))

	return router
}
