package main

import (
	"flag"
	"go/build"
	"log"
	"os"
)

const templateDir = "github.com/zabawaba99/cachegen/template"

var (
	keyType   string
	valueType string
)

func init() {
	flag.StringVar(&keyType, "key-type", "", "The type that will be used to add and retrieve items from the cache")
	flag.StringVar(&valueType, "value-type", "", "The type that will be cached")
	flag.Parse()

	if keyType == "" {
		log.Fatal("Must specify '-key-type'")
	}

	if valueType == "" {
		log.Fatal("Must specify '-value-type'")
	}
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("cachegen: ")

	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Could not fetch working dir %v", err)
	}

	pkg := getPackage()
	generate(wd, pkg)
}

// getPackage finds the package name of all the go files in the
// current working directory
func getPackage() string {
	p, err := build.Default.Import(".", ".", build.ImportMode(0))
	if err != nil {
		log.Fatalf("Could not identify package: %v", err)
	}
	return p.Name
}
