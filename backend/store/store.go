package store

import (
	"encoding/json"
	"github.com/dlyahov/startuplink-web-go/backend/model"
	"github.com/jessevdk/go-flags"
	"github.com/pkg/errors"
	bolt "go.etcd.io/bbolt"
	"log"
	"os"
	"time"
)

const dbFileName = "app.db"

type BoltConfig struct {
	Path    string `long:"path" env:"STORE_BOLT_DB" default:"." description:"parent dir for bolt files"`
	Timeout int    `long:"timeout" env:"TIMEOUT" default:"3000" description:"bolt timeout in milliseconds"`
}

type Storage interface {
	SaveUser(user *model.User) error
	FindUser(id string) (*model.User, error)
}

var (
	ErrUserNotFound = errors.New("user not found")
)

type BoltDb struct {
	db *bolt.DB
}

const userBuckets = "users"

func NewStorage() (*BoltDb, error) {
	boltConfig := &BoltConfig{}
	if _, err := flags.Parse(boltConfig); err != nil {
		log.Println("Could not parse bolt db configuration")
		return nil, err
	}

	if _, err := os.Stat(boltConfig.Path); os.IsNotExist(err) {
		if err := os.Mkdir(boltConfig.Path, 0700); err != nil {
			log.Println("Cannot create bolt db path. ", err)
			return nil, err
		}
	}
	dbPath := boltConfig.Path + "/" + dbFileName
	options := &bolt.Options{Timeout: time.Duration(boltConfig.Timeout) * time.Millisecond}
	log.Println("Bolt options: ", options)
	boltDb, err := bolt.Open(dbPath, 0666, options)
	if err != nil {
		log.Println("Could not open database file")
		return nil, err
	}
	err = boltDb.Update(func(tx *bolt.Tx) error {
		if _, e := tx.CreateBucketIfNotExists([]byte(userBuckets)); e != nil {
			return errors.Wrapf(e, "failed to create bucket %s", userBuckets)
		}
		return nil
	})

	if err != nil {
		return nil, errors.Wrap(err, "failed to create bucket)")
	}

	log.Printf("Bolt store created")
	return &BoltDb{boltDb}, nil
}

func (boltDb *BoltDb) SaveUser(user *model.User) error {
	db := boltDb.db
	err := db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(userBuckets))
		userJson, err := json.Marshal(user)
		if err != nil {
			log.Println("Could not marshal user to json")
			return err
		}
		err = bucket.Put([]byte(user.Id), userJson)
		if err != nil {
			return errors.Wrapf(err, "failed to save user data. User id %s", user.Id)
		}
		return nil
	})

	if err != nil {
		log.Println("failed to save user. ", err)
		return err
	}
	return nil
}

func (boltDb *BoltDb) FindUser(id string) (*model.User, error) {
	db := boltDb.db

	user := model.User{}
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(userBuckets))
		userBytes := bucket.Get([]byte(id))
		if userBytes == nil {
			return ErrUserNotFound
		}

		err := json.Unmarshal(userBytes, &user)
		if err != nil {
			return errors.Wrapf(err, "Could not marshal user with id %s", id)
		}
		return nil
	})

	if err != nil {
		log.Println("failed to find user. ", err)
		return &user, err
	}
	return &user, nil
}
