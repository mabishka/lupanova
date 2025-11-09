package fileloader

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"os"
	"sync"

	"github.com/mabishka/lupanova/internal/model"
)

type FileLoader struct {
	*sync.Mutex
	fileName string
	fileSize int64
}

func New(fileName string) *FileLoader {
	return &FileLoader{Mutex: &sync.Mutex{}, fileName: fileName}
}

// return map [short string] full string
func (p *FileLoader) Load(ctx context.Context) (map[string]string, error) {

	if err := p.create(); err != nil {
		return nil, err
	}

	p.Lock()
	defer p.Unlock()

	content, err := os.ReadFile(p.fileName)
	if err != nil {
		return nil, err
	}

	var data []model.StoreItem

	if err := json.Unmarshal(content, &data); err != nil {
		return nil, err
	}

	response := make(map[string]string)
	for _, v := range data {
		response[v.Short] = v.Full
	}

	return response, nil
}

func (p *FileLoader) Store(ctx context.Context, full, short string) error {

	if err := p.create(); err != nil {
		return err
	}

	p.Lock()
	defer p.Unlock()

	file, err := os.OpenFile(p.fileName, os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	filesize := int64(p.fileSize) - 1
	if err := file.Truncate(filesize); err != nil {
		return err
	}
	p.fileSize = filesize

	file.Seek(0, io.SeekEnd)

	defer func() {
		n, err := file.Write([]byte("]"))
		if err != nil {
			return
		}
		p.fileSize += int64(n)
	}()

	fileData := model.StoreItem{
		Short: short,
		Full:  full,
	}

	sendData, err := json.Marshal(&fileData)
	if err != nil {
		return err
	}

	if p.fileSize > 2 {
		sendData = append([]byte(",\n"), sendData...)
	}

	sendData = append(sendData, ']')
	n, err := file.Write(sendData)
	if err != nil {
		return err
	}
	p.fileSize += int64(n)

	return nil
}

func (p *FileLoader) StoreList(ctx context.Context, list []model.StoreItem) error {

	if err := p.create(); err != nil {
		return err
	}

	p.Lock()
	defer p.Unlock()

	file, err := os.OpenFile(p.fileName, os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	filesize := int64(p.fileSize) - 1
	if err := file.Truncate(filesize); err != nil {
		return err
	}
	p.fileSize = filesize

	file.Seek(0, io.SeekEnd)

	defer func() {
		n, err := file.Write([]byte("]"))
		if err != nil {
			return
		}
		p.fileSize += int64(n)
	}()

	buffer := bufio.NewWriter(file)
	var size int

	for _, v := range list {

		sendData, err := json.Marshal(&v)
		if err != nil {
			return err
		}

		if p.fileSize > 2 {
			n, err := buffer.Write([]byte(",\n"))
			if err != nil {
				return err
			}
			size += n
		}

		n, err := buffer.Write(sendData)
		if err != nil {
			return err
		}
		size += n
	}

	if err := buffer.Flush(); err != nil {
		return err
	}
	p.fileSize += int64(size)

	return nil
}

func (p *FileLoader) create() error {

	exist, err := p.exist()
	if err != nil {
		return err
	}
	if exist {
		return nil
	}

	p.Lock()
	defer p.Unlock()

	if err := os.WriteFile(p.fileName, []byte("[]"), 0644); err != nil {
		return err
	}
	p.fileSize = 2

	return nil
}

func (p *FileLoader) exist() (bool, error) {

	if p.fileSize != 0 {
		return true, nil
	}

	p.Lock()
	defer p.Unlock()

	stat, err := os.Stat(p.fileName)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	p.fileSize = stat.Size()

	return p.fileSize != 0, nil
}
