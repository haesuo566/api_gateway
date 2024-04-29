package main

import (
	"github.com/labstack/echo/v4"
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

	// e.Use()
	// e.Use()
	// e.Use()

	router.SetRouter(e)

	e.Logger.Fatal(e.Start(":1323"))
	// e.Group("/asd", )
	// app := middleware.New()

	// // Handlers
	// googleAuthHandler := google.NewHandler()
	// naverAuthHandler := naver.NewHandler()

	// // google auth
	// app.Get("/auth/google", googleAuthHandler.Signin)
	// app.Get("/auth/google/callback", googleAuthHandler.Callback)

	// // naver auth
	// app.Get("/auth/naver", naverAuthHandler.Signin)
	// app.Get("/auth/naver/callback", naverAuthHandler.Callback)

	// if err := app.Run(":12121"); err != nil {
	// 	log.Println(err)
	// }
}
