package configs

import (
	"os"
	
	"github.com/urfave/cli"
)


type configs struct {
	File      string
	Port      int
	Hostnames []string
}

var Values configs

var flags = []cli.Flag {
	cli.StringFlag{
		Name: "config, c",
		Usage: "YAML config `FILE`.",
		Destination: &Values.File,
	},
	cli.IntFlag{
		Name: "port, p",
		Usage: "Port to listen on.",
		Value: 80,
		Destination: &Values.Port,
	},
	cli.StringSliceFlag{
		Name: "hostnames, n",
		Usage: "Hostname(s) to listen to.",
	},
}

func Init() error {
	app        := cli.NewApp()
	app.Name    = "Book Review"
	app.Usage   = "Better book reviews!"
	app.Version = "0.0.1"
	app.Flags   = flags
	
	app.Action = func(c *cli.Context) error {
		Values.Hostnames = c.StringSlice("hostnames")
		return nil
	}
	
	err := app.Run(os.Args)
	return err
}
