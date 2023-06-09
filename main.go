package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"

	"github.com/Juniper/go-netconf/netconf"
	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

// SystemInformation provides a representation of the system-information container
type SystemInformation struct {
	HardwareModel string `xml:"system-information>hardware-model"`
	OsName        string `xml:"system-information>os-name"`
	OsVersion     string `xml:"system-information>os-version"`
	SerialNumber  string `xml:"system-information>serial-number"`
	HostName      string `xml:"system-information>host-name"`
}

var (
	host         = flag.String("host", "192.168.86.3", "Hostname")
	username     = flag.String("username", "", "Username")
	key          = flag.String("key", os.Getenv("HOME")+"/.ssh/id_rsa", "SSH private key file")
	passphrase   = flag.String("passphrase", "", "SSH private key passphrase (cleartext)")
	nopassphrase = flag.Bool("nopassphrase", false, "SSH private key does not contain a passphrase")
	pubkey       = flag.Bool("pubkey", false, "Use SSH public key authentication")
	agent        = flag.Bool("agent", false, "Use SSH agent for public key authentication")
)

func getRPC(devIP string, sshConfig *ssh.ClientConfig, rpcCommand string) *netconf.RPCReply {

	fmt.Printf("Connecting: %s\n", devIP)
	s, err := netconf.DialSSH(devIP, sshConfig)
	if err != nil {
		log.Fatalf("Error connecting: %s", err)
	}

	defer s.Close()

	// Sends raw XML
	res, err2 := s.Exec(netconf.RawMethod(rpcCommand))
	if err2 != nil {
		panic(err2)
	}
	return res
}

func BuildConfig() *ssh.ClientConfig {
	var config *ssh.ClientConfig
	var pass string
	if *pubkey {
		if *agent {
			var err error
			config, err = netconf.SSHConfigPubKeyAgent(*username)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			if *nopassphrase {
				pass = "\n"
			} else {
				if *passphrase != "" {
					pass = *passphrase
				} else {
					var readpass []byte
					var err error
					fmt.Printf("Enter Passphrase for %s: ", *key)
					readpass, err = term.ReadPassword(syscall.Stdin)
					if err != nil {
						log.Fatal(err)
					}
					pass = string(readpass)
					fmt.Println()
				}
			}
			var err error
			config, err = netconf.SSHConfigPubKeyFile(*username, *key, pass)
			if err != nil {
				log.Fatal(err)
			}
		}
	} else {
		fmt.Printf("Enter Password: ")
		bytePassword, err := term.ReadPassword(syscall.Stdin)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println()

		config = netconf.SSHConfigPassword(*username, string(bytePassword))
	}
	return config
}

func main() {
	flag.Parse()

	sshConfig := BuildConfig()

	fmt.Println("\nConnecting....")
	s, err := netconf.DialSSH(*host, sshConfig)
	if err != nil {
		log.Fatal(err)
	}
	defer s.Close()

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\nRPC command to send: ")
	sendCmd, _ := reader.ReadString('\n')
	sendCmd = strings.Trim(sendCmd, "\n")

	fmt.Println(sendCmd)

	fmt.Println("Calling getData function")
	RPCreply := getRPC(*host, sshConfig, sendCmd)
	fmt.Println("Finished getting data.... ")

	fmt.Println(RPCreply)

}
