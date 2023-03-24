package redirect

import (
	"fmt"
	"github.com/dlclark/regexp2"
	"github.com/sirupsen/logrus"
	"net"
	"redi301/config"
	"strings"
)

func Process(conn net.Conn) {
	var RequestHead string
	var Path string
	i := 0
	for {
		i++
		var buf [128]byte
		//接受数据
		n, err := conn.Read(buf[:])
		if err != nil {
			logrus.Errorf("[redirect] read from connect failed, err: %v\n", err)
			break
		}
		logrus.Debugf("Receive data: \n%s\n", string(buf[:n]))
		RequestHead += string(buf[:n])
		logrus.Debug(RequestHead)
		if strings.Contains(RequestHead, "HTTP") {
			re := regexp2.MustCompile(`(GET|HEAD?) (.*?) HTTP`, 0)
			if matchArr, err := re.FindStringMatch(RequestHead); err != nil {
				Path = "/"
			} else {
				if matchArr != nil {
					if len(matchArr.Groups()) > 2 {
						Path = matchArr.Groups()[2].String()
					} else {
						Path = "/"
					}
				} else {
					Path = "/"
				}

			}
			if _, err = conn.Write([]byte("HTTP/1.1 301\r\n" +
				"Content-Type: text/html\r\n" +
				"Cache-Control: max-age=86400\r\n" +
				"Content-Length: " + "0" + "\r\n" +
				"Connection: close\r\n" +
				fmt.Sprintf("Mark: %v\r\n", i) +
				fmt.Sprintf("Location: %v%v\r\n\r\n", config.Target, Path))); err != nil {
				logrus.Errorf("write to client failed, err: %v\n", err)
				break
			}
			logrus.Debug("[redirect] Response sent...")
			// force to close the connection
			//if err := conn.(*net.TCPConn).SetLinger(0); err != nil {
			//	logrus.Debugf("failed to set linger, err: %v\n", err)
			//}
			//break
		}
	}

	err := conn.Close()
	if err == nil {
		logrus.Debug("succeed to close")
	}
}
