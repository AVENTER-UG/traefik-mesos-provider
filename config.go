package traefik_mesos_provider

import (
	"encoding/json"
	"fmt"

	"github.com/traefik/genconf/dynamic"
	"github.com/traefik/genconf/dynamic/tls"
)

func (p *Provider) buildConfiguration() *dynamic.Configuration {
	configuration := &dynamic.Configuration{
		HTTP: &dynamic.HTTPConfiguration{
			Routers:           make(map[string]*dynamic.Router),
			Middlewares:       make(map[string]*dynamic.Middleware),
			Services:          make(map[string]*dynamic.Service),
			ServersTransports: make(map[string]*dynamic.ServersTransport),
		},
		TCP: &dynamic.TCPConfiguration{
			Routers:  make(map[string]*dynamic.TCPRouter),
			Services: make(map[string]*dynamic.TCPService),
		},
		TLS: &dynamic.TLSConfiguration{
			Stores:  make(map[string]tls.Store),
			Options: make(map[string]tls.Options),
		},
		UDP: &dynamic.UDPConfiguration{
			Routers:  make(map[string]*dynamic.UDPRouter),
			Services: make(map[string]*dynamic.UDPService),
		},
	}

	for i, tasks := range p.mesosConfig {
		task := tasks
		// The first Task is the leading one

		p.buildTCPServiceConfiguration(i, configuration.TCP)
		p.buildUDPServiceConfiguration(i, configuration.UDP)
		p.buildHTTPServiceConfiguration(i, configuration.HTTP)

		res2B, _ := json.Marshal(configuration)
		fmt.Println(string(res2B))
	}

	return configuration
}
