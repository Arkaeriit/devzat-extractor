package main

import (
	"os"
    "fmt"

	api "github.com/quackduck/devzat/devzatapi"
)

func main() {
	s, err := api.NewSession("devzat.hackclub.com:5556", os.Getenv("DEVZAT_TOKEN"))
	if err != nil {
		panic(err)
	}

    messageChan, _, err := s.RegisterListener(false, false, "")
	if err != nil {
		panic(err)
	}

	for {
		select {
		case err = <-s.ErrorChan:
			panic(err)
		case msg := <-messageChan:
            fmt.Println(msg);
        }
    }

}
