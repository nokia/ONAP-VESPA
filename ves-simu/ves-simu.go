/*
	Copyright 2019 Nokia

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/

package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Simulator parameters
var (
	serverRoot    = flag.String("server-root", "", "The path before the /eventListener part of the POST URL")
	user          = flag.String("user", "user", "The username for authenticating incoming requests")
	pass          = flag.String("passwd", "pass", "The password for authenticating incoming requests")
	port          = flag.Int("port", 8443, "The port to bind VES simulator to")
	topic         = flag.String("topic", "", "An optional topic")
	https         = flag.Bool("https", false, "Enable HTTPS instead of plain HTTP")
	certFile      = flag.String("cert", "", "Path to certificate file if HTTPS is enabled")
	keyFile       = flag.String("key", "", "Path to certificate's private key file if HTTPS is enabled")
	maxEventsKeep = flag.Int("max-event", 4096, "Maximum number of events to keep in memory, 0 meaning there's no limit")
	eventMaxSize  = flag.Int("event-size", 2000000, "Maximum size in bytes allowed for events")
	logfile       = flag.String("logfile", "", "Path to a file where to write logs outputs, additionally to standard output")
	jsonlog       = flag.Bool("jsonlog", false, "Format log messaegs into JSON")
	debug         = flag.Uint("d", 0, "Verbosity level. 0, 1 or 2")
)

func initialize() {
	flag.Parse()
	var formatter log.Formatter = &log.TextFormatter{
		FullTimestamp: true,
	}
	if *jsonlog {
		formatter = &log.JSONFormatter{}
	}
	log.SetFormatter(formatter)
	if *logfile != "" {
		lf, err := os.OpenFile(*logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		log.SetOutput(io.MultiWriter(os.Stdout, lf))
	}
	switch *debug {
	case 0:
		log.SetLevel(log.InfoLevel)
	case 1:
		log.SetLevel(log.DebugLevel)
	case 2:
		log.SetLevel(log.TraceLevel)
	default:
		log.Fatalf("Unsupported verbosity level %d", *debug)
	}
}

func main() {
	initialize()
	log.Infof("Starting VES simulator (PID=%d)", os.Getpid())
	*topic = strings.TrimLeft(*topic, "/")
	if *topic != "" && !strings.HasPrefix(*topic, "/") {
		*topic = "/" + *topic
	}
	*serverRoot = strings.TrimLeft(*serverRoot, "/")
	if *serverRoot != "" && !strings.HasPrefix(*serverRoot, "/") {
		*serverRoot = "/" + *serverRoot
	}
	*serverRoot = strings.TrimRight(*serverRoot, "/")

	router := initRoutes()

	addr := fmt.Sprintf("0.0.0.0:%d", *port)
	log.Infof("Listenning on %s", addr)
	var err error
	if !*https {
		err = http.ListenAndServe(addr, router)
	} else {
		if *certFile == "" {
			log.Fatal("Missing certificate file")
		}
		if *keyFile == "" {
			log.Fatal("Missing certificate's private key file")
		}
		err = http.ListenAndServeTLS(addr, *certFile, *keyFile, router)
	}
	if err != nil {
		if err != http.ErrServerClosed {
			log.Error(err.Error())
		} else {
			log.Warn("Server is shutting down")
		}
	}
}
