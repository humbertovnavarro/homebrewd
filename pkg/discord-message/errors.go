package discordmessage

type DiscordMessageSendError = string

const (
	EmptyErrorMessage        DiscordMessageSendError = `HTTP 400 Bad Request, {"message": "Cannot send an empty message", "code": 50006}`
	UnauthorizedErrorMessage DiscordMessageSendError = `HTTP 401 Unauthorized, {"message": "401: Unauthorized", "code": 0}`
)
