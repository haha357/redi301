package lagran

import (
	"context"
	"github.com/coreos/go-iptables/iptables"
	"github.com/florianl/go-nfqueue"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/panjf2000/ants/v2"
	"github.com/sirupsen/logrus"
	"redi301/config"
	"regexp"
	"strings"
	"sync"
	"time"
)

const (
	EnableSynAck       = true
	EnableAck          = true
	EnablePshAck       = true
	EnableFinAck       = true
	WindowSizeOfSynAck = 15
	WindowSizeOfAck    = 5
	WindowSizeOfPshAck = 5
	WindowSizeOfFinAck = 5
)

func Run() {
	setIptable(config.HttpPort)
	var wg sync.WaitGroup
	if EnableSynAck {
		p1, _ := ants.NewPoolWithFunc(128, func(i interface{}) {
			packetHandle(i.(int))
			wg.Done()
		})
		defer p1.Release()
		logrus.Debug("[lagran service] Starting Task p1")
		for i := 1000; i < 1128; i++ {
			wg.Add(1)
			_ = p1.Invoke(int(i))
		}
	}
	if EnableAck {
		p2, _ := ants.NewPoolWithFunc(128, func(i interface{}) {
			packetHandle(i.(int))
			wg.Done()
		})
		defer p2.Release()
		logrus.Debug("[lagran service] Starting Task p2")
		for i := 2000; i < 2128; i++ {
			wg.Add(1)
			_ = p2.Invoke(int(i))
		}
	}
	if EnablePshAck {
		p3, _ := ants.NewPoolWithFunc(128, func(i interface{}) {
			packetHandle(i.(int))
			wg.Done()
		})
		defer p3.Release()
		logrus.Debug("[lagran service] Starting Task p3")
		for i := 3000; i < 3128; i++ {
			wg.Add(1)
			_ = p3.Invoke(int(i))
		}
	}
	if EnableFinAck {
		p4, _ := ants.NewPoolWithFunc(128, func(i interface{}) {
			packetHandle(i.(int))
			wg.Done()
		})
		defer p4.Release()
		logrus.Debug("[lagran service] Starting Task p4")
		for i := 4000; i < 4128; i++ {
			wg.Add(1)
			_ = p4.Invoke(int(i))
		}
	}
}

