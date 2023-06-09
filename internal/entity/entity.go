package entity

import (
	"io"
	"os"

	log "github.com/sirupsen/logrus"
)

const (
	TypeLoginAndPassword RecordType = iota
	TypeFile
	TypeText
	TypeCreditCard
)

// UserCredentials struct for user authorization.
type UserCredentials struct {
	Login, Password string
	MasterKey       []byte
}

// UserID is unique identificator of user.
type UserID string

// AuthToken is authorization token of user. Should store userID.
type AuthToken string

// Record is struct for decrypted or encrypted information.
type Record struct {
	ID, Metadata string
	Type         RecordType
	Data         []byte
}

type RecordType int32

func (r RecordType) String() string {
	switch r {
	case TypeLoginAndPassword:
		return "Login + password"
	case TypeFile:
		return "Binary file"
	case TypeText:
		return "Text"
	case TypeCreditCard:
		return "Credit card"
	default:
		return "Unknown"
	}
}

// LoginAndPassword for encrypted login and password.
type LoginAndPassword struct {
	Login, Password string
}

// Bytes implementation of Data interface.
func (data *LoginAndPassword) Bytes() ([]byte, error) {
	return []byte(data.Login + ":" + data.Password), nil
}

// TextData for encrypted text data.
type TextData struct {
	Text string
}

// Bytes gets bytes of information.
func (data *TextData) Bytes() ([]byte, error) {
	return []byte(data.Text), nil
}

// BinaryFile for encrypted file.
type BinaryFile struct {
	FilePath string
	File     *os.File
}

// Bytes gets bytes of information.
func (data *BinaryFile) Bytes() ([]byte, error) {
	file, err := os.Open(data.FilePath)
	if err != nil {
		log.Infoln(err)

		return nil, err
	}
	data.File = file

	return io.ReadAll(data.File)
}

// CreditCard for encrypted credit card.
type CreditCard struct {
	CardNumber, ExpirationDate, CVCCode string
}

// Bytes gets bytes of information.
func (data *CreditCard) Bytes() ([]byte, error) {
	return []byte(data.CardNumber + "|" + data.ExpirationDate + "|" + data.CVCCode), nil
}
