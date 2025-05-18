package server

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Details string `json:"details"`
}

type DeviceListResponse struct {
	Devices []*Device `json:"devices"`
}
