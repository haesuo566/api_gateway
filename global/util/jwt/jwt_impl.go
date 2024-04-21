package jwt

const (
	ACCESS_TOKEN = iota
	REFRESH_TOKEN
)

type IJwtUtil interface {
	GenerateAccessToken(email string) (string, error)
	GenerateRefreshToken() (string, error)
	ValidateToken(token string) error
}
