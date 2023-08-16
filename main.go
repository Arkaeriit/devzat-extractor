package main

import (
	"fmt"
	"os"
	//"strconv"
	"time"

	//"github.com/gin-gonic/gin"
	api "github.com/quackduck/devzat/devzatapi"
)

type TimedMsg struct {
	msg api.Message
	ts  int64
}

// A circular buffer used to store the latest messages
type messageBank struct {
	size  int
	index int
	msgs  []TimedMsg
}

// Generates a circular buffer of the given size
func makeBank(size int) messageBank {
	return messageBank{
		size:  size,
		index: 0,
		msgs:  make([]TimedMsg, size),
	}
}

// Return the nth latest message. readNthPreviousMsg(0) returns the last
// message, readNthPreviousMsg(1) returns the one before; and so on.
// If no suitable message is found, return nil.
func (b *messageBank) readNthPreviousMsg(index int) *TimedMsg {
	index += 1
	if index >= b.size {
		return nil
	}
	if b.index-index < 0 {
		return nil
	}
	return &b.msgs[(b.index-index)%b.size]
}

// Adds a new message to the buffer.
func (b *messageBank) addMessage(msg TimedMsg) {
	b.msgs[b.index%b.size] = msg
	b.index = b.index + 1
}

// Takes all the messages starting by the last one and finishing by the
// `count`th and format them together. If fromRoom is a name, it will only
// takes the messages from that room and if it is an empty string, it will
// get them from every rooms.
func (b *messageBank) compilePreviousMsg(count int, fromRoom string) string {
	ret := ""
	i := b.size - 1
	for i >= 0 && count > 0 {
		msg := b.readNthPreviousMsg(i)
		if msg != nil && (fromRoom == "" || msg.msg.Room == fromRoom) {
			ret = ret + formatMsg(msg.msg)
			count = count - 1
		}
		i = i - 1
	}
	return ret
}

// Format a single message.
func formatMsg(msg api.Message) string {
	return fmt.Sprintf("%v: %v  \n", msg.From, msg.Data)
}

// Takes a message and adds the current timestamp to it.
func timeMessage(msg api.Message) TimedMsg {
	return TimedMsg{
		msg: msg,
		ts:  time.Now().Unix(),
	}
}

// Takes a string representing a duration and returns the timestamp that was
// this duration ago. Return nil if the duration is not valid.
func timestampWhenDuration(msg string) *int64 {
	if msg[0] != '-' {
		msg = "-" + msg
	}
	duration, err := time.ParseDuration(msg)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	now := time.Now()
	then := now.Add(duration).Unix()
	return &then
}

// Returned the compiled messages between the two timestamps.
func (b *messageBank) messagesBetween(ts_start int64, ts_stop int64, fromRoom string) string {
	ret := ""
	for i := b.size - 1; i >= 0; i-- {
		msg := b.readNthPreviousMsg(i)
		if msg == nil {
			continue
		}
		if ts_stop < msg.ts {
			break
		}
		if ts_start < msg.ts {
			if fromRoom == "" || fromRoom == msg.msg.Room {
				ret = ret + formatMsg(msg.msg)
			}
		}
	}
	return ret
}

func main() {
	session, err := api.NewSession("devzat.hackclub.com:5556", os.Getenv("DEVZAT_TOKEN"))
	if err != nil {
		panic(err)
	}
	bank := makeBank(50)

	// Read all incoming messages
	go func() {
		messageChan, _, err := session.RegisterListener(false, false, "")
        fmt.Println(messageChan)
        fmt.Println(err)
		if err != nil {
			panic(err)
		}

		for {
			select {
			case err = <-session.ErrorChan:
				panic(err)
			case msg := <-messageChan:
                fmt.Println(msg)
				bank.addMessage(timeMessage(msg))
			}
		}
	}()

	// tmp
	err = session.RegisterCmd("extract", "duration", "Extract the messages posted in `duration`",
		func(cmdCall api.CmdCall, err error) {
			fmt.Println("???")
			from := timestampWhenDuration(cmdCall.Args)
			if from == nil {
				err := session.SendMessage(api.Message{Room: cmdCall.Room, From: "Devzat-extractor", Data: "Error, invalid duration", DMTo: ""})
				if err != nil {
					panic(err)
				}
			}
			url := fmt.Sprintf("http://localhost:8080/timespan/%v/%v/%v", cmdCall.Room, *from, *timestampWhenDuration("-1"))
            er := session.SendMessage(api.Message{Room: cmdCall.Room, From: "Devzat-extractor", Data: url, DMTo: ""})
			if er != nil {
				panic(err)
			}
		})
	if err != nil {
		panic(err)
	}

    /*
	go func() {
		router := gin.Default()
		router.GET("/timespan/:room/:from/:to", func(c *gin.Context) {
			room := c.Param("room")
			if room == "all" {
				room = ""
			}
			from, err := strconv.ParseInt(c.Param("from"), 10, 64)
			if err != nil {
				c.String(400, "Error: %v", err)
			}
			to, err := strconv.ParseInt(c.Param("to"), 10, 64)
			if err != nil {
				c.String(400, "Error: %v", err)
			}

			c.String(200, "%v", bank.messagesBetween(from, to, room))
		})

		router.Run("localhost:8080")
	}()
    */

	// Debug
	for {
		time.Sleep(10 * time.Second)
        fmt.Printf("<%v>\n", bank.compilePreviousMsg(30, ""))
	}

}

