package server

import (
	"ameresii_smart_home/pkg/my_heat"
	"ameresii_smart_home/pkg/smart_home"
	"ameresii_smart_home/pkg/telegram"
	"fmt"
	"os"
)

var DevicesMap map[string]Device

type Server struct {
	smartHomeApp smart_home.App
}

type Device struct {
	Name      string `json:"name"`
	ID        int    `json:"ID"`
	State     int    `json:"state"`
	StateDesc string `json:"stateDesc"`
}

func NewServer(app smart_home.App) Server {
	return Server{
		smartHomeApp: app,
	}
}

func ReceiveDevicesList() error {
	login := os.Getenv("MYHEAT_LOGIN")
	key := os.Getenv("MYHEAT_API_KEY")

	deviceList, err := my_heat.MyHeatGetDevices(login, key)
	if err != nil {
		return err
	}

	for _, device := range deviceList {
		DevicesMap[device.Name] = Device{
			Name:      device.Name,
			ID:        device.ID,
			State:     device.Severity,
			StateDesc: device.SeverityDesc,
		}
	}

	return nil
}

func SendSystemStartingNotification() error {
	message := "Система AMeresii_SMART_HOME успешно запущена!\nДоступные объекты:\n"
	for _, device := range DevicesMap {
		deviceText := fmt.Sprintf("- %s:\n   Статус: %s\n", device.Name, device.StateDesc)
		message += deviceText
	}

	if err := telegram.SendMessage(message); err != nil {
		return err
	}

	return nil
}
