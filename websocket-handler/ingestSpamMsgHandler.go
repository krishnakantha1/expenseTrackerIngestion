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
	STATE_START_SPAM      = "START"
	CLIENT_START_MSG_SPAM = "START"
	ACK_START_SPAM        = "ACK_START"

	STATE_NOM_SPAM      = "NUMBER_OF_MESSAGE"
	CLIENT_NOM_MSG_SPAM = "NUMBER_OF_MESSAGE"
	ACK_NOM_SPAM        = "ACK_NOM"

	STATE_READ_ENTRY_SPAM  = "READ_ENTRY"
	CLIENT_READ_ENTRY_SPAM = "READ_ENTRY"
	ACK_READ_ENTRY_SPAM    = "ACK_READ_ENTRY"

	STATE_SAVE_DONE_SPAM = "SAVE_DONE"
	MSG_SAVE_DONE_SPAM   = "SAVE_DONE"
	ACK_SAVE_DONE_SPAM   = "ACK_SAVE_ENTRY"
)

func IngestSpamMsgReadLoop(db *mongo.Client, ws *websocket.Conn) {
	defer ws.Close()

	buf := make([]byte, 4096)
	done := false
	state := STATE_START_SPAM
	messageCount := 0
	curCount := 0
	var rawMessages []*t.RawMessage

	for !done {
		n, err := ws.Read(buf)

		if err != nil {
			break
		}

		msg := string(buf[:n])
		event := util.ParseEvent(msg)

		switch state {
		case STATE_START_SPAM:
			if event == CLIENT_START_MSG_SPAM {
				state = STATE_NOM_SPAM
				messageCount = 0
				curCount = 0

				if err != nil {
					log.Println(err)
					util.HandleStateMisMatch(ws)
					return
				}

				ws.Write([]byte(ACK_START_SPAM))

			} else {
				util.HandleStateMisMatch(ws)
				return
			}

		case STATE_NOM_SPAM:
			if event == CLIENT_NOM_MSG_SPAM {
				messageCount, err = util.ParseCount(msg)

				if err != nil {
					util.HandleStateMisMatch(ws)
					return
				}

				rawMessages = make([]*t.RawMessage, 0, messageCount)
				state = STATE_READ_ENTRY_SPAM

				ws.Write([]byte(ACK_NOM_SPAM))
			} else {
				util.HandleStateMisMatch(ws)
				return
			}

		case STATE_READ_ENTRY_SPAM:
			if event == CLIENT_READ_ENTRY_SPAM {
				rawMessages = append(rawMessages, util.ParseRawMessage(msg))

				ws.Write([]byte(fmt.Sprintf("%s %d", ACK_READ_ENTRY_SPAM, curCount)))
				curCount++

				if curCount == messageCount {
					dw.InsertRawMessages(db, rawMessages, "SPAM")
					state = STATE_SAVE_DONE_SPAM
					ws.Write([]byte(MSG_SAVE_DONE_SPAM))
				}
			} else {
				util.HandleStateMisMatch(ws)
				return
			}

		case STATE_SAVE_DONE_SPAM:
			if event == ACK_SAVE_DONE_SPAM {
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
