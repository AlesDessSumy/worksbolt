package main

import (
	"BGQ"
	"Citys"
	"database/sql"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
	"trs"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/civil"
)

// поиск поездок водителя за период
func Search_Driver_Trips_Period(n, d1, d2 string) {

	trs.Rec_Clear("1_EWAl5qFMMWh0xbr3gzJ_KM87QNXhm7wJy89Vrs_zDs", "'Driver'!E1", [][]string{{"У процесі"}}, false, false)

	n = trs.Clear_name_driver(n)
	dr := Drivers{}
	dr.Read_Accaunt(n)

	r := []Trips{}
	day := List_Day(d1, d2)

	if len(dr) == 1 {
		r = get_Driver_Trips(dr[0].Name, day)
	} else if len(dr) == 2 {
		if dr[0].Name == dr[1].Name {
			r = get_Driver_Trips(dr[0].Name, day)
		} else {
			r1 := get_Driver_Trips(dr[0].Name, day)
			r2 := get_Driver_Trips(dr[1].Name, day)
			r = append(r, r1...)
			r = append(r, r2...)
		}
	} else if len(dr) == 3 {
		if dr[0].Name == dr[1].Name && dr[0].Name == dr[2].Name {
			r = get_Driver_Trips(dr[0].Name, day)
		} else if dr[0].Name == dr[1].Name && dr[0].Name != dr[2].Name {
			r1 := get_Driver_Trips(dr[0].Name, day)
			r2 := get_Driver_Trips(dr[2].Name, day)
			r = append(r, r1...)
			r = append(r, r2...)
		} else if dr[0].Name == dr[2].Name && dr[0].Name != dr[1].Name {
			r1 := get_Driver_Trips(dr[0].Name, day)
			r2 := get_Driver_Trips(dr[1].Name, day)
			r = append(r, r1...)
			r = append(r, r2...)
		} else if dr[1].Name == dr[2].Name && dr[1].Name != dr[0].Name {
			r1 := get_Driver_Trips(dr[0].Name, day)
			r2 := get_Driver_Trips(dr[1].Name, day)
			r = append(r, r1...)
			r = append(r, r2...)
		} else if dr[0].Name != dr[1].Name && dr[0].Name != dr[2].Name && dr[1].Name != dr[2].Name {
			r1 := get_Driver_Trips(dr[0].Name, day)
			r2 := get_Driver_Trips(dr[1].Name, day)
			r3 := get_Driver_Trips(dr[2].Name, day)
			r = append(r, r1...)
			r = append(r, r2...)
			r = append(r, r3...)
		}
	}

	// сортируем по времени
	sort.Slice(r, func(i, j int) (less bool) {
		return r[j].Dat.After(r[i].Dat)
	})

	total := Trips{Name: n, City: "TOTAL"}

	rs := make([][]string, 0)
	res := make([][]string, 0)
	for _, i := range r {
		rs = append(rs, i.SliseString())
		total.Add(i)

	}
	res = append(res, total.SliseString())
	res = append(res, rs...)
	trs.Rec_Clear("1_EWAl5qFMMWh0xbr3gzJ_KM87QNXhm7wJy89Vrs_zDs", "'Driver'!A3", res, true, true)
	trs.Rec_Clear("1_EWAl5qFMMWh0xbr3gzJ_KM87QNXhm7wJy89Vrs_zDs", "'Driver'!E1", [][]string{{"Готово"}}, false, false)
}

func Check_Trips() {

	db, err := sql.Open("sqlite3", path+"sqlite/osnova.db")
	if err != nil {
		errorLog.Println(path+"sqlite/osnova.db", err)
	}
	defer db.Close()

	citys := Citys.Create_sitys_List()

	for _, cit := range citys {
		for _, dn := range last_days(28) {
			zap := fmt.Sprintf("Select * from trips where Date = '%s' and City = '%s'", dn.Format(TNS), cit.Id_BGQ)
			//fmt.Println(zap)
			rows, err := db.Query(zap)

			if err != nil {
				fmt.Println(err, zap)
			}
			trips := 0
			for rows.Next() {
				var p [7]string
				err = rows.Scan(&p[0], &p[1], &p[2], &p[3], &p[4], &p[5], &p[6])
				if err != nil {
					fmt.Println(err)
				}
				r := Trips{}
				r.Create(p)
				trips = trips + r.Trip

			}
			fmt.Println(dn.Format(TNS), cit.Name, trips)
		}
	}
}

func Put_Trip_DB_V2(period []time.Time, ignor_control bool) {

	OS(NS("Run ---  Put_Trip_DB_V2(last_days(1)) ", 40))

	control := control_execute{}
	control.data = time.Now().Format("02.01.2006 15-04")
	control.name = "Put_Trip_DB_V2()"
	control.coment = "completed"

	if control.read() {
		fmt.Println("Uclon Trip Read completed successfully")
		if !ignor_control {
			return
		}
	}

	db, err := sql.Open("sqlite3", path+"sqlite/osnova.db")
	if err != nil {
		errorLog.Println(path+"sqlite/osnova.db", err)
	}
	defer db.Close()

	del := "DELETE FROM trips"

	for n, tm := range period {
		if n == 0 {
			del = del + fmt.Sprintf(" WHERE Date = '%s'", tm.Format("02.01.2006"))
		} else {
			del = del + fmt.Sprintf(" OR Date = '%s'", tm.Format("02.01.2006"))
		}
	}

	fmt.Println(del)

	_, err = db.Exec(del)
	if err != nil {
		errorLog.Println("Ошибка удаления", del, err)
		fmt.Println(err)
	}

	trips := Get_Bgq_Uber(period)
	OS(fmt.Sprintf("Поездок Убер записано %d", len(trips)))

	for _, i := range trips {
		//	i.Print()
		name := strings.ReplaceAll(i.Name, "'", "|")
		zap := fmt.Sprintf("INSERT INTO trips (Name, Trip, Servise, Cash, Date, City) values ('%s', '%d','%s','%.2f','%s','%s')", name, i.Trip, i.Servise, i.Cash, i.Date, i.City)
		_, err := db.Exec(zap)
		if err != nil {
			errorLog.Println(err, zap)
		}
	}

	trips = Get_Bgq_Bolt(period)
	OS(fmt.Sprintf("Поездок Bolt записано %d", len(trips)))

	for _, i := range trips {
		//	i.Print()
		name := strings.ReplaceAll(i.Name, "'", "|")
		zap := fmt.Sprintf("INSERT INTO trips (Name, Trip, Servise, Cash, Date, City) values ('%s', '%d','%s','%.2f','%s','%s')", name, i.Trip, i.Servise, i.Cash, i.Date, i.City)
		_, err := db.Exec(zap)
		if err != nil {
			errorLog.Println(err, zap)
		}
	}

	trips = Get_Bgq_Uklon(period)
	OS(fmt.Sprintf("Поездок Uclon записано %d", len(trips)))

	for _, i := range trips {
		//i.Print()
		name := strings.ReplaceAll(i.Name, "'", "|")
		zap := fmt.Sprintf("INSERT INTO trips (Name, Trip, Servise, Cash, Date, City) values ('%s', '%d','%s','%.2f','%s','%s')", name, i.Trip, i.Servise, i.Cash, i.Date, i.City)
		_, err := db.Exec(zap)
		if err != nil {
			errorLog.Println(err, zap)
		}
	}

	zap := fmt.Sprintf("SELECT * FROM trips WHERE Date == '%s'", time.Now().AddDate(0, 0, -1).Format("02.01.2006"))
	rows, err := db.Query(zap)

	if err != nil {
		errorLog.Println(err, zap)
	}

	res := make([][]string, 0)
	for rows.Next() {
		var p [7]string
		err = rows.Scan(&p[0], &p[1], &p[2], &p[3], &p[4], &p[5], &p[6])
		if err != nil {
			errorLog.Println(err)
		}
		//fmt.Println(p)
		res = append(res, []string{p[5], p[6], p[3], p[1], p[2], strings.ReplaceAll(p[4], ".", ",")})
	}
	sort.Slice(res, func(i, j int) (less bool) {
		return res[i][2] < res[j][2]
	})
	sort.Slice(res, func(i, j int) (less bool) {
		return res[i][3] < res[j][3]
	})
	sort.Slice(res, func(i, j int) (less bool) {
		return res[i][1] < res[j][1]
	})
	trs.Rec("1boU2to5mT77THBlyIL1ZsCJLsqUKE393DcgJga3Z0_w", "'V2'!A2", res, false, true)

	control.save()
}

// для записи в базу
type Trips struct {
	Name    string
	Trip    int
	Bolt    int
	Uklon   int
	Uber    int
	Servise string
	Cash    float64
	Date    string
	City    string
	Dat     time.Time
}

func (a *Trips) Report() (res []string) {
	res = append(res, a.City)
	res = append(res, a.Name)
	res = append(res, fmt.Sprint(a.Trip))
	return
}

func (a *Trips) SliseString() (res []string) {
	res = append(res, a.City)
	res = append(res, a.Name)
	res = append(res, a.Date)
	res = append(res, a.Servise)
	res = append(res, fmt.Sprint(a.Trip))
	res = append(res, fmt.Sprint(a.Cash))
	res = append(res, fmt.Sprint(a.Bolt))
	res = append(res, fmt.Sprint(a.Uklon))
	res = append(res, fmt.Sprint(a.Uber))
	return
}

func (a *Trips) Create(b [7]string) {
	a.Name = b[1]
	a.Trip = string_to_int(b[2])
	a.Servise = b[3]
	a.Cash = string_to_float(b[4])
	a.Date = b[5]
	a.City = b[6]
	a.Dat, _ = time.Parse(TNS, a.Date)
}

// слияние двух элементов
func (a *Trips) Add(b Trips) {
	a.Trip = a.Trip + b.Trip
	a.Cash = a.Cash + b.Cash
	if b.Servise == "Bolt" {
		a.Bolt = a.Bolt + b.Trip
	} else if b.Servise == "Uklon" {
		a.Uklon = a.Uklon + b.Trip
	} else if b.Servise == "Uber" {
		a.Uber = a.Uber + b.Trip
	}
}

func (a *Trips) Print() {
	fmt.Println("Trips")
	fmt.Println("a.Name             :", a.Name)
	fmt.Println("a.Trip             :", a.Trip)
	fmt.Println("a.Bolt             :", a.Bolt)
	fmt.Println("a.Uklon            :", a.Uklon)
	fmt.Println("a.Uber             :", a.Uber)
	fmt.Println("a.Servise          :", a.Servise)
	fmt.Println("a.Cash             :", a.Cash)
	fmt.Println("a.Date             :", a.Date)
	fmt.Println("a.City             :", a.City)
	fmt.Println()
}

type list_Trips struct {
	s []Trips
}

// по имени, дате , городу отдает поездки
func (a *list_Trips) Calc(n1, n2, d, cit string) (r Trips) {
	if n1 == n2 {
		for _, i := range a.s {
			if i.Name == n1 && i.Date == d && cit == i.City {
				r.Add(i)
			}
		}
	} else {
		for _, i := range a.s {
			if (i.Name == n1 || i.Name == n2) && i.Date == d && cit == i.City {
				r.Add(i)
			}
		}
	}
	return
}

type Trip_uber struct {
	Name            string
	Payment_id      string
	Partner_id      string
	Driver_id       string
	Category        string
	Payment_amount  float64
	Cash_collected  float64
	Breakdown_other float64
	Riders_fee      float64
	Trip_id         string
	Currency_code   string
	Event_time      string
	Date            time.Time
	City            string
}

func (a *Trip_uber) Add(b Trip_uber) {
	a.Payment_amount = a.Payment_amount + b.Payment_amount
	a.Cash_collected = a.Cash_collected + b.Cash_collected
	a.Breakdown_other = a.Breakdown_other + b.Breakdown_other
	a.Riders_fee = a.Riders_fee + b.Riders_fee
	a.Date = b.Date
}

func (a *Trip_uber) CreateTrips() Trips {
	return Trips{Trip: 1, Name: trs.Clear_32(a.Name), Servise: "Uber", Cash: a.Payment_amount, Date: a.Date.Format("02.01.2006"), City: a.City}
}

func (a *Trip_bolt) CreateTrips() Trips {
	return Trips{Trip: 1, Name: trs.Clear_32(a.Drivers_name), Servise: "Bolt", Cash: a.The_price_of_the_trip, Date: a.Date.Format("02.01.2006"), City: a.City}
}

func (a *Trip_uklon_1) CreateTrips() Trips {
	return Trips{Trip: int(a.trips), Name: trs.Clear_32(a.first_name + " " + a.last_name), Servise: "Uklon", Cash: a.total_tariff, Date: a.date.Format("02.01.2006"), City: a.city}
}

func (a *Trip_uklon_2) CreateTrips() Trips {
	return Trips{Trip: int(a.Trips), Name: trs.Clear_32(a.First_Name), Servise: "Uklon", Cash: a.gross_income, Date: a.date.Format("02.01.2006"), City: a.city}
}

func (a *Trip_uber) Create(i []bigquery.Value) {
	if i[0] != nil {
		a.Payment_id = i[0].(string)
	}
	if i[1] != nil {
		a.Partner_id = i[1].(string)
	}
	if i[2] != nil {
		a.Driver_id = i[2].(string)
	}
	if i[3] != nil {
		a.Category = i[3].(string)
	}
	if i[4] != nil {
		a.Payment_amount = i[4].(float64)
	}
	if i[5] != nil {
		a.Cash_collected = i[5].(float64)
	}
	if i[6] != nil {
		a.Breakdown_other = i[6].(float64)
	}
	if i[7] != nil {
		a.Riders_fee = i[7].(float64)
	}
	if i[8] != nil {
		a.Trip_id = i[8].(string)
	}
	if i[9] != nil {
		a.Currency_code = i[9].(string)
	}
	if i[10] != nil {
		a.Event_time = i[10].(string)
	}
	if i[11] != nil {
		date := i[11].(civil.Date)
		a.Date = time.Date(date.Year, date.Month, date.Day, 7, 7, 7, 0, time.UTC)
	}
	if i[12] != nil {
		a.City = i[12].(string)
	}
}

func (a *Trip_uber) Print() {
	fmt.Println("Trip Uber")
	fmt.Println("a.Name             :", a.Name)
	fmt.Println("a.payment_id       :", a.Payment_id)
	fmt.Println("a.partner_id       :", a.Partner_id)
	fmt.Println("a.driver_id        :", a.Driver_id)
	fmt.Println("a.category         :", a.Category)
	fmt.Println("a.payment_amount   :", a.Payment_amount)
	fmt.Println("a.cash_collected   :", a.Cash_collected)
	fmt.Println("a.breakdown_other  :", a.Breakdown_other)
	fmt.Println("a.riders_fee       :", a.Riders_fee)
	fmt.Println("a.trip_id          :", a.Trip_id)
	fmt.Println("a.currency_code    :", a.Currency_code)
	fmt.Println("a.event_time       :", a.Event_time)
	fmt.Println("a.date             :", a.Date.Format("02.01.2006 15:04:05"))
	fmt.Println("a.city             :", a.City)
	fmt.Println()
}

// поездки по юбер
func Get_Bgq_Uber(period []time.Time) map[string]Trips {

	var drivers = make(map[string]BGQ.Drivers, 0)

	r := BGQ.Get_Drivers()
	for _, i := range r {
		drivers[i.Taxi_id] = i
	}

	dataset := "DB_External"
	table := "uber_payments_daily"

	rs := make(map[string]Trips, 0)

	s := fmt.Sprintf("SELECT * FROM %s.%s WHERE ", dataset, table)

	for n, tm := range period {
		if n == 0 {
			s = s + fmt.Sprintf(" date = '%s' ", tm.Format("2006-01-02"))
		} else {
			s = s + fmt.Sprintf(" OR date = '%s' ", tm.Format("2006-01-02"))
		}
	}

	list, _ := BGQ.Read_BGQ(s)

	for _, row := range list {
		//	fmt.Println(row)
		r := Trip_uber{}
		r.Create(row)
		dr := drivers[r.Driver_id]
		r.Name = dr.Name
		//r.Print()
		tr := r.CreateTrips()
		key := tr.Name + "|" + tr.Date + "|" + tr.City + "|" + tr.Servise
		temp, ok := rs[key]

		if ok {
			temp.Add(tr)
			rs[key] = temp
		} else {
			rs[key] = tr
		}
	}

	return rs

}

type Report_bolt struct {
	The_driver                          string
	Drivers_phone                       string
	Period                              string
	General_tariff                      float64
	Cancellation_fee                    float64
	Authorization_payment_payment       float64
	Authorization_payment_deductions    float64
	Additional_fee                      float64
	Bolt_Commission                     float64
	Cash_trips_cash_collected           float64
	Bolt_discount_amount_for_cash_trips float64
	Driver_bonus                        float64
	Compensation                        float64
	Refund                              float64
	Tip                                 float64
	Weekly_balance                      float64
	Online_h                            float64
	Utilisation                         float64
	City                                string
	Csv_date                            time.Time
}

func (a *Report_bolt) Add(b Report_bolt) {
	a.General_tariff = a.General_tariff + b.General_tariff
	a.Cancellation_fee = a.Cancellation_fee + b.Cancellation_fee
	a.Authorization_payment_payment = a.Authorization_payment_payment + b.Authorization_payment_payment
	a.Authorization_payment_deductions = a.Authorization_payment_deductions + b.Authorization_payment_deductions
	a.Additional_fee = a.Additional_fee + b.Additional_fee
	a.Bolt_Commission = a.Bolt_Commission + b.Bolt_Commission
	a.Cash_trips_cash_collected = a.Cash_trips_cash_collected + b.Cash_trips_cash_collected
	a.Bolt_discount_amount_for_cash_trips = a.Bolt_discount_amount_for_cash_trips + b.Bolt_discount_amount_for_cash_trips
	a.Driver_bonus = a.Driver_bonus + b.Driver_bonus
	a.Compensation = a.Compensation + b.Compensation
	a.Refund = a.Refund + b.Refund
	a.Tip = a.Tip + b.Tip
	a.Weekly_balance = a.Weekly_balance + b.Weekly_balance
	a.Online_h = a.Online_h + b.Online_h
	a.Utilisation = a.Utilisation + b.Utilisation
	a.Csv_date = b.Csv_date
	a.Period = b.Period
}

func (a *Report_bolt) Create(i []bigquery.Value) {
	if i[0] != nil {
		a.The_driver = i[0].(string)
	}
	if i[1] != nil {
		a.Drivers_phone = i[1].(string)
	}
	if i[2] != nil {
		a.Period = i[2].(string)
	}
	if i[3] != nil {
		a.General_tariff = i[3].(float64)
	}
	if i[4] != nil {
		a.Cancellation_fee = i[4].(float64)
	}
	if i[5] != nil {
		a.Authorization_payment_payment = i[5].(float64)
	}
	if i[6] != nil {
		a.Authorization_payment_deductions = i[6].(float64)
	}
	if i[7] != nil {
		a.Additional_fee = i[7].(float64)
	}
	if i[8] != nil {
		a.Bolt_Commission = i[8].(float64)
	}
	if i[9] != nil {
		a.Cash_trips_cash_collected = i[9].(float64)
	}
	if i[10] != nil {
		a.Bolt_discount_amount_for_cash_trips = i[10].(float64)
	}
	if i[11] != nil {
		a.Driver_bonus = i[11].(float64)
	}
	if i[12] != nil {
		a.Compensation = i[12].(float64)
	}
	if i[13] != nil {
		a.Refund = i[13].(float64)
	}
	if i[14] != nil {
		a.Tip = i[14].(float64)
	}
	if i[15] != nil {
		a.Weekly_balance = i[15].(float64)
	}
	if i[16] != nil {
		a.Online_h = i[16].(float64)
	}
	if i[17] != nil {
		a.Utilisation = i[17].(float64)
	}
	if i[18] != nil {
		a.City = i[18].(string)
	}
	if i[19] != nil {
		date := i[19].(civil.Date)
		a.Csv_date = time.Date(date.Year, date.Month, date.Day, 7, 7, 7, 0, time.UTC)
	}
}

func (a *Report_bolt) Print() {
	fmt.Println()
	fmt.Println("Trip_bolt                             :")
	fmt.Println("The_driver                            :", a.The_driver)
	fmt.Println("Drivers_phone                         :", a.Drivers_phone)
	fmt.Println("Period                                :", a.Period)
	fmt.Println("General_tariff                        :", float_to_string(a.General_tariff))
	fmt.Println("Cancellation_fee                      :", float_to_string(a.Cancellation_fee))
	fmt.Println("Authorization_payment_payment         :", float_to_string(a.Authorization_payment_payment))
	fmt.Println("Authorization_payment_deductions      :", float_to_string(a.Authorization_payment_deductions))
	fmt.Println("Additional_fee                        :", float_to_string(a.Additional_fee))
	fmt.Println("Bolt_Commission                       :", float_to_string(a.Bolt_Commission))
	fmt.Println("Cash_trips_cash_collected             :", float_to_string(a.Cash_trips_cash_collected))
	fmt.Println("Bolt_discount_amount_for_cash_trips   :", float_to_string(a.Bolt_discount_amount_for_cash_trips))
	fmt.Println("Driver_bonus                          :", float_to_string(a.Driver_bonus))
	fmt.Println("Compensation                          :", float_to_string(a.Compensation))
	fmt.Println("Refund                                :", float_to_string(a.Refund))
	fmt.Println("Tip                                   :", float_to_string(a.Tip))
	fmt.Println("Weekly_balance                        :", float_to_string(a.Weekly_balance))
	fmt.Println("Online_h                              :", float_to_string(a.Online_h))
	fmt.Println("Utilisation                           :", float_to_string(a.Utilisation))
	fmt.Println("City                                  :", a.City)
	fmt.Println("Csv_date                              :", a.Csv_date.Format(TNS))
}

type Trip_bolt struct {
	Drivers_name          string
	Drivers_phone         string
	Date                  time.Time
	Payment_confirmed     string
	Landing               string
	Payment_method        string
	Asked                 string
	The_price_of_the_trip float64
	Authorization_payment float64
	Additional_fee        float64
	Cancellation_fee      float64
	Tip                   float64
	Order_Status          string
	Car_Model             string
	Car_Reg_Number        string
	City                  string
	Csv_date              time.Time
}

func (a *Trip_bolt) Add(b Trip_bolt) {
	a.The_price_of_the_trip = a.The_price_of_the_trip + b.The_price_of_the_trip
	a.Authorization_payment = a.Authorization_payment + b.Authorization_payment
	a.Additional_fee = a.Additional_fee + b.Additional_fee
	a.Cancellation_fee = a.Cancellation_fee + b.Cancellation_fee
	a.Tip = a.Tip + b.Tip
	a.Date = b.Date
}

func (a *Trip_bolt) Create(i []bigquery.Value) {
	if i[0] != nil {
		a.Drivers_name = i[0].(string)
	}
	if i[1] != nil {
		a.Drivers_phone = i[1].(string)
	}
	if i[2] != nil {
		date := i[2].(civil.Date)
		a.Date = time.Date(date.Year, date.Month, date.Day, 7, 7, 7, 0, time.UTC)
	}
	if i[3] != nil {
		a.Payment_confirmed = i[3].(string)
	}
	if i[4] != nil {
		a.Landing = i[4].(string)
	}
	if i[5] != nil {
		a.Payment_method = i[5].(string)
	}
	if i[6] != nil {
		a.Asked = i[6].(string)
	}
	if i[7] != nil {
		a.The_price_of_the_trip = i[7].(float64)
	}
	if i[8] != nil {
		a.Authorization_payment = i[8].(float64)
	}
	if i[9] != nil {
		a.Additional_fee = i[9].(float64)
	}
	if i[10] != nil {
		a.Cancellation_fee = i[10].(float64)
	}
	if i[11] != nil {
		a.Tip = i[11].(float64)
	}
	if i[12] != nil {
		a.Order_Status = i[12].(string)
	}
	if i[13] != nil {
		a.Car_Model = i[13].(string)
	}
	if i[14] != nil {
		a.Car_Reg_Number = i[14].(string)
	}
	if i[15] != nil {
		a.City = i[15].(string)
	}
	if i[16] != nil {
		date := i[16].(civil.Date)
		a.Date = time.Date(date.Year, date.Month, date.Day, 7, 7, 7, 0, time.UTC)
	}
}

func (a *Trip_bolt) Print() {
	fmt.Println()
	fmt.Println("Trip_bolt        :")
	fmt.Println("Drivers_name :", a.Drivers_name)
	fmt.Println("Drivers_phone :", a.Drivers_phone)
	fmt.Println("Date :", a.Date)
	fmt.Println("Payment_confirmed :", a.Payment_confirmed)
	fmt.Println("Landing :", a.Landing)
	fmt.Println("Payment_method :", a.Payment_method)
	fmt.Println("Asked :", a.Asked)
	fmt.Println("The_price_of_the_trip :", a.The_price_of_the_trip)
	fmt.Println("Authorization_payment :", a.Authorization_payment)
	fmt.Println("Additional_fee :", a.Additional_fee)
	fmt.Println("Cancellation_fee :", a.Cancellation_fee)
	fmt.Println("Tip :", a.Tip)
	fmt.Println("Order_Status :", a.Order_Status)
	fmt.Println("Car_Model :", a.Car_Model)
	fmt.Println("Car_Reg_Number :", a.Car_Reg_Number)
	fmt.Println("City :", a.City)
}

// записывает поездки по БОЛТ
// списку городов за список дат
func Get_Bgq_Bolt(period []time.Time) map[string]Trips {

	rs := make(map[string]Trips, 0)

	for _, cit := range create_sitys() {

		if cit.bolt_trip.name_dataset == "" || cit.bolt_trip.name_table == "" {
			continue
		}

		s := fmt.Sprintf("SELECT * FROM %s.%s WHERE ", cit.bolt_trip.name_dataset, cit.bolt_trip.name_table)

		for n, tm := range period {
			if n == 0 {
				s = s + fmt.Sprintf(" Csv_date = '%s' ", tm.Format("2006-01-02"))
			} else {
				s = s + fmt.Sprintf(" OR Csv_date = '%s' ", tm.Format("2006-01-02"))
			}
		}

		//fmt.Println(s)

		list, _ := BGQ.Read_BGQ(s)

		for _, row := range list {
			r := Trip_bolt{}
			r.Create(row)
			//r.Print()
			tr := r.CreateTrips()
			key := tr.Name + "|" + tr.Date + "|" + tr.City + "|" + tr.Servise
			temp, ok := rs[key]

			if ok {
				temp.Add(tr)
				rs[key] = temp
			} else {
				rs[key] = tr
			}

		}
	}
	return rs
}

type Trip_uklon_1 struct {
	id            string
	first_name    string
	last_name     string
	signal        string
	trips         int64
	distance      float64
	total_tariff  float64
	fee           float64
	cash          float64
	wallet        float64
	wallet_fee    float64
	total_balance float64
	tips          float64
	city          string
	date          time.Time
}

func (a *Trip_uklon_1) Add(b Trip_uklon_1) {
	a.trips = a.trips + b.trips
	a.distance = a.distance + b.distance
	a.total_tariff = a.total_tariff + b.total_tariff
	a.fee = a.fee + b.fee
	a.cash = a.cash + b.cash
	a.wallet = a.wallet + b.wallet
	a.wallet_fee = a.wallet_fee + b.wallet_fee
	a.total_balance = a.total_balance + b.total_balance
	a.tips = a.tips + b.tips
	a.date = b.date
}

func (a *Trip_uklon_1) Create_II(i []bigquery.Value) {
	if i[0] != nil {
		a.id = i[0].(string)
	}
	if i[1] != nil {
		n := i[1].(string)
		nn := strings.Split(n, " ")
		if len(nn) == 2 {
			a.first_name = nn[0]
			a.last_name = nn[1]
		} else if len(nn) == 3 {
			a.first_name = nn[0]
			a.last_name = nn[1] + " " + nn[2]
		} else if len(nn) == 1 {
			a.first_name = nn[0]
		}

	}
	if i[2] != nil {
		a.signal = fmt.Sprint(i[2].(int64))
	}

	if i[4] != nil {
		a.trips = i[4].(int64)
	}
	if i[5] != nil {
		a.distance = i[5].(float64)
	}
	if i[6] != nil {
		a.total_tariff = i[6].(float64)
	}
	if i[7] != nil {
		a.fee = -i[7].(float64)
	}
	if i[8] != nil {
		a.cash = i[8].(float64)
	}
	// if i[9] != nil {
	// 	a.gross_amount_cashless = i[9].(float64)
	// }
	// if i[10] != nil {
	// 	a.profit = i[10].(float64)
	// }
	// if i[11] != nil {
	// 	a.Driver_bonus = i[11].(float64)
	// }
	if i[12] != nil {
		a.city = i[12].(string)
	}
	if i[13] != nil {
		date := i[13].(civil.Date)
		a.date = time.Date(date.Year, date.Month, date.Day, 7, 7, 7, 0, time.UTC)
	}
}

func (a *Trip_uklon_1) Create(i []bigquery.Value) {
	if i[0] != nil {
		a.id = i[0].(string)
	}
	if i[1] != nil {
		a.first_name = i[1].(string)
	}
	if i[2] != nil {
		a.last_name = i[2].(string)
	}
	if i[3] != nil {
		a.signal = i[3].(string)
	}
	if i[4] != nil {
		a.trips = i[4].(int64)
	}
	if i[5] != nil {
		a.distance = i[5].(float64)
	}
	if i[6] != nil {
		a.total_tariff = i[6].(float64)
	}
	if i[7] != nil {
		a.fee = i[7].(float64)
	}
	if i[8] != nil {
		a.cash = i[8].(float64)
	}
	if i[9] != nil {
		a.wallet = i[9].(float64)
	}
	if i[10] != nil {
		a.wallet_fee = i[10].(float64)
	}
	if i[11] != nil {
		a.total_balance = i[11].(float64)
	}
	if i[12] != nil {
		a.tips = i[12].(float64)
	}
	if i[13] != nil {
		a.city = i[13].(string)
	}
	if i[14] != nil {
		date := i[14].(civil.Date)
		a.date = time.Date(date.Year, date.Month, date.Day, 7, 7, 7, 0, time.UTC)
	}
}

func (a *Trip_uklon_1) Print() {
	fmt.Println()
	fmt.Println("Trip_uklon_1   :")
	fmt.Println("id             :", a.id)
	fmt.Println("name           :", a.first_name+" "+a.last_name)
	fmt.Println("signal         :", a.signal)
	fmt.Println("trips          :", a.trips)
	fmt.Println("distance       :", a.distance)
	fmt.Println("total_tariff   :", a.total_tariff)
	fmt.Println("fee            :", a.fee)
	fmt.Println("cash           :", a.cash)
	fmt.Println("wallet         :", a.wallet)
	fmt.Println("wallet_fee     :", a.wallet_fee)
	fmt.Println("total_balance  :", a.total_balance)
	fmt.Println("tips           :", a.tips)
	fmt.Println("city           :", a.city)
	fmt.Println("date           :", a.date)
}

type Trip_uklon_2 struct {
	driver_id             string
	First_Name            string
	signal                int64
	Car_Reg_Number        string
	Trips                 int64
	distance              float64
	gross_income          float64
	fee_income            float64
	Cash_Colected         float64
	gross_amount_cashless float64
	profit                float64
	Driver_bonus          float64
	city                  string
	date                  time.Time
}

func (a *Trip_uklon_2) Create(i []bigquery.Value) {
	if i[0] != nil {
		a.driver_id = i[0].(string)
	}
	if i[1] != nil {
		a.First_Name = i[1].(string)
	}
	if i[2] != nil {
		a.signal = i[2].(int64)
	}
	if i[3] != nil {
		a.Car_Reg_Number = i[3].(string)
	}
	if i[4] != nil {
		a.Trips = i[4].(int64)
	}
	if i[5] != nil {
		a.distance = i[5].(float64)
	}
	if i[6] != nil {
		a.gross_income = i[6].(float64)
	}
	if i[7] != nil {
		a.fee_income = i[7].(float64)
	}
	if i[8] != nil {
		a.Cash_Colected = i[8].(float64)
	}
	if i[9] != nil {
		a.gross_amount_cashless = i[9].(float64)
	}
	if i[10] != nil {
		a.profit = i[10].(float64)
	}
	if i[11] != nil {
		a.Driver_bonus = i[11].(float64)
	}
	if i[12] != nil {
		a.city = i[12].(string)
	}
	if i[13] != nil {
		date := i[13].(civil.Date)
		a.date = time.Date(date.Year, date.Month, date.Day, 7, 7, 7, 0, time.UTC)
	}
}

func (a *Trip_uklon_2) Print() {
	fmt.Println()
	fmt.Println("Trip_uklon_2            :")
	fmt.Println("driver_id               :", a.driver_id)
	fmt.Println("First_Name              :", a.First_Name)
	fmt.Println("signal                  :", a.signal)
	fmt.Println("Car_Reg_Number          :", a.Car_Reg_Number)
	fmt.Println("Trips                   :", a.Trips)
	fmt.Println("distance                :", a.distance)
	fmt.Println("gross_income            :", a.gross_income)
	fmt.Println("fee_income              :", a.fee_income)
	fmt.Println("Cash_Colected           :", a.Cash_Colected)
	fmt.Println("gross_amount_cashless   :", a.gross_amount_cashless)
	fmt.Println("profit                  :", a.profit)
	fmt.Println("Driver_bonus            :", a.Driver_bonus)
	fmt.Println("city                    :", a.city)
	fmt.Println("date                    :", a.date)
}

// записывает поездки по UKLON
// списку городов за список дат
func Get_Bgq_Uklon(period []time.Time) map[string]Trips {

	rs := make(map[string]Trips)

	s := "SELECT * FROM DB_External.uklon_report_orders  WHERE "
	for n, tm := range period {
		if n == 0 {
			s = s + fmt.Sprintf(" date = '%s' ", tm.Format("2006-01-02"))
		} else {
			s = s + fmt.Sprintf(" OR date = '%s' ", tm.Format("2006-01-02"))
		}
	}
	//fmt.Println(s)

	list, _ := BGQ.Read_BGQ(s)

	for _, row := range list {
		r := Trip_uklon_1{}
		r.Create(row)

		//r.Print()
		tr := r.CreateTrips()

		key := tr.Name + "|" + tr.Date + "|" + tr.City + "|" + tr.Servise
		temp, ok := rs[key]

		if ok {
			temp.Add(tr)
			rs[key] = temp
		} else {
			rs[key] = tr
		}
	}

	for _, cit := range create_sitys() {

		if cit.uclon_cash.name_dataset == "" || cit.uclon_cash.name_table == "" {
			continue
		}

		s := fmt.Sprintf("SELECT * FROM %s.%s WHERE ", cit.uclon_cash.name_dataset, cit.uclon_cash.name_table)
		for n, tm := range period {
			if n == 0 {
				s = s + fmt.Sprintf(" date = '%s'", tm.Format("2006-01-02"))
			} else {
				s = s + fmt.Sprintf(" OR date = '%s' ", tm.Format("2006-01-02"))
			}
		}

		list, _ := BGQ.Read_BGQ(s)

		for _, row := range list {
			r := Trip_uklon_2{}
			r.Create(row)
			tr := r.CreateTrips()
			if tr.Name == "" {
				continue
			}
			key := tr.Name + "|" + tr.Date + "|" + tr.City + "|" + tr.Servise
			temp, ok := rs[key]

			if ok {
				temp.Add(tr)
				rs[key] = temp
			} else {
				rs[key] = tr
			}
		}
	}
	return rs
}

// записывает поездки по БОЛТ
// списку городов за список дат
func get_bgq_bolt(sitys []sity, days []time.Time, ignor_control bool) {

	OS(NS("Run ---  get_bgq_bolt(sitys, day, false) ", 40))

	control := control_execute{}
	control.data = time.Now().Format("02.01.2006 15-04")
	control.name = "get_bgq_bolt"
	control.coment = "completed"

	if control.read() {
		fmt.Println("Bolt Trip Read completed successfully")
		if !ignor_control {
			return
		}
	}

	rec_sity_day := func(sit sity, dayT time.Time) (res_trip [][]string) {

		dt := make([]string, 0)

		day := dayT.Format("2006-01-02")

		s := "SELECT * FROM " + sit.bolt_trip.name_dataset + "." + sit.bolt_trip.name_table + " WHERE Csv_date = " + "\"" + day + "\""

		list, _ := BGQ.Read_BGQ(s)

		res := make([]element, 0)

		for _, row := range list {

			if row[0] == nil {
				continue
			}
			var rs element
			rs.sity = sit.bolt_trip.name_dataset
			rs.trip = 1
			rs.name = row[0].(string)
			dat := row[2].(civil.Date)
			rs.data = dat.String()
			rs.data = fmt.Sprintf("%02d.%02d.%04d", dat.Day, dat.Month, dat.Year)
			rs.time = row[3].(string)
			rs.pay = int(row[7].(float64))
			fl := 0
			for _, chek := range res {
				if chek == rs {
					fl = 1
					break
				}
			}
			if fl == 0 {
				res = append(res, rs)
				if !comp_string(rs.data, dt) {
					dt = append(dt, rs.data)
				}
			}
		}

		// получаем занятость за день
		imp := get_inployment(sit, dayT.Format("02.01.2006"))

		create_vod := func(a vod_day) (n, m string) {
			res := ""
			res = strconv.Itoa(a.trip_sum) + ";" + strconv.Itoa(a.cash_summ)
			res = res + ";"
			for _, i := range a.trip {
				res = res + strconv.Itoa(i) + ","
			}
			res = res + ";"
			for _, i := range a.cash {
				res = res + strconv.Itoa(i) + ","
			}
			return a.name, res
		}

		create_day := func(a []element) [][4]string {
			res_cd := make([][4]string, 0)
			//создать список водителей за день
			vod := make([]string, 0)
			for _, i := range a {
				if !comp_string(i.name, vod) {
					vod = append(vod, i.name)
				}
			}
			// просчет водителей
			for _, i := range vod {
				res_vd := vod_day{}
				res_vd.name = i
				res_vd.data = a[0].data

				for _, j := range a {
					if i == j.name {
						res_vd.cash_summ = res_vd.cash_summ + j.pay
						res_vd.trip_sum++
						hour := string([]byte(j.time)[:2])
						switch hour {
						case "00":
							res_vd.trip[0]++
							res_vd.cash[0] = res_vd.cash[0] + j.pay
						case "01":
							res_vd.trip[1]++
							res_vd.cash[1] = res_vd.cash[1] + j.pay
						case "02":
							res_vd.trip[2]++
							res_vd.cash[2] = res_vd.cash[2] + j.pay
						case "03":
							res_vd.trip[3]++
							res_vd.cash[3] = res_vd.cash[3] + j.pay
						case "04":
							res_vd.trip[4]++
							res_vd.cash[4] = res_vd.cash[4] + j.pay
						case "05":
							res_vd.trip[5]++
							res_vd.cash[5] = res_vd.cash[5] + j.pay
						case "06":
							res_vd.trip[6]++
							res_vd.cash[6] = res_vd.cash[6] + j.pay
						case "07":
							res_vd.trip[7]++
							res_vd.cash[7] = res_vd.cash[7] + j.pay
						case "08":
							res_vd.trip[8]++
							res_vd.cash[8] = res_vd.cash[8] + j.pay
						case "09":
							res_vd.trip[9]++
							res_vd.cash[9] = res_vd.cash[9] + j.pay
						case "10":
							res_vd.trip[10]++
							res_vd.cash[10] = res_vd.cash[10] + j.pay
						case "11":
							res_vd.trip[11]++
							res_vd.cash[11] = res_vd.cash[11] + j.pay
						case "12":
							res_vd.trip[12]++
							res_vd.cash[12] = res_vd.cash[12] + j.pay
						case "13":
							res_vd.trip[13]++
							res_vd.cash[13] = res_vd.cash[13] + j.pay
						case "14":
							res_vd.trip[14]++
							res_vd.cash[14] = res_vd.cash[14] + j.pay
						case "15":
							res_vd.trip[15]++
							res_vd.cash[15] = res_vd.cash[15] + j.pay
						case "16":
							res_vd.trip[16]++
							res_vd.cash[16] = res_vd.cash[16] + j.pay
						case "17":
							res_vd.trip[17]++
							res_vd.cash[17] = res_vd.cash[17] + j.pay
						case "18":
							res_vd.trip[18]++
							res_vd.cash[18] = res_vd.cash[18] + j.pay
						case "19":
							res_vd.trip[19]++
							res_vd.cash[19] = res_vd.cash[19] + j.pay
						case "20":
							res_vd.trip[20]++
							res_vd.cash[20] = res_vd.cash[20] + j.pay
						case "21":
							res_vd.trip[21]++
							res_vd.cash[21] = res_vd.cash[21] + j.pay
						case "22":
							res_vd.trip[22]++
							res_vd.cash[22] = res_vd.cash[22] + j.pay
						case "23":
							res_vd.trip[23]++
							res_vd.cash[23] = res_vd.cash[23] + j.pay
						}
					}
				}
				res_v := [4]string{}
				res_v[0] = res_vd.data
				res_v[1], res_v[2] = create_vod(res_vd)
				for _, im := range imp {
					if im.name == i {
						res_v[3] = im.data_s
					}
				}
				res_cd = append(res_cd, res_v)
				res_trip = append(res_trip, []string{day, sit.name, res_vd.name, fmt.Sprint(res_vd.trip_sum), fmt.Sprint(res_vd.cash_summ)})
			}
			return res_cd
		}

		for _, i := range dt {
			if i == time.Now().Format("02.01.2006") {
				continue
			}

			m := 0

			day := make([]element, 0)
			for _, j := range res {
				if i == j.data {
					day = append(day, j)
				}
			}

			d := create_day(day)

			db, err := sql.Open("sqlite3", path+"sqlite/trip.db")
			if err != nil {
				fmt.Println("Запись", err)
			}
			defer db.Close()

			// проверить есть ли сегодняшний день в базе
			t, _ := time.Parse("02.01.2006", i)
			year, ned := t.ISOWeek()

			name_table := sit.ident + "_ned_" + strconv.Itoa(ned) + "_" + strconv.Itoa(year)
			zapr := "SELECT count(*) FROM sqlite_master WHERE type='table' AND name='" + name_table + "';"
			//fmt.Println(zapr)
			tables := db.QueryRow(zapr)
			nn := 0
			tables.Scan(&nn)

			zapr1 := "CREATE TABLE " + name_table + " (id INTEGER PRIMARY KEY AUTOINCREMENT ,data VARCHAR(10), name NVARCHAR(70), dbase VARCHAR(135), imp VARCHAR(155), type VARCHAR(5));"
			//fmt.Println(zapr1)
			if nn == 0 {

				for {
					_, err := db.Exec(zapr1)
					if err != nil {
						fmt.Println("Запись поездок за ", t, " невозможно")
						fmt.Println("Functhion saveBD_Trip, file trip.go, str 346")
						time.Sleep(3 * time.Second)
					} else {
						break
					}
				}
			}

			for _, dd := range d {
				//fmt.Println(dd)

				name := string(code([]byte(dd[1])))

				zap := "SELECT * FROM " + name_table + " WHERE data == " + "'" + dd[0] + "'" + " AND name = " + "'" + name + "'" + " AND type = " + "'bolt'"
				//fmt.Println(zap)
				pp := db.QueryRow(zap)

				var p [6]string
				err = pp.Scan(&p[0], &p[1], &p[2], &p[3], &p[4], &p[5])
				//fmt.Println(p)
				if err == nil {
					del := "DELETE FROM " + name_table + " WHERE data == " + "'" + dd[0] + "'" + " AND name = " + "'" + name + "'" + " AND type = " + "'bolt'"
					_, err := db.Exec(del)
					if err != nil {
						fmt.Println("Ошибка удаления", err)
					}
				}
				m++

				zapr := "insert into " + name_table + " (data, name, dbase, imp, type) values ("
				zapr = zapr + "'" + dd[0] + "', '" + name + "', '" + dd[2] + "', '" + dd[3] + "', 'bolt')"

				//fmt.Println(zapr)

				for {
					_, err := db.Exec(zapr)
					if err != nil {
						errorLog.Println(err)
						time.Sleep(3 * time.Second)
					} else {
						break
					}
				}
			}

			send := fmt.Sprint("Записываем BOLT  ", sit.norma_name(), ",  день: ", i, "- Сделано ", m, " записей")
			fmt.Println(send)

		}

		return res_trip

	}

	for _, dayT := range days {

		res := make([][]string, 0)

		for _, sit := range sitys {

			res = append(res, []string{"--", "--", "--", "--", "--"})

			if sit.bolt_trip.name_dataset == "" {
				continue
			}
			if sit.bolt_trip.name_table == "" {
				continue
			}
			r := rec_sity_day(sit, dayT)
			res = append(res, r...)

		}
		trs.Rec("1boU2to5mT77THBlyIL1ZsCJLsqUKE393DcgJga3Z0_w", "'bolt'!A2", res, false, true)
	}

	control.save()

}

// записывает поездки по UKLON
// списку городов за список дат
func get_bgq_uclon_V2(days []time.Time, ignor_control bool) {

	OS(NS("Run ---  get_bgq_uclon_V2(day, false) ", 40))

	control := control_execute{}
	control.data = time.Now().Format("02.01.2006 15-04")
	control.name = "get_bgq_uclon"
	control.coment = "completed"

	if control.read() {
		fmt.Println("Uclon Trip Read completed successfully")
		if !ignor_control {
			return
		}
	}

	rec_days := func(t time.Time) {

		res_trip := make([][]string, 0)

		rs_allsity := get_bgq_uclon_all_city([]time.Time{t})

		day := t.Format("2006-01-02")

		var wg sync.WaitGroup

		summa_day := func(sit sity, day string) {
			defer wg.Done()
			// получаем сумму за день
			s := fmt.Sprintf("SELECT * FROM %s.%s WHERE date = '%s'", sit.uclon_cash.name_dataset, sit.uclon_cash.name_table, day)
			//fmt.Println(s)

			list, _ := BGQ.Read_BGQ(s)

			for _, row := range list {
				//fmt.Println("**", sit.id_BGQ, row)

				if row[1] == nil || row[3] == nil || row[4] == nil || row[5] == nil || row[6] == nil || row[13] == nil { //
					continue
				}
				//fmt.Println("--", sit.id_BGQ, row)
				var rs uclon_sum
				rs.city = sit.id_BGQ
				rs.name = row[1].(string)
				rs.name = trs.Clear_32(rs.name)
				rs.car = row[3].(string)
				rs.trip = row[4].(int64)
				rs.travel = row[5].(float64)
				rs.cash = row[6].(float64)
				dat := row[13].(civil.Date)
				rs.data = dat.String()
				rs.data = fmt.Sprintf("%02d.%02d.%04d", dat.Day, dat.Month, dat.Year)

				// if rs.city == "Kryvyi" { // || strings.Contains(rs.name, "Дмитро Аршулік")
				// 	fmt.Println("*", sit.name)
				// 	fmt.Println(row)
				// 	fmt.Println(rs)
				// 	fmt.Println()
				// }
				// fmt.Println(rs)
				rs_allsity = append(rs_allsity, rs)
			}
		}

		// for _, i := range rs_allsity {
		// 	if i.city == "Kryvyi" {
		// 		fmt.Println(i)
		// 	}
		// }

		for _, sit := range create_sitys() {

			if sit.uclon_cash.name_dataset == "" || sit.uclon_cash.name_table == "" {
				continue
			}
			wg.Add(1)
			go summa_day(sit, day)
			<-time.After(50 * time.Millisecond)
		}

		wg.Wait()

		create_vod := func(a vod_day) (n, m string) {
			res := ""
			res = strconv.Itoa(a.trip_sum) + ";" + strconv.Itoa(a.cash_summ)
			res = res + ";"
			for _, i := range a.trip {
				res = res + strconv.Itoa(i) + ","
			}
			res = res + ";"
			for _, i := range a.cash {
				res = res + strconv.Itoa(i) + ","
			}
			return a.name, res
		}

		create_day := func(a, cit string) (res_cd [][4]string, trip, cash int) {

			res_trip = append(res_trip, []string{"--", "--", "--", "--", "--"})
			// просчет водителей
			//fmt.Println(a)
			temp := make(map[string]uclon_sum, 0)

			for _, i := range rs_allsity {

				if i.city == cit {
					name := "(UK) " + i.name
					_, ok := temp[name]
					if !ok {
						temp[name] = i
					} else {
						//fmt.Println("ok", temp[name])
						r := temp[name]
						r.cash = r.cash + i.cash
						r.trip = r.trip + i.trip
						r.travel = r.travel + i.travel
						temp[name] = r
						// fmt.Println("ok", temp[name])
						// fmt.Println()
					}
				}

			}

			for _, i := range temp {
				res_vd := vod_day{}
				res_vd.name = "(UK) " + i.name
				res_vd.data = i.data
				res_vd.trip_sum = int(i.trip)
				res_vd.cash_summ = int(i.cash)
				res_v := [4]string{}
				res_v[0] = a
				res_v[1], res_v[2] = create_vod(res_vd)
				res_cd = append(res_cd, res_v)
				trip = trip + int(i.trip)
				cash = cash + int(i.cash)
				res_trip = append(res_trip, []string{res_vd.data, psevdonimy[cit], res_vd.name, fmt.Sprint(res_vd.trip_sum), fmt.Sprint(res_vd.cash_summ)})

			}
			return
		}

		for _, sit := range create_sitys() {

			m := 0

			d, trip, _ := create_day(day, sit.id_BGQ)

			db, err := sql.Open("sqlite3", path+"sqlite/trip.db")
			if err != nil {
				fmt.Println("Запись", err)
			}
			defer db.Close()

			// проверить есть ли сегодняшний день в базе
			t, _ := time.Parse("2006-01-02", day)
			year, ned := t.ISOWeek()

			name_table := sit.ident + "_ned_" + strconv.Itoa(ned) + "_" + strconv.Itoa(year)
			//fmt.Println(name_table)
			zapr := "SELECT count(*) FROM sqlite_master WHERE type='table' AND name='" + name_table + "';"
			//fmt.Println(zapr)
			tables := db.QueryRow(zapr)
			nn := 0
			tables.Scan(&nn)

			zapr1 := "CREATE TABLE " + name_table + " (id INTEGER PRIMARY KEY AUTOINCREMENT, data VARCHAR(10), name NVARCHAR(70), dbase VARCHAR(135), imp VARCHAR(155), type VARCHAR(5));"
			//fmt.Println(zapr1)
			if nn == 0 {
				for {
					_, err := db.Exec(zapr1)
					if err != nil {
						fmt.Println("Запись поездок за", t, "невозможно")
						fmt.Println("Functhion saveBD_Trip, file trip.go, str 1156")
						time.Sleep(3 * time.Second)
					} else {
						break
					}
				}
			}

			// удаляем день перед записью
			del := fmt.Sprintf("DELETE FROM %s WHERE data = '%s'", name_table, t.Format("02.01.2006"))
			//fmt.Println(del)
			_, err = db.Exec(del)
			if err != nil {
				fmt.Println("Ошибка удаления", err)
			}

			for _, dd := range d {

				tm, _ := time.Parse("2006-01-02", dd[0])
				data_corr := tm.Format("02.01.2006")
				//	fmt.Println(data_corr, dd)
				name := string(code([]byte(dd[1])))

				zap := "SELECT * FROM " + name_table + " WHERE data == " + "'" + dd[0] + "'" + " AND name = " + "'" + name + "'" + " AND type = " + "'uclon'"
				pp := db.QueryRow(zap)

				var p [6]string
				err = pp.Scan(&p[0], &p[1], &p[2], &p[3], &p[4], &p[5])
				if err == nil {
					del := "DELETE FROM " + name_table + " WHERE data == " + "'" + dd[0] + "'" + " AND name = " + "'" + name + "'" + " AND type = " + "'uclon'"
					_, err := db.Exec(del)
					if err != nil {
						fmt.Println("Ошибка удаления", err)
					}
				}
				m++

				zapr := "insert into " + name_table + " (data, name, dbase, imp, type) values ("
				zapr = zapr + "'" + data_corr + "', '" + name + "', '" + dd[2] + "', '" + dd[3] + "', 'uclon')"
				//fmt.Println(name_table, dd)

				//fmt.Println(zapr)

				for {
					_, err := db.Exec(zapr)
					if err != nil {
						errorLog.Println(err)
						time.Sleep(3 * time.Second)
					} else {
						break
					}
				}

			}

			send := fmt.Sprintf("Записываем UCLON %s,  день: %s, Сделано %d записей, поездок %d", sit.norma_name(), day, m, trip)
			fmt.Println(send)
			if zapusk_bot {
				send_bot <- send
			}
		}
		trs.Rec("1boU2to5mT77THBlyIL1ZsCJLsqUKE393DcgJga3Z0_w", "'uklon'!A2", res_trip, false, true)
	}

	for _, t := range days {
		rec_days(t)

	}

	control.save()
}

func delete_trip_table() {

	ned := 45
	year := 2022
	for _, sit := range create_sitys() {

		db, err := sql.Open("sqlite3", path+"sqlite/trip.db")
		if err != nil {
			fmt.Println("Запись", err)
		}
		defer db.Close()

		name_table := sit.ident + "_ned_" + strconv.Itoa(ned) + "_" + strconv.Itoa(year)

		//zapr1 := "CREATE TABLE " + name_table + " (id INTEGER PRIMARY KEY AUTOINCREMENT, data VARCHAR(10), name NVARCHAR(70), dbase VARCHAR(135), imp VARCHAR(155), type VARCHAR(5));"
		zapr2 := "drop table " + name_table + ";"

		fmt.Println(zapr2)

		_, err = db.Exec(zapr2)
		fmt.Println(sit.name, err)

	}
}

// записывает поездки по UKLON
// по всем городам
// списку городов за список дат
func get_bgq_uclon_all_city(days []time.Time) []uclon_sum {

	res := make([]uclon_sum, 0)

	for _, dayT := range days {

		day := dayT.Format("2006-01-02")

		// получаем поездки за день
		s := "SELECT * FROM DB_External.uklon_report_orders WHERE date = " + "\"" + day + "\""
		//fmt.Println(s)

		list, _ := BGQ.Read_BGQ(s)

		for _, row := range list {

			if row[0] == nil || row[1] == nil || row[13] == nil || row[14] == nil || row[4] == nil || row[5] == nil || row[8] == nil {
				continue
			}
			if row[2] == "canceled" {
				continue
			}
			var r uclon_sum
			r.name = row[1].(string) + " " + row[2].(string)
			r.name = trs.Clear_32(r.name)
			r.city = row[13].(string)

			dat := row[14].(civil.Date)
			r.data = dat.String()
			r.data = fmt.Sprintf("%02d.%02d.%04d", dat.Day, dat.Month, dat.Year)

			r.trip = row[4].(int64)
			r.travel = row[5].(float64)
			// не совсем понятно, какую сумму
			// надо писать
			// в табеле записан  тотал баланс ячейка 11
			//  мы запишем ячейку 6
			r.cash = row[6].(float64)

			// if r.city == "Lutsk" {
			// 	fmt.Println(r.city)
			// 	fmt.Println(row)
			// 	fmt.Println(r)
			// 	fmt.Println()
			// }

			res = append(res, r)
		}
	}
	return res
}
