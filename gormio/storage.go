package gormio

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

/*数据库驱动*/
type databaseWrapper interface {
	Open(dsn string) error
}

/*加载数据库驱动*/
func loadDatabaseWrapper(driver string) (databaseWrapper, error) {
	switch driver {
	case "test":
		return &testDatabaseWrapper{}, nil
	}
	return nil, nil
}

type testDatabaseWrapper struct {
	databaseWrapper *gorm.DB
}

func (t *testDatabaseWrapper) Open(dsn string) error {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		return err
	}
	t.databaseWrapper = db
	return nil
}
