package operations

import (
	"log"

	lxd "github.com/lxc/lxd/client"
	"github.com/lxc/lxd/shared/api"
)

type ContainerCreate struct {
	operationBase
	Name      string `json:"name"`
	BaseImage string `json:"base_image"`
	Repo      string `json:"repo"`
}

func NewContainerCreate(name, baseImage, repo string) *ContainerCreate {

	return &ContainerCreate{
		Name:      name,
		BaseImage: baseImage,
		Repo:      repo,
	}
}

func (c *ContainerCreate) log(s string) {
	c.loglines = append(c.loglines, s)

}

func (c *ContainerCreate) setStatus(status Status) {
	c.status = status

}

func (c *ContainerCreate) Run() error {

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

	bi := "dlxbase"

	source := api.ContainerSource{
		Type: "image",
		//Server:   "https://cloud-images.ubuntu.com/daily",
		//Alias:    getImage(),
		Alias: bi,
		//Protocol: "simplestreams",
	}

	c.log("Creating Container")
	req := api.ContainersPost{
		Name: name,
		ContainerPut: api.ContainerPut{
			Profiles: []string{"default"}, // TODO: ? support command line adding profiles
		},
		Source: source,
	}

	// Get LXD to create the container (background operation)
	op, err := server.CreateContainer(req)
	if err != nil {
		c.setStatus(StatusError)
		c.log(err.Error())
		log.Println(err.Error())
		return err
	}

	c.log("Waiting for Container")
	// Wait for the operation to complete
	err = op.Wait()
	if err != nil {
		c.setStatus(StatusError)
		c.log(err.Error())
		log.Println(err.Error())
		return err
	}

	c.log("Starting Container")
	// Get LXD to start the container (background operation)
	reqState := api.ContainerStatePut{
		Action:  "start",
		Timeout: -1,
	}

	op, err = server.UpdateContainerState(name, reqState, "")
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
func (c *ContainerCreate) Log() []string {
	return c.loglines
}
func (c *ContainerCreate) Status() Status {
	return c.status
}

func (c *ContainerCreate) SetStatus(status Status) {
	c.status = status
}
func (c *ContainerCreate) SetID(id int64) {
	c.id = id
}
func (c *ContainerCreate) ID() int64 {
	return c.id
}
