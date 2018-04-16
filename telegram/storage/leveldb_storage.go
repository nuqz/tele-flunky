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

	if err = db.Put(user.StorageKey(), value, nil); err != nil {
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

	u := new(access.User)
	if err := json.Unmarshal(userBs, u); err != nil {
		return nil, err
	}

	return u, nil
}

func (db *LevelDB) DeleteUser(user *access.User) error {
	return db.Delete(user.StorageKey(), nil)
}

func (db *LevelDB) SetUserNextChatMessageHandler(
	user *access.User,
	chatID int64,
	handlerName string,
) error {
	return db.Put(
		[]byte(fmt.Sprintf("user_%d_chat_%d_message_handler", user.ID, chatID)),
		[]byte(handlerName),
		nil,
	)
}

func (db *LevelDB) GetUserNextChatMessageHandler(
	user *access.User,
	chatID int64,
) (string, error) {
	bs, err := db.Get(
		[]byte(fmt.Sprintf("user_%d_chat_%d_message_handler", user.ID, chatID)),
		nil,
	)

	if err == leveldb.ErrNotFound {
		return "", nil
	} else if err != nil {
		return "", err
	}

	return string(bs), nil
}
