package main

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

// запрос:
// список городов
// список дат
var trips_dbase = make(chan gets_trips_dbase, 0)
var trips_driver_dbase = make(chan gets_driver_trips_dbase, 0)

type gets_trips_dbase struct {
	city []sity
	days []time.Time
	x    chan []Trips
}

type gets_driver_trips_dbase struct {
	name string
	days []time.Time
	x    chan []Trips
}

// запрос:
// элемент графика
var trips_dbase_graf = make(chan gets_trips_dbase_graf, 0)

type gets_trips_dbase_graf struct {
	city string
	days time.Time
	g    *grafik
	x    chan [2]int
}

// поездки по городу за день
// с разбивкой по службам
type trips_Citys_days struct {
	date          string
	name          string
	name_Cyrrilic string
	bolt          int
	uklon         int
	uber          int
	trips         int
	cash_bolt     float64
	cash_uklon    float64
	cash_uber     float64
	cash          float64
}

func (a *trips_Citys_days) Print() {
	fmt.Println()
	fmt.Println("type trips_Citys_days struct :")
	fmt.Println("date                         :", a.date)
	fmt.Println("name                         :", a.name_Cyrrilic)
	fmt.Println("name                         :", a.name)
	fmt.Println("Bolt                         :", a.bolt)
	fmt.Println("Uklon                        :", a.uklon)
	fmt.Println("Uber                         :", a.uber)
	fmt.Println("Всего                        :", a.trips)
	fmt.Println("Cash_Bolt                    :", float_to_string(a.cash_bolt))
	fmt.Println("Cash_Uklon                   :", float_to_string(a.cash_uklon))
	fmt.Println("Cash_Uber                    :", float_to_string(a.cash_uber))
	fmt.Println("Сумма:                       :", float_to_string(a.cash))

}

func (a *trips_Citys_days) Create(cit sity, b []Trips, t time.Time) {
	a.date = t.Format("02.01.2006")
	a.name = cit.id_BGQ
	a.name_Cyrrilic = cit.name
	for _, i := range b {
		if i.Date == t.Format("02.01.2006") && i.City == cit.id_BGQ {
			if i.Servise == "Bolt" {
				a.bolt = a.bolt + i.Trip
				a.cash_bolt = a.cash_bolt + i.Cash
			} else if i.Servise == "Uklon" {
				a.uklon = a.uklon + i.Trip
				a.cash_uklon = a.cash_uklon + i.Cash
			} else if i.Servise == "Uber" {
				a.uber = a.uber + i.Trip
				a.cash_uber = a.cash_uber + i.Cash
			}
		}
	}
	a.trips = a.bolt + a.uklon + a.uber
	a.cash = a.cash_bolt + a.cash_uklon + a.cash_uber
}

// создать запрос
// на получение поездок по
// списку городов и
// списку дат
func (a *gets_trips_dbase) CreateGet1() (s string) {
	s = "SELECT * FROM trips WHERE "
	s = s + "("
	for n, i := range a.city {
		if n == 0 {
			s = s + fmt.Sprintf("City = '%s'", i.id_BGQ)
		} else {
			s = s + fmt.Sprintf(" OR City = '%s'", i.id_BGQ)
		}
	}
	s = s + ") AND ("
	for n, i := range a.days {
		if n == 0 {
			s = s + fmt.Sprintf("Date = '%s'", i.Format("02.01.2006"))
		} else {
			s = s + fmt.Sprintf(" OR Date = '%s'", i.Format("02.01.2006"))
		}
	}
	s = s + ")"
	return
}

// создать запрос
// driver
// списку дат
func (a *gets_driver_trips_dbase) CreateGet() (s string) {
	s = fmt.Sprintf("SELECT * FROM trips WHERE Name = '%s' AND (", a.name)
	for n, i := range a.days {
		if n == 0 {
			s = s + fmt.Sprintf("Date = '%s'", i.Format("02.01.2006"))
		} else {
			s = s + fmt.Sprintf(" OR Date = '%s'", i.Format("02.01.2006"))
		}
	}
	s = s + ")"
	return
}

