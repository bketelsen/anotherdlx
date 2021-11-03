package operations

import (
	"log"

	lxd "github.com/lxc/lxd/client"
	"github.com/lxc/lxd/shared/api"
)

type ContainerStop struct {
	operationBase
	Name string `json:"name"`
}

func NewContainerStop(name string) *ContainerStop {

	return &ContainerStop{
		Name: name,
	}
}

func (c *ContainerStop) log(s string) {
	c.loglines = append(c.loglines, s)

}

func (c *ContainerStop) setStatus(status Status) {
	c.status = status

}

func (c *ContainerStop) Run() error {

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

	c.log("Stopping Container")
	// Get LXD to start the container (background operation)
	reqState := api.ContainerStatePut{
		Action:  "stop",
		Timeout: -1,
		Force:   true,
	}

	op, err := server.UpdateContainerState(name, reqState, "")
	if err != nil {
		c.setStatus(StatusError)
		c.log(err.Error())
		log.Println(err.Error())
		return err
	}

	c.log("Waiting for Container Stop")
	// Wait for the operation to complete
	err = op.Wait()
	if err != nil {
		c.setStatus(StatusError)
		c.log(err.Error())
		log.Println(err.Error())
		return err
	}
	c.log("Container Stopped")

	c.setStatus(StatusFinished)
	return nil

}
func (c *ContainerStop) Log() []string {
	return c.loglines
}
func (c *ContainerStop) Status() Status {
	return c.status
}

func (c *ContainerStop) SetStatus(status Status) {
	c.status = status
}
func (c *ContainerStop) SetID(id int64) {
	c.id = id
}
func (c *ContainerStop) ID() int64 {
	return c.id
}
