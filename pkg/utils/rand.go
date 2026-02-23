// Package utils вспомогательные функции
package utils

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
)

// CreateShort формирование рандомной строки длиной n.
func CreateShort(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b)[:n], nil
}

// ErrConflict ошибка
var ErrConflict = errors.New("already exist")

// Ошибка удаленного значения.
var ErrorDeleted = errors.New("item deleted")
