package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

var serverConfig struct {
	listenPort     uint
	welcomeMessage string
}

func main() {
	fmt.Println("Starting dummy HTTP test server for GoLoBa")

	var port uint
	var welcome string

	// commnad line params
	flag.UintVar(&port, "p", 8080, "Listening port")
	hostname, _ := os.Hostname()
	flag.StringVar(&welcome, "m", hostname, "Welcome message (printed in HTML)")
	flag.Parse()

	fmt.Println("... listening on port: ", port)
	fmt.Println("... welcome message: ", welcome)

	serverConfig.listenPort = port
	serverConfig.welcomeMessage = welcome

	http.HandleFunc("/", dummyHandler)
	status := http.ListenAndServe(":"+strconv.Itoa(int(serverConfig.listenPort)), nil)
	fmt.Println("Server launching status: ", status)
}

func dummyHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "<html>\n<body>\n<h2>Dummy server message: "+serverConfig.welcomeMessage+"</h2>\n</body>\n</head>\n")
}
