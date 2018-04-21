package storage

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"

	"github.com/nuqz/tele-flunky/telegram/access"
)

type LevelDB struct {
	*leveldb.DB
}

func NewLevelDBStorage(path string) (BotStorage, error) {
	opts := opt.Options{Filter: filter.NewBloomFilter(10)}
	db, err := leveldb.OpenFile(path, &opts)
	if err != nil && errors.IsCorrupted(err) {
		db, err = leveldb.RecoverFile(path, &opts)
		log.Println("Recovered corrupted LevedDB file")
	}

	if err != nil {
		return nil, err
	}

	return &LevelDB{db}, nil
}

func (db *LevelDB) PutUser(user *access.User) error {
	value, err := json.Marshal(user)
	if err != nil {
		return err
	}

	if err = db.Put([]byte(fmt.Sprintf("user_%d", user.ID)), value, nil); err != nil {
		return err
	}

	return nil
}

func (db *LevelDB) HasUser(user *tgbotapi.User) (bool, error) {
	return db.Has([]byte(fmt.Sprintf("user_%d", user.ID)), nil)
}

func (db *LevelDB) GetUser(user *tgbotapi.User) (*access.User, error) {
	userBs, err := db.Get([]byte(fmt.Sprintf("user_%d", user.ID)), nil)
	if err != nil {
		return nil, err
	}

	return access.NewUserFromJSON(userBs)
}

func (db *LevelDB) DeleteUser(user *access.User) error {
	return db.Delete([]byte(fmt.Sprintf("user_%d", user.ID)), nil)
}

func (db *LevelDB) SetUserNextChatMessageHandler(
	user *access.User,
	chat *tgbotapi.Chat,
	handlerName string,
) error {
	return db.Put(
		[]byte(fmt.Sprintf("user_%d_chat_%d_message_handler", user.ID, chat.ID)),
		[]byte(handlerName),
		nil,
	)
}

func (db *LevelDB) GetUserNextChatMessageHandler(
	user *access.User,
	chat *tgbotapi.Chat,
) (string, error) {
	bs, err := db.Get(
		[]byte(fmt.Sprintf("user_%d_chat_%d_message_handler", user.ID, chat.ID)),
		nil,
	)

	if err == leveldb.ErrNotFound {
		return "", nil
	} else if err != nil {
		return "", err
	}

	return string(bs), nil
}
