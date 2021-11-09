package main

import (
	"log"
	"time"

	"github.com/bketelsen/localspaces/operations"
	"github.com/yuriizinets/kyoto"
)

type ComponentLaunch struct {
	Page        kyoto.Page `json:"-"`
	Repo        string
	Name        string
	Message     string
	OperationID int64
	Status      operations.Status
	Logs        []string
}

func (c *ComponentLaunch) Init(p kyoto.Page) {
	c.Page = p
}
func (c *ComponentLaunch) Actions(p kyoto.Page) kyoto.ActionMap {

	return kyoto.ActionMap{
		"Submit": func(args ...interface{}) {
			for arg := range args {
				log.Println("Launch args", arg)
			}
			// Do what you want here
			log.Println("Launch values", c.Repo, c.Name)

			container := operations.NewContainerCreate(c.Name, "dlxbase", c.Repo)
			id, err := operations.Submit(container)
			if err != nil {
				log.Println("Error submitting container", err)
			}
			log.Println("Submitted Job:", id)
			// ensure the host has the mount paths for project file storage
			/*
				err = project.CreateMountPath()
				if err != nil {
					return nil, errors.Wrap(err, "creating mount path on host")
				}
				err = project.CreateCommonMountPath()
				if err != nil {
					return nil, errors.Wrap(err, "creating common mount path on host")
				}

				// Mount the project directory into container FS
				devname := "persist"
				devSource := "source=" + project.InstanceMountPath(name)
				devPath := "path=" + project.ContainerMountPath()
				log.Println("Source: ", devSource)
				log.Println("Path: ", devPath)

				err = project.CreateInstanceMountPath(name)
				if err != nil {
					return nil, errors.Wrap(err, "failed to create host mount path")
				}
				err = addDevice(d, name, []string{devname, "disk", devSource, devPath})
				if err != nil {
					return nil, errors.Wrap(err, "failed to mount project directory")
				}
			*/
			log.Println("Submitted Job:", id)
			c.OperationID = id
			c.Status, err = operations.GetStatus(c.OperationID)
			if err != nil {
				log.Println("Error getting status", err)
				c.Logs = append(c.Logs, "Error getting status")
				c.Logs = append(c.Logs, err.Error())
			}
			c.Logs, _ = operations.GetLogs(c.OperationID)
			if err != nil {
				log.Println("Error getting operation logs", err)
				c.Logs = append(c.Logs, "Error getting operation logs")
				c.Logs = append(c.Logs, err.Error())
			}
			kyoto.SSAFlush(c.Page, c)
			for {
				c.Status, err = operations.GetStatus(c.OperationID)
				if err != nil {
					log.Println("Error getting status", err)
					c.Logs = append(c.Logs, "Error getting status")
					c.Logs = append(c.Logs, err.Error())
				}
				c.Logs, _ = operations.GetLogs(c.OperationID)
				if err != nil {
					log.Println("Error getting operation logs", err)
					c.Logs = append(c.Logs, "Error getting operation logs")
					c.Logs = append(c.Logs, err.Error())
				}
				kyoto.SSAFlush(c.Page, c)

				if c.Status == operations.StatusRunning {
					time.Sleep(time.Second)
				} else {
					break
				}

			}
			kyoto.Redirect(&kyoto.RedirectParameters{
				Page:              c.Page,
				ResponseWriterKey: "internal:rw",
				RequestKey:        "internal:r",
				Target:            "/containers",
			})
		},
	}
}
