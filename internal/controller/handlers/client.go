package handlers

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"os"
	"sync"

	"github.com/bbt-t/lets-go-keep/internal/controller"
	"github.com/bbt-t/lets-go-keep/internal/entity"
	"github.com/bbt-t/lets-go-keep/internal/storage"
	"github.com/bbt-t/lets-go-keep/pkg"

	log "github.com/sirupsen/logrus"
)

// client struct for client handlers.
type client struct {
	conn      ClientConnection
	authToken entity.AuthToken
	masterKey []byte
	*sync.Mutex
}

// newClientHandlers returns new client handlers.
func newClientHandlers(connection ClientConnection) *client {
	return &client{
		conn:  connection,
		Mutex: &sync.Mutex{},
	}
}

// Login logins user by login and password.
func (c *client) Login(credentials entity.UserCredentials) error {
	if credentials.Login == "" || credentials.Password == "" || len(credentials.MasterKey) == 0 {
		return controller.ErrFieldIsEmpty
	}
	authToken, err := c.conn.Login(credentials)
	if err != nil {
		log.Infoln(err)

		return err
	}

	c.Lock()
	defer c.Unlock()

	c.authToken = entity.AuthToken(authToken)
	sha := sha256.New()

	if _, err = sha.Write(c.masterKey); err != nil {
		return storage.ErrUnknown
	}

	key := sha.Sum(nil)
	c.masterKey = key

	return nil
}

// Register creates new user by login and password.
func (c *client) Register(credentials entity.UserCredentials) error {
	if credentials.Login == "" || credentials.Password == "" || len(credentials.MasterKey) == 0 {
		return controller.ErrFieldIsEmpty
	}
	authToken, err := c.conn.Register(credentials)
	if err != nil {
		log.Infoln(err)

		return err
	}

	c.Lock()
	defer c.Unlock()

	c.authToken = entity.AuthToken(authToken)

	sha := sha256.New()
	_, err = sha.Write(c.masterKey)
	if err != nil {
		log.Infoln(err)

		return storage.ErrUnknown
	}

	key := sha.Sum(nil)
	c.masterKey = key

	return nil
}

// GetRecordsInfo gets all records.
func (c *client) GetRecordsInfo() ([]entity.Record, error) {
	c.Lock()
	defer c.Unlock()

	return c.conn.GetRecordsInfo(c.authToken)
}

// GetRecord gets record by recordID and decodes it.
func (c *client) GetRecord(recordID string) (entity.Record, error) {
	c.Lock()
	defer c.Unlock()

	record, errGetRecord := c.conn.GetRecord(c.authToken, recordID)
	if errGetRecord != nil {
		log.Infoln(errGetRecord)

		return record, errGetRecord
	}

	aesBlock, errNewCipher := aes.NewCipher(c.masterKey)
	if errNewCipher != nil {
		log.Infoln(errNewCipher)

		return record, controller.ErrWrongMasterKey
	}

	aesGCM, errNewGCM := cipher.NewGCM(aesBlock)
	if errNewGCM != nil {
		log.Infoln(errNewGCM)

		return record, storage.ErrUnknown
	}

	nonce := record.Data[:aesGCM.NonceSize()]

	decoded, err := aesGCM.Open(nil, nonce, record.Data[aesGCM.NonceSize():], nil)
	if err != nil {
		log.Infoln(err)

		return record, storage.ErrUnknown
	}

	record.Data = decoded

	if record.Type == entity.TypeFile {
		file, err := os.Create(record.Metadata)
		if err != nil {
			log.Infoln(err)

			return record, storage.ErrUnknown
		}

		_, err = file.Write(record.Data)
		if err != nil {
			log.Infoln(err)

			return record, storage.ErrUnknown
		}
		record.Data = []byte("Saved file successfully to " + record.Metadata + ".")
	}

	return record, nil
}

// DeleteRecord deletes record by his ID.
func (c *client) DeleteRecord(recordID string) error {
	c.Lock()
	defer c.Unlock()

	return c.conn.DeleteRecord(c.authToken, recordID)
}

// CreateRecord creates new record.
func (c *client) CreateRecord(record entity.Record) error {
	c.Lock()
	defer c.Unlock()

	aesBlock, errNewCipher := aes.NewCipher(c.masterKey)

	if errNewCipher != nil {
		log.Infoln(errNewCipher)

		return controller.ErrWrongMasterKey
	}

	aesGCM, errNewGCM := cipher.NewGCM(aesBlock)
	if errNewGCM != nil {
		log.Infoln(errNewGCM)

		return storage.ErrUnknown
	}

	nonce, errGenerateRandom := pkg.GenerateRandom(aesGCM.NonceSize())
	if errGenerateRandom != nil {
		log.Infoln(errGenerateRandom)

		return storage.ErrUnknown
	}
	// Encryption:
	out := aesGCM.Seal(nil, nonce, record.Data, nil)

	record.Data = append(nonce, out...)

	return c.conn.CreateRecord(c.authToken, record)
}
