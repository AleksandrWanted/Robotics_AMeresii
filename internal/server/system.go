package server

import (
	"ameresii_smart_home/internal/err_stack"
	"ameresii_smart_home/pkg/config_manager"
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
	"log"
	"net/http"
	"os"
)

func (s Server) HandlerSystemRun(ctx *fasthttp.RequestCtx) {
	confFileName, isExist := os.LookupEnv("CONF_FILE_NAME")
	if !isExist {
		confFileName = "main"
	}

	cfgPath := fmt.Sprintf("configs/%s.yaml", confFileName)

	if err := config_manager.CheckCfgAvailability(cfgPath); err != nil {
		log.Print(err_stack.WithStack(fmt.Errorf("configuration file main.yaml not found: %v", err)))

		ctx.Response.Header.Set("Content-Type", "application/json; charset=UTF-8")
		ctx.Response.SetStatusCode(http.StatusInternalServerError)

		res := &ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "Internal Server Error",
			Details: fmt.Sprintf("configuration file not found"),
		}

		err = json.NewEncoder(ctx).Encode(res)
		if err != nil {
			log.Print(err_stack.WithStack(fmt.Errorf("error json encode: %v", err)))
		}

		return
	}

	if err := s.smartHomeApp.Run(ctx, cfgPath); err != nil {
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
