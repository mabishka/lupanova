package fileloader

import (
	"encoding/json"
	"os"
	"sync"
)

type FileLoader struct {
	*sync.Mutex
	fileName string
	fileSize int64
}

type fileData struct {
	Short string `json:"short"`
	Full  string `json:"full"`
}

func New(fileName string) *FileLoader {
	return &FileLoader{Mutex: &sync.Mutex{}, fileName: fileName}
}

// return map [short string] full string
func (p *FileLoader) Load() (map[string]string, error) {
	p.create()

	p.Lock()
	defer p.Unlock()

	content, err := os.ReadFile(p.fileName)
	if err != nil {
		return nil, err
	}

	var data []fileData

	if err := json.Unmarshal(content, &data); err != nil {
		return nil, err
	}

	response := make(map[string]string)
	for _, v := range data {
		response[v.Short] = v.Full
	}

	return response, nil
}

func (p *FileLoader) Store(full, short string) error {

	p.create()

	p.Lock()
	defer p.Unlock()

	if err := os.Truncate(p.fileName, int64(p.fileSize)-1); err != nil {
		return err
	}
	p.fileSize--

	file, err := os.OpenFile(p.fileName, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	fileData := fileData{
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

func (p *FileLoader) create() error {

	if p.exist() {
		return nil
	}

	p.Lock()
	defer p.Unlock()

	os.WriteFile(p.fileName, []byte("[]"), 0644)
	p.fileSize = 2

	return nil
}

func (p *FileLoader) exist() bool {

	if p.fileSize != 0 {
		return true
	}

	p.Lock()
	defer p.Unlock()

	stat, err := os.Stat(p.fileName)
	if err != nil {
		return false
	}
	p.fileSize = stat.Size()

	return p.fileSize != 0
}
