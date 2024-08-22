package emailverification

type EmailRepository interface {
	Save(email, token string) error
	FindByToken(token string) (string, error)
}

type QueueRepository interface {
	Publish(message string) error
}
