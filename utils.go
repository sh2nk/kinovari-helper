package main

import (
	"log"
	"os"

	"github.com/SevereCloud/vksdk/v2/api/params"
)

// Gets some values from env vars, otherwise returns fallback value
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// Check user preferences
func isConverstationAdmin(peerID int, fromID int) bool {
	b := params.NewMessagesGetConversationMembersBuilder()
	b.PeerID(peerID)

	users, err := VK.MessagesGetConversationMembers(b.Params)
	if err != nil {
		log.Fatal(err)
	}

	for _, user := range users.Items {
		if user.IsAdmin {
			if user.MemberID == fromID {
				return true
			}
		}
	}
	return false
}
