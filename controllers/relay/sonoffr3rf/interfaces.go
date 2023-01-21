package sonoffr3rf

type OSServicesProvider interface {
	GetIPFromMAC(mac string) (string, error)
}
