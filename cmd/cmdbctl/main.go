// Copyright 2020 Zhizhesihai (Beijing) Technology Limited.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/juju/loggo"
	"github.com/urfave/cli/v2"
	v1 "github.com/zhihu/cmdb/pkg/api/v1"
	"github.com/zhihu/cmdb/pkg/signals"
	"google.golang.org/grpc"
)

var Version = "dev"

var log = loggo.GetLogger("")

const FlagServer = "server"

func main() {
	ctx := signals.SignalHandler(context.Background())
	app := cli.NewApp()
	app.Name = "cmdbctl"
	app.Version = Version
	app.Usage = "cmdbctl"
	app.Description = `cmdbctl. See https://github.com/zhihu/cmdb for details`
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:     FlagServer,
			Aliases:  []string{"s"},
			EnvVars:  []string{"CMDB_SERVER"},
			Value:    "127.0.0.1:8080",
			Required: false,
		},
	}
	var globalClient *Client
	{
		var watch = &cli.Command{
			Name:    "watch",
			Aliases: []string{"w"},
		}
		watch.Action = func(c *cli.Context) error {
			if c.Args().Len() == 0 {
				return errors.New("no type name")
			}
			return globalClient.Watch(c.Args().Get(0), c.Context)
		}

		var get = &cli.Command{
			Name:    "get",
			Aliases: []string{"g"},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "filter",
					Aliases: []string{"f"},
				},
			},
		}
		get.Action = func(c *cli.Context) error {
			var args = c.Args()
			if args.Len() == 0 {
				return errors.New("no type name")
			}
			return globalClient.Get(args.Get(0), args.Get(1), c.Context)
		}
		app.Commands = append(app.Commands, watch, get)
	}
	app.Before = func(c *cli.Context) (err error) {
		globalClient, err = NewClient(c.Context, c.String(FlagServer))
		return
	}
	err := app.RunContext(ctx, os.Args)
	if err != nil {
		log.Errorf("application run error: %s", err)
		os.Exit(1)
	}
}

type Client struct {
	conn              *grpc.ClientConn
	ObjectsClient     v1.ObjectsClient
	ObjectTypesClient v1.ObjectTypesClient
}

func NewClient(ctx context.Context, s string) (c *Client, err error) {
	c = &Client{}
	//ctx, cancel := context.WithTimeout(ctx, time.Second)
	//defer cancel()
	conn, err := grpc.DialContext(ctx, s, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Errorf("did not connect: %v", err)
		return nil, err
	}
	c.conn = conn
	c.ObjectsClient = v1.NewObjectsClient(conn)
	c.ObjectTypesClient = v1.NewObjectTypesClient(conn)
	return c, nil
}

func (c *Client) Watch(name string, ctx context.Context) error {
	req := &v1.ObjectListRequest{}
	req.Type = name
	w, err := c.ObjectsClient.Watch(ctx, req)
	if err != nil {
		return err
	}
	for {
		event, err := w.Recv()
		if err != nil {
			return err
		}
		json := &runtime.JSONPb{OrigName: true, EmitDefaults: true}
		msg, err := json.Marshal(event)
		fmt.Printf("event: %s %v \n", msg, err)
	}
}

func (c *Client) Get(name string, query string, ctx context.Context) error {
	req := &v1.ObjectListRequest{}
	req.Type = name
	req.Query = query
	req.View = v1.ObjectView_NORMAL
	resp, err := c.ObjectsClient.List(ctx, req)
	if err != nil {
		return err
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	fmt.Fprintf(w, "NAME\tSTATUS\tSTATE\tMETAS\tDESC\n")
	for _, object := range resp.Objects {
		buf := bytes.NewBuffer(nil)
		for name, value := range object.Metas {
			if buf.Len() > 0 {
				buf.WriteString(",")
			}
			buf.WriteString(name)
			buf.WriteString(": ")
			buf.WriteString(value.Value)
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", object.Name, object.Status, object.State, buf.String(), object.Description)
	}
	w.Flush()
	return err
}
