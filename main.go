package main

import (
	"bufio"
	"fmt"
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

	lastActiveClass := ""
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
