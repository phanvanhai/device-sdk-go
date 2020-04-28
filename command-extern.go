package device

import (
	"fmt"

	"github.com/edgexfoundry/device-sdk-go/internal/common"
	"github.com/edgexfoundry/device-sdk-go/internal/handler"
	contract "github.com/edgexfoundry/go-mod-core-contracts/models"
)

// CommandRequest :mo rong SDK cho phep Device-Service gui lenh cho cac doi tuong trong cung Device-Service
func CommandRequest(devName string, cmName string, body string, method string, queryParams string) (*contract.Event, error) {
	vars := make(map[string]string, 2)
	vars[common.NameVar] = devName
	vars[common.CommandVar] = cmName

	event, appErr := handler.CommandHandler(vars, body, method, queryParams)

	if appErr != nil {
		return nil, fmt.Errorf(appErr.Message())
	} else if event != nil {
		return &event.Event, nil
	}
	return nil, nil
}
