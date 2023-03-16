package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"trs"
	"unsafe"
)

func Print_time_slise(a *[]time.Time) {
	for _, i := range *a {
		fmt.Println(i.Format(TNSF))
	}
}

var test_driver = map[string]bool{
	"Суми|Дмитро Коваленко":      true,
	"Суми|(UK) Дмитро Коваленко": true,
	"Суми|Руслан Панахов":        true,
	"Суми|UK) Руслан Панахов":    true,
	"Sumy|Дмитро Коваленко":      true,
	"Sumy|(UK) Дмитро Коваленко": true,
	"Sumy|Руслан Панахов":        true,
	"Sumy|UK) Руслан Панахов":    true,
}

var list_name = map[string]string{
	"1vXvaa-FVmDkJxhVvbdkuZ5BZDGRW9EoIm3skYsgUzGE": "List SITY",
	"1XLH9bF5TeuSaMr1PO-PVtv9_8H5i03fA8Y7BvyYePzo": "График за сегодня для чтения",
	"172ZsPP13-yVv7EVUZrymXpX1JjuSJG2RmD34p7PsM9I": "График за сегодня",
	"1Y-E8coo3Td3BwysXQzIECO0FrXcsOHikhN26OpMk8tg": "ДМ. Робота з водіями.",
	"1IJgBM5qObVF-FRfc64nIZX8-w6GSCFVzRlokIXN2r1E": "История ДМ работа с водителями",
	"1Gv6kmVA0wmOxaISB59laiUYWLauddvQmBP9Zm4pJvQw": "КПД 90 дней",
	"18hK0457uyVcfDQu9N5fhKzreHIfAiFkFRPVevVxiw8o": "КПД АВТОПАРКОВ",
	"18fpt9H8pR2P3_NY_IAeTpGZGXPav5mgqC-6m6fJg5nc": "КПД Отчет Autoreport V3",
	"10jUdHs6lYQwb6vEpWhfUY_f9BmMnTFcLuNlXPrB9oxA": "КПД Отчет Autoreport V4",
	"1uvZoC0bOoY4rPJVO62J0MX3uTh6rrC_cmwbYhePyZqE": "КПИ ДМ",
	"1uUlPjoPjjiSyC-fvx8_rF73BB9IaFwI9lsMroS718lw": "Отчет НОВЫЙ ВОДИТЕЛЬ",
	"1Dws1iUuLbS1RoEGxIqXduJqWId27eif3U6kCLc9rqOs": "Авто на 2 службах",
	"1r7-bcbpZmzWWTeHLm4WGJeYRyZxoDTSQVyFpDu3OYiM": "Справочник водителей на сегодня",
	"1GP5CJ1APDlC4-y_zjMhpVhmwl2Nco6SIC2biEUROdO0": "Средний ПРОБЕГ НОВЫЙ",
	"18DSsXmfP1OizUcbDL8Re1qOKnKt36viAXNSk3TPEWP4": "Тижневий звіт",
	"1siTwp9RYsuVcLyppCPcJjjkKpvawqUIxuaFZoPADYYQ": "Фин отчет по городам",
	"1ymg8HjgV5hIkNYWMl2AsFKOdGWfYXwOtcg6MYaY3Syg": "Speed Control",
	"1KaPDbkJr7rJi7BYQJ_6nzWyOyRnpWVs9PokD40X69q4": "Касса города Мукачево",
	"1kW4Oi_mdQWwpSnbPP9wJhgfJa5S0AEVHs6P7oVjErZc": "Касса города Івано Франківськ",
	"1usL5ixIR4fxXm3gL8NVeCUi8zEzJLpiYoaJRDY0NLe8": "Касса города Суми",
	"1btFuFzjWQyqnL-3KxctcYOJLpg-2SMcFdSj6RrId7po": "Касса города Луцьк",
	"1qntN3DIlp3uBGARt-MVlB1CIPtRs_wP59QdwJRM88UQ": "Касса города Тернопіль",
	"1r4AwzCU4fSOxOnlyDS8UtBJIhz5sadaM9-Tz72_olU4": "Касса города Київ",
	"1Hji5OX-MXpblrDx760ZmvOA8bWVvtSZDzbgILL8ZRbY": "Касса города Запоріжжя",
	"1jpYdBDTVPYonQKZL67nTEH3viAYefmZ-tNxNy4hUT-U": "Касса города Дніпро",
	"1DuuqCPGk4P3uRE7lOjZmBg3l2N9j0JmpuigRX-aGhcY": "Касса города Одеса",
	"1kSl8ujUqcqEljVtiwKPPC9gupyj44xk9nTBCGAStKUc": "Касса города Ужгород",
	"1uSSrsQTheP4Atu5m6PEZJa2lErCEDNdGu498_CEQPj0": "Касса города Черкаси",
	"1BpMiT8AbMMusJVwlhlQ428OVR7a54yhqDp7CCDkeUkQ": "Касса города Чернігів",
	"1aLCsYT9eUwUNuagJA4xHIX8ht41B7M9OUIoIBQK9AxE": "Касса города Чернівці",
	"1FrmSwFuaAebzsP-DaY50HfsQzL1-1jk7jrNOkSYKOhs": "Касса города Полтава",
	"1EfAVZzfSe0Du9sej1xpayL0LyGHZ3494j9eM7SNGPB8": "Касса города Львів",
	"1oS2yQqjM_NYsNwBb34_RXKzYFPIL8TQxbBlRB5dgKfs": "Касса города Кривий Ріг",
	"1XFsWEz0vBFqe3LzhhtMnJOEYt_sme5trT-XfKj-j7-g": "Касса города Рівне",
	"1kqbRGXtJxmH7d2JCDYIlZJiI83ManHANQdkxtKQSOmw": "Касса города Вінниця",
	"1BWnaHdrd0tBS7JXCDapPfByO7v_qFYT-RJCAdglzd-E": "Касса города Житомир",
	"1dKWGIhBhVDZ_XIUJ3J0n0ul2QRbPSGHjhGo2x8LPHh8": "Касса города Хмельницький",
	"18EJ2d5T9-QzAYAwrk_MzVVbuvmEh7_IFDElz4t43D9Y": "Касса города Харків",
	"1IQKfnTbgJLKM8n8JFdOkpOpreNe178fsRB1RUDwJlKw": "Касса города Мукачево",
	"1JET8OdmDfZSEpaAO7nJjAV3WBgEsXDqDpFgl-_Z5x3A": "Касса города Івано Франківськ",
	"1oIo3tdzR5RhqwOuXU9GAnsjw9AV2RjgUEQCoSDIuOuY": "Касса города Суми",
	"1uP_n6mOqn6-Gi_az_YmopDMS9Ranz5OCEHfgx-6Gas8": "Касса города Луцьк",
	"1l_RfbBJp7iPXglj6bkRmO9oXk45v372ivaRQCUV41zY": "Касса города Тернопіль",
	"1qCs8qjK-K_N4_uSngCwwcTtNo2RNZXn78hn0kbA5Jrg": "Касса города Київ",
	"1doiMprUYKdZEPLpEcTtvtkjszcZLKZGaseNnVlSHUps": "Касса города Запоріжжя",
	"1w1QljcTbKor040d7rRLVjfb8NGfnBXVsA88O9IYj5pU": "Касса города Дніпро",
	"1O-DcAs5GreUENH2gjPss5g0PuQy4TT33pQgGsnl93dY": "Касса города Одеса",
	"13NYz9CyWLHlJvo2AXUOJ1i0QZdeaQZJrpo2SKRkaswQ": "Касса города Ужгород",
	"13VtpsPQjmUNA5yDp4ONsdjnGhicEFed_9vD7KOzLKzM": "Касса города Черкаси",
	"109XwhdNzCxVubub1d8cm22LG_L6VVt8Oh5bIB95jIYo": "Касса города Чернігів",
	"1wezVTQN9Ufn01xBn8usiwO9Ty1ISt47hmir5WFrDdRs": "Касса города Чернівці",
	"1V6Vuy2QB0qikMplPNaeCKA_g8QG6jAgZjCyFDCDESj0": "Касса города Полтава",
	"1PcJj89wBtUkdc3rEIFwKgnB7BakvRvSKZK6p1057-64": "Касса города Львів",
	"1o3-CVKWiu0TDR91-2d8vNdRxvwWYRXPkp9F6w03ipUo": "Касса города Кривий Ріг",
	"16UC1w8ZEpZYJsfeSOXCjkJQspinGbmViNeqaFWL5-ow": "Касса города Рівне",
	"1AV168uV1cb0S6nbElUHSP9nUxRWI0i_SNxcq7QJpksI": "Касса города Вінниця",
	"11ODi4dzBcqYJbjoSEk0tbheCJtzf7k8IitXIUN-D8xc": "Касса города Житомир",
	"1PVwCRm9tg3mOs6nSer7nzssfdc0fSByxlWUmezxz0a8": "Касса города Хмельницький",
	"1wU5aMG_vSSc7A9aDHzlivcGrvITzydKPrwi473rTGAc": "Касса города Харків",
	"18usVSob1OcMH6UKKVa2JxuzqDe1DFL10QxYLOU3SmnE": "Касса города Херсон",
	"1Ooh6CRu0RFH6kjU4gA3rRGJSx0QsTB6KTOIwf7YGTzc": "Касса города Варшава",
	//"1AV168uV1cb0S6nbElUHSP9nUxRWI0i_SNxcq7QJpksI": "График города Винница",
	//"": "",
	//"": "",
	//"": "",
	//"": "",
	//"": "",
	//"": "",
	//"": "",
	//"": "",
	//"": "",
	//"": "",
	//"": "",
	//"": "",
	//"": "",
	//"": "",
	//"": "",
	//"": "",
	//"": "",
	//"": "",
	//"": "",
	//"": "",
	//"": "",
	//"": "",
	//"": "",
	//"": "",
	//"": "",
	//"": "",
}

