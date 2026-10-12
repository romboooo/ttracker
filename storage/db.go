package storage

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/romboooo/ttracker/tracker"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Session struct {
	gorm.Model
	Class     string
	StartTime time.Time
	EndTime   time.Time
}

type Storage struct {
	db *gorm.DB
}

func (s *Storage) SessionsBetween(ctx context.Context, from, to time.Time) ([]tracker.Session, error) {
	var sessions []tracker.Session
	if err := s.db.WithContext(ctx).Where("start_time < ? AND end_time > ?", to, from).Find(&sessions).Error; err != nil {
		return nil, err
	}
	return sessions, nil
}

func Connect() (*Storage, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	dbPath := filepath.Join(homeDir, ".local", "share", "ttracker", "ttracker.db")

	dir := filepath.Dir(dbPath)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&Session{}); err != nil {
		return nil, err
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveSession(ctx context.Context, session tracker.Session) error {

	if session.StartTime.Equal(session.EndTime) || !session.EndTime.After(session.StartTime) {
		return fmt.Errorf("end time should be after start time")
	}

	if err := s.db.WithContext(ctx).Create(&Session{
		Class:     session.Class,
		StartTime: session.StartTime.UTC(),
		EndTime:   session.EndTime.UTC(),
	}).Error; err != nil {
		return err
	}

	return nil
}

func (s *Storage) Close() error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
