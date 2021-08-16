package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"

	"github.com/Juniper/go-netconf/netconf"
	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
	// "golang.org/x/crypto/ssh/terminal"
)

func credentials() (string, string) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter Username: ")
	username, _ := reader.ReadString('\n')

	fmt.Print("Enter Password: ")
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		fmt.Println("\nPassword typed: " + string(bytePassword))
	}
	password := string(bytePassword)

	return strings.TrimSpace(username), strings.TrimSpace(password)
}

func getRPC(devIP string, devUser string, devPass string, rpcCommand string) *netconf.RPCReply {
	sshConfig := &ssh.ClientConfig{
		User:            devUser,
		Auth:            []ssh.AuthMethod{ssh.Password(devPass)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         60,
	}
	fmt.Printf("Connecting: %s\n", devIP)
	s, err := netconf.DialSSH(devIP, sshConfig)
	if err != nil {
		log.Fatalf("Error connecting: %s", err)
	}

	defer s.Close()

	// fmt.Println(s.ServerCapabilities)
	// fmt.Println(s.SessionID)

	// Sends raw XML
	res, err2 := s.Exec(netconf.RawMethod(rpcCommand))
	if err2 != nil {
		panic(err2)
	}
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Device IP/Hostname: ")
	devAddr, _ := reader.ReadString('\n')
	devAddr = strings.Trim(devAddr, "\n")

	User, Passwd := credentials()

	reader = bufio.NewReader(os.Stdin)
	fmt.Print("\nRPC command to send: ")
	sendCmd, _ := reader.ReadString('\n')
	sendCmd = strings.Trim(sendCmd, "\n")

	fmt.Println("Calling getData function")
	RPCreply := getRPC(devAddr, User, Passwd, sendCmd)
	fmt.Println("Finished getting data.... ")

	fmt.Println(RPCreply)

}
