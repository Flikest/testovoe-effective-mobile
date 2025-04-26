package main

import (
	"flag"
)

func main() {
	var port string
	flag.StringVar(&port, "port", "", "Port for the application")

}
