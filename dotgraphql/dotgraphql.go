package dotgraphql

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
)

type GqlRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

type GqlError struct {
	Message string   `json:"message"`
	Path    []string `json:"path"`
}

type GqlResponse struct {
	Data   *json.RawMessage `json:"data,omitempty"`
	Errors *[]GqlError      `json:"errors,omitempty"`
}

type GqlClient struct {
	Endpoint string
	Headers  map[string]string
	queries  map[string]string
}

func NewGqlClient(endpointURL string, headers map[string]string) *GqlClient {
	return &GqlClient{
		Endpoint: endpointURL,
		Headers:  headers,
	}
}

func (gc *GqlClient) LoadFilesFromPath(path string, recursive bool) error {
	items, err := ioutil.ReadDir(path)
	if err != nil {
		return fmt.Errorf("failed to load from path %s with error: %v", path, err)
	}

	fails := make([]string, 0, len(items))

	for _, item := range items {
		if item.IsDir() && recursive {
			err = gc.LoadFilesFromPath(path+"/"+item.Name(), true)
		} else {
			if !(len(item.Name()) > 8 && item.Name()[len(item.Name())-8:] == ".graphql") {
				continue
			}
			err = gc.LoadFile(path + "/" + item.Name())
		}

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

// Delete this eventually...
func (gc *GqlClient) PrintLoadedFiles() {
	fmt.Println("Printing loaded files...")
	for name, file := range gc.queries {
		fmt.Printf("name: %s, content: %s", name, file)
	}
}

func (gc GqlClient) Exec(filename string, queryVars map[string]interface{}) (*GqlResponse, error) {
	query, ok := gc.queries[filename]
	if !ok {
		return nil, fmt.Errorf("this query has not been loaded yet. Please load query before attempting to use")
	}
	request := GqlRequest{
		Query:     query,
		Variables: queryVars,
	}
	bodyBytes, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal GraphQL request with error %v", err)
	}

	req, err := http.NewRequest("POST", gc.Endpoint, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("faield to create HTTP request with error %v", err)
	}

	for header, value := range gc.Headers {
		req.Header.Set(header, value)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to POST with error %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("server returned non-OK response status code %d with message %s", res.StatusCode, res.Status)
	}

	var gqlRes GqlResponse

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body of server request with error: %v", err)
	}

	err = json.Unmarshal(body, &gqlRes)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body into graphQL response with error %v", err)
	}

	return &gqlRes, nil
}

func (gc GqlClient) ExecAndBindResult(filename string, queryVars map[string]interface{}, result interface{}) error {
	if result == nil {
		return fmt.Errorf("received nil pointer as result. result must a pointer be of some type")
	}

	if reflect.ValueOf(result).Kind() != reflect.Ptr {
		return fmt.Errorf("received non-pointer type. result must be a pointer to some type")
	}
	res, err := gc.Exec(filename, queryVars)
	if err != nil {
		return fmt.Errorf("failed to execute GraphQL query with error: %v", err)
	}

	if res != nil {
		err = json.Unmarshal(*res.Data, result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal response data with error %v", err)
		}
	} else {
		return fmt.Errorf("received nil response data")
	}

	return nil
}
