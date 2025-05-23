package server

import (
	"ameresii_smart_home/internal/err_stack"
	"ameresii_smart_home/pkg/telegram"
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
	"log"
	"net/http"
	"time"
)

var PirSensorActivity bool

const (
	bathroom1LeakSensorLabel   = "bathroom_first_floor"
	bathroom2LeakSensorLabel   = "bathroom_second_floor"
	kitchenLeakSensorLabel     = "kitchen"
	LeakDetectedStatus         = "Leak Detected"
	mainEntrancePirSensorLabel = "main_entrance"
	backyardPirSensorLabel     = "backyard"
	balconyPirSensorLabel      = "balcony"
	PirMotionStatus            = "Motion Detected"
)

func (s Server) HandlerArduinoLeakSensorState(ctx *fasthttp.RequestCtx) {
	ctx.SetConnectionClose()

	var request ArduinoSensorStateRequest
	if err := json.Unmarshal(ctx.Request.Body(), &request); err != nil {
		ctx.Response.Header.Set("Content-Type", "application/json; charset=UTF-8")
		ctx.Response.SetStatusCode(http.StatusInternalServerError)

		res := &ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "Internal Server Error",
			Details: err.Error(),
		}

		err = json.NewEncoder(ctx).Encode(res)
		if err != nil {
			log.Print(err_stack.WithStack(fmt.Errorf("error json encode: %v", err)))
		}

		return
	}

	if request.Room == "" || request.SensorState == "" {
		ctx.Response.Header.Set("Content-Type", "application/json; charset=UTF-8")
		ctx.Response.SetStatusCode(http.StatusBadRequest)

		res := &ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "Bad Request",
			Details: fmt.Sprintf("Fields room and sensor_state are required"),
		}

		err := json.NewEncoder(ctx).Encode(res)
		if err != nil {
			log.Print(err_stack.WithStack(fmt.Errorf("error json encode: %v", err)))
		}

		return
	}

	ctx.Response.Header.Set("Content-Type", "application/json; charset=UTF-8")
	ctx.Response.SetStatusCode(http.StatusNoContent)

	if request.SensorState == LeakDetectedStatus {
		var err error
		var roomName string
		switch request.Room {
		case bathroom1LeakSensorLabel:
			roomName = "санузел первого этажа"
		case bathroom2LeakSensorLabel:
			roomName = "санузел второго этажа"
		case kitchenLeakSensorLabel:
			roomName = "кухня"
		default:
			return
		}

		for i := 0; i < 3; i++ {
			if err = telegram.SendMessage(fmt.Sprintf("Тревога!!! Обнаружена протечка в помещении %s",
				roomName)); err != nil {
				log.Print(err_stack.WithStack(fmt.Errorf("error while send telegram message %v", err)))
				time.Sleep(5 * time.Second)
				continue
			}
			break
		}
	}
}

func (s Server) HandlerArduinoPirSensorState(ctx *fasthttp.RequestCtx) {
	ctx.SetConnectionClose()

	var request ArduinoSensorStateRequest
	if err := json.Unmarshal(ctx.Request.Body(), &request); err != nil {
		ctx.Response.Header.Set("Content-Type", "application/json; charset=UTF-8")
		ctx.Response.SetStatusCode(http.StatusInternalServerError)

		res := &ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "Internal Server Error",
			Details: err.Error(),
		}

		err = json.NewEncoder(ctx).Encode(res)
		if err != nil {
			log.Print(err_stack.WithStack(fmt.Errorf("error json encode: %v", err)))
		}

		return
	}

	if request.Room == "" || request.SensorState == "" {
		ctx.Response.Header.Set("Content-Type", "application/json; charset=UTF-8")
		ctx.Response.SetStatusCode(http.StatusBadRequest)

		res := &ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "Bad Request",
			Details: fmt.Sprintf("Fields room and sensor_state are required"),
		}

		err := json.NewEncoder(ctx).Encode(res)
		if err != nil {
			log.Print(err_stack.WithStack(fmt.Errorf("error json encode: %v", err)))
		}

		return
	}

	ctx.Response.Header.Set("Content-Type", "application/json; charset=UTF-8")
	ctx.Response.SetStatusCode(http.StatusNoContent)

	if !PirSensorActivity {
		return
	}

	if request.SensorState == PirMotionStatus {
		var err error
		var roomName string
		switch request.Room {
		case mainEntrancePirSensorLabel:
			roomName = "главный вход"
		case backyardPirSensorLabel:
			roomName = "вход задний двор"
		case balconyPirSensorLabel:
			roomName = "вход балкон"
		default:
			return
		}

		for i := 0; i < 3; i++ {
			if err = telegram.SendMessage(fmt.Sprintf("Тревога!!! Обнаружено движение в зоне %s",
				roomName)); err != nil {
				log.Print(err_stack.WithStack(fmt.Errorf("error while send telegram message %v", err)))
				time.Sleep(3 * time.Second)
				continue
			}
			break
		}
	}
}
