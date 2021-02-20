package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/Kaporos/Tunnelio/internal/server/containers"
	"github.com/Kaporos/Tunnelio/internal/shared"
	"net/http"
	"strings"
)





type ApiHandler struct {}

func (ah ApiHandler) Handle(w http.ResponseWriter, r *http.Request, dm *containers.ContainerManager) {
	switch strings.Split(r.RequestURI,"?")[0] {
	case "/forward":
		Forward(w,r,dm)
		return
	case "/free":
		Free(w,r,dm)
		return

	default:
		break
	}


	fmt.Fprint(w,"Hello from API")
}

func Forward(w http.ResponseWriter, r *http.Request, dm *containers.ContainerManager) {

	chid, chusr, chpass, tunnId, port := dm.GetFreeShare()

	var response = shared.ForwardResponse{
		ChiselId:      chid,
		ChiselUsername: chusr,
		ChiselPassword: chpass,
		AllowedPort:    port,
		TunnelUID: tunnId,
	}

	var respString, _ = json.Marshal(response)
	fmt.Fprint(w,string(respString))
}

func Free(w http.ResponseWriter, r *http.Request, dm *containers.ContainerManager) {
	uids, ok := r.URL.Query()["uid"]

	if !ok || len(uids[0]) < 1 {
		fmt.Fprintf(w,"Url Param 'port' is missing")
		return
	}

	// Query()["key"] will return an array of items,
	// we only want the single item.
	uid := uids[0]
	fmt.Fprintf(w,	dm.StopShare(uid))


}