package main

import (
	"os"

	"gitlab.com/shipink/common/krakend/parser"
)

func main() {
	env := os.Getenv("ENVIRONMENT")

	if env == "dev" {
		files, err := os.ReadDir("_endpoints/dev/")
		if err != nil {
			panic(err)
		}
		for _, file := range files {
			parser.Parse("_endpoints/dev/" + file.Name())
		}

		parser.Concat("dev", "dev")
	}

	if env == "test" {
		files, err := os.ReadDir("_endpoints/test/")
		if err != nil {
			panic(err)
		}
		for _, file := range files {
			parser.Parse("_endpoints/test/" + file.Name())
		}

		parser.Concat("test", "test")
	}

	if env == "prod" {
		files, err := os.ReadDir("_endpoints/prod/")
		if err != nil {
			panic(err)
		}
		for _, file := range files {
			parser.Parse("_endpoints/prod/" + file.Name())
		}

		parser.Concat("prod", "prod")
	}

}
