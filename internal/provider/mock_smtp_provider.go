package provider

type MockSMTPClient struct{}

func (m *MockSMTPClient) SendConfirmationToken(email, token, city string) error {
	return nil
}
