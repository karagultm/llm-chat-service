package chat

type Chat struct { //chatdto
	Message   string
	SessionID string
}
type MessageKind string

const (
	UserPrompt MessageKind = "USER_PROMPT"
	LLMOutput  MessageKind = "LLM_OUTPUT"
)

type ChatMessage struct { //direkt chat olmalı adı bence.
	ID        int
	Kind      MessageKind
	Message   string
	Timestamp int64
	SessionID string
}
