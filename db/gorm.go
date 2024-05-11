package db

import "gorm.io/gorm"

type GormDB struct {
	DB *gorm.DB
}

func (g *GormDB) Migrate(source string) error {
	err := g.Open(source)
	if err != nil {
		return err
	}
	return nil
}

func (g *GormDB) Open(source string) error {
	db, err := gorm.Open(nil, &gorm.Config{})
	if err != nil {
		return err
	}
	g.DB = db
	return nil
}

func (g *GormDB) Close() error {
	//TODO implement me
	panic("implement me")
}

func (g *GormDB) Backend() {
	//TODO implement me
	panic("implement me")
}

func (g *GormDB) Tx() error {
	//TODO implement me
	panic("implement me")
}

func (g *GormDB) Commit() error {
	//TODO implement me
	panic("implement me")
}
