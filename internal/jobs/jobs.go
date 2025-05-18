package jobs

import (
	"ameresii_smart_home/internal/err_stack"
	"ameresii_smart_home/internal/server"
	"ameresii_smart_home/pkg/my_heat"
	"ameresii_smart_home/pkg/telegram"
	"fmt"
	"log"
	"os"
)

func CheckDevicesState() {
	login := os.Getenv("MYHEAT_LOGIN")
	key := os.Getenv("MYHEAT_API_KEY")

	for _, device := range server.DevicesMap {
		deviceInfo, err := my_heat.MyHeatGetDeviceInfo(login, key, device.ID)
		if err != nil {
			log.Print(err_stack.WithStack(fmt.Errorf(
				"receiving information for device %s failed: %v", device.Name, err)))
			if err = telegram.SendMessage(fmt.Sprintf("Не удалось получить информацию о состоянии объекта %s",
				device.Name)); err != nil {
				log.Print(err_stack.WithStack(fmt.Errorf(
					"failed to send telegram message: %v", err)))
			}
		}

		var message string
		switch deviceInfo.Severity {
		case my_heat.SeverityNormalState:
			continue
		case my_heat.SeverityUnknownState:
			message = fmt.Sprintf("Состояние объекта %s неизвестно:\nДетали:\n", device.Name)
		case my_heat.SeverityWarningState:
			message = fmt.Sprintf("Состояние объекта %s требует внимания:\nДетали:\n", device.Name)
		case my_heat.SeverityCriticalState:
			message = fmt.Sprintf("Состояние объекта %s критическое:\nДетали:\n", device.Name)
		}

		if len(deviceInfo.Alarms) > 0 {
			message += prepareAlarmsMessage(deviceInfo)
		}

		message += prepareObjStateMessage(deviceInfo)

		if err = telegram.SendMessage(message); err != nil {
			log.Print(err_stack.WithStack(fmt.Errorf(
				"failed to send telegram message: %v", err)))
		}
	}

}

func prepareObjStateMessage(data *my_heat.Data) string {
	message := "Текущее состояние системы:\n"

	message += "Котлы:\n"
	for _, heater := range data.Heaters {
		heaterStateMessage := fmt.Sprintf("- название: %s\n", heater.Name)

		if heater.Disabled {
			heaterStateMessage += "статус: включен\n"
		} else {
			heaterStateMessage += "статус: отключен\n"
		}

		heaterStateMessage += fmt.Sprintf("состояние: %s\n", heater.SeverityDesc)

		if heater.BurnerHeating {
			heaterStateMessage += "нагрев помещения: да\n"
		} else {
			heaterStateMessage += "нагрев помещения: нет\n"
		}

		if heater.BurnerWater {
			heaterStateMessage += "нагрев воды: да\n"
		} else {
			heaterStateMessage += "нагрев воды: нет\n"
		}

		heaterStateMessage += fmt.Sprintf("Температура подачи контура: %f\n", *heater.FlowTemp)

		if heater.ReturnTemp != nil {
			heaterStateMessage += fmt.Sprintf("Температура обратки контура: %f\n", *heater.ReturnTemp)
		}

		if heater.Pressure != nil {
			heaterStateMessage += fmt.Sprintf("Давление контура: %f\n", *heater.Pressure)
		}

		heaterStateMessage += fmt.Sprintf("Целевая температура: %d\n", heater.TargetTemp)

		message += heaterStateMessage
	}

	if data.Envs != nil {
		message += "Датчики температуры:\n"
		for _, tSensor := range data.Envs {
			tSensorStateMessage := fmt.Sprintf("- название: %s\n", tSensor.Name)

			tSensorStateMessage += fmt.Sprintf("тип: %s\n", tSensor.Type)

			tSensorStateMessage += fmt.Sprintf("состояние: %s\n", tSensor.SeverityDesc)

			if tSensor.Value != nil {
				tSensorStateMessage += fmt.Sprintf("температура помещения: %f\n", *tSensor.Value)
			}

			if tSensor.Target != nil {
				tSensorStateMessage += fmt.Sprintf("целевая температура: %f\n", *tSensor.Target)
			}

			message += tSensorStateMessage
		}
	}

	if data.Engs != nil {
		message += "Инженерное оборудование:\n"
		for _, eng := range data.Engs {
			engStateMessage := fmt.Sprintf("- название: %s\n", eng.Name)

			engStateMessage += fmt.Sprintf("тип: %s\n", eng.Type)

			engStateMessage += fmt.Sprintf("состояние: %s\n", eng.SeverityDesc)

			if eng.TurnedOn {
				engStateMessage += "статус: включен\n"
			} else {
				engStateMessage += "статус: отключен\n"
			}

			message += engStateMessage
		}
	}

	return message
}

func prepareAlarmsMessage(data *my_heat.Data) string {
	message := "Тревожные сообщения:\n"
	for _, alarm := range data.Alarms {
		message += fmt.Sprintf("-%s\n", alarm)
	}
	return message
}
