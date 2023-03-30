package edgegpt

import (
	"fmt"
)

type ConversationStyle string

const (
	Creative ConversationStyle = "h3relaxedimg"
	Balanced ConversationStyle = "galileo"
	Precise  ConversationStyle = "h3precise"
)

// Request object for ChatHub
type ChatHubRequest struct {
	Struct                map[string]interface{}
	ConversationSignature string
	ClientId              string
	ConversationId        string
	InvocationId          int
}

func NewChatHubRequest(conversationSignature, clientId, conversationId string, invocationId int) *ChatHubRequest {
	req := &ChatHubRequest{
		Struct:                map[string]interface{}{},
		ConversationSignature: conversationSignature,
		ConversationId:        conversationId,
		ClientId:              clientId,
		InvocationId:          invocationId,
	}
	return req
}

// Updates request object
func (req *ChatHubRequest) Update(prompt string, conversation_style ConversationStyle, options ...string) {
	if len(options) == 0 {
		options = []string{
			"deepleo",
			"enable_debug_commands",
			"disable_emoji_spoken_text",
			"enablemm",
		}
	}
	if conversation_style != "" {
		options = []string{
			"nlu_direct_response_filter",
			"deepleo",
			"disable_emoji_spoken_text",
			"responsible_ai_policy_235",
			"enablemm",
			string(conversation_style),
			"dtappid",
			"cricinfo",
			"cricinfov2",
			"dv3sugg",
		}
	}
	req.Struct = map[string]interface{}{
		"arguments": []interface{}{
			map[string]interface{}{
				"source":      "cib",
				"optionsSets": options,
				"sliceIds": []interface{}{
					"222dtappid",
					"225cricinfo",
					"224locals0",
				},
				"traceId":          GetRandomHex(32),
				"isStartOfSession": req.InvocationId == 0,
				"message": map[string]interface{}{
					"author":      "user",
					"inputMethod": "Keyboard",
					"text":        prompt,
					"messageType": "Chat",
				},
				"conversationSignature": req.ConversationSignature,
				"participant": map[string]interface{}{
					"id": req.ClientId,
				},
				"conversationId": req.ConversationId,
			},
		},
		"invocationId": fmt.Sprintf("%d", req.InvocationId),
		"target":       "chat",
		"type":         4,
	}
	req.InvocationId += 1
}
