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


type configs struct {
	File      string    `yaml:"File"`
	Port      int       `yaml:"Port"`
	Hostnames []string  `yaml:"Hostnames"`
}

var Values configs

var flags = []cli.Flag {
	cli.StringFlag{
		Name: "config, c",
		Usage: "YAML config `FILE`.",
		Destination: &Values.File,
	},
	cli.IntFlag{
		Name: "Port, p",
		Usage: "`PORT` to listen on.",
		Value: 80,
		Destination: &Values.Port,
	},
	cli.StringSliceFlag{
		Name: "Hostnames, n",
		Usage: "`HOSTNAME`s to listen to.",
		Value: &cli.StringSlice{"localhost"},
	},
}

func LoadFile(c *configs) error {
	if Values.File == "" {
		// No config file was provided...
		log.Printf("No config file provided...\n")
		return nil
	}
	log.Printf("Loading config file from %v", Values.File)
	yamlFile, err := ioutil.ReadFile(Values.File)
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

func UpdateValues(fileConfigs *configs, c *cli.Context) {
	fileFields := reflect.TypeOf(*fileConfigs)
	fileValues := reflect.ValueOf(*fileConfigs)
	vv := reflect.ValueOf(&Values).Elem()
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

func Init() error {
	app        := cli.NewApp()
	app.Name    = "Book Review"
	app.Usage   = "Better book reviews!"
	app.Version = "0.0.1"
	app.Flags   = flags

	app.Action = func(c *cli.Context) error {
		var yamlConfigs configs
		LoadFile(&yamlConfigs)

		// Override config values passed as params
		Values.Hostnames = c.StringSlice("Hostnames")
		// XXX: StringSlice hack exists due to:
		// https://github.com/urfave/cli/issues/790
		if c.IsSet("Hostnames") || c.IsSet("n") {
			Values.Hostnames = Values.Hostnames[1:]
		} else if len(yamlConfigs.Hostnames) > 0 {
			Values.Hostnames = yamlConfigs.Hostnames
		}

		UpdateValues(&yamlConfigs, c)
		log.Printf(
			"Configs loaded:\nPort: %d\nHostnames: %v\n",
			Values.Port,
			Values.Hostnames,
		)		
		
		return nil
	}
	
	err := app.Run(os.Args)
	return err
}
