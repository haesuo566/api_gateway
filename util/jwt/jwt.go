package jwt

type JwtUtil struct {
}

func New() IJwtUtil {
	return &JwtUtil{}
}
