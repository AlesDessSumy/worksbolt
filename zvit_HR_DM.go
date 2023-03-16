package main

import (
	"Citys"
	"database/sql"
	"fmt"
	"sort"
	"time"
	"trs"
)

type Zvit_HRDM_Avto_Car struct {
	City     string
	Name     string
	Car_01   int
	Car_yst  int
	Car_aver float64
}

// количество авто за месяц, среднее, первый и последний день
// для плана
func Zvit_HRDM_Create_Time() {

	res := make([][]string, 0)
	res = append(res, []string{"Місто", "Дата порушення", "Час запису звіту до бази даних", "Спізнення, хвилин", "Штраф на менеджерів парку"})

	db, err := sql.Open("sqlite3", path+"sqlite/ZvitHRDM.db")
	if err != nil {
		fmt.Println("Запись", err)
	}
	defer db.Close()

	day := mounth_until_today()

	r := Read_DB_ZvitHRDM(db, day)

	for _, cit := range Citys.Create_sitys_List() {
		for _, i := range r {
			if cit.Name == i.City {
				if i.CreatedTime.Hour() > 19 {
					fmt.Println(i.City, i.Dates, (i.CreatedTime.Hour()-20)*60+i.CreatedTime.Minute())
					res = append(res, []string{i.City, i.Dates, i.CreatedTime.Format("15:04:05"), fmt.Sprint((i.CreatedTime.Hour()-20)*60 + i.CreatedTime.Minute()), "200"})
				}
				// if i.CreatedTime.Hour() < 16 {
				// 	fmt.Println(i.City, i.Dates, (i.CreatedTime.Hour()-20)*60+i.CreatedTime.Minute())
				// 	res = append(res, []string{i.City, i.Dates, i.CreatedTime.Format("15:04:05"), fmt.Sprint(-(i.CreatedTime.Hour()-16)*60 - i.CreatedTime.Minute()), "100"})
				// }
			}
		}
	}

	trs.Rec("1hKEXOZENKgAERhZDponXcZT7gxKdqPFrLUir0nmlf9k", "'Pain'!A1", res, true, true)
}

// количество авто за месяц, среднее, первый и последний день
// для плана
func Zvit_HRDM_Avto() (res []Zvit_HRDM_Avto_Car) {
	db, err := sql.Open("sqlite3", path+"sqlite/ZvitHRDM.db")
	if err != nil {
		fmt.Println("Запись", err)
	}
	defer db.Close()

	day := mounth_until_today()

	r := Read_DB_ZvitHRDM(db, day)

	for _, cit := range Citys.Create_sitys_List() {
		rs := Zvit_HRDM_Avto_Car{City: cit.Id_BGQ, Name: cit.Name}
		rss := make([]ZvitHRDM, 0)
		for _, dt := range day {
			for _, i := range r {
				if int(dt.Weekday()) == 6 {
					if cit.Name == i.City && dt.AddDate(0, 0, -1).Format(TNS) == i.Dates {
						rss = append(rss, i)
					}
				} else if int(dt.Weekday()) == 0 {
					if cit.Name == i.City && dt.AddDate(0, 0, -2).Format(TNS) == i.Dates {
						rss = append(rss, i)
					}
				} else {
					if cit.Name == i.City && dt.Format(TNS) == i.Dates {
						rss = append(rss, i)
					}
				}
			}
		}

		sort.Slice(rss, func(i, j int) (less bool) {
			return rss[i].Date.Before(rss[j].Date)
		})
		if len(rss) == 0 {
			continue
		}
		rs.Car_01 = rss[0].Prm11
		rs.Car_yst = rss[len(rss)-1].Prm11
		n := 0
		summ := 0
		for _, i := range rss {
			// fmt.Println(i)
			n++
			summ = summ + i.Prm11
		}
		rs.Car_aver = float64(summ) / float64(n)
		res = append(res, rs)
		// fmt.Println(rs)
	}
	return
}

// отчет история за период
func Zvit_HRDM_History(z, d1, d2 string) {

	db, err := sql.Open("sqlite3", path+"sqlite/ZvitHRDM.db")
	if err != nil {
		fmt.Println("Запись", err)
	}
	defer db.Close()

	r := Read_DB_ZvitHRDM(db, List_Day(d1, d2))

	if z == "Усі міста сума за період" {

		r = Parse_ZvitHRDM_Summa(r)
		r = Result_ZvitHRDM(r, Total_ZvitHRDM(r))
		res := Result_SliseString_ZvitHRDM(r)
		trs.Rec_Clear("1hKEXOZENKgAERhZDponXcZT7gxKdqPFrLUir0nmlf9k", "'History_ZvitHRDM'!A3", res, true, true)

	} else if z == "Усі міста середнє за період" {
		r = Parse_ZvitHRDM_Aver(r)
		r = Result_ZvitHRDM(r, Total_ZvitHRDM(r))
		res := Result_SliseString_ZvitHRDM(r)
		trs.Rec_Clear("1hKEXOZENKgAERhZDponXcZT7gxKdqPFrLUir0nmlf9k", "'History_ZvitHRDM'!A3", res, true, true)
	} else {
		r = Parse_ZvitHRDM_City(r, z)
		tot := Total_ZvitHRDM(r)
		tot.Div_sm(len(r))
		r = Result_ZvitHRDM(r, tot)
		res := Result_SliseString_ZvitHRDM(r)
		trs.Rec_Clear("1hKEXOZENKgAERhZDponXcZT7gxKdqPFrLUir0nmlf9k", "'History_ZvitHRDM'!A3", res, true, true)
	}
}

// отчет история за период сокращенный
func Zvit_HRDM_History_Small(z, d1, d2 string) {

	db, err := sql.Open("sqlite3", path+"sqlite/ZvitHRDM.db")
	if err != nil {
		fmt.Println("Запись", err)
	}
	defer db.Close()

	r := Read_DB_ZvitHRDM(db, List_Day(d1, d2))

	if z == "Усі міста сума за період" {

		r = Parse_ZvitHRDM_Summa(r)
		r = Result_ZvitHRDM(r, Total_ZvitHRDM(r))
		res := Result_SliseString_ZvitHRDM_Small(r)
		trs.Rec_Clear("1hKEXOZENKgAERhZDponXcZT7gxKdqPFrLUir0nmlf9k", "'History_Small'!A3", res, true, true)

	} else if z == "Усі міста середнє за період" {
		r = Parse_ZvitHRDM_Aver(r)
		r = Result_ZvitHRDM(r, Total_ZvitHRDM(r))
		res := Result_SliseString_ZvitHRDM_Small(r)
		trs.Rec_Clear("1hKEXOZENKgAERhZDponXcZT7gxKdqPFrLUir0nmlf9k", "'History_Small'!A3", res, true, true)
	} else {
		r = Parse_ZvitHRDM_City(r, z)
		r = Result_ZvitHRDM(r, Total_ZvitHRDM(r))
		res := Result_SliseString_ZvitHRDM_Small(r)
		trs.Rec_Clear("1hKEXOZENKgAERhZDponXcZT7gxKdqPFrLUir0nmlf9k", "'History_Small'!A3", res, true, true)
	}
}

// сервер отчета HR и рекрутера
func Server_ZvitHRDM() {

	db, err := sql.Open("sqlite3", path+"sqlite/ZvitHRDM.db")
	if err != nil {
		fmt.Println("Запись", err)
	}
	defer db.Close()
	n := 0
	for {
		t := time.Now()
		if t.Hour() == 5 && t.Minute() == 0 {
			Clear_ZvitHRDM()
			<-time.After(100 * time.Second)
		}
		if t.Hour() < 9 {
			<-time.After(30 * time.Second)
			continue
		}
		fmt.Println("Check Server_ZvitHRDM()", n, time.Now().Format(TNSF))

		for _, cit := range List_City_ZvitHRDM {
			// if cit != "Суми" {
			// 	continue
			// }
			rang := fmt.Sprintf("'%s'!A1:D", cit)
			rangA2 := fmt.Sprintf("'%s'!A2", cit)
			rangD2 := fmt.Sprintf("'%s'!D2", cit)
			list, err := trs.Read_sheets_CASHBOX_err("1Ee-aLEC19a2So8B_4NFSYaBmxzdfzYfOy_dl3qEK9Tk", rang, 0)
			if err != nil {
				fmt.Println(err)
			}
			if len(list) == 0 {
				continue
			}
			if len(list) < 52 {
				continue
			}
			if len(list[0]) < 4 {
				continue
			}

			r := ZvitHRDM{}
			r.Create(list)
			nayavn := r.Presence_DB_ZvitHRDM(db)

			// если запись уже есть в базе
			if list[1][0] == "Звіт записано до бази даних" && nayavn {
				if r.RecordDB {
					trs.Rec("1Ee-aLEC19a2So8B_4NFSYaBmxzdfzYfOy_dl3qEK9Tk", rangD2, [][]string{{"Ні"}}, false, false)
				}
				continue
			} else if list[1][0] != "Звіт записано до бази даних" && nayavn {
				trs.Rec("1Ee-aLEC19a2So8B_4NFSYaBmxzdfzYfOy_dl3qEK9Tk", rangA2, [][]string{{"Звіт записано до бази даних"}}, false, false)
				if r.RecordDB {
					trs.Rec("1Ee-aLEC19a2So8B_4NFSYaBmxzdfzYfOy_dl3qEK9Tk", rangD2, [][]string{{"Ні"}}, false, false)
				}
				continue
			}

			// fmt.Println(r)
			rsumm := r.Suma()
			send, ok := r.Check()
			fmt.Println(cit, send, ok)

			// если лист не заполнен
			if rsumm == 0 || r.Prm11 == 0 {
				trs.Rec("1Ee-aLEC19a2So8B_4NFSYaBmxzdfzYfOy_dl3qEK9Tk", rangA2, [][]string{{"Заповніть будь ласка форму"}}, false, false)
				continue
			} else {
				trs.Rec("1Ee-aLEC19a2So8B_4NFSYaBmxzdfzYfOy_dl3qEK9Tk", rangA2, [][]string{{send}}, false, false)
			}

			if nayavn && r.RecordDB {
				trs.Rec("1Ee-aLEC19a2So8B_4NFSYaBmxzdfzYfOy_dl3qEK9Tk", rangA2, [][]string{{"Звіт записано до бази даних"}}, false, false)
				trs.Rec("1Ee-aLEC19a2So8B_4NFSYaBmxzdfzYfOy_dl3qEK9Tk", rangD2, [][]string{{"Ні"}}, false, false)
			} else if rsumm != 0 && ok && r.RecordDB && r.Prm11 != 0 {
				fmt.Println("Пишем в базу", cit, send, ok)
				err := r.RecDB(db)
				if err == nil {
					trs.Rec("1Ee-aLEC19a2So8B_4NFSYaBmxzdfzYfOy_dl3qEK9Tk", rangA2, [][]string{{"Звіт записано до бази даних"}}, false, false)
					trs.Rec("1Ee-aLEC19a2So8B_4NFSYaBmxzdfzYfOy_dl3qEK9Tk", rangD2, [][]string{{"Ні"}}, false, false)
				} else {
					fmt.Println(err)
				}
			} else if (!ok || rsumm == 0 || r.Prm11 == 0) && r.RecordDB {
				trs.Rec("1Ee-aLEC19a2So8B_4NFSYaBmxzdfzYfOy_dl3qEK9Tk", rangD2, [][]string{{"Ні"}}, false, false)
			}
		}

		fmt.Println(n, n%10)
		if n%10 != 0 {

			<-time.After(8 * time.Second)
			n++
			continue
		}
		n++

		r := Read_DB_ZvitHRDM(db, last_days(0))

		rs := make([][]string, 0)
		res := make([][]string, 0)
		res = append(res, zag_ZvitHRDM)

		rs_sm := make([][]string, 0)
		res_sm := make([][]string, 0)
		res_sm = append(res_sm, zag_ZvitHRDM_Small)
		total := ZvitHRDM{City: "Total", Dates: time.Now().Format(TNS)}

		for _, cit := range List_City_ZvitHRDM {
			fl := 0
			for _, i := range r {
				if cit == i.City {
					rs = append(rs, i.Slise_String())
					rs_sm = append(rs_sm, i.Slise_String_Small())
					total.Summ(i)
					fl = 1
					break
				}
			}
			if fl == 0 {
				str := make([]string, 53)
				str_sm := make([]string, 16)
				// fmt.Println(str)
				str[0] = cit
				str[1] = time.Now().Format(TNS)
				str[2] = "false"

				str_sm[0] = cit
				str_sm[1] = time.Now().Format(TNS)
				str_sm[2] = "false"
				// fmt.Println(str)
				rs = append(rs, str)
				rs_sm = append(rs_sm, str_sm)
			}
		}
		res = append(res, total.Slise_String())
		res = append(res, rs...)

		res_sm = append(res_sm, total.Slise_String_Small())
		res_sm = append(res_sm, rs_sm...)

		trs.Rec("1hKEXOZENKgAERhZDponXcZT7gxKdqPFrLUir0nmlf9k", "'ZvitHRDM'!A1", res, true, true)
		trs.Rec("1hKEXOZENKgAERhZDponXcZT7gxKdqPFrLUir0nmlf9k", "'Small'!A1", res_sm, true, true)

		del, _ := trs.Read_sheets_CASHBOX_err("1hKEXOZENKgAERhZDponXcZT7gxKdqPFrLUir0nmlf9k", "'Delete'!A1:B", 0)
		for n, i := range del {
			if len(i) < 2 {
				continue
			}
			if i[0] == "" {
				break
			}
			if i[1] == "Y" {
				Delete_ZvitHRDM(db, i[0], time.Now().Format(TNS))
				rang := fmt.Sprintf("'Delete'!B%d", n+1)
				rang1 := fmt.Sprintf("'%s'!A2", i[0])
				trs.Rec("1hKEXOZENKgAERhZDponXcZT7gxKdqPFrLUir0nmlf9k", rang, [][]string{{"N"}}, false, false)
				trs.Rec("1Ee-aLEC19a2So8B_4NFSYaBmxzdfzYfOy_dl3qEK9Tk", rang1, [][]string{{"Заповніть будь ласка форму"}}, false, false)

			}
		}

		<-time.After(8 * time.Second)
	}
}

// очистка форм ДМ и рекрутера
func Clear_ZvitHRDM() {

	for _, cit := range List_City_ZvitHRDM {
		res := Clear_ZvitHRDM_Slise
		res[0][0] = time.Now().Format(TNS)
		res[0][1] = cit
		rang := fmt.Sprintf("'%s'!A1", cit)
		trs.Rec("1Ee-aLEC19a2So8B_4NFSYaBmxzdfzYfOy_dl3qEK9Tk", rang, res, false, false)
		fmt.Println("Clear_ZvitHRDM", cit, "Passsed OK")
	}

}

type ZvitHRDM struct {
	City        string
	Date        time.Time
	CreatedTime time.Time
	Dates       string
	RecordDB    bool
	Prm11       int // Кількість авто у парку
	Prm12       int // кількість авто в робочому стані на 17.00
	Prm13       int // кількість авто не в робочому стані на 17.00
	Prm14       int // авто прийшло у парк на 17.00
	Prm15       int // авто вибуло з парку на 17.00
	Prm16       int // скільки авто вийшло у роботу сьогодні на 17.00
	Prm17       int // скільки авто вийшло з роботи сьогодні на 17.00
	Prm21       int // Авто у роботі
	Prm22       int // авто без водія на 17.00
	Prm23       int // вихідних водіїв сьогодні  на 17.00
	Prm24       int // водіїв на лікарняному на 17.00
	Prm31       int // Звільнено водіїв на 17.00
	Prm32       int // знайшов іншу роботу
	Prm33       int // не влаштовує рівень зарабітку у парку
	Prm34       int // перехід на авто у оренду, або власне авто
	Prm35       int // переїзд до іншого міста
	Prm36       int // призив до лав ЗСУ
	Prm37       int // не його робота (не хоче бути водієм таксі)
	Prm38       int // не вдається поєднувати з іншою роботою
	Prm390      int // алкогольна або наркозалежність
	Prm391      int // систематичне порушення режиму роботи
	Prm41       int // Здали авто тимчасово водіїв на 17.00
	Prm42       int // водій захворів
	Prm43       int // тимчасово не буде у місті
	Prm44       int // відпустка водія
	Prm45       int // за сімейних обставин
	Prm46       int // навчання водія
	Prm47       int // Видано авто сьогодні на 17.00
	Prm48       int // виїхало нових водіїв сьогодні на 17.00
	Prm49       int // виїхало з вже працюючих водіїв сьогодні на 17.00
	Prm50       int // Виїхало нових водіїв сьогодні на 17.00
	Prm51       int // віїхало з незаписаних на співбесіду:
	Prm52       int // віїхало з призначених на інші дні:
	Prm53       int // віїхало з призначених на сьогодні:
	Prm54       int // перенесено виїзд на інший день:
	Prm61       int // Призначено співбесід на сьгодні:
	Prm62       int // проведено співбесід з призначених на сьогодні:
	Prm63       int // проведено співбесід з незаписаних на співбесіду:
	Prm64       int // проведено співбесід з призначених на інші дні:
	Prm65       int // перенесено співбесід на інший день:
	Prm66       int // кандидата переведено до резерву:
	Prm68       int // кандидат не прийшов:
	Prm69       int // кандидат переніс співбесіду:
	Prm71       int // Кандидат відмовився працевлаштовуватися:
	Prm72       int // кандидат не підійшов парку:
	Prm73       int // кандидата не задовільняє рівень заробітку:
	Prm74       int // кандидату вакансія неактуальна:
	Prm75       int // не підійшли умови праці:
	Prm76       int // у кандидата немає грошей на депозит:
	Prm77       int // у кандидата немає грошей на пальне:
}

