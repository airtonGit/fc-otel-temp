package domain

type CEPByTempClient interface {
	DoRequest(cep string) (string, error)
}
type service struct {
	cepByTempClient CEPByTempClient
}

func NewService(cepByTempClient CEPByTempClient) *service {
	return &service{
		cepByTempClient: cepByTempClient,
	}
}

func (s *service) GetCEP(cep string) (string, error) {
	return s.cepByTempClient.DoRequest(cep)
}
