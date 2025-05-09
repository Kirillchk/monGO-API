package main

import (
	"flag"
	"fmt"
	"net/http"
)

func main() {
	port := flag.Int("port", 8080, "Port to listen on")
	DBport := flag.Int("DBport", 27017, "Port to listen on")
	flag.Parse()
	portstring := fmt.Sprintf(":%d", *port)
	DBportstring := fmt.Sprintf("%d", *DBport)
	InitDB(DBportstring)
	InitHandlers()
	fmt.Printf("http://localhost%s/", portstring)
	http.ListenAndServe(portstring, nil)
}
