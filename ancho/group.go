package ancho

type JwtValidate func()

type UrlGroup struct {
	Url string
}

func Group(url string) *UrlGroup {
	return &UrlGroup{
		Url: url,
	}
}

func (u *UrlGroup) Use() {

}