func (a *ZvitHRDM) Suma() int {
	r := a.Prm11 + a.Prm12 + a.Prm13 + a.Prm14 + a.Prm15 + a.Prm16 + a.Prm17
	r = r + a.Prm21 + a.Prm22 + a.Prm23 + a.Prm24
	r = r + a.Prm31 + a.Prm32 + a.Prm33 + a.Prm34 + a.Prm35 + a.Prm36 + a.Prm37 + a.Prm38 + a.Prm390 + a.Prm391
	r = r + a.Prm41 + a.Prm42 + a.Prm43 + a.Prm44 + a.Prm45 + a.Prm46 + a.Prm47 + a.Prm48 + a.Prm49
	r = r + a.Prm50 + a.Prm51 + a.Prm52 + a.Prm53 + a.Prm54
	r = r + a.Prm61 + a.Prm62 + a.Prm63 + a.Prm64 + a.Prm65 + a.Prm66 + a.Prm68 + a.Prm69
	r = r + a.Prm71 + a.Prm72 + a.Prm73 + a.Prm74 + a.Prm75 + a.Prm76 + a.Prm77
	return r
}

func (a *ZvitHRDM) Slise_String() []string {
	b := make([]string, 53)
	b[0] = a.City
	b[1] = a.Dates
	b[2] = fmt.Sprint(a.RecordDB)
	b[3] = fmt.Sprint(a.Prm11)
	b[4] = fmt.Sprint(a.Prm12)
	b[5] = fmt.Sprint(a.Prm13)
	b[6] = fmt.Sprint(a.Prm14)
	b[7] = fmt.Sprint(a.Prm15)
	b[8] = fmt.Sprint(a.Prm16)
	b[9] = fmt.Sprint(a.Prm17)
	b[10] = fmt.Sprint(a.Prm21)
	b[11] = fmt.Sprint(a.Prm22)
	b[12] = fmt.Sprint(a.Prm23)
	b[13] = fmt.Sprint(a.Prm24)
	b[14] = fmt.Sprint(a.Prm31)
	b[15] = fmt.Sprint(a.Prm32)
	b[16] = fmt.Sprint(a.Prm33)
	b[17] = fmt.Sprint(a.Prm34)
	b[18] = fmt.Sprint(a.Prm35)
	b[19] = fmt.Sprint(a.Prm36)
	b[20] = fmt.Sprint(a.Prm37)
	b[21] = fmt.Sprint(a.Prm38)
	b[22] = fmt.Sprint(a.Prm390)
	b[23] = fmt.Sprint(a.Prm391)
	b[24] = fmt.Sprint(a.Prm41)
	b[25] = fmt.Sprint(a.Prm42)
	b[26] = fmt.Sprint(a.Prm43)
	b[27] = fmt.Sprint(a.Prm44)
	b[28] = fmt.Sprint(a.Prm45)
	b[29] = fmt.Sprint(a.Prm46)
	b[30] = fmt.Sprint(a.Prm47)
	b[31] = fmt.Sprint(a.Prm48)
	b[32] = fmt.Sprint(a.Prm49)
	b[33] = fmt.Sprint(a.Prm50)
	b[34] = fmt.Sprint(a.Prm51)
	b[35] = fmt.Sprint(a.Prm52)
	b[36] = fmt.Sprint(a.Prm53)
	b[37] = fmt.Sprint(a.Prm54)
	b[38] = fmt.Sprint(a.Prm61)
	b[39] = fmt.Sprint(a.Prm62)
	b[40] = fmt.Sprint(a.Prm63)
	b[41] = fmt.Sprint(a.Prm64)
	b[42] = fmt.Sprint(a.Prm65)
	b[43] = fmt.Sprint(a.Prm66)
	b[44] = fmt.Sprint(a.Prm68)
	b[45] = fmt.Sprint(a.Prm69)
	b[46] = fmt.Sprint(a.Prm71)
	b[47] = fmt.Sprint(a.Prm72)
	b[48] = fmt.Sprint(a.Prm73)
	b[49] = fmt.Sprint(a.Prm74)
	b[50] = fmt.Sprint(a.Prm75)
	b[51] = fmt.Sprint(a.Prm76)
	b[52] = fmt.Sprint(a.Prm77)
	return b
}

// % уволенных водителей от авто в работе
// % уволенных водителей от авто в парке
// % сдавших водителей от авто в работе
// % сдавших водителей от авто в парке

// % выехавших водителей от авто в работе
// % выехавшихводителей от авто в парке

// % приходу от призначених
// % приходу загалом от призначених

// % кандидат не подошел парку от пришедших загалом

// % авто в работе от авто в парке

// %  лыкарняних
// % вихыдних
// % вихыдних лыкарняних

// Кандидат відмовився працевлаштовуватися:
// кандидат не підійшов парку:
// кандидата не задовільняє рівень заробітку:
// кандидату вакансія неактуальна:
// не підійшли умови праці:
// у кандидата немає грошей на депозит:
// у кандидата немає грошей на пальне:

