package my_heat

const (
	SeverityUnknownState  = 0
	SeverityNormalState   = 1
	SeverityWarningState  = 32
	SeverityCriticalState = 64

	TypeCircuitTemperature = "circuit_temperature"
	TypeRoomTemperature    = "room_temperature"
	TypeBoilerTemperature  = "boiler_temperature"
)
