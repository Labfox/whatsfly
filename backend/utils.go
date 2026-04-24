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

func getNewsletterJid(user string) types.JID {
	return types.JID{
		User:   user,
		Server: types.NewsletterServer,
	}
}
