package main

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/novel/auth/router"
)

/*
체크리스트
1. architecture -> directory 구조 -> https://github.com/eldimious/golang-api-showcase 이거 참고해서 변경합시다 (Best임)
the bestest clean architecture
공통 repository 같은건 common/repository 이게 맞음 , 그리고 DI를 main에서 주입해주는거 생각...
2. CreateUserToken, GenerateCsrfToken 메서드 분리 어디다가 할지
-----------------------------------------------------------------------------------------------------
3. database index, constraint, view(이거는 거의 알아서 뭐...), table lock(이거 중요 격리수준과 연관깊음) 찾아서 학습
4. golang 문법 틈틈히
5. redis, kafka도 해봐야하는데 redis는 나중에 다시 해보고 mas로 만들거니까 kafka 먼저 해보는게 좋을 것 같음
6. 나중에 grpc로 migration 해보는 것도 좋은데 이건 많이 나중에
*/

// '/' 경로는 Api Gateway 추가예정
func main() {
	e := echo.New()

	e.HTTPErrorHandler = func(err error, c echo.Context) {
		e.Logger.Error(err)
		c.JSON(http.StatusInternalServerError, "Internal Server Error")
	}

	customLogger := middleware.LoggerConfig{
		Output: os.Stdout,
		Format: `{"time":"${time_rfc3339_nano}","method":"${method}","uri":"${uri}","status":${status},"latency_human":"${latency_human}"}'` + "\n",
	}

	e.Use(
		middleware.LoggerWithConfig(customLogger),
		// middleware.Recover(), // recover는 production 환경에서 하는게 좋을 듯 함
	)

	router.SetRouter(e)

	e.Logger.Fatal(e.Start(":12121"))
}
