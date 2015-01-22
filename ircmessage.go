package main

import (
	"strings"
)

type IRCMessageType int
const (
	_ = iota
	IRCMessagePRIVMSG
	IRCMessageJOIN
	IRCMessagePART
	IRCMessageNOTICE
	IRCMessagePING
	IRCMessageQUIT
	IRCMessageKICK
	IRCMessageRAW
	IRCMessageERROR
	IRCMessageMODE
)

type IRCMessage struct {
	Type	IRCMessageType
	From	string		// Who/where message was sent from
	To	string		// Who/where message was sent to
	Dest	string		// Calculated reply destination
	Data	string
	KNick	string		// Nick being kicked
	Raw	string		// Raw code
	Full	string		// Full, un-modified message
}

func IRCMessageParse(in string) (m IRCMessage) {
	// [:prefix] <command> [params] [:trailing]
	hasPrefix := false
	secIndex := 0

	in = strings.TrimSuffix(in, "\r\n")
	m.Full = in

	if strings.HasPrefix(in, ":") {
		hasPrefix = true
		secIndex++
	}

	sections := strings.Split(in, ":")
	prefixCmdParams := strings.Fields(sections[secIndex]);
	secIndex++
	pcmIndex := 0

	if hasPrefix {
		m.From = prefixCmdParams[pcmIndex];
		pcmIndex++
	}
	m.Type = IRCMessageTypeParse(prefixCmdParams[pcmIndex])
	pcmIndex++

	switch m.Type {
	case IRCMessagePRIVMSG:
		m.To = prefixCmdParams[pcmIndex]
		m.Data = sections[secIndex]
	case IRCMessagePART:
		m.To = prefixCmdParams[pcmIndex]
		m.Data = sections[secIndex]
	case IRCMessageNOTICE:
		m.To = prefixCmdParams[pcmIndex]
		m.Data = sections[secIndex]
	case IRCMessageJOIN:
		m.To = prefixCmdParams[pcmIndex]
	case IRCMessagePING:
		m.Data = sections[secIndex]
	case IRCMessageQUIT:
		m.Data = sections[secIndex]
	case IRCMessageKICK:
		m.To = prefixCmdParams[pcmIndex]
		m.KNick = prefixCmdParams[pcmIndex + 1]
		m.Data = sections[secIndex]
	case IRCMessageRAW:
		m.Raw = prefixCmdParams[pcmIndex - 1]
	case IRCMessageERROR:
		m.Data = sections[secIndex]
	case IRCMessageMODE:
		m.To = prefixCmdParams[pcmIndex]
		m.Data = sections[secIndex]
	}

	return
}

func IRCMessageTypeParse(in string) (t IRCMessageType) {
	switch in {
	case "PRIVMSG":
		return IRCMessagePRIVMSG
	case "JOIN":
		return IRCMessageJOIN
	case "PART":
		return IRCMessagePART
	case "NOTICE":
		return IRCMessageNOTICE
	case "PING":
		return IRCMessagePING
	case "QUIT":
		return IRCMessageQUIT
	case "KICK":
		return IRCMessageKICK
	case "ERROR":
		return IRCMessageERROR
	case "MODE":
		return IRCMessageMODE
	default:
		return IRCMessageRAW
	}
}

// vim: set noexpandtab ts=8 sw=8 sts=8:
