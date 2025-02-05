package files

import (
	"encoding/gob"
	"errors"
	"example/hello/lib/e"
	"example/hello/storage"
	"fmt"
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
		return nil, storage.ErrNoSavedPages
	}

	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(len(files))

	file := files[n]

	return s.decodePage(filepath.Join(s.basePath, file.Name()))
}

func (s Storage) Remove(page *storage.Page) (err error) {
	fName, err := fileName(page)
	if err != nil {
		msg := fmt.Sprintf("can't search  for remove page %s", fName)
		return e.Wrap(msg, err)
	}

	filePath := filepath.Join(s.basePath, page.UserName, fName)
	if err = os.Remove(filePath); err != nil {
		msg := fmt.Sprintf("can't remove page by path:  %s", filePath)
		return e.Wrap(msg, err)
	}

	return nil
}

func (s Storage) IsExists(page *storage.Page) (bool, error) {
	fName, err := fileName(page)
	if err != nil {
		msg := fmt.Sprintf("can't search page %s", fName)
		return false, e.Wrap(msg, err)
	}

	filePath := filepath.Join(s.basePath, page.UserName, fName)

	switch _, err = os.Stat(filePath); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		msg := fmt.Sprintf("can't check if file exists: %s", filePath)
		return false, e.Wrap(msg, err)
	}

	return true, nil
}

func (s Storage) decodePage(filePath string) (p *storage.Page, err error) {
	defer func() { err = e.WrapIfErr("can't decode file", err) }()

	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()

	var page storage.Page

	if err := gob.NewDecoder(f).Decode(&page); err != nil {
		return nil, err
	}

	return &page, nil
}

func fileName(page *storage.Page) (string, error) {
	return page.Hash()
}
