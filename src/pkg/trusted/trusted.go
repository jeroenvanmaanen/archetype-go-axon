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

    grpcExample "github.com/jeroenvm/archetype-go-axon/src/pkg/grpc/example"
)

var KeyManagers map[string]string
var TrustedKeys map[string]string
var privateKey rsa.PrivateKey
var privateKeyName string

func SetPrivateKey(name string, pemString string) error {
    var e error

    encodedPublicKey := TrustedKeys[name]
    log.Printf("Trusted: Set private key: public key: %v: %v", name, encodedPublicKey)

    publicKey, e := getTrustedKey(name)
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
    log.Printf("Trusted: Set private key: private key: %v", name)
    return nil
}

func AddTrustedKey(request *grpcExample.TrustedKeyRequest, nonce []byte) error {
    name := request.PublicKey.Name
    publicKey := request.PublicKey.PublicKey
    signatureName := request.SignatureName
    protoSignature := request.Signature
    isKeyManager := request.IsKeyManager

    _, e := parsePublicKey(publicKey)
    if e != nil {
        log.Printf("Trusted: Add trusted key: Unable to parse new trusted key: %v", e)
        return errors.New("Invalid trusted key")
    }

    signatureKey, e := getKeyManagerKey(signatureName)
    if e != nil {
        log.Printf("Trusted: Add trusted key: Unable to parse signature key: %v", e)
        return errors.New("Invalid trusted key")
    }

    signature := ssh.Signature{
        Format: protoSignature.Format,
        Blob: protoSignature.Blob,
        Rest: protoSignature.Rest,
    }

    e = signatureKey.Verify(nonce, &signature)
    if e != nil {
        log.Printf("Trusted: Add trusted key: Unable to verify signature of nonce: %v", e)
        return errors.New("Invalid trusted key")
    }

    TrustedKeys[name] = publicKey
    if isKeyManager {
        KeyManagers[name] = publicKey
    }
    return nil
}

func getKeyManagerKey(name string) (ssh.PublicKey, error) {
    var e error

    encodedPublicKey := KeyManagers[name]
    log.Printf("Trusted: Get key manager key: %v: %v", name, encodedPublicKey)
    publicKey, e := parsePublicKey(encodedPublicKey)
    return publicKey, e
}

func getTrustedKey(name string) (ssh.PublicKey, error) {
    var e error

    encodedPublicKey := TrustedKeys[name]
    log.Printf("Trusted: Get trusted key: %v: %v", name, encodedPublicKey)
    publicKey, e := parsePublicKey(encodedPublicKey)
    return publicKey, e
}

func parsePublicKey(encodedPublicKey string) (ssh.PublicKey, error) {
    publicKey, _, _, _, e := ssh.ParseAuthorizedKey([]byte("ssh-rsa " + encodedPublicKey))
    if e != nil {
        log.Printf("Trusted: Unable to parse public key: %v", e)
        return nil, e
    }
    return publicKey, nil
}