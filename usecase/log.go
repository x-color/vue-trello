package usecase

import (
	"errors"
	"fmt"

	"github.com/x-color/vue-trello/model"
)

func logError(logger Logger, err error) {
	if errors.Is(err, model.ServerError{}) {
		logger.Error(err.Error())
	} else {
		logger.Info(err.Error())
	}
	return
}

func formatLogMsg(id, msg string) string {
	return fmt.Sprintf("%s %s", id, msg)
}
