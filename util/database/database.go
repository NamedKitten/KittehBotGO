package database


import (
	"github.com/tidwall/buntdb"
	//log "github.com/sirupsen/logrus"
)


var db *buntdb.DB

func init() {
	database, _ := buntdb.Open("data.db")
	db = database
}

func Set(key string, value string) {
	db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(key, value, nil)
		return err
	})
}

func Get(key string) string {
	var value string
	db.View(func(tx *buntdb.Tx) error {
		val, err := tx.Get(key)
		if err == buntdb.ErrNotFound {
			value = ""
		} else {
			value = val
		}
		return nil
	})
	return value
}