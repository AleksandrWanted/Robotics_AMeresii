package my_heat

type MyHeatRequestBody struct {
	Action     string `json:"action,omitempty"`
	Login      string `json:"login,omitempty"`
	Key        string `json:"key,omitempty"`
	DeviceID   int    `json:"deviceId,omitempty"`
	ObjID      int    `json:"objId,omitempty"`
	Goal       int    `json:"goal,omitempty"`
	ChangeMode int    `json:"changeMode,omitempty"`
}

type MyHeatResponseBody struct {
	Data        *Data `json:"data,omitempty"`
	Err         int   `json:"err,omitempty"`
	RefreshPage bool  `json:"refreshPage,omitempty"`
	Schedule    int   `json:"schedule,omitempty"`
	HeatingMode int   `json:"heatingMode,omitempty"`
}

type Data struct {
	Devices      []*Device `json:"devices,omitempty"`
	Heaters      []*Heater `json:"heaters,omitempty"`
	Envs         []*Env    `json:"envs,omitempty"`
	Engs         []*Eng    `json:"engs,omitempty"`
	Alarms       []string  `json:"alarms,omitempty"`
	DataActual   bool      `json:"dataActual,omitempty"`
	Severity     int       `json:"severity"`
	SeverityDesc string    `json:"severityDesc,omitempty"`
	WeatherTemp  *string   `json:"weatherTemp"`
	City         string    `json:"city,omitempty"`
}

type Device struct {
	ID           int    `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	City         string `json:"city,omitempty"`
	Severity     int    `json:"severity"`
	SeverityDesc string `json:"severityDesc,omitempty"`
}

type Heater struct {
	ID            int      `json:"id,omitempty"`
	Name          string   `json:"name,omitempty"`
	Disabled      bool     `json:"disabled,omitempty"`
	FlowTemp      *float64 `json:"flowTemp,omitempty"`
	ReturnTemp    *float64 `json:"returnTemp"`
	Pressure      *float64 `json:"pressure"`
	TargetTemp    int      `json:"targetTemp"`
	BurnerWater   bool     `json:"burnerWater,omitempty"`
	BurnerHeating bool     `json:"burnerHeating,omitempty"`
	Modulation    *float64 `json:"modulation"`
	Severity      int      `json:"severity"`
	SeverityDesc  string   `json:"severityDesc,omitempty"`
}

type Env struct {
	ID           int      `json:"id,omitempty"`
	Type         string   `json:"type,omitempty"`
	Name         string   `json:"name,omitempty"`
	Value        *float64 `json:"value"`
	Target       *float64 `json:"target"`
	Demand       bool     `json:"demand,omitempty"`
	Severity     int      `json:"severity"`
	SeverityDesc string   `json:"severityDesc,omitempty"`
}

type Eng struct {
	ID           int    `json:"id,omitempty"`
	Type         string `json:"type,omitempty"`
	Name         string `json:"name,omitempty"`
	TurnedOn     bool   `json:"turnedOn,omitempty"`
	Severity     int    `json:"severity"`
	SeverityDesc string `json:"severityDesc,omitempty"`
}
