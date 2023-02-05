package network

import "strings"

const DefaultPort string = "6660"

type Config struct {
	Addr         string
	Port         string
	KeyBase      string
	InstancePath string
	RelayPath    string
}

func NewConfig(base, addr string) Config {
	split := strings.SplitN(addr, ":", 2)
	port := DefaultPort
	if len(split) > 1 {
		port = split[1]
	}
	return Config{
		Addr: addr,
		Port: port,
		// TODO proper key
		KeyBase:      "API-KEY-BASE",
		RelayPath:    "/" + base,
		InstancePath: "/_" + base,
	}
}
