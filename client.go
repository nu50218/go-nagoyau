package nagoyau

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

const loginAuthURL = "https://auth.nagoya-u.ac.jp/cas/login"

func makeHTTPError(statusCode int, method, url string) error {
	return fmt.Errorf("ステータスコードが異常です: %d (Method: %s, URL: %s)", statusCode, method, url)
}

// NewClient servicesで指定した名古屋大学のサービスにログイン済みの*http.Clientを返してくれます
func NewClient(username, password string, services ...Service) (*http.Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{
		Jar: jar,
	}

	// 一回GETする
	res, err := client.Get(loginAuthURL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode >= 300 {
		return nil, makeHTTPError(res.StatusCode, http.MethodGet, loginAuthURL)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	// フォームの値をセットする
	values := url.Values{}
	values.Set("username", username)
	values.Set("password", password)
	for _, name := range []string{"lt", "_eventId", "submit"} {
		value, exists := doc.Find(fmt.Sprintf("input[name=%s]", name)).Attr("value")
		if !exists {
			return nil, fmt.Errorf("ログイン画面のinputに存在しません: %v", name)
		}
		values.Set(name, value)
	}

	// ログイン
	res, err = client.PostForm(loginAuthURL, values)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode >= 300 {
		return nil, makeHTTPError(res.StatusCode, http.MethodPost, loginAuthURL)
	}

	// 認証成功したかをチェックする
	// CASTGCがクッキーにセットされていたら成功
	cookies := res.Cookies()
	foundCASTGC := false
	for _, cookie := range cookies {
		if cookie.Name == "CASTGC" {
			foundCASTGC = true
		}
	}
	if !foundCASTGC {
		return nil, errors.New("認証に失敗しました")
	}

	// servicesのうちログインに特別な操作が必要なものの操作を行う
	for _, service := range services {
		switch service {
		case CT:
			url := makeServiceLoginURL(CT)
			res, err := client.Get(url)
			if err != nil {
				return nil, err
			}
			defer res.Body.Close()
			if res.StatusCode >= 300 {
				return nil, makeHTTPError(res.StatusCode, http.MethodGet, url)
			}
		default:
			// 何もしない
		}
	}

	return client, nil
}
