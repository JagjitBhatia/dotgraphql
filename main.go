package main

import (
	"github.com/JagjitBhatia/dotgraphql"
)

func main() {
	client := &dotgraphql.GqlClient{}

	client.LoadFilesFromPath(".", false)

	client.PrintLoadedFiles()

}