func (a *ZvitHRDM) Slise_String_Small() []string {
	b := make([]string, 53)
	b[0] = a.City
	b[1] = a.Dates
	b[2] = fmt.Sprint(a.RecordDB)
	b[3] = fmt.Sprint(a.Prm11)
	b[4] = fmt.Sprint(a.Prm12)
	b[5] = fmt.Sprint(a.Prm13)
	b[6] = fmt.Sprint(a.Prm21)
	b[7] = fmt.Sprint(a.Prm22)
	b[8] = fmt.Sprint(a.Prm31)
	b[9] = fmt.Sprint(a.Prm41)
	b[10] = fmt.Sprint(a.Prm47)
	b[11] = fmt.Sprint(a.Prm48)
	b[12] = fmt.Sprint(a.Prm49)
	b[13] = fmt.Sprint(a.Prm61)
	b[14] = fmt.Sprint(a.Prm62 + a.Prm63 + a.Prm64)
	b[15] = fmt.Sprint(a.Prm51 + a.Prm52 + a.Prm53)
	b[16] = fmt.Sprint(a.Prm71)
	// %
	if a.Prm21 != 0 { // % уволенных водителей от авто в работе
		b[17] = trs.Float_to_string_2(float64(a.Prm31) / float64(a.Prm21) * 100)
	}
	if a.Prm11 != 0 { // % уволенных водителей от авто в парке
		b[18] = trs.Float_to_string_2(float64(a.Prm31) / float64(a.Prm11) * 100)
	}
	if a.Prm21 != 0 { // % сдавших водителей от авто в работе
		b[19] = trs.Float_to_string_2(float64(a.Prm41) / float64(a.Prm21) * 100)
	}
	if a.Prm11 != 0 { // % сдавших водителей от авто в парке
		b[20] = trs.Float_to_string_2(float64(a.Prm41) / float64(a.Prm11) * 100)
	}
	if a.Prm21 != 0 { // %  выехавших водителей от авто в работе
		b[21] = trs.Float_to_string_2(float64(a.Prm48) / float64(a.Prm21) * 100)
	}
	if a.Prm11 != 0 { // % выехавшихводителей от авто в парке
		b[22] = trs.Float_to_string_2(float64(a.Prm48) / float64(a.Prm11) * 100)
	}
	if a.Prm61 != 0 { // % приходу от призначених
		b[23] = trs.Float_to_string_2(float64(a.Prm62) / float64(a.Prm61) * 100)
	}
	if a.Prm61 != 0 { // % приходу загалом от призначених
		b[24] = trs.Float_to_string_2(float64(a.Prm62+a.Prm63+a.Prm64) / float64(a.Prm61) * 100)
	}
	if a.Prm61 != 0 { //  % выехавших от назначенных на сегодня
		b[25] = trs.Float_to_string_2(float64(a.Prm51+a.Prm52+a.Prm53) / float64(a.Prm61) * 100)
	}
	if a.Prm62 != 0 { //  % выехавших от пришедших на сегодня из записаных на сегодня
		b[26] = trs.Float_to_string_2(float64(a.Prm51+a.Prm52+a.Prm53) / float64(a.Prm62) * 100)
	}
	if a.Prm62-a.Prm72 != 0 { // // % выехавших от (пришедших на сегодня из записаных на сегодня - не подошел)
		b[27] = trs.Float_to_string_2(float64(a.Prm51+a.Prm52+a.Prm53) / float64(a.Prm62-a.Prm72) * 100)
	}
	if a.Prm62+a.Prm63+a.Prm64 != 0 { // % кандидат не подошел парку от пришедших загалом
		b[28] = trs.Float_to_string_2(float64(a.Prm72) / float64(a.Prm62+a.Prm63+a.Prm64) * 100)
	}
	if a.Prm11 != 0 { // % авто в работе от авто в парке
		b[29] = trs.Float_to_string_2(float64(a.Prm21) / float64(a.Prm11) * 100)
	}
	if a.Prm21 != 0 { // % вихыдних
		b[30] = trs.Float_to_string_2(float64(a.Prm23) / float64(a.Prm21) * 100)
	}
	if a.Prm21 != 0 { // лыкарняних
		b[31] = trs.Float_to_string_2(float64(a.Prm24) / float64(a.Prm21) * 100)
	}
	if a.Prm21 != 0 { // лыкарняних + виходных
		b[32] = trs.Float_to_string_2(float64(a.Prm23+a.Prm24) / float64(a.Prm21) * 100)
	}
	if a.Prm21 != 0 { // без водителя
		b[33] = trs.Float_to_string_2(float64(a.Prm22) / float64(a.Prm21) * 100)
	}
	if a.Prm62+a.Prm63+a.Prm64 != 0 {
		b[34] = trs.Float_to_string_2(float64(a.Prm71) / float64(a.Prm62+a.Prm63+a.Prm64) * 100)
		b[35] = trs.Float_to_string_2(float64(a.Prm72) / float64(a.Prm62+a.Prm63+a.Prm64) * 100)
		b[36] = trs.Float_to_string_2(float64(a.Prm73) / float64(a.Prm62+a.Prm63+a.Prm64) * 100)
		b[37] = trs.Float_to_string_2(float64(a.Prm74) / float64(a.Prm62+a.Prm63+a.Prm64) * 100)
		b[38] = trs.Float_to_string_2(float64(a.Prm75) / float64(a.Prm62+a.Prm63+a.Prm64) * 100)
		b[39] = trs.Float_to_string_2(float64(a.Prm76) / float64(a.Prm62+a.Prm63+a.Prm64) * 100)
		b[40] = trs.Float_to_string_2(float64(a.Prm77) / float64(a.Prm62+a.Prm63+a.Prm64) * 100)
	}
	return b
}

func (a *ZvitHRDM) Create(b [][]string) {
	if len(b) == 0 {
		return
	}
	if len(b) < 52 {
		return
	}
	if len(b[0]) < 3 {
		return
	}
	a.City = b[0][1]
	a.CreatedTime = time.Now()
	//eur, _ := time.LoadLocation("Europe/Vienna")

	t, _ := time.Parse(TNS, b[0][0])
	a.Date = time.Date(t.Year(), t.Month(), t.Day(), 17, 0, 0, 0, a.CreatedTime.Location())

	a.Dates = b[0][0]
	if b[1][3] == "Так" {
		a.RecordDB = true
	}
	a.Prm11 = trs.String_to_int(b[2][3])   // Кількість авто у парку
	a.Prm12 = trs.String_to_int(b[3][3])   // кількість авто в робочому стані на 17.00
	a.Prm13 = trs.String_to_int(b[4][3])   // кількість авто не в робочому стані на 17.00
	a.Prm14 = trs.String_to_int(b[5][3])   // авто прийшло у парк на 17.00
	a.Prm15 = trs.String_to_int(b[6][3])   // авто вибуло з парку на 17.00
	a.Prm16 = trs.String_to_int(b[7][3])   // скільки авто вийшло у роботу сьогодні на 17.00
	a.Prm17 = trs.String_to_int(b[8][3])   // скільки авто вийшло з роботи сьогодні на 17.00
	a.Prm21 = trs.String_to_int(b[9][3])   // Авто у роботі
	a.Prm22 = trs.String_to_int(b[10][3])  // авто без водія на 17.00
	a.Prm23 = trs.String_to_int(b[11][3])  // вихідних водіїв сьогодні  на 17.00
	a.Prm24 = trs.String_to_int(b[12][3])  // водіїв на лікарняному на 17.00
	a.Prm31 = trs.String_to_int(b[13][3])  // Звільнено водіїв на 17.00
	a.Prm32 = trs.String_to_int(b[14][3])  // знайшов іншу роботу
	a.Prm33 = trs.String_to_int(b[15][3])  // не влаштовує рівень зарабітку у парку
	a.Prm34 = trs.String_to_int(b[16][3])  // перехід на авто у оренду, або власне авто
	a.Prm35 = trs.String_to_int(b[17][3])  // переїзд до іншого міста
	a.Prm36 = trs.String_to_int(b[18][3])  // призив до лав ЗСУ
	a.Prm37 = trs.String_to_int(b[19][3])  // не його робота (не хоче бути водієм таксі)
	a.Prm38 = trs.String_to_int(b[20][3])  // не вдається поєднувати з іншою роботою
	a.Prm390 = trs.String_to_int(b[21][3]) // алкогольна або наркозалежність
	a.Prm391 = trs.String_to_int(b[22][3]) // систематичне порушення режиму роботи
	a.Prm41 = trs.String_to_int(b[23][3])  // Здали авто тимчасово водіїв на 17.00
	a.Prm42 = trs.String_to_int(b[24][3])  // водій захворів
	a.Prm43 = trs.String_to_int(b[25][3])  // тимчасово не буде у місті
	a.Prm44 = trs.String_to_int(b[26][3])  // відпустка водія
	a.Prm45 = trs.String_to_int(b[27][3])  // за сімейних обставин
	a.Prm46 = trs.String_to_int(b[28][3])  // навчання водія
	a.Prm47 = trs.String_to_int(b[29][3])  // Видано авто сьогодні на 17.00
	a.Prm48 = trs.String_to_int(b[30][3])  // виїхало нових водіїв сьогодні на 17.00
	a.Prm49 = trs.String_to_int(b[31][3])  // виїхало з вже працюючих водіїв сьогодні на 17.00
	a.Prm50 = trs.String_to_int(b[32][2])  // Виїхало нових водіїв сьогодні на 17.00
	a.Prm51 = trs.String_to_int(b[33][2])  // віїхало з незаписаних на співбесіду:
	a.Prm52 = trs.String_to_int(b[34][2])  // віїхало з призначених на інші дні:
	a.Prm53 = trs.String_to_int(b[35][2])  // віїхало з призначених на сьогодні:
	a.Prm54 = trs.String_to_int(b[36][2])  // перенесено виїзд на інший день:
	a.Prm61 = trs.String_to_int(b[37][2])  // Призначено співбесід на сьгодні:
	a.Prm62 = trs.String_to_int(b[38][2])  // проведено співбесід з призначених на сьогодні:
	a.Prm63 = trs.String_to_int(b[39][2])  // проведено співбесід з незаписаних на співбесіду:
	a.Prm64 = trs.String_to_int(b[40][2])  // проведено співбесід з призначених на інші дні:
	a.Prm65 = trs.String_to_int(b[41][2])  // перенесено співбесід на інший день:
	a.Prm66 = trs.String_to_int(b[42][2])  // кандидата переведено до резерву:
	a.Prm68 = trs.String_to_int(b[43][2])  // кандидат не прийшов:
	a.Prm69 = trs.String_to_int(b[44][2])  // кандидат переніс співбесіду:
	a.Prm71 = trs.String_to_int(b[45][2])  // Кандидат відмовився працевлаштовуватися:
	a.Prm72 = trs.String_to_int(b[46][2])  // кандидат не підійшов парку:
	a.Prm73 = trs.String_to_int(b[47][2])  // кандидата не задовільняє рівень заробітку:
	a.Prm74 = trs.String_to_int(b[48][2])  // кандидату вакансія неактуальна:
	a.Prm75 = trs.String_to_int(b[49][2])  // не підійшли умови праці:
	a.Prm76 = trs.String_to_int(b[50][2])  // у кандидата немає грошей на депозит:
	a.Prm77 = trs.String_to_int(b[51][2])  // у кандидата немає грошей на пальне:
}

