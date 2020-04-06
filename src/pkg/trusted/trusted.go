package trusted

import (
    errors "errors"
    log "log"

    hex "encoding/hex"
    pem "encoding/pem"
    rand "crypto/rand"
    rsa "crypto/rsa"
    x509 "crypto/x509"

    ssh "golang.org/x/crypto/ssh"
)

var TrustedKeys map[string]string
var privateKey rsa.PrivateKey
var privateKeyName string

func SetPrivateKey(name string, pemString string) error {
    var e error

    encodedPublicKey := TrustedKeys[name]
    log.Printf("Trusted: Set private key: public key: %v: %v", name, encodedPublicKey)

    publicKey, _, _, _, e := ssh.ParseAuthorizedKey([]byte("ssh-rsa " + encodedPublicKey))
    if e != nil {
        log.Printf("Trusted: Set private key: Unable to parse public key: %v", e)
        return errors.New("Invalid public key")
    }

    privatePem, _ := pem.Decode([]byte(pemString))
    if privatePem.Type != "RSA PRIVATE KEY" {
        log.Printf("RSA private key is of the wrong type: %v", privatePem.Type)
    }

    privatePemBytes := privatePem.Bytes
    if privatePemBytes == nil {
        log.Printf("Trusted: Set private key: empty PEM")
        return errors.New("Empty PEM")
    }

    var parsedKey interface{}
    if parsedKey, e = x509.ParsePKCS1PrivateKey(privatePemBytes); e != nil {
        if parsedKey, e = x509.ParsePKCS8PrivateKey(privatePemBytes); e != nil { // note this returns type `interface{}`
            log.Printf("Trusted: Set private key: Unable to parse RSA private key: %v", e)
            return errors.New("Invalid private key")
        }
    }

    var givenPrivateKey *rsa.PrivateKey
    var ok bool
    givenPrivateKey, ok = parsedKey.(*rsa.PrivateKey)
    if !ok {
        log.Printf("Type assertion for RSA private key failed")
        return errors.New("Invalid private key")
    }

    signer, e := ssh.NewSignerFromKey(givenPrivateKey)
    if e != nil {
        log.Printf("Trusted: Set private key: Unable to create signer from private key: %v", e)
        return errors.New("Invalid private key")
    }
    nonce := make([]byte, 64)
    rand.Reader.Read(nonce)
    hexNonce := hex.EncodeToString(nonce)
    log.Printf("Nonce: %v", hexNonce)
    signature, e := signer.Sign(rand.Reader, nonce)
    if e != nil {
        log.Printf("Trusted: Set private key: Unable to sign nonce: %v", e)
        return errors.New("Invalid private key")
    }

    e = publicKey.Verify(nonce, signature)
    if e != nil {
        log.Printf("Trusted: Set private key: Unable to verify signature of nonce: %v", e)
        return errors.New("Invalid private key")
    }

    privateKey = *givenPrivateKey
    privateKeyName = name
    log.Printf("Trusted: Set private key: private key: %v: %v", name, privateKey)
    return nil
}
