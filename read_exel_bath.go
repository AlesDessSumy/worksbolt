package main

import (
	"fmt"

	"google.golang.org/api/sheets/v4"
)

// читает лист "Пустые строки"
// и скидывает в канал
// элеенты "empty_string"
func read_empt_strings(multicoin, multicoin_range string) {

	//ctx := context.Background()

	// srv, err := sheets.NewService(ctx, option.WithCredentialsFile(path+"credentials.json"))
	// if err != nil {
	// 	log.Fatalf("Unable to retrieve Sheets client: %v", err)
	// }
	srv := get_SCHEETS_client()

	n := 0

	wrr := "'Лист с таблицей'!A5136"

	n++
	fmt.Println("Start round N:", n)

	sheet, _ := srv.Spreadsheets.Get(multicoin).IncludeGridData(true).Ranges(multicoin_range).Do()

	// if err != nil {
	// 	errorLog.Println(err)
	// 	<-time.After(5 * time.Second)
	// 	continue
	// }

	data := sheet.Sheets

	dt := make([]*sheets.RowData, 0)

	for _, i := range data {
		//	fmt.Println("read data")
		for _, g := range i.Data {
			for _, row := range g.RowData {
				*row.Values[0].UserEnteredValue.NumberValue = 55578
				*row.Values[2].UserEnteredValue.NumberValue = 0
				*row.Values[3].UserEnteredValue.NumberValue = 0
				*row.Values[4].UserEnteredValue.NumberValue = 0
				*row.Values[5].UserEnteredValue.NumberValue = 0
				*row.Values[6].UserEnteredValue.NumberValue = 0
				*row.Values[7].UserEnteredValue.NumberValue = 0
				*row.Values[8].UserEnteredValue.NumberValue = 0
				*row.Values[9].UserEnteredValue.NumberValue = 0
				*row.Values[10].UserEnteredValue.NumberValue = 0
				*row.Values[11].UserEnteredValue.NumberValue = 0
				*row.Values[12].UserEnteredValue.NumberValue = 0
				*row.Values[13].UserEnteredValue.NumberValue = 0
				*row.Values[14].UserEnteredValue.NumberValue = 0
				*row.Values[15].UserEnteredValue.NumberValue = 0
				*row.Values[16].UserEnteredValue.NumberValue = 0
				*row.Values[17].UserEnteredValue.NumberValue = 0
				*row.Values[18].UserEnteredValue.NumberValue = 0
				dt = append(dt, row)
			}
		}
	}
	var vr sheets.ValueRange
	myval := make([][]interface{}, 0)
	for _, row := range dt {
		mv := make([]interface{}, 0)
		for n, g := range row.Values {
			fmt.Println(n)
			if n == 19 {
				break
			}
			if n == 1 {
				mv = append(mv, g.UserEnteredValue.StringValue)
			} else {
				mv = append(mv, g.UserEnteredValue.NumberValue)
			}

		}
		myval = append(myval, mv)
	}
	for _, i := range myval {
		fmt.Println(i)
	}

	vr.Values = append(vr.Values, myval...)
	_, err := srv.Spreadsheets.Values.Append(multicoin, wrr, &vr).ValueInputOption("RAW").Do() // USER_ENTERED
	errorLog.Println(err)
}
