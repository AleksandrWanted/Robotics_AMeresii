package jobs

import (
	"fmt"
	"github.com/AleksandrWanted/AMeresii_SMART_HOME/internal/err_stack"
	"github.com/AleksandrWanted/AMeresii_SMART_HOME/internal/server"
	"github.com/AleksandrWanted/AMeresii_SMART_HOME/pkg/config_manager"
	"github.com/AleksandrWanted/AMeresii_SMART_HOME/pkg/my_heat"
	"github.com/AleksandrWanted/AMeresii_SMART_HOME/pkg/telegram"
	"log"
	"os"
	"strconv"
	"strings"
)

func ControlTemperature() {
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

		for _, obj := range config_manager.SmartHomeConfig.Config().SystemGoalParams.Objects {
			if device.ID == obj.ID {
				if err = heatingAdjustment(device.ID, deviceInfo, *obj); err != nil {
					if err = telegram.SendMessage(fmt.Sprintf("Неудачная попытка регулировки температуры объекта %s",
						device.Name)); err != nil {
						log.Print(err_stack.WithStack(fmt.Errorf(
							"failed to send telegram message: %v", err)))
					}
				}
			}
		}
	}
}

func heatingAdjustment(deviceID int, currentState *my_heat.Data, targetState config_manager.Object) error {
	login := os.Getenv("MYHEAT_LOGIN")
	key := os.Getenv("MYHEAT_API_KEY")
	var currentRoomTemp float64
	var currentBoilTemp float64
	var currentTargetCircuitTemp float64
	var circuitTempEnvID int

	for _, env := range currentState.Envs {
		switch env.Type {
		case my_heat.TypeRoomTemperature:
			currentRoomTemp = *env.Value
		case my_heat.TypeBoilerTemperature:
			currentBoilTemp = *env.Value
		case my_heat.TypeCircuitTemperature:
			currentTargetCircuitTemp = *env.Target
			circuitTempEnvID = env.ID
		}
	}

	switch strings.ToLower(targetState.HeatingMode) {
	case "fast":
		if currentRoomTemp < float64(targetState.RoomTemperature) || currentBoilTemp < float64(targetState.BoilerTemperature) {
			targetCircuitTempStr := os.Getenv("FAST_HEATING_MODE_CIRCUIT_TEMP")
			targetCircuitTemp, _ := strconv.Atoi(targetCircuitTempStr)

			if int(currentTargetCircuitTemp) != targetCircuitTemp {
				if _, err := my_heat.MyHeatSetEnvGoal(&my_heat.MyHeatRequestBody{
					Login:    login,
					Key:      key,
					DeviceID: deviceID,
					ObjID:    circuitTempEnvID,
					Goal:     targetCircuitTemp,
				}); err != nil {
					return err
				}
			}

		} else {
			if currentTargetCircuitTemp != currentRoomTemp {
				if _, err := my_heat.MyHeatSetEnvGoal(&my_heat.MyHeatRequestBody{
					Login:    login,
					Key:      key,
					DeviceID: deviceID,
					ObjID:    circuitTempEnvID,
					Goal:     int(currentRoomTemp),
				}); err != nil {
					return err
				}
			}
		}

	case "slow":
		if currentRoomTemp < float64(targetState.RoomTemperature) || currentBoilTemp < float64(targetState.BoilerTemperature) {
			targetCircuitTempStr := os.Getenv("SLOW_HEATING_MODE_CIRCUIT_TEMP")
			targetCircuitTemp, _ := strconv.Atoi(targetCircuitTempStr)

			if int(currentTargetCircuitTemp) != targetCircuitTemp {
				if _, err := my_heat.MyHeatSetEnvGoal(&my_heat.MyHeatRequestBody{
					Login:    login,
					Key:      key,
					DeviceID: deviceID,
					ObjID:    circuitTempEnvID,
					Goal:     targetCircuitTemp,
				}); err != nil {
					return err
				}
			}

		} else {
			if currentTargetCircuitTemp != currentRoomTemp {
				if _, err := my_heat.MyHeatSetEnvGoal(&my_heat.MyHeatRequestBody{
					Login:    login,
					Key:      key,
					DeviceID: deviceID,
					ObjID:    circuitTempEnvID,
					Goal:     int(currentRoomTemp),
				}); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
