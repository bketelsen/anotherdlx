package main

import (
	"github.com/yuriizinets/kyoto"
)

type ComponentMain struct {
	Message    string
	Containers kyoto.Component
}

func (c *ComponentMain) Init(p kyoto.Page) {

	c.Containers = kyoto.RegC(p, &ComponentContainers{})

}