// проверка формы
func (a *ZvitHRDM) Check() (string, bool) {
	if a.Prm11 == 0 {
		return "Заповніть будь ласка форму", false
	}
	if a.Prm11 != a.Prm12+a.Prm13 {
		return "Кількість авто у парку повинна дорівнювати сумі (авто в робочому стані) + (авто не в робочому стані)", false
	}
	if a.Prm12 != a.Prm21+a.Prm22 {
		return "Кількість авто в робочому стані повинна дорівнювати сумі (авто у роботі) + (авто без водія)", false
	}
	if a.Prm31 != a.Prm32+a.Prm33+a.Prm34+a.Prm35+a.Prm36+a.Prm37+a.Prm38+a.Prm390+a.Prm391 {
		return "Кількість звільнених водіїв повинна дорівнювати сумі ((знайшов іншу роботу) + (не влаштовує рівень зарабітку у парку) + ...)", false
	}
	if a.Prm41 != a.Prm42+a.Prm43+a.Prm44+a.Prm45+a.Prm46 {
		return "Кількість здали авто тимчасово водіїв повинна дорівнювати сумі ((водій захворів) + (тимчасово не буде у місті) + (відпустка водія) + ...)", false
	}
	if a.Prm47 != a.Prm48+a.Prm49 {
		return "Видано авто сьогодні повинна дорівнювати сумі (виїхало нових водіїв сьогодні) + (виїхало з вже працюючих водіїв)", false
	}
	if a.Prm48 != a.Prm50 {
		return "Дані по виїзду нових водіїв ДМ та рекрутера не співпадають", false
	}
	if a.Prm50 != a.Prm51+a.Prm52+a.Prm53 {
		return "Виїхало нових водіїв сьогодні повинна дорівнювати сумі ((віїхало з незаписаних на співбесіду) + (віїхало з призначених на інші дні) + (віїхало з призначених на сьогодні)", false
	}
	if a.Prm61 != a.Prm62+a.Prm65+a.Prm68 {
		return "Призначено співбесід на сьгодні повинна дорівнювати сумі ((проведено співбесід з призначених на сьогодні) + (перенесено співбесід на інший день) + (кандидат не прийшов)", false
	}
	if a.Prm22 > 0 && a.Prm66 > 0 {
		return "Не можна переводити кандидата до резерву, якщо є авто без водія", false
	}
	if a.Prm53 != a.Prm62-(a.Prm66+a.Prm54)-a.Prm71 {
		return "Виїхало з призначених на сьогодні повинна дорівнювати виразу ((проведено співбесід з призначених на сьогодні) - (перенесено виїзд на інший день + кандидата переведено до резерву) - (Кандидат відмовився працевлаштовуватися)", false
	}
	if a.Prm71 != a.Prm72+a.Prm73+a.Prm74+a.Prm75+a.Prm76+a.Prm77 {
		return "Кандидат відмовився працевлаштовуватися повинна дорівнювати сумі ((кандидат не підійшов парку) - (кандидата не задовільняє рівень заробітку) - (кандидату вакансія неактуальна) + ...", false
	}

	return "Усе вірно!!!", true
}

// все выехавшие общие
// % выехавших от назначенных на сегодня
// % выехавших от пришедших на сегодня из записаных на сегодня
// % выехавших от (пришедших на сегодня из записаных на сегодня - не подошел)
// поставить за Х

// список городов
var List_City_ZvitHRDM = []string{"Київ", "Дніпро", "Запоріжжя", "Одеса", "Івано Франківськ", "Суми", "Полтава", "Луцьк", "Ужгород", "Черкаси", "Чернігів",
	"Чернівці", "Львів", "Кривий Ріг", "Рівне", "Вінниця", "Житомир", "Хмельницький", "Тернопіль", "Мукачево", "Камянец", "Харків"}

// чистая форма
var Clear_ZvitHRDM_Slise = [][]string{
	{"", "", "Рекрутер", "ДМ"},
	{"Заповніть будь ласка форму", "", "Записати до Бази Даних?", "Ні"},
	{"Кількість авто у парку", "", "", "0"},
	{"", "кількість авто в робочому стані на 17.00", "", "0"},
	{"", "кількість авто не в робочому стані на 17.00", "", "0"},
	{"", "авто прийшло у парк на 17.00", "", "0"},
	{"", "авто вибуло з парку на 17.00", "", "0"},
	{"", "скільки авто вийшло у роботу сьогодні на 17.00", "", "0"},
	{"", "скільки авто вийшло з роботи сьогодні на 17.00", "", "0"},
	{"Авто у роботі", "", "", "0"},
	{"", "авто без водія на 17.00", "", "0"},
	{"", "вихідних водіїв сьогодні  на 17.00", "", "0"},
	{"", "водіїв на лікарняному на 17.00", "", "0"},
	{"Звільнено водіїв на 17.00", "", "", "0"},
	{"", "знайшов іншу роботу", "", "0"},
	{"", "не влаштовує рівень зарабітку у парку", "", "0"},
	{"", "перехід на авто у оренду, або власне авто", "", "0"},
	{"", "переїзд до іншого міста", "", "0"},
	{"", "призив до лав ЗСУ", "", "0"},
	{"", "не його робота (не хоче бути водієм таксі)", "", "0"},
	{"", "не вдається поєднувати з іншою роботою", "", "0"},
	{"", "алкогольна або наркозалежність", "", "0"},
	{"", "систематичне порушення режиму роботи", "", "0"},
	{"Здали авто тимчасово водіїв на 17.00", "", "", "0"},
	{"", "водій захворів", "", "0"},
	{"", "тимчасово не буде у місті", "", "0"},
	{"", "відпустка водія", "", "0"},
	{"", "за сімейних обставин", "", "0"},
	{"", "навчання водія", "", "0"},
	{"Видано авто сьогодні на 17.00", "", "", "0"},
	{"", "виїхало нових водіїв сьогодні на 17.00", "", "0"},
	{"", "виїхало з вже працюючих водіїв сьогодні на 17.00", "", "0"},
	{"Виїхало нових водіїв сьогодні на 17.00", "", "0", ""},
	{"", "виїхало з незаписаних на співбесіду:", "0", ""},
	{"", "виїхало з призначених на інші дні:", "0", ""},
	{"", "виїхало з призначених на сьогодні:", "0", ""},
	{"", "перенесено виїзд на інший день:", "0", ""},
	{"Призначено співбесід на сьгодні:", "", "0", ""},
	{"", "проведено співбесід з призначених на сьогодні:", "0", ""},
	{"", "проведено співбесід з незаписаних на співбесіду:", "0", ""},
	{"", "проведено співбесід з призначених на інші дні:", "0", ""},
	{"", "перенесено співбесід на інший день:", "0", ""},
	{"", "кандидата переведено до резерву:", "0", ""},
	{"", "кандидат не прийшов:", "0", ""},
	{"", "___________________:", "0", ""},
	{"Кандидат не працевлаштувався:", "", "0", ""},
	{"", "кандидат не підійшов парку:", "0", ""},
	{"", "кандидата не задовільняє рівень заробітку:", "0", ""},
	{"", "кандидату вакансія неактуальна:", "0", ""},
	{"", "не підійшли умови праці:", "0", ""},
	{"", "у кандидата немає грошей на депозит:", "0", ""},
	{"", "у кандидата немає грошей на пальне:", "0", ""},
	{"", "", "", ""},
}

