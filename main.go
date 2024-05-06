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
1. database index, constraint, view(이거는 거의 알아서 뭐...), table lock(이거 중요 격리수준과 연관깊음) 찾아서 학습
2. golang 문법 틈틈히
3. redis, kafka도 해봐야하는데 redis는 나중에 다시 해보고 mas로 만들거니까 kafka 먼저 해보는게 좋을 것 같음
4. 나중에 grpc로 migration 해보는 것도 좋은데 이건 많이 나중에
5. logout(나중에 구현) 은 redis 로 처리합시다. 나중에 캐시데이터도 redis 에 넣는걸로
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
