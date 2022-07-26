package main

import (
	"fmt"
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

func getTempString(t float32) string {
	if t >= 0 {
		return fmt.Sprintf("+%.1f°", t)
	} else {
		return fmt.Sprintf("-%.1f°", t)
	}
}

func makeWarningMessage(msg string) string {
	return fmt.Sprintf("⚠️ %s", msg)
}

func makeErrorMessage(msg string) string {
	return fmt.Sprintf("⛔ %s", msg)
}
