package main

import (
	"bufio"
	"context"
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

	store, err := storage.Connect()
	if err != nil {
		log.Fatalf("%s %v", "db error: ", err)
	}
	defer store.Close()
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
			if err := store.SaveSession(ctx, session); err != nil {
				log.Printf("save session: %v", err)
			}
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
		saveCtx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		if err := store.SaveSession(saveCtx, session); err != nil {
			log.Printf("save session: %v", err)
		}
	}
}
