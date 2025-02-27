package db

import "gorm.io/gorm"

type GormClient struct {
	DB *gorm.DB
}

func NewGormClient() *GormClient {
	return &GormClient{}
}

func (g *GormClient) Migrate(source string) error {
	//err := g.Open(source)
	//if err != nil {
	//	return err
	//}
	return nil
}

func (g *GormClient) Close() error {
	//TODO implement me
	panic("implement me")
}

func (g *GormClient) Backend() {
	//TODO implement me
	panic("implement me")
}

func (g *GormClient) Tx() error {
	//TODO implement me
	panic("implement me")
}

func (g *GormClient) Commit() error {
	//TODO implement me
	panic("implement me")
}
