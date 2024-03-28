package handler

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"
)

// csrf token
func generateState(w http.ResponseWriter) string {
	// 하루
	expiration := time.Now().Add(1 * 24 * time.Hour)

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)

	cookie := &http.Cookie{Name: "state", Value: state, Expires: expiration}
	http.SetCookie(w, cookie)
	return state
}
