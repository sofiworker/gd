package gerr

type GErr interface {
	String() string
	Wrap(error) GErr
}
