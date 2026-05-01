package main

import (
	"context"
	"encoding/json"
)

func (w *WhatsAppClient) GetGroupInviteLink(group string, reset bool, returnid string) int {
	numberObj := getJid(group, true)

	if !w.wpClient.IsConnected() {
		err := w.wpClient.Connect()
		if err != nil {
			return 1
		}
	}

	link, err := w.wpClient.GetGroupInviteLink(context.Background(), numberObj, reset)
	w.addEventToQueue(OutGoingEvent{
		Type:         "groupInviteLink",
		JIDConcerned: w.wpClient.Store.GetJID().ADString(),
		Content: struct {
			Group string `json:"group"`
			Link  string `json:"link"`
		}{
			Group: group,
			Link:  link,
		},
	})
	w.addEventToQueue(OutGoingEvent{
		Type:         "methodReturn",
		JIDConcerned: w.wpClient.Store.GetJID().ADString(),
		Content: ReturnFunc{
			ID:      returnid,
			Content: link,
		},
	})
	if err != nil {
		return 1
	}
	return 0
}

func (w *WhatsAppClient) JoinGroupWithInviteLink(link string) int {
	if !w.wpClient.IsConnected() {
		err := w.wpClient.Connect()
		if err != nil {
			return 1
		}
	}

	_, err := w.wpClient.JoinGroupWithLink(context.Background(), link)
	if err != nil {
		return 1
	}
	return 0
}

func (w *WhatsAppClient) SetGroupAnnounce(group string, announce bool) int {
	numberObj := getJid(group, true)

	if !w.wpClient.IsConnected() {
		err := w.wpClient.Connect()
		if err != nil {
			return 1
		}
	}

	err := w.wpClient.SetGroupAnnounce(context.Background(), numberObj, announce)
	if err != nil {
		return 1
	}
	return 0
}

func (w *WhatsAppClient) SetGroupLocked(group string, locked bool) int {
	numberObj := getJid(group, true)

	if !w.wpClient.IsConnected() {
		err := w.wpClient.Connect()
		if err != nil {
			return 1
		}
	}

	err := w.wpClient.SetGroupLocked(context.Background(), numberObj, locked)
	if err != nil {
		return 1
	}
	return 0
}

func (w *WhatsAppClient) SetGroupName(group string, name string) int {
	numberObj := getJid(group, true)

	if !w.wpClient.IsConnected() {
		err := w.wpClient.Connect()
		if err != nil {
			return 1
		}
	}

	err := w.wpClient.SetGroupName(context.Background(), numberObj, name)
	if err != nil {
		return 1
	}
	return 0
}

func (w *WhatsAppClient) SetGroupTopic(group string, topic string) int {
	numberObj := getJid(group, true)

	if !w.wpClient.IsConnected() {
		err := w.wpClient.Connect()
		if err != nil {
			return 1
		}
	}

	err := w.wpClient.SetGroupTopic(context.Background(), numberObj, "", "", topic)
	if err != nil {
		return 1
	}
	return 0
}

func (w *WhatsAppClient) GetGroupInfo(group string, return_id string) int {
	numberObj := getJid(group, true)

	if !w.wpClient.IsConnected() {
		err := w.wpClient.Connect()
		if err != nil {
			return 1
		}
	}

	groupinfo, err := w.wpClient.GetGroupInfo(context.Background(), numberObj)
	if err != nil {
		return 1
	}

	b, err := json.Marshal(&groupinfo)
	if err != nil {
		return 1
	}

	w.addEventToQueue(OutGoingEvent{
		Type:         "methodReturn",
		JIDConcerned: w.wpClient.Store.GetJID().ADString(),
		Content: ReturnFunc{
			ID:      return_id,
			Content: b,
		},
	})

	return 0
}
