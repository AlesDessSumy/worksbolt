package main

import (
	"BGQ"
	"fmt"
	"image/png"
	"os"
	"trs"

	// создание  QrCode
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"

	// чтение  QrCode
	"github.com/tuotoo/qrcode"
)

func QrCode_Save() {

	// справочник авто в квери
	r, _ := BGQ.Read_BGQ("select * from Bitrix.cars")
	cars := make(map[string]string, 0)
	for _, i := range r {
		cars[trs.Trs(i[1].(string))] = i[7].(string)

	}

	for _, sit := range create_sitys() {
		grf := get_grafik(sit, teck_days())
		for _, g := range grf {
			if cars[g.car_nomer] == "" {
				fmt.Println("Ошибка", sit.name, g.car_nomer, cars[g.car_nomer])
				continue
			}
			f := "./QrCode/" + sit.name + "_" + g.car_nomer + ".png"

			s := cars[g.car_nomer]

			qrCode, _ := qr.Encode(s, qr.M, qr.Auto)

			qrCode, _ = barcode.Scale(qrCode, 256, 256)

			file, _ := os.Create(f)
			defer file.Close()

			png.Encode(file, qrCode)
		}
	}
}

func QrCode_Read() {

	fi, err := os.Open("./QrCode/qr2.png")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer fi.Close()
	qrmatrix, err := qrcode.Decode(fi)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(qrmatrix.Content)
}
