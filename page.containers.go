package main

import (
	"html/template"
	"log"
	"net/http"
	"strings"

	lxd "github.com/lxc/lxd/client"
	"github.com/yuriizinets/kyoto"
)

type PageContainers struct {
	// API
	Server lxd.InstanceServer

	// State
	Page string

	// Components
	CMeta      kyoto.Component
	Nav        kyoto.Component
	Header     kyoto.Component
	TopNav     kyoto.Component
	Containers kyoto.Component
}

func (p *PageContainers) Template() *template.Template {
	fm := template.FuncMap{"divide": func(a, b int64) int64 {
		return a / b
	}}
	return template.Must(template.New("page.containers.html").Funcs(fm).Funcs(kyoto.TFuncMap()).ParseGlob("*.html"))
}

func (p *PageContainers) Meta() kyoto.Meta {
	return kyoto.Meta{
		Title:       "dlux - LXD Development Environments",
		Description: "Fast and flexible development environments in LXD",
	}
}

func (p *PageContainers) Init() {

	r := kyoto.GetContext(p, "internal:r").(*http.Request)
	p.Page = strings.TrimPrefix(r.URL.Path, "/")

	//p.Page = strings.ReplaceAll(r.URL.Path, "/", "")
	// Init components
	log.Println("ContainersRequest", p.Page)
	server, err := lxd.ConnectLXDUnix("", nil)
	if err != nil {
		log.Println(err.Error())
		return
	}
	kyoto.SetContext(p, "server", server)

	p.CMeta = kyoto.RegC(p, &ComponentMeta{})
	p.TopNav = kyoto.RegC(p, &ComponentTopNav{})
	p.Header = kyoto.RegC(p, &ComponentHeader{})
	p.Containers = kyoto.RegC(p, &ComponentContainers{})
	p.Nav = kyoto.RegC(p, &ComponentNav{
		Title: "dlux",
		Items: []ComponentNavbarHref{
			{
				Title:   "Dashboard",
				Href:    "/",
				Image:   "<i class='fas fa-tv mr-2 text-sm opacity-75'></i>",
				Current: p.Page == "",
			},
			{
				Title: "Containers",
				Href:  "/containers",
				Image: "<i class='fas fa-tools mr-2 text-sm text-blueGray-300'></i>",

				Current: p.Page == "containers",
			},
			{
				Title: "New Container",
				Href:  "/launch",
				Image: "<i class='fas fa-table mr-2 text-sm text-blueGray-300'></i>",

				Current: p.Page == "launch",
			},
			{
				Title: "Automatic Mounts",
				Href:  "/mounts",
				Image: "<i class='fas fa-table mr-2 text-sm text-blueGray-300'></i>",

				Current: p.Page == "mounts",
			},
		},
	})
}
