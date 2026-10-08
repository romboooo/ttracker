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
	"github.com/romboooo/ttracker/storage"
	"github.com/romboooo/ttracker/tracker"
)

func main() {
	fmt.Print("!")
	storage, err := storage.Connect()
	if err != nil {
		log.Fatalf("%s %v", "db error: ", err)
	}
	defer storage.Close()
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

		session := tracker.Session{
			Class:     lastActiveClass,
			StartTime: intervalStart,
			EndTime:   now,
		}
		if lastActiveClass != "" {
			tracker.PrintSession(session)

		}
		lastActiveClass = currAppClass
		intervalStart = now
	}

	if err := scanner.Err(); err != nil && ctx.Err() == nil {
		log.Fatalf("%v", err)
	}

	session := tracker.Session{
		Class:     lastActiveClass,
		StartTime: intervalStart,
		EndTime:   time.Now(),
	}

	if ctx.Err() != nil && session.Class != "" {

		tracker.PrintSession(session)
	}
}
