package main

import (
	"fmt"
	"os"
	"time"

	"github.com/x-color/vue-trello/interface/controller/api"
	"github.com/x-color/vue-trello/interface/presenter/logging"
	"github.com/x-color/vue-trello/interface/repository/rdb"
	"github.com/x-color/vue-trello/model"
	"github.com/x-color/vue-trello/usecase"
)

func main() {
	dbm, err := rdb.NewDBManager()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Add init tags.
	tx := dbm.TransactionManager.BeginTransaction(false)
	dbm.TagDBManager.Create(tx, model.Tag{ID: "0", Name: "t1", Color: model.RED})
	dbm.TagDBManager.Create(tx, model.Tag{ID: "1", Name: "t2", Color: model.YELLOW})
	dbm.TagDBManager.Create(tx, model.Tag{ID: "2", Name: "t3", Color: model.GREEN})
	dbm.TagDBManager.Create(tx, model.Tag{ID: "3", Name: "t4", Color: model.BLUE})

	location, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		location = time.FixedZone("Asia/Tokyo", 9*60*60)
	}
	logger, err := logging.NewLogger(os.Stdout, location)
	if err != nil {
		fmt.Println(err)
		return
	}

	itemIntera, err := usecase.NewItemInteractor(
		&dbm.TransactionManager,
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
		&dbm.TransactionManager,
		&dbm.ItemDBManager,
		&dbm.ListDBManager,
		&dbm.BoardDBManager,
		&logger,
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	boardIntera, err := usecase.NewBoardInteractor(
		&dbm.TransactionManager,
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
		&dbm.TransactionManager,
		&dbm.UserDBManager,
		&logger,
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	resourceIntera, err := usecase.NewResourceInteractor(
		&dbm.TransactionManager,
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
