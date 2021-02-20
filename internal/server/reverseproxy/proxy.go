package reverseproxy

import (
	"fmt"
	"github.com/Kaporos/Tunnelio/internal/server/containers"
	"github.com/Kaporos/Tunnelio/internal/server/handlers"
	"net/http"
	"strings"
)

type ReverseProxy struct {
	ChiselDomain string
}

var apiHandler handlers.ApiHandler
var chiselHandler handlers.ChiselHandler
var tunnelHandler handlers.TunnelHandler
var DockerManager containers.ContainerManager


func (rp *ReverseProxy) Handle(w http.ResponseWriter, r *http.Request) {

	var domain = strings.Split(r.Host, ":")[0]
	var splittedDomains = strings.Split(domain,".")


	if len(splittedDomains) == 2 { //No subdomain, ex: tunnelio.io
		apiHandler.Handle(w,r, &DockerManager)
		return

	}
	if splittedDomains[1] == rp.ChiselDomain{
		chiselHandler.Handle(w,r,splittedDomains[0], DockerManager)
		return
	}

	if len(splittedDomains) == 3 {
		tunnelHandler.Handle(w,r,splittedDomains[0], DockerManager)
		return
	}
	fmt.Fprint(w,"Invalid request")



}


func init(){
	apiHandler = handlers.ApiHandler{}
	chiselHandler = handlers.ChiselHandler{}
	tunnelHandler = handlers.TunnelHandler{}
	DockerManager = *containers.NewContainerManager()
	//DockerManager.GetContainer("test").SharePortWithUID(5000,"demotunnel")
}