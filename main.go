package main

import (
    "fmt"
    "os"
    "time"

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

func (b *messageBank) readNthPreviousMsg(index int) (TimedMsg, bool) {
    index += 1
    if index >= b.size {
        return b.msgs[0], false
    }
    if b.index - index < 0 {
        return b.msgs[0], false
    }
    return b.msgs[(b.index - index) % b.size], true
}

func (b *messageBank) addMessage(msg TimedMsg) {
    b.msgs[b.index % b.size] = msg
    b.index = b.index + 1
}

func (b *messageBank) compilePreviousMsg(count int) string {
    ret := ""
    i := count - 1
    for i >= 0 {
        msg, Ok := b.readNthPreviousMsg(i)
        if Ok {
            ret = ret + formatMsg(msg.msg)
        }
        i = i - 1
    }
    return ret
}

func formatMsg(msg api.Message) string {
    return fmt.Sprintf("%v %v: %v\n", msg.Room, msg.From, msg.Data)
}

func timeMessage(msg api.Message) TimedMsg {
    return TimedMsg{
        msg: msg,
        ts: time.Now(),
    }
}

func main() {
    s, err := api.NewSession("devzat.hackclub.com:5556", os.Getenv("DEVZAT_TOKEN"))
    if err != nil {
        panic(err)
    }

    bank := makeBank(5)

    go func() {
        messageChan, _, err := s.RegisterListener(false, false, "")
        if err != nil {
            panic(err)
        }

        for {
            select {
            case err = <-s.ErrorChan:
                panic(err)
            case msg := <-messageChan:
                bank.addMessage(timeMessage(msg))
            }
        }
    }()

    for {
        time.Sleep(10* time.Second)
        fmt.Printf("<%v>\n", bank.compilePreviousMsg(3))
    }

}

