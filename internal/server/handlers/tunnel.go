package handlers

import (
	"fmt"
	"github.com/Kaporos/Tunnelio/internal/server/containers"
	"net/http"
	"net/http/httputil"
	"net/url"
)


type TunnelHandler struct {}

func (ah TunnelHandler) Handle(w http.ResponseWriter, r *http.Request, tun string, dm containers.ContainerManager) {

	var pubUrl, ok = dm.GetURLFromSharedUid(tun)
	if !ok {
		fmt.Fprintf(w,"Tunnel %v not found.", tun)
		return
	}

	remote, _ := url.Parse(pubUrl)
	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.ServeHTTP(w,r)

}