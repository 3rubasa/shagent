package sonoffr3rf

const (
	//relaymDNSName = "eWeLink_10012ff7ab.local" // Fails to resolve for some reason
	relayPort = 8081
)

const (
	relayInfoPath   = "zeroconf/info"
	relaySwitchPath = "zeroconf/switch"
)

const (
	relayInfoBody = `{ 
		"deviceid": "", 
		"data": { } 
	 }`

	relaySwitchOnBody = `{ 
		"deviceid": "", 
		"data": {
			"switch": "on" 
		} 
	 }`

	relaySwitchOffBody = `{ 
		"deviceid": "", 
		"data": {
			"switch": "off" 
		} 
	 }`
)
