package api

type MessageSender interface {
	SendVerifyMessageToRMQ(string) (string, error)
}
