package capture

type Handler interface {
	//Raw() []byte
	Start() error
	Stop()
	Packets() <-chan []byte
}

type CaptureOptionFunc func(*CaptureOption)

type CaptureOption struct {
}

type Capture struct {
	Options *CaptureOption
	handler Handler
}

//func NewCapture(opts ...CaptureOptionFunc) (*Capture, error) {
//	var options CaptureOption
//	for _, opt := range opts {
//		opt(&options)
//	}
//	var h Handler
//	switch runtime.GOOS {
//	case "linux":
//		h = &LinuxHandler{}
//	case "windows":
//		h = &WindowsHandler{}
//	default:
//		return nil, fmt.Errorf("no support os")
//	}
//	return &Capture{Options: &options, handler: h}, nil
//}

func (c *Capture) Raw() {

}
