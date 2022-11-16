package informer

import (
	"fmt"
	"time"

	"github.com/6za/k1-watcher/pkg/shell"
	"go.uber.org/zap"
)

func TestGeneric(configFile string, loggerIn *zap.Logger) error {
	logger = loggerIn
	go func() {
		_, _, err := shell.ExecShellReturnStrings("kubectl", "wait", "--for=condition=Ready", "pod/sample", "--timeout=15s")
		if err != nil {
			logger.Info(err.Error())
		} else {
			logger.Info("Condition met")
		}

	}()
	time.Sleep(30 * time.Second)

	return fmt.Errorf("timeout - Failed to meet exit condition")

}
