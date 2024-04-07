package jwt

const (
	ACCESS_TOKEN = iota
	REFRESH_TOKEN
)

type ResposneToken struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type IJwtUtil interface {
	GenerateAccessToken(email string) (string, error)
	GenerateRefreshToken() (string, error)
	ValidateToken(token string) error
}
