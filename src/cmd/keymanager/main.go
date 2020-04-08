package main

import (
    bufio "bufio"
    context "context"
    fmt "fmt"
    log "log"
    os "os"
    strings "strings"

    rand "crypto/rand"
    rsa "crypto/rsa"

    grpc "google.golang.org/grpc"
    ssh "golang.org/x/crypto/ssh"

    grpcExample "github.com/jeroenvm/archetype-go-axon/src/pkg/grpc/example"
    trusted "github.com/jeroenvm/archetype-go-axon/src/pkg/trusted"
)

func main() {
    log.Printf("Key Manager -- reading from standard input")

    conn, e := grpc.Dial("localhost:8181", grpc.WithInsecure())
    if e != nil {
        log.Printf("Error when dialing application: %v", e)
        panic("Could not connect to application")
    }
    defer conn.Close()

    client := grpcExample.NewGreeterServiceClient(conn)
    log.Printf("Greeter service client: %v", client)

    stream, e := client.ChangeTrustedKeys(context.Background())
    if e != nil {
        log.Printf("Error when opening stream for ChangeTrustedKeys: %v", e)
        panic("Could not call ChangeTrustedKeys")
    }
    log.Printf("ChangeTrustedKeys stream: %v", stream)

    request := grpcExample.TrustedKeyRequest{}
    e = stream.Send(&request)
    if e != nil {
        log.Printf("Error when sending first (empty) request for ChangeTrustedKeys: %v", e)
        panic("Could not send first request")
    }

    response, e := stream.Recv()
    if e != nil {
        log.Printf("Error when receiving first response for ChangeTrustedKeys: %v", e)
        panic("Could not receive first response]")
    }
    log.Printf("First response: %v", response)

    reader := bufio.NewReader(os.Stdin)
    line := readLine(reader)
    var name string
    var pem string
    var managerPrivateKey *rsa.PrivateKey
    var signer ssh.Signer
    for true {
        if strings.HasPrefix(line, ">>> Manager: ") {
            name = getName(line)
            pem, line = readPem(reader)
            log.Printf("Manager name: %v: %d", name, len(pem))
            managerPrivateKey, e = trusted.ParsePrivateKey(pem)
            if e != nil {
                log.Printf("Error when parsing manager private key: %v", e)
                panic("Could not parse manager private key")
            }
            signer, e = ssh.NewSignerFromKey(managerPrivateKey)
            if e != nil {
                log.Printf("Unable to create signer from private key: %v", e)
                panic("Could not create signer")
            }
        } else if strings.HasPrefix(line, ">>> Trusted") {
            for true {
                line = readLine(reader)
                if strings.HasPrefix(line, ">>> ") {
                    break;
                }
                parts := strings.Split(strings.Trim(line, "\n"), " ")
                if len(parts) < 3 {
                    log.Printf("Not enough parts: %v", line)
                    continue
                }
                log.Printf("Trusted public key: %v: %v", parts[2], parts[1])
                publicKey := grpcExample.PublicKey{
                    Name: parts[2],
                    PublicKey: parts[1],
                }
                log.Printf("Public key: %v", publicKey)
                nonce := response.Nonce
                signature, e := signer.Sign(rand.Reader, nonce)
                if e != nil {
                    log.Printf("Unable to sign nonce: %v", e)
                    panic("Could not sign nonce")
                }
                log.Printf("Signature: %v", signature)
                grpcSignature := grpcExample.Signature{}
                grpcSignature.Format = signature.Format
                grpcSignature.Blob = signature.Blob
                grpcSignature.Rest = signature.Rest
                request.PublicKey = &publicKey
                request.Nonce = nonce
                request.SignatureName = name
                request.Signature = &grpcSignature
                request.IsKeyManager = false
                e = stream.Send(&request)
                if e != nil {
                    log.Printf("Error when sending add key request for ChangeTrustedKeys: %v", e)
                    panic("Could not send first request")
                }

                response, e := stream.Recv()
                if e != nil {
                    log.Printf("Error when receiving response for ChangeTrustedKeys: %v", e)
                    panic("Could not receive first response]")
                }
                log.Printf("Response: %v", response)
            }
        } else if strings.HasPrefix(line, ">>> Identity Provider: ") {
            name = getName(line)
            pem, line = readPem(reader)
            log.Printf("Identity Provider name: %v: %d", name, len(pem))
        } else if strings.HasPrefix(line, ">>> End") {
            log.Printf("End")
            break
        } else {
            panic("Expected: >>> { Manager: <name> | Trusted | Identity Provider: <name> | End }: got: [" + line + "]")
        }
    }
}

func readLine(reader *bufio.Reader) string {
    text, e := reader.ReadString('\n')
    if e != nil {
        m, _ := fmt.Printf("Error reading string from stdin: %v", e)
        panic(m)
    }
    return text
}

func getName(line string) string {
    parts := strings.SplitN(line, ":", 2)
    if len(parts) < 2 {
        return "???"
    } else {
        return strings.Trim(parts[1], " \t\r\n")
    }
}

func readPem(reader *bufio.Reader) (pem string, nextLine string) {
    pem = ""
    var builder strings.Builder
    for true {
        nextLine = readLine(reader)
        if strings.HasPrefix(nextLine, ">>> ") {
            break
        }
        builder.WriteString(nextLine)
    }
    pem = builder.String()
    return
}