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
	"context"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/google/gops/agent"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/juju/loggo"
	"github.com/soheilhy/cmux"
	"github.com/urfave/cli/v2"
	v1 "github.com/zhihu/cmdb/pkg/api/v1"
	"github.com/zhihu/cmdb/pkg/signals"
	"github.com/zhihu/cmdb/pkg/storage/tidb"
	"github.com/zhihu/cmdb/pkg/tools/grpcserver"
	"github.com/zhihu/cmdb/pkg/tools/httpserver"
	"github.com/zhihu/cmdb/pkg/tools/logger"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

var Version = "dev"

var log = loggo.GetLogger("")

const (
	FlagLogConf = "log_conf"
	FlagAddr    = "addr"
	FlagDSN     = "dsn"
	FlagGOPS    = "gops"
)

func main() {
	// handle term signal
	ctx := signals.SignalHandler(context.Background())

	app := cli.NewApp()
	app.Name = "cmdb"
	app.Version = Version
	app.Usage = "A programmable CMDB"
	app.Description = `Configuration Management Database (CMDB) is the source of truth and knowledge about all assets, whether is virtual or physical. See https://github.com/zhihu/cmdb for details`
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    FlagAddr,
			Usage:   "set serve address",
			EnvVars: []string{"LISTEN_ADDR"},
			Value:   ":8080",
		},
		&cli.StringFlag{
			Name:     FlagDSN,
			Usage:    "set data source name, see https://github.com/go-sql-driver/mysql",
			EnvVars:  []string{"TIDB_DSN"},
			Required: true,
		},
		&cli.StringFlag{
			Name:    FlagLogConf,
			FilePath: "./log.conf",
			Usage:   "config loggers level, see http://github.com/juju/loggo",
			EnvVars: []string{"LOG_CONF"},
			Value:   "<root>=INFO",
		},
		&cli.BoolFlag{
			Name:     FlagGOPS,
			Usage:    "start gops agent",
			EnvVars:  []string{"GOPS"},
			Required: true,
		},
	}
	app.Action = func(c *cli.Context) error {
		err := logger.Setup(c.String(FlagLogConf))
		if err != nil {
			return err
		}
		log.Debugf("cmdb version: %s", Version)
		return Run(c.Context, c.String(FlagAddr), c.String(FlagDSN), c.Bool(FlagGOPS))
	}

	err := app.RunContext(ctx, os.Args)
	if err != nil {
		log.Errorf("application run error: %s", err)
		os.Exit(1)
	}
}

func Run(ctx context.Context, addr string, dsn string, pprofAgent bool) error {
	if pprofAgent {
		defer agent.Close()
		go func() {
			_ = agent.Listen(agent.Options{})
		}()
	}
	// use cmux to support grpc and http on one port.
	listener, err := net.Listen("tcp", addr)
	cm := cmux.New(listener)
	grpcL := cm.Match(cmux.HTTP2HeaderField("content-type", "application/grpc"))
	httpL := cm.Match(cmux.HTTP1Fast())


	// create grpc server
	grpcServer := grpc.NewServer()
	// create grpc-gateway server
	gateway := runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{OrigName: true, EmitDefaults: true}))

	// init server instance.
	srv, err := InitServer(ctx, tidb.DSN(dsn))
	if err != nil {
		return err
	}
	srv.Register(grpcServer, gateway)

	group, ctx := errgroup.WithContext(ctx)
	// start http server
	group.Go(func() error {
		m := http.NewServeMux()
		// support swagger
		m.HandleFunc("/swagger", func(writer http.ResponseWriter, request *http.Request) {
			writer.Header().Set("Content-Type", "application/json")
			_, _ = writer.Write([]byte(v1.Swagger))
		})
		// support CORS requests
		m.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
			writer.Header().Set("Access-Control-Allow-Origin", "*")
			// grpc-gateway serve
			gateway.ServeHTTP(writer, request)
		})
		return httpserver.Run(ctx, httpL, m)
	})

	// start grpc server
	group.Go(func() error {
		return grpcserver.Run(ctx, grpcL, grpcServer)
	})

	group.Go(func() error {
		err = cm.Serve()
		// this error means that listener closed by term signal.
		if err != nil && strings.Contains(err.Error(), "use of closed network connection") {
			return nil
		}
		return err
	})
	return group.Wait()
}
