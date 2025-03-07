package ghttp

type Resolver interface {
	Resolve(string) string
}

type DefaultResolver struct {
	baseUrl string
}

func NewDefaultResolver(baseUrl string) *DefaultResolver {
	return &DefaultResolver{baseUrl: baseUrl}
}
