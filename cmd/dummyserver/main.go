package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/markamdev/goloba/pkg/logger"
	"github.com/namsral/flag"
)

var (
	port    = flag.Int("port", 8070, "Dummy server listetning port")
	message = flag.String("message", "", "Server welcome message")
)

func main() {
	logger.SetLevel(logger.GlbDebug)
	logger.Debug("Starting dummy HTTP test server for GoLoBa")
	flag.Parse()

	logger.Debug("... listening on port: ", *port)
	logger.Debug("... welcome message: ", *message)

	http.HandleFunc("/", dummyHandler)
	status := http.ListenAndServe(":"+strconv.Itoa(*port), nil)
	logger.Debug("Server launching status: ", status)
}

func dummyHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "<html>\n<body>\n<h2>Dummy server message: "+*message+"</h2>\n</body>\n</head>\n")
}
