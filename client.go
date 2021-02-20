package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/Kaporos/Tunnelio/internal/shared"
	"github.com/gookit/color"
	chclient "github.com/jpillora/chisel/client"
	"github.com/jpillora/chisel/share/cos"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)


var publicEndpoint string
var url string
var uid string


func main(){

	port := flag.Int("port", 0, "Port to forward through tunnel")
	urla := flag.String("url", "", "Your instance url ( ex: tunnelio.local:8080 )")
	flag.Parse()
	if *port == 0 {
		log.Fatal("You have to specify a port (-port)")
	}
	if *urla == ""{
		log.Fatal("You have to specify an instance domain (ex: tunnelio.local ) (-url)")

	}
	url = *urla


	resp, err := http.Get(fmt.Sprintf("http://%v/forward", url))
	if err != nil {
		log.Fatal(err)
	}

	var forwardResp shared.ForwardResponse
	body, err := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &forwardResp)


	fmt.Println("Hello from Tunnelio Client")
	if forwardResp.ChiselId == "" {
		fmt.Println("Server is full !")
		os.Exit(0)
	}

	//fmt.Println(fmt.Sprintf("Chisel url : http://%v.chisel.%v", forwardResp.ChiselId, url))
	//fmt.Println(fmt.Sprintf("Identifiers : %v:%v", forwardResp.ChiselUsername, forwardResp.ChiselPassword))
	//fmt.Println("Port to use : ", forwardResp.AllowedPort)
	uid = forwardResp.TunnelUID
	publicEndpoint = fmt.Sprintf("http://%v.%v ",forwardResp.TunnelUID, url)
	var chiselArgs = fmt.Sprintf("--max-retry-count 0 --auth %v:%v http://%v.chisel.%v R:%v:%v", forwardResp.ChiselUsername, forwardResp.ChiselPassword, forwardResp.ChiselId, url, forwardResp.AllowedPort, *port)
	NewClient(chiselArgs)

}


func client(args []string) {
	flags := flag.NewFlagSet("client", flag.ContinueOnError)
	config := chclient.Config{Headers: http.Header{}}
	flags.StringVar(&config.Fingerprint, "fingerprint", "", "")
	flags.StringVar(&config.Auth, "auth", "", "")
	flags.DurationVar(&config.KeepAlive, "keepalive", 25*time.Second, "")
	flags.IntVar(&config.MaxRetryCount, "max-retry-count", -1, "")
	flags.DurationVar(&config.MaxRetryInterval, "max-retry-interval", 0, "")
	flags.StringVar(&config.Proxy, "proxy", "", "")
	flags.StringVar(&config.TLS.CA, "tls-ca", "", "")
	flags.BoolVar(&config.TLS.SkipVerify, "tls-skip-verify", false, "")
	flags.StringVar(&config.TLS.Cert, "tls-cert", "", "")
	flags.StringVar(&config.TLS.Key, "tls-key", "", "")
	hostname := flags.String("hostname", "", "")
	verbose := flags.Bool("v", false, "")

	flags.Parse(args)
	//pull out options, put back remaining args
	args = flags.Args()
	if len(args) < 2 {
		log.Fatalf("A server and least one remote is required")
	}
	config.Server = args[0]
	config.Remotes = args[1:]
	//default auth
	if config.Auth == "" {
		config.Auth = os.Getenv("AUTH")
	}
	//move hostname onto headers
	if *hostname != "" {
		config.Headers.Set("Host", *hostname)
	}
	//ready
	c, err := chclient.NewClient(&config)
	if err != nil {
		log.Fatal(err)
	}
	c.Debug = *verbose
	go cos.GoStats()
	ctx := cos.InterruptContext()
	if err := c.Start(ctx); err != nil {
		log.Fatal(err)
	}

	color.Cyan.Printf("ENDPOINT: %v\n", publicEndpoint)

	if err := c.Wait(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Exiting...")
	resp, _ := http.Get(fmt.Sprintf("http://%v/free?uid=%v", url, uid))
	body, err := ioutil.ReadAll(resp.Body)
	color.Cyan.Printf(string(body)+"\n")

	fmt.Println("Goodbye !")

}

func NewClient(args string) {

	client(strings.Split(args, " "))

}

