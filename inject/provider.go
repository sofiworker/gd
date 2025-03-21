package inject

type Provider struct {
	Key         string
	Target      interface{}
	Constructor interface{}
	Single      bool
}
