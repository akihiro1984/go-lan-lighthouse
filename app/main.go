package main

/**
LAN内の子機接続のための灯台
信頼性の高い通信方式へ移行することを想定
信頼性の低い通信データをやりとりする場合はそのままUDP MQTTを利用するに留めるデータとする


port=24444
1, 子機からのブロードキャストを受信する
2, 子機アドレスに向けてUDP接続チャンネルを開く
3, TCP接続に必要な情報を送信する
4, 子機から切断されない限りは保持する、この状態はコンソール上で確認できるデータとする
5, フックできるデータを検出した場合、こちらで記録する
 */

import (
    "lighthouse/tower"
    "lighthouse/ship"
    "lighthouse/code"
    "flag"
    "log"
    "fmt"
)

func main() {
    server := flag.Bool("server", false, "Server Mode")
    viewIdent := flag.Bool("ident", false, "Ident Mode")
    flag.Parse()

    switch {
    case *server:
        tower.Start()
    case *viewIdent:
        serialCode := code.GetSerialNumber()
        log.Printf("Identify Serial Code: %x", serialCode)
        localIp := code.GetLocalIp()
        log.Printf("Local IP: %s", localIp)
    default:
        result := ship.Resolve()
        fmt.Println(fmt.Sprintf("%x %s", result.SerialNumber, result.IpAddr))
    }
}

