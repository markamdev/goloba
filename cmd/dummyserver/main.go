package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

const (
	envNamePort     = "DUMMY_PORT"
	envNameMessage  = "DUMMY_MESSAGE"
	flagNamePort    = "p"
	flagNameMessage = "m"
)

type serverConfig struct {
	listenPort     uint
	welcomeMessage string
}

var _config serverConfig

func main() {
	fmt.Println("Starting dummy HTTP test server for GoLoBa")

	var err error
	_config, err = loadConfig()
	if err != nil {
		panic(fmt.Errorf("Failed to load config: %v", err))
	}

	fmt.Println("... listening on port: ", _config.listenPort)
	fmt.Println("... welcome message: ", _config.welcomeMessage)

	http.HandleFunc("/", dummyHandler)
	status := http.ListenAndServe(":"+strconv.Itoa(int(_config.listenPort)), nil)
	fmt.Println("Server launching status: ", status)
}

func dummyHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "<html>\n<body>\n<h2>Dummy server message: "+_config.welcomeMessage+"</h2>\n</body>\n</head>\n")
}

func loadConfig() (serverConfig, error) {
	// commnad line params
	hostname, _ := os.Hostname()
	fPort := flag.Uint(flagNamePort, 9000, "Listening port")
	fMessage := flag.String(flagNameMessage, hostname, "Welcome message (printed in HTML)")
	flag.Parse()

	result := serverConfig{}

	port, err := loadInt(envNamePort, fPort)
	if err != nil {
		return serverConfig{}, err
	}
	result.listenPort = port

	message, err := loadString(envNameMessage, fMessage)
	if err != nil {
		return serverConfig{}, err
	}
	result.welcomeMessage = message

	return result, nil
}

func loadInt(env string, flg *uint) (uint, error) {
	eValue := os.Getenv(env)
	if len(eValue) > 0 {
		// parse and use env var
		iValue, err := strconv.ParseInt(eValue, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("Failed to fetch int from %s=%s", env, eValue)
		}
		return uint(iValue), nil
	}
	// check flag
	if flg == nil || *flg < 1 {
		return 0, fmt.Errorf("Failed to load value from flag")
	}
	return *flg, nil
}

func loadString(env string, flg *string) (string, error) {
	eValue := os.Getenv(env)
	if len(eValue) > 0 {
		return eValue, nil
	}
	if flg == nil || len(*flg) == 0 {
		return "", fmt.Errorf("Failed to load string from flag")
	}
	return *flg, nil
}
