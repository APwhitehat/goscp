package main

import (
	"flag"
	"net/http"
	"runtime"

	"github.com/APwhitehat/goscp"
	"github.com/sirupsen/logrus"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	port := flag.String("port", "8080", "Port")
	flag.Parse()

	logrus.Info("starting on port: ", *port)

	http.HandleFunc("/scp", goscp.ScpHandler)
	http.ListenAndServe(":"+*port, nil)
}
