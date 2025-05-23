package server

import (
	"ameresii_smart_home/internal/err_stack"
	"ameresii_smart_home/pkg/my_heat"
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
	"log"
	"net/http"
	"os"
)

func (s Server) HandlerMyheatDevicesList(ctx *fasthttp.RequestCtx) {
	login := os.Getenv("MYHEAT_LOGIN")
	key := os.Getenv("MYHEAT_API_KEY")

	ctx.Response.Header.Set("Content-Type", "application/json; charset=UTF-8")

	devicesList, err := my_heat.MyHeatGetDevices(login, key)
	if err != nil {
		ctx.Response.SetStatusCode(http.StatusInternalServerError)

		res := &ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "Internal Server Error",
			Details: fmt.Sprintf("myheat endpoint not available"),
		}

		err = json.NewEncoder(ctx).Encode(res)
		if err != nil {
			log.Print(err_stack.WithStack(fmt.Errorf("error json encode: %v", err)))
		}

		return
	}

	ctx.Response.SetStatusCode(http.StatusOK)

	resBody := DeviceListResponse{Devices: make([]*Device, 0)}
	for _, device := range devicesList {
		newDevice := &Device{
			Name:      device.Name,
			ID:        device.ID,
			State:     device.Severity,
			StateDesc: device.SeverityDesc,
		}

		resBody.Devices = append(resBody.Devices, newDevice)
		DevicesMap[device.Name] = *newDevice
	}

	if err = json.NewEncoder(ctx).Encode(&resBody); err != nil {
		log.Print(err_stack.WithStack(fmt.Errorf("error json encode: %v", err)))
	}

	return

}
