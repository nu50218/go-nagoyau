# go-nagoyau

名古屋大学のサービスにログイン済みの*http.Clientを返してくれます

## 使い方

```go
// NewClient servicesで指定した名古屋大学のサービスにログイン済みの*http.Clientを返してくれます
func NewClient(username, password string, services ...Service)
```

を使います。

たとえば名大ポータルとNUCTにログイン済みの`*http.Client`を得たい場合は、以下のようにサービスを指定します。

```go
client, err := nagoyau.NewClient("username", "password", nagoyau.Portal, nagoyau.CT)
```

### 対応サービス一覧

| サービス   | `Service`              |
| ------ | ---------------- |
| 名大ポータル | `nagoyau.Portal` |
| NUCT   | `nagoyau.CT`     |

ログイン時のURLのクエリパラメーターの`service`を`service.go`に追加し、`NewClient`の`switch`の部分を書き足すだけなので、追加したいものがあればPRしてください！
