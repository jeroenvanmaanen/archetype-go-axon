package authentication

import (
    bytes "bytes"
    log "log"
    strings "strings"

    base64 "encoding/base64"
    rand "crypto/rand"
    sha256 "crypto/sha256"

    ssh "golang.org/x/crypto/ssh"

    grpcExample "github.com/jeroenvm/archetype-go-axon/src/pkg/grpc/example"
    trusted "github.com/jeroenvm/archetype-go-axon/src/pkg/trusted"
)

var acceptedCredentials map[string]string

func Init() {
    acceptedCredentials = map[string]string{}
}

func SetCredentials(credentials *grpcExample.Credentials) error {
    payload := credentials.Identifier + "=" + credentials.Secret
    signatureKey, e := trusted.GetKeyManagerKey(credentials.Signature.SignatureName)
    if e != nil {
        return e
    }
    signature := ssh.Signature{
        Format: credentials.Signature.Format,
        Blob: credentials.Signature.Blob,
        Rest: credentials.Signature.Rest,
    }
    signatureKey.Verify([]byte(payload), &signature)
    log.Printf("Set credentials: %v: %v", credentials.Identifier, credentials.Secret)
    acceptedCredentials[credentials.Identifier] = credentials.Secret
    return nil
}

func Authenticate(username string, password string) bool {
    hashedPassword := acceptedCredentials[username]
    parts := strings.Split(hashedPassword, ":")
    if len(parts) != 3 {
        return false
    }

    if parts[0] != "sha256" {
        return false
    }

    salt, e := base64.RawStdEncoding.DecodeString(parts[1])
    if e != nil {
        return false
    }

    storedHash, e := base64.RawStdEncoding.DecodeString(parts[2])
    if e != nil {
        return false
    }

    blob := append(salt, ([]byte(password))...)
    givenHash := sha256.Sum256(blob)

    return bytes.Compare(givenHash[:], storedHash) == 0
}

func Encode(password string) string {
    salt := make([]byte, 32)
    rand.Reader.Read(salt)
    passwordBytes := []byte(password)
    blob := append(salt, passwordBytes...)
    hash := sha256.Sum256(blob)
    var builder strings.Builder
    builder.WriteString("sha256:")
    builder.WriteString(base64.RawStdEncoding.EncodeToString(salt))
    builder.WriteString(":")
    builder.WriteString(base64.RawStdEncoding.EncodeToString(hash[:]))
    return builder.String()
}
