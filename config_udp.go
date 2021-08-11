package traefik_mesos_provider

import (
	"net"
	"strconv"

	"github.com/traefik/genconf/dynamic"
)

// buildUDPServiceConfiguration buid the UDP Service of the Mesos Taks
// containerName.
func (p *Provider) buildUDPServiceConfiguration(nr int, configuration *dynamic.UDPConfiguration) {
	if len(configuration.Routers) == 0 {
		return
	}
	if len(configuration.Services) == 0 {
		configuration.Services = make(map[string]*dynamic.UDPService)
	}

	for _, service := range configuration.Routers {
		// search all different ports by name and create a Loadbalancer configuration for traefik
		task := p.mesosConfig[nr]
		if len(task.Discovery.Ports.Ports) > 0 {
			for _, port := range task.Discovery.Ports.Ports {
				if len(port.Name) == 0 || port.Protocol != "udp" {
					continue
				}
				if port.Name != service.Service {
					continue
				}
				lb := &dynamic.UDPServersLoadBalancer{}
				lb.Servers = p.getUDPServers(port.Name, nr)

				lbService := &dynamic.UDPService{
					LoadBalancer: lb,
				}

				configuration.Services[service.Service] = lbService
			}
		}
	}
}

// getUDPServers search all IP addresses to the given portName of
// the Mesos Task with the containerName.
func (p *Provider) getUDPServers(portName string, nr int) []dynamic.UDPServer {
	var servers []dynamic.UDPServer
	name := p.mesosConfig[nr].Name
	for _, task := range p.mesosConfig {
		// ever take the first IP in the list
		ip := task.Statuses[0].ContainerStatus.NetworkInfos[0].IPAddresses[0].IPAddress
		if name == task.Name && len(task.Discovery.Ports.Ports) > 0 {
			for _, port := range task.Discovery.Ports.Ports {
				if portName != port.Name || port.Protocol != "udp" {
					continue
				}
				po := strconv.Itoa(port.Number)
				server := dynamic.UDPServer{
					Address: net.JoinHostPort(ip, po),
				}
				servers = append(servers, server)
			}
		}
	}
	return servers
}