// заголовок полного отчета
var zag_ZvitHRDM = []string{
	"Місто",
	"Дата",
	"Записано до БД",
	"Кількість авто у парку",
	"кількість авто в робочому стані на 17.00",
	"кількість авто не в робочому стані на 17.00",
	"авто прийшло у парк на 17.00",
	"авто вибуло з парку на 17.00",
	"скільки авто вийшло у роботу сьогодні на 17.00",
	"скільки авто вийшло з роботи сьогодні на 17.00",
	"Авто у роботі",
	"авто без водія на 17.00",
	"вихідних водіїв сьогодні  на 17.00",
	"водіїв на лікарняному на 17.00",
	"Звільнено водіїв на 17.00",
	"знайшов іншу роботу",
	"не влаштовує рівень зарабітку у парку",
	"перехід на авто у оренду, або власне авто",
	"переїзд до іншого міста",
	"призив до лав ЗСУ",
	"не його робота (не хоче бути водієм таксі)",
	"не вдається поєднувати з іншою роботою",
	"алкогольна або наркозалежність",
	"систематичне порушення режиму роботи",
	"Здали авто тимчасово водіїв на 17.00",
	"водій захворів",
	"тимчасово не буде у місті",
	"відпустка водія",
	"за сімейних обставин",
	"навчання водія",
	"Видано авто сьогодні на 17.00",
	"виїхало нових водіїв сьогодні на 17.00",
	"виїхало з вже працюючих водіїв сьогодні на 17.00",
	"Виїхало нових водіїв сьогодні на 17.00",
	"виїхало з незаписаних на співбесіду:",
	"виїхало з призначених на інші дні:",
	"виїхало з призначених на сьогодні:",
	"перенесено виїзд на інший день:",
	"Призначено співбесід на сьгодні:",
	"проведено співбесід з призначених на сьогодні:",
	"проведено співбесід з незаписаних на співбесіду:",
	"проведено співбесід з призначених на інші дні:",
	"перенесено співбесід на інший день:",
	"кандидата переведено до резерву:",
	"кандидат не прийшов:",
	"кандидат переніс співбесіду:",
	"Кандидат не працевлаштувався:",
	"кандидат не підійшов парку:",
	"кандидата не задовільняє рівень заробітку:",
	"кандидату вакансія неактуальна:",
	"не підійшли умови праці:",
	"у кандидата немає грошей на депозит:",
	"у кандидата немає грошей на пальне:",
}

// заголовок отчета сокращенный
var zag_ZvitHRDM_Small = []string{
	"Місто",
	"Дата",
	"Записано до БД",
	"Кількість авто у парку",
	"кількість авто в робочому стані на 17.00",
	"кількість авто не в робочому стані на 17.00",
	"Авто у роботі",
	"авто без водія на 17.00",
	"Звільнено водіїв на 17.00",
	"Здали авто тимчасово водіїв на 17.00",
	"Видано авто сьогодні на 17.00",
	"виїхало нових водіїв сьогодні на 17.00",
	"виїхало з вже працюючих водіїв сьогодні на 17.00",
	"Призначено співбесід на сьгодні:",
	"Проведено співбесід за сьогодні:",
	"Виїхало загалом сьогодні:",
	"Кандидат не працевлаштувався:",
	"% звільнених водіїв від авто в роботі",
	"% звільнених водіїв від авто в парку",
	"% сдавших водіїв від авто в роботі",
	"% сдавших водіїв від авто в парку",
	"% виїхавших водіїв від авто в роботі",
	"% виїхавших водіїв від авто в парку",
	"% приходу від призначених",
	"% приходу загалом від призначених",
	"% виїхали від призначених на сьогодні",
	"% що виїхали від тих, що прийшли на сьогодні, із записаних на сьогодні",
	"% виїхали від (прийшли сьогодні із записаних сьогодні - не підійшов)",
	"% кандидат не підійшов парку від пришедших загалом",
	"% авто в работі від авто в парку",
	"% вихідних від авто в роботі",
	"% лікарняних від авто в роботі",
	"% вихідних та лікарняних від авто в роботі",
	"% авто без водія від авто в роботі",
	"% кандидат відмовився працевлаштовуватися від пришедших загалом",
	"% кандидат не підійшов парку від пришедших загалом",
	"% кандидата не задовільняє рівень заробітку від пришедших загалом",
	"% кандидату вакансія неактуальна від пришедших загалом",
	"% не підійшли умови праці від пришедших загалом",
	"% у кандидата немає грошей на депозит від пришедших загалом",
	"% у кандидата немає грошей на пальне від пришедших загалом",
}

// Печать элементаZvitHRDM
func (a *ZvitHRDM) Print() {
	fmt.Println()
	fmt.Println("ZvitHRDM            :")
	fmt.Println("City                :", a.City)
	fmt.Println("Date                :", a.Date)
	fmt.Println("CreatedTime         :", a.CreatedTime)
	fmt.Println("Dates               :", a.Dates)
	fmt.Println("RecordDB            :", a.RecordDB)
	fmt.Println("Prm11               :", a.Prm11)
	fmt.Println("Prm12               :", a.Prm12)
	fmt.Println("Prm13               :", a.Prm13)
	fmt.Println("Prm14               :", a.Prm14)
	fmt.Println("Prm15               :", a.Prm15)
	fmt.Println("Prm16               :", a.Prm16)
	fmt.Println("Prm17               :", a.Prm17)
	fmt.Println("Prm21               :", a.Prm21)
	fmt.Println("Prm22               :", a.Prm22)
	fmt.Println("Prm23               :", a.Prm23)
	fmt.Println("Prm24               :", a.Prm24)
	fmt.Println("Prm31               :", a.Prm31)
	fmt.Println("Prm32               :", a.Prm32)
	fmt.Println("Prm33               :", a.Prm33)
	fmt.Println("Prm34               :", a.Prm34)
	fmt.Println("Prm35               :", a.Prm35)
	fmt.Println("Prm36               :", a.Prm36)
	fmt.Println("Prm37               :", a.Prm37)
	fmt.Println("Prm38               :", a.Prm38)
	fmt.Println("Prm390              :", a.Prm390)
	fmt.Println("Prm391              :", a.Prm391)
	fmt.Println("Prm41               :", a.Prm41)
	fmt.Println("Prm42               :", a.Prm42)
	fmt.Println("Prm43               :", a.Prm43)
	fmt.Println("Prm44               :", a.Prm44)
	fmt.Println("Prm45               :", a.Prm45)
	fmt.Println("Prm46               :", a.Prm46)
	fmt.Println("Prm47               :", a.Prm47)
	fmt.Println("Prm48               :", a.Prm48)
	fmt.Println("Prm49               :", a.Prm49)
	fmt.Println("Prm50               :", a.Prm50)
	fmt.Println("Prm51               :", a.Prm51)
	fmt.Println("Prm52               :", a.Prm52)
	fmt.Println("Prm53               :", a.Prm53)
	fmt.Println("Prm54               :", a.Prm54)
	fmt.Println("Prm61               :", a.Prm61)
	fmt.Println("Prm62               :", a.Prm62)
	fmt.Println("Prm63               :", a.Prm63)
	fmt.Println("Prm64               :", a.Prm64)
	fmt.Println("Prm65               :", a.Prm65)
	fmt.Println("Prm66               :", a.Prm66)
	fmt.Println("Prm68               :", a.Prm68)
	fmt.Println("Prm69               :", a.Prm69)
	fmt.Println("Prm71               :", a.Prm71)
	fmt.Println("Prm72               :", a.Prm72)
	fmt.Println("Prm73               :", a.Prm73)
	fmt.Println("Prm74               :", a.Prm74)
	fmt.Println("Prm75               :", a.Prm75)
	fmt.Println("Prm76               :", a.Prm76)
	fmt.Println("Prm77               :", a.Prm77)
}

