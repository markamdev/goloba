package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/markamdev/goloba/pkg/utils"
	"github.com/namsral/flag"
	"github.com/sirupsen/logrus"
)

var (
	port    = flag.Int("port", 8070, "Dummy server listetning port")
	message = flag.String("message", "", "Server welcome message")
)

func main() {
	utils.SetupLogger()
	logrus.Debugln("Starting dummy HTTP test server for GoLoBa")
	flag.Parse()

	logrus.Debugln("... listening on port: ", *port)
	logrus.Debugln("... welcome message: ", *message)

	http.HandleFunc("/", dummyHandler)
	status := http.ListenAndServe(":"+strconv.Itoa(*port), nil)
	logrus.Debugln("Server launching status: ", status)
}

func dummyHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "<html>\n<body>\n<h2>Dummy server message: "+*message+"</h2>\n</body>\n</head>\n")
}
