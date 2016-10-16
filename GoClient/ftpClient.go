package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

var (
	reader    bufio.Reader
	connected bool = false
	text      string
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("ftp>")
	text, _ := reader.ReadString('\n')

	tcpAddr, err := net.ResolveTCPAddr("tcp4", text)
	checkError(err)

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	defer conn.Close()
	checkError(err)

	connectionReader := bufio.NewReader(conn)

	fmt.Printf("ftp> Connected to : %s  on port : %d \n", tcpAddr.IP, tcpAddr.Port)

	for {
		// data := make([]byte, 99999)
		// n, err := conn.Read(data)
		n, _, err := connectionReader.ReadLine()
		checkError(err)

		fmt.Printf("ftp> %s\n", string(n))

		text, _ = reader.ReadString('\n')

		_, err = conn.Write([]byte(text))
		checkError(err)
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stdout, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

// text, _ := reader.ReadString('\n')
