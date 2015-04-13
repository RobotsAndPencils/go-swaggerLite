package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"go/ast"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/RobotsAndPencils/go-swaggerLite/markup"
	"github.com/RobotsAndPencils/go-swaggerLite/parser"
)

const (
	AVAILABLE_FORMATS = "go|markdown"
)

var apiPackage = flag.String("apiPackage", "", "The package that implements the API controllers, relative to $GOPATH/src")
var mainApiFile = flag.String("mainApiFile", "", "The file that contains the general API annotations, relative to $GOPATH/src")
var basePath = flag.String("basePath", "", "Web service base path")
var outputFormat = flag.String("format", "go", "Output format type for the generated files: "+AVAILABLE_FORMATS)
var output = flag.String("output", "generatedSwaggerSpec.go", "The opitonal name of the output file to be generated")
var generatedPackage = flag.String("package", "main", "The opitonal package name of the output file to be generated")

var generatedFileTemplate = `package {{generagedPackage}}
//This file is generated automatically. Do not edit it manually.

import (
	"net/http"
	"strings"
)

func SwaggerApiHandler(prefix string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resource := strings.TrimPrefix(r.URL.Path, prefix)
		resource = strings.Trim(resource, "/")

		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "accept, authorization, content-type")
			w.Header().Set("Access-Control-Max-Age", "1800")
			w.WriteHeader(204)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET")

		if resource == "" {
			w.Write([]byte(swaggerResourceListing))
			return
		}

		if json, ok := swaggerApiDescriptions[resource]; ok {
			w.Write([]byte(json))
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}
}


var swaggerResourceListing = {{resourceListing}}
var swaggerApiDescriptions = {{apiDescriptions}}
`

// It must return true if funcDeclaration is controller. We will try to parse only comments before controllers
// Stubbed out for now
func IsController(funcDeclaration *ast.FuncDecl) bool {
	return true
}

func generateSwaggerDocs(parser *parser.Parser) {
	fd, err := os.Create(path.Join("./", *output))
	if err != nil {
		log.Fatalf("Can not create document file: %v\n", err)
	}
	defer fd.Close()

	var apiDescriptions bytes.Buffer
	for apiKey, apiDescription := range parser.TopLevelApis {
		apiDescriptions.WriteString("\"" + apiKey + "\":")

		apiDescriptions.WriteString("`")
		json, err := json.MarshalIndent(apiDescription, "", "    ")
		if err != nil {
			log.Fatalf("Can not serialise []ApiDescription to JSON: %v\n", err)
		}
		apiDescriptions.Write(json)
		apiDescriptions.WriteString("`,")
	}

	doc := strings.Replace(generatedFileTemplate, "{{resourceListing}}", "`"+string(parser.GetResourceListingJson())+"`", -1)
	doc = strings.Replace(doc, "{{apiDescriptions}}", "map[string]string{"+apiDescriptions.String()+"}", -1)
	doc = strings.Replace(doc, "{{generagedPackage}}", *generatedPackage, -1)

	fd.WriteString(doc)
}

func InitParser() *parser.Parser {
	parser := parser.NewParser()

	parser.BasePath = *basePath
	parser.IsController = IsController

	parser.TypesImplementingMarshalInterface["NullString"] = "string"
	parser.TypesImplementingMarshalInterface["NullInt64"] = "int"
	parser.TypesImplementingMarshalInterface["NullFloat64"] = "float"
	parser.TypesImplementingMarshalInterface["NullBool"] = "bool"

	return parser
}

func main() {
	flag.Parse()

	if *mainApiFile == "" {
		*mainApiFile = *apiPackage + "/main.go"
	}
	if *apiPackage == "" {
		flag.PrintDefaults()
		return
	}

	parser := InitParser()
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		log.Fatalf("Please, set $GOPATH environment variable\n")
	}

	log.Println("Start parsing")
	gopaths := filepath.SplitList(gopath)
	var err error
	var errs string
	for _, gop := range gopaths {
		err = parser.ParseGeneralAPIInfo(path.Join(gop, "src", *mainApiFile))
		if err != nil {
			errs += fmt.Sprintf("    %s\n", err)
		} else {
			break
		}
	}
	if err != nil {
		log.Fatalf("Error locating main API File:\n%s", errs)
	}
	parser.ParseApi(*apiPackage)
	log.Println("Finish parsing")

	format := strings.ToLower(*outputFormat)
	switch format {
	case "go":
		generateSwaggerDocs(parser)
		log.Println("Doc file generated")
	case "markdown":
		markup.GenerateMarkup(parser, new(markup.MarkupMarkDown), output, ".md")
		log.Println("MarkDown file generated")
	default:
		log.Fatalf("Invalid -format specified. Must be one of %v.", AVAILABLE_FORMATS)
	}

}
