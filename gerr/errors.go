package gerr

import "fmt"

var (
	NotFoundMethod            = fmt.Errorf("not found method")
	NotAllowMultiLayerPointer = fmt.Errorf("not allow multi layer pointer")
	ParamsMustBeFunc          = fmt.Errorf("params must be func")
	InitDefaultLoggerFailed   = fmt.Errorf("init default logger failed")
)
