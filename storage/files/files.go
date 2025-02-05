package files

import (
	"encoding/gob"
	"errors"
	"example/hello/lib/e"
	"example/hello/storage"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

type Storage struct {
	basePath string
}

const (
	defaultDirPerm = 0774
)

var ErrNoSavedPages = errors.New("no saved pages")

func New(basePath string) Storage {
	return Storage{basePath: basePath}
}

func (s Storage) Save(page *storage.Page) (err error) {
	defer func() { err = e.WrapIfErr("can't save page", err) }()

	filePath := filepath.Join(s.basePath, page.UserName)
	if err = os.MkdirAll(filePath, defaultDirPerm); err != nil {
		return err
	}

	fName, err := fileName(page)
	if err != nil {
		return err
	}

	filePath = filepath.Join(filePath, fName)

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}

	defer func() { _ = file.Close() }()

	if err := gob.NewEncoder(file).Encode(page); err != nil {
		return err
	}

	return nil
}

func (s Storage) PickRandom(userName string) (page *storage.Page, err error) {
	defer func() { err = e.WrapIfErr("can't pick random", err) }()

	filePath := filepath.Join(s.basePath, page.UserName)

	files, err := os.ReadDir(filePath)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, ErrNoSavedPages
	}

	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(len(files))

	file := files[n]

	//open

}

func (s Storage) decodePage(filePath string) {}

func fileName(page *storage.Page) (string, err) {
	return page.Hash()
}
