package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

const NETWORK_NAME = "unix"

func main() {

	xdgRuntimeDir := os.Getenv("XDG_RUNTIME_DIR")
	instanceSig := os.Getenv("HYPRLAND_INSTANCE_SIGNATURE")
	if len(xdgRuntimeDir) == 0 || len(instanceSig) == 0 {
		log.Fatal("make sure you ran hyprland session")
	}
	// /run/user/1000/hypr/efb50993780079460b0cbed1363e2166a2de1d9f_1791459651_1868869060/.socket2.sock
	socketPath := xdgRuntimeDir + "/hypr/" + instanceSig + "/.socket2.sock"

	conn, err := net.Dial(NETWORK_NAME, socketPath)

	if err != nil {
		log.Fatalf("%s %v", "Connection error: ", err)
	}
	defer conn.Close()

	scanner := bufio.NewScanner(conn)

	lastActiveClass, err := getActiveClass()
	if err != nil {
		log.Fatalf("%s %v", "Error with getting active class: ", err)
	}
	intervalStart := time.Now()

	for scanner.Scan() {
		line := scanner.Text()
		now := time.Now()
		_, after, found := strings.Cut(line, "activewindow>>")
		if !found {
			continue
		}
		currAppClass, _, _ := strings.Cut(after, ",")

		if lastActiveClass == currAppClass {
			continue
		}

		duration := now.Sub(intervalStart)
		fmt.Printf("%s %v\n", "active window: "+lastActiveClass+": ", duration)
		lastActiveClass = currAppClass
		intervalStart = now
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("%v", err)
	}
}

func getActiveClass() (string, error) {

	xdgRuntimeDir := os.Getenv("XDG_RUNTIME_DIR")
	instanceSig := os.Getenv("HYPRLAND_INSTANCE_SIGNATURE")
	if len(xdgRuntimeDir) == 0 || len(instanceSig) == 0 {
		return "", fmt.Errorf("make sure you ran hyprland session")
	}
	socketPath := xdgRuntimeDir + "/hypr/" + instanceSig + "/.socket.sock"
	timeout := 2 * time.Second
	conn, err := net.DialTimeout(NETWORK_NAME, socketPath, timeout)

	if err != nil {
		return "", fmt.Errorf("%s %v\n", "Connection error: ", err)
	}
	defer conn.Close()
	if err := conn.SetDeadline(time.Now().Add(timeout)); err != nil {
		return "", fmt.Errorf("%s %v\n", "Connection error: ", err)
	}
	command := "/activewindow"

	_, err = conn.Write([]byte(command))

	if err != nil {
		return "", fmt.Errorf("%s %v", "Error with writing a command to socket: ", err)
	}

	var response bytes.Buffer
	_, err = io.Copy(&response, conn)
	if err != nil {
		return "", fmt.Errorf("Failed to read from socket: %v", err)
	}

	scanner := bufio.NewScanner(&response)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "class: ") {
			return strings.TrimPrefix(line, "class: "), nil
		}
	}

	return "", fmt.Errorf("property 'class' not found in hyprland response")
}
