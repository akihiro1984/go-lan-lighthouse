package tower

import (
    "lighthouse/code"
    //
    "golang.org/x/sys/unix"
    "encoding/gob"
    "context"
    "syscall"
    "bytes"
    "time"
    "log"
    "net"
    "os"
)

const (
    UDP_PACKET_SIZE = 2048
)

type ControlCommand struct {
    Type string
    Mode string
    Value string
}

func createServiceSocket() net.PacketConn {
    listenConfig := &net.ListenConfig{
        Control: func(network, address string, c syscall.RawConn) (err error) {
            return c.Control(func(fd uintptr) {
                iFd := int(fd)
                syscall.SetsockoptInt(iFd, unix.SOL_SOCKET, unix.SO_REUSEADDR, 1)
                syscall.SetsockoptInt(iFd, unix.SOL_SOCKET, unix.SO_REUSEPORT, 1)
            })
        },
    }

    addr := ":" + os.Getenv("RECEIVE_PORT")

    conn, err := listenConfig.ListenPacket(context.Background(), "udp", addr)
    if err != nil {
        log.Fatal(err)
    }

    return conn
}

func getServerIdent() []byte {

    serverIdent := code.UnitIdentify{
        SerialNumber: code.GetSerialNumber(),
        IpAddr: code.GetLocalIp(),
        CallCode: 0xFF,
        Version: "0.1a",
    }

    output := bytes.NewBuffer(nil)

    err := gob.NewEncoder(output).Encode(&serverIdent)
    if err != nil {
        panic(err)
    }

    return output.Bytes()
}

func attemptConnectClient(remoteAddr net.Addr) {

    addr, _ := net.ResolveUDPAddr(remoteAddr.Network(), remoteAddr.String())

    addr.Port = 24444

    dialer := net.Dialer{
        Control: func(network, address string, c syscall.RawConn) (err error) {
            return c.Control(func(fd uintptr) {
                iFd := int(fd)
                syscall.SetsockoptInt(iFd, unix.SOL_SOCKET, unix.SO_REUSEADDR, 1)
                syscall.SetsockoptInt(iFd, unix.SOL_SOCKET, unix.SO_REUSEPORT, 1)
            })
        },
    }
    conn, err := dialer.DialContext(context.Background(), addr.Network(), addr.String())
    if err != nil {
        log.Print(err)
        return
    }
    defer conn.Close()

    // SERVER SIGNATURE
    a,e := conn.Write(getServerIdent())
    if e != nil {
        log.Print(e)
        return
    }
    log.Printf(">>>> : %v %v %v", addr, a, e)

    buf := make([]byte, UDP_PACKET_SIZE)
    for {
        // timeout in 30 seconds
        conn.SetDeadline(time.Now().Add(30 * time.Second))

        readLen, err1 := conn.Read(buf[:])
        if err1 != nil {
            if err2, ok := err1.(net.Error); ok && err2.Timeout() {
                break
            }
            log.Print(err1)
            break
        }

        if readLen == 0 {
            log.Println("close connection", conn.RemoteAddr())
            break
        }

        log.Printf(">>>> : %v %v %v", addr, readLen, buf[:readLen])
    }
}

func Start() {
    conn := createServiceSocket()

    var ident code.UnitIdentify
    buf := make([]byte, UDP_PACKET_SIZE)
    for {
        recvLength, remoteAddr, err := conn.ReadFrom(buf[:])
        if err != nil {
            log.Print(err)
            break
        }

        if err != nil || recvLength == 0 {
            // ignore
            continue
        }

        // @todo: 暗号化するか？
        buf := bytes.NewBuffer(buf)
        err = gob.NewDecoder(buf).Decode(&ident)
        if err != nil {
            // extra data in buffer
            continue
        }

        log.Printf("client joined : %v %v", recvLength, remoteAddr)

        go attemptConnectClient(remoteAddr)
    }
}
