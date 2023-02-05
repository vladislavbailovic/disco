package network

import "strings"

const DefaultHost string = "127.0.0.1"
const DefaultPort string = "6660"

type Config struct {
	Host         string
	Port         string
	KeyBase      string
	InstancePath string
	RelayPath    string
}

func NewConfig(base, addr string) Config {
	host := DefaultHost
	port := DefaultPort
	split := strings.SplitN(addr, ":", 2)
	if len(split) > 0 {
		host = split[0]
	}
	if len(split) > 1 {
		port = split[1]
	}
	return Config{
		Host: host,
		Port: port,
		// TODO proper key
		KeyBase:      "API-KEY-BASE",
		RelayPath:    "/" + base,
		InstancePath: "/_" + base,
	}
}
