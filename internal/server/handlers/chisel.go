package handlers

import (
	"fmt"
	"github.com/Kaporos/Tunnelio/internal/server/containers"
	"net/http"
	"net/http/httputil"
	"net/url"
)


type ChiselHandler struct {}




func (ah ChiselHandler) Handle(w http.ResponseWriter, r *http.Request, ch string, dm containers.ContainerManager) {

	var container = dm.GetContainer(ch)


	if container == nil {
		fmt.Fprintf(w,"This chisel instance (%v) dont exist !", ch)
		return

	}

	remote, _ := url.Parse(container.GetChiselUrl())
	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.ServeHTTP(w,r)




}