var status = map[string]string{
	"ДТП":                   "ДТП",
	"NW-ДТП":                "ДТП",
	"Ремонт":                "Ремонт",
	"NW-Ремонт":             "Ремонт",
	"ТО-Сервіс":             "Сервіс",
	"NW-ТО-Сервіс":          "Сервіс",
	"Выходной":              "Вихідний",
	"Вихідний":              "Вихідний",
	"ВЫХОДНОЙ":              "Вихідний",
	"Відновлення-авто":      "Відновлення-авто",
	"ЕВАКУАЦІЯ - Відправка": "ЕВАКУАЦІЯ - Відправка",
	"ЕВАКУАЦІЯ-Відправка":   "ЕВАКУАЦІЯ - Відправка",
	"ЕВАКУАЦІЯ - Прийом":    "ЕВАКУАЦІЯ - Прийом",
	"ЕВАКУАЦІЯ-Прийом":      "ЕВАКУАЦІЯ - Прийом",
	"ВИКОНАНО":              "ВИКОНАНО",
	"Штрафмайданчик":        "Штрафмайданчик",
	"Штрафмайданчик ":       "Штрафмайданчик",
	"NW-Штрафплощадка":      "Штрафмайданчик",
	"Документи":             "Документи",
	"NW-Документы":          "Документи",
	"Оренда":                "Оренда",
	"NW-Аренда":             "Оренда",
	"ПРОДАЖА":               "Оренда",
	"Викуп":                 "Оренда",
	"Угон":                  "Угон",
	"NW-Угон":               "Угон",
	"Ключі":                 "Ключі",
	"NW-Ключи":              "Ключі",
	"Волонтер":              "Волонтер",
	"NW-Волонтер":           "Волонтер",
	"Списання":              "Списання",
	"СПИСАННЯ":              "Списання",
	"NW-Списання":           "Списання",
	"Лікарняний":            "Лікарняний",
	"Больничный":            "Лікарняний",
	"NW-Лікарняний":         "Лікарняний",
	"Перегін":               "Перегін",
}

var NAME_SITY_SPEED = map[string]string{
	"Kyiv":                "Kyiv",
	"Brovary":             "Kyiv",
	"Uzhhorod":            "Uzhhorod",
	"Lutsk":               "Lutsk",
	"Rivne":               "Rivne",
	"Vinnytsia":           "Vinnytsia",
	"Вінниця":             "Vinnytsia",
	"Lviv":                "Lviv",
	"Ivano-Frankivsk":     "Ivano-Frankivsk", // Ivano
	"Ivano":               "Ivano-Frankivsk", // Ivano
	"Cherkasy":            "Cherkasy",
	"Chernivtsi":          "Chernivtsi",
	"Khmelnytskyi":        "Khmelnytskyi",
	"Dnipropetrovsk":      "Dnipropetrovsk",
	"Dnipro":              "Dnipropetrovsk",
	"Дніпро́":             "Dnipropetrovsk",
	"Kryvyi Rih":          "Kryvyi Rih", // KrRig
	"KrRig":               "Kryvyi Rih", // KrRig
	"Kryvyi":              "Kryvyi Rih", // KrRig
	"Odesa":               "Odesa",
	"Одеса":               "Odesa",
	"Zaporizhzhia":        "Zaporizhzhia", // Zaporizhia
	"Zaporizhia":          "Zaporizhzhia", // Zaporizhia
	"Zhytomyr":            "Zhytomyr",
	"Poltava":             "Poltava",
	"Sumy":                "Sumy",
	"Chernihiv":           "Chernihiv",
	"Kharkiv":             "Kharkiv",
	"Kherson":             "Kherson",
	"Kam":                 "Kamianets_Podilskyi",
	"Kamianets_Podilskyi": "Kamianets_Podilskyi",
	"ternopil":            "Ternopil",
	"Ternopil":            "Ternopil",
	"Mukachevo":           "Mukachevo",
}

type control_execute struct {
	data   string
	name   string
	coment string
}

// control := control_execute{}
// 	control.data = time.Now().Format("02.01.2006 15-04")
// 	control.name = ""
// 	control.coment = "completed"

// 	if control.read() {
// 		fmt.Println("Uclon Trip Read completed successfully")
// 		if !ignor_control {
// 			return
// 		}
// 	}

// defer control.save()

func (ctr control_execute) save() {
	f, err := os.OpenFile(path+"CONTROL/contexec.blt", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		errorLog.Println("No open file ", path+"CONTROL/contexec.blt")
	}
	defer f.Close()

	if _, err = f.WriteString(ctr.data + ":" + ctr.name + ":" + ctr.coment + "\n"); err != nil {
		fmt.Println("No write file ", path+"CONTROL/contexec.blt")
	}
	//OS(NS("control_execute ... "+ctr.name, 40))
}

func (ctr control_execute) read() bool {

	f, err := os.Open(path + "CONTROL/contexec.blt")
	if err != nil {
		errorLog.Println(err)
		return false
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		s := scanner.Text()
		sl := strings.Split(s, ":")
		if len(sl) != 3 {
			continue
		}
		if len(sl[0]) < 10 {
			continue
		}
		if sl[0][:10] != ctr.data[:10] {
			continue
		}
		if sl[1] != ctr.name {
			continue
		}
		if sl[2] != ctr.coment {
			continue
		}
		return true
	}
	return false
}

type behavior struct {
	city               string
	data               string
	car                string
	name               string
	tg_id              string
	score              float64
	eco_speed_class    string
	eco_speed_score    int
	braking_class      string
	braking_score      int
	acceltration_class string
	acceltration_score int
	turn_class         string
	turn_score         int
	comment            string
}

// создание строки для отчета
// ДМ робота з водыями
func (a *behavior) SliseString() []string {
	r := make([]string, 20)
	r[0] = a.city
	r[1] = a.data
	r[7] = float_to_string(a.score)
	r[8] = a.car
	r[9] = "Водій"
	r[10] = a.name
	r[11] = a.eco_speed_class
	r[12] = fmt.Sprint(a.eco_speed_score)
	r[13] = a.braking_class
	r[14] = fmt.Sprint(a.braking_score)
	r[15] = a.acceltration_class
	r[16] = fmt.Sprint(a.acceltration_score)
	r[17] = a.turn_class
	r[18] = fmt.Sprint(a.turn_score)
	r[19] = a.comment
	return r
}

// создание строки для отчета
// ДМ робота з водыями
func (a *behavior) SliseStringBOT() []string {
	r := make([]string, 4)
	r[0] = a.city
	r[1] = a.data
	r[2] = a.tg_id
	r[3] = a.comment
	return r
}

// коммент
func (a *behavior) Comment() {
	s := ""
	s = s + "Поведінка водіння водія " + a.name + " за " + a.data + " авто " + a.car + " місто " + a.city + ": "
	if a.score < 50 {
		s = s + "загальна оцінка - " + float_to_string(a.score) + ": "
	}
	if a.eco_speed_class == "F" || a.eco_speed_class == "G" {
		s = s + " еко швидкість  - оцінка " + a.eco_speed_class + " число випадків " + fmt.Sprint(a.eco_speed_score) + ": "
	}
	if a.braking_class == "F" || a.braking_class == "G" {
		s = s + " різке гальмуванн - оцінка " + a.braking_class + " число випадків " + fmt.Sprint(a.braking_score) + ": "
	}
	if a.acceltration_class == "F" || a.acceltration_class == "G" {
		s = s + " різке прискорення - оцінка " + a.acceltration_class + " число випадків " + fmt.Sprint(a.acceltration_score) + ": "
	}
	if a.turn_class == "F" || a.turn_class == "G" {
		s = s + " різкі повороти - оцінка " + a.turn_class + " число випадків " + fmt.Sprint(a.turn_score) + ": "
	}

	s = s + " Наразі Вам необхідно змінити стиль водіння на більш спокійний."

	a.comment = s
}

// Определение проблемности элемента behavior
func (a *behavior) def_problem() bool {
	if a.score < 50 {
		return true
	}
	if a.eco_speed_class == "F" || a.eco_speed_class == "G" {
		return true
	}
	if a.braking_class == "F" || a.braking_class == "G" {
		return true
	}
	if a.acceltration_class == "F" || a.acceltration_class == "G" {
		return true
	}
	if a.turn_class == "F" || a.turn_class == "G" {
		return true
	}
	return false
}

// Создание элемента behavior
func (a *behavior) Create(b []string, s, n, tg string) {
	if len(b) < 11 {
		return
	}
	a.city = s
	a.data = b[0]
	a.car = b[1]
	a.name = n
	a.tg_id = tg
	score64, err := strconv.ParseFloat(b[2], 64)
	if err != nil {
		a.score = 0
	} else {
		a.score = score64
	}
	a.eco_speed_class = b[3]
	score, err := strconv.Atoi(b[4])
	if err != nil {
		a.eco_speed_score = 0
	} else {
		a.eco_speed_score = score
	}
	a.braking_class = b[5]
	score, err = strconv.Atoi(b[6])
	if err != nil {
		a.braking_score = 0
	} else {
		a.braking_score = score
	}
	a.acceltration_class = b[7]
	score, err = strconv.Atoi(b[8])
	if err != nil {
		a.acceltration_score = 0
	} else {
		a.acceltration_score = score
	}
	a.turn_class = b[9]
	score, err = strconv.Atoi(b[10])
	if err != nil {
		a.turn_score = 0
	} else {
		a.turn_score = score
	}
}

type uclon_sum struct {
	city   string
	name   string
	data   string
	car    string
	trip   int64
	travel float64
	cash   float64
}

type uclon_trip struct {
	name string
	data string
	time string
}

type b_q_get struct {
	name_dataset string
	name_table   string
}

type speed struct {
	nomer  string
	data   string
	time   string
	park   string
	adress string
	speed  int
	id     int
}

// данные из справочника водителей
type vod_direct struct {
	name    string
	data    string
	time    string
	tel     string
	tg_id   string
	sity    string
	car     string
	id      string
	mail    string
	dogovor string
	status  string
	nach    string
	con     string
	prych   string
	lastned string
	aktiv   bool
	ofis    bool
	bolt    bool
	uklon   bool
	top     bool
	podmena bool
}

type report_taras_22022022 struct {
	park       string
	day        int
	vod        int
	cash       int
	clear_cash int
	dayvod     int
	tripvod    int
	cashvod    int
	day6       int
	day12      int
	day21      int
	day50      int
	day100     int
	day999     int
	// техническая составляющая
	ttw   int
	ttt   int
	ttc   int
	ttvod int
}

type peregon struct {
	park      []string
	CARNUMBER string
	MODEL     string
	MAPONID   string
	KOD       string
	misto     string
	mistoid   string
	misto_out string
	misto_in  string
	tel       string
	fio       string
	priniato  string
	vydaly    bool
}

type peregon_gr struct {
	CARNUMBER string
	MAPONID   string
	KOD       string
	misto     string
	mistoid   string
	sending   bool
	received  bool
	done      bool
	uroboti   bool
}

type list_car struct {
	nomer      string
	pereprobeg int
}

type horly_rep struct {
	hr    [24]float64
	data  time.Time
	nomer string
}

// отчет за день
type day struct {
	res_work []vod_day // результаты работы водителей за день
	data     string
	travel   float32
	hour     float32
	trip     int
	pay      int
}

// финансы водителя
type vod_finanse struct {
	name     string
	stag     int
	balans   float64
	dtp_borg float64
	deposit  float64
}

// отчет по водителю за день
type vod_day struct {
	impl        [24][3]int  // занятость по часам
	travel      [24]float64 // пробег по часам за день
	trip        [24]int     // поездки по часам за день
	cash        [24]int     // денег в час
	sity        string
	name        string
	name_clear  string
	data        string
	avto        string
	status      string // статус водителя (активный, выходной и т.п.)
	hour        float64
	travel_summ float64
	cashtrip    float64
	traveltrip  float64
	cashtravel  float64
	pereprobeg  float64

	procent_bolt    int
	procent_vozvrat int
	komissia        int
	proccent_vod    int
	cena_gas        int
	rasxod_avto     int

	cash_vod     float64
	cash_park    float64
	fuel_costs   float64
	zarplata_vod float64
	id           int
	pay          int
	trip_sum     int
	cash_summ    int
}

func (a vod_day) Print() {
	fmt.Println()
	fmt.Printf("Имя %s, %s\n", a.name, a.data)
	fmt.Printf("Поездок %d, пробег %.1f, кэш %d, авто %s\n", a.trip_sum, a.travel_summ, a.cash_summ, a.avto)
	fmt.Printf("Имя %s\n", a.name)
}

// просчет перепробега за день
func (a vod_day) rass_pereprobeg() int {
	pereprobeg := 0
	hour_no_trip := 0
	hour_with_trip := 0
	hour_expensive_trip := 0
	coeff_expensive_trip := 0
	for j := 0; j < 24; j++ {
		if a.travel[j] > 1 && a.trip[j] == 0 {
			pereprobeg = pereprobeg + int(a.travel[j])
		}
		if a.trip[j] == 0 {
			hour_no_trip++
		} else {
			hour_with_trip++
		}

		if a.trip[j] == 0 && a.cash[j] > 80 {
			hour_expensive_trip++
			coeff_expensive_trip = coeff_expensive_trip + a.cash[j]/90
		} else if a.cash[j] > a.trip[j]*80 {
			if a.cash[j]/a.trip[j] > 90 {
				hour_expensive_trip++
				coeff_expensive_trip = coeff_expensive_trip + (a.cash[j]-a.trip[j]*65)/90
			}
		}
	}

	if a.trip_sum > 0 {
		// скидка по дорогим поездкам
		sk1 := float64(coeff_expensive_trip) * a.travel_summ / float64(a.trip_sum)
		// скидка по 30км личного пробега
		sk2 := 30 * (float64(a.trip_sum) / 20)
		if sk2 > 30 {
			sk2 = 30
		}
		pereprobeg = pereprobeg - int(sk1)
		pereprobeg = pereprobeg - int(sk2)
	}

	if pereprobeg < 0 {
		return 0
	} else {
		return pereprobeg
	}
}

// выполняет слияние данных
// за один день
func merge_too_day(a, b day) (res day) {
	res.data = a.data
	res.travel = a.travel
	res.trip = a.trip + b.trip
	res.hour = a.hour + b.hour
	res.pay = a.pay + b.pay
	if len(a.res_work) == 1 && len(b.res_work) == 1 {
		res.res_work = append(res.res_work, summ_Vod_day(a.res_work[0], b.res_work[0]))
	} else if len(a.res_work) == 1 {
		res.res_work = append(res.res_work, a.res_work[0])
	} else if len(b.res_work) == 1 {
		res.res_work = append(res.res_work, b.res_work[0])
	}

	return
}