// Просчет суммы значения
// по массиву а + б ZvitHRDM
func (a *ZvitHRDM) Summ(b ZvitHRDM) {
	a.Prm11 = a.Prm11 + b.Prm11
	a.Prm12 = a.Prm12 + b.Prm12
	a.Prm13 = a.Prm13 + b.Prm13
	a.Prm14 = a.Prm14 + b.Prm14
	a.Prm15 = a.Prm15 + b.Prm15
	a.Prm16 = a.Prm16 + b.Prm16
	a.Prm17 = a.Prm17 + b.Prm17
	a.Prm21 = a.Prm21 + b.Prm21
	a.Prm22 = a.Prm22 + b.Prm22
	a.Prm23 = a.Prm23 + b.Prm23
	a.Prm24 = a.Prm24 + b.Prm24
	a.Prm31 = a.Prm31 + b.Prm31
	a.Prm32 = a.Prm32 + b.Prm32
	a.Prm33 = a.Prm33 + b.Prm33
	a.Prm34 = a.Prm34 + b.Prm34
	a.Prm35 = a.Prm35 + b.Prm35
	a.Prm36 = a.Prm36 + b.Prm36
	a.Prm37 = a.Prm37 + b.Prm37
	a.Prm38 = a.Prm38 + b.Prm38
	a.Prm390 = a.Prm390 + b.Prm390
	a.Prm391 = a.Prm391 + b.Prm391
	a.Prm41 = a.Prm41 + b.Prm41
	a.Prm42 = a.Prm42 + b.Prm42
	a.Prm43 = a.Prm43 + b.Prm43
	a.Prm44 = a.Prm44 + b.Prm44
	a.Prm45 = a.Prm45 + b.Prm45
	a.Prm46 = a.Prm46 + b.Prm46
	a.Prm47 = a.Prm47 + b.Prm47
	a.Prm48 = a.Prm48 + b.Prm48
	a.Prm49 = a.Prm49 + b.Prm49
	a.Prm50 = a.Prm50 + b.Prm50
	a.Prm51 = a.Prm51 + b.Prm51
	a.Prm52 = a.Prm52 + b.Prm52
	a.Prm53 = a.Prm53 + b.Prm53
	a.Prm54 = a.Prm54 + b.Prm54
	a.Prm61 = a.Prm61 + b.Prm61
	a.Prm62 = a.Prm62 + b.Prm62
	a.Prm63 = a.Prm63 + b.Prm63
	a.Prm64 = a.Prm64 + b.Prm64
	a.Prm65 = a.Prm65 + b.Prm65
	a.Prm66 = a.Prm66 + b.Prm66
	a.Prm68 = a.Prm68 + b.Prm68
	a.Prm69 = a.Prm69 + b.Prm69
	a.Prm71 = a.Prm71 + b.Prm71
	a.Prm72 = a.Prm72 + b.Prm72
	a.Prm73 = a.Prm73 + b.Prm73
	a.Prm74 = a.Prm74 + b.Prm74
	a.Prm75 = a.Prm75 + b.Prm75
	a.Prm76 = a.Prm76 + b.Prm76
	a.Prm77 = a.Prm77 + b.Prm77
}

// Просчет суммы значения
// по массиву  а / т ZvitHRDM
func (a *ZvitHRDM) Div_sm(b int) {
	a.Prm11 = int(round(float64(a.Prm11) / float64(b)))
	a.Prm12 = int(round(float64(a.Prm12) / float64(b)))
	a.Prm13 = int(round(float64(a.Prm13) / float64(b)))
	a.Prm21 = int(round(float64(a.Prm21) / float64(b)))
	a.Prm22 = int(round(float64(a.Prm22) / float64(b)))
}

// Просчет суммы значения
// по массиву  а / т ZvitHRDM
func (a *ZvitHRDM) Div(b int) {
	a.Prm11 = int(round(float64(a.Prm11) / float64(b)))
	a.Prm12 = int(round(float64(a.Prm12) / float64(b)))
	a.Prm13 = int(round(float64(a.Prm13) / float64(b)))
	a.Prm14 = int(round(float64(a.Prm14) / float64(b)))
	a.Prm15 = int(round(float64(a.Prm15) / float64(b)))
	a.Prm16 = int(round(float64(a.Prm16) / float64(b)))
	a.Prm17 = int(round(float64(a.Prm17) / float64(b)))
	a.Prm21 = int(round(float64(a.Prm21) / float64(b)))
	a.Prm22 = int(round(float64(a.Prm22) / float64(b)))
	a.Prm23 = int(round(float64(a.Prm23) / float64(b)))
	a.Prm24 = int(round(float64(a.Prm24) / float64(b)))
	a.Prm31 = int(round(float64(a.Prm31) / float64(b)))
	a.Prm32 = int(round(float64(a.Prm32) / float64(b)))
	a.Prm33 = int(round(float64(a.Prm33) / float64(b)))
	a.Prm34 = int(round(float64(a.Prm34) / float64(b)))
	a.Prm35 = int(round(float64(a.Prm35) / float64(b)))
	a.Prm36 = int(round(float64(a.Prm36) / float64(b)))
	a.Prm37 = int(round(float64(a.Prm37) / float64(b)))
	a.Prm38 = int(round(float64(a.Prm38) / float64(b)))
	a.Prm390 = int(round(float64(a.Prm390) / float64(b)))
	a.Prm391 = int(round(float64(a.Prm391) / float64(b)))
	a.Prm41 = int(round(float64(a.Prm41) / float64(b)))
	a.Prm42 = int(round(float64(a.Prm42) / float64(b)))
	a.Prm43 = int(round(float64(a.Prm43) / float64(b)))
	a.Prm44 = int(round(float64(a.Prm44) / float64(b)))
	a.Prm45 = int(round(float64(a.Prm45) / float64(b)))
	a.Prm46 = int(round(float64(a.Prm46) / float64(b)))
	a.Prm47 = int(round(float64(a.Prm47) / float64(b)))
	a.Prm48 = int(round(float64(a.Prm48) / float64(b)))
	a.Prm49 = int(round(float64(a.Prm49) / float64(b)))
	a.Prm50 = int(round(float64(a.Prm50) / float64(b)))
	a.Prm51 = int(round(float64(a.Prm51) / float64(b)))
	a.Prm52 = int(round(float64(a.Prm52) / float64(b)))
	a.Prm53 = int(round(float64(a.Prm53) / float64(b)))
	a.Prm54 = int(round(float64(a.Prm54) / float64(b)))
	a.Prm61 = int(round(float64(a.Prm61) / float64(b)))
	a.Prm62 = int(round(float64(a.Prm62) / float64(b)))
	a.Prm63 = int(round(float64(a.Prm63) / float64(b)))
	a.Prm64 = int(round(float64(a.Prm64) / float64(b)))
	a.Prm65 = int(round(float64(a.Prm65) / float64(b)))
	a.Prm66 = int(round(float64(a.Prm66) / float64(b)))
	a.Prm68 = int(round(float64(a.Prm68) / float64(b)))
	a.Prm69 = int(round(float64(a.Prm69) / float64(b)))
	a.Prm71 = int(round(float64(a.Prm71) / float64(b)))
	a.Prm72 = int(round(float64(a.Prm72) / float64(b)))
	a.Prm73 = int(round(float64(a.Prm73) / float64(b)))
	a.Prm74 = int(round(float64(a.Prm74) / float64(b)))
	a.Prm75 = int(round(float64(a.Prm75) / float64(b)))
	a.Prm76 = int(round(float64(a.Prm76) / float64(b)))
	a.Prm77 = int(round(float64(a.Prm77) / float64(b)))
}

// Восстановление элемента ZvitHRDM
// из строки БД
func (a *ZvitHRDM) Reestablish(p []string) {
	//fmt.Println(p)
	a.City = p[1]
	dat, err := time.Parse(TNSF2, p[2][:19])
	if err != nil {
		fmt.Println(err)
	}
	a.Date = dat
	crt, err := time.Parse(TNSF2, p[3][:19])
	if err != nil {
		fmt.Println(err)
	}
	a.CreatedTime = crt
	a.Dates = p[4]
	if p[5] == "true" {
		a.RecordDB = true
	}
	a.Prm11 = trs.String_to_int(p[6])
	a.Prm12 = trs.String_to_int(p[7])
	a.Prm13 = trs.String_to_int(p[8])
	a.Prm14 = trs.String_to_int(p[9])
	a.Prm15 = trs.String_to_int(p[10])
	a.Prm16 = trs.String_to_int(p[11])
	a.Prm17 = trs.String_to_int(p[12])
	a.Prm21 = trs.String_to_int(p[13])
	a.Prm22 = trs.String_to_int(p[14])
	a.Prm23 = trs.String_to_int(p[15])
	a.Prm24 = trs.String_to_int(p[16])
	a.Prm31 = trs.String_to_int(p[17])
	a.Prm32 = trs.String_to_int(p[18])
	a.Prm33 = trs.String_to_int(p[19])
	a.Prm34 = trs.String_to_int(p[20])
	a.Prm35 = trs.String_to_int(p[21])
	a.Prm36 = trs.String_to_int(p[22])
	a.Prm37 = trs.String_to_int(p[23])
	a.Prm38 = trs.String_to_int(p[24])
	a.Prm390 = trs.String_to_int(p[25])
	a.Prm391 = trs.String_to_int(p[26])
	a.Prm41 = trs.String_to_int(p[27])
	a.Prm42 = trs.String_to_int(p[28])
	a.Prm43 = trs.String_to_int(p[29])
	a.Prm44 = trs.String_to_int(p[30])
	a.Prm45 = trs.String_to_int(p[31])
	a.Prm46 = trs.String_to_int(p[32])
	a.Prm47 = trs.String_to_int(p[33])
	a.Prm48 = trs.String_to_int(p[34])
	a.Prm49 = trs.String_to_int(p[35])
	a.Prm50 = trs.String_to_int(p[36])
	a.Prm51 = trs.String_to_int(p[37])
	a.Prm52 = trs.String_to_int(p[38])
	a.Prm53 = trs.String_to_int(p[39])
	a.Prm54 = trs.String_to_int(p[40])
	a.Prm61 = trs.String_to_int(p[41])
	a.Prm62 = trs.String_to_int(p[42])
	a.Prm63 = trs.String_to_int(p[43])
	a.Prm64 = trs.String_to_int(p[44])
	a.Prm65 = trs.String_to_int(p[45])
	a.Prm66 = trs.String_to_int(p[46])
	a.Prm68 = trs.String_to_int(p[47])
	a.Prm69 = trs.String_to_int(p[48])
	a.Prm71 = trs.String_to_int(p[49])
	a.Prm72 = trs.String_to_int(p[50])
	a.Prm73 = trs.String_to_int(p[51])
	a.Prm74 = trs.String_to_int(p[52])
	a.Prm75 = trs.String_to_int(p[53])
	a.Prm76 = trs.String_to_int(p[54])
	a.Prm77 = trs.String_to_int(p[55])
}

