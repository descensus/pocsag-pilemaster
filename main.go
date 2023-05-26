package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	server := client.OptionsReader()
	log.Println("[POCSAG-Pilemaster] Connected to:", server.Servers()[0])
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	server := client.OptionsReader()
	log.Printf("[POCSAG-Pilemaster] Connect lost to %s: %v", server.Servers()[0], err)
}

func publish(client mqtt.Client, msg string, topic string) {

	log.Printf("[POCSAG-Pilemaster] Publishing to %s: \n%s\n", topic, msg)
	token := client.Publish(topic, 2, false, msg)
	token.Wait()
	time.Sleep(time.Second)

}

func main() {
	var broker string = os.Getenv("POCSAGPILEMASTER_BROKER")
	var port, err = strconv.Atoi(os.Getenv("POCSAGPILEMASTER_PORT"))
	if err != nil {
		log.Fatalln("Port is not an integer")
	}
	var clientID string = os.Getenv("POCSAGPILEMASTER_CLIENTID")
	var username string = os.Getenv("POCSAGPILEMASTER_USERNAME")
	var password string = os.Getenv("POCSAGPILEMASTER_PASSWORD")
	var debug string = os.Getenv("POCSAGPILEMASTER_DEBUG")
	var topic string = os.Getenv("POCSAGPILEMASTER_TOPIC")

	log.Printf("[POCSAG-Pilemaster] Connecting to %s:%d as %s with user %s\n", broker, port, clientID, username)
	c := make(chan os.Signal, 1)
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetClientID(clientID)
	opts.SetUsername(username)
	opts.SetPassword(password)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	var text string

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		text = line
		var words []string = strings.Fields(text)

		var addressStr string
		var protocol string
		var function string
		var message string

		// Somewhat arbitrarily ignore lines that are too short to be processed
		if len(words) >= 5 {

			protocol = words[2]
			addressStr = words[4]
			function = words[6]

			protocol = strings.Replace(protocol, ":", "", -1)

			address, err := strconv.Atoi(strings.Replace(addressStr, ":", "", -1))

			if err == nil {

				message_fix := strings.Join(words[8:], " ")

				// Obnoxious maze of esoteric newline corrections etc
				message = strings.Replace(message_fix, "<CR><LF>", "\n", -1)
				message = strings.Replace(message, "<NUL>", "", -1)
				message = strings.Replace(message, "<LF>", "\n", -1)
				message = strings.Replace(message, "<CR>", "\n", -1)

				if debug == "YES" {
					log.Printf("DEBUG: protocol: %s address: %s function: %s msg: %s \n", protocol, addressStr, function, message)
				}
				// Skip these strange messages.
				if message_fix != "AB" && message_fix != "CD" {

					// What the heck is function 3 anyways? At least it contains useful data.
					if function == "3" {
						msg := fmt.Sprintf("%s,%d,%s", protocol, address, message)
						publish(client, msg, topic)
					}

				}
			}
		} else {
			fmt.Println("Malformed input. Skipping..")
		}

	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Error reading standard input:", err)
	}

	<-c
}
