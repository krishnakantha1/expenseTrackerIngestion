package websockethandler

import (
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/websocket"

	dw "github.com/krishnakantha1/expenseTrackerIngestion/database/data-wrapper"
	t "github.com/krishnakantha1/expenseTrackerIngestion/types"
	util "github.com/krishnakantha1/expenseTrackerIngestion/util"
)

const (
	STATE_START      = "START"
	CLIENT_START_MSG = "START"
	ACK_START        = "ACK_START"

	STATE_NOM      = "NUMBER_OF_MESSAGE"
	CLIENT_NOM_MSG = "NUMBER_OF_MESSAGE"
	ACK_NOM        = "ACK_NOM"

	STATE_READ_ENTRY  = "READ_ENTRY"
	CLIENT_READ_ENTRY = "READ_ENTRY"
	ACK_READ_ENTRY    = "ACK_READ_ENTRY"

	STATE_SAVE_DONE = "SAVE_DONE"
	MSG_SAVE_DONE   = "SAVE_DONE"
	ACK_SAVE_DONE   = "ACK_SAVE_ENTRY"
)

func IngestExpenseReadLoop(db *mongo.Client, ws *websocket.Conn) {
	defer ws.Close()

	buf := make([]byte, 2048)
	done := false
	state := STATE_START
	messageCount := 0
	curCount := 0
	var user *t.UserData
	var expenseMessages []*t.ExpenseMessage

	for !done {
		n, err := ws.Read(buf)

		if err != nil {
			break
		}

		msg := string(buf[:n])
		event := util.ParseEvent(msg)

		switch state {
		case STATE_START:
			if event == CLIENT_START_MSG {
				log.Println("in start")
				state = STATE_NOM
				messageCount = 0
				curCount = 0
				user, err = util.ParseUser(msg)

				if err != nil {
					log.Println(err)
					util.HandleStateMisMatch(ws)
					return
				}

				ws.Write([]byte(ACK_START))
			} else {
				util.HandleStateMisMatch(ws)
				return
			}

		case STATE_NOM:
			if event == CLIENT_NOM_MSG {
				messageCount, err = util.ParseCount(msg)
				log.Println("in nom", messageCount)

				if err != nil {
					util.HandleStateMisMatch(ws)
					return
				}
				expenseMessages = make([]*t.ExpenseMessage, 0, messageCount)
				state = STATE_READ_ENTRY

				ws.Write([]byte(ACK_NOM))
			} else {
				util.HandleStateMisMatch(ws)
				return
			}

		case STATE_READ_ENTRY:
			if event == CLIENT_READ_ENTRY {
				log.Println("in client read", curCount)
				expenseMessages[curCount] = util.ParseExpenseMessaeg(msg)
				ws.Write([]byte(fmt.Sprintf("%s %d", ACK_READ_ENTRY, curCount)))
				curCount++

				if curCount == messageCount {
					dw.UpsertExpenseMessages(db, expenseMessages, user.ID)
					state = STATE_SAVE_DONE
					ws.Write([]byte(MSG_SAVE_DONE))
				}
			} else {
				util.HandleStateMisMatch(ws)
				return
			}

		case STATE_SAVE_DONE:
			if event == ACK_SAVE_DONE {
				return
			} else {
				util.HandleStateMisMatch(ws)
				return
			}

		default:
			util.HandleStateMisMatch(ws)
			return
		}

	}
}
