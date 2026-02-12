package audit

import (
	"context"
	"os"
)

const observerFileName = "file"

// FileObserver Хранилище аудита в файле.
type FileObserver struct {
	name string
}

// NewFileObserver создание хранилища аудита в файле.
func NewFileObserver(name string) *FileObserver {
	// TODO: проверить существование файла, попоробовать открыть его для записи
	return &FileObserver{name: name}
}

// GetName имя аудита.
func (p *FileObserver) GetName() string {
	return observerFileName
}

// Send отправка аудита в файл.
func (p *FileObserver) Send(ctx context.Context, data []byte) error {

	f, err := os.OpenFile(p.name, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	defer f.Close()
	data = append(data, '\n')

	if _, err = f.Write(data); err != nil {
		return err
	}
	return nil
}
