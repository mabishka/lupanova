package model

import "errors"

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

var ErrorDeleted = errors.New("item deleted")

type Request struct {
	Full string `json:"url"`
}

type Response struct {
	Short string `json:"result"`
}

type ShortItem struct {
	Corr  string `json:"correlation_id"`
	Short string `json:"short_url"`
}

type FullItem struct {
	Corr string `json:"correlation_id"`
	Full string `json:"original_url"`
}

type StoreItem struct {
	Short string `json:"short_url"`
	Full  string `json:"original_url"`
}
