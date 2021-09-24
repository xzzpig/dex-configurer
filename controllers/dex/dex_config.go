package dex

import (
	v1 "github.com/xzzpig/dex-configurer/apis/dex/v1"
	"gopkg.in/yaml.v2"
)

type DexConfigMap = map[string]interface{}

type DexConfig struct {
	m DexConfigMap
}

type StaticClient = v1.DexAuthClientSpec

func GetConfigFromMap(m DexConfigMap) *DexConfig {
	return &DexConfig{m}
}

func GetConfigFromString(s string) (*DexConfig, error) {
	m := make(DexConfigMap)

	err := yaml.Unmarshal([]byte(s), &m)
	if err != nil {
		return nil, err
	}

	return GetConfigFromMap(m), nil
}

func GetConfigFromBytes(s []byte) (*DexConfig, error) {
	m := make(DexConfigMap)

	err := yaml.Unmarshal(s, &m)
	if err != nil {
		return nil, err
	}

	return GetConfigFromMap(m), nil
}

func (config *DexConfig) GetStaticClients() ([]StaticClient, error) {
	mm := config.m["staticClients"]
	if mm == nil {
		return make([]StaticClient, 0), nil
	}
	m := mm.([]interface{})
	d, err := yaml.Marshal(m)
	if err != nil {
		return nil, err
	}

	staticClients := make([]StaticClient, 0)
	err = yaml.Unmarshal(d, &staticClients)
	if err != nil {
		return nil, err
	}
	return staticClients, nil
}

func (config *DexConfig) SetStaticClients(clients []StaticClient) *DexConfig {
	config.m["staticClients"] = clients
	return config
}

func (config *DexConfig) ToString() (string, error) {
	out, err := yaml.Marshal(config.m)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func (config *DexConfig) GoString() (s string) {
	s, _ = config.ToString()
	return
}

func (config *DexConfig) SetStaticClient(client *StaticClient) error {
	clients, err := config.GetStaticClients()
	if err != nil {
		return err
	}
	found := false
	for i, c := range clients {
		if c.Id != client.Id {
			continue
		}
		found = true
		client.DeepCopyInto(&clients[i])
	}
	if !found {
		clients = append(clients, *client)
	}
	config.SetStaticClients(clients)
	return nil
}

func (config *DexConfig) RemoveStaticClient(id string) error {
	clients, err := config.GetStaticClients()
	if err != nil {
		return err
	}
	if len(clients) == 0 {
		return nil
	}
	index := 0
	for index < len(clients) {
		if clients[index].Id == id {
			clients = append(clients[:index], clients[index+1:]...)
			continue
		}
		index++
	}
	config.SetStaticClients(clients)
	return nil
}