// суммирование 2х vod_day
func summ_Vod_day(a, b vod_day) (rs vod_day) {
	rs = a
	rs.name = a.name_clear
	rs.name_clear = a.name_clear
	rs.travel = a.travel
	rs.trip = summ_24int(a.trip, b.trip)
	rs.cash = summ_24int(a.cash, b.cash)
	rs.data = a.data
	rs.avto = a.avto
	rs.status = a.status
	rs.hour = a.hour + b.hour
	rs.cash_summ = a.cash_summ + b.cash_summ
	rs.trip_sum = a.trip_sum + b.trip_sum
	return
}

// суммирование   [24]int
func summ_24int(a, b [24]int) (rs [24]int) {
	for i := 0; i < 24; i++ {
		rs[i] = a[i] + b[i]
	}
	return
}

// суммирование 2х vod_day
func (a *vod_day) Summ(b vod_day) {
	a.trip_sum = a.trip_sum + b.trip_sum
	a.travel_summ = a.travel_summ + b.travel_summ
	a.hour = a.hour + b.hour
	a.cash_summ = a.cash_summ + b.cash_summ
	a.cash_park = a.cash_park + b.cash_park
	a.cash_vod = a.cash_vod + b.cash_vod
	a.fuel_costs = a.fuel_costs + b.fuel_costs
	a.zarplata_vod = a.zarplata_vod + b.zarplata_vod
	a.komissia = a.komissia + b.komissia
}

// просчет финансовых показателей по водителю
func (vod_day) Calculation(a vod_day) vod_day {
	a.cash_vod = (float64(a.cash_summ) - float64(a.cash_summ)*(float64(a.procent_bolt-a.procent_vozvrat)/100)) * (float64(a.proccent_vod) / 100)
	a.cash_park = float64(a.cash_summ) - a.cash_vod - float64(a.cash_summ)*float64(a.procent_bolt)/100
	a.fuel_costs = a.travel_summ * float64(a.rasxod_avto) * float64(a.cena_gas) / 100
	a.zarplata_vod = a.cash_vod - a.fuel_costs
	a.komissia = int(float64(a.cash_summ) * (float64(a.procent_bolt-a.procent_vozvrat) / 100))
	return a
}

// печать финансовых показателей
func (a vod_day) PrintCash() {
	fmt.Println()
	fmt.Println("Имя", a.name)
	fmt.Println("Поезок", a.trip_sum)
	fmt.Println("Часо онлайн", a.hour)
	fmt.Println("Пробег", a.travel_summ)
	fmt.Println("Касса всего", a.cash_summ)
	fmt.Println("Зарабтал водитель грязными", a.cash_vod)
	fmt.Println("Заработал водитель чистыми", a.zarplata_vod)
	fmt.Println("Потатил водитель на газ", a.fuel_costs)
	fmt.Println("Осталось парку", a.cash_park)
	//fmt.Println("", a.)
}

// готовит строку для печати
func (a vod_day) StrtoPrint() (res []string) {
	res = append(res, a.name, fmt.Sprint(a.trip_sum), fmt.Sprint(int(a.travel_summ)), fmt.Sprint(int(a.cash_summ)), fmt.Sprint(int(a.procent_bolt-a.procent_vozvrat)), fmt.Sprint(int(a.cash_summ*(a.procent_bolt-a.procent_vozvrat)/100)), fmt.Sprint(int(a.cash_park)), fmt.Sprint(int(a.cash_vod)), fmt.Sprint(int(a.fuel_costs)), fmt.Sprint(int(a.zarplata_vod)))
	return
}

// готовит строку для печати
func (a vod_day) StrtoTotakPrint() (res []string) {
	res = append(res, a.name, fmt.Sprint(a.trip_sum), fmt.Sprint(int(a.travel_summ)), fmt.Sprint(int(a.cash_summ)), fmt.Sprint(int(a.procent_bolt-a.procent_vozvrat)), fmt.Sprint(int(a.cash_summ-int(a.cash_park)-int(a.cash_vod))), fmt.Sprint(int(a.cash_park)), fmt.Sprint(int(a.cash_vod)), fmt.Sprint(int(a.fuel_costs)), fmt.Sprint(int(a.zarplata_vod)))
	return
}

type car_direct struct {
	TYPE         string
	MODEL        string
	CARNUMBER    string
	VINCODE      string
	DOCID        string
	KOD          string
	COMPANY      string
	SITY         string
	comment      string
	MAPONID      int64
	Nomer_grafik string
	Id_grafik    string
	Sity_grafik  string
	Nomer_Mapon  string
	Id_mapon     string
	Sity_mapon   string
}

// сокращенный вариант
type vod_day_trip struct {
	name string
	data string
	trip int
}

// элемент для bd
type element struct {
	sity string
	data string
	time string
	name string
	pay  int
	trip int
}

