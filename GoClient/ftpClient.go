package main

import (
	"bufio"
	"fmt"
	//"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"
)

var (
	reader      bufio.Reader
	connected   bool = false
	text        string
	pasv        bool = false
	result      string
	pasvPort    int
	dataReader  *bufio.Reader
	commandConn *net.TCPConn
	dataConn    *net.TCPConn
	dataResult  []byte
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("ftp>")
	text, _ := reader.ReadString('\n')

	tcpAddr, err := net.ResolveTCPAddr("tcp4", text)
	checkError(err)

	commandConn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)

	connectionReader := bufio.NewReader(commandConn)

	fmt.Printf("ftp> Connected to : %s  on port : %d \n", tcpAddr.IP, tcpAddr.Port)
	result, _, err := connectionReader.ReadLine()
	checkError(err)
	fmt.Printf("ftp> %s \nftp>", string(result))

	for {

		text, _ = reader.ReadString('\n')

		switch text[:4] {
		case "FEAT":
			_, err = commandConn.Write([]byte(text))
			checkError(err)
			byteData := []byte{}

			for {
				currentByte, _ := connectionReader.ReadByte()

				byteData = append(byteData, currentByte)
				if len(byteData) > 3 {
					if string(byteData[len(byteData)-3:]) == "End" {
						break
					}
				}

			}

			fmt.Printf("ftp> %s \nftp>", string(byteData))

			break

		case "PWD":
			_, err = commandConn.Write([]byte(text))
			checkError(err)

			result, _, err = connectionReader.ReadLine()
			checkError(err)
		case "PASV":

			pasv = true
			_, err = commandConn.Write([]byte(text))
			checkError(err)

			result, _, err = connectionReader.ReadLine()
			checkError(err)
			fmt.Printf("ftp> %s \nftp>", string(result))
			passiveIpAndPort := strings.Split(string(result[27:len(result)-2]), ",")
			octet5, _ := strconv.Atoi(passiveIpAndPort[4])
			octet6, _ := strconv.Atoi(passiveIpAndPort[5])
			pasvPort = (octet5 * 256) + octet6
			pasvPortstr := strconv.Itoa(pasvPort)

			tcpAddrData, err := net.ResolveTCPAddr("tcp4", strings.Replace(string(result[27:len(result)-9]), ",", ".", -1)+":"+pasvPortstr)
			checkError(err)

			dataConn, err := net.DialTCP("tcp", nil, tcpAddrData)
			checkError(err)

			dataReader = bufio.NewReader(dataConn)
			checkError(err)
			defer dataConn.Close()
			break

		case "LIST":
			if !pasv {
				_, err = commandConn.Write([]byte(text))
				checkError(err)
				result, _, err = connectionReader.ReadLine()
				fmt.Printf("ftp> %s \nftp>", string(result))
				break
			}
			_, err = commandConn.Write([]byte(text))
			checkError(err)

			result, _, err = connectionReader.ReadLine()
			checkError(err)
			fmt.Println("1")
			fmt.Printf("ftp> %s \n", string(result))
			dataResult = []byte{}

			for {
				currentByte, err := dataReader.ReadByte()
				dataResult = append(dataResult, currentByte)
				if err != nil {
					break
				}
			}

			//			dataResult, _, err = dataReader.ReadLine()

			fmt.Printf("\n%s \n", string(dataResult))

			result, _, err = connectionReader.ReadLine()

			fmt.Printf("ftp> %s \nftp>", string(result))
			// fullList, err := ioutil.ReadAll(connectionReader)
			// result, _, err = connectionReader.ReadLine()
			// checkError(err)

			//fmt.Printf("ftp> %s \nftp>", string(fullList))
			pasv = false

		default:
			_, err = commandConn.Write([]byte(text))
			checkError(err)

			result, _, err = connectionReader.ReadLine()
			checkError(err)
			fmt.Printf("ftp> %s \nftp>", string(result))
		}

		// if dataConnected && awaitingData {
		// 	dataResult, _, err = dataReader.ReadLine()
		// 	checkError(err)
		// 	awaitingData = false
		// }
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stdout, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

// text, _ := reader.ReadString('\n')
