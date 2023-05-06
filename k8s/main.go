package main

import (
	"embed"
	"fmt"
	"os"
	"encoding/base64"

	"github.com/flosch/pongo2/v6"
	"github.com/joho/godotenv"
)

var version string = "v1.0.0-alpha"

func init() {
	if _, err := os.Stat(".env"); err == nil {
		err = godotenv.Load()
		if err != nil {
			os.Stderr.WriteString(err.Error())
		}
	}
	if _, err := os.Stat(".env.mysql"); err == nil {
		err = godotenv.Load(".env.mysql")
		if err != nil {
			os.Stderr.WriteString(err.Error())
		}
	}
}

//go:embed files
var files embed.FS
//go:embed templates
var templates embed.FS

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "version" || os.Args[1] == "-v" || os.Args[1] == "--version") {
		fmt.Println(version)
		return
	}

	var context pongo2.Context

	if os.Getenv("BASE_IMAGE") != "rayshoo/petclinic-base" {
		// cm.yaml
		fmt.Println(readEmbedFile("files/cm.yaml"))
		// job.yaml
		context = getContext([]string{"BASE_IMAGE","BASE_IMAGE_TAG","BASE_IMAGE_PUSH_SECRET"})
		fmt.Println(templating("templates/job.yaml",context))
	}
	// secert.yaml
	context = getBase64EncodedContext([]string{
		"BASE_IMAGE",
		"BASE_IMAGE_TAG",
		"WAS_IMAGE",
		"WAS_IMAGE_TAG",
		"MYSQL_PORT",
		"MYSQL_DATABASE",
		"MYSQL_USER",
		"MYSQL_PASS",
		"MYSQL_ROOT_PASSWORD",
	})
	context["MYSQL_URL"] = base64.StdEncoding.EncodeToString([]byte("petclinic-" + os.Getenv("MYSQL_URL")))
	fmt.Println(templating("templates/secret.yaml",context))

	// deploy.yaml
	context = getContext([]string{"WAS_IMAGE","WAS_IMAGE_TAG","WAS_IMAGE_PULL_SECRET"})
	fmt.Println(templating("templates/deploy.yaml",context))

	// svc.yaml
	context = getContext([]string{"MYSQL_URL","MYSQL_PORT"})
	fmt.Println(templating("templates/svc.yaml",context))

	// ing.yaml	
	fmt.Println(readEmbedFile("files/ing.yaml"))
}

func getBase64EncodedContext(keys []string) pongo2.Context {
	context := pongo2.Context{}
	for i := range keys {
		context[keys[i]] = base64.StdEncoding.EncodeToString([]byte(os.Getenv(keys[i])))
	}
	return context
}

func getContext(keys []string) pongo2.Context {
	context := pongo2.Context{}
	for i := range keys {
		context[keys[i]] = os.Getenv(keys[i])
	}
	return context
}

func readEmbedFile(filePath string) string {
	fileString, err := files.ReadFile(filePath)
	if err != nil {
		handleError(err)
	}
	return string(fileString)
}

func templating(filePath string, context pongo2.Context) string {
	templateString, err := templates.ReadFile(filePath)
	if err != nil {
		handleError(err)
	}
	template, err := pongo2.FromString(string(templateString))
	if err != nil {
		handleError(err)
	}
	result, err := template.Execute(context)
	if err != nil {
		handleError(err)
	}
	return result
}

func handleError(err error) {
	os.Stderr.WriteString(err.Error())
	os.Exit(1)
}