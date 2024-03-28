package handler

import "github.com/novel/auth/ancho"

type AuthHandler struct {
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

func (a *AuthHandler) Signin(ctx *ancho.Ctx) {
	ctx.W.Write([]byte("HelloWorld"))
}