func setIptable(sport string) {
	ipt, err := iptables.New()
	if err != nil {
		logrus.Fatalf("[lagran service] Iptabels new error:%v\n", err)
	}

	logrus.Debug("Starting to set iptables OUTPUT chain.")
	//if err := ipt.ClearChain("filter", "OUTPUT"); err != nil {
	//	color.Redln("[lagran service] Failed to clear iptables OUTPUT chain.")
	//}
	if EnableSynAck {
		_ = ipt.AppendUnique("filter", "OUTPUT", "-p", "tcp", "-m", "multiport", "--sport", sport, "--tcp-flags", "SYN,RST,ACK,FIN,PSH", "SYN,ACK", "-j", "NFQUEUE", "--queue-balance", "1000:1127")
	}
	if EnableAck {
		_ = ipt.AppendUnique("filter", "OUTPUT", "-p", "tcp", "-m", "multiport", "--sport", sport, "--tcp-flags", "SYN,RST,ACK,FIN,PSH", "ACK", "-j", "NFQUEUE", "--queue-balance", "2000:2127")
	}
	if EnablePshAck {
		_ = ipt.AppendUnique("filter", "OUTPUT", "-p", "tcp", "-m", "multiport", "--sport", sport, "--tcp-flags", "SYN,RST,ACK,FIN,PSH", "PSH,ACK", "-j", "NFQUEUE", "--queue-balance", "3000:3127")
	}
	if EnableFinAck {
		_ = ipt.AppendUnique("filter", "OUTPUT", "-p", "tcp", "-m", "multiport", "--sport", sport, "--tcp-flags", "SYN,RST,ACK,FIN,PSH", "FIN,ACK", "-j", "NFQUEUE", "--queue-balance", "4000:4127")
	}
}
func UnsetIptable(sport string) {
	ipt, err := iptables.New()
	if err != nil {
		logrus.Fatalf("[lagran service] Iptabels new error:%v", err)
	}
	if EnableSynAck {
		_ = ipt.Delete("filter", "OUTPUT", "-p", "tcp", "-m", "multiport", "--sport", sport, "--tcp-flags", "SYN,RST,ACK,FIN,PSH", "SYN,ACK", "-j", "NFQUEUE", "--queue-balance", "1000:1127")
	}
	if EnableAck {
		_ = ipt.Delete("filter", "OUTPUT", "-p", "tcp", "-m", "multiport", "--sport", sport, "--tcp-flags", "SYN,RST,ACK,FIN,PSH", "ACK", "-j", "NFQUEUE", "--queue-balance", "2000:2127")
	}
	if EnablePshAck {
		_ = ipt.Delete("filter", "OUTPUT", "-p", "tcp", "-m", "multiport", "--sport", sport, "--tcp-flags", "SYN,RST,ACK,FIN,PSH", "PSH,ACK", "-j", "NFQUEUE", "--queue-balance", "3000:3127")
	}
	if EnableFinAck {
		_ = ipt.Delete("filter", "OUTPUT", "-p", "tcp", "-m", "multiport", "--sport", sport, "--tcp-flags", "SYN,RST,ACK,FIN,PSH", "FIN,ACK", "-j", "NFQUEUE", "--queue-balance", "4000:4127")
	}
}
func packetHandle(queueNum int) {
	nfqconfig := nfqueue.Config{
		NfQueue:      uint16(queueNum),
		MaxPacketLen: 0xFFFF,
		MaxQueueLen:  0xFF,
		Copymode:     nfqueue.NfQnlCopyPacket,
		WriteTimeout: 15 * time.Millisecond,
	}

	nf, err := nfqueue.Open(&nfqconfig)
	if err != nil {
		logrus.Fatalf("[lagran] could not open nfqueue socket:", err)
	}

	defer nf.Close()

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	fn := func(a nfqueue.Attribute) int {
		id := *a.PacketID
		packet := gopacket.NewPacket(*a.Payload, layers.LayerTypeIPv4, gopacket.Default)

		if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
			tcp, _ := tcpLayer.(*layers.TCP)
			ports := strings.Split(config.HttpPort, ",")
			reg := regexp.MustCompile(`\d+`)
			sport := reg.FindString(tcp.SrcPort.String())
			var matchedPort = false
			for _, port := range ports {
				if port == sport {
					matchedPort = true
					break
				}
			}
			if matchedPort {
				var ok1 = EnableSynAck && tcp.SYN && tcp.ACK
				var ok2 = EnableAck && tcp.ACK && !tcp.PSH && !tcp.FIN && !tcp.SYN && !tcp.RST
				var ok3 = EnablePshAck && tcp.PSH && tcp.ACK
				var ok4 = EnableFinAck && tcp.FIN && tcp.ACK
				var windowSize uint16
				if ok1 || ok2 || ok3 || ok4 {
					if ok1 {
						windowSize = WindowSizeOfSynAck
						logrus.Debug("[lagran] Handle SYN=1 and ACK=1")
					}
					if ok2 {
						windowSize = WindowSizeOfAck
						logrus.Debug("[lagran] Handle ACK=1")
					}
					if ok3 {
						windowSize = WindowSizeOfPshAck
						logrus.Debug("[lagran] Handle PSH=1 and ACK=1")
					}
					if ok4 {
						windowSize = WindowSizeOfFinAck
						logrus.Debug("[lagran] Handle FIN=1 and ACK=1")
					}
					packet.TransportLayer().(*layers.TCP).Window = windowSize
					err := packet.TransportLayer().(*layers.TCP).SetNetworkLayerForChecksum(packet.NetworkLayer())
					if err != nil {
						logrus.Errorf("[lagran] SetNetworkLayerForChecksum error: %v\n", err)
					}
					buffer := gopacket.NewSerializeBuffer()
					options := gopacket.SerializeOptions{
						ComputeChecksums: true,
						FixLengths:       true,
					}
					if err := gopacket.SerializePacket(buffer, options, packet); err != nil {
						logrus.Errorf("[lagran] SerializePacket error: %v\n", err)
					}
					packetBytes := buffer.Bytes()
					logrus.Debugf("[lagran] Set TCP window size to %d\n", windowSize)
					err = nf.SetVerdictModPacket(id, nfqueue.NfAccept, packetBytes)
					if err != nil {
						logrus.Errorf("[lagran] SetVerdictModified error: %v\n", err)
					}
					return 0
				}
				err := nf.SetVerdict(id, nfqueue.NfAccept)
				if err != nil {
					logrus.Errorf("[lagran] SetVerdictModified error: %v\n", err)
				}
				return 0
			}
			err := nf.SetVerdict(id, nfqueue.NfAccept)
			if err != nil {
				logrus.Errorf("[lagran] SetVerdictModified error: %v\n", err)
			}
			return 0
		}
		err := nf.SetVerdict(id, nfqueue.NfAccept)
		if err != nil {
			logrus.Errorf("[lagran] SetVerdictModified error: %v\n", err)
		}
		return 0
	}

	// Register your function to listen on nflqueue queue 100
	err = nf.RegisterWithErrorFunc(ctx, fn, func(e error) int {
		if e != nil {
			logrus.Errorf("[lagran] RegisterWithErrorFunc Error:%v\n", e)
		}
		return 0
	})
	if err != nil {
		logrus.Fatalf("[lagran] error: %v\n", err)
	}
	<-ctx.Done()
}
