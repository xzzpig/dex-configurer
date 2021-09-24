package dex_test

import (
	"fmt"
	"testing"

	"github.com/xzzpig/dex-configurer/controllers/dex"
)

func TestSetClient(t *testing.T) {
	var data = `
enablePasswordDB: true
issuer: https://auth.example.com
staticClients:
- id: example-app
  name: Example App
  redirectURIs:
  - http://127.0.0.1:5555/callback
  secret: ZXhhbXBsZS1hcHAtc2VjcmV0
- id: gitea
  name: Gitea Example
  redirectURIs:
  - https://git.example.com/user/oauth2/auth/callback
  secret: ZXhhbXBsZS1hcHAtc2VjcmV0
storage:
  config:
    inCluster: true
  type: kubernetes
`

	config, err := dex.GetConfigFromString(data)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	fmt.Println(config.GoString())
	fmt.Println("========================================")

	staticClients, err := config.GetStaticClients()
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	c := staticClients[0].DeepCopy()
	c.Name = "Test"
	c.Id = "test"
	config.SetStaticClient(c)
	fmt.Println(config.GoString())
	fmt.Println("========================================")
	config, err = dex.GetConfigFromString(config.GoString())
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	config.RemoveStaticClient("gitea")
	fmt.Println(config.GoString())
}

func TestSetClient2(t *testing.T) {
	var data = `
enablePasswordDB: true
issuer: https://auth.example.com
storage:
  config:
    inCluster: true
  type: kubernetes
`

	config, err := dex.GetConfigFromString(data)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	fmt.Println(config.GoString())
	fmt.Println("========================================")

	c := dex.StaticClient{
		Id:           "test",
		Name:         "Test",
		Secret:       "aaa",
		RedirectURIs: []string{"https://git.fae2ly.com:8443/user/oauth2/auth/callback"},
	}
	config.SetStaticClient(&c)
	fmt.Println(config.GoString())
	fmt.Println("========================================")
	config, err = dex.GetConfigFromString(config.GoString())
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	config.RemoveStaticClient("gitea")
	fmt.Println(config.GoString())
	fmt.Println("========================================")
	config, err = dex.GetConfigFromString(config.GoString())
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	config.RemoveStaticClient("test")
	fmt.Println(config.GoString())
}
