package client

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/hyperhq/hyper/engine"

	gflag "github.com/jessevdk/go-flags"
)

func (cli *HyperClient) HyperCmdList(args ...string) error {
	var opts struct {
		Aux bool   `short:"x" long:"aux" default:"false" value-name:"false" description:"show the auxiliary containers"`
		Pod string `short:"p" long:"pod" value-name:"\"\"" description:"only list the specified pod"`
	}

	var parser = gflag.NewParser(&opts, gflag.Default|gflag.IgnoreUnknown)
	parser.Usage = "list [OPTIONS] [pod|container]\n\nlist all pods or container information"
	args, err := parser.Parse()
	if err != nil {
		if !strings.Contains(err.Error(), "Usage") {
			return err
		} else {
			return nil
		}
	}
	var item string
	if len(args) == 1 {
		item = "pod"
	} else {
		item = args[1]
	}

	if item != "pod" && item != "vm" && item != "container" {
		return fmt.Errorf("Error, the %s can not support %s list!", os.Args[0], item)
	}

	v := url.Values{}
	v.Set("item", item)
	if opts.Aux {
		v.Set("auxiliary", "yes")
	}
	if opts.Pod != "" {
		v.Set("pod", opts.Pod)
	}
	body, _, err := readBody(cli.call("GET", "/list?"+v.Encode(), nil, nil))
	if err != nil {
		return err
	}
	out := engine.NewOutput()
	remoteInfo, err := out.AddEnv()
	if err != nil {
		return err
	}

	if _, err := out.Write(body); err != nil {
		fmt.Printf("Error reading remote info: %s", err)
		return err
	}
	out.Close()

	var (
		vmResponse        = []string{}
		podResponse       = []string{}
		containerResponse = []string{}
	)
	if remoteInfo.Exists("item") {
		item = remoteInfo.Get("item")
	}
	if remoteInfo.Exists("Error") {
		return fmt.Errorf("Found an error while getting %s list: %s", item, remoteInfo.Get("Error"))
	}

	if item == "vm" {
		vmResponse = remoteInfo.GetList("vmData")
	}
	if item == "pod" {
		podResponse = remoteInfo.GetList("podData")
	}
	if item == "container" {
		containerResponse = remoteInfo.GetList("cData")
	}

	//fmt.Printf("Item is %s\n", item)
	if item == "vm" {
		fmt.Printf("%15s%20s\n", "VM name", "Status")
		for _, vm := range vmResponse {
			fields := strings.Split(vm, ":")
			fmt.Printf("%15s%20s\n", fields[0], fields[2])
		}
	}

	if item == "pod" {
		fmt.Printf("%15s%30s%20s%10s\n", "POD ID", "POD Name", "VM name", "Status")
		for _, p := range podResponse {
			fields := strings.Split(p, ":")
			var podName = fields[1]
			if len(fields[1]) > 27 {
				podName = fields[1][:27]
			}
			fmt.Printf("%15s%30s%20s%10s\n", fields[0], podName, fields[2], fields[3])
		}
	}

	if item == "container" {
		fmt.Printf("%-66s%-20s%15s%10s\n", "Container ID", "Name", "POD ID", "Status")
		for _, c := range containerResponse {
			fields := strings.Split(c, ":")
			name := fields[1]
			if len(name) > 0 {
				if name[0] == '/' {
					name = name[1:]
				}
			}
			fmt.Printf("%-66s%-20s%15s%10s\n", fields[0], name, fields[2], fields[3])
		}
	}
	return nil
}
