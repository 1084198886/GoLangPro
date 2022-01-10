package gormio

import (
	"fmt"
	"gopkg.in/natefinch/lumberjack.v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"path"
	"sync"
	"time"
)

const (
	logPath     = "F:/log"
	logFileName = "gormio.log"
)

type User struct {
	Id       int       `gorm:"auto_increment;primarykey;comment:'主键'"`
	Name     string    `gorm:"type:varchar(10);default:'default'"`
	Stuempno string    `gorm:"type:char(18);unique;not null"`
	UUID     string    `gorm:"type:varchar(18)"`
	Fee      float64   `gorm:"type decimal(5,2);default 0"`
	UpdateAt int64     `gorm:"autoUpdateTime:nano"`
	CreateAt time.Time `gorm:"comment:'创建时间';"`
}

//func (u *User) BeforeDelete(db *gorm.DB) error {
//	if u.Stuempno == "002" {
//		return errors.New("002 not allowed to delete")
//	}
//	return nil
//}
//
//func (u *User) BeforeCreate(db *gorm.DB) error {
//	u.UUID = fmt.Sprintf("%d", rand.Intn(100))
//	if u.Stuempno == "002" {
//		return errors.New("002 not allowed to create")
//	}
//	return nil
//}

type gormioStorage struct {
	db databaseWrapper
	mu sync.Mutex
}

func (s *gormioStorage) cfg() *gorm.DB {
	return s.db.User()
}

func (s *gormioStorage) close() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.db != nil {
		_ = s.db.Close()
	}
}

func (s *gormioStorage) connect(dsn string) error {
	wrapper, err := loadDatabaseWrapper(dsn)
	if err != nil {
		return err
	}
	err = wrapper.Open(dsn)
	if err != nil {
		return err
	}
	s.db = wrapper
	return s.migrate()
}

func (s *gormioStorage) migrate() error {
	if err := s.cfg().AutoMigrate(&User{}); err != nil {
		return err
	}
	return nil
}

type databaseWrapper interface {
	Open(dsn string) error
	User() *gorm.DB
	Close() error
}

func loadDatabaseWrapper(driver string) (databaseWrapper, error) {
	if driver == "test" {
		return &testDatabaseWrapper{}, nil
	} else {
		panic(fmt.Sprintf("not support db driver %v", driver))
	}
}

type testDatabaseWrapper struct {
	databaseWrapper
	db *gorm.DB
}

func (s *testDatabaseWrapper) Open(dsn string) error {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		return err
	}
	s.db = db
	return nil
}

func (s *testDatabaseWrapper) User() *gorm.DB {
	return s.db
}

func (s *testDatabaseWrapper) Close() error {
	if s.db != nil {
		db, err := s.db.DB()
		if err != nil {
			return err
		}
		db.Close()
		s.db = nil
	}
	return nil
}

func SetLogOutput() {
	if len(logPath) == 0 {
		return
	}
	log.SetOutput(&lumberjack.Logger{
		Filename:   path.Join(logPath, logFileName),
		MaxSize:    500,
		MaxBackups: 3,
		MaxAge:     30,
		Compress:   true,
	})
}
