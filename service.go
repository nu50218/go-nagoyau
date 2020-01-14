package nagoyau

import "net/url"

// Service 名大ポータルやNUCTなど
type Service string

const (
	// Portal 名大ポータル
	Portal Service = "https://portal.nagoya-u.ac.jp/login"
	// CT NUCT
	CT Service = "https://ct.nagoya-u.ac.jp/sakai-login-tool/container"
)

func makeServiceLoginURL(service Service) string {
	u, _ := url.Parse(loginAuthURL)

	q := u.Query()
	q.Set("service", string(service))
	u.RawQuery = q.Encode()

	return u.String()
}
