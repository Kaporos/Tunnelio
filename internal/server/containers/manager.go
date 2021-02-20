package containers

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	dockerclient "github.com/docker/docker/client"
	"github.com/rs/xid"
	"time"
)

type Forward struct {

	Port int
	Uid string

}

type Container struct {
	MaxForwards int
	Forwards []Forward
	Ip string
	Uid string
	Start int
	Id string
}

type ContainerManager struct {
	MaxContainers int
	Containers []Container
	DockerClient *dockerclient.Client
	MaxForwardPerContainer int
}

func (cm *ContainerManager) AddContainer() {
	fmt.Println("Creating container...")

	ctx := context.Background()
	networkConfig := &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{},
	}

	gatewayConfig := &network.EndpointSettings{}
	networkConfig.EndpointsConfig["tunnelio"] = gatewayConfig


	resp, err := cm.DockerClient.ContainerCreate(ctx, &container.Config{
		Image: "tunnelio:test",
	}, nil, networkConfig, nil, "")

	if err != nil {
		panic(err)
	}
	fmt.Println("Starting container...")

	if err := cm.DockerClient.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	cont, err := cm.DockerClient.ContainerInspect(ctx, resp.ID)
	if err != nil {
		panic(err)
	}
	var ip = cont.NetworkSettings.Networks["tunnelio"].IPAddress



	cm.Containers = append(cm.Containers, Container{
		Ip: ip,
		MaxForwards: cm.MaxForwardPerContainer,
		Forwards: []Forward{},
		Uid: xid.New().String(),
		Start: 8081,
		Id: resp.ID,
	})
}



func (cm *ContainerManager) CleanDocker() {
	removeOptions := types.ContainerRemoveOptions{
		RemoveVolumes: true,
		Force:         true,
	}

	for _, c := range cm.Containers {
		ctx := context.Background()

		if err := cm.DockerClient.ContainerStop(ctx, c.Id, nil); err != nil {
			fmt.Println("Unable to stop container %s: %s", c.Id, err)
		}

		if err := cm.DockerClient.ContainerRemove(ctx, c.Id, removeOptions); err != nil {
			fmt.Println("Unable to remove container %s: %s", c.Id, err)
		}
		fmt.Println(c.Id+" successfuly deleted !")

	}
}



func (cm *ContainerManager) GetContainer(uid string) *Container {
	var containerId int = -1
	for i, c := range cm.Containers {
		if c.Uid == uid {
			containerId = i
		}
	}
	if containerId != -1 {
		return &cm.Containers[containerId]
	}
	return nil
}



func (cm *ContainerManager) GetFreeShare() (string, string, string, string, int) {

	var containerId = -1
	for i, c := range cm.Containers {
		if len(c.Forwards) < c.MaxForwards {

			containerId = i
			break
		}
	}

	if containerId == -1 {
		fmt.Println("Max capacity reached! Let's create another container")
		if len(cm.Containers) < cm.MaxContainers {
			cm.AddContainer()
			time.Sleep(2 * time.Second) //WAITING CHISEL TO START
			return cm.GetFreeShare()
		} else {
			fmt.Println("Server is full !")
			return "","","","",0
		}


	}

	c := &cm.Containers[containerId]

	chid, chusr, chpass, port := c.getFreeShare()
	var uid = c.SharePort(port)
	return chid, chusr, chpass, uid, port


}

func (cm *ContainerManager) getContainerAndForwardFromForwardUid(uid string) (*Container,*Forward) {
	var contId = -1
	var forId = -1
	for i, c := range cm.Containers {
		for j, p := range c.Forwards {
			fmt.Println(p)
			if p.Uid == uid {
				contId = i
				forId = j
			}
		}
	}
	if contId == -1 {
		return nil, nil
	} else {
		var cont = &cm.Containers[contId]
		return cont, &cont.Forwards[forId]
	}
}

func (cm *ContainerManager) GetURLFromSharedUid(uid string) (string, bool){
	var cont, forw = cm.getContainerAndForwardFromForwardUid(uid)
	if cont == nil {
		return "", false
	}
	publicURL := fmt.Sprintf("http://%v:%v",cont.Ip , forw.Port)

	return publicURL, true
}
func indexOf(element Forward, data []Forward) (int) {
	for k, v := range data {
		if element == v {
			return k
		}
	}
	return -1    //not found.
}


func (cm *ContainerManager) StopShare(uid string) string {
	var cont, forw = cm.getContainerAndForwardFromForwardUid(uid)
	if cont == nil {
		return fmt.Sprintf("Tunnel %v not found. Can't delete it ", uid)
	}
	var index = indexOf(*forw, cont.Forwards)
	cont.Forwards = append(cont.Forwards[:index], cont.Forwards[index+1:]...)
	return fmt.Sprintf("Tunnel %v successfuly deleted", uid)
}


func (c *Container) sharePort(port int, id string) { //PRIVATE FUNCTION
	c.Forwards = append(c.Forwards, Forward{
		Port: port,
		Uid:  id,
	})
}

func (c *Container) SharePort(port int) string{
	id := xid.New().String()
	c.sharePort(port, id)
	return id
}

func (c *Container) SharePortWithUID(port int, uid string) string{
	c.sharePort(port, uid)
	return uid
}

func (c *Container) GetChiselUrl() string{
	return "http://"+c.Ip+":8080"
}

func (c *Container) getFreeShare() (string, string, string, int) {
	var port = c.Start + len(c.Forwards)
	if port > c.Start + c.MaxForwards {
		return "","","", 0
	}
	return c.Uid, "chisel", "chisel", port
}


func NewContainerManager() *ContainerManager {

	cli, err := dockerclient.NewClientWithOpts(dockerclient.FromEnv, dockerclient.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}


	return &ContainerManager{
		MaxContainers: 100,
		Containers:    []Container{},
		DockerClient: cli,
		MaxForwardPerContainer: 1000,
	}
}