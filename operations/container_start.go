package operations

import (
	"log"

	lxd "github.com/lxc/lxd/client"
	"github.com/lxc/lxd/shared/api"
)

type ContainerStart struct {
	operationBase
	Name string `json:"name"`
}

func NewContainerStart(name string) *ContainerStart {

	return &ContainerStart{
		Name: name,
	}
}

func (c *ContainerStart) log(s string) {
	c.loglines = append(c.loglines, s)

}

func (c *ContainerStart) setStatus(status Status) {
	c.status = status

}

func (c *ContainerStart) Run() error {

	c.setStatus(StatusRunning)

	//d = d.UseProject(project.Name)
	c.log("Connecting to Container Server")
	server, err := lxd.ConnectLXDUnix("", nil)
	if err != nil {
		c.setStatus(StatusError)
		c.log(err.Error())
		log.Println(err.Error())
		return err
	}

	name := c.Name

	c.log("Starting Container")
	// Get LXD to start the container (background operation)
	reqState := api.ContainerStatePut{
		Action:  "start",
		Timeout: -1,
	}

	op, err := server.UpdateContainerState(name, reqState, "")
	if err != nil {
		c.setStatus(StatusError)
		c.log(err.Error())
		log.Println(err.Error())
		return err
	}

	c.log("Waiting for Container Start")
	// Wait for the operation to complete
	err = op.Wait()
	if err != nil {
		c.setStatus(StatusError)
		c.log(err.Error())
		log.Println(err.Error())
		return err
	}
	c.log("Container Started")

	c.setStatus(StatusFinished)
	return nil

}
func (c *ContainerStart) Log() []string {
	return c.loglines
}
func (c *ContainerStart) Status() Status {
	return c.status
}

func (c *ContainerStart) SetStatus(status Status) {
	c.status = status
}
func (c *ContainerStart) SetID(id int64) {
	c.id = id
}
func (c *ContainerStart) ID() int64 {
	return c.id
}
