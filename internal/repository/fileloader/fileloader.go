package fileloader

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/mabishka/lupanova/internal/config"
	"github.com/mabishka/lupanova/internal/model"
	"github.com/mabishka/lupanova/pkg/rand"
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

func (p *FileLoader) GetShortList(ctx context.Context, fullList []model.FullItem) (map[string]string, error) {
	if err := p.create(); err != nil {
		return nil, err
	}

	p.Lock()
	defer p.Unlock()

	file, err := os.OpenFile(p.fileName, os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	if err := p.preSave(file); err != nil {
		return nil, err
	}
	defer p.postSave(file)

	buffer := bufio.NewWriter(file)
	var size int

	storeList := make(map[string]string)
	for _, v := range fullList {

		short, n, err := p.writeItem(buffer, v.Full)
		if err != nil {
			return nil, err
		}

		size += n
		storeList[v.Full] = short
	}

	if err := buffer.Flush(); err != nil {
		return nil, err
	}
	p.fileSize += int64(size)

	return storeList, nil
}

func (p *FileLoader) GetShort(ctx context.Context, full string) (string, error) {

	if err := p.create(); err != nil {
		return "", err
	}

	p.Lock()
	defer p.Unlock()

	file, err := os.OpenFile(p.fileName, os.O_WRONLY, 0644)
	if err != nil {
		return "", err
	}
	defer file.Close()

	if err := p.preSave(file); err != nil {
		return "", err
	}
	defer p.postSave(file)

	buffer := bufio.NewWriter(file)

	short, n, err := p.writeItem(buffer, full)
	if err != nil {
		return "", err
	}
	if err := buffer.Flush(); err != nil {
		return "", err
	}

	p.fileSize += int64(n)

	return short, nil
}

func (p *FileLoader) GetFull(ctx context.Context, short string) (string, error) {
	return "", fmt.Errorf("full not found for short %s", short)
}

func (p *FileLoader) preSave(file *os.File) error {
	filesize := int64(p.fileSize) - 1
	if err := file.Truncate(filesize); err != nil {
		return err
	}
	p.fileSize = filesize

	file.Seek(0, io.SeekEnd)
	return nil
}

func (p *FileLoader) postSave(file *os.File) {
	n, err := file.Write([]byte("]"))
	if err != nil {
		return
	}
	p.fileSize += int64(n)
}

func (p *FileLoader) writeItem(buffer *bufio.Writer, full string) (string, int, error) {
	short, err := rand.CreateShort(config.ShortLen)
	if err != nil {
		return "", 0, err
	}

	item := model.StoreItem{
		Full:  full,
		Short: short,
	}

	sendData, err := json.Marshal(&item)
	if err != nil {
		return "", 0, err
	}

	size := 0
	if p.fileSize > 2 {
		n, err := buffer.Write([]byte(",\n"))
		if err != nil {
			return "", 0, err
		}
		size += n
	}

	n, err := buffer.Write(sendData)
	if err != nil {
		return "", 0, err
	}
	size += n
	return short, size, nil
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
