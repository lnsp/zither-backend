package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/zither-oss/zither-backend/player"
	"github.com/zither-oss/zither-backend/routes"
)

var (
	hostname = flag.String("h", "localhost", "listen for clients at host")
	port     = flag.String("p", "8080", "listen for clients at port")
)

const helpText = `usage:
	zitherd [-p 8080] [-h localhost] [remote-host] [remote-port]`

func showHelpAndExit() {
	fmt.Println(helpText)
	os.Exit(0)
}

func main() {
	flag.Usage = showHelpAndExit
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		showHelpAndExit()
	}

	remoteHost, remotePort := args[0], args[1]
	remote, err := player.Connect(remoteHost, remotePort)
	if err != nil {
		log.Fatal(err)
	}

	router := routes.New(remote)
	http.Handle("/", router)
	if err := http.ListenAndServe(remoteHost+":"+remotePort, nil); err != nil {
		log.Fatal(err)
	}
}
