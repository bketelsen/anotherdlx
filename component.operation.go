package main

import (
	"log"
	"time"

	"github.com/bketelsen/localspaces/operations"
	"github.com/yuriizinets/kyoto"
)

type ComponentOperation struct {
	Page   kyoto.Page `json:"-"`
	ID     int64
	Status operations.Status
	Logs   []string
}

func (c *ComponentOperation) Init(p kyoto.Page) {
	c.Page = p
}

func (c *ComponentOperation) Async(p kyoto.Page) error {
	var err error
	c.Status, err = operations.GetStatus(c.ID)
	return err
}

func (c *ComponentOperation) Actions(p kyoto.Page) kyoto.ActionMap {

	return kyoto.ActionMap{
		"Status": func(args ...interface{}) {
			log.Printf("Fetching status for %d \n", c.ID)
			c.Status, _ = operations.GetStatus(c.ID)
			c.Logs, _ = operations.GetLogs(c.ID)

			if c.Status == operations.StatusFinished {
				// let the viewer see the logs -- maybe this isn't needed
				time.Sleep(time.Second * 2)
				kyoto.Redirect(&kyoto.RedirectParameters{
					Page:              c.Page,
					ResponseWriterKey: "internal:rw",
					RequestKey:        "internal:r",
					Target:            "/containers",
				})
			}
		},
	}
}
