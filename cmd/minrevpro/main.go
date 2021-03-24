package main

import (
	"flag"
	"log"
	"net/url"

	"github.com/thlcodes/minrevpro"
)

const appName = "minrevpro"
const appVersion = "v1"

const secretHeader = "X-SECRET"

var debugFlag = flag.Bool("debug", false, "Verbose output")
var targetFlag = flag.String("target", "", "Target URI")
var portFlag = flag.String("port", "8080", "Proxy port")
var hostFlag = flag.String("host", "", "Proxy host")
var secretFlag = flag.String("secret", "", "Proxy secret")
var basePathFlag = flag.String("basepath", "", "base path")

func main() {
	flag.Parse()

	if *targetFlag == "" {
		log.Fatalf("ERROR: Target es empty")
	}
	target, err := url.Parse(*targetFlag)
	if err != nil {
		log.Fatalf("ERROR: Target is invalid: " + err.Error())
	}

	addr := *hostFlag + ":" + *portFlag
	opts := []minrevpro.OptionFunc{
		minrevpro.WithAddr(addr),
		minrevpro.Debug(*debugFlag),
	}
	if *secretFlag != "" {
		opts = append(opts, minrevpro.WithSecret(secretHeader, *secretFlag))
	}
	if *basePathFlag != "" {
		opts = append(opts, minrevpro.WithBasePath(*basePathFlag))
	}
	proxy := minrevpro.NewReverseProxy(
		target,
		opts...,
	)

	log.Printf("Starting %s %s on %s:%s proxying %s", appName, appVersion, *hostFlag, *portFlag, *targetFlag)
	if err := proxy.Start(); err != nil {
		log.Fatalf("ERROR: could not start proxy: " + err.Error())
	}
}
