package auth

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"
)

func GenerateState(w http.ResponseWriter) string {
	expires := time.Now().Add(time.Hour * 24)

	data := make([]byte, 16)
	rand.Read(data)
	state := base64.URLEncoding.EncodeToString(data)
	cookie := &http.Cookie{
		Name:    "state",
		Value:   state,
		Expires: expires,
	}
	http.SetCookie(w, cookie)
	return state
}
