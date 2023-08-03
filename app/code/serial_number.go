package code

import (
    "hash/fnv"
    "os/exec"
    "strings"
    "io/ioutil"
    "log"
)

func getRaspberryPi() string {
    out ,err := exec.Command("sh","-c","/bin/cat /proc/cpuinfo | grep Serial | cut -d ' ' -f 2").CombinedOutput()
    if err != nil || len(out) == 0 {
        return ""
    }
    return strings.TrimSuffix(string(out), "\n")
}

func getMacAddress() string {
    out ,err := exec.Command("sh","-c","/bin/cat `find /sys/devices/ -name eth0 2>/dev/null`/address").CombinedOutput()
    log.Printf("%s %v", out, err)
    if err != nil || len(out) == 0 {
        return ""
    }
    return strings.TrimSuffix(string(out), "\n")
}

func getSerialNumberFromHostname() string {
    out ,err := exec.Command("sh","-c","hostname").CombinedOutput()
    if err != nil || len(out) == 0 {
        return ""
    }
    return strings.TrimSuffix(string(out), "\n")
}

func getProductUuid() string {
    out ,err := exec.Command("sh","-c","/bin/cat /sys/class/dmi/id/product_uuid").CombinedOutput()
    if err != nil || len(out) == 0 {
        return ""
    }
    return strings.TrimSuffix(string(out), "\n")
}

func getMachineId() string {
    bytes, err := ioutil.ReadFile("/etc/machine-id")
    if err != nil || len(bytes) == 0 {
        return ""
    }
    return strings.TrimSuffix(string(bytes), "\n")
}

func getProductSerial() string {
    bytes, err := ioutil.ReadFile("/sys/class/dmi/id/product_serial")
    if err != nil || len(bytes) == 0 {
        return ""
    }
    return strings.TrimSuffix(string(bytes), "\n")
}

func getBoardSerial() string {
    bytes, err := ioutil.ReadFile("/sys/class/dmi/id/board_serial")
    if err != nil || len(bytes) == 0 {
        return ""
    }
    return strings.TrimSuffix(string(bytes), "\n")
}

func getMachineHardwareInfo() string {
    out ,err := exec.Command("sh","-c","echo '$(fdisk --list)$(lshw -short)' | sha1sum | cut --delimiter=' ' --fields=1").CombinedOutput()
    if err != nil || len(out) == 0 {
        return ""
    }
    return strings.TrimSuffix(string(out), "\n")
}

func hash(s string) uint64 {
    h := fnv.New64a()
    h.Write([]byte(s))
    return h.Sum64()
}

func GetSerialNumber() uint64 {
    var sid string

    sid = getRaspberryPi()

    if sid == "" {
        sid = getProductUuid()
    }

    if sid == "" {
        sid = getMachineId()
    }

    if sid == "" {
        sid = getMacAddress()
    }

    if sid == "" {
        sid = getProductSerial()
    }

    if sid == "" {
        sid = getBoardSerial()
    }

    if sid == "" {
        sid = getMachineHardwareInfo()
    }

    if sid == "" {
        sid = getSerialNumberFromHostname()
    }

    if sid == "" {
        panic("no ident serial number")
    }

    return hash(sid)
}

