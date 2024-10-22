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
	relayInfoBodyTmpl = `{ 
		"deviceid": "%s", 
		"data": { } 
	 }`

	relaySwitchOnBodyTmpl = `{ 
		"deviceid": "%s", 
		"data": {
			"switch": "on" 
		} 
	 }`

	relaySwitchOffBodyTmpl = `{ 
		"deviceid": "%s", 
		"data": {
			"switch": "off" 
		} 
	 }`
)
