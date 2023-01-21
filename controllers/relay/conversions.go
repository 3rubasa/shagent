package relay

import "fmt"

const (
	relayStateNeutral RelayState = iota
	relayStateOn
	relayStateOff
)

func StringToRelayState(state string) (RelayState, error) {
	switch state {
	case "on":
		return relayStateOn, nil
	case "off":
		return relayStateOff, nil
	default:
		return relayStateNeutral, fmt.Errorf("unexpected value for relay state: %s", state)
	}
}

func RelayStateToString(state RelayState) (string, error) {
	switch state {
	case relayStateOn:
		return "on", nil
	case relayStateOff:
		return "off", nil
	default:
		return "", fmt.Errorf("unexpected relay state value")
	}
}
