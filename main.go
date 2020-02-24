package main

import (
	"fmt"
	"os"

	"github.com/x-color/vue-trello/interface/controller/api"
	"github.com/x-color/vue-trello/interface/presenter/logging"
	"github.com/x-color/vue-trello/interface/repository/rdb"
	"github.com/x-color/vue-trello/usecase"
)

func main() {
	dbm, err := rdb.NewDBManager()
	if err != nil {
		fmt.Println(err)
		return
	}

	logger := logging.NewLogger(os.Stdout)

	itemIntera, err := usecase.NewItemInteractor(
		&dbm.ItemDBManager,
		&dbm.ListDBManager,
		&dbm.TagDBManager,
		&logger,
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	listIntera, err := usecase.NewListInteractor(
		&dbm.ListDBManager,
		&dbm.BoardDBManager,
		&logger,
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	boardIntera, err := usecase.NewBoardInteractor(
		&dbm.BoardDBManager,
		&dbm.ListDBManager,
		&dbm.ItemDBManager,
		&logger,
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	userIntera, err := usecase.NewUserInteractor(
		&dbm.UserDBManager,
		&logger,
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	resourceIntera, err := usecase.NewResourceInteractor(
		&dbm.TagDBManager,
		&logger,
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	interaBox, err := api.NewInteraBox(
		&itemIntera,
		&listIntera,
		&boardIntera,
		&userIntera,
		&resourceIntera,
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	router := api.NewRouter(interaBox)
	router.Logger.Fatal(router.Start(":8080"))
}
