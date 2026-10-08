package hyprland

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"
)

const NETWORK_NAME = "unix"

func ConnectEvents() (net.Conn, error) {
	xdgRuntimeDir := os.Getenv("XDG_RUNTIME_DIR")
	instanceSig := os.Getenv("HYPRLAND_INSTANCE_SIGNATURE")
	if xdgRuntimeDir == "" || instanceSig == "" {
		return nil, fmt.Errorf("make sure you ran hyprland session")
	}

	socketPath := xdgRuntimeDir + "/hypr/" + instanceSig + "/.socket2.sock"

	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		return nil, fmt.Errorf("connect to hyprland event socket: %w", err)
	}

	return conn, nil
}

func GetActiveClass() (string, error) {

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

func ParseActiveWindow(line string) (string, bool) {
	data, found := strings.CutPrefix(line, "activewindow>>")
	if !found {
		return "", false
	}

	class, _, found := strings.Cut(data, ",")
	if !found {
		return "", false
	}

	return class, true
}
