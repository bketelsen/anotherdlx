package main

import (
	"html/template"
	"log"

	lxd "github.com/lxc/lxd/client"
	"github.com/lxc/lxd/shared/api"
	"github.com/yuriizinets/kyoto"
)

type ComponentTopNav struct {
	Title    string
	Projects []api.Project
	Items    []ComponentNavbarHref
}

type ComponentTopNavbarHref struct {
	Title     template.HTML
	Image     template.HTML
	Href      string
	AriaLabel string
	Current   bool
}

func (c *ComponentTopNav) Async(p kyoto.Page) error {
	server := kyoto.GetContext(p, "server").(lxd.InstanceServer)

	projects, err := server.GetProjects()
	if err != nil {
		log.Println(err.Error())
		return err
	}
	c.Projects = projects
	return nil
}
