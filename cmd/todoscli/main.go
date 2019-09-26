package main

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"github.com/urfave/cli"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
	"todos-cnvt/owl"
	"todos-cnvt/todos"
)

func main() {
	app := cli.NewApp()
	app.Name = "owlconvercli"
	app.Version = "2019.5.1"
	app.Compiled = time.Now()
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Vitaliy Marchenko",
			Email: "vmarchenko@gmail.com",
		},
	}
	app.Copyright = "(c) 2019 VMarchenko"
	app.HelpName = "owlconvercli"
	app.Usage = "Convert different owl format"
	app.UsageText = "owlconvercli - use for convert owl documents in to different formats"

	app.Flags = []cli.Flag{

		cli.StringFlag{
			Name:  "src",
			Value: "pizza-functional.owl",
			Usage: "--src pizza-functional.owl",
		},
		cli.StringFlag{
			Name:  "dst",
			Value: "pizza-functional.xml",
			Usage: "--dst pizza-functional.xml",
		},
		cli.StringFlag{
			Name:  "cfg",
			Value: "config.toml",
			Usage: "--cfg config.toml",
		},
	}

	app.Action = func(c *cli.Context) error {

		dir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(dir)
		src := c.String("src")
		fmt.Println(" Source file", src)
		dst := c.String("dst")
		fmt.Println("Destination file", dst)

		srcExt := filepath.Ext(src)
		dstExt := filepath.Ext(dst)
		if len(srcExt) == 0 {
			fmt.Println("Source file not found", srcExt)
			return nil
		}
		if len(dstExt) == 0 {
			fmt.Println("Destination file not found", dstExt)
			return nil
		}
		if srcExt == dstExt {
			fmt.Println("Trying convert files same types", srcExt, dstExt)
			return nil
		}
		f, err := os.Open(src)
		if err != nil {
			fmt.Println("Open source file error!", err)
			return nil
		}
		defer f.Close()
		if srcExt == ".owl" {

		} else if srcExt == ".xml" {
			tod := todos.Ontology{}
			err := tod.OntologyParser(f)
			if err != nil {
				fmt.Println("Open source file error!", err)
				return nil
			}
			if dstExt == ".owl" {
				name := filepath.Base(dst)
				name = strings.Trim(name, dstExt)
				f, err := os.Create(dst)
				defer f.Close()

				if err != nil {
					fmt.Println(err)
					return err
				}
				w := bufio.NewWriter(f)
				_, err = w.WriteString(xml.Header)
				cfg:= c.String("cfg")
				err = owl.TodosFileBuilder(cfg, name, &tod, w)
				if err != nil {
					fmt.Println(err)
				}
			}
		}

		return nil
	}
	app.Run(os.Args)
}
