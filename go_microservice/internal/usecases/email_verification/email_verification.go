// saving email token into db
// find row by token and return email
// send token to email
// send email to queue
package emailverification

type EmailSender func(to string, token string) error
type EmailUseCase struct {
	emailRepo   EmailRepository
	queue       QueueRepository
	emailSender EmailSender
}

func NewEmailUseCase(emailRepo EmailRepository, queue QueueRepository, emailSender EmailSender) *EmailUseCase {
	return &EmailUseCase{emailRepo, queue, emailSender}
}
