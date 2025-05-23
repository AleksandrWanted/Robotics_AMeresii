package server

import (
	"ameresii_smart_home/internal/err_stack"
	"ameresii_smart_home/pkg/telegram"
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
	"log"
	"net/http"
)

func (s Server) HandlerSystemRun(ctx *fasthttp.RequestCtx) {
	if err := s.smartHomeApp.Run(ctx); err != nil {
		log.Print(err_stack.WithStack(fmt.Errorf("failed to run method: %v", err)))
		ctx.Response.Header.Set("Content-Type", "application/json; charset=UTF-8")
		ctx.Response.SetStatusCode(http.StatusInternalServerError)

		res := &ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "Internal Server Error",
			Details: fmt.Sprintf("failed to run method: %v", err),
		}

		err = json.NewEncoder(ctx).Encode(res)
		if err != nil {
			log.Print(err_stack.WithStack(fmt.Errorf("error json encode: %v", err)))
		}

		return
	}

	ctx.Response.Header.Set("Content-Type", "application/json; charset=UTF-8")
	ctx.Response.SetStatusCode(http.StatusNoContent)
}

func (s Server) HandlerSystemPirSwitchActivityState(ctx *fasthttp.RequestCtx) {
	ctx.SetConnectionClose()

	var request SwitchPirSensorActivityStateRequest
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

	ctx.Response.Header.Set("Content-Type", "application/json; charset=UTF-8")
	ctx.Response.SetStatusCode(http.StatusNoContent)

	if request.Activity != PirSensorActivity {
		PirSensorActivity = request.Activity

		var messageText string
		switch request.Activity {
		case true:
			messageText = "Датчик движения включен!"
		default:
			messageText = "Датчик движения выключен!"
		}

		if err := telegram.SendMessage(messageText); err != nil {
			log.Print(err_stack.WithStack(fmt.Errorf("error while send telegram message %v", err)))
		}
	}
}
