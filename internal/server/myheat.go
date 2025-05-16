package server

import (
	"ameresii_smart_home/internal/err_stack"
	"ameresii_smart_home/pkg/telegram"
	"fmt"
	"github.com/valyala/fasthttp"
	"log"
	"net/http"
)

func (s Server) HandlerMyheatDevicesList(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Content-Type", "application/json; charset=UTF-8")
	ctx.Response.SetStatusCode(http.StatusNoContent)
	err := telegram.SendMessage("some text")
	log.Print(err_stack.WithStack(err))
	log.Print(err_stack.WithStack(fmt.Errorf("some text")))
}
