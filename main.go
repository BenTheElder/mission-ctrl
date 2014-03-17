package main

import (
	"flag"
	"fmt"
	"net"
	"log"
	"time"
	"net/http"
	_"io/ioutil"
	_"strings"
	"os"
	"os/signal"
	"syscall"
	"github.com/thebenjaneer/mission-ctrl/stats"
)

const(
	VERSION = "Mission-Ctrl v0.0.1"
)


type mcHttpHandler struct{}

func (h *mcHttpHandler)  ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		if r.URL.Path != "/favicon.ico"{
			fmt.Fprintf(w, "Nothing to see here...")
		}else{
			w.WriteHeader(http.StatusNotFound)
		}
	}
}

/*func statCollector(){

}
*/
func main(){
//Get startup time
	START_TIME := time.Now()
	
//Handle ctrl-c etc.
	signalChannel := make(chan os.Signal, 2)
    signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
    go func() {
        sig := <-signalChannel
        switch sig {
        case os.Interrupt:
            //handle SIGINT
            println("SIGINT")
            os.Exit(0)
        case syscall.SIGTERM:
            //handle SIGTERM
            println("SIGTERM")
            os.Exit(0)
        }
    }()
	
//Get options
//62875 = MCTRL
	var portFlag = flag.String("port", "62875", "port to serve on")
	flag.Parse()
	var HTTP_ADDR = ":"+*portFlag

//Print Welcome Message
	println("  __  __  ____  ___  ___  ____  _____  _  _        ___  ____  ____  __")  
	println(" (  \\/  )(_  _)/ __)/ __)(_  _)(  _  )( \\( ) ___  / __)(_  _)(  _ \\(  )")  
	println("  )    (  _)(_ \\__ \\\\__ \\ _)(_  )(_)(  )  ( (___)( (__   )(   )   / )(__ ")
	println(" (_/\\/\\_)(____)(___/(___/(____)(_____)(_)\\_)      \\___) (__) (_)\\_)(____)")
	println("==========================================================================")
	print("\n")
	println(VERSION)
	println("Started @ "+START_TIME.Format(time.UnixDate))

//Get and print network interface addresses	
	addrs, err := net.InterfaceAddrs()
	if err != nil{
		log.Fatal(err)
	}else{
		println("Detected Addresses:")
		for _, addr := range addrs {
			fmt.Println("\t\t\t", addr.String())
		}
	}

//Create and launch web server
	s := &http.Server{
		Addr:           HTTP_ADDR,
		Handler:        &mcHttpHandler{},
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	fmt.Printf("Launching HTTP Server on %s\n", HTTP_ADDR)
	log.Fatal(s.ListenAndServe())
}