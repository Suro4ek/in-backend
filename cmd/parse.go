package main

import (
	"context"
	"fmt"
	"github.com/xuri/excelize/v2"
	"in-backend/internal/composites"
	"in-backend/internal/config"
	"in-backend/internal/items"
	"in-backend/pkg/logging"
	"in-backend/pkg/postgres"
	"strings"
)

func main1() {
	logging.Init()
	logger := logging.GetLogger()
	logger.Info("create router")

	cfg := config.GetConfig()

	client := postgres.NewClient(context.Background(), 5, cfg.Postgres)

	userComposite, _ := composites.NewUserComposite(client, &logger, cfg)
	item, _ := composites.NewItemComposite(client, userComposite.Repository, &logger)
	f, err := excelize.OpenFile("rest.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	rows, err := f.GetRows("Реестр")
	if err != nil {
		fmt.Println(err)
		return
	}
	for index, row := range rows {
		if index < 8 {
			continue
		}
		if row == nil {
			continue
		}
		var itm *items.Item
		if len(row) == 4 {
			itm = &items.Item{
				ProductName:  row[2],
				Name:         row[3],
				SerialNumber: "нету",
			}
		} else {
			if strings.TrimSpace(row[4]) == "" {
				itm = &items.Item{
					ProductName:  row[2],
					Name:         row[3],
					SerialNumber: "нету",
				}
			} else {
				itm = &items.Item{
					ProductName:  row[2],
					Name:         row[3],
					SerialNumber: row[4],
				}
			}
		}

		fmt.Print(itm.ProductName+" "+itm.SerialNumber+" "+itm.Name+" ", "\t")
		err := item.Repository.Create(context.TODO(), itm)
		if err != nil {
			return
		}
		//for index1, colCell := range row {
		//	if index1 < 2 || index1 > 4 {
		//		continue
		//	}
		//	if strings.TrimSpace(colCell) == "" {
		//		continue
		//	}
		//
		//	fmt.Print(colCell, "\t")
		//}
	}
}
