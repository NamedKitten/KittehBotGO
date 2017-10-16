package database

import (
	"fmt"
	"github.com/tidwall/buntdb"
	"github.com/tidwall/redcon"
	"log"
	"strings"
)

func Start(file string, port int) {
	db, _ := buntdb.Open(file)
	addr := fmt.Sprintf(":%d", port)
	go log.Printf("started server at %s", addr)
	err := redcon.ListenAndServe(addr,
		func(conn redcon.Conn, cmd redcon.Command) {
			switch strings.ToLower(string(cmd.Args[0])) {
			default:
				conn.WriteError("ERR unknown command '" + string(cmd.Args[0]) + "'")
			case "detach":
				hconn := conn.Detach()
				log.Printf("connection has been detached")
				go func() {
					defer hconn.Close()
					hconn.WriteString("OK")
					hconn.Flush()
				}()
				return
			case "ping":
				conn.WriteString("PONG")
			case "quit":
				conn.WriteString("OK")
				conn.Close()
			case "set":
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
			// use this function to accept or deny the connection.
			// log.Printf("accept: %s", conn.RemoteAddr())
			return true
		},
		func(conn redcon.Conn, err error) {
			// this is called when the connection has been closed
			// log.Printf("closed: %s, err: %v", conn.RemoteAddr(), err)
		},
	)
	if err != nil {
		log.Fatal(err)
	}
}
