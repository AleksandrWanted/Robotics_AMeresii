package server

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Details string `json:"details"`
}

type DeviceListResponse struct {
	Devices []*Device `json:"devices"`
}

type SwitchPirSensorActivityStateRequest struct {
	Activity bool `json:"activity"`
}

type ArduinoSensorStateRequest struct {
	Room        string `json:"room"`
	SensorState string `json:"sensor_state"`
}
