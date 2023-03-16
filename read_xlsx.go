package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	"trs"

	"github.com/tealeg/xlsx/v3"
)

var path_behavior = path + "exel/behavior"
var path_brend_bolt = path + "exel/brend_bolt"
var path_brend_uber = path + "exel/brend_uber"
var path_brend_uklon = path + "exel/brend_uklon"
var path_speed = path + "exel/speed"
var path_traval = path + "exel/traval"

func Read_Brend_Uber() {

	citycar := make(map[string][2]string)

	day := last_mounth_voskr()

	for _, cit := range CIT {
		gr := get_grafik(cit, day[:len(day)-1])
		for _, i := range gr {
			citycar[i.car_nomer] = [2]string{cit.name, cit.id_BGQ}
		}
	}

	n := 0

	car := make([]Brend_Uber, 0)

	files, err := ioutil.ReadDir(path_brend_uber)
	if err != nil {
		fmt.Println(err)
	}

	Cit := make([]City_Brend_Uber, 0)
	Car := make([]Car_Brend_Uber, 0)

	for _, file := range files {

		file, err := os.Open(path_brend_uber + "/" + file.Name())
		if err != nil {
			log.Fatal(err)
		}

		scanner := bufio.NewScanner(file)

		n := 0
		for scanner.Scan() {
			n++
			s := scanner.Text()
			ss := strings.Split(s, ",")
			if n == 1 {
				continue
			}
			if strings.Contains(ss[0], "all") {
				r := City_Brend_Uber{}
				r.Create(ss, n)
				Cit = append(Cit, r)
			} else {
				r := Car_Brend_Uber{}
				r.Create(ss, n)
				Car = append(Car, r)
			}
		}

		file.Close()

		poteriashka := make([]Car_Brend_Uber, 0)

		// присваиваем города автомобилям
		for j := 0; j < len(Car); j++ {
			nom, ok := citycar[Car[j].Nomer]
			if ok {
				Car[j].Up_City(nom[1])
			} else {
				fmt.Println("no search car")
				poteriashka = append(poteriashka, Car[j])
			}
		}

		// if len(poteriashka) > 0 {
		// 	for _, ct := range Cit {
		// 		fmt.Println(ct.Trip, ct.City, City_Name[psevdonimy[ct.City]])
		// 		trip := 0
		// 		for _, cr := range Car {
		// 			if cr.City == City_Name[psevdonimy[ct.City]] {
		// 				trip = trip + cr.Trip
		// 			}
		// 		}
		// 		fmt.Println(trip, ct.Trip)

		// 	}
		// }

		// <-time.After(time.Hour)

		// присваиваем бренд автомобилям
		for i := 0; i < len(Cit); i++ {
			trip := 0
			for _, j := range Car {
				if j.City == City_Name[psevdonimy[Cit[i].City]] {
					trip = trip + j.Trip
				}
			}
			paytrip := Cit[i].Pay / float64(trip)
			for j := 0; j < len(Car); j++ {
				if Car[j].City == City_Name[psevdonimy[Cit[i].City]] {
					Car[j].Up_Pay(paytrip * float64(Car[j].Trip))
				}
			}
		}

		// for n, m := range Car {
		// 	fmt.Println(n, m)
		// }

		// <-time.After(time.Hour)

		vsc := float64(len(day))

		for _, d := range day {
			dt := d.Format(TNS)
			for _, i := range Cit {
				r := Brend_Uber{}
				y, n := d.ISOWeek()
				r.Year = y
				r.Ned = n
				r.City = psevdonimy[i.City]
				r.Park = City_Name[r.City]
				r.Date = dt
				r.Pay = i.Pay / vsc
				r.Trip = int(float64(i.Trip)/vsc + 0.5)
				r.Type = "city"
				car = append(car, r)
			}
			for _, i := range Car {
				r := Brend_Uber{}
				y, n := d.ISOWeek()
				r.Year = y
				r.Ned = n
				r.City = psevdonimy[i.City]
				r.Park = City_Name[r.City]
				r.Date = dt
				r.Pay = i.Pay / vsc
				r.Trip = int(float64(i.Trip)/vsc + 0.5)
				r.Car = i.Nomer
				r.Type = "car"
				car = append(car, r)
			}
		}

		db, err := sql.Open("sqlite3", "C:/bolt_db/sqlite/profit.db")
		if err != nil {
			fmt.Println(path+"sqlite/profit.db", err)
		}

		dat := make([]string, 0)
		for _, i := range car {
			if !comp_string(i.Date, dat) {
				dat = append(dat, i.Date)
			}
		}

		zap := "DELETE FROM Brend_Uber WHERE "
		for n, t := range dat {
			if n == 0 {
				zap = zap + fmt.Sprintf("Date = '%s'", t)
			} else {
				zap = zap + fmt.Sprintf(" OR Date = '%s'", t)
			}
		}

		// fmt.Println(zap)

		_, err = db.Exec(zap)
		if err != nil {
			fmt.Println(err, zap)
			return
		}

		for _, i := range car {
			i.RecDB(db)
			n++
		}

		db.Close()

		// fmt.Println(file.Name())

		err = os.Remove(file.Name())
		if err != nil {
			fmt.Println(err)
		}
	}

	OS(NS(fmt.Sprintf("Бренд Убер сделано %d записей", n), 40))
}

func read_xlsx() (res_c, res_d [][][]string) {

	f, err := xlsx.OpenFile(path + "mail/trip.xlsx")
	if err != nil {
		fmt.Println("Error Open trip.xlsx", err)
	} else {
		sh := f.Sheets[0]
		row := sh.MaxRow

		rn := 0
		rk := 0

		rn = 8
		rk = row - 1

		c_dat, _ := sh.Cell(1, 0)
		dat, _ := time.Parse("02.01.2006", string([]byte(c_dat.Value)[len([]byte(c_dat.Value))-10:]))

		res_f := make([][]string, 0)
		for i := rn; i < rk; i++ {
			var hr []string
			hr = append(hr, dat.Format("02.01.2006"))
			n_dat, _ := sh.Cell(i, 0)
			hr = append(hr, trs.Trs(n_dat.Value))
			db := ""
			for j := 2; j < 26; j++ {
				c, _ := sh.Cell(i, j)
				s := c.Value
				f, _ := strconv.ParseFloat(s, 64)
				db = db + strconv.Itoa(int(f*10)) + ","
			}
			hr = append(hr, db)
			res_f = append(res_f, hr)
		}
		res_c = append(res_c, res_f)
	}

	fs, err := xlsx.OpenFile(path + "mail/opov.xlsx")
	if err != nil {
		fmt.Println("Error Open opov.xlsx", err)
	} else {
		sh := fs.Sheets[0]
		rows := sh.MaxRow
		cols := sh.MaxCol

		rns := 2
		rks := rows
		res := make([][]string, 0)
		for i := rns; i < rks; i++ {
			rs := make([]string, 0)
			for j := 0; j < cols; j++ {
				c, _ := sh.Cell(i, j)
				rs = append(rs, c.Value)
			}
			if !strings.Contains(rs[6], "Превышение скорости") {
				continue
			}
			adr, tr := get_speed_adress(rs[3])
			if !tr {
				continue
			}
			rs[3] = adr
			res = append(res, rs)
		}
		res_d = append(res_d, res)

	}
	return
}
