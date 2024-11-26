package frr

import (
	"fmt"
	"magic-wan/pkg"
	"magic-wan/pkg/osUtil"
	"strings"
)

const emptyFRRLine = "!"

type Network struct {
	cidr string
	area string
}

type OSPFRouter struct {
	routerId string
	networks []*Network
	areas    []*Network
}

func newOSPFRouter(routerId string) *OSPFRouter {
	return &OSPFRouter{
		routerId: routerId,
		networks: make([]*Network, 0),
		areas:    make([]*Network, 0),
	}
}

func (r *OSPFRouter) AddNetwork(cidr, area string) *OSPFRouter {
	r.networks = append(r.networks, &Network{cidr: cidr, area: area})
	return r
}
func (r *OSPFRouter) AddArea(cidr, area string) *OSPFRouter {
	r.areas = append(r.areas, &Network{cidr: cidr, area: area})
	return r
}

func (r *OSPFRouter) buildConfigText() string {
	configLines := make([]string, 0)

	configLines = append(configLines, "router ospf")

	configLines = append(configLines, fmt.Sprintf(" ospf router-id %s", r.routerId), " "+emptyFRRLine)

	// default-disable ospf on all interfaces
	configLines = append(configLines, " passive-interface default")

	/*
		Looks like the below code only produces Problems in combination with "systemctl reload frr"
		Since it also does not seem to be required let's just not add it... ^^

		for _, network := range r.networks {
			configLines = append(configLines, fmt.Sprintf(" network %s area %s", network.cidr, network.area))
		}
	*/
	for _, area := range r.areas {
		// this line allows for route summarization
		configLines = append(configLines, fmt.Sprintf(" area %s range %s", area.area, area.cidr))
	}

	configLines = append(configLines, "exit", emptyFRRLine)
	return strings.Join(configLines, "\n")
}

type Interface struct {
	name    string
	area    string
	passive bool
}

func (i *Interface) buildConfigText() string {
	configLines := make([]string, 0)

	configLines = append(configLines, fmt.Sprintf("interface %s", i.name))

	if i.passive {
		configLines = append(configLines, " ip ospf passive")
	} else {
		configLines = append(configLines, " no ip ospf passive")
	}
	configLines = append(configLines, fmt.Sprintf(" ip ospf area %s", i.area))

	configLines = append(configLines, "exit", emptyFRRLine)
	return strings.Join(configLines, "\n")
}

type Config struct {
	logging    []string
	hostname   string
	Router     *OSPFRouter
	interfaces []*Interface
}

func NewConfig(hostname, routerId string) *Config {
	return &Config{
		logging:    make([]string, 0),
		hostname:   hostname,
		Router:     newOSPFRouter(routerId),
		interfaces: make([]*Interface, 0),
	}
}

func (c *Config) AddLogging(destination, level string) *Config {
	c.logging = append(c.logging, fmt.Sprintf("%s %s", destination, level))
	return c
}

func (c *Config) AddInterface(name, area string, passive bool) *Config {
	c.interfaces = append(c.interfaces, &Interface{
		name:    name,
		area:    area,
		passive: passive,
	})
	return c
}

func (c *Config) buildConfigText() string {
	// IMPROVMENT: allow for broadcasting (+ summarizing) additional networks
	// IMPROVMENT: allow for including additional interfaces
	// IMPROVMENT: allow for SNAT / DNAT to be done for external interfaces

	configLines := make([]string, 0)

	// header infos
	for _, log := range c.logging {
		configLines = append(configLines, fmt.Sprintf("log %s", log))
	}
	configLines = append(configLines, "frr defaults traditional")
	configLines = append(configLines, fmt.Sprintf("hostname %s", c.hostname))
	configLines = append(configLines, emptyFRRLine)

	// main router config
	configLines = append(configLines, c.Router.buildConfigText())
	configLines = append(configLines, emptyFRRLine)

	// config for the individual interfaces
	for _, iface := range c.interfaces {
		configLines = append(configLines, iface.buildConfigText())
	}
	configLines = append(configLines, emptyFRRLine)

	return strings.Join(configLines, "\n")
}

func (c *Config) WriteConfFile() error {
	return osUtil.WriteFile(DEFAULT_CONFIG_PATH, c.buildConfigText())
}

func (c *Config) PushToFrr() error {
	err := c.WriteConfFile()
	if err != nil {
		return err
	}
	return c.RestartFRR()
}

func (c *Config) RestartFRR() error {
	return pkg.FrrService.Reload()
}

func (c *Config) StartFRR() error {
	err := c.WriteConfFile()
	if err != nil {
		return err
	}
	return pkg.FrrService.Start()
}
