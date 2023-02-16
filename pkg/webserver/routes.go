package webserver

import "net/http"

func SetupRoutes(mux *http.ServeMux, h *handlers) {
	mux.HandleFunc("/boiler/get_state", h.BoilerHandler.GetState)
	mux.HandleFunc("/boiler/turn_on", h.BoilerHandler.TurnOn)
	mux.HandleFunc("/boiler/turn_off", h.BoilerHandler.TurnOff)

	mux.HandleFunc("/room_light/get_state", h.RoomLightHandler.GetState)
	mux.HandleFunc("/room_light/turn_on", h.RoomLightHandler.TurnOn)
	mux.HandleFunc("/room_light/turn_off", h.RoomLightHandler.TurnOff)

	mux.HandleFunc("/cam_light/get_state", h.CamLightHandler.GetState)
	mux.HandleFunc("/cam_light/turn_on", h.CamLightHandler.TurnOn)
	mux.HandleFunc("/cam_light/turn_off", h.CamLightHandler.TurnOff)

	mux.HandleFunc("/cell/balance", h.LTEModuleHandler.GetAccountBalance)
	mux.HandleFunc("/cell/inet", h.LTEModuleHandler.GetInetBalance)
	mux.HandleFunc("/cell/tariff", h.LTEModuleHandler.GetTariff)
}
