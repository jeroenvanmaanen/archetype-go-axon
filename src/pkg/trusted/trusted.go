package trusted

import (
    errors "errors"
    log "log"

    bigmath "math/big"
    hex "encoding/hex"
    pem "encoding/pem"
    rand "crypto/rand"
    rsa "crypto/rsa"
    x509 "crypto/x509"

    grpc "google.golang.org/grpc"
    jwt "github.com/pascaldekloe/jwt"
    ssh "golang.org/x/crypto/ssh"

    axonserver "github.com/jeroenvm/archetype-go-axon/src/pkg/grpc/axonserver"
    axonutils "github.com/jeroenvm/archetype-go-axon/src/pkg/axonutils"
    grpcExample "github.com/jeroenvm/archetype-go-axon/src/pkg/grpc/example"
)

var keyManagers map[string]string
var trustedKeys map[string]string
var privateKey rsa.PrivateKey
var privateKeyName string

func SetPrivateKey(name string, pemString string) error {
    var e error

    encodedPublicKey := trustedKeys[name]
    log.Printf("Trusted: Set private key: public key: %v: %v", name, encodedPublicKey)

    publicKey, e := getTrustedKey(name)
    if e != nil {
        log.Printf("Trusted: Set private key: Unable to parse public key: %v", e)
        return errors.New("Invalid public key")
    }

    givenPrivateKey, e := ParsePrivateKey(pemString)
    if e != nil {
        return e
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

func AddTrustedKey(request *grpcExample.TrustedKeyRequest, nonce []byte, conn *grpc.ClientConn, clientInfo *axonserver.ClientIdentification) error {
    name := request.PublicKey.Name
    publicKey := request.PublicKey.PublicKey
    protoSignature := request.Signature
    signatureName := protoSignature.SignatureName
    isKeyManager := request.IsKeyManager

    _, e := parsePublicKey(publicKey)
    if e != nil {
        log.Printf("Trusted: Add trusted key: Unable to parse new trusted key: %v", e)
        return errors.New("Invalid trusted key")
    }

    signatureKey, e := GetKeyManagerKey(signatureName)
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

    trustedKeys[name] = publicKey
    if isKeyManager {
        command := grpcExample.RegisterKeyManagerCommand{
            PublicKey: request.PublicKey,
        }
        e = axonutils.SendCommand("RegisterKeyManagerCommand", &command, conn, clientInfo)
        if e != nil {
            log.Printf("Trusted: Error when sending RegisterKeyManagerCommand: %v", e)
        }
    } else {
        command := grpcExample.RegisterTrustedKeyCommand{
            PublicKey: request.PublicKey,
        }
        e = axonutils.SendCommand("RegisterTrustedKeyCommand", &command, conn, clientInfo)
        if e != nil {
            log.Printf("Trusted: Error when sending RegisterTrustedKeyCommand: %v", e)
        }
    }
    log.Printf("Added public key: %v: %v", name, publicKey)
    return nil
}

func UnsafeSetTrustedKey(publicKey *grpcExample.PublicKey) {
    trustedKeys[publicKey.Name] = publicKey.PublicKey
}

func UnsafeSetKeyManager(publicKey *grpcExample.PublicKey) {
    keyManagers[publicKey.Name] = publicKey.PublicKey
}

func GetKeyManagerKey(name string) (ssh.PublicKey, error) {
    var e error

    encodedPublicKey := keyManagers[name]
    log.Printf("Trusted: Get key manager key: %v: %v", name, encodedPublicKey)
    publicKey, e := parsePublicKey(encodedPublicKey)
    return publicKey, e
}

func getTrustedKey(name string) (ssh.PublicKey, error) {
    var e error

    encodedPublicKey := trustedKeys[name]
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

func ParsePrivateKey(pemString string) (givenPrivateKey *rsa.PrivateKey, e error) {
    privatePem, _ := pem.Decode([]byte(pemString))
    if privatePem.Type != "PRIVATE KEY" && privatePem.Type != "RSA PRIVATE KEY" {
        log.Printf("RSA private key is of the wrong type: %v", privatePem.Type)
        return nil, errors.New("Wrong PEM type")
    }

    privatePemBytes := privatePem.Bytes
    if privatePemBytes == nil {
        log.Printf("Trusted: Set private key: empty PEM")
        return nil, errors.New("Empty PEM")
    }

    var parsedKey interface{}
    if parsedKey, e = x509.ParsePKCS1PrivateKey(privatePemBytes); e != nil {
        if parsedKey, e = x509.ParsePKCS8PrivateKey(privatePemBytes); e != nil { // note this returns type `interface{}`
            log.Printf("Trusted: Set private key: Unable to parse RSA private key: %v", e)
            return nil, errors.New("Invalid private key")
        }
    }

    givenPrivateKey, ok := parsedKey.(*rsa.PrivateKey)
    if !ok {
        log.Printf("Type assertion for RSA private key failed")
        return nil, errors.New("Invalid private key")
    }
    return
}

func CreateJWT(claims jwt.Claims) (token string, e error) {
    tokenBuffer, e := claims.RSASign("RS256", &privateKey)
    token = string(tokenBuffer)
    return
}

func GetRsaPublicKey() (*rsa.PublicKey, error) {
    publicKey, e := parsePublicKey(trustedKeys[privateKeyName])
    buffer := publicKey.Marshal()
    w := getInt(buffer)
    keyType := string(buffer[4:w+4])
    if keyType != "ssh-rsa" {
        return nil, errors.New("Not an ssh-rsa key: " + keyType)
    }
    buffer = buffer[w+4:]

    w = getInt(buffer)
    exponent := int(getBigInt(buffer[4:w+4]).Int64())
    buffer = buffer[w+4:]

    w = getInt(buffer)
    modulus := getBigInt(buffer[4:w+4])

    rsaPublicKey := rsa.PublicKey{
        N: modulus,
        E: exponent,
    }

    return &rsaPublicKey, e
}

func getInt(buffer []byte) (result uint32) {
    result = 0
    for i := 0; i < 4; i++ {
        result = result * 256 + uint32(buffer[i])
    }
    return
}

func getBigInt(buffer []byte) (*bigmath.Int) {
    result := bigmath.NewInt(0)
    for _, v := range buffer {
        result.Mul(result, bigmath.NewInt(256))
        result.Add(result, bigmath.NewInt(int64(uint32(v))))
    }
    return result
}

func GetTrustedKeys() (map[string]string) {
    var result = make(map[string]string)
    for name, key := range trustedKeys {
        result[name] = key
    }
    return result
}