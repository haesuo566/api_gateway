package handler

import "github.com/novel/auth/middleware"

type AuthHandler struct {
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

func (a *AuthHandler) Signin(ctx *middleware.Ctx) {
	ctx.W.Write([]byte("HelloWorld"))
}
