package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/romboooo/ttracker/hyprland"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	conn, err := hyprland.ConnectEvents()

	if err != nil {
		log.Fatalf("%s %v", "Connection error: ", err)
	}
	defer conn.Close()

	go func() {
		<-ctx.Done()
		conn.Close()
	}()

	scanner := bufio.NewScanner(conn)

	lastActiveClass, err := hyprland.GetActiveClass()
	if err != nil {
		log.Fatalf("%s %v", "Error with getting active class: ", err)
	}

	intervalStart := time.Now()
	for scanner.Scan() {
		currAppClass, ok := hyprland.ParseActiveWindow(scanner.Text())
		if !ok {
			continue
		}

		now := time.Now()
		if lastActiveClass == currAppClass {
			continue
		}

		duration := now.Sub(intervalStart)
		fmt.Printf("%s %v\n", "active window: "+lastActiveClass+": ", duration)
		lastActiveClass = currAppClass
		intervalStart = now
	}

	if err := scanner.Err(); err != nil && ctx.Err() == nil {
		log.Fatalf("%v", err)
	}

	if ctx.Err() != nil {
		fmt.Printf("%s %v\n", "active window: "+lastActiveClass+": ", time.Since(intervalStart))
	}
}
