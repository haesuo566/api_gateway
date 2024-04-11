package main

import (
	"log"

	"github.com/novel/auth/handler/google"
	"github.com/novel/auth/handler/naver"
	"github.com/novel/auth/handler/novel"
	"github.com/novel/auth/middleware"
)

/*
체크리스트
1. architecture 곰곰히 생각 -> 폴더 구조
2. CreateUserToken, GenerateCsrfToken 메서드 분리 어디다가 할지
3. database index, constraint 찾아서 학습
4. golang 문법 틈틈히
*/

// '/' 경로는 Api Gateway 추가예정
func main() {
	app := middleware.New()

	// Handlers
	authHandler := novel.NewAuthHandler()
	googleAuthHandler := google.New()
	naverAuthHandler := naver.New()

	app.Get("/signin", authHandler.Signin)

	// google auth
	app.Get("/auth/google", googleAuthHandler.Signin)
	app.Get("/auth/google/callback", googleAuthHandler.Callback)

	// naver auth
	app.Get("/auth/naver", naverAuthHandler.Signin)
	app.Get("/auth/naver/callback", naverAuthHandler.Callback)

	if err := app.Run(":12121"); err != nil {
		log.Println(err)
	}
}
