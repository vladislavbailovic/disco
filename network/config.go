package network

import "strings"

type Config struct {
	Addr         string
	Port         string
	InstancePath string
	RelayPath    string
}

func NewConfig(base, addr string) Config {
	split := strings.SplitN(addr, ":", 2)
	return Config{
		Addr:         addr,
		Port:         split[1],
		InstancePath: "/" + base,
		RelayPath:    "/_" + base,
	}
}
