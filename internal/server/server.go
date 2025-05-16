package server

import "ameresii_smart_home/pkg/smart_home"

type Server struct {
	smartHomeApp smart_home.App
}

func NewServer(app smart_home.App) Server {
	return Server{
		smartHomeApp: app,
	}
}
