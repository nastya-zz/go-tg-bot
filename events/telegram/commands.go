package telegram

import (
	"errors"
	"example/hello/lib/e"
	"example/hello/storage"
	"log"
	"net/url"
	"strings"
)

const (
	RndCmd   = "/rnd"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

func (p *Processor) doCmd(text string, chatId int, username string) error {
	text = strings.TrimSpace(text)
	log.Printf("got new message '%s' from %s", text, username)

	if isAddCmd(text) {
		return p.savePage(chatId, text, username)
	}

	switch text {
	case RndCmd:
		return p.sendRandom(chatId, text, username)
	case HelpCmd:
		return p.sendHelp(chatId)
	case StartCmd:
		return p.sendStart(chatId)
	default:
		return p.tg.SendMassage(chatId, msgUnknownCommand)
	}
}

func (p *Processor) savePage(chatID int, pageURL string, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: save page", err) }()

	page := &storage.Page{
		URL:      pageURL,
		UserName: username,
	}

	isExist, err := p.storage.IsExist(page)
	if err != nil {
		return err
	}

	if isExist {
		return p.tg.SendMassage(chatID, msgAlreadyExists)
	}

	if err = p.storage.Save(page); err != nil {
		return err
	}
	if err = p.tg.SendMassage(chatID, msgSaved); err != nil {
		return err
	}

	return nil

}

func (p *Processor) sendRandom(chatID int, pageURL string, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: send random page", err) }()

	page, err := p.storage.PickRandom(username)
	if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {
		return err
	}
	if errors.Is(err, storage.ErrNoSavedPages) {
		return p.tg.SendMassage(chatID, msgNoSavedPages)
	}

	if err = p.tg.SendMassage(chatID, page.URL); err != nil {
		return err
	}
	return p.storage.Remove(page)
}

func (p *Processor) sendHelp(chatID int) (err error) {
	return p.tg.SendMassage(chatID, msgHelp)
}

func (p *Processor) sendStart(chatID int) (err error) {
	return p.tg.SendMassage(chatID, msgHello)
}

func isAddCmd(text string) bool {
	return isUrl(text)
}

func isUrl(text string) bool {
	u, err := url.Parse(text)
	return err == nil && u.Scheme != "" && u.Host != ""
}
