package util

import (
	"encoding/json"
	"strconv"
	"strings"

	"golang.org/x/net/websocket"

	t "github.com/krishnakantha1/expenseTrackerIngestion/types"
)

func ParseEvent(str string) string {
	len := len(str)
	space := strings.Index(str, " ")

	if space != -1 {
		len = space
	}

	return str[:len]
}

func ParseUser(message string) (*t.UserData, error) {
	message = message[strings.Index(message, " ")+1:]

	return DecodeJWT(message)
}

func ParseCount(str string) (int, error) {
	return strconv.Atoi(str[strings.Index(str, " ")+1:])
}

func ParseExpenseMessaeg(message string) *t.ExpenseMessage {
	message = message[strings.Index(message, " ")+1:]

	em := new(t.ExpenseMessage)
	json.Unmarshal([]byte(message), em)

	return em
}

func HandleStateMisMatch(ws *websocket.Conn) {
	ws.Write([]byte("Server Error. Did not recive expected event message."))
}