type sity struct {
	imployment      b_q_get
	bolt_trip       b_q_get
	uclon_trip      b_q_get
	uclon_cash      b_q_get
	uclon_today     b_q_get
	grafik          b_q_get
	vod             b_q_get
	kassa           b_q_get
	save_kassa_BGQ  bool
	ident           string
	id_BGQ          string
	name            string
	name_vod        string
	name_online     string
	name_grafik     string
	sity_report     string
	online_diapazon string
	peregon         string
	peregon_list    string
	nov_vod         string
	mapon_bolt      string
	mapon_uclon     string
	mapon_bolt_id   string
	mapon_uclon_id  string
	srez            string
	DM              string
	Kvoliti         string
	norma_trip_hour float64
	top_vodila      int
	srpob           int
	preprobeg       int
	workhour        int
	id              int
	procent_bolt    int
	procent_uclon   int
	procent_vozvrat int
	procent_vod     int
	cena_gas        int
	rasxod_avto     int
	cmp_grf         bool
	ext_rep         bool
	in_report       bool
}

func (a *sity) Print() {
	fmt.Println()
	fmt.Println("*sity            :")
	fmt.Println("imployment       :", a.imployment)
	fmt.Println("bolt_trip        :", a.bolt_trip)
	fmt.Println("uclon_trip       :", a.uclon_trip)
	fmt.Println("uclon_cash       :", a.uclon_cash)
	fmt.Println("uclon_today      :", a.uclon_today)
	fmt.Println("grafik           :", a.grafik)
	fmt.Println("vod              :", a.vod)
	fmt.Println("kassa            :", a.kassa)
	fmt.Println("save_kassa_BGQ   :", a.save_kassa_BGQ)
	fmt.Println("ident            :", a.ident)
	fmt.Println("id_BGQ           :", a.id_BGQ)
	fmt.Println("name             :", a.name)
	fmt.Println("name_vod         :", a.name_vod)
	fmt.Println("name_online      :", a.name_online)
	fmt.Println("name_grafik      :", a.name_grafik)
	fmt.Println("sity_report      :", a.sity_report)
	fmt.Println("online_diapazon  :", a.online_diapazon)
	fmt.Println("peregon          :", a.peregon)
	fmt.Println("peregon_list     :", a.peregon_list)
	fmt.Println("nov_vod          :", a.nov_vod)
	fmt.Println("mapon_bolt       :", a.mapon_bolt)
	fmt.Println("mapon_uclon      :", a.mapon_uclon)
	fmt.Println("mapon_bolt_id    :", a.mapon_bolt_id)
	fmt.Println("mapon_uclon_id   :", a.mapon_uclon_id)
	fmt.Println("srez             :", a.srez)
	fmt.Println("DM               :", a.DM)
	fmt.Println("Kvoliti          :", a.Kvoliti)
	fmt.Println("norma_trip_hour  :", a.norma_trip_hour)
	fmt.Println("top_vodila       :", a.top_vodila)
	fmt.Println("srpob            :", a.srpob)
	fmt.Println("preprobeg        :", a.preprobeg)
	fmt.Println("workhour         :", a.workhour)
	fmt.Println("id               :", a.id)
	fmt.Println("procent_bolt     :", a.procent_bolt)
	fmt.Println("procent_uclon    :", a.procent_uclon)
	fmt.Println("procent_vozvrat  :", a.procent_vozvrat)
	fmt.Println("procent_vod      :", a.procent_vod)
	fmt.Println("cena_gas         :", a.cena_gas)
	fmt.Println("rasxod_avto      :", a.rasxod_avto)
	fmt.Println("cmp_grf          :", a.cmp_grf)
	fmt.Println("ext_rep          :", a.ext_rep)
	fmt.Println("in_report        :", a.in_report)

}

func (cit *sity) Search(s string) {
	City_Park_G_Car.RLock()
	for _, m := range City_Park_G_Car.m {
		if s == m.id_BGQ || s == m.name {
			*cit = m
			break
		}
	}
	City_Park_G_Car.RUnlock()
}

func (cit *sity) Create(r []string) {
	cit.name = r[0]
	cit.id_BGQ = r[1]

	cit.imployment.name_dataset = r[3]
	cit.imployment.name_table = r[4]

	cit.bolt_trip.name_dataset = r[5]
	cit.bolt_trip.name_table = r[6]

	cit.uclon_trip.name_dataset = r[7]
	cit.uclon_trip.name_table = r[8]

	cit.uclon_cash.name_dataset = r[9]
	cit.uclon_cash.name_table = r[10]

	cit.vod.name_dataset = r[11]
	cit.vod.name_table = r[12]

	cit.kassa.name_dataset = r[13]
	cit.kassa.name_table = r[14]

	if r[15] == "TRUE" {
		cit.save_kassa_BGQ = true
	}
	cit.ident = r[16]
	cit.id_BGQ = r[17]
	cit.name = r[18]
	cit.name_vod = r[19]
	cit.name_online = r[20]
	cit.name_grafik = r[21]
	cit.sity_report = r[22]
	cit.online_diapazon = r[23]
	cit.peregon = r[24]
	cit.peregon_list = r[25]
	cit.nov_vod = r[26]
	cit.mapon_bolt = r[27]
	cit.mapon_uclon = r[28]
	cit.mapon_bolt_id = r[29]
	cit.mapon_uclon_id = r[30]
	cit.srez = r[31]
	cit.DM = r[32]
	cit.Kvoliti = r[33]

	cit.norma_trip_hour = trs.String_to_float(r[34])
	cit.top_vodila = trs.String_to_int(r[35])
	cit.srpob = trs.String_to_int(r[36])
	cit.preprobeg = trs.String_to_int(r[37])
	cit.workhour = trs.String_to_int(r[38])
	cit.id = trs.String_to_int(r[39])
	cit.procent_bolt = trs.String_to_int(r[40])
	cit.procent_uclon = trs.String_to_int(r[41])
	cit.procent_vozvrat = trs.String_to_int(r[42])
	cit.procent_vod = trs.String_to_int(r[43])
	cit.cena_gas = trs.String_to_int(r[44])
	cit.rasxod_avto = trs.String_to_int(r[45])

	if r[46] == "TRUE" {
		cit.cmp_grf = true
	}
	if r[47] == "TRUE" {
		cit.ext_rep = true
	}
	if r[48] == "TRUE" {
		cit.in_report = true
	}

	cit.grafik.name_dataset = r[49]
	cit.grafik.name_table = r[50]

	cit.uclon_today.name_dataset = r[51]
	cit.uclon_today.name_table = r[52]

}

