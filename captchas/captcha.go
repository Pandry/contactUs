package captchas

type ICaptcha interface {
	GetResponseFieldName() string
	VerifyCaptcha(response string, ip string) bool
}
