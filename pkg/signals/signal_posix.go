/*
@Time : 2022/11/26 22:11
@Author : lianyz
@Description :
*/

package signals

import (
	"os"
	"syscall"
)

var shutdownSignals = []os.Signal{os.Interrupt, syscall.SIGTERM}
