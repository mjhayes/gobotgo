package main

import (
	"fmt"
	"os"
	"os/signal"
)

func main() {
	s, err := IRCSessCreate("gobotgo", "gobotgo_", "hi", "hey", "irc.freenode.net:6667", 240, EventHandler)
	if err != nil {
		fmt.Println("Session creation failed:", err)
		return
	}

	defer s.Destroy("Leaving")

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc)

	for sig := range sigc {
		switch sig {
		case os.Interrupt:
			return
		case os.Kill:
			return
		}
	}
}

func EventHandler(s *IRCSess, e IRCSessEvent) {
	switch e.Type {
	case IRCSessEventConnect:
		fmt.Println("-- Connected to", s.ServerAddress, "--")
	case IRCSessEventDisconnect:
		fmt.Println("-- Disconnected from", s.ServerAddress, "--")
	case IRCSessEventMessage:
		HandleMessage(s, e.Message)
	}
}

func HandleMessage(s *IRCSess, m IRCMessage) {	
	switch m.Type {
	case IRCMessageRAW:
		switch m.Raw {
		case "001":
			s.Send("JOIN #notaboutlegos")
		}
	}
}

// vim: set noexpandtab ts=8 sw=8 sts=8:
