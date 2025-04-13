package kit

var (
	defaultCopier *Copier
)

func init() {
	defaultCopier, _ = NewCopier()
}

type CopierOptsFunc func(*CopierOpts)

type CopierOpts struct {
	tagName string
}

type Copier struct {
	opts *CopierOpts
}

func NewCopier(opts ...CopierOptsFunc) (*Copier, error) {
	var copierOpts CopierOpts
	for _, opt := range opts {
		opt(&copierOpts)
	}
	return &Copier{
		opts: &copierOpts,
	}, nil
}

func (c *Copier) Copy(src, dst interface{}) error {
	return nil
}
