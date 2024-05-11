package gredis

type Option struct {

}

type Redis interface {
	NewClient() error
}
