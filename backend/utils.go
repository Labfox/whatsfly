package main

import "go.mau.fi/whatsmeow/types"

func getJid(user string, is_group bool) types.JID {
	server := types.DefaultUserServer
	if is_group {
		server = types.GroupServer
	}

	return types.JID{
		User:   user,
		Server: server,
	}
}
