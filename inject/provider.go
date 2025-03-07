package inject

type NameOption string

func (o NameOption) ApplyProvideOption(opt *ProvideOptions) {
	opt.Name = string(o)
}

func WithName(name string) ProvideOption {
	return NameOption(name)
}

type GroupOption string

func (o GroupOption) ApplyProvideOption(opt *ProvideOptions) {
	opt.Group = string(o)
}

func WithGroup(group string) ProvideOption {
	return GroupOption(group)
}

type ProvideAsOption []interface{}

func (o ProvideAsOption) ApplyProvideOption(opts *ProvideOptions) {
	opts.As = append(opts.As, o...)
}

func WithAs(data ...interface{}) ProvideOption {
	return ProvideAsOption(data)
}
