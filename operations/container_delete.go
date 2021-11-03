package operations

import (
	"log"

	lxd "github.com/lxc/lxd/client"
)

type ContainerDelete struct {
	operationBase
	Name string `json:"name"`
}

func NewContainerDelete(name string) *ContainerDelete {

	return &ContainerDelete{
		Name: name,
	}
}

func (c *ContainerDelete) log(s string) {
	c.loglines = append(c.loglines, s)

}

func (c *ContainerDelete) setStatus(status Status) {
	c.status = status

}

func (c *ContainerDelete) Run() error {

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

	c.log("Deleting Container")

	op, err := server.DeleteContainer(name)
	if err != nil {
		c.setStatus(StatusError)
		c.log(err.Error())
		log.Println(err.Error())
		return err
	}

	c.log("Waiting for Container Deletion")
	// Wait for the operation to complete
	err = op.Wait()
	if err != nil {
		c.setStatus(StatusError)
		c.log(err.Error())
		log.Println(err.Error())
		return err
	}
	c.log("Container Deleted")

	c.setStatus(StatusFinished)
	return nil

}
func (c *ContainerDelete) Log() []string {
	return c.loglines
}
func (c *ContainerDelete) Status() Status {
	return c.status
}

func (c *ContainerDelete) SetStatus(status Status) {
	c.status = status
}
func (c *ContainerDelete) SetID(id int64) {
	c.id = id
}
func (c *ContainerDelete) ID() int64 {
	return c.id
}
