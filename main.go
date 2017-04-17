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
	timeout1  = 1
	times     = 100000
	concurent = 32
)

func sum(arrays []int, ch chan int) {
	//fmt.Println(arrays)
	sum := 0
	for _, array := range arrays {
		sum += array
	}
	ch <- sum
}

func main22() {
	arrayChan := make(chan int, 20)
	arrayInt := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	for t := 0; t < 10; t++ {
		length := len(arrayInt)
		go sum(arrayInt[length-t:], arrayChan)
	}

	arrayResult := [10]int{0}
	for i := 0; i < 10; i++ {
		arrayResult[i] = <-arrayChan
	}
	fmt.Println(arrayResult)
}

func run_a_c(local_c int, ch chan int) {
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

		c.Write([]byte(s))
		result = 2

		buf := make([]byte, 256)
		ll, _ := c.Read(buf)

		if ll > 0 {
			aa := fmt.Sprintf("LEM helloACK 3\r\n%09dACK", local_c)
			buf = buf[:ll]
			if 0 == strings.Compare(aa, string(buf)) {
				result = 0
			}
		}
	}
}

func main() {

	count1 := 0
	count_returned := 0
	count_send := 0
	count1_recv_ok := 0

	// start enough coroutines first
	a_chan := make(chan int, concurent)

	for i := 0; i < concurent; i++ {
		count1 += 1
		go run_a_c(count1, a_chan)
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
		}

		if count1 < times {
			count1++
			go run_a_c(count1, a_chan)
		} else {
			if count_returned < count1 {
				time.Sleep(time.Second)
			} else {
				fmt.Println(count1, count_send, count1_recv_ok)
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
