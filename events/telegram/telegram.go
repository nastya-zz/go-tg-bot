package telegram

import (
	"errors"
	"example/hello/clients/telegram"
	"example/hello/events"
	"example/hello/lib/e"
	"example/hello/storage"
)

type Processor struct {
	tg      *telegram.Client
	offset  int
	storage storage.Storage
}

type Meta struct {
	ChatId   int    `json:"chat_id"`
	UserName string `json:"username"`
}

var ErrUnknownEventType = errors.New("can't process message")
var ErrUnknownMetaType = errors.New("can't get meta")

func New(client *telegram.Client, storage storage.Storage) *Processor {
	return &Processor{
		tg:      client,
		storage: storage,
	}
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, e.Wrap("can't get events", err)
	}

	if len(updates) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(updates))

	for _, u := range updates {
		res = append(res, event(u))
	}

	p.offset = updates[len(updates)-1].ID + 1
	return res, nil
}

func (p *Processor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(event)
	default:
		return e.Wrap("can't process message", ErrUnknownEventType)
	}
}

func (p *Processor) processMessage(event events.Event) (err error) {
	defer func() { err = e.WrapIfErr("can't process message", err) }()

	metaType, err := meta(event)
	if err != nil {
		return err
	}

	if err := p.doCmd(event.Text, metaType.ChatId, metaType.UserName); err != nil {
		return err
	}

	return nil

}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, e.Wrap("can't get meta", ErrUnknownMetaType)
	}
	return res, nil
}

func event(u telegram.Update) events.Event {
	uType := fetchType(u)
	res := events.Event{
		Type: uType,
		Text: fetchText(u),
	}

	if uType == events.Message {
		res.Meta = Meta{
			ChatId:   u.Message.Chat.ID,
			UserName: u.Message.From.UserName,
		}
	}

	return res
}

func fetchText(u telegram.Update) string {
	if u.Message == nil {
		return ""
	}
	return u.Message.Text
}

func fetchType(u telegram.Update) events.Type {
	if u.Message == nil {
		return events.Unknown
	}
	return events.Message
}
