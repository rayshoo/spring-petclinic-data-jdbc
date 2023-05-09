package main

import (
	"embed"
	"encoding/base64"
	"fmt"
	"os"
	"strconv"

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

	envs := []string{
		"IMAGE",
		"IMAGE_TAG",
		"WAS_PORT",
		"IMAGE_REPO_SECRET",
		"MYSQL_URL",
		"MYSQL_PORT",
		"MYSQL_DATABASE",
		"MYSQL_USER",
		"MYSQL_PASS",
		"MYSQL_ROOT_PASSWORD",
	}

	for i := range envs {
		if os.Getenv(envs[i]) == "" {
			handleError(fmt.Errorf("The " + envs[i] + " environment variable needs a value."))
		}
	}
	portCheck("WAS_PORT")
	portCheck("MYSQL_PORT")

	context := getContext(&envs)
	b64Context := getBase64EncodedContext(&context)

	// cm.yaml
	fmt.Println(templating("templates/cm.yaml", context))
	// job.yaml
	fmt.Println(templating("templates/job.yaml", context))

	// secert.yaml
	b64Context["MYSQL_URL"] = base64.StdEncoding.EncodeToString([]byte("petclinic-" + os.Getenv("MYSQL_URL")))
	fmt.Println(templating("templates/secret.yaml", b64Context))

	// pvc.yaml
	fmt.Println(readEmbedFile("files/pvc.yaml"))

	// deploy.yaml
	fmt.Println(templating("templates/deploy.yaml", context))

	// svc.yaml
	fmt.Println(templating("templates/svc.yaml", context))

	// ing.yaml
	fmt.Println(readEmbedFile("files/ing.yaml"))
}

func getBase64EncodedContext(context *pongo2.Context) pongo2.Context {
	result := pongo2.Context{}
	for k, v := range *context {
		result[k] = base64.StdEncoding.EncodeToString([]byte(v.(string)))
	}
	return result
}

func getContext(keys *[]string) pongo2.Context {
	context := pongo2.Context{}
	for i := range *keys {
		context[(*keys)[i]] = os.Getenv((*keys)[i])
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

func portCheck(key string) {
	_, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		handleError(fmt.Errorf("You must assign a numeric value to the environment variable '" + key + "'. - current value: " + os.Getenv(key)))
	}
}

func handleError(err error) {
	os.Stderr.WriteString(err.Error() + "\n")
	os.Exit(1)
}
