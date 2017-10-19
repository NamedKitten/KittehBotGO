package database

import (
	"fmt"
	"github.com/tidwall/buntdb"
	"github.com/tidwall/redcon"
	"log"
	"strings"
)

var Password string = ""

type Conn struct {
	IsAuthenticated bool
}

func Start(file string, port int) {
	db, _ := buntdb.Open(file)
	addr := fmt.Sprintf("0.0.0.0:%d", port)
	go log.Printf("InternalDB: Started server at %s", addr)
	err := redcon.ListenAndServe(addr,
		func(conn redcon.Conn, cmd redcon.Command) {
			switch strings.ToLower(string(cmd.Args[0])) {
			default:
				buff := ""
				for index, arg := range cmd.Args {
					buff += string(arg)
					if index+1 != len(cmd.Args) {
						buff += " "
					}
				}
				conn.WriteError("ERR unknown command '" + buff + "'")
			case "detach":
				hconn := conn.Detach()
				go log.Println("InternalDB: Cnnection has been detached.")
				go func() {
					defer hconn.Close()
					hconn.WriteString("OK")
					hconn.Flush()
				}()
				return
			case "select":
				go log.Println("DB select not implemented.")
				conn.WriteString("OK")
				return
			case "ping":
				conn.WriteString("PONG")
				return
			case "quit":
				conn.WriteString("OK")
				conn.Close()
				return
			case "auth":
				if Password == "" {
					conn.WriteString("OK")
					return
				}
				if len(cmd.Args) != 2 {
					conn.WriteError("ERR wrong number of arguments for '" + string(cmd.Args[0]) + "' command")
					return
				}
				if string(cmd.Args[1]) == Password {
					conn.WriteString("OK")
					conn.SetContext(true)
				} else {
					conn.WriteError("ERR invalid password")
				}
				return
			case "set":
				if conn.Context() != true {
					conn.WriteError("NOAUTH Authentication required.")
					return
				}
				if len(cmd.Args) != 3 {
					conn.WriteError("ERR wrong number of arguments for '" + string(cmd.Args[0]) + "' command")
					return
				}
				err := db.Update(func(tx *buntdb.Tx) error {
					_, _, err := tx.Set(string(cmd.Args[1]), string(cmd.Args[2]), nil)
					return err
				})
				//mu.Lock()
				//items[string(cmd.Args[1])] = cmd.Args[2]
				//mu.Unlock()
				if err == nil {
					conn.WriteString("OK")
				}
			case "get":
				if conn.Context() != true {
					conn.WriteError("NOAUTH Authentication required.")
					return
				}
				if len(cmd.Args) != 2 {
					conn.WriteError("ERR wrong number of arguments for '" + string(cmd.Args[0]) + "' command")
					return
				}
				err := db.View(func(tx *buntdb.Tx) error {
					val, err := tx.Get(string(cmd.Args[1]))
					if err != nil {
						conn.WriteNull()
					} else {
						conn.WriteString(val)
					}
					return nil
				})
				if err != nil {
					panic(err)
				}

			case "del":
				if conn.Context() != true {
					conn.WriteError("(error) NOAUTH Authentication required.")
					return
				}
				if len(cmd.Args) != 2 {
					conn.WriteError("ERR wrong number of arguments for '" + string(cmd.Args[0]) + "' command")
					return
				}
				err := db.Update(func(tx *buntdb.Tx) error {
					_, err := tx.Delete(string(cmd.Args[1]))
					return err
				})
				if err != nil {
					conn.WriteInt(0)
				} else {
					conn.WriteInt(1)
				}
			}
		},
		func(conn redcon.Conn) bool {
			go conn.SetContext(false)
			go log.Printf("InternalDB: Accept %s", conn.RemoteAddr())
			return true
		},
		func(conn redcon.Conn, err error) {
			log.Printf("InternalDB: Disconnect %s, err: %v", conn.RemoteAddr(), err)
		},
	)
	if err != nil {
		log.Fatal(err)
	}
}
