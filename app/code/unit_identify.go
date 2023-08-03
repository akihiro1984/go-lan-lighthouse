package code

const (
    UDP_PACKET_SIZE = 1024
)

type UnitIdentify struct {
    SerialNumber uint64
    IpAddr string
    CallCode byte
    Version string
    Payload []byte
}
