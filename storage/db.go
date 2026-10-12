package storage

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/romboooo/ttracker/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type SessionRecord struct {
	gorm.Model
	Class     string
	StartTime time.Time
	EndTime   time.Time
}

func (SessionRecord) TableName() string {
	return "sessions"
}

type Storage struct {
	db *gorm.DB
}

func (s *Storage) SessionsBetween(ctx context.Context, from, to time.Time) ([]model.Session, error) {
	var records []SessionRecord
	if err := s.db.WithContext(ctx).Where("start_time < ? AND end_time > ?", to, from).Find(&records).Error; err != nil {
		return nil, err
	}

	sessions := make([]model.Session, 0, len(records))
	for _, record := range records {
		sessions = append(sessions, model.Session{
			Class:     record.Class,
			StartTime: record.StartTime,
			EndTime:   record.EndTime,
		})
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

	if err := db.AutoMigrate(&SessionRecord{}); err != nil {
		return nil, err
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveSession(ctx context.Context, session model.Session) error {

	if session.StartTime.Equal(session.EndTime) || !session.EndTime.After(session.StartTime) {
		return fmt.Errorf("end time should be after start time")
	}

	if err := s.db.WithContext(ctx).Create(&SessionRecord{
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
