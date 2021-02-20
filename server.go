package main

import (
	"fmt"
	"github.com/Kaporos/Tunnelio/internal/server/reverseproxy"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)


const (
	chiselSubdomain = "chisel"
)

func main(){
	fmt.Println("Hello from Tunnelio Server")
	fmt.Printf("Starting server at port 8080\n")

	var reverseProxy = reverseproxy.ReverseProxy{
		ChiselDomain: chiselSubdomain,
	}

	ch := make(chan os.Signal, 3)

	signal.Notify(ch, os.Interrupt,syscall.SIGTERM,syscall.SIGINT)


	go func() {
		_ = <-ch
		signal.Stop(ch)
		fmt.Println("Exit command received. Exiting...")
		reverseproxy.DockerManager.CleanDocker()
		os.Exit(0)

	}()




	http.HandleFunc("/",reverseProxy.Handle)



	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}