package signatureutil

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"time"
)

func GenerateToken(accountId string) string {
	return EncryptWithSalt(fmt.Sprintf("%s_%s_%d", accountId, uuid.NewV4().String(), time.Now().UnixNano()))
}
