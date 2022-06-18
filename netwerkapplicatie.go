package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strconv"
	"time"
)

//Alle strings van het programma.
var (
	applicatie string
	ip         string
	poort      int
)

//De basis van de applicatie.
func main() {
	netwerk()
}

//De vragen aan de gebruiker.
func userinput() {
	fmt.Println("Welkom bij mijn netwerkapplicatie!")
	fmt.Println("Wat wilt u doen?")
	fmt.Println("Voor een poortscan, toets 1")
	fmt.Println("Voor een IP-adres ping, toets 2")
	fmt.Println("Voor een IP-adres ping van meerdere IP-adressen, toets 3")
	fmt.Println("Voor een traceroute van een IP-adres, toets 4")
	fmt.Println("Om de netwerkapplicatie af te sluiten, toets 9")
	fmt.Scan(&applicatie)
}

//De applicaties zelf
func netwerk() {
	for {
		userinput()
		switch applicatie {
		//Poortscan applicatie
		case ("1"):
			fmt.Println("Voer een IP-adres in:")
			fmt.Scan(&ip)
			fmt.Println("Voer een poort in")
			fmt.Scan(&poort)
			fmt.Println("Poort wordt gescand...")
			open := scanPort("tcp", ip, poort)
			fmt.Printf("Is de gekozen port open?: %t\n", open)
			//IP-adres applicatie
		case ("2"):
			var IP string
			fmt.Printf("Welk IP-adres wilt u pingen? \n")
			fmt.Scan(&IP)
			Command := fmt.Sprintf("ping -c 1 %v > /dev/null && echo Het pingen van het gekozen IP-adres is gelukt || echo 1 Het pingen van het gekozen IP-adres is mislukt", IP)
			output, err := exec.Command("/bin/sh", "-c", Command).Output()
			fmt.Print(string(output))
			fmt.Print(err)
			//Subnet scannen (NIET AF)
		case ("3"):
			//var IP string
			//fmt.Printf("Welk IP-adres wilt u pingen? \n")
			//fmt.Scan(&IP)
			Command := fmt.Sprintf("ping -c 1 192.168.10.253 > /dev/null && echo Het pingen van IP-adres 1 is gelukt || echo Het pingen van IP-adres 1 is mislukt && ping -c 1 192.168.10.51 > /dev/null && echo Het pingen van IP-adres 2 is gelukt || echo Het pingen van IP-adres 2 is mislukt")
			output, err := exec.Command("/bin/sh", "-c", Command).Output()
			fmt.Print(string(output))
			fmt.Print(err)
			//Traceroute applicatie
		case ("4"):
			var trace string
			fmt.Print("Voor welk IP-adres moet er een traceroute worden uitgevoerd? ")
			fmt.Scan(&trace)
			//Draait de traceroute functie
			RunTraceroute(trace)
			//Sluit de applicatie af
		case ("9"):
			os.Exit(0)
		}
	}
}

//Poortscan
func scanPort(protocol, hostname string, port int) bool {
	address := hostname + ":" + strconv.Itoa(port)
	conn, err := net.DialTimeout(protocol, address, 60*time.Second)
	//Output als de poort dicht is
	if err != nil {
		return false
	}
	//Output als de poort open is
	defer conn.Close()
	return true
}

//Traceroute
func RunTraceroute(host string) {
	errch := make(chan error, 1)
	cmd := exec.Command("traceroute", host)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	go func() {
		errch <- cmd.Wait()
	}()

	go func() {
		for _, char := range "|/-\\" {
			fmt.Printf("\r%s...%c", "Running traceroute", char)
			time.Sleep(100 * time.Millisecond)
		}
		scanner := bufio.NewScanner(stdout)
		fmt.Println("")
		for scanner.Scan() {
			line := scanner.Text()
			log.Println(line)
		}
	}()

	select {
	case <-time.After(time.Second * 100):
		log.Println("Timeout hit..")
		return
	case err := <-errch:
		if err != nil {
			log.Println("traceroute failed:", err)
		}
	}
}
