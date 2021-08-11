// Package traefik_mesos_provider is a traefik provider plugin for Apache Mesos
package traefik_mesos_provider

import (
	"context"
	cTls "crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/m7shapan/njson"
	"github.com/traefik/genconf/dynamic"
	ptypes "github.com/traefik/paerser/types"
)

// DefaultTemplateRule The default template for the default rule.
const DefaultTemplateRule = "Host(`{{ normalize .Name }}`)"

// Provider holds configuration of the provider.
type Provider struct {
	Endpoint     string          `description:"Mesos server endpoint. You can also specify multiple endpoint for Mesos"`
	Principal    string          `Description:"Principal to authorize against Mesos Manager"`
	Secret       string          `Description:"Secret authorize against Mesos Manager"`
	PollInterval ptypes.Duration `description:"Polling interval for endpoint." json:"pollInt"`
	PollTimeout  ptypes.Duration `description:"Polling timeout for endpoint." json:"pollTime"`
	mesosConfig  []*MesosTask
	cancel       func()
	name         string
}

// New creates a new Provider plugin.
func New(ctx context.Context, config *Config, name string) (*Provider, error) {
	return &Provider{
		Endpoint:     config.Endpoint,
		PollInterval: ptypes.Duration(time.Second),
		PollTimeout:  ptypes.Duration(time.Second),
		Secret:       config.Secret,
		Principal:    config.Principal,
		name:         name,
	}, nil
}

// Config the plugin configuration.
type Config struct {
	Endpoint     string `description:"Mesos server endpoint. You can also specify multiple endpoint for Mesos"`
	PollInterval string `json:"pollInterval,omitempty"`
	Principal    string `Description:"Principal to authorize against Mesos Manager"`
	Secret       string `Description:"Secret authorize against Mesos Manager"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		PollInterval: "10s",
	}
}

// Init the provider.
func (p *Provider) Init() error {
	p.name = "Apache Mesos Provider"
	p.mesosConfig = []*MesosTask{}
	return nil
}

// Provide creates and send dynamic configuration.
func (p *Provider) Provide(cfgChan chan<- json.Marshaler) error {
	ctx, cancel := context.WithCancel(context.Background())
	p.cancel = cancel

	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Print(err)
			}
		}()

		p.loadConfiguration(ctx, cfgChan)
	}()

	return nil
}

// Stop to stop the provider and the related go routines.
func (p *Provider) Stop() error {
	p.cancel()
	return nil
}

func (p *Provider) loadConfiguration(ctx context.Context, cfgChan chan<- json.Marshaler) {
	ticker := time.NewTicker(time.Duration(p.PollInterval))
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			data, _ := p.getTasks()

			var tasks MesosTasks
			if err := njson.Unmarshal(data, &tasks); err != nil {
				log.Print(err)
			}

			// collect all mesos tasks and combine the belong one.
			for _, task := range tasks.Tasks {
				res2B, _ := json.Marshal(task)
				fmt.Println(string(res2B))
				log.Print(task.State)
				if task.State == "TASK_RUNNING" {
					if p.checkTraefikLabels(task) {
						mesosCfg := &MesosTask{}
						mesosCfg = &task
						p.mesosConfig = append(p.mesosConfig, mesosCfg)
					}
				}
			}

			// build the treafik configuration
			if len(p.mesosConfig) > 0 {
				configuration := p.buildConfiguration()
				cfgChan <- &dynamic.JSONPayload{Configuration: configuration}
			}

			// cleanup old data
			p.mesosConfig = []*MesosTask{}
		case <-ctx.Done():
			return
		}
	}
}

func (p *Provider) checkTraefikLabels(task MesosTask) bool {
	for _, label := range task.Labels {
		if strings.Contains(label.Key, "traefik.") {
			return true
		}
	}
	return false
}

func (p *Provider) getTasks() ([]byte, error) {
	client := &http.Client{}
	client.Transport = &http.Transport{
		TLSClientConfig: &cTls.Config{InsecureSkipVerify: true},
	}
	req, _ := http.NewRequest("GET", p.Endpoint+"/tasks?order=asc&limit=-1", nil)
	req.SetBasicAuth(p.Principal, p.Secret)
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		log.Print(err)
	}

	if res.StatusCode != http.StatusOK {
		log.Printf("received non-ok response code: %d", res.StatusCode)
		return nil, err
	}

	return io.ReadAll(res.Body)
}
