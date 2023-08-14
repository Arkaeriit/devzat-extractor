package main

import (
    "os"
    "fmt"
    "time"
    "strconv"

    api "github.com/quackduck/devzat/devzatapi"
)

type TimedMsg struct {
    msg api.Message
    ts  time.Time
}

type messageBank struct {
    size  int
    index int
    msgs  []TimedMsg
}

func makeBank(size int) messageBank {
    return messageBank{
        size: size,
        index: 0,
        msgs: make([]TimedMsg, size),
    }
}

func (b *messageBank) readNthPreviousMsg(index int) *TimedMsg {
    index += 1
    if index >= b.size {
        return nil
    }
    if b.index - index < 0 {
        return nil
    }
    return &b.msgs[(b.index - index) % b.size]
}

func (b *messageBank) addMessage(msg TimedMsg) {
    b.msgs[b.index % b.size] = msg
    b.index = b.index + 1
}

func (b *messageBank) compilePreviousMsg(count int, fromRoom string) string {
    ret := ""
    i := count - 1
    for i >= 0 {
        msg := b.readNthPreviousMsg(i)
        if msg != nil && (fromRoom == "" || msg.msg.Room == fromRoom) {
            ret = ret + formatMsg(msg.msg)
        }
        i = i - 1
    }
    return ret
}

func formatMsg(msg api.Message) string {
    return fmt.Sprintf("%v: %v  \n", msg.From, msg.Data)
}

func timeMessage(msg api.Message) TimedMsg {
    return TimedMsg{
        msg: msg,
        ts: time.Now(),
    }
}

var bank messageBank
var session *api.Session

func extractCmd(cmdCall api.CmdCall, err error) {
    if err != nil {
        panic(err)
    }
    count, err := strconv.Atoi(cmdCall.Args)
    if err != nil {
        panic(err)
    }
    err = session.SendMessage(api.Message{Room: cmdCall.Room, From: "Devzat-extractor", Data: bank.compilePreviousMsg(count, cmdCall.Room), DMTo: ""})
    if err != nil {
        panic(err)
    }
}

func main() {
    var err error
    session, err = api.NewSession("devzat.hackclub.com:5556", os.Getenv("DEVZAT_TOKEN"))
    if err != nil {
        panic(err)
    }

    bank = makeBank(5)

    go func() {
        messageChan, _, err := session.RegisterListener(false, false, "")
        if err != nil {
            panic(err)
        }

        for {
            select {
            case err = <-session.ErrorChan:
                panic(err)
            case msg := <-messageChan:
                bank.addMessage(timeMessage(msg))
            }
        }
    }()

    err = session.RegisterCmd("extract", "duration", "Extract the messages posted in `duration`", extractCmd)
    if err != nil {
        panic(err)
    }

    for {
        //time.Sleep(10* time.Second)
        //fmt.Printf("<%v>\n", bank.compilePreviousMsg(3))
    }

}

