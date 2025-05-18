package main

import (
	"ameresii_smart_home/internal/err_stack"
	"ameresii_smart_home/internal/jobs"
	"ameresii_smart_home/internal/server"
	"ameresii_smart_home/pkg/dotenv"
	"ameresii_smart_home/pkg/jobs_manager"
	"ameresii_smart_home/pkg/smart_home"
	"context"
	"fmt"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"log"
	"sync"
	"time"
)

var mu sync.Mutex

func main() {
	server.DevicesMap = make(map[string]server.Device)
	ctx := context.Background()

	dotenv.Load()

	if err := server.ReceiveDevicesList(); err != nil {
		log.Print(err_stack.WithStack(fmt.Errorf("receive devices list failed: %v", err)))
	}

	if err := server.SendSystemStartingNotification(); err != nil {
		log.Print(err_stack.WithStack(fmt.Errorf("send start notification failed: %v", err)))
	}

	srv := server.NewServer(smart_home.NewSmartHomeApp(ctx))
	smartHomeServer := router.New()

	smartHomeServer.POST("/api/system/run", srv.HandlerSystemRun)

	smartHomeServer.POST("/api/myheat/devices/list", srv.HandlerMyheatDevicesList)

	httpServer := &fasthttp.Server{
		Name:                               "AMeresii_SMART_home",   // название (похоже что используется только в заголовке Server)
		Handler:                            smartHomeServer.Handler, // хандлер запросов
		ErrorHandler:                       nil,                     // хандлер ошибок во время разбора запроса
		ConnState:                          nil,                     // хандлер на получение статуса подключений
		HeaderReceived:                     nil,                     // не используем
		MaxRequestBodySize:                 1024 * 1024,             // максимальный размер тела запроса
		ContinueHandler:                    nil,                     // Continue requests не используем
		Concurrency:                        256 * 1024,              // максимальное кол-во подключений
		ReadBufferSize:                     4096,                    // размер буфера на чтение
		WriteBufferSize:                    4096,                    // размер буфера на запись
		MaxConnsPerIP:                      0,                       // кол-во подключений на ip не лимитируем
		MaxRequestsPerConn:                 0,                       // кол-во запросов на подключение не лимитируем
		Logger:                             nil,                     // отключаем логи
		LogAllErrors:                       false,                   // отключаем логи
		SecureErrorLogMessage:              false,                   // отключаем логи
		DisableKeepalive:                   false,                   // Keepalive включаем http сессии
		TCPKeepalive:                       true,                    // Keepalive включаем tcp сессии
		TCPKeepalivePeriod:                 time.Second * 3600,      // Keepalive время жизни tcp сессии
		IdleTimeout:                        time.Second * 60,        // Keepalive таймаут между запросами
		ReadTimeout:                        0,                       // Keepalive не лимитируем чтение
		WriteTimeout:                       0,                       // не лимитируем запись
		ReduceMemoryUsage:                  false,                   // максимально используем память
		GetOnly:                            false,                   // разрешаем все методы в запросах
		DisablePreParseMultipartForm:       true,                    // парсинг multipart отключаем
		DisableHeaderNamesNormalizing:      false,                   // включаем нормализацию названий заголовков
		SleepWhenConcurrencyLimitsExceeded: 0,                       // отключаем паузы при достижении лимита на колл-во параллельных запросов
		NoDefaultServerHeader:              true,                    // отключаем дефолтные заголовки сервера
		NoDefaultDate:                      true,                    // отключаем дефолтные заголовки даты
		NoDefaultContentType:               false,                   // включаем дефолтный тип контента
		StreamRequestBody:                  false,                   // стримы не используем
		KeepHijackedConns:                  false,                   // отключаем, так как ws не используем
		CloseOnShutdown:                    true,                    // сервер нормально закроет keepalive подключения при shutdown событии
	}

	if err := httpServer.ListenAndServe(":19930"); err != nil {
		log.Print(err_stack.WithStack(fmt.Errorf("failed listen server: %v", err)))
	}

}

func init() {
	jobs_manager.Register(jobs.ExampleJob)
	jobs_manager.Register(jobs.CheckDevicesState)
}
