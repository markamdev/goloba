package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/namsral/flag"
)

var (
	port    = flag.Int("port", 8070, "Dummy server listetning port")
	message = flag.String("message", "", "Server welcome message")
)

func main() {
	fmt.Println("Starting dummy HTTP test server for GoLoBa")
	flag.Parse()

	fmt.Println("... listening on port: ", *port)
	fmt.Println("... welcome message: ", *message)

	http.HandleFunc("/", dummyHandler)
	status := http.ListenAndServe(":"+strconv.Itoa(*port), nil)
	fmt.Println("Server launching status: ", status)
}

func dummyHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "<html>\n<body>\n<h2>Dummy server message: "+*message+"</h2>\n</body>\n</head>\n")
}
