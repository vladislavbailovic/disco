package network

import "disco/logging"

type DiscoveryStatus uint8

const (
	Init DiscoveryStatus = iota
	EstablishingQuorum
	Ready
)

func (x DiscoveryStatus) String() string {
	switch x {
	case Init:
		return "Initialized Discovery"
	case EstablishingQuorum:
		return "Establishing Quorum"
	case Ready:
		return "Ready"
	}
	logging.Get().Fatal("Unknown discovery status: %d", x)
	panic("Unknown discovery status")
}
