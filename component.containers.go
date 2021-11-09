package main

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/bketelsen/localspaces/operations"
	lxd "github.com/lxc/lxd/client"
	"github.com/lxc/lxd/shared"
	"github.com/lxc/lxd/shared/api"
	"github.com/yuriizinets/kyoto"
)

type InstanceDetail struct {
	api.InstanceFull
	IPAddresses      []string
	CleanIPAddresses []string
}

type ComponentContainers struct {
	Page        kyoto.Page `json:"-"`
	Instances   []InstanceDetail
	OperationID int64
	Status      operations.Status
	Logs        []string
}

func (c *ComponentContainers) Init(p kyoto.Page) {
	c.Page = p
}

func (c *ComponentContainers) Async(p kyoto.Page) error {
	err := c.getInstances()
	if err != nil {
		return err
	}
	if c.OperationID != 0 {
		var err error
		c.Status, err = operations.GetStatus(c.OperationID)
		return err
	}
	return nil
}

func (c *ComponentContainers) getInstances() error {
	server, err := lxd.ConnectLXDUnix("", nil)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	instances, err := server.GetInstancesFull(api.InstanceTypeAny)
	if err != nil {

		return err
	}
	c.Instances = []InstanceDetail{}

	for _, instance := range instances {
		i := InstanceDetail{}
		i.InstanceFull = instance
		i.IPAddresses = IP4ColumnData(instance)
		for _, ip := range i.IPAddresses {
			log.Println("IP", ip)
			cleanIP := strings.Split(ip, " ")
			log.Println("Clean IP", cleanIP[0])

			i.CleanIPAddresses = append(i.CleanIPAddresses, cleanIP[0])
		}
		c.Instances = append(c.Instances, i)
	}
	return nil
}

//
// Copied from github.com/lxc/lxd
func IP4ColumnData(cInfo api.InstanceFull) []string {
	if cInfo.IsActive() && cInfo.State != nil && cInfo.State.Network != nil {
		ipv4s := []string{}
		for netName, net := range cInfo.State.Network {
			if net.Type == "loopback" {
				continue
			}

			for _, addr := range net.Addresses {
				if shared.StringInSlice(addr.Scope, []string{"link", "local"}) {
					continue
				}

				if addr.Family == "inet" {
					ipv4s = append(ipv4s, fmt.Sprintf("%s (%s)", addr.Address, netName))
				}
			}
		}
		sort.Sort(sort.Reverse(sort.StringSlice(ipv4s)))
		return ipv4s
	}

	return []string{}
}
func (c *ComponentContainers) Actions() kyoto.ActionMap {

	return kyoto.ActionMap{
		"Start": func(args ...interface{}) {
			for arg := range args {
				log.Println("Start args", arg)
			}

			start := operations.NewContainerStart(args[0].(string))
			id, err := operations.Submit(start)
			if err != nil {
				log.Println("Error submitting start", err)
			}
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
			err = c.getInstances()
			if err != nil {
				log.Println(err)
			}
			kyoto.SSAFlush(c.Page, c)

		},

		"Stop": func(args ...interface{}) {
			for arg := range args {
				log.Println("Start args", arg)
			}

			stop := operations.NewContainerStop(args[0].(string))
			id, err := operations.Submit(stop)
			if err != nil {
				log.Println("Error submitting stop", err)
			}

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
			err = c.getInstances()
			if err != nil {
				log.Println(err)
			}
			kyoto.SSAFlush(c.Page, c)
		},

		"Delete": func(args ...interface{}) {
			for arg := range args {
				log.Println("Start args", arg)
			}
			dlt := operations.NewContainerDelete(args[0].(string))
			id, err := operations.Submit(dlt)
			if err != nil {
				log.Println("Error submitting delete", err)
			}

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
			err = c.getInstances()
			if err != nil {
				log.Println(err)
			}
			kyoto.SSAFlush(c.Page, c)
		},
	}
}