// Запись элемента ZvitHRDM в БД
func (a *ZvitHRDM) RecDB(db *sql.DB) (err error) {

	rc := "false"
	if a.RecordDB {
		rc = "true"
	}

	name := "City, Date, CreatedTime, Dates, RecordDB, Prm11, Prm12, Prm13, Prm14, Prm15, Prm16, Prm17, Prm21, Prm22, Prm23, Prm24, Prm31, Prm32, Prm33, Prm34, Prm35, Prm36, Prm37, Prm38, Prm390, Prm391, Prm41, Prm42, Prm43, Prm44, Prm45, Prm46, Prm47, Prm48, Prm49, Prm50, Prm51, Prm52, Prm53, Prm54, Prm61, Prm62, Prm63, Prm64, Prm65, Prm66, Prm68, Prm69, Prm71, Prm72, Prm73, Prm74, Prm75, Prm76, Prm77"

	z := fmt.Sprintf("INSERT INTO ZvitHRDM (%s) VALUES ('%s', '%s', '%s', '%s', '%s', '%d', '%d', '%d', '%d', '%d', '%d', '%d', '%d', '%d', '%d', '%d', '%d', '%d', '%d', '%d', '%d', '%d', '%d', '%d', '%d', '%d', '%d', '%d', '%d', '%d', '%d', '%d', '%d', '%d', '%d', '%d', '%d', '%d', '%d', '%d', '%d', '%d', '%d', '%d', '%d', '%d', '%d', '%d', '%d', '%d', '%d', '%d', '%d', '%d', '%d')", name, a.City, a.Date, a.CreatedTime, a.Dates, rc, a.Prm11, a.Prm12, a.Prm13, a.Prm14, a.Prm15, a.Prm16, a.Prm17, a.Prm21, a.Prm22, a.Prm23, a.Prm24, a.Prm31, a.Prm32, a.Prm33, a.Prm34, a.Prm35, a.Prm36, a.Prm37, a.Prm38, a.Prm390, a.Prm391, a.Prm41, a.Prm42, a.Prm43, a.Prm44, a.Prm45, a.Prm46, a.Prm47, a.Prm48, a.Prm49, a.Prm50, a.Prm51, a.Prm52, a.Prm53, a.Prm54, a.Prm61, a.Prm62, a.Prm63, a.Prm64, a.Prm65, a.Prm66, a.Prm68, a.Prm69, a.Prm71, a.Prm72, a.Prm73, a.Prm74, a.Prm75, a.Prm76, a.Prm77)

	fmt.Println(z)
	_, err = db.Exec(z)
	return
}

// проверка наличия
// отчета за день по городу
func (a *ZvitHRDM) Presence_DB_ZvitHRDM(db *sql.DB) bool {
	zap := fmt.Sprintf("SELECT * FROM ZvitHRDM WHERE City = '%s' AND Dates = '%s'", a.City, a.Dates)

	rows, err := db.Query(zap)

	if err != nil {
		return false
	}
	n := 0
	for rows.Next() {
		n++
	}
	//fmt.Println(n, err, zap)
	if n == 0 {
		return false
	} else {
		return true
	}
}

// Чтение элемента ZvitHRDM из БД
func Read_DB_ZvitHRDM(db *sql.DB, date []time.Time) (res []ZvitHRDM) {

	zap := "SELECT * FROM ZvitHRDM WHERE "
	for n, t := range date {
		if n == 0 {
			zap = zap + fmt.Sprintf("Dates = '%s'", t.Format(TNS))
		} else {
			zap = zap + fmt.Sprintf(" OR Dates = '%s'", t.Format(TNS))
		}
	}
	// fmt.Println(zap)

	rows, err := db.Query(zap)
	if err != nil {
		fmt.Println(err, zap)
		return
	}

	for rows.Next() {
		p := make([]string, 56)
		err = rows.Scan(&p[0], &p[1], &p[2], &p[3], &p[4], &p[5], &p[6], &p[7], &p[8], &p[9], &p[10], &p[11], &p[12], &p[13], &p[14], &p[15], &p[16], &p[17], &p[18], &p[19], &p[20], &p[21], &p[22], &p[23], &p[24], &p[25], &p[26], &p[27], &p[28], &p[29], &p[30], &p[31], &p[32], &p[33], &p[34], &p[35], &p[36], &p[37], &p[38], &p[39], &p[40], &p[41], &p[42], &p[43], &p[44], &p[45], &p[46], &p[47], &p[48], &p[49], &p[50], &p[51], &p[52], &p[53], &p[54], &p[55])
		if err != nil {
			fmt.Println(err)
		}
		r := ZvitHRDM{}
		r.Reestablish(p)
		res = append(res, r)
	}
	return
}

// парсит []ZvitHRDM по городам
// отдает список сумм городов
func Parse_ZvitHRDM_Summa(a []ZvitHRDM) (b []ZvitHRDM) {
	cit := make([]string, 0)
	for _, i := range a {
		if !comp_string(i.City, cit) {
			cit = append(cit, i.City)
		}
	}
	for _, c := range cit {
		ln := 0
		r := ZvitHRDM{City: c}
		for _, i := range a {
			if c == i.City {
				r.Summ(i)
				ln++
			}
		}
		r.Div_sm(ln)
		b = append(b, r)
	}
	return
}

// парсит []ZvitHRDM по городам
// отдает список средних значений городов
func Parse_ZvitHRDM_Aver(a []ZvitHRDM) (b []ZvitHRDM) {
	cit := make([]string, 0)
	for _, i := range a {
		if !comp_string(i.City, cit) {
			cit = append(cit, i.City)
		}
	}
	for _, c := range cit {
		ln := 0
		r := ZvitHRDM{City: c}
		for _, i := range a {
			if c == i.City {
				r.Summ(i)
				ln++
			}
		}
		r.Div(ln)
		b = append(b, r)
	}
	return
}

// парсит []ZvitHRDM по городн
// отдает список дней города
func Parse_ZvitHRDM_City(a []ZvitHRDM, c string) (b []ZvitHRDM) {
	for _, i := range a {
		if c == i.City {
			b = append(b, i)
		}
	}
	return
}

// просчет общего значения по списку
func Total_ZvitHRDM(a []ZvitHRDM) (b ZvitHRDM) {
	b.City = "TOTAL"
	for _, i := range a {
		b.Summ(i)
	}
	return
}

// добавляет тотал в верх списка
func Result_ZvitHRDM(a []ZvitHRDM, b ZvitHRDM) (res []ZvitHRDM) {
	res = append(res, b)
	res = append(res, a...)
	return
}

// подготовка слайса строк для записи на лист
func Result_SliseString_ZvitHRDM(a []ZvitHRDM) (res [][]string) {
	res = append(res, zag_ZvitHRDM)
	for _, i := range a {
		res = append(res, i.Slise_String())
	}
	return
}

// подготовка слайса строк для записи на лист короткий вариант
func Result_SliseString_ZvitHRDM_Small(a []ZvitHRDM) (res [][]string) {
	res = append(res, zag_ZvitHRDM_Small)
	for _, i := range a {
		res = append(res, i.Slise_String_Small())
	}
	return
}

func Delete_ZvitHRDM(db *sql.DB, c, d string) {
	zap := fmt.Sprintf("DELETE FROM ZvitHRDM WHERE City = '%s' AND Dates = '%s'", c, d)
	fmt.Println(zap)
	db.Exec(zap)
}
