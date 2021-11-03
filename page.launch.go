package main

import (
	"html/template"
	"log"
	"net/http"
	"strings"

	lxd "github.com/lxc/lxd/client"
	"github.com/yuriizinets/kyoto"
)

type PageLaunch struct {
	// API
	Server lxd.InstanceServer

	// State
	Page string

	// Components
	CMeta              kyoto.Component
	Nav                kyoto.Component
	TopNav             kyoto.Component
	Launch             kyoto.Component
	DemoEmailValidator kyoto.Component
}

func (p *PageLaunch) Template() *template.Template {
	return template.Must(template.New("page.launch.html").Funcs(kyoto.TFuncMap()).ParseGlob("*.html"))
}

func (p *PageLaunch) Meta() kyoto.Meta {
	return kyoto.Meta{
		Title:       "dlux - LXD Development Environments",
		Description: "Fast and flexible development environments in LXD",
	}
}

func (p *PageLaunch) Init() {

	r := kyoto.GetContext(p, "internal:r").(*http.Request)
	p.Page = strings.TrimPrefix(r.URL.Path, "/")

	//p.Page = strings.ReplaceAll(r.URL.Path, "/", "")
	// Init components
	log.Println("Launch Request", p.Page)
	server, err := lxd.ConnectLXDUnix("", nil)
	if err != nil {
		log.Println(err.Error())
		return
	}
	kyoto.SetContext(p, "server", server)

	p.CMeta = kyoto.RegC(p, &ComponentMeta{})
	p.TopNav = kyoto.RegC(p, &ComponentTopNav{})
	p.Launch = kyoto.RegC(p, &ComponentLaunch{})

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
		},
	})
}
