package helper

import "encoding/base64"

type encrypt struct {
}

type Encrypt interface {
}

func NewEncryption() *encrypt {
	return &encrypt{}
}

func (e *encrypt) EncryptString(value string) string {
	return base64.StdEncoding.EncodeToString([]byte(value))
}

func (e *encrypt) DecodeString(value string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(value)

	if err != nil {
		return "", err
	}

	return string(decoded), nil
}
