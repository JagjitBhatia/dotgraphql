package main

import (
	"fmt"

	"github.com/JagjitBhatia/dotgraphql/dotgraphql"
)

type Member struct {
	User  *User  `json:"user"`
	Role  string `json:"role"`
	Title string `json:"title"`
}

type Org struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Institution string    `json:"institution"`
	OrgPicURL   *string   `json:"org_pic_url"`
	Members     []*Member `json:"members"`
}

type User struct {
	ID          string  `json:"id"`
	Username    string  `json:"username"`
	FirstName   string  `json:"firstName"`
	LastName    string  `json:"lastName"`
	Institution string  `json:"institution"`
	PfpURL      *string `json:"pfp_url"`
}

type Response struct {
	Orgs []Org `json:"orgs"`
}

func main() {
	client := dotgraphql.NewGqlClient("http://localhost:8080/query", map[string]string{})

	client.LoadFile("testFiles/test1.graphql")

	var res Response

	err := client.ExecAndBindResult("testFiles/test1.graphql", map[string]interface{}{}, &res)
	if err != nil {
		fmt.Printf("Error executing query: %v\n", err)
	} else {
		fmt.Printf("Received response: %#v\n", res)
	}

	//client.PrintLoadedFiles()

}
