package businesslogic

type RoomLightController interface {
	Start() error
	Stop() error
}
