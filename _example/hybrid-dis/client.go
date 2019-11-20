package main

import (
    "os"
    "fmt"
    "bufio"
    "strings"
    "github.com/number571/gopeer"
)

const (
    TITLE_MESSAGE = "[TITLE:MESSAGE]"
    MODE_READ = "[MODE:READ]"
)

func init() {
    gopeer.SettingsSet(gopeer.SettingsType{
        "IS_DISTRIB": true,
        "HAS_CRYPTO": true,
        "HAS_ROUTING": true,
        // "HANDLE_ROUTING": true,
    })
}

func main() {
    node := gopeer.NewNode(gopeer.SettingsGet("CLIENT_NAME").(string)).GeneratePrivate(2048)
    node.Run(handleServer, handleClient)
}

func handleServer(node *gopeer.Node, pack *gopeer.Package) {
    switch pack.Head.Title {
        case TITLE_MESSAGE:
            switch pack.Head.Mode {
                case MODE_READ:
                    message := strings.TrimLeft(pack.Body.Data[0], " ")
                    if message == "" { return }
                    fmt.Printf("[%s]: %s\n", pack.From.Address, message)
            }
    }
}

func handleClient(node *gopeer.Node) {
    node.Connect(":8080")
    for {
        handleCLI(node, strings.Split(inputString(), " "))
    }
}

func inputString() string {
    msg, _ := bufio.NewReader(os.Stdin).ReadString('\n')
    return strings.Replace(msg, "\n", "", -1)
}

func handleCLI(node *gopeer.Node, message []string) {
    switch message[0] {
        case "/exit": os.Exit(0)
        case "/whoami": fmt.Println("|", node.Hashname)
        case "/hidden": node.HiddenConnect(strings.Join(message[1:], " "))
        case "/connect": node.MergeConnect(strings.Join(message[1:], " "))
        case "/network": fmt.Println(node.GetConnections(gopeer.RelationAll))
        case "/send":
            switch len(message[1:]) {
                case 0, 1: fmt.Println("[connect] need > 0, 1 arguments")
                default: node.SendInitRedirect(&gopeer.Package{
                    To: gopeer.To{
                        Address: message[1],
                    },
                    Head: gopeer.Head{
                        Title: TITLE_MESSAGE,
                        Mode: MODE_READ,
                    },
                    Body: gopeer.Body{
                        Data: [gopeer.DATA_SIZE]string{strings.Join(message[2:], " ")},
                    },
                })
            }
        default: node.SendToAll(&gopeer.Package{
            Head: gopeer.Head{
                Title: TITLE_MESSAGE,
                Mode: MODE_READ,
            },
            Body: gopeer.Body{
                Data: [gopeer.DATA_SIZE]string{strings.Join(message, " ")},
            },
        })
    }
}
