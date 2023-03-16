package main

import (
	"database/sql"
	"fmt"
	"sort"
	"strings"
	"time"
	"trs"
)

var Vidhileno_Str = []string{"", "", "", "", "", "", "", "", "", "", "На розгляді"}
var Tracking_Empty = []string{"", "", "", "Ні", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""}
var Archive_Empty = []string{"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""}

const Norma_Trip_Referal = 500

// запуск в работу
// проверяет поданых на регистрацию
// пишет в базу прошедших проверку
// проверяет добавление аккаунтов и запись в архив
func Treatment_Referals(timing_1, timing_2, timing_3 int) {

	db, err := sql.Open("sqlite3", path+"sqlite/referal.db")
	if err != nil {
		fmt.Println(path+"sqlite/referal.db", err)
	}
	defer db.Close()

	// первый запуск
	Referals_Register(db)
	Treatment_Referals_Tracing(db)

	ticker1 := time.NewTicker(time.Duration(timing_1) * time.Second)
	ticker2 := time.NewTicker(time.Duration(timing_2) * time.Second)
	//ticker3 := time.NewTicker(time.Duration(timing_3) * time.Second)

	for {
		select {
		case <-ticker1.C:
			Referals_Register(db)
		case <-ticker2.C:
			Treatment_Referals_Tracing(db)
			//case <-ticker3.C:

		}
	}
}

// пишет на пустую строку на лист
func (a *Referal) Rec_Tracking_List() {
	list, err := trs.Read_sheets_CASHBOX_err("1_EWAl5qFMMWh0xbr3gzJ_KM87QNXhm7wJy89Vrs_zDs", "'Tracing'!A3:U", 0)
	if err != nil {
		fmt.Println(err)
		return
	}
	for n, i := range list {
		if i[0] == "" && i[1] == "" {
			trs.Rec("1_EWAl5qFMMWh0xbr3gzJ_KM87QNXhm7wJy89Vrs_zDs", fmt.Sprintf("'Tracing'!A%d", n+3), [][]string{a.SliseString()}, false, false)
			//fmt.Println(n, i)
			break
		}
	}
}

// обновляет инфо на листе отслеживания
func (a *Referal) Upgrade_Tracking_List() {
	list, err := trs.Read_sheets_CASHBOX_err("1_EWAl5qFMMWh0xbr3gzJ_KM87QNXhm7wJy89Vrs_zDs", "'Tracing'!A3:U", 0)
	if err != nil {
		fmt.Println(err)
		return
	}
	for n, i := range list {
		if a.City == i[0] && a.Referal_Pib == i[9] {
			trs.Rec("1_EWAl5qFMMWh0xbr3gzJ_KM87QNXhm7wJy89Vrs_zDs", fmt.Sprintf("'Tracing'!A%d", n+3), [][]string{a.SliseString()}, false, false)
			//fmt.Println(n, i)
			break
		}
	}
}

// обновляет инфо на листе отслеживания
// пустая строка
func (a *Referal) Upgrade_Tracking_List_Clear() {
	list, err := trs.Read_sheets_CASHBOX_err("1_EWAl5qFMMWh0xbr3gzJ_KM87QNXhm7wJy89Vrs_zDs", "'Tracing'!A3:U", 0)
	if err != nil {
		fmt.Println(err)
		return
	}
	for n, i := range list {
		if a.City == i[0] && a.Referal_Pib == i[9] {
			str := make([]string, 21)
			str[19] = "Ні"
			str[20] = "Ні"
			trs.Rec("1_EWAl5qFMMWh0xbr3gzJ_KM87QNXhm7wJy89Vrs_zDs", fmt.Sprintf("'Tracing'!A%d", n+3), [][]string{str}, false, false)
			break
		}
	}
}

func Referal_morning(ignor_control bool) {

	control := control_execute{}
	control.data = time.Now().Format("02.01.2006 15-04")
	control.name = "Referal_morning()()"
	control.coment = "completed"

	if control.read() {
		fmt.Println("Referal_morning()() completed successfully")
		if !ignor_control {
			return
		}
	}

	db, err := sql.Open("sqlite3", path+"sqlite/referal.db")
	if err != nil {
		fmt.Println(path+"sqlite/referal.db", err)
	}
	defer db.Close()
	Referals_Tracking(db)
	Referals_Tracking_Archive(db)
	control.save()
}

// проверка  рефералов в базе
// запись на лист Tracing
func Referals_Tracking(db *sql.DB) {

	r := Referals{}
	r.Read(db, "Activ")

	sort.Slice(r, func(i, j int) (less bool) {
		return r[i].Date_Start.Before(r[j].Date_Start)
	})

	sort.Slice(r, func(i, j int) (less bool) {
		return r[i].City > r[j].City
	})

	res := make([][]string, 0)
	for _, i := range r {
		i.Trips(db)
		res = append(res, i.SliseString())
	}
	if len(res) == 0 {
		for i := 0; i < 200; i++ {
			res = append(res, Tracking_Empty)
		}
	}

	trs.Rec_Clear("1_EWAl5qFMMWh0xbr3gzJ_KM87QNXhm7wJy89Vrs_zDs", fmt.Sprintf("'Tracing'!A%d", 3), res, false, true)
}

// проверка  рефералов в базе
// запись на лист Archive
func Referals_Tracking_Archive(db *sql.DB) {
	r := Referals{}
	r.Read(db, "Archive")
	res := make([][]string, 0)
	for _, i := range r {
		res = append(res, i.SliseString_Archive())
	}
	if len(res) == 0 {
		res = append(res, Archive_Empty)
	}
	trs.Rec_Clear("1_EWAl5qFMMWh0xbr3gzJ_KM87QNXhm7wJy89Vrs_zDs", fmt.Sprintf("'Archive'!A%d", 3), res, false, true)
}

// чистит лист регистрации в час ночи
// по понедельникам
func Clear_Referals_Register() {
	for {
		t := time.Now()
		if t.Weekday() == 1 && t.Hour() == 1 && t.Minute() == 0 {
			res := make([][]string, 0)
			for i := 0; i < 200; i++ {
				res = append(res, Vidhileno_Str)
			}
			trs.Rec_Clear("1_EWAl5qFMMWh0xbr3gzJ_KM87QNXhm7wJy89Vrs_zDs", fmt.Sprintf("'Register'!A%d", 3), res, false, false)
		}
	}
}

// проверяет лист сопровождения раз в час
func Treatment_Referals_Tracing(db *sql.DB) {

	//fmt.Println("Treatment_Referals_Tracing")

	list, err := trs.Read_sheets_CASHBOX_err("1_EWAl5qFMMWh0xbr3gzJ_KM87QNXhm7wJy89Vrs_zDs", "'Tracing'!A3:U", 0)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, i := range list {
		if i[0] == "" {
			continue
		}
		r := Referal{}
		r.Reestablish_Tracking(i)

		if r.Archive {
			r.Up_Archive(db)
			r.Upgrade_Tracking_List_Clear()
		}
		if r.Upd {
			r.Delete(db)
			r.Trips(db)
			r.RecDB(db)
			r.Upgrade_Tracking_List()
		}
		if r.Del {
			r.Delete(db)
			r.Upgrade_Tracking_List_Clear()
		}
	}
}

func (a *Referal) Trips(db *sql.DB) {
	trips := 0
	dat := List_Day(a.Date_Start.Format(TNS), a.Date_Finish.Format(TNS))
	cit := sity{}
	cit.Search(a.City)
	graf := get_grafik(cit, dat)
	trip := get_Trips([]sity{cit}, dat)
	if a.Acc_bolt == "" {
		a.Acc_bolt = "no accaunt"
	}
	if a.Acc_uber == "" {
		a.Acc_uber = "no accaunt"
	}
	if a.Acc_uklon == "" {
		a.Acc_uklon = "no accaunt"
	}
	for n, dt := range dat {
		in_car := 0
		tripsday := 0
		for _, gr := range graf {
			if dt.Format(TNS) != gr.data.Format(TNS) {
				continue
			}
			if gr.name_vod1 == a.Acc_bolt || gr.name_vod2 == a.Acc_bolt || gr.name_vod1 == a.Acc_uber || gr.name_vod2 == a.Acc_uber || gr.name_vod1 == a.Acc_uklon || gr.name_vod2 == a.Acc_uklon {
				in_car++
				for _, tr := range trip {
					if dt.Format(TNS) != tr.Date {
						continue
					}
					if n == 0 { // первый день работы
						if trs.Clear_name_driver(a.Acc_bolt) == tr.Name || trs.Clear_name_driver(a.Acc_uklon) == tr.Name || trs.Clear_name_driver(a.Acc_uber) == tr.Name {
							tripsday = tripsday + tr.Trip
						}
					} else if trs.Clear_name_driver(gr.name_vod1) == trs.Clear_name_driver(gr.name_vod2) { // если один аккаунт
						if trs.Clear_name_driver(gr.name_vod1) == tr.Name {
							tripsday = tripsday + tr.Trip
						}
					} else {
						if trs.Clear_name_driver(gr.name_vod1) == tr.Name || trs.Clear_name_driver(gr.name_vod2) == tr.Name { // ekipag
							tripsday = tripsday + tr.Trip
						}
					}
				}
			}
		}
		if in_car == 0 { // проверяем поездки если водителя нет в графике за сегодня
			for _, tr := range trip {
				if dt.Format(TNS) != tr.Date {
					continue
				}
				if trs.Clear_name_driver(a.Acc_bolt) == tr.Name || trs.Clear_name_driver(a.Acc_uklon) == tr.Name || trs.Clear_name_driver(a.Acc_uber) == tr.Name {
					tripsday = tripsday + tr.Trip
				}
			}
		}
		if in_car > 0 { // если за день работал более чем на 1 авто
			trips = trips + tripsday/in_car
		}
	}

	a.Trip = trips
	err := a.Up_Trip(db, fmt.Sprint(a.Trip))
	if err != nil {
		fmt.Println(err)
	}

	if a.Trip+1 > Norma_Trip_Referal {
		a.Status = "Умови виконано"
		if a.City == "Київ" {
			a.Cash = 2000
		} else {
			a.Cash = 1500
		}
		err := a.Up_Status(db, "Умови виконано")
		if err != nil {
			fmt.Println(err)
		}
		err = a.Up_Cash(db, fmt.Sprint(a.Cash))
		if err != nil {
			fmt.Println(err)
		}
	}

	if a.Date_Finish.Before(time.Now().AddDate(0, 0, -1)) && a.Trip < 500 {
		a.Status = "Умови не виконано"
		a.Up_Status(db, "Умови не виконано")
	}
}

// проверяет поданых на регистрацию
// пишет в базу прошедших проверку
func Referals_Register(db *sql.DB) {

	//fmt.Println("Referals_Register")

	list, err := trs.Read_sheets_CASHBOX_err("1_EWAl5qFMMWh0xbr3gzJ_KM87QNXhm7wJy89Vrs_zDs", "'Register'!A3:L", 0)
	if err != nil {
		fmt.Println(err)
		return
	}
	for n, i := range list {

		if i[0] == "" {
			continue
		}

		if len(i) < 11 {
			continue
		}
		//	fmt.Println(n, i)
		r := Referal{}
		stat, _ := r.Reestablish_Register(i)
		r.Print()

		pom, ok := r.Check_Register()
		if !ok {
			trs.Rec("1_EWAl5qFMMWh0xbr3gzJ_KM87QNXhm7wJy89Vrs_zDs", fmt.Sprintf("'Register'!L%d", n+3), [][]string{{pom}}, false, false)
			continue
		} else {
			trs.Rec("1_EWAl5qFMMWh0xbr3gzJ_KM87QNXhm7wJy89Vrs_zDs", fmt.Sprintf("'Register'!L%d", n+3), [][]string{{""}}, false, false)
		}

		if stat == "Погоджено" {
			r.Status = "У процесі"
			err := r.RecDB(db)
			if err == nil {
				trs.Rec("1_EWAl5qFMMWh0xbr3gzJ_KM87QNXhm7wJy89Vrs_zDs", fmt.Sprintf("'Register'!A%d", n+3), [][]string{Vidhileno_Str}, false, false)
			}
			r.Rec_Tracking_List()
		}
	}
}

type Referals []Referal

func (a *Referals) Read(db *sql.DB, s string) {
	*a = Read_DB_Referal(db, s)
}

func (a *Referals) Print() {
	for _, i := range *a {
		i.Print()
	}
}

// Статуси:
// У процесі
// Умови виконано
// Умови не виконано

type Referal struct {
	City            string    // Місто
	Beneficiary_Pib string    // Вигодонабувач реферала
	Beneficiary_tel string    // Вигодонабувач реферала
	Status          string    // Activ - Passiv
	Cash            int       // до виплати
	Referal_Pib     string    // ПІБ реферала
	Referal_tel     string    // ПІБ реферала
	Date_Start      time.Time // Старт програми реферала
	Date_Finish     time.Time // Фініш програми реферала
	Trip            int       // Поїздок за період
	TG_id           int64     // Telegram ID
	Pasport         string    // № паспорта
	Posvidch        string    // № водійського посвідчення
	Linc_Bitrix     string    // посилання на Bitrix
	ID_cod          string    // Ідентефікаційний код
	Acc_bolt        string    // Акаунт Болт
	Acc_uklon       string    // Акаунт Уклон
	Acc_uber        string    // Акаунт УБЕР
	Archive         bool      // Archive
	Upd             bool      // обновить
	Del             bool      // удалить
}

// подготовка строки к записи на вкладку
// Tracing
func (a Referal) SliseString() (res []string) {
	res = append(res, a.City)             // Місто
	res = append(res, fmt.Sprint(a.Cash)) // До виплати, грн
	res = append(res, a.Status)           // Результат програми реферала
	if a.Archive {
		res = append(res, "Так")
	} else {
		res = append(res, "Ні")
	}
	res = append(res, a.Date_Start.Format(TNS))  // Старт програми реферала
	res = append(res, fmt.Sprint(a.Trip))        // Поїздок за період
	res = append(res, a.Date_Finish.Format(TNS)) // Фініш програми реферала
	res = append(res, a.Beneficiary_Pib)         // Вигодонабувач реферала PIB
	res = append(res, a.Beneficiary_tel)         // Вигодонабувач реферала tel
	res = append(res, a.Referal_Pib)             // ПІБ реферала
	res = append(res, a.Referal_tel)             // tel реферала
	res = append(res, a.Pasport)                 // № паспорта
	res = append(res, a.Posvidch)                // № водійського посвідчення
	res = append(res, a.ID_cod)                  // Ідентефікаційний код
	res = append(res, a.Linc_Bitrix)             // Посилання на угоду у Bitrix
	if a.Acc_bolt == "" {
		res = append(res, "no accaunt")
	} else {
		res = append(res, a.Acc_bolt) // Акаунт Болт
	}
	if a.Acc_uklon == "" {
		res = append(res, "no accaunt")
	} else {
		res = append(res, a.Acc_uklon) // Акаунт Уклон
	}
	if a.Acc_uber == "" {
		res = append(res, "no accaunt")
	} else {
		res = append(res, a.Acc_uber) // Акаунт УБЕР
	}
	res = append(res, fmt.Sprint(a.TG_id)) // Telegram ID
	res = append(res, "Ні")                // обновить
	res = append(res, "Ні")                // удалить
	return
}

// подготовка строки к записи на вкладку
// Archive
func (a Referal) SliseString_Archive() (res []string) {
	res = append(res, a.City)                    // Місто
	res = append(res, fmt.Sprint(a.Cash))        // До виплати, грн
	res = append(res, a.Status)                  // Результат програми реферала
	res = append(res, a.Date_Start.Format(TNS))  // Старт програми реферала
	res = append(res, fmt.Sprint(a.Trip))        // Поїздок за період
	res = append(res, a.Date_Finish.Format(TNS)) // Фініш програми реферала
	res = append(res, a.Beneficiary_Pib)         // Вигодонабувач реферала PIB
	res = append(res, a.Beneficiary_tel)         // Вигодонабувач реферала tel
	res = append(res, a.Referal_Pib)             // ПІБ реферала
	res = append(res, a.Referal_tel)             // tel реферала
	res = append(res, a.Pasport)                 // № паспорта
	res = append(res, a.Posvidch)                // № водійського посвідчення
	res = append(res, a.ID_cod)                  // Ідентефікаційний код
	res = append(res, a.Linc_Bitrix)             // Посилання на угоду у Bitrix
	res = append(res, a.Acc_bolt)                // Акаунт Болт
	res = append(res, a.Acc_uklon)               // Акаунт Уклон
	res = append(res, a.Acc_uber)                // Акаунт УБЕР
	res = append(res, fmt.Sprint(a.TG_id))       // Telegram ID
	return
}

func (a *Referal) Check_Register() (string, bool) {
	if a.City == "" {
		return "Не заповнено поле \"Місто\"", false
	}
	if a.Beneficiary_Pib == "" {
		return "Не заповнено поле \"Вигодонабувач ПІБ (як в касі парку)\"", false
	}
	if a.Beneficiary_tel == "" {
		return "Не заповнено поле \"Вигодонабувач Телефон\"", false
	}
	if strings.Count(a.Referal_Pib, " ") < 1 {
		return "Не заповнено поле \"ПІБ реферала\"", false
	}
	if a.Pasport == "" {
		return "Не заповнено поле \"№ паспорта\"", false
	}
	if a.Posvidch == "" {
		return "Не заповнено поле \"№ водійського посвідчення\"", false
	}
	if a.ID_cod == "" {
		return "Не заповнено поле \"Ідентефікаційний код\"", false
	}
	if a.TG_id == 0 {
		return "Некоректний або відсутній TG_id", false
	}
	if a.Acc_bolt == "" && a.Acc_uber == "" && a.Acc_uklon == "" {
		return "Не знайдено акаунтів за наданим TG_id", false
	}

	return "", true
}

// Печать элемента Referal
func (a *Referal) Print() {
	fmt.Println()
	fmt.Println("Referal             :")
	fmt.Println("City                :", a.City)
	fmt.Println("Beneficiary_Pib     :", a.Beneficiary_Pib)
	fmt.Println("Beneficiary_tel     :", a.Beneficiary_tel)
	fmt.Println("Status              :", a.Status)
	fmt.Println("Cash                :", a.Cash)
	fmt.Println("Referal_Pib         :", a.Referal_Pib)
	fmt.Println("Referal_tel         :", a.Referal_tel)
	fmt.Println("Date_Start          :", a.Date_Start.Format(TNS))
	fmt.Println("Date_Finish         :", a.Date_Finish.Format(TNS))
	fmt.Println("Trip                :", a.Trip)
	fmt.Println("TG_id               :", a.TG_id)
	fmt.Println("Pasport             :", a.Pasport)
	fmt.Println("Posvidch            :", a.Posvidch)
	fmt.Println("Linc_Bitrix         :", a.Linc_Bitrix)
	fmt.Println("ID_cod              :", a.ID_cod)
	fmt.Println("Acc_bolt            :", a.Acc_bolt)
	fmt.Println("Acc_uklon           :", a.Acc_uklon)
	fmt.Println("Acc_uber            :", a.Acc_uber)
	fmt.Println("Archive             :", a.Archive)
	fmt.Println("Update              :", a.Upd)
	fmt.Println("Delete              :", a.Del)
}

// Просчет суммы значения
// по массиву []Referal
func (a *Referal) Add(b Referal) {
	a.Trip = a.Trip + b.Trip
	a.TG_id = a.TG_id + b.TG_id
}

// Восстановление элемента Referal
// из строки БД
func (a *Referal) Reestablish(p []string) {
	a.City = p[1]
	a.Beneficiary_Pib = p[2]
	a.Beneficiary_tel = p[3]
	a.Status = p[4]
	a.Cash = trs.String_to_int(p[5])
	a.Referal_Pib = p[6]
	a.Referal_tel = p[7]
	a.Date_Start, _ = time.Parse(TNS, p[8])
	a.Date_Finish, _ = time.Parse(TNS, p[9])
	a.Trip = trs.String_to_int(p[10])
	a.TG_id = int64(trs.String_to_int(p[11]))
	a.Pasport = p[12]
	a.Posvidch = p[13]
	a.ID_cod = p[14]
	a.Linc_Bitrix = p[15]
	a.Acc_bolt = p[16]
	a.Acc_uklon = p[17]
	a.Acc_uber = p[18]
	if p[19] == "yes" {
		a.Archive = true
	} else {
		a.Archive = false
	}
}

// Восстановление элемента Referal
// из строки листа  Register
func (a *Referal) Reestablish_Register(p []string) (Status string, ok bool) {
	a.City = p[0]
	a.Beneficiary_Pib = p[1]
	a.Beneficiary_tel = p[2]
	a.Referal_Pib = p[3]
	a.Referal_tel = p[4]
	a.Pasport = p[5]
	a.Posvidch = p[6]
	a.ID_cod = p[7]
	a.Linc_Bitrix = p[8]

	a.TG_id = int64(trs.String_to_int(strings.Trim(p[9], " ")))
	a.Archive = false
	Status = p[10]
	if Status == "Погоджено" || Status == "Відхилено" {
		ok = true
	}
	a.Date_Start = time.Now().AddDate(0, 0, 1)
	a.Date_Finish = time.Now().AddDate(0, 1, 1)
	if a.TG_id > 0 {
		r := Drivers{}
		r.Read_TgId(a.TG_id)
		for _, i := range r {
			if i.Taxi_type == "BOLT" {
				a.Acc_bolt = i.Accaunt
			} else if i.Taxi_type == "UKLON" {
				a.Acc_uklon = i.Accaunt
			} else if i.Taxi_type == "UBER" {
				a.Acc_uber = i.Accaunt
			}
		}
	}
	return
}

func (a *Referal) ArchiveSet() {
	a.Archive = true
}

// Восстановление элемента Referal
// из строки листа  Tracking
func (a *Referal) Reestablish_Tracking(p []string) {
	if len(p) < 19 {
		return
	}
	a.City = p[0]
	a.Cash = string_to_int(p[1])
	a.Status = p[2]
	a.Date_Start, _ = time.Parse(TNS, p[4])
	a.Trip = trs.String_to_int(p[5])
	a.Date_Finish, _ = time.Parse(TNS, p[6])
	a.Beneficiary_Pib = p[7]
	a.Beneficiary_tel = p[8]
	a.Referal_Pib = p[9]
	a.Referal_tel = p[10]
	a.Pasport = p[11]
	a.Posvidch = p[12]
	a.ID_cod = p[13]
	a.Linc_Bitrix = p[14]
	a.Acc_bolt = p[15]
	a.Acc_uklon = p[16]
	a.Acc_uber = p[17]
	a.TG_id = int64(trs.String_to_int(p[18]))
	if p[3] == "Так" {
		a.Archive = true
	} else {
		a.Archive = false
	}
	if p[19] == "Так" {
		a.Upd = true
	} else {
		a.Upd = false
	}
	if p[20] == "Так" {
		a.Del = true
	} else {
		a.Del = false
	}

}

// Обновление TG id элемента Referal в БД
func (a *Referal) Up_TgId(db *sql.DB, id int64) (err error) {
	zap := fmt.Sprintf("UPDATE referal SET TG_id  = '%d' WHERE City = '%s' AND  Referal_Pib = '%s'", id, a.City, a.Referal_Pib)
	_, err = db.Exec(zap)
	return
}

// Обновление дата начала элемента Referal в БД
func (a *Referal) Up_Date(db *sql.DB, dn, dk string) (err error) {
	zap := fmt.Sprintf("UPDATE referal SET Date_Start  = '%s', Date_Finish = '%s' WHERE City = '%s' AND  Referal_Pib = '%s'", dn, dk, a.City, a.Referal_Pib)
	_, err = db.Exec(zap)
	return
}

// Обновление аккаунта Болт элемента Referal в БД
func (a *Referal) Up_Bolt(db *sql.DB, s string) (err error) {
	zap := fmt.Sprintf("UPDATE referal SET Acc_bolt = '%s' WHERE City = '%s' AND  Referal_Pib = '%s'", s, a.City, a.Referal_Pib)
	_, err = db.Exec(zap)
	return
}

// Обновление аккаунта Уклон элемента Referal в БД
func (a *Referal) Up_Uklon(db *sql.DB, s string) (err error) {
	zap := fmt.Sprintf("UPDATE referal SET Acc_uklon = '%s' WHERE City = '%s' AND  Referal_Pib = '%s' ", s, a.City, a.Referal_Pib)
	_, err = db.Exec(zap)
	return
}

// Обновление аккаунта Убер элемента Referal в БД
func (a *Referal) Up_Uber(db *sql.DB, s string) (err error) {
	zap := fmt.Sprintf("UPDATE referal SET Acc_uber = '%s' WHERE City = '%s' AND  Referal_Pib = '%s'", s, a.City, a.Referal_Pib)
	_, err = db.Exec(zap)
	return
}

// Отправка в архив элемента Referal в БД
func (a *Referal) Up_Archive(db *sql.DB) (err error) {
	if a.Status == "У процесі" {
		a.Status = "Архів"
	}
	zap := fmt.Sprintf("UPDATE referal SET Archive = '%s' WHERE City = '%s' AND  Referal_Pib = '%s' AND  Referal_tel = '%s'", "yes", a.City, a.Referal_Pib, a.Referal_tel)
	_, err = db.Exec(zap)
	return
}

// Обновление Status элемента Referal в БД
func (a *Referal) Up_Status(db *sql.DB, s string) (err error) {
	zap := fmt.Sprintf("UPDATE referal SET Status = '%s' WHERE City = '%s' AND  Referal_Pib = '%s' AND  Referal_tel = '%s'", s, a.City, a.Referal_Pib, a.Referal_tel)
	_, err = db.Exec(zap)

	return
}

// Обновление Cash элемента Referal в БД
func (a *Referal) Up_Cash(db *sql.DB, s string) (err error) {
	zap := fmt.Sprintf("UPDATE referal SET Cash = '%s' WHERE City = '%s' AND  Referal_Pib = '%s' AND  Referal_tel = '%s'", s, a.City, a.Referal_Pib, a.Referal_tel)
	_, err = db.Exec(zap)
	return
}

// Обновление Cash элемента Referal в БД
func (a *Referal) Up_Trip(db *sql.DB, s string) (err error) {
	zap := fmt.Sprintf("UPDATE referal SET Trip = '%s' WHERE City = '%s' AND  Referal_Pib = '%s' AND  Referal_tel = '%s'", s, a.City, a.Referal_Pib, a.Referal_tel)
	_, err = db.Exec(zap)
	return
}

// Delete элемента Referal в БД
func (a *Referal) Delete(db *sql.DB) (err error) {
	zap := fmt.Sprintf("DELETE FROM referal WHERE City = '%s' AND  Referal_Pib = '%s'", a.City, a.Referal_Pib)
	_, err = db.Exec(zap)
	return
}

// Запись элемента Referal в БД
func (a *Referal) RecDB(db *sql.DB) (err error) {

	name := "City, Beneficiary_Pib, Beneficiary_tel, Status, Cash, Referal_Pib, Referal_tel, Date_Start, Date_Finish, Trip, TG_id, Pasport, Posvidch, ID_cod, Linc_Bitrix, Acc_bolt, Acc_uklon, Acc_uber, Archive"
	arc := "no"
	if a.Archive {
		arc = "yes"
	}
	recdb := fmt.Sprintf("INSERT INTO Referal (%s) VALUES ('%s', '%s', '%s', '%s', '%d', '%s', '%s', '%s', '%s', '%d', '%d', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s')", name, a.City, a.Beneficiary_Pib, a.Beneficiary_tel, a.Status, a.Cash, a.Referal_Pib, a.Referal_tel, a.Date_Start.Format(TNS), a.Date_Finish.Format(TNS), a.Trip, a.TG_id, a.Pasport, a.Posvidch, a.ID_cod, a.Linc_Bitrix, a.Acc_bolt, a.Acc_uklon, a.Acc_uber, arc)

	_, err = db.Exec(recdb)
	if err != nil {
		fmt.Println(err, recdb)
	}

	return err
}

func (a *Referal) Read_Db(db *sql.DB, city, pip, tel string) (err error) {
	zap := fmt.Sprintf("SELECT * FROM referal WHERE City = '%s' AND  Referal_Pib = '%s' AND  Referal_tel = '%s'", city, pip, tel)

	row := db.QueryRow(zap)

	p := make([]string, 20)
	err = row.Scan(&p[0], &p[1], &p[2], &p[3], &p[4], &p[5], &p[6], &p[7], &p[8], &p[9], &p[10], &p[11], &p[12], &p[13], &p[14], &p[15], &p[16], &p[17], &p[18], &p[19])

	if err != nil {
		fmt.Println(err)
		return
	}
	a.Reestablish(p)
	return
}

// Чтение элемента Referal из БД
func Read_DB_Referal(db *sql.DB, s string) (res []Referal) {

	zap := ""

	if s == "Activ" {
		zap = "SELECT * FROM referal WHERE Archive = 'no'"
	} else if s == "Archive" {
		zap = "SELECT * FROM referal WHERE Archive = 'yes'"
	}

	//fmt.Println(zap)

	rows, err := db.Query(zap)
	if err != nil {
		fmt.Println(err, zap)
		return
	}

	for rows.Next() {
		p := make([]string, 20)
		err = rows.Scan(&p[0], &p[1], &p[2], &p[3], &p[4], &p[5], &p[6], &p[7], &p[8], &p[9], &p[10], &p[11], &p[12], &p[13], &p[14], &p[15], &p[16], &p[17], &p[18], &p[19])
		if err != nil {
			fmt.Println(err)
		}
		r := Referal{}
		r.Reestablish(p)
		res = append(res, r)
	}
	return
}

// Чтение элемента Referal из БД
func Read_DB_Referal_Pib(db *sql.DB, s string) (res []Referal) {

	zap := "SELECT * FROM referal WHERE Referal_Pib = '" + s + "'"

	//fmt.Println(zap)

	rows, err := db.Query(zap)
	if err != nil {
		fmt.Println(err, zap)
		return
	}

	for rows.Next() {
		p := make([]string, 20)
		err = rows.Scan(&p[0], &p[1], &p[2], &p[3], &p[4], &p[5], &p[6], &p[7], &p[8], &p[9], &p[10], &p[11], &p[12], &p[13], &p[14], &p[15], &p[16], &p[17], &p[18], &p[19])
		if err != nil {
			fmt.Println(err)
		}
		r := Referal{}
		r.Reestablish(p)
		res = append(res, r)
	}
	return
}

func Repair_Referal_Accaunt() {

	db, err := sql.Open("sqlite3", path+"sqlite/referal.db")
	if err != nil {
		fmt.Println(path+"sqlite/referal.db", err)
	}
	defer db.Close()

	rr := Referals{}
	rr.Read(db, "Activ")

	for _, i := range rr {

		ii := i
		r := Drivers{}
		r.Read_TgId(i.TG_id)
		for _, j := range r {
			if j.Taxi_type == "BOLT" {
				i.Acc_bolt = j.Accaunt
			} else if j.Taxi_type == "UKLON" {
				i.Acc_uklon = j.Accaunt
			} else if j.Taxi_type == "UBER" {
				i.Acc_uber = j.Accaunt
			}
		}
		if i.Acc_bolt != ii.Acc_bolt || i.Acc_uber != ii.Acc_uber || i.Acc_uklon != ii.Acc_uklon {
			i.Print()
			ii.Print()
			i.Delete(db)
			i.RecDB(db)
		}
	}
}
