package main

import (
    bufio "bufio"
    fmt "fmt"
    log "log"
    os "os"
    strings "strings"
)

func main() {
    log.Printf("Key Manager -- reading from standard input")
    reader := bufio.NewReader(os.Stdin)

    line := readLine(reader)
    var name string
    var pem string
    for true {
        if strings.HasPrefix(line, ">>> Manager: ") {
            name = getName(line)
            pem, line = readPem(reader)
            log.Printf("Manager name: %v: %d", name, len(pem))
        } else if strings.HasPrefix(line, ">>> Trusted") {
            for true {
                line = readLine(reader)
                if strings.HasPrefix(line, ">>> ") {
                    break;
                }
                log.Printf("Trusted public key: %v", line)
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