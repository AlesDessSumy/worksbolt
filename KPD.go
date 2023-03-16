package main

import (
	"Citys"
	"database/sql"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
	"trs"
)

type Grag_hist struct {
	Name      string
	Sity      string
	Data      string
	Cars_graf int
	Cars_map  int
	Cars_work int
}

// запись дней в БД
func REC_KPD_DAY_SITY_DB(ld []time.Time, ignor_control bool) {

	OS(NS("Run ---  REC_KPD_DAY_SITY_DB(day) ", 40))

	control := control_execute{}
	control.data = time.Now().Format("02.01.2006 15-04")
	control.name = "REC_KPD_DAY_SITY_DB(day)"
	control.coment = "completed"

	if control.read() {
		fmt.Println("Uclon Trip Read completed successfully")
		if !ignor_control {
			return
		}
	}

	file := path + "sqlite/KPD.DB"

	db, err := sql.Open("sqlite3", file)
	if err != nil {
		errorLog.Println("Запись", err)
	}
	defer db.Close()

	for _, day := range ld {

		//fmt.Println(day.Format("02.01.2006 15:04"))

		for _, sit := range create_sitys() {
			//fmt.Println(sit.name)
			<-time.After(time.Second)
			report := KPD_REPORT_SITY_DAY(sit, day)
			//fmt.Println(report)

			rs := report.Slise_String_DB()

			dt := day.Format("02.01.2006")

			// удаляем перед записью
			del := fmt.Sprintf("DELETE FROM  day WHERE data = '%s' AND sity = '%s'", dt, sit.name)
			_, err = db.Exec(del)
			if err != nil {
				errorLog.Println(err)
			}
			// запись дня в БД
			// fmt.Println(report)
			// fmt.Println(rs)
			zapr := fmt.Sprintf("insert into day (data, sity, dbase) values ('%s', '%s', '%s')", rs[1], rs[0], rs[2])
			_, err = db.Exec(zapr)
			if err != nil {
				errorLog.Println(err)
			}
		}
	}
	control.save()
}

// просчет средних поездок
// за три месяца
// + текущий месяц
func AVERAGE_TRIP_FOR_DAY(ignor_control bool) {

	OS(NS("Run ---  AVERAGE_TRIP_FOR_DAY() ", 40))

	control := control_execute{}
	control.data = time.Now().Format("02.01.2006 15-04")
	control.name = "AVERAGE_TRIP_FOR_DAY()"
	control.coment = "completed"

	if control.read() {
		fmt.Println("Uclon Trip Read completed successfully")
		if !ignor_control {
			return
		}
	}

	defer control.save()

	t := time.Now()

	t1 := t.AddDate(0, -1, 0)
	t2 := t.AddDate(0, -2, 0)
	t3 := t.AddDate(0, -3, 0)

	res := make([][]string, 0)
	res = append(res, []string{"Среднее количество поездок на авто в паркеб шт.", "", "", "", "", "", "", "", "", ""})
	res = append(res, []string{"Город", mounth_name_RUS(int(t3.Month())), mounth_name_RUS(int(t2.Month())), mounth_name_RUS(int(t1.Month())), "Среднее значение", "План на " + mounth_name_RUS(int(t.Month())), "Показатели текущего месяца", "план (% выполнения )", "план +5% (% выполнения )", "план + 10% (% выполнения )"})

	type work struct {
		data string
		city string
		trip int
		car  int
	}

	rs := make([]work, 0)

	file := path + "sqlite/KPD.DB"

	db, err := sql.Open("sqlite3", file)
	if err != nil {
		fmt.Println("Запись", err)
	}

	zapr := "SELECT * FROM day"

	rows, err := db.Query(zapr)
	if err != nil {
		errorLog.Println(err)

	}

	for rows.Next() {
		p := make([]string, 4)
		err = rows.Scan(&p[0], &p[1], &p[2], &p[3])
		if err != nil {
			errorLog.Println(err)
			continue
		}
		s := strings.Split(p[3], "|")
		car := string_to_int(s[38])
		trip := string_to_int(s[1])
		if p[1] == "" {
			continue
		}
		r := work{data: p[1], city: p[2], trip: trip, car: car}
		//fmt.Println(r)
		rs = append(rs, r)
	}

	fmt.Println(len(rs))

	create := func(cit string, month, year int) (trip, car int, res float64) {
		days := mounth(month, year)
		//	fmt.Println(month, year, len(days))
		for _, day := range days {
			s := day.Format("02.01.2006")
			for _, i := range rs {
				if s == i.data && cit == i.city {
					trip = trip + i.trip
					car = car + i.car
					break
				}
			}
		}
		if car > 0 {
			res = float64(trip) / float64(car)
		}

		return
	}

	for _, cit := range create_sitys() {
		fmt.Println(cit.norma_name())
		tr3, c3, r3 := create(cit.name, int(t3.Month()), t3.Year())
		tr2, c2, r2 := create(cit.name, int(t2.Month()), t2.Year())
		tr1, c1, r1 := create(cit.name, int(t1.Month()), t1.Year())

		_, _, r0 := create(cit.name, int(t.Month()), t.Year())

		r := make([]string, 10)
		r[0] = cit.name
		r[1] = float_to_string(r3)
		r[2] = float_to_string(r2)
		r[3] = float_to_string(r1)
		var aver float64
		if c1+c2+c3 > 0 {
			aver = float64(tr1+tr2+tr3) / float64(c1+c2+c3)
		}
		plan := aver * 1
		r[4] = float_to_string(aver)
		r[5] = float_to_string(plan)
		r[6] = float_to_string(r0)
		if plan > 0 {
			r[7] = float_to_string((r0/(plan*1) - 1) * 100)
		}
		if plan > 0 {
			r[8] = float_to_string((r0/(plan*1.05) - 1) * 100)
		}
		if plan > 0 {
			r[9] = float_to_string((r0/(plan*1.11) - 1) * 100)
		}

		//fmt.Println(r)
		res = append(res, r)
	}
	trs.Rec("172ZsPP13-yVv7EVUZrymXpX1JjuSJG2RmD34p7PsM9I", "'norma_trip'!A1", res, true, true)

}

// читает кпд городов из базы
func READ_KPD_DAY_SITY_DB(s string, ld []time.Time) (res []sity_report) {

	file := path + "sqlite/KPD.DB"

	db, err := sql.Open("sqlite3", file)
	if err != nil {
		fmt.Println("Запись", err)
	}

	defer db.Close()

	for _, day := range ld {

		rs := sity_report{}

		dt := day.Format("02.01.2006")

		zapr := ""
		if s == "*" {
			zapr = fmt.Sprintf("SELECT * FROM day WHERE data = '%s'", dt)
		} else {
			zapr = fmt.Sprintf("SELECT * FROM day WHERE data = '%s' AND sity = '%s'", dt, s)
		}

		rows, err := db.Query(zapr)
		if err != nil {
			errorLog.Println(err)
			continue
		}

		for rows.Next() {
			p := make([]string, 4)
			err = rows.Scan(&p[0], &p[1], &p[2], &p[3])
			if err != nil {
				errorLog.Println(err)
				continue
			}

			rs.DB_String_Slise(p[1:])
			res = append(res, rs)
		}
	}
	return
}

// Запись КПД по городам
func KPD_REPORT_ALL(ignor_control bool) {

	OS(NS("Run ---  KPD_REPORT_ALL() ", 40))

	control := control_execute{}
	control.data = time.Now().Format("02.01.2006 15-04")
	control.name = "KPD_REPORT_ALL"
	control.coment = "completed"

	if control.read() {
		fmt.Println("KPD_REPORT_ALL completed successfully")
		if !ignor_control {
			return
		}
	}
	KPD_REPORT_ALL_SITY_7_DAY()
	for _, sit := range create_sitys() {
		KPD_REPORT_SITY_FOUR_NED(sit)
	}

	OS(NS("Finish .... KPD_REPORT_ALL() ", 40))
	control.save()
}

// наполняет значениями GGRAFIK_HISTORY
// обновляется раз в час
func Create_GGRAFIK_HISTORY() {

	OS(NS("Run ---  go Create_GGRAFIK_HISTORY()", 40))

	for {

		if !work_time() {
			<-time.After(5 * time.Minute)
			continue
		}

		list, err := trs.Read_sheets_CASHBOX_err("172ZsPP13-yVv7EVUZrymXpX1JjuSJG2RmD34p7PsM9I", "'nw_list'!A3:F", 0)

		if err != nil {
			<-time.After(time.Minute)
			continue
		}

		for _, i := range list {
			r := Grag_hist{}
			r.Name = i[0] + " " + i[1]
			r.Sity = i[1]
			r.Data = i[0]
			r.Cars_graf, _ = strconv.Atoi(i[3])
			r.Cars_map, _ = strconv.Atoi(i[4])
			r.Cars_work, _ = strconv.Atoi(i[5])
			GRAFIK_HISTORY[r.Name] = r
		}
		<-time.After(time.Hour)
	}
}

// отчет города
type sity_report struct {
	Data          time.Time
	Sity          string
	Day           string
	Trip_car_work float64
	Trip_car_trip float64
	Hour          float64
	Tipbolt_hour  float64
	Waitbolt      float64
	Trip          int
	Trip_Bolt     int
	Trip_Uklon    int
	Trip_Uber     int
	Probeg        int
	Preprobeg     int
	Cash          int
	Cash_Bolt     int
	Cash_Uklon    int
	Cash_Uber     int
	Ecip          int // Экипаж
	Ecip1212      int
	Dvaakk        int // два аакаунта
	Vyh           int // выходной
	Progul        int // прогул
	Bezvod        int // без водителя
	Remont        int // ремонт
	Dtp           int // ДТП
	Evak          int // Эвакуация
	Naprodag      int // Продано
	Shtraff       int // Штрафмайданчик
	Doki          int // Документи
	Ugon          int // Угон
	Keys          int // Ключі
	Volonter      int // Волонтер
	Spysan        int // списання
	Peregon       int // Перегін
	Pomylka       int // Помилка
	Dop4          int // доп4
	Dop5          int // доп5
	Dop6          int // доп6
	Dop7          int // доп7
	Dop8          int // доп8
	Dop9          int // доп9
	Dop10         int // доп10
	Dop11         int // доп11
	Dop12         int // больше 30 поездок
	Cars          int // всего авто
	Car_work      int // авто в работе
	Car_trip      int // авто с поездками
	Car_0_trip    int // процентовка
	Car_5_trip    int
	Car_10_trip   int
	Car_15_trip   int
	Car_20_trip   int
	Car_25_trip   int // больше 25 поездок
	Car_99_trip   int
	day_report    int // количество дней в отчете за неделю
}

func (a *sity_report) Procent_Trip(b int) {
	if b > 0 {
		a.Car_trip++
	}
	if b == 0 {
		a.Car_0_trip++
	} else if b < 6 {
		a.Car_5_trip++
	} else if b < 11 {
		a.Car_10_trip++
	} else if b < 16 {
		a.Car_15_trip++
	} else if b < 21 {
		a.Car_20_trip++
	} else if b < 26 {
		a.Car_25_trip++
	} else {
		a.Car_99_trip++
	}
}

// подготовка элемента к записи в БД
func (a sity_report) Slise_String_DB() (res []string) {
	res = append(res, a.Sity, a.Day)
	str := ""
	str = str + fmt.Sprintf("%.2f", a.Hour) + "|" + fmt.Sprint(a.Trip) + "|" + fmt.Sprint(a.Trip_Bolt) + "|" + fmt.Sprint(a.Trip_Uklon) + "|" + fmt.Sprint(a.Probeg) + "|" + fmt.Sprint(a.Preprobeg)
	str = str + "|" + fmt.Sprint(a.Cash) + "|" + fmt.Sprint(a.Cash_Bolt) + "|" + fmt.Sprint(a.Cash_Uklon)
	str = str + "|" + fmt.Sprint(a.Ecip) + "|" + fmt.Sprint(a.Ecip1212) + "|" + fmt.Sprint(a.Dvaakk) + "|" + fmt.Sprint(a.Vyh) + "|" + fmt.Sprint(a.Progul) + "|" + fmt.Sprint(a.Bezvod)
	str = str + "|" + fmt.Sprint(a.Remont) + "|" + fmt.Sprint(a.Dtp) + "|" + fmt.Sprint(a.Evak) + "|" + fmt.Sprint(a.Naprodag) + "|" + fmt.Sprint(a.Shtraff) + "|" + fmt.Sprint(a.Doki)
	str = str + "|" + fmt.Sprint(a.Ugon) + "|" + fmt.Sprint(a.Keys) + "|" + fmt.Sprint(a.Volonter) + "|" + fmt.Sprint(a.Spysan) + "|" + fmt.Sprint(a.Pomylka)
	str = str + "|" + fmt.Sprint(a.Trip_Uber) + "|" + fmt.Sprint(a.Cash_Uber) + "|" + fmt.Sprint(float_to_string(a.Waitbolt)) + "|" + fmt.Sprint(a.Dop4) + "|" + fmt.Sprint(a.Dop5) + "|" + fmt.Sprint(a.Dop6)
	str = str + "|" + fmt.Sprint(a.Dop7) + "|" + fmt.Sprint(a.Dop8) + "|" + fmt.Sprint(a.Dop9) + "|" + fmt.Sprint(a.Dop10) + "|" + fmt.Sprint(a.Dop11) + "|" + fmt.Sprint(a.Dop12)
	str = str + "|" + fmt.Sprint(a.Cars) + "|" + fmt.Sprint(a.Car_work) + "|" + fmt.Sprint(a.Car_trip)
	str = str + "|" + fmt.Sprint(a.Car_0_trip) + "|" + fmt.Sprint(a.Car_5_trip) + "|" + fmt.Sprint(a.Car_10_trip) + "|" + fmt.Sprint(a.Car_15_trip) + "|" + fmt.Sprint(a.Car_20_trip) + "|" + fmt.Sprint(a.Car_25_trip) + "|" + fmt.Sprint(a.Car_99_trip) + "|"
	res = append(res, str)
	return
}

// распаковка строки в элемент КПД
func (a *sity_report) DB_String_Slise(b []string) {
	var err error
	s := make([]string, 0)
	s = append(s, b[1])
	s = append(s, b[0])
	s = append(s, strings.Split(b[2], "|")...)
	a.Sity = s[0]
	a.Day = s[1]
	a.Data, err = time.Parse("02.01.2006", s[1])
	if err != nil {
		errorLog.Println(err)
	}
	a.Hour, _ = strconv.ParseFloat(s[2], 64)
	a.Trip, _ = strconv.Atoi(s[3])
	a.Trip_Bolt, _ = strconv.Atoi(s[4])
	a.Trip_Uklon, _ = strconv.Atoi(s[5])

	a.Probeg, _ = strconv.Atoi(s[6])
	a.Preprobeg, _ = strconv.Atoi(s[7])
	a.Cash, _ = strconv.Atoi(s[8])
	a.Cash_Bolt, _ = strconv.Atoi(s[9])
	a.Cash_Uklon, _ = strconv.Atoi(s[10])

	a.Ecip, _ = strconv.Atoi(s[11])
	a.Ecip1212, _ = strconv.Atoi(s[12])
	a.Dvaakk, _ = strconv.Atoi(s[13])
	a.Vyh, _ = strconv.Atoi(s[14])
	a.Progul, _ = strconv.Atoi(s[15])
	a.Bezvod, _ = strconv.Atoi(s[16])
	a.Remont, _ = strconv.Atoi(s[17])
	a.Dtp, _ = strconv.Atoi(s[18])
	a.Evak, _ = strconv.Atoi(s[19])
	a.Naprodag, _ = strconv.Atoi(s[20])
	a.Shtraff, _ = strconv.Atoi(s[21])
	a.Doki, _ = strconv.Atoi(s[22])
	a.Ugon, _ = strconv.Atoi(s[23])
	a.Keys, _ = strconv.Atoi(s[24])
	a.Volonter, _ = strconv.Atoi(s[25])
	a.Spysan, _ = strconv.Atoi(s[26])
	a.Pomylka, _ = strconv.Atoi(s[27])
	a.Trip_Uber, _ = strconv.Atoi(s[28])
	a.Cash_Uber, _ = strconv.Atoi(s[29])
	a.Waitbolt = string_to_float(s[30])
	a.Dop4, _ = strconv.Atoi(s[31])
	a.Dop5, _ = strconv.Atoi(s[32])
	a.Dop6, _ = strconv.Atoi(s[33])
	a.Dop7, _ = strconv.Atoi(s[34])
	a.Dop8, _ = strconv.Atoi(s[35])
	a.Dop9, _ = strconv.Atoi(s[36])
	a.Dop10, _ = strconv.Atoi(s[37])
	a.Dop11, _ = strconv.Atoi(s[38])
	a.Dop12, _ = strconv.Atoi(s[39])
	a.Cars, _ = strconv.Atoi(s[40])
	a.Car_work, _ = strconv.Atoi(s[41])
	a.Car_trip, _ = strconv.Atoi(s[42])
	a.Car_0_trip, _ = strconv.Atoi(s[43])
	a.Car_5_trip, _ = strconv.Atoi(s[44])
	a.Car_10_trip, _ = strconv.Atoi(s[45])
	a.Car_15_trip, _ = strconv.Atoi(s[46])
	a.Car_20_trip, _ = strconv.Atoi(s[47])
	a.Car_25_trip, _ = strconv.Atoi(s[48])
	a.Car_99_trip, _ = strconv.Atoi(s[49])
	a.day_report, _ = strconv.Atoi(s[50])
}

// Печать элемента sity_report
func (a *sity_report) Print() {
	fmt.Println()
	fmt.Println("sity_report         :")
	fmt.Println("Data                :", a.Data.Format(TNS))
	fmt.Println("Sity                :", a.Sity)
	fmt.Println("Day                 :", a.Day)
	fmt.Println("Trip_car_work       :", a.Trip_car_work)
	fmt.Println("Trip_car_trip       :", a.Trip_car_trip)
	fmt.Println("Hour                :", a.Hour)
	fmt.Println("Tipbolt_hour        :", a.Tipbolt_hour)
	fmt.Println("Waitbolt            :", a.Waitbolt)
	fmt.Println("Trip                :", a.Trip)
	fmt.Println("Trip_Bolt           :", a.Trip_Bolt)
	fmt.Println("Trip_Uklon          :", a.Trip_Uklon)
	fmt.Println("Trip_Uber           :", a.Trip_Uber)
	fmt.Println("Probeg              :", a.Probeg)
	fmt.Println("Preprobeg           :", a.Preprobeg)
	fmt.Println("Cash                :", a.Cash)
	fmt.Println("Cash_Bolt           :", a.Cash_Bolt)
	fmt.Println("Cash_Uklon          :", a.Cash_Uklon)
	fmt.Println("Cash_Uber           :", a.Cash_Uber)
	fmt.Println("Ecip                :", a.Ecip)
	fmt.Println("Ecip1212            :", a.Ecip1212)
	fmt.Println("Dvaakk              :", a.Dvaakk)
	fmt.Println("Vyh                 :", a.Vyh)
	fmt.Println("Progul              :", a.Progul)
	fmt.Println("Bezvod              :", a.Bezvod)
	fmt.Println("Remont              :", a.Remont)
	fmt.Println("Dtp                 :", a.Dtp)
	fmt.Println("Evak                :", a.Evak)
	fmt.Println("Naprodag            :", a.Naprodag)
	fmt.Println("Shtraff             :", a.Shtraff)
	fmt.Println("Doki                :", a.Doki)
	fmt.Println("Ugon                :", a.Ugon)
	fmt.Println("Keys                :", a.Keys)
	fmt.Println("Volonter            :", a.Volonter)
	fmt.Println("Spysan              :", a.Spysan)
	fmt.Println("Peregon             :", a.Peregon)
	fmt.Println("Pomylka             :", a.Pomylka)
	fmt.Println("Dop4                :", a.Dop4)
	fmt.Println("Dop5                :", a.Dop5)
	fmt.Println("Dop6                :", a.Dop6)
	fmt.Println("Dop7                :", a.Dop7)
	fmt.Println("Dop8                :", a.Dop8)
	fmt.Println("Dop9                :", a.Dop9)
	fmt.Println("Dop10               :", a.Dop10)
	fmt.Println("Dop11               :", a.Dop11)
	fmt.Println("Dop12               :", a.Dop12)
	fmt.Println("Cars                :", a.Cars)
	fmt.Println("Car_work            :", a.Car_work)
	fmt.Println("Car_trip            :", a.Car_trip)
	fmt.Println("Car_0_trip          :", a.Car_0_trip)
	fmt.Println("Car_5_trip          :", a.Car_5_trip)
	fmt.Println("Car_10_trip         :", a.Car_10_trip)
	fmt.Println("Car_15_trip         :", a.Car_15_trip)
	fmt.Println("Car_20_trip         :", a.Car_20_trip)
	fmt.Println("Car_25_trip         :", a.Car_25_trip)
	fmt.Println("Car_99_trip         :", a.Car_99_trip)
	fmt.Println("day_report          :", a.day_report)
}

// печать элемента
// основные параметры
func (a sity_report) Print_cor() { // %s
	fmt.Printf("По городу  %s на дату %s часов %.1f\n  поездок %d болт, поездок %d уклон,\n авто %d всего, авто %d в работе, авто %d без водителя, пробег всего %d\n\n ", a.Sity, a.Day, a.Hour, a.Trip_Bolt, a.Trip_Uklon, a.Cars, a.Car_work, a.Bezvod, int(a.Probeg))
}

// подготовка строки для отчета
func (a sity_report) Strings_KPD(ln int, sel bool) (res []string) {
	if a.Sity == "" {
		res = make([]string, 39)
		return
	}
	if ln > 0 {
		a.Hour = a.Hour / float64(ln)
		a.Bezvod = a.Bezvod / ln
		a.Vyh = a.Vyh / ln
		a.Dvaakk = a.Dvaakk / ln
		a.Ecip = a.Ecip / ln
		a.Remont = a.Remont / ln
		a.Dtp = a.Dtp / ln
		a.Evak = a.Evak / ln
		a.Cars = a.Cars / ln
		a.Car_work = a.Car_work / ln
		a.Car_trip = a.Car_trip / ln
		a.Trip_car_work = a.Trip_car_work / float64(ln)
		a.Trip_car_trip = a.Trip_car_trip / float64(ln)
	}
	if sel {
		res = append(res, a.Sity)
	} else {
		if strings.Contains(a.Day, "Статистика") {
			res = append(res, "Усього:")
		} else {
			dnd := int(a.Data.Weekday())
			switch {
			case dnd == 0:
				res = append(res, "Вс")
			case dnd == 1:
				res = append(res, "Пн")
			case dnd == 2:
				res = append(res, "Вт")
			case dnd == 3:
				res = append(res, "Ср")
			case dnd == 4:
				res = append(res, "Чт")
			case dnd == 5:
				res = append(res, "Пт")
			case dnd == 6:
				res = append(res, "Сб")
			}
		}
	}
	if sel {
		res = append(res, a.Day)
	} else {
		res = append(res, a.Data.Format("02.01.2006"))
	}

	res = append(res, fmt.Sprint(a.Bezvod))
	res = append(res, fmt.Sprint(a.Dtp))
	res = append(res, fmt.Sprint(a.Remont))
	res = append(res, fmt.Sprint(a.Naprodag))
	res = append(res, fmt.Sprint(a.Shtraff))
	res = append(res, fmt.Sprint(a.Keys))
	res = append(res, fmt.Sprint(a.Volonter))
	res = append(res, fmt.Sprint(a.Spysan))
	res = append(res, fmt.Sprint(a.Pomylka))

	res = append(res, fmt.Sprint(a.Vyh))
	res = append(res, fmt.Sprint(a.Car_work))
	res = append(res, fmt.Sprint(a.Car_trip))

	res = append(res, fmt.Sprint(a.Trip_Bolt))
	res = append(res, fmt.Sprint(a.Trip_Uklon))
	res = append(res, fmt.Sprint(a.Trip_Uber))
	res = append(res, fmt.Sprint(a.Trip))
	res = append(res, strings.Replace(fmt.Sprintf("%.1f", float64(a.Hour)), ".", ",", 1))

	var trphour float64
	if a.Hour > 0 {
		trphour = float64(a.Trip_Bolt) / a.Hour
	}

	res = append(res, strings.Replace(fmt.Sprintf("%.1f", trphour), ".", ",", 1))
	res = append(res, float_to_string(a.Waitbolt*100))

	res = append(res, strings.Replace(fmt.Sprintf("%.1f", float64(a.Cash_Bolt)), ".", ",", 1))
	res = append(res, strings.Replace(fmt.Sprintf("%.1f", float64(a.Cash_Uklon)), ".", ",", 1))
	res = append(res, strings.Replace(fmt.Sprintf("%.1f", float64(a.Cash)), ".", ",", 1))

	if a.Cars > 0 {
		res = append(res, strings.Replace(fmt.Sprintf("%.1f", float64(a.Cash)/float64(a.Cars)), ".", ",", 1))
	} else {
		res = append(res, "0")
	}

	if a.Car_work > 0 {
		res = append(res, fmt.Sprint(a.Trip/a.Car_work))
	} else {
		res = append(res, "0")
	}

	if a.Car_trip > 0 {
		res = append(res, fmt.Sprint(a.Trip/a.Car_trip))
	} else {
		res = append(res, "0")
	}

	if a.Cars > 0 {
		res = append(res, fmt.Sprint(a.Trip/a.Cars))
	} else {
		res = append(res, "0")
	}

	// Блок пробега
	res = append(res, fmt.Sprint(a.Probeg))
	if a.Car_work > 0 {
		res = append(res, fmt.Sprint(a.Probeg/a.Car_work))
	} else {
		res = append(res, "0")
	}
	if a.Car_trip > 0 {
		res = append(res, fmt.Sprint(a.Probeg/a.Car_trip))
	} else {
		res = append(res, "0")
	}
	if a.Cars > 0 {
		res = append(res, fmt.Sprint(a.Probeg/a.Cars))
	} else {
		res = append(res, "0")
	}
	if a.Trip > 0 {
		res = append(res, strings.Replace(fmt.Sprintf("%.2f", float64(a.Probeg)/float64(a.Trip)), ".", ",", 1))
	} else {
		res = append(res, "0")
	}

	// блок кеша
	if a.Car_trip > 0 {
		res = append(res, fmt.Sprint(a.Cash/a.Car_trip))
	} else {
		res = append(res, "0")
	}
	if a.Trip > 0 {
		res = append(res, fmt.Sprint(a.Cash/a.Trip))
	} else {
		res = append(res, "0")
	}
	if a.Probeg > 0 {
		res = append(res, strings.Replace(fmt.Sprintf("%.2f", float64(a.Cash)/float64(a.Probeg)), ".", ",", 1))
	} else {
		res = append(res, "0")
	}

	res = append(res, fmt.Sprint(a.Dvaakk))
	res = append(res, fmt.Sprint(a.Ecip))

	if a.day_report == 0 {
		a.day_report = 1
	}

	// норма 15
	if a.Car_work > 0 {
		norma := float64(a.Car_15_trip) / float64(a.Car_work) * 100 / float64(a.day_report)
		res = append(res, strings.Replace(fmt.Sprintf("%.2f", norma), ".", ",", 1))
	} else {
		res = append(res, "0")
	}
	// норма 20
	if a.Car_work > 0 {
		norma := float64(a.Car_20_trip) / float64(a.Car_work) * 100 / float64(a.day_report)
		res = append(res, strings.Replace(fmt.Sprintf("%.2f", norma), ".", ",", 1))
	} else {
		res = append(res, "0")
	}
	//  норма 25
	if a.Car_work > 0 {
		norma := float64(a.Car_25_trip) / float64(a.Car_work) * 100 / float64(a.day_report)
		res = append(res, strings.Replace(fmt.Sprintf("%.2f", norma), ".", ",", 1))
	} else {
		res = append(res, "0")
	}

	res = append(res, fmt.Sprint(a.Cars))
	return
}

// отчет за 6 недель
func KPD_REPORT_6_ned() {

	res := make([][]string, 0)

	t := time.Now()
	dn := int(t.Weekday())
	if dn == 0 {
		dn = 7
	}
	days := last_days(42 + dn - 1)

	fmt.Println(days[0].Format("02.01.2006"))

	sitys := create_sitys()

	for _, sit := range sitys {
		kpds := READ_KPD_DAY_SITY_DB(sit.name, days[:42])
		rs := make([]int, 6)
		for n, i := range kpds {
			rs[n/7] = rs[n/7] + i.Trip
		}
		rpr := func(a, b int) string {
			return float_to_string((float64(a) - float64(b)) / float64(b) * 100)
		}
		r := make([]string, 12)
		r[0] = sit.name
		r[1] = fmt.Sprint(rs[0])
		r[2] = rpr(rs[1], rs[0])
		r[3] = fmt.Sprint(rs[1])
		r[4] = rpr(rs[2], rs[1])
		r[5] = fmt.Sprint(rs[2])
		r[6] = rpr(rs[3], rs[2])
		r[7] = fmt.Sprint(rs[3])
		r[8] = rpr(rs[4], rs[3])
		r[9] = fmt.Sprint(rs[4])
		r[10] = rpr(rs[5], rs[4])
		r[11] = fmt.Sprint(rs[5])
		res = append(res, r)
	}
	list := "'city6ned'!A3"
	fmt.Println(res)
	trs.Rec("172ZsPP13-yVv7EVUZrymXpX1JjuSJG2RmD34p7PsM9I", list, res, true, true)
}

type PlanCar struct {
	City     string
	Car_01   int
	Car_yest int
}

// КПД щтчет по городу
// всем городам текущий месяц
func KPD_REPORT_Plan(day []time.Time) (res []sity_report, car []PlanCar, landay int) {

	sitys := create_sitys()

	if len(day) == 0 {
		return
	}

	landay = len(day)

	for _, sit := range sitys {

		rp := sity_report{}
		rc := PlanCar{}
		rc.City = sit.name
		for i := 0; i < len(day); i++ {
			d := day[i]

			kpd := sity_report{}
			kpds := READ_KPD_DAY_SITY_DB(sit.name, []time.Time{d})
			if len(kpds) > 0 {
				kpd = kpds[0]
			} else {
				continue
			}
			if i == 0 {
				rc.Car_01 = kpd.Cars
			}
			if i == len(day)-1 {
				rc.Car_yest = kpd.Cars
			}

			rp.Sity = sit.name

			rp.Trip_car_work = rp.Trip_car_work + kpd.Trip_car_work
			rp.Trip_car_trip = rp.Trip_car_trip + kpd.Trip_car_trip
			rp.Hour = rp.Hour + kpd.Hour

			rp.Trip = rp.Trip + kpd.Trip
			rp.Trip_Bolt = rp.Trip_Bolt + kpd.Trip_Bolt
			rp.Trip_Uklon = rp.Trip_Uklon + kpd.Trip_Uklon
			rp.Trip_Uber = rp.Trip_Uber + kpd.Trip_Uber

			rp.Probeg = rp.Probeg + kpd.Probeg
			//rs.Preprobeg
			rp.Cash = rp.Cash + kpd.Cash
			rp.Cash_Bolt = rp.Cash_Bolt + kpd.Cash_Bolt
			rp.Cash_Uklon = rp.Cash_Uklon + kpd.Cash_Uklon
			rp.Cash_Uber = rp.Cash_Uber + kpd.Cash_Uber

			// Экипаж
			rp.Ecip = rp.Ecip + kpd.Ecip
			// два аакаунта
			rp.Dvaakk = rp.Dvaakk + kpd.Dvaakk
			// выходной
			rp.Vyh = rp.Vyh + kpd.Vyh
			// прогул
			rp.Progul = rp.Progul + kpd.Progul
			// без водителя
			rp.Bezvod = rp.Bezvod + kpd.Bezvod
			// ремонт
			rp.Remont = rp.Remont + kpd.Remont
			// ДТП
			rp.Dtp = rp.Dtp + kpd.Dtp
			// Эвакуация
			rp.Evak = rp.Evak + kpd.Evak
			// Продано
			rp.Naprodag = rp.Naprodag + kpd.Naprodag
			//Штрафмайданчик
			rp.Shtraff = rp.Shtraff + kpd.Shtraff
			//Документи
			rp.Doki = rp.Doki + kpd.Doki
			//Угон
			rp.Ugon = rp.Ugon + kpd.Ugon
			//Ключі
			rp.Keys = rp.Keys + kpd.Keys
			//Волонтер
			rp.Volonter = rp.Volonter + kpd.Volonter
			//списання
			rp.Spysan = rp.Spysan + kpd.Spysan
			// Помилка
			rp.Pomylka = rp.Pomylka + kpd.Pomylka
			// всего авто
			rp.Cars = rp.Cars + kpd.Cars
			// авто в работе
			rp.Car_work = rp.Car_work + kpd.Car_work
			// авто с поездками
			rp.Car_trip = rp.Car_trip + kpd.Car_trip
			// процентовка
			rp.Car_0_trip = rp.Car_0_trip + kpd.Car_0_trip
			rp.Car_5_trip = rp.Car_5_trip + kpd.Car_5_trip
			rp.Car_10_trip = rp.Car_10_trip + kpd.Car_10_trip
			rp.Car_15_trip = rp.Car_15_trip + kpd.Car_15_trip
			rp.Car_20_trip = rp.Car_20_trip + kpd.Car_20_trip
			rp.Car_25_trip = rp.Car_25_trip + kpd.Car_25_trip
			rp.Car_99_trip = rp.Car_99_trip + kpd.Car_99_trip

		}
		// if sit.name == "Київ" {
		// 	rc.Car_01 = rc.Car_01 - 3
		// 	rc.Car_yest = rc.Car_yest - 3
		// 	//fmt.Println(rc.Car_01, rc.Car_yest, sit.name)
		// }
		// if sit.name == "Житомир" {
		// 	rc.Car_01 = rc.Car_01 - 1
		// 	//rc.Car_yest = rc.Car_yest - 5
		// 	//fmt.Println(rc.Car_01, rc.Car_yest, sit.name)
		// }
		// if sit.name == "Запоріжжя" {
		// 	rc.Car_01 = rc.Car_01 - 1
		// 	//rc.Car_yest = rc.Car_yest - 5
		// 	//fmt.Println(rc.Car_01, rc.Car_yest, sit.name)
		// }
		// if sit.name == "Черкаси" {
		// 	rc.Car_01 = rc.Car_01 - 1
		// 	//rc.Car_yest = rc.Car_yest - 5
		// 	//fmt.Println(rc.Car_01, rc.Car_yest, sit.name)
		// }

		res = append(res, rp)
		car = append(car, rc)
		//fmt.Println(rc, rp)
	}
	return
}

