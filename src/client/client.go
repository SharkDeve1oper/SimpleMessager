package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

// прием данных и ввести на печать
func readSock(conn net.Conn) {

	if conn == nil {
		panic("Connections is nil")
	}

	buf := make([]byte, 256)
	eof_count := 0
	for {
		// очистка буфера
		for i := 0; i < 256; i++ {
			buf[i] = 0
		}

		readed_len, err := conn.Read(buf)
		if err != nil {
			if err.Error() == "EOF" {
				eof_count++
				time.Sleep(time.Second * 4)
				fmt.Println("EOF")
				if eof_count > 5 {
					fmt.Println("Timeout connection")
					break
				}

				continue
			}

			if strings.Index(err.Error(), "use of closed network connection") > 0 {
				fmt.Println("Connection not exists or closed")
				continue
			}

			panic(err.Error())
		}

		if readed_len > 0 {
			fmt.Println(string(buf))
		}
	}
}

func readConsole(ch chan string) {
	for {
		line, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		if len(line) > 200 {
			fmt.Println("Error: message is very large")
			continue
		}
		fmt.Print(">")
		out := line[:len(line)-1] // Убираем каретку
		ch <- out                 // Отправляем данные в канал
	}
}

func main() {

	ch := make(chan string)

	defer close(ch)

	// Подключаемся к сокету
	conn, _ := net.Dial("tcp", "127.0.0.2:8080")

	if conn == nil {
		panic("Connection is nil")
	}

	go readConsole(ch)
	go readSock(conn)

	for {
		val, ok := <-ch

		if ok {
			out := []byte(val)
			_, err := conn.Write(out)
			if err != nil {
				fmt.Println("Write error:", err.Error())
				break
			}
		} else {
			// данных нет, устроим задержку
			time.Sleep(time.Second * 2)
		}
	}

	fmt.Println("Finished...")

	conn.Close()
}
