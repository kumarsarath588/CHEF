package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	chefc "github.com/go-chef/chef"
)

//Resource struct defines node attributes and actions
type Resource struct {
	Node   chefc.Node
	Action string
}

//ChefConfig is chef config struct
type ChefConfig struct {
	Provider chefc.Config
	Resource Resource
	client   *chefc.Client
}

//NewClient creates new client for chef server
func NewClient(c *ChefConfig) (*chefc.Client, error) {
	config := &chefc.Config{
		Name:    c.Provider.Name,
		Key:     c.Provider.Key,
		SkipSSL: c.Provider.SkipSSL,
		BaseURL: c.Provider.BaseURL,
	}
	client, err := chefc.NewClient(config)
	if err != nil {
		return nil, err
	}
	return client, nil
}

//DeleteNode is used to delete the give chef node
func (c *ChefConfig) DeleteNode() error {
	n, err := c.client.Nodes.Get(c.Resource.Node.Name)
	if n.Name == "" {
		return fmt.Errorf("Node Dosen't exists")
	}
	if err != nil {
		return err
	}
	if err := c.client.Nodes.Delete(n.Name); err != nil {
		return err
	}
	if err := c.client.Clients.Delete(n.Name); err != nil {
		return err
	}
	return nil
}

//ChefWebhookHomePage is homepage for chef
func ChefWebhookHomePage(w http.ResponseWriter, r *http.Request) {
	config := new(ChefConfig)
	err := json.NewDecoder(r.Body).Decode(config)
	if err != nil {
		panic(err)
	}
	key, err := ioutil.ReadFile("/go/bin/key.pem")
	if err != nil {
		fmt.Println("Couldn't read key file:", err)
		os.Exit(1)
	}
	config.Provider.Key = string(key)
	client, err := NewClient(config)
	if err != nil {
		fmt.Printf("Unable to connect to client")
	}
	config.client = client
	if config.Resource.Action == "Delete" {
		if err := config.DeleteNode(); err != nil {
			fmt.Printf("%s Deletion failed: %s\n", config.Resource.Node.Name, err)
		}
	}
}

func main() {
	http.HandleFunc("/", ChefWebhookHomePage)
	http.ListenAndServe(":3002", nil)
}

