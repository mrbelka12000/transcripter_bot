package transcriber

import "errors"

type TranscriberService interface {
	TranscribeAndSave(string, int64) error
}

// TODO add services that will be used
type ServiceImpl struct {
}

// my convertation and adding to DB
func (s *ServiceImpl) TranscribeAndSave(str string, num int64) error {
	return errors.New("Method not implemated")
}
