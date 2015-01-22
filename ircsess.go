package main

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

type IRCSess struct {
	Nick		string
	AlternateNick	string	
	UserName	string
	RealName	string
	ServerAddress	string
	ServerTimeout	int
	EventHandler	IRCSessEventFunc

	ConnInfo struct {
		Socket		net.Conn
		Connected	bool
		TimeConnected	time.Time
		TimeLastEvent	time.Time
	}
}

type IRCSessEventType int
const (
	_ = iota
	IRCSessEventConnect
	IRCSessEventDisconnect
	IRCSessEventMessage
)

type IRCSessEvent struct {
	Type	IRCSessEventType
	Message	IRCMessage
}

type IRCSessEventFunc func(*IRCSess, IRCSessEvent)


func IRCSessCreate(nick, alternateNick, userName, realName, serverAddress string, serverTimeout int, eventHandler IRCSessEventFunc) (s IRCSess, e error) {
	s.Nick = nick
	s.AlternateNick = alternateNick
	s.UserName = userName
	s.RealName = realName
	s.ServerAddress = serverAddress
	s.ServerTimeout = serverTimeout
	s.EventHandler = eventHandler

	e = s.Connect()
	if e != nil {
		return
	}

	go s.ConnectionThread()

	return
}

func (s *IRCSess) Destroy(quitMessage string) (e error) {
	if s.ConnInfo.Connected {
		s.Send("QUIT :", quitMessage)
		time.Sleep(time.Second)
		s.Disconnect(false)
		return
	}

	return
}

func (s *IRCSess) Connect() (e error) {
	if s.ConnInfo.Connected {
		return
	}

	s.ConnInfo.Socket, e = net.Dial("tcp", s.ServerAddress)
	if e != nil {
		return
	}

	s.ConnInfo.Connected = true
	s.ConnInfo.TimeConnected = time.Now()
	s.ConnInfo.TimeLastEvent = time.Now()

	s.Send("NICK ", s.Nick)
	s.Send("USER ", s.UserName, " 0 * :", s.RealName)

	ev := IRCSessEvent{ Type:IRCSessEventConnect }
	go s.EventHandler(s, ev)

	return
}

func (s *IRCSess) Disconnect(launchEventHandler bool) {
	if !s.ConnInfo.Connected {
		return
	}

	s.ConnInfo.Socket.Close()
	s.ConnInfo.Connected = false

	if launchEventHandler {
		ev := IRCSessEvent{ Type:IRCSessEventDisconnect }
		go s.EventHandler(s, ev)
	}
}

func (s *IRCSess) Send(v ...interface{}) {
	writer := bufio.NewWriter(s.ConnInfo.Socket)
	out := fmt.Sprintf("%s\r\n", fmt.Sprint(v...))

	_, err := writer.WriteString(out)
	if err != nil {
		fmt.Println("Couldn't write:", err)
		return
	}
	writer.Flush()

	fmt.Print(">> ", out)
}

func (s *IRCSess) ConnectionThread() {
	reader := bufio.NewReader(s.ConnInfo.Socket)

	for {
		s.ConnInfo.Socket.SetReadDeadline(time.Now().Add(time.Duration(s.ServerTimeout) * time.Second))
		in, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			s.Disconnect(true)
			return
		}

		s.ConnInfo.TimeLastEvent = time.Now()

		m := IRCMessageParse(in)
		fmt.Println("<<", m.Full)

		switch m.Type {
		case IRCMessagePING:
			s.Send("PONG ", m.Data)
		}

		ev := IRCSessEvent{ IRCSessEventMessage, m }
		go s.EventHandler(s, ev)
	}
}

// vim: set noexpandtab ts=8 sw=8 sts=8:
