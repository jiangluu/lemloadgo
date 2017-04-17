package main

import (
	"fmt"
	"net"
	//"bytes"
	"strings"
	"time"
)

const (
	s_port    = ":22222"
	try_times = -1
	timeout1  = 1
	times     = 99999
)

func main() {

	count1 := 0
	count_send := 0
	count1_recv_ok := 0

	for {
		count1 += 1

		conn, err := net.DialTimeout("tcp", s_port, time.Second*timeout1)
		if err != nil {
			fmt.Println(err)
		} else {
			go func(c net.Conn, local_c int) {
				s := fmt.Sprintf("LEM hello 9 CUSTOM\r\n%09d", local_c)

				c.Write([]byte(s))
				count_send += 1

				buf := make([]byte, 256)
				ll, _ := c.Read(buf)

				c.Close()

				if ll > 0 {

					aa := fmt.Sprintf("LEM helloACK 3\r\n%09dACK", local_c)
					buf = buf[:ll]
					if 0 == strings.Compare(aa, string(buf)) {
						count1_recv_ok += 1
						fmt.Println(local_c, count_send, count1_recv_ok)
					}
				}
			}(conn, count1)
		}

		if count1 >= times {
			time.Sleep(time.Second)
		}
	}
}