// КПД щтчет по городу
// всем городам за 7 дней
func KPD_REPORT_ALL_SITY_7_DAY() {
	res := make([][]string, 0)

	day := last_days(7)
	sitys := create_sitys()
	for i := 6; i > -1; i-- {
		d := day[i]
		res = append(res, append_probel([]string{"", "", "Статистика по містам України за " + d.Format("02.01.2006")}, 37))
		rs := make([]sity_report, 0)
		rp := sity_report{}
		t := 0
		for _, sit := range sitys {

			//kpd := KPD_REPORT_SITY_DAY(sit, d)
			kpd := sity_report{}
			kpds := READ_KPD_DAY_SITY_DB(sit.name, []time.Time{d})
			if len(kpds) > 0 {
				kpd = kpds[0]
			}

			rs = append(rs, kpd)

			//на сколько дней делить
			if rp.Trip > 0 {
				t++
			}

			rp.Sity = "Усього"
			rp.Day = fmt.Sprintf("Статистика по містам УКРАЇНИ за %s", d.Format("02.01"))
			rp.Trip_car_work = rp.Trip_car_work + kpd.Trip_car_work
			rp.Trip_car_trip = rp.Trip_car_trip + kpd.Trip_car_trip
			rp.Hour = rp.Hour + kpd.Hour

			rp.Trip = rp.Trip + kpd.Trip
			rp.Trip_Bolt = rp.Trip_Bolt + kpd.Trip_Bolt
			rp.Trip_Uklon = rp.Trip_Uklon + kpd.Trip_Uklon
			rp.Trip_Uber = rp.Trip_Uber + kpd.Trip_Uber

			rp.Probeg = rp.Probeg + kpd.Probeg
			//rs.Preprobeg
			rp.Cash = rp.Cash + kpd.Cash
			rp.Cash_Bolt = rp.Cash_Bolt + kpd.Cash_Bolt
			rp.Cash_Uklon = rp.Cash_Uklon + kpd.Cash_Uklon
			rp.Cash_Uber = rp.Cash_Uber + kpd.Cash_Uber

			// Экипаж
			rp.Ecip = rp.Ecip + kpd.Ecip
			// два аакаунта
			rp.Dvaakk = rp.Dvaakk + kpd.Dvaakk
			// выходной
			rp.Vyh = rp.Vyh + kpd.Vyh
			// прогул
			rp.Progul = rp.Progul + kpd.Progul
			// без водителя
			rp.Bezvod = rp.Bezvod + kpd.Bezvod
			// ремонт
			rp.Remont = rp.Remont + kpd.Remont
			// ДТП
			rp.Dtp = rp.Dtp + kpd.Dtp
			// Эвакуация
			rp.Evak = rp.Evak + kpd.Evak
			// Продано
			rp.Naprodag = rp.Naprodag + kpd.Naprodag
			//Штрафмайданчик
			rp.Shtraff = rp.Shtraff + kpd.Shtraff
			//Документи
			rp.Doki = rp.Doki + kpd.Doki
			//Угон
			rp.Ugon = rp.Ugon + kpd.Ugon
			//Ключі
			rp.Keys = rp.Keys + kpd.Keys
			//Волонтер
			rp.Volonter = rp.Volonter + kpd.Volonter
			//списання
			rp.Spysan = rp.Spysan + kpd.Spysan
			// Помилка
			rp.Pomylka = rp.Pomylka + kpd.Pomylka
			// всего авто
			rp.Cars = rp.Cars + kpd.Cars
			// авто в работе
			rp.Car_work = rp.Car_work + kpd.Car_work
			// авто с поездками
			rp.Car_trip = rp.Car_trip + kpd.Car_trip
			// процентовка
			rp.Car_0_trip = rp.Car_0_trip + kpd.Car_0_trip
			rp.Car_5_trip = rp.Car_5_trip + kpd.Car_5_trip
			rp.Car_10_trip = rp.Car_10_trip + kpd.Car_10_trip
			rp.Car_15_trip = rp.Car_15_trip + kpd.Car_15_trip
			rp.Car_20_trip = rp.Car_20_trip + kpd.Car_20_trip
			rp.Car_25_trip = rp.Car_25_trip + kpd.Car_25_trip
			rp.Car_99_trip = rp.Car_99_trip + kpd.Car_99_trip

		}

		sort.Slice(rs, func(i, j int) (less bool) {
			return rs[i].Trip > rs[j].Trip
		})

		for i := 0; i < 23; i++ {
			if i < len(rs)-1 {
				str := rs[i].Strings_KPD(0, false)
				str[0] = rs[i].Sity
				res = append(res, str)
			} else {
				res = append(res, append_probel([]string{}, 40))
			}

		}
		str := rp.Strings_KPD(0, false)
		// fmt.Println(rp)
		// fmt.Println(str)
		res = append(res, str)

	}
	list := "'Україна'!A4"
	trs.Rec("18fpt9H8pR2P3_NY_IAeTpGZGXPav5mgqC-6m6fJg5nc", list, res, false, false)

}

// КПД щтчет по городу за 4 недели
func KPD_REPORT_SITY_FOUR_NED(sit sity) {
	res := make([][]string, 0)
	// строим недели по дням
	// три последних + текущая
	t := int(time.Now().Weekday())
	if t == 0 {
		t = 7
	}

	day := last_days(28 + t - 1)

	if t > 1 {
		day = last_days(28 + t - 1)
	}

	day1 := day[7:14]
	day2 := day[14:21]
	day3 := day[21:28]
	day4 := teck_days()

	if t == 1 {
		day1 = last_days(28)[:7]
		day2 = last_days(21)[:7]
		day3 = last_days(14)[:7]
		day4 = last_days(7)
	}

	rp1, zag1 := KPD_REPORT_SITY_NED(sit, day1)
	rp2, zag2 := KPD_REPORT_SITY_NED(sit, day2)
	rp3, zag3 := KPD_REPORT_SITY_NED(sit, day3)
	rp4, zag4 := KPD_REPORT_SITY_NED(sit, day4)

	res = append(res, KPD_REPORT_STRINGS_NED(rp1, zag1)...)
	res = append(res, KPD_REPORT_STRINGS_NED(rp2, zag2)...)
	res = append(res, KPD_REPORT_STRINGS_NED(rp3, zag3)...)
	res = append(res, KPD_REPORT_STRINGS_NED(rp4, zag4)...)

	list := "'" + sit.name + "'!A4"

	trs.Rec("18fpt9H8pR2P3_NY_IAeTpGZGXPav5mgqC-6m6fJg5nc", list, res, false, false)

}

func KPD_REPORT_STRINGS_NED(r []sity_report, rs sity_report) (res [][]string) {

	str := append_probel([]string{"", "", rs.Day, ""}, 37)
	res = append(res, str)
	for _, i := range r {
		//fmt.Println(i.Strings_KPD())
		res = append(res, i.Strings_KPD(0, false))
	}
	//fmt.Println(rs.Strings_KPD())
	res = append(res, rs.Strings_KPD(0, false))
	return
}

func append_probel(s []string, n int) []string {
	for i := 0; i < n; i++ {
		s = append(s, "")
	}
	return s
}

// КПД щтчет по городу за 1 нед
func KPD_REPORT_SITY_NED(sit sity, ned []time.Time) (r []sity_report, rs sity_report) {
	t := 0
	for _, n := range ned {

		tm := time.Now()
		rp := sity_report{}

		if tm.Format("02.01.2006") != n.Format("02.01.2006") {
			r := READ_KPD_DAY_SITY_DB(sit.name, []time.Time{n})
			if len(r) > 0 {
				rp = r[0]
			}
		}

		r = append(r, rp)

		//на сколько дней делить
		if rp.Trip > 0 {
			t++
		}

		rs.Sity = sit.name
		rs.Day = fmt.Sprintf("Статистика %s-%s", ned[0].Format("02.01"), ned[len(ned)-1].Format("02.01"))
		rs.Trip_car_work = rs.Trip_car_work + rp.Trip_car_work
		rs.Trip_car_trip = rs.Trip_car_trip + rp.Trip_car_trip
		rs.Hour = rs.Hour + rp.Hour

		rs.Trip = rs.Trip + rp.Trip
		rs.Trip_Bolt = rs.Trip_Bolt + rp.Trip_Bolt
		rs.Trip_Uklon = rs.Trip_Uklon + rp.Trip_Uklon

		rs.Probeg = rs.Probeg + rp.Probeg
		//rs.Preprobeg
		rs.Cash = rs.Cash + rp.Cash
		rs.Cash_Bolt = rs.Cash_Bolt + rp.Cash_Bolt
		rs.Cash_Uklon = rs.Cash_Uklon + rp.Cash_Uklon
		// Экипаж
		rs.Ecip = rs.Ecip + rp.Ecip
		// два аакаунта
		rs.Dvaakk = rs.Dvaakk + rp.Dvaakk
		// выходной
		rs.Vyh = rs.Vyh + rp.Vyh
		// прогул
		rs.Progul = rs.Progul + rp.Progul
		// без водителя
		rs.Bezvod = rs.Bezvod + rp.Bezvod
		// ремонт
		rs.Remont = rs.Remont + rp.Remont
		// ДТП
		rs.Dtp = rs.Dtp + rp.Dtp
		// Эвакуация
		rs.Evak = rs.Evak + rp.Evak
		// Продано
		rs.Naprodag = rs.Naprodag + rp.Naprodag
		//Штрафмайданчик
		rs.Shtraff = rs.Shtraff + rp.Shtraff
		//Документи
		rs.Doki = rs.Doki + rp.Doki
		//Угон
		rs.Ugon = rs.Ugon + rp.Ugon
		// Перегін
		rs.Peregon = rs.Peregon + rp.Peregon
		//Ключі
		rs.Keys = rs.Keys + rp.Keys
		//Волонтер
		rs.Volonter = rs.Volonter + rp.Volonter
		//списання
		rs.Spysan = rs.Spysan + rp.Spysan
		// Помилка
		rs.Pomylka = rs.Pomylka + rp.Pomylka
		// всего авто
		rs.Cars = rs.Cars + rp.Cars
		// авто в работе
		rs.Car_work = rs.Car_work + rp.Car_work
		// авто с поездками
		rs.Car_trip = rs.Car_trip + rp.Car_trip
		// процентовка
		rs.Car_0_trip = rs.Car_0_trip + rp.Car_0_trip
		rs.Car_5_trip = rs.Car_5_trip + rp.Car_5_trip
		rs.Car_10_trip = rs.Car_10_trip + rp.Car_10_trip
		rs.Car_15_trip = rs.Car_15_trip + rp.Car_15_trip
		rs.Car_20_trip = rs.Car_20_trip + rp.Car_20_trip
		rs.Car_25_trip = rs.Car_25_trip + rp.Car_25_trip
		rs.Car_99_trip = rs.Car_99_trip + rp.Car_99_trip
	}

	// строим недельный отчет
	if t == 1 {
		rs.day_report = 1
		return
	}
	//fmt.Println(t, ned[0].Format("02.01.2006"))

	tt := float64(t)
	// усреднение недельных показателей
	if t > 0 {
		rs.Cars = int(float64(rs.Cars) / tt)
		rs.Vyh = int(float64(rs.Vyh) / tt)
		rs.Bezvod = int(float64(rs.Bezvod) / tt)
		rs.Dtp = int(float64(rs.Dtp) / tt)
		rs.Remont = int(float64(rs.Remont) / tt)
		rs.Naprodag = int(float64(rs.Naprodag) / tt)
		rs.Shtraff = int(float64(rs.Shtraff) / tt)
		rs.Keys = int(float64(rs.Keys) / tt)
		rs.Volonter = int(float64(rs.Volonter) / tt)
		rs.Spysan = int(float64(rs.Spysan) / tt)
		rs.Pomylka = int(float64(rs.Pomylka) / tt)
		rs.Car_work = int(float64(rs.Car_work) / tt)
		rs.Car_trip = int(float64(rs.Car_trip) / tt)
		rs.Dvaakk = int(float64(rs.Dvaakk) / tt)
		rs.Ecip = int(float64(rs.Ecip) / tt)
		//rs.Hour = rs.Hour / float64(t)
	}
	rs.day_report = t

	return
}

// КПД щтчет по городу за 1 день
func KPD_REPORT_SITY_DAY(sit sity, dy time.Time) (r sity_report) {

	// поездки за указанные дни
	report := readBD_Trip(sit, []time.Time{dy}, "*")

	// поездки по новому
	trips := get_Trips([]sity{sit}, []time.Time{dy})

	b := trips_Citys_days{}
	b.Create(sit, trips, dy)

	dl := get_ALL_DispetcherLog(sit.id_BGQ, dy.Format("2006-01-02"), dy.Add(24*time.Hour).Format("2006-01-02"))

	// пробег по часам за указанные дни
	trip_h := trip_hour([]time.Time{dy})

	// ищет запись по дате
	// и номеру машины
	// отдает [24]float32
	searh_car_ch := func(n string) ([24]float64, float64) {
		var res [24]float64
		var rs float64
		if n == "" {
			return res, 0
		}
		for _, i := range trip_h {
			if i.nomer == "" {
				continue
			}
			if n == i.nomer {
				res = i.hr
				break
			}
		}
		for _, i := range res {
			rs = rs + i
		}
		return res, rs
	}

	search_vod_day := func(n string) vod_day {
		var rs vod_day
		if len(report) == 0 {
			return rs
		}
		for _, vod := range report[0].res_work {
			if vod.name == n {
				return vod
			}
		}

		return rs
	}

	// просчет работы онлайн
	h_online := func(a [24][3]int, tr [24]int) (rs float64) {

		summ := 0
		for i := 0; i < 24; i++ {
			summ = summ + a[i][0] + a[i][1] + a[i][2]
		}

		for i := 0; i < 24; i++ {
			imp := a[i]
			if (imp[0] + imp[1] + imp[2]) > 0 {
				min := 60 / (imp[0] + imp[1] + imp[2]) * (imp[0] + imp[1])
				rs = rs + float64(min)
			}
		}
		rs = rs / 60
		if rs > 20 {
			rs = 0
		}

		return
	}

	cars := make([]string, 0)
	id_map := make([]int64, 0)

	// данные графика за указанный период
	// за один день, по одному городу
	graf := get_grafik(sit, []time.Time{dy})

	//n := 0

	for _, i := range graf {

		if i.car_nomer == "" {
			continue
		}

		// проверка на повтор авто в списке
		if !comp_string(i.car_nomer, cars) {
			cars = append(cars, i.car_nomer)
		} else {
			continue
		}

		// проверка на id_mapon в списке
		if !comp_int64(i.id_mapon, id_map) {
			id_map = append(id_map, i.id_mapon)
		} else {
			continue
		}
		// n++
		// fmt.Println(n, i.car_nomer, i.id_mapon)

		v1 := i.name_vod1
		v2 := i.name_vod2

		r.Sity = sit.name
		r.Cars++
		r.Day = dy.Format("02.01.2006")
		r.Data = dy

		stat := i.Status()
		tr1 := search_vod_day(v1)
		tr2 := search_vod_day(v2)
		_, pr := searh_car_ch(i.car_nomer)
		r.Probeg = r.Probeg + int(pr)
		var hour, hour1, hour2 float64

		switch {
		case stat == "Новий водій":
			r.Car_work++
			//tc := get_Trips_Graf(i)
			if strings.Contains(v1, "(UK) ") {
				r.Trip_Uklon = r.Trip_Uklon + tr1.trip_sum
				r.Cash_Uklon = r.Cash_Uklon + tr1.cash_summ
			} else {
				r.Trip_Bolt = r.Trip_Bolt + tr1.trip_sum
				r.Cash_Bolt = r.Cash_Bolt + tr1.cash_summ
				hour1 = h_online(tr1.impl, tr1.trip)
			}
			if strings.Contains(v2, "(UK) ") {
				r.Trip_Uklon = r.Trip_Uklon + tr2.trip_sum
				r.Cash_Uklon = r.Cash_Uklon + tr2.cash_summ
			} else {
				r.Trip_Bolt = r.Trip_Bolt + tr2.trip_sum
				r.Cash_Bolt = r.Cash_Bolt + tr2.cash_summ
				hour2 = h_online(tr2.impl, tr2.trip)
			}
			// просчет часов
			hour = hour1 + hour2
			r.Procent_Trip(tr1.trip_sum + tr2.trip_sum)
		case stat == "Екіпаж":
			r.Ecip++
			r.Car_work++
			if strings.Contains(v1, "(UK) ") {
				r.Trip_Uklon = r.Trip_Uklon + tr1.trip_sum
				r.Cash_Uklon = r.Cash_Uklon + tr1.cash_summ
			} else {
				r.Trip_Bolt = r.Trip_Bolt + tr1.trip_sum
				r.Cash_Bolt = r.Cash_Bolt + tr1.cash_summ
				hour1 = h_online(tr1.impl, tr1.trip)
			}
			if strings.Contains(v2, "(UK) ") {
				r.Trip_Uklon = r.Trip_Uklon + tr2.trip_sum
				r.Cash_Uklon = r.Cash_Uklon + tr2.cash_summ
			} else {
				r.Trip_Bolt = r.Trip_Bolt + tr2.trip_sum
				r.Cash_Bolt = r.Cash_Bolt + tr2.cash_summ
				hour2 = h_online(tr2.impl, tr2.trip)
			}
			// просчет часов
			hour = hour1 + hour2
			r.Procent_Trip(tr1.trip_sum + tr2.trip_sum)
		case stat == "Два аккаунти":
			//	i.Print()
			r.Dvaakk++
			r.Car_work++
			if strings.Contains(v1, "(UK) ") {
				//	fmt.Println("*********************")
				r.Trip_Uklon = r.Trip_Uklon + tr1.trip_sum
				r.Cash_Uklon = r.Cash_Uklon + tr1.cash_summ
			} else {
				r.Trip_Bolt = r.Trip_Bolt + tr1.trip_sum
				r.Cash_Bolt = r.Cash_Bolt + tr1.cash_summ
				hour1 = h_online(tr1.impl, tr1.trip)
			}
			if strings.Contains(v2, "(UK) ") {
				//		fmt.Println("------------------------")
				r.Trip_Uklon = r.Trip_Uklon + tr2.trip_sum
				r.Cash_Uklon = r.Cash_Uklon + tr2.cash_summ
			} else {
				r.Trip_Bolt = r.Trip_Bolt + tr2.trip_sum
				r.Cash_Bolt = r.Cash_Bolt + tr2.cash_summ
				hour2 = h_online(tr2.impl, tr2.trip)
			}
			// просчет часов
			hour = hour1 + hour2
			r.Procent_Trip(tr1.trip_sum + tr2.trip_sum)
		case stat == "Персональний":
			r.Car_work++
			if strings.Contains(v1, "(UK) ") {
				r.Trip_Uklon = r.Trip_Uklon + tr1.trip_sum
				r.Cash_Uklon = r.Cash_Uklon + tr1.cash_summ
			} else {
				r.Trip_Bolt = r.Trip_Bolt + tr1.trip_sum
				r.Cash_Bolt = r.Cash_Bolt + tr1.cash_summ
				hour1 = h_online(tr1.impl, tr1.trip)
			}
			if strings.Contains(v2, "(UK) ") {
				r.Trip_Uklon = r.Trip_Uklon + tr2.trip_sum
				r.Cash_Uklon = r.Cash_Uklon + tr2.cash_summ
			} else {
				r.Trip_Bolt = r.Trip_Bolt + tr2.trip_sum
				r.Cash_Bolt = r.Cash_Bolt + tr2.cash_summ
				hour2 = h_online(tr2.impl, tr2.trip)
			}
			// просчет часов
			hour = hour1 + hour2
			r.Procent_Trip(tr1.trip_sum + tr2.trip_sum)
		case stat == "Ремонт":
			r.Remont++
		case stat == "ДТП":
			r.Dtp++
		case stat == "Волонтер":
			r.Volonter++
		case stat == "Ключі":
			r.Keys++
		case stat == "Угон":
			r.Ugon++
		case stat == "Документи":
			r.Doki++
		case stat == "Штрафмайданчик":
			r.Shtraff++
		case stat == "Перегін":
			r.Peregon++
		case stat == "Списання":
			r.Spysan++
		case stat == "Оренда":
			r.Car_work++
			if strings.Contains(v1, "(UK) ") {
				r.Trip_Uklon = r.Trip_Uklon + tr1.trip_sum
				r.Cash_Uklon = r.Cash_Uklon + tr1.cash_summ
			} else {
				r.Trip_Bolt = r.Trip_Bolt + tr1.trip_sum
				r.Cash_Bolt = r.Cash_Bolt + tr1.cash_summ
				hour1 = h_online(tr1.impl, tr1.trip)
			}
			if strings.Contains(v2, "(UK) ") {
				r.Trip_Uklon = r.Trip_Uklon + tr2.trip_sum
				r.Cash_Uklon = r.Cash_Uklon + tr2.cash_summ
			} else {
				r.Trip_Bolt = r.Trip_Bolt + tr2.trip_sum
				r.Cash_Bolt = r.Cash_Bolt + tr2.cash_summ
				hour2 = h_online(tr2.impl, tr2.trip)
			}
			// просчет часов
			hour = hour1 + hour2
			r.Procent_Trip(tr1.trip_sum + tr2.trip_sum)
		case stat == "Вихідний":
			r.Car_work++
			r.Vyh++
			if strings.Contains(v1, "(UK) ") {
				r.Trip_Uklon = r.Trip_Uklon + tr1.trip_sum
				r.Cash_Uklon = r.Cash_Uklon + tr1.cash_summ
			} else {
				r.Trip_Bolt = r.Trip_Bolt + tr1.trip_sum
				r.Cash_Bolt = r.Cash_Bolt + tr1.cash_summ
				hour1 = h_online(tr1.impl, tr1.trip)
			}
			if strings.Contains(v2, "(UK) ") {
				r.Trip_Uklon = r.Trip_Uklon + tr2.trip_sum
				r.Cash_Uklon = r.Cash_Uklon + tr2.cash_summ
			} else {
				r.Trip_Bolt = r.Trip_Bolt + tr2.trip_sum
				r.Cash_Bolt = r.Cash_Bolt + tr2.cash_summ
				hour2 = h_online(tr2.impl, tr2.trip)
			}
			// просчет часов
			hour = hour1 + hour2
			r.Procent_Trip(tr1.trip_sum + tr2.trip_sum)
		case stat == "Лікарняний":
			r.Car_work++
			r.Vyh++
			if strings.Contains(v1, "(UK) ") {
				r.Trip_Uklon = r.Trip_Uklon + tr1.trip_sum
				r.Cash_Uklon = r.Cash_Uklon + tr1.cash_summ
			} else {
				r.Trip_Bolt = r.Trip_Bolt + tr1.trip_sum
				r.Cash_Bolt = r.Cash_Bolt + tr1.cash_summ
				hour1 = h_online(tr1.impl, tr1.trip)
			}
			if strings.Contains(v2, "(UK) ") {
				r.Trip_Uklon = r.Trip_Uklon + tr2.trip_sum
				r.Cash_Uklon = r.Cash_Uklon + tr2.cash_summ
			} else {
				r.Trip_Bolt = r.Trip_Bolt + tr2.trip_sum
				r.Cash_Bolt = r.Cash_Bolt + tr2.cash_summ
				hour2 = h_online(tr2.impl, tr2.trip)
			}
			// просчет часов
			hour = hour1 + hour2
			r.Procent_Trip(tr1.trip_sum + tr2.trip_sum)
		case stat == "Без водія":
			r.Bezvod++
		case stat == "Некоректний Вихідний":
			r.Car_work++
			r.Pomylka++
		case stat == "Помилка":
			r.Pomylka++
			//default:
		}

		r.Hour = r.Hour + hour

		if r.Bezvod < 0 {
			r.Bezvod = 0
		}

	}

	r.Trip_Bolt = b.bolt
	r.Trip_Uklon = b.uklon
	r.Trip_Uber = b.uber
	r.Trip = b.trips

	r.Cash_Bolt = int(b.cash_bolt)
	r.Cash_Uklon = int(b.cash_uklon)
	r.Cash_Uber = int(b.cash_uber)
	r.Cash = int(b.cash)

	r.Hour = dl.work
	if dl.work > 0 {
		r.Tipbolt_hour = float64(b.bolt) / dl.work
		r.Waitbolt = dl.wait / dl.work
	}

	r.Cash_Bolt = int(float64(r.Cash_Bolt) * float64(100-sit.procent_bolt) * 0.45 / 100)
	r.Cash_Uklon = int(float64(r.Cash_Uklon) * float64(100-sit.procent_uclon) * 0.45 / 100)
	r.Cash_Uber = int(float64(r.Cash_Uber) * float64(100-sit.procent_uclon) * 0.45 / 100)
	r.day_report = 1
	fmt.Println(r.Sity, r.Cars, r.Bezvod, r.Car_work, r.Car_trip)
	return
}

// отчет для тараса 3
func KPD_report_read_A() {

	control := control_execute{}
	control.data = time.Now().Format("02.01.2006 15-04")
	control.name = "KPD_report()"
	control.coment = "completed"

	if control.read() {
		fmt.Println("KPD Report completed successfully")
		return
	}

	sitys := create_sitys()

	res := make([][]sity_report, 0)
	resd := make([][]sity_report, 0)

	trvr := func(i sity) {

		fmt.Print("Формирование отчета по кпд ")
		rs := travel_rep(i, last_days(7), false)
		res = append(res, rs)
	}

	for _, i := range sitys {

		trvr(i)
	}

	for n := 0; n < 7; n++ {
		rsd := make([]sity_report, 0)
		for _, i := range res {

			for nn, j := range i {
				if nn != n {
					continue
				}
				rsd = append(rsd, j)
			}

		}
		resd = append(resd, rsd)
	}

	// сортировка по поездкам
	for _, rs := range resd {
		sort.Slice(rs, func(i, j int) (less bool) {
			return rs[i].Trip > rs[j].Trip
		})
	}

	str := make([][]string, 0)
	str = append(str, []string{"Звіт по перепробігу сформовано: " + time.Now().Format("02.01.2006 15:04")})

	st := make([]string, 0)
	st = append(st, "Назва міста")

	for _, i := range res[0] {
		st = append(st, i.Day)
	}
	str = append(str, st)

	for _, i := range res {
		st := make([]string, 0)
		st = append(st, i[0].Sity)
		for _, j := range i {
			st = append(st, strconv.Itoa(j.Preprobeg))
		}
		str = append(str, st)
	}

	str1 := make([][]string, 0)
	str1 = append(str1, []string{"Звіт по перепробігу сформовано: " + time.Now().Format("02.01.2006 15:04")})

	st1 := make([]string, 0)
	st1 = append(st1, "Назва міста", "Авто", "Перепробіг")
	str1 = append(str1, st1)

	for _, i := range res {
		for n, j := range i {
			if n > 0 {
				continue
			}
			str1 = append(str1, []string{j.Day})

		}
	}

	str2 := make([][]string, 0)
	str2 = append(str2, []string{"Звіт КПД сформований : " + time.Now().Format("02.01.2006 15:04"), "", "", "", "", "", "", "", "", "", "дохід автопарку на 1 авто з поїздками", "", "не у роботі", "", "", "", "", "", "", "", "", "", "", "", "", "", ""})

	st2 := make([]string, 0)
	st2 = append(st2, "Назва міста", "% виконання плану на усі авто", "% виконання плану на авто у роботі", "Поїздок", "Авто всього", "Авто у роботі", "Авто з поїздками")
	st2 = append(st2, "Поїздок на годину", "Час онлайн по парку усього", "Середній час онлайн", "за 1 годину", "за добу")
	st2 = append(st2, "Вихідний", "Без водія", "Ремонт", "ДТП", "Евакуація")
	st2 = append(st2, "Поїздок на авто у роботі", "Поїздок на усі авто")
	st2 = append(st2, "% 0 поїздок", "% 1 - 5 поїздок", "% 6 - 10 поїздок", "%  < 10 поїздок", "% 11 - 15 поїздок", "% 16 - 20 поїздок", "% < 20 поїздок", "% 21 - 25 поїздок", "% > 25 поїздок", "% > 20 поїздок")
	str2 = append(str2, st2)

	for n, i := range resd {
		str2 = append(str2, []string{"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""})
		str2 = append(str2, []string{resd[n][0].Day, "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""})

		for _, j := range i {
			//fmt.Println(j)

			st := make([]string, 0)
			st = append(st, j.Sity)
			if j.Car_work*norma_trips > 0 && j.Cars*norma_trips > 0 {
				st = append(st, comm_krapka(fmt.Sprint(j.Trip*100/(j.Cars*norma_trips))))
				st = append(st, comm_krapka(fmt.Sprint(j.Trip*100/(j.Car_work*norma_trips))))
			} else {
				st = append(st, "")
				st = append(st, "")
			}
			st = append(st, fmt.Sprint(j.Trip))
			st = append(st, fmt.Sprint(j.Cars))
			st = append(st, fmt.Sprint(j.Car_work))
			st = append(st, fmt.Sprint(j.Car_trip))
			if j.Hour > 0 && j.Car_work > 0 && j.Car_trip > 0 {
				st = append(st, comm_krapka(fmt.Sprintf("%.1f", float64(j.Trip)/j.Hour)))
				st = append(st, comm_krapka(fmt.Sprintf("%.1f", j.Hour)))
				st = append(st, comm_krapka(fmt.Sprintf("%.1f", j.Hour/float64(j.Car_trip))))
				st = append(st, comm_krapka(fmt.Sprintf("%.1f", 0.45*float64(j.Cash)/j.Hour)))
				st = append(st, comm_krapka(fmt.Sprintf("%.1f", 0.45*float64(j.Cash)/float64(j.Car_trip))))
			} else {
				st = append(st, "", "")
				st = append(st, "", "")
				st = append(st, "")
			}

			st = append(st, fmt.Sprint(j.Vyh))
			st = append(st, fmt.Sprint(j.Bezvod))
			st = append(st, fmt.Sprint(j.Remont))
			st = append(st, fmt.Sprint(j.Dtp))
			st = append(st, fmt.Sprint(j.Evak))

			st = append(st, comm_krapka(fmt.Sprintf("%.1f", j.Trip_car_trip)))
			st = append(st, comm_krapka(fmt.Sprintf("%.1f", j.Trip_car_work)))

			if j.Car_work > 0 && j.Car_trip > 0 {

				st = append(st, comm_krapka(fmt.Sprintf("%.1f", float64(j.Car_0_trip)/float64(j.Car_work)*100)))
				st = append(st, comm_krapka(fmt.Sprintf("%.1f", float64(j.Car_5_trip)/float64(j.Car_work)*100)))
				st = append(st, comm_krapka(fmt.Sprintf("%.1f", float64(j.Car_10_trip)/float64(j.Car_work)*100)))
				st = append(st, comm_krapka(fmt.Sprintf("%.1f", float64(j.Car_10_trip+j.Car_5_trip+j.Car_0_trip)/float64(j.Car_work)*100)))
				st = append(st, comm_krapka(fmt.Sprintf("%.1f", float64(j.Car_15_trip)/float64(j.Car_work)*100)))
				st = append(st, comm_krapka(fmt.Sprintf("%.1f", float64(j.Car_20_trip)/float64(j.Car_work)*100)))
				st = append(st, comm_krapka(fmt.Sprintf("%.1f", float64(j.Car_10_trip+j.Car_5_trip+j.Car_15_trip+j.Car_20_trip+j.Car_0_trip)/float64(j.Car_work)*100)))
				st = append(st, comm_krapka(fmt.Sprintf("%.1f", float64(j.Car_25_trip)/float64(j.Car_work)*100)))
				st = append(st, comm_krapka(fmt.Sprintf("%.1f", float64(j.Car_99_trip)/float64(j.Car_work)*100)))
				st = append(st, comm_krapka(fmt.Sprintf("%.1f", float64(j.Car_99_trip+j.Car_25_trip)/float64(j.Car_work)*100)))

			} else {
				st = append(st, "", "", "", "", "", "", "", "", "", "")
			}

			str2 = append(str2, st)
		}

	}

	//time.Sleep(20 * time.Minute)

	// запись суммы перепробегов по паркам
	trs.Rec("18hK0457uyVcfDQu9N5fhKzreHIfAiFkFRPVevVxiw8o", "'Pereprobeg'!A1", str, true, true)
	// запись перепробегов по номерам машин
	trs.Rec("18hK0457uyVcfDQu9N5fhKzreHIfAiFkFRPVevVxiw8o", "'Pld'!A1", str1, true, true)
	// запись КПД
	trs.Rec("18hK0457uyVcfDQu9N5fhKzreHIfAiFkFRPVevVxiw8o", "'КПД'!A1", str2, true, true)

	for _, i := range res {

		str2 := make([][]string, 0)
		str2 = append(str2, []string{"Звіт КПД сформований : " + time.Now().Format("02.01.2006 15:04"), "", "", "", "", "", "", "", "", "", "дохід автопарку на 1 авто з поїздками", "", "не у роботі", "", "", "", "", "", "", "", "", "", "", "", "", "", ""})

		st2 := make([]string, 0)
		st2 = append(st2, i[0].Sity, "% виконання плану на усі авто", "% виконання плану на авто у роботі", "Поїздок", "Авто всього", "Авто у роботі", "Авто з поїздками")
		st2 = append(st2, "Поїздок на годину", "Час онлайн по парку усього", "Середній час онлайн", "за 1 годину", "за добу")
		st2 = append(st2, "Вихідний", "Без водія", "Ремонт", "ДТП", "Евакуація")
		st2 = append(st2, "Поїздок на авто у роботі", "Поїздок на усі авто")
		st2 = append(st2, "% 0 поїздок", "% 1 - 5 поїздок", "% 6 - 10 поїздок", "%  < 10 поїздок", "% 11 - 15 поїздок", "% 16 - 20 поїздок", "% < 20 поїздок", "% 21 - 25 поїздок", "% > 25 поїздок", "% > 20 поїздок")
		str2 = append(str2, st2)

		itogi := itogi_sity_report(i)

		i = append(i, itogi)

		for _, j := range i {

			st := make([]string, 0)
			st = append(st, krapka_tire(j.Day))
			if j.Car_work*norma_trips > 0 && j.Cars*norma_trips > 0 {
				st = append(st, comm_krapka(fmt.Sprint(j.Trip*100/(j.Cars*norma_trips))))
				st = append(st, comm_krapka(fmt.Sprint(j.Trip*100/(j.Car_work*norma_trips))))
			} else {
				st = append(st, "")
				st = append(st, "")
			}
			st = append(st, fmt.Sprint(j.Trip))
			st = append(st, fmt.Sprint(j.Cars))
			st = append(st, fmt.Sprint(j.Car_work))
			st = append(st, fmt.Sprint(j.Car_trip))
			if j.Hour > 0 && j.Car_work > 0 && j.Car_trip > 0 {
				st = append(st, comm_krapka(fmt.Sprintf("%.1f", float64(j.Trip)/j.Hour)))
				st = append(st, comm_krapka(fmt.Sprintf("%.1f", j.Hour)))
				st = append(st, comm_krapka(fmt.Sprintf("%.1f", j.Hour/float64(j.Car_trip))))
				st = append(st, comm_krapka(fmt.Sprintf("%.1f", 0.45*float64(j.Cash)/j.Hour)))
				st = append(st, comm_krapka(fmt.Sprintf("%.1f", 0.45*float64(j.Cash)/float64(j.Car_trip))))
			} else {
				st = append(st, "", "")
				st = append(st, "", "")
				st = append(st, "")
			}

			st = append(st, fmt.Sprint(j.Vyh))
			st = append(st, fmt.Sprint(j.Bezvod))
			st = append(st, fmt.Sprint(j.Remont))
			st = append(st, fmt.Sprint(j.Dtp))
			st = append(st, fmt.Sprint(j.Evak))

			st = append(st, comm_krapka(fmt.Sprintf("%.1f", j.Trip_car_trip)))
			st = append(st, comm_krapka(fmt.Sprintf("%.1f", j.Trip_car_work)))

			if j.Car_work > 0 && j.Car_trip > 0 {

				st = append(st, comm_krapka(fmt.Sprintf("%.1f", float64(j.Car_0_trip)/float64(j.Car_work)*100)))
				st = append(st, comm_krapka(fmt.Sprintf("%.1f", float64(j.Car_5_trip)/float64(j.Car_work)*100)))
				st = append(st, comm_krapka(fmt.Sprintf("%.1f", float64(j.Car_10_trip)/float64(j.Car_work)*100)))
				st = append(st, comm_krapka(fmt.Sprintf("%.1f", float64(j.Car_10_trip+j.Car_5_trip+j.Car_0_trip)/float64(j.Car_work)*100)))
				st = append(st, comm_krapka(fmt.Sprintf("%.1f", float64(j.Car_15_trip)/float64(j.Car_work)*100)))
				st = append(st, comm_krapka(fmt.Sprintf("%.1f", float64(j.Car_20_trip)/float64(j.Car_work)*100)))
				st = append(st, comm_krapka(fmt.Sprintf("%.1f", float64(j.Car_10_trip+j.Car_5_trip+j.Car_15_trip+j.Car_20_trip+j.Car_0_trip)/float64(j.Car_work)*100)))
				st = append(st, comm_krapka(fmt.Sprintf("%.1f", float64(j.Car_25_trip)/float64(j.Car_work)*100)))
				st = append(st, comm_krapka(fmt.Sprintf("%.1f", float64(j.Car_99_trip)/float64(j.Car_work)*100)))
				st = append(st, comm_krapka(fmt.Sprintf("%.1f", float64(j.Car_99_trip+j.Car_25_trip)/float64(j.Car_work)*100)))

			} else {
				st = append(st, "", "", "", "", "", "", "", "", "", "")
			}

			str2 = append(str2, st)
		}

		list := "'КПД_" + i[0].Sity + "'!A1"
		//fmt.Println(list)

		// запись КПД
		trs.Rec("18hK0457uyVcfDQu9N5fhKzreHIfAiFkFRPVevVxiw8o", list, str2, true, true)

	}

	fmt.Println("END")
	control.save()
}

// просчет отчета города за период
func itogi_sity_report(a []sity_report) sity_report {
	var res sity_report
	ln := len(a)
	if ln == 0 {
		return res
	}
	res.Sity = a[0].Sity
	res.Day = "За період:"
	for _, i := range a {
		res.Trip = res.Trip + i.Trip
		res.Hour = res.Hour + i.Hour
		res.Probeg = res.Probeg + i.Probeg
		res.Preprobeg = res.Preprobeg + i.Preprobeg
		res.Cash = res.Cash + i.Cash
		res.Ecip = res.Ecip + i.Ecip
		res.Vyh = res.Vyh + i.Vyh
		res.Bezvod = res.Bezvod + i.Bezvod
		res.Remont = res.Remont + i.Remont
		res.Dtp = res.Dtp + i.Dtp
		res.Evak = res.Evak + i.Evak
		res.Cars = res.Cars + i.Cars
		res.Car_work = res.Car_work + i.Car_work
		res.Car_trip = res.Car_trip + i.Car_trip
		res.Trip_car_work = res.Trip_car_work + i.Trip_car_work
		res.Trip_car_trip = res.Trip_car_trip + i.Trip_car_trip
		res.Car_0_trip = res.Car_0_trip + i.Car_0_trip
		res.Car_5_trip = res.Car_5_trip + i.Car_5_trip
		res.Car_10_trip = res.Car_10_trip + i.Car_10_trip
		res.Car_15_trip = res.Car_15_trip + i.Car_15_trip
		res.Car_20_trip = res.Car_20_trip + i.Car_20_trip
		res.Car_25_trip = res.Car_25_trip + i.Car_25_trip
		res.Car_99_trip = res.Car_99_trip + i.Car_99_trip
	}
	res.Trip = res.Trip / ln
	res.Hour = res.Hour / float64(ln)
	res.Cash = res.Cash / ln
	res.Bezvod = res.Bezvod / ln
	res.Remont = res.Remont / ln
	res.Dtp = res.Dtp / ln
	res.Evak = res.Evak / ln
	res.Cars = res.Cars / ln
	res.Car_work = res.Car_work / ln
	res.Car_trip = res.Car_trip / ln
	res.Trip_car_work = res.Trip_car_work / float64(ln)
	res.Trip_car_trip = res.Trip_car_trip / float64(ln)
	res.Car_0_trip = res.Car_0_trip / ln
	res.Car_5_trip = res.Car_5_trip / ln
	res.Car_10_trip = res.Car_10_trip / ln
	res.Car_15_trip = res.Car_15_trip / ln
	res.Car_20_trip = res.Car_20_trip / ln
	res.Car_25_trip = res.Car_25_trip / ln
	res.Car_99_trip = res.Car_99_trip / ln

	return res
}

// просчет отчета города за период
func summa_sity_report(a []sity_report, day []time.Time, titul string) sity_report {
	var res sity_report
	ln := len(a)
	if ln == 0 {
		return res
	}
	res.Sity = titul
	if len(day) == 1 {
		res.Day = "За " + day[0].Format("02.01.2006") + " число"
	} else {
		res.Day = "За період з " + day[0].Format("02.01.2006") + " по " + day[len(day)-1].Format("02.01.2006")
	}
	res.day_report = ln
	for _, i := range a {
		res.Trip = res.Trip + i.Trip
		res.Trip_Bolt = res.Trip_Bolt + i.Trip_Bolt
		res.Trip_Uklon = res.Trip_Uklon + i.Trip_Uklon
		res.Hour = res.Hour + i.Hour
		res.Probeg = res.Probeg + i.Probeg
		res.Preprobeg = res.Preprobeg + i.Preprobeg
		res.Cash = res.Cash + i.Cash
		res.Cash_Bolt = res.Cash_Bolt + i.Cash_Bolt
		res.Cash_Uklon = res.Cash_Uklon + i.Cash_Uklon
		res.Ecip = res.Ecip + i.Ecip
		res.Vyh = res.Vyh + i.Vyh
		res.Bezvod = res.Bezvod + i.Bezvod
		res.Remont = res.Remont + i.Remont
		res.Dvaakk = res.Dvaakk + i.Dvaakk
		res.Dtp = res.Dtp + i.Dtp
		res.Naprodag = res.Naprodag + i.Naprodag
		res.Shtraff = res.Shtraff + i.Shtraff
		res.Doki = res.Doki + i.Doki
		res.Keys = res.Keys + i.Keys
		res.Ugon = res.Ugon + i.Ugon
		res.Volonter = res.Volonter + i.Volonter
		res.Spysan = res.Spysan + i.Spysan
		res.Evak = res.Evak + i.Evak
		res.Cars = res.Cars + i.Cars
		res.Car_work = res.Car_work + i.Car_work
		res.Car_trip = res.Car_trip + i.Car_trip
		res.Trip_car_work = res.Trip_car_work + i.Trip_car_work
		res.Trip_car_trip = res.Trip_car_trip + i.Trip_car_trip
		res.Car_0_trip = res.Car_0_trip + i.Car_0_trip
		res.Car_5_trip = res.Car_5_trip + i.Car_5_trip
		res.Car_10_trip = res.Car_10_trip + i.Car_10_trip
		res.Car_15_trip = res.Car_15_trip + i.Car_15_trip
		res.Car_20_trip = res.Car_20_trip + i.Car_20_trip
		res.Car_25_trip = res.Car_25_trip + i.Car_25_trip
		res.Car_99_trip = res.Car_99_trip + i.Car_99_trip
	}

	return res
}

// просчет отчета города за период
func average_sity_report(a []sity_report, day []time.Time, titul string) sity_report {
	var res sity_report
	ln := len(a)
	if ln == 0 {
		return res
	}
	res.Sity = titul
	if len(day) == 1 {
		res.Day = "За " + day[0].Format("02.01.2006") + " число"
	} else {
		res.Day = "За період з " + day[0].Format("02.01.2006") + " по " + day[len(day)-1].Format("02.01.2006")
	}
	for _, i := range a {
		res.Trip = res.Trip + i.Trip
		res.Trip_Bolt = res.Trip_Bolt + i.Trip_Bolt
		res.Trip_Uklon = res.Trip_Uklon + i.Trip_Uklon
		res.Hour = res.Hour + i.Hour
		res.Probeg = res.Probeg + i.Probeg
		res.Preprobeg = res.Preprobeg + i.Preprobeg
		res.Cash = res.Cash + i.Cash
		res.Cash_Bolt = res.Cash_Bolt + i.Cash_Bolt
		res.Cash_Uklon = res.Cash_Uklon + i.Cash_Uklon
		res.Ecip = res.Ecip + i.Ecip
		res.Vyh = res.Vyh + i.Vyh
		res.Bezvod = res.Bezvod + i.Bezvod
		res.Remont = res.Remont + i.Remont
		res.Dvaakk = res.Dvaakk + i.Dvaakk
		res.Dtp = res.Dtp + i.Dtp
		res.Naprodag = res.Naprodag + i.Naprodag
		res.Shtraff = res.Shtraff + i.Shtraff
		res.Doki = res.Doki + i.Doki
		res.Keys = res.Keys + i.Keys
		res.Ugon = res.Ugon + i.Ugon
		res.Volonter = res.Volonter + i.Volonter
		res.Spysan = res.Spysan + i.Spysan
		res.Evak = res.Evak + i.Evak
		res.Cars = res.Cars + i.Cars
		res.Car_work = res.Car_work + i.Car_work
		res.Car_trip = res.Car_trip + i.Car_trip
		res.Trip_car_work = res.Trip_car_work + i.Trip_car_work
		res.Trip_car_trip = res.Trip_car_trip + i.Trip_car_trip
		res.Car_0_trip = res.Car_0_trip + i.Car_0_trip
		res.Car_5_trip = res.Car_5_trip + i.Car_5_trip
		res.Car_10_trip = res.Car_10_trip + i.Car_10_trip
		res.Car_15_trip = res.Car_15_trip + i.Car_15_trip
		res.Car_20_trip = res.Car_20_trip + i.Car_20_trip
		res.Car_25_trip = res.Car_25_trip + i.Car_25_trip
		res.Car_99_trip = res.Car_99_trip + i.Car_99_trip
	}
	res.Trip = res.Trip / ln
	res.Trip_Bolt = res.Trip_Bolt / ln
	res.Trip_Uklon = res.Trip_Uklon / ln
	res.Hour = res.Hour / float64(ln)
	res.Probeg = res.Probeg / ln
	res.Preprobeg = res.Preprobeg / ln
	res.Cash = res.Cash / ln
	res.Cash_Bolt = res.Cash_Bolt / ln
	res.Cash_Uklon = res.Cash_Uklon / ln
	res.Ecip = res.Ecip / ln
	res.Vyh = res.Vyh / ln
	res.Bezvod = res.Bezvod / ln
	res.Remont = res.Remont / ln
	res.Dvaakk = res.Dvaakk / ln
	res.Dtp = res.Dtp / ln
	res.Naprodag = res.Naprodag / ln
	res.Shtraff = res.Shtraff / ln
	res.Doki = res.Doki / ln
	res.Keys = res.Keys / ln
	res.Ugon = res.Ugon / ln
	res.Volonter = res.Volonter / ln
	res.Spysan = res.Spysan / ln
	res.Evak = res.Evak / ln
	res.Cars = res.Cars / ln
	res.Car_work = res.Car_work / ln
	res.Car_trip = res.Car_trip / ln
	res.Trip_car_work = res.Trip_car_work / float64(ln)
	res.Trip_car_trip = res.Trip_car_trip / float64(ln)
	res.Car_0_trip = res.Car_0_trip / ln
	res.Car_5_trip = res.Car_5_trip / ln
	res.Car_10_trip = res.Car_10_trip / ln
	res.Car_15_trip = res.Car_15_trip / ln
	res.Car_20_trip = res.Car_20_trip / ln
	res.Car_25_trip = res.Car_25_trip / ln
	res.Car_99_trip = res.Car_99_trip / ln

	return res
}

// процент водителей
// сделавших 25 и 30 поездок
func Report_Procent_Trips(dn int) {
	res := make([][]string, 0)
	day := last_days(dn)
	r := READ_KPD_DAY_SITY_DB("*", day)
	rd := make([]string, 0)
	rd = append(rd, "Date")
	for _, i := range day {
		rd = append(rd, i.Format(TNS))
	}
	fmt.Println(rd)
	res = append(res, rd)
	for _, i := range Citys.Create_sitys_List() {
		r25 := make([]string, 61)
		r30 := make([]string, 61)
		rTrip := make([]string, 61)
		rwr := make([]string, 61)
		rwt := make([]string, 61)
		rtc := make([]string, 61)
		r25[0] = i.Name + "|25"
		r30[0] = i.Name + "|30"
		rTrip[0] = i.Name + "|Trip"
		rwr[0] = i.Name + "|Rwr"
		rwt[0] = i.Name + "|Rwt"
		rtc[0] = i.Name + "|Rtc"
		for _, j := range r {
			if i.Name != j.Sity {
				continue
			}
			for n, d := range rd {
				if d == j.Day {
					r25[n] = fmt.Sprint(j.Car_25_trip)
					r30[n] = fmt.Sprint(j.Car_99_trip)
					rTrip[n] = fmt.Sprint(j.Trip)
					rwr[n] = fmt.Sprint(j.Car_work)
					rwt[n] = fmt.Sprint(j.Car_trip)
					rtc[n] = trs.Float_to_string_2(float64(j.Trip) / float64(j.Car_work))
				}
			}
		}
		res = append(res, r25)
		res = append(res, r30)
		res = append(res, rTrip)
		res = append(res, rwr)
		res = append(res, rwt)
		res = append(res, rtc)
		fmt.Println(r25)
		fmt.Println(r30)
		fmt.Println(rTrip)
		fmt.Println(rwr)
		fmt.Println(rwt)
		fmt.Println(rtc)
	}

	// за последние 9 недель

	day = last_ned(9)
	r = READ_KPD_DAY_SITY_DB("*", day)
	rd = make([]string, 0)
	rd = append(rd, "Ned")
	for _, i := range day {
		_, ned := i.ISOWeek()
		if !comp_string(fmt.Sprint(ned), rd) {
			rd = append(rd, fmt.Sprint(ned))
		}
	}
	fmt.Println(rd)
	res = append(res, rd)
	for _, i := range Citys.Create_sitys_List() {
		rn := make([]string, 10)
		rn[0] = i.Name + "|ned"
		for n, ned := range rd {
			if n == 0 {
				continue
			}
			trip := 0
			for _, tr := range r {
				if i.Name != tr.Sity {
					continue
				}
				dt, _ := time.Parse(TNS, tr.Day)
				_, nn := dt.ISOWeek()
				if fmt.Sprint(nn) == ned {
					trip = trip + tr.Trip
				}
			}
			rn[n] = fmt.Sprint(trip)
		}
		fmt.Println(rn)
		res = append(res, rn)
	}

	trs.Rec_Clear("1x8bp9hRXzGF-iSdv5EtyVIRQfcUT9A7poyhoA74vWi4", "'QBdata'!A1", res, false, true)
}
