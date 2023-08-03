package ship

import (
    "lighthouse/code"
    //
    "encoding/gob"
    "context"
    "bytes"
    "time"
    "log"
    "os"
    // "log"
    "net"
    // "os/exec"
    // "strconv"
    // "strings"
)

/**
 * CONNECTION FOR SEND
 */
func createConnectionForSend() *net.UDPConn {

    addr := "255.255.255.255:" + os.Getenv("RECEIVE_PORT")

    resolvedAddr, err := net.ResolveUDPAddr("udp", addr)
    if err != nil {
        return nil
    }

    conn, err := net.DialUDP("udp", nil, resolvedAddr)
    if err != nil {
        return nil
    }

    return conn
}

/**
 * CONNECTION FOR RECEIVE
 */
 func createConnectionForReceive() *net.UDPConn {

    addr := ":" + os.Getenv("RECEIVE_PORT")

    resolvedAddr, err := net.ResolveUDPAddr("udp", addr)
    if err != nil {
        panic(err)
        return nil
    }

    conn, err := net.ListenUDP("udp", resolvedAddr)
    if err != nil {
        panic(err)
        return nil
    }

    return conn
}

func send(ctx context.Context) {
    serverIdent := code.UnitIdentify{
        SerialNumber: code.GetSerialNumber(),
        IpAddr: code.GetLocalIp(),
        CallCode: 0xFF,
        Version: "0.2a",
    }

    output := bytes.NewBuffer(nil)
    err := gob.NewEncoder(output).Encode(&serverIdent)
    if err != nil {
        panic(err)
    }

    conn := createConnectionForSend()
    defer conn.Close()

    ticker := time.NewTicker(time.Second * 4)
    defer ticker.Stop()

    loop:
    for {
        select {
        case <-ctx.Done():
            break loop
        case <-ticker.C:
            conn.Write(output.Bytes())
        }
    }
}

func receive(done chan code.UnitIdentify) {
    conn := createConnectionForReceive()
    defer conn.Close()

    for {
        conn.SetReadDeadline(time.Now().Add(time.Second * 60))

        buffer := make([]byte, code.UDP_PACKET_SIZE)
        len, addr, err := conn.ReadFrom(buffer)
        if err != nil {
            break
        }

        // 送信元が自分なので落としておく
        if code.CompareSelfIP(addr.(*net.UDPAddr).IP.String()) {
            continue
        }

        clientIdent := code.UnitIdentify{}
        {
            buf := bytes.NewBuffer(buffer[:len])
            err := gob.NewDecoder(buf).Decode(&clientIdent)
            if err != nil {
                log.Println(err)
            }
        }

        done <- clientIdent

        break
    }

    return
}

func Resolve() code.UnitIdentify {
    resultChan := make(chan code.UnitIdentify, 1)

    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    go receive(resultChan)
    go send(ctx)

    return <-resultChan
}
