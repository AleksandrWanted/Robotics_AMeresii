package server

import (
	"ameresii_smart_home/internal/err_stack"
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
