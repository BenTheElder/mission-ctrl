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
	"github.com/thebenjaneer/mission-ctrl/binarydata"
)

const(
	VERSION = "Mission-Ctrl v0.0.1"
	MAIN_HTML ="<html><style>#canvas-container{width: 100%;text-align:center;}"+
"canvas {display: inline;}</style>"+
"<body style=\"background-color: rgb(56, 104, 146)\">"+
"<script type=\"text/javascript\" src=\"smoothie.js\"></script>"+
"<h1 style=\"text-align:center\">Cpu Usage (%):</h1>"+
"<div id=\"canvas-container\">"+
"<canvas id=\"cpucanvas\" width=\"400\" height=\"100\"></canvas>"+
"</div>"+
"<h1 style=\"text-align:center\">Memory Usage (%):</h1>"+
"<div id=\"canvas-container\">"+
"<canvas id=\"memcanvas\" width=\"400\" height=\"100\"></canvas>"+
"</div>"+
"<script type=\"text/javascript\">"+
"var line1 = new TimeSeries();"+
"var line2 = new TimeSeries();"+
"setInterval(function() {"+
"var xmlHttp = new XMLHttpRequest();"+
"xmlHttp.open( \"GET\", \"/data\", false );"+
"xmlHttp.send( null );"+
"var resp = xmlHttp.responseText;"+
"var data = JSON.parse(resp);"+
"line1.append(new Date().getTime(), data.cpu);"+
"line2.append(new Date().getTime(), data.mem);"+
"}, 500);"+
"var smoothie = new SmoothieChart({ grid: { strokeStyle: 'rgb(125, 0, 0)', fillStyle: 'rgb(60, 0, 0)', lineWidth: 1, millisPerLine: 250, verticalSections: 6 } });"+
"smoothie.addTimeSeries(line1, { strokeStyle: 'rgb(0, 255, 0)', fillStyle: 'rgba(0, 255, 0, 0.4)', lineWidth: 3 });"+
"smoothie.streamTo(document.getElementById(\"cpucanvas\"), 1000);"+
"var smoothie2 = new SmoothieChart({ grid: { strokeStyle: 'rgb(125, 0, 0)', fillStyle: 'rgb(60, 0, 0)', lineWidth: 1, millisPerLine: 250, verticalSections: 6 } });"+
"smoothie2.addTimeSeries(line2, { strokeStyle: 'rgb(0, 255, 0)', fillStyle: 'rgba(0, 255, 0, 0.4)', lineWidth: 3 });"+
"smoothie2.streamTo(document.getElementById(\"memcanvas\"), 1000);"+
"</script>"+
"</body>"+
"</html>"
)

type sysStats struct{ cpuPercent, memPercent float32 }

type mcHttpHandler struct{ ch chan sysStats}


func (h *mcHttpHandler)  ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		if r.URL.Path != "/favicon.ico"{
			if r.URL.Path == "/data" {
				sstats := <-h.ch
				fmt.Fprintf(w, "{\"cpu\": %f,\"mem\":%f}", sstats.cpuPercent, sstats.memPercent)
			} else if (r.URL.Path == "/smoothie.js") || (r.URL.Path == "smoothie.js") {
				data, _:= binarydata.Asset("smoothiecharts/smoothie.js")
				fmt.Fprint(w, string(data))
			} else {
				fmt.Fprint(w, MAIN_HTML)
			}
		}else{
			w.WriteHeader(http.StatusNotFound)
		}
	}
}



func main(){
//Get startup time
	START_TIME := time.Now()

//Create channels and tickers
	statsChan := make(chan sysStats)
	statsTicker := time.NewTicker(time.Millisecond * 500)
	stopChan := make(chan int)

//Handle ctrl-c etc.
	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		sig := <-signalChannel
		switch sig {
		case os.Interrupt:
			//handle SIGINT
			println("Caught SIGINT, exiting.")
			statsTicker.Stop()
			close(stopChan)
			os.Exit(0)
		case syscall.SIGTERM:
			//handle SIGTERM
			println("Caught SIGTERM, exiting.")
			statsTicker.Stop()
			close(stopChan)
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
	}
	println("Detected Addresses:")
	for _, addr := range addrs {
		fmt.Println("\t\t\t", addr.String())
	}

	
//Collect stats regularly
	sstat := sysStats{-1,-1}
	go func() {
		for _ = range statsTicker.C {
			i, j := stats.GetStats()
			sstat = sysStats{i,j}
		}
	}()

//Pass stats along to http server as needed
	go func(){	
		for{
			select {
			case statsChan <- sstat:
			case <- stopChan:
				return
			default:
			}
			time.Sleep(time.Millisecond * 100)
		}
	}()

//Create and launch web server
	s := &http.Server{
		Addr:           HTTP_ADDR,
		Handler:        &mcHttpHandler{statsChan},
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	fmt.Printf("Launching HTTP Server on %s\n", HTTP_ADDR)
	log.Fatal(s.ListenAndServe())
}