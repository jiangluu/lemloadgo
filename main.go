package main

import (
	"fmt"
	"net"
	//"bytes"
	"flag"
	//"runtime"
	"strings"
	"time"
)

var (
	s_port string
)

func run_a_c(local_c int, ch chan int) {
	result := -1
	defer func() {
		//fmt.Println("defer", result)
		ch <- result
	}()

	c, err := net.Dial("tcp", s_port)
	if err != nil {
		result = 1
	} else {
		defer c.Close()

		s := fmt.Sprintf("LEM hello 9 CUSTOM\r\n%09d", local_c)

		c.Write([]byte(s))
		result = 2

		buf := make([]byte, 256)
		ll, _ := c.Read(buf)

		if ll > 0 {
			aa := fmt.Sprintf("LEM helloACK 3\r\n%09dACK", local_c)
			buf = buf[:ll]

			if 0 == strings.Compare(aa, string(buf)) {
				result = 0
			} else {
				fmt.Println(string(buf))
				fmt.Println(aa)
			}
		}

	}
}

func run_a_c_2(local_c int, ch chan int) {
	result := -1
	defer func() {
		ch <- result
	}()

	c, err := net.Dial("tcp", s_port)
	if err != nil {
		result = 1
	} else {
		defer c.Close()

		s := fmt.Sprintf("LEM hello 9 CUSTOM\r\n%09d", local_c)
		//should_re := fmt.Sprintf("LEM helloACK 3\r\n%09dACK", local_c)
		count := 0
		buf := make([]byte, 256)
		const N = 100

		for i := 0; i < N; i++ {
			c.Write([]byte(s))
			result = 2

			ll, _ := c.Read(buf)

			if ll > 0 {
				count++
			}
		}
		if count == N {
			result = 0
		}

	}
}

func main() {

	count1 := 0
	count_returned := 0
	count_send := 0
	count1_recv_ok := 0

	concurent := flag.Int("c", 32, "concurent connections")
	times := flag.Int("n", 10000, "number of tasks")
	pp := flag.String("p", ":22222", "host and port")

	flag.Parse()

	s_port = *pp

	// start enough coroutines first
	a_chan := make(chan int, *concurent)

	time1 := time.Now()

	for i := 0; i < *concurent; i++ {
		count1 += 1
		go run_a_c_2(count1, a_chan)
	}

	// if some coroutine returned, make one
	for {
		r1 := <-a_chan
		count_returned++

		if 0 == r1 {
			count1_recv_ok++
			count_send++
		} else if 1 == r1 {
			count_send++
		}

		if 0 == count1%100 {
			fmt.Println(count1, count_send, count1_recv_ok)
			time2 := time.Now()
			fmt.Println("Time used:", time2.Sub(time1).String())
		}

		if count1 < *times {
			count1++
			go run_a_c_2(count1, a_chan)
		} else {
			if count_returned < count1 {
				// time.Sleep(time.Second)
			} else {
				// Over HERE
				time2 := time.Now()
				fmt.Println(count1, count_send, count1_recv_ok)
				fmt.Println("Time used:", time2.Sub(time1).String())
				break
			}
		}
	}

	/*for {
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
	}*/
}
