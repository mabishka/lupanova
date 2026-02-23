package model

// Константы отправки запросов
const (
	HeaderContentType     = "Content-Type"
	HeaderContentEncoding = "Content-Encoding"
	HeaderAcceptEncoding  = "Accept-Encoding"
	HeaderLocation        = "Location"
	HeaderAuth            = "Authorisation"
	ContentTypeText       = "text/plain"
	ContentTypeJSON       = "application/json"
	ContentTypeHTML       = "text/html"

	CookieAuth = "Auth"
	CookieUser = "User"

	ContextValueUser = "User"
)

// Запрос POST /api/shorten.
// generate:reset
type Request struct {
	Full string `json:"url"`
}

// Ответ POST /api/shorten.
// generate:reset
type Response struct {
	Short string `json:"result"`
}

// отыет POST /api/shorten/batch.
// generate:reset
type ShortItem struct {
	Corr  string `json:"correlation_id"`
	Short string `json:"short_url"`
}

// запрос POST /api/shorten/batch.
// generate:reset
type FullItem struct {
	Corr string `json:"correlation_id"`
	Full string `json:"original_url"`
}

// ответ GET /api/user/urls.
// generate:reset
type StoreItem struct {
	Short string `json:"short_url"`
	Full  string `json:"original_url"`
}

// данные аудита.
//
//	{
//	 "ts": 12345678,        // unix timestamp события
//	 "action": "shorten",   // действие: shorten (создание) или follow (прохождение по ссылке)
//	 "user_id": "12315134", // идентификатор пользователя, если есть
//	 "url": "https://mylongdomain.com/my/long/path/to/shorten/" // оригинальный (не сокращенный) URL
//	}
//
// generate:reset
type AuditData struct {
	Created int64  `json:"ts"`
	Action  string `json:"action"`
	User    string `json:"user_id"`
	Address string `json:"url"`
}

// активность аудита
const (
	ActionShorten = "shorten"
	ActionFollow  = "follow"
)
