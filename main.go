package main

import (
	"os"
	"runtime"

	"github.com/codegangsta/cli"
	"github.com/quorumsco/application"
	"github.com/quorumsco/cmd"
	"github.com/quorumsco/logs"
	"github.com/quorumsco/router"
	"github.com/quorumsco/settings"
	"github.com/quorumsco/wercker/controllers"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	cmd := cmd.New()
	cmd.Name = "contacts"
	cmd.Usage = "quorums contacts backend"
	cmd.Version = "0.0.1"
	cmd.Before = serve
	cmd.Flags = append(cmd.Flags, []cli.Flag{
		cli.StringFlag{Name: "config, c", Usage: "configuration file", EnvVar: "CONFIG"},
		cli.BoolFlag{Name: "debug, d", Usage: "debug log level"},
		cli.HelpFlag,
	}...)
	cmd.RunAndExitOnError()
}

func serve(ctx *cli.Context) error {
	var app = application.New()

	config, err := settings.Parse(ctx.String("config"))
	if err != nil && ctx.String("config") != "" {
		logs.Error(err)
	}

	if ctx.Bool("debug") || config.Debug() {
		logs.Level(logs.DebugLevel)
	}

	app.Components["Mux"] = router.New()

	if ctx.Bool("debug") || config.Debug() {
		app.Use(router.Logger)
	}

	app.Use(app.Apply)
	app.Post("/trigger/:id", controllers.TriggerBuild)

	server, err := config.Server()
	if err != nil {
		logs.Critical(err)
		os.Exit(1)
	}
	return app.Serve(server.String())
}
