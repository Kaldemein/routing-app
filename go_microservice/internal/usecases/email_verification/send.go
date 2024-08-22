package emailverification

import "log"

func (uc *EmailUseCase) SendVerifyMessageToRMQ(token string) (string, error) {
	log.Print("Searching meail by token...")
	email, err := uc.emailRepo.FindByToken(token)
	if err != nil {
		return "", err
	}
	log.Printf("email finded: %v", email)

	// Publish email to RabbitMQ
	err = uc.queue.Publish(email)
	if err != nil {
		return "", err
	}
	log.Printf("email published to queue: %v", email)

	return email, nil
}