// получение поездок по
// списку городов и
// списку дат
func get_Trips(cit []sity, day []time.Time) (res []Trips) {
	x := gets_trips_dbase{}
	x.city = cit
	x.days = day
	x.x = make(chan []Trips)
	trips_dbase <- x
	r := <-x.x
	return r
}

// получение поездок по
// списку городов и
// списку дат
func get_Driver_Trips(n string, day []time.Time) (res []Trips) {
	x := gets_driver_trips_dbase{}
	x.name = n
	x.days = day
	x.x = make(chan []Trips)
	trips_driver_dbase <- x
	r := <-x.x
	return r
}

// получение поездок по
// графику
func get_Trips_Graf(g grafik) (res [2]int) {
	x := gets_trips_dbase_graf{}
	x.city = g.city
	x.days = g.data
	x.g = &g
	x.x = make(chan [2]int)
	trips_dbase_graf <- x
	r := <-x.x
	return r
}

func server_trips_dbase() {

	OS("Run ---  go server_trips_dbase()")

	db, err := sql.Open("sqlite3", path+"sqlite/osnova.db")
	if err != nil {
		fmt.Println("Запись", err)
	}
	defer db.Close()

	for {

		select {
		case msg := <-trips_driver_dbase:
			rows, err := db.Query(msg.CreateGet())
			fmt.Println(msg.CreateGet())
			if err != nil {
				errorLog.Println(err)
			}
			res := make([]Trips, 0)
			for rows.Next() {
				var p [7]string
				err = rows.Scan(&p[0], &p[1], &p[2], &p[3], &p[4], &p[5], &p[6])
				if err != nil {
					errorLog.Println(err)
				}
				r := Trips{}
				r.Create(p)
				r.Name = strings.ReplaceAll(r.Name, "|", "'")
				res = append(res, r)
			}

			msg.x <- res

		case msg := <-trips_dbase:

			rows, err := db.Query(msg.CreateGet1())

			if err != nil {
				errorLog.Println(err)
			}

			res := make([]Trips, 0)
			for rows.Next() {
				var p [7]string
				err = rows.Scan(&p[0], &p[1], &p[2], &p[3], &p[4], &p[5], &p[6])
				if err != nil {
					errorLog.Println(err)
				}
				r := Trips{}
				r.Create(p)
				r.Name = strings.ReplaceAll(r.Name, "|", "'")
				res = append(res, r)
			}

			msg.x <- res

		case msg := <-trips_dbase_graf:

			n1 := strings.ReplaceAll(msg.g.name_vod1, "(UK) ", "")
			n1 = strings.ReplaceAll(n1, "(UB) ", "")
			n2 := strings.ReplaceAll(msg.g.name_vod2, "(UK) ", "")
			n2 = strings.ReplaceAll(n2, "(UB) ", "")
			d := msg.g.data.Format("02.01.2006")
			s := ""
			if n1 == n2 {
				s = fmt.Sprintf("SELECT * FROM trips WHERE City = '%s' AND Date = '%s' AND Name = '%s'", msg.g.city, d, n1)
			} else {
				s = fmt.Sprintf("SELECT * FROM trips WHERE City = '%s' AND Date = '%s' AND (Name = '%s' OR Name = '%s')", msg.g.city, d, n1, n2)
			}
			//fmt.Println(s)
			rows, err := db.Query(s)

			if err != nil {
				errorLog.Println(err)
				fmt.Println(err)
				msg.x <- [2]int{0, 0}
				continue
			}

			trips := 0
			cash := 0
			for rows.Next() {
				var p [7]string
				err = rows.Scan(&p[0], &p[1], &p[2], &p[3], &p[4], &p[5], &p[6])
				if err != nil {
					errorLog.Println(err)
				}
				r := Trips{}
				r.Create(p)
				//r.Print()
				trips = trips + r.Trip
				cash = cash + int(r.Cash)
			}

			msg.x <- [2]int{trips, cash}

		}

		<-time.After(3 * time.Millisecond)
	}

}