func (sit sity) norma_name() string {
	l := 25 - len([]rune(sit.name))
	s := sit.name
	if l > 0 {
		for i := 0; i < l; i++ {
			s = s + " "
		}
	}
	return s
}

func create_sity_List() {

	CIT = make([]sity, 0)

	list, err := trs.Read_sheets_CASHBOX_err("10JOrpQU63VyK6VL7EEdj1iiyk6sNzbDaGzW7Fqg2eMs", "'Citys'", 0)
	if err != nil {
		fmt.Println(err)
		return

	}

	for _, i := range list {
		if i[2] != "YES" {
			continue
		}
		cit := sity{}
		cit.Create(i)
		CIT = append(CIT, cit)
	}

}

func create_sitys_List() {

	for {
		res := make([]sity, 0)
		list, err := trs.Read_sheets_CASHBOX_err("10JOrpQU63VyK6VL7EEdj1iiyk6sNzbDaGzW7Fqg2eMs", "'Citys'", 0)
		if err != nil {
			errorLog.Println(err)
			<-time.After(time.Minute)
			continue
		}

		for _, i := range list {
			if i[2] != "YES" {
				continue
			}
			cit := sity{}
			cit.Create(i)
			res = append(res, cit)
		}

		City_Park_G_Car.Lock()
		City_Park_G_Car.m = City_Park_G_Car.m[:0]

		City_Park_G_Car.m = res

		City_Park_G_Car.Unlock()
		<-time.After(time.Hour)
	}
}

type employment struct {
	employ    [24][3]int
	date      string
	name      string
	data_s    string
	work_time float64
}

type mapon_groups struct {
	groups []string
	id     string
	nomer  string
}

type TerraformResource struct {
	Cloud               string    // 16 Bytes
	Name                time.Time // 16 Bytes
	HaveDSL             bool      //  1 Byte
	PluginVersion       string    // 16 Bytes
	IsVersionControlled bool      //  1 Byte
	TerraformVersion    string    // 16 Bytes
	ModuleVersionMajor  int32     //  4 Bytes

}

func optimized() {
	var d TerraformResource
	d.Cloud = "aws"
	d.Name = time.Now()
	d.HaveDSL = true
	d.PluginVersion = "3.64"
	d.TerraformVersion = "1.1"
	d.ModuleVersionMajor = 1
	d.IsVersionControlled = true
	fmt.Println("==============================================================")
	fmt.Printf("Total Memory Usage StructType:d %T => [%d]\n", d, unsafe.Sizeof(d))
	fmt.Println("==============================================================")
	fmt.Printf("Cloud Field StructType:d.Cloud %T => [%d]\n", d.Cloud, unsafe.Sizeof(d.Cloud))
	fmt.Printf("Name Field StructType:d.Name %T => [%d]\n", d.Name, unsafe.Sizeof(d.Name))
	fmt.Printf("HaveDSL Field StructType:d.HaveDSL %T => [%d]\n", d.HaveDSL, unsafe.Sizeof(d.HaveDSL))
	fmt.Printf("PluginVersion Field StructType:d.PluginVersion %T => [%d]\n", d.PluginVersion, unsafe.Sizeof(d.PluginVersion))
	fmt.Printf("ModuleVersionMajor Field StructType:d.IsVersionControlled %T => [%d]\n", d.IsVersionControlled, unsafe.Sizeof(d.IsVersionControlled))
	fmt.Printf("TerraformVersion Field StructType:d.TerraformVersion %T => [%d]\n", d.TerraformVersion, unsafe.Sizeof(d.TerraformVersion))
	fmt.Printf("ModuleVersionMajor Field StructType:d.ModuleVersionMajor %T => [%d]\n", d.ModuleVersionMajor, unsafe.Sizeof(d.ModuleVersionMajor))
}
