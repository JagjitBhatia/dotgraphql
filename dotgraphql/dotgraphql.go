package dotgraphql

import (
	"fmt"
	"io/ioutil"
)

type GqlClient struct {
	queries map[string]string
}

func (gc *GqlClient) LoadFilesFromPath(path string, recursive bool) error {
	items, err := ioutil.ReadDir(path)
	if err != nil {
		return fmt.Errorf("failed to load from path %s with error: %v", path, err)
	}

	fails := make([]string, 0, len(items))

	for _, item := range items {
		if item.IsDir() && recursive {
			return gc.LoadFilesFromPath(path+"/"+item.Name(), true)
		}

		if !(len(item.Name()) > 8 && item.Name()[len(item.Name())-8:] == ".graphql") {
			continue
		}
		err = gc.LoadFile(path + "/" + item.Name())
		if err != nil {
			fails = append(fails, path+"/"+item.Name())
		}
	}

	if len(fails) > 0 {
		return fmt.Errorf("error loading from path %s: failed to load the following graphQL files: %v", path, fails)
	}

	return nil
}

func (gc *GqlClient) LoadFile(path string) error {
	fileBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to load file %s with error %v", path, err)
	}

	if gc.queries == nil {
		gc.queries = make(map[string]string)
	}

	gc.queries[path] = string(fileBytes)

	return nil
}

func (gc *GqlClient) PrintLoadedFiles() {
	fmt.Println("Printing loaded files...")
	for name, file := range gc.queries {
		fmt.Printf("name: %s, content: %s", name, file)
	}
}
