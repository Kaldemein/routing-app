package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	mock_api "service/internal/api/mocks"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestServeHTTP(t *testing.T) {
	tokenUrl := fmt.Sprintf("/verify?token=%s", "token")
	req, err := http.NewRequest("GET", tokenUrl, nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	mockCtrl := gomock.NewController(t)
	emailSenderMock := mock_api.NewMockMessageSender(mockCtrl)
	handler := NewHttpHandler(emailSenderMock)

	emailSenderMock.EXPECT().SendVerifyMessageToRMQ("token").Times(1).Return("email", nil)
	handler.ServeHTTP(rr, req)
}
