package model

type Request struct {
	Full string `json:"url"`
}

type Response struct {
	Short string `json:"result"`
}
