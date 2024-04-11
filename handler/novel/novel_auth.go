package novel

import "github.com/novel/auth/middleware"

type AuthHandler struct {
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

func (a *AuthHandler) Signin(ctx *middleware.Ctx) {
	ctx.W.Write([]byte("HelloWorld"))
}

func (a *AuthHandler) Signout(ctx *middleware.Ctx) {

}

func (a *AuthHandler) Signup(ctx *middleware.Ctx) {

}
