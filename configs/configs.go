package configs

import (
	"os"
	"log"
	"reflect"
	"errors"
	
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)


type Definitions struct {
	File      string    `yaml:"File"`
	Port      int       `yaml:"Port"`
	Hostnames []string  `yaml:"Hostnames"`
}

var values Definitions

var flags = []cli.Flag {
	cli.StringFlag{
		Name: "config, c",
		Usage: "YAML config `FILE`.",
		Destination: &values.File,
	},
	cli.IntFlag{
		Name: "Port, p",
		Usage: "`PORT` to listen on.",
		Value: 80,
		Destination: &values.Port,
	},
	cli.StringSliceFlag{
		Name: "Hostnames, n",
		Usage: "`HOSTNAME`s to listen to.",
		Value: &cli.StringSlice{"localhost"},
	},
}

func LoadFile(c *Definitions) error {
	if values.File == "" {
		// No config file was provided...
		log.Printf("No config file provided...\n")
		return nil
	}
	log.Printf("Loading config file from %v", values.File)
	yamlFile, err := ioutil.ReadFile(values.File)
	if err != nil {
		log.Printf("yamlFile.Get err #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	log.Printf("Config file loaded!\n")
	
	return nil
}

func UpdateValues(fileConfigs *Definitions, c *cli.Context) {
	fileFields := reflect.TypeOf(*fileConfigs)
	fileValues := reflect.ValueOf(*fileConfigs)
	vv := reflect.ValueOf(&values).Elem()
	for i := 0; i < fileFields.NumField(); i++ {
		field := fileFields.Field(i)
		if !c.IsSet(field.Name) {
			v := vv.FieldByName(field.Name)
			switch v.Kind() {
			case reflect.String:
				x := fileValues.FieldByName(field.Name)
				y := x.Interface().(string)
				if y == "" {
					continue
				}
				v.SetString(y)
			case reflect.Int:
				x := fileValues.FieldByName(field.Name)
				y := x.Interface().(int)
				if y == 0 {
					continue
				}
				v.SetInt(int64(y))
			case reflect.Slice:
				// Handled elsewhere due to urfave pkg bug
				break
			default:
				panic(errors.New("Invalid config type recvd"))
			}
		}
	}
}

func Init() (*Definitions, error) {
	app        := cli.NewApp()
	app.Name    = "Book Review"
	app.Usage   = "Better book reviews!"
	app.Version = "0.0.1"
	app.Flags   = flags

	app.Action = func(c *cli.Context) error {
		var yamlConfigs Definitions
		LoadFile(&yamlConfigs)

		// Override config values passed as params
		values.Hostnames = c.StringSlice("Hostnames")
		// XXX: StringSlice hack exists due to:
		// https://github.com/urfave/cli/issues/790
		if c.IsSet("Hostnames") || c.IsSet("n") {
			values.Hostnames = values.Hostnames[1:]
		} else if len(yamlConfigs.Hostnames) > 0 {
			values.Hostnames = yamlConfigs.Hostnames
		}

		UpdateValues(&yamlConfigs, c)
		log.Printf(
			"Configs loaded:\nPort: %d\nHostnames: %v\n",
			values.Port,
			values.Hostnames,
		)		
		
		return nil
	}
	
	err := app.Run(os.Args)
	return &values, err
}
