package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	lxd "github.com/lxc/lxd/client"
	"github.com/yuriizinets/kyoto"
)

type PageOperation struct {
	// API
	Server lxd.InstanceServer

	// State
	Page string

	// Components
	CMeta     kyoto.Component
	Nav       kyoto.Component
	TopNav    kyoto.Component
	Header    kyoto.Component
	Operation kyoto.Component
}

func (p *PageOperation) Template() *template.Template {
	fm := template.FuncMap{"divide": func(a, b int64) int64 {
		return a / b
	}}
	return template.Must(template.New("page.operation.html").Funcs(fm).Funcs(kyoto.TFuncMap()).ParseGlob("*.html"))
}

func (p *PageOperation) Meta() kyoto.Meta {
	return kyoto.Meta{
		Title:       "dlux - LXD Development Environments",
		Description: "Fast and flexible development environments in LXD",
	}
}

func (p *PageOperation) Init() {

	r := kyoto.GetContext(p, "internal:r").(*http.Request)
	p.Page = strings.TrimPrefix(r.URL.Path, "/")

	log.Println(r.URL.Query())
	id := r.URL.Query().Get("id")
	convertedID, err := strconv.Atoi(id)
	if err != nil {
		log.Println(err.Error())
		return
	}
	//p.Page = strings.ReplaceAll(r.URL.Path, "/", "")
	// Init components
	log.Println("Operation Request", p.Page)
	server, err := lxd.ConnectLXDUnix("", nil)
	if err != nil {
		log.Println(err.Error())
		return
	}
	kyoto.SetContext(p, "server", server)

	p.CMeta = kyoto.RegC(p, &ComponentMeta{})
	p.TopNav = kyoto.RegC(p, &ComponentTopNav{})
	p.Header = kyoto.RegC(p, &ComponentHeader{})

	p.Operation = kyoto.RegC(p, &ComponentOperation{
		ID: int64(convertedID),
	})

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
