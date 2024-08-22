package emailverification

import (
	"crypto/rand"
	"encoding/hex"
	"log"
)

func (uc *EmailUseCase) GenerateAndSendLink(email string) error {
	token := generateSecureToken(20)
	log.Printf("Generated token: %v", token)
	err := uc.emailRepo.Save(email, token)
	if err != nil {
		return err
	}

	err = uc.emailSender(email, token)
	if err != nil {
		return err
	}
	log.Printf("email sended")

	return nil
}

func generateSecureToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}
