package driver

import (
	"github.com/webhkp/godft/internal/consts"
)

// Common interface for all the drivers
type Driver interface {
	Execute(*consts.FlowDataSet)
	Validate() bool
	GetInput() (string, bool)

	Read(data *consts.FlowDataSet)
	Write(data *consts.FlowDataSet)
}

