package main

import (
	"fmt"
	"sort"
	"strings"
	"time"
	"trs"
)

// Отчет за последние 7 дней
func NEW_REPORT_DM_NED() {

	OS(NS("Run ---  NEW_REPORT_DM_NED() ", 40))

	// if int(time.Now().Weekday()) != 1 {
	// 	//return
	// }

	// dned := int(time.Now().Weekday())

	res := make([][]string, 0)
	res_p := make([][]string, 0)
	day := last_days(7) //[:7]
	day_pr := last_days(41)
	// res = append(res, []string{fmt.Sprintf("Отчет обновлен %s", time.Now().Format("02.01.2006  15:04")), "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""})
	res = append(res, []string{"Місто", "Дата", "Тип відповідального", "Прізвище відповідального", "Корегуванняштрафу", "Сумма  штрафа до внесення", "Комментар ДМ", "Водитель", "Поездок всего", "Начал работать", "Закончил работать", "Штраф ДМ", "Штраф водителя", "Дней в работе", "Выходных", "Больничный", "Прогул", "Недоработка", day[0].Format(TNS), day[1].Format(TNS), day[2].Format(TNS), day[3].Format(TNS), day[4].Format(TNS), day[5].Format(TNS), day[6].Format(TNS)})
	res_p = append(res_p, []string{"Місто", "Дата", "Тип відповідального", "Прізвище відповідального", "Корегуванняштрафу", "Сумма  штрафа до внесення", "Комментар ДМ", "Водитель", "Поездок всего", "Начал работать", "Закончил работать", "Штраф ДМ", "Штраф водителя", "Дней в работе", "Выходных", "Больничный", "Прогул", "Недоработка", day[0].Format(TNS), day[1].Format(TNS), day[2].Format(TNS), day[3].Format(TNS), day[4].Format(TNS), day[5].Format(TNS), day[6].Format(TNS)})

	for _, sit := range create_sitys() {
		// данные графика за указанный период
		graf := get_grafik(sit, day)
		graf_pr := get_grafik(sit, day_pr)

		// поездки за указанные дни
		report := readBD_Trip(sit, day, "*")

		list_vod := make([]string, 0)

		for _, d := range report {
			for _, vd := range d.res_work {
				s1 := vd.name
				s2 := strings.Replace(s1, "(UK) ", "", 1)
				if !comp_string(s1, list_vod) && !strings.Contains(s1, "(UK) ") {
					list_vod = append(list_vod, s1)
				}
				if !comp_string(s2, list_vod) && !strings.Contains(s2, "(UK) ") {
					list_vod = append(list_vod, s2)
				}
			}
		}

		sort.Strings(list_vod)

		search_vod_day := func(n, d string) vod_day {
			var rs vod_day
			if len(report) == 0 {
				return rs
			}
			for _, dn := range report {
				if dn.data != d {
					continue
				}
				for _, vod := range dn.res_work {
					if vod.name == n {
						return vod
					}
				}
			}

			return rs
		}

		novvod := func(s, d string) bool {
			s1 := strings.Replace(s, "(UK) ", "", 1)
			if status[s] != "" {
				return false
			}
			t, _ := time.Parse(TNS, d)
			for _, i := range graf_pr {
				if t.Before(i.data) {
					continue
				}
				if i.name_vod1 == s || i.name_vod2 == s || strings.Replace(i.name_vod1, "(UK) ", "", 1) == s || strings.Replace(i.name_vod2, "(UK) ", "", 1) == s {
					return false
				}
				if i.name_vod1 == s1 || i.name_vod2 == s1 { // || strings.Replace(i.name_vod1, "(UK) ", "", 1) == s || strings.Replace(i.name_vod2, "(UK) ", "", 1) == s
					return false
				}
			}
			return true
		}

		oldvod := func(s string) string {

			dt := ""
			s1 := strings.Replace(s, "(UK) ", "", 1)
			if status[s] != "" {
				return dt
			}

			for _, i := range graf {

				if i.name_vod1 == s || i.name_vod2 == s || strings.Replace(i.name_vod1, "(UK) ", "", 1) == s || strings.Replace(i.name_vod2, "(UK) ", "", 1) == s {
					dt = i.data.Format(TNS)
				}
				if i.name_vod1 == s1 || i.name_vod2 == s1 { // || strings.Replace(i.name_vod1, "(UK) ", "", 1) == s || strings.Replace(i.name_vod2, "(UK) ", "", 1) == s
					dt = i.data.Format(TNS)
				}
			}
			if dt == day[len(day)-1].Format(TNS) { // .Add(-24*time.Hour)
				return ""
			} else {
				return dt
			}

		}

		type vodit struct {
			sity         string
			name         string
			car          []string
			trip         [7]int
			vyh_         [7]int
			lic_         [7]int
			prog_        [7]int
			nedor_       [7]int
			trav         [7]float64
			trips        int
			travals      float64
			hour         float64
			wr           int
			vyx          int
			lic          int
			norma        int
			prog         int
			nedor        int
			datnov       string
			datkon       string
			nov          int
			uv           int
			ocenka       string
			pain         int
			pain_manager int
		}

		v_str := func(a vodit) (res []string) {
			res = append(res, a.sity, tns(), "", "", "FALSE", "", "")
			res = append(res, a.name)
			res = append(res, fmt.Sprint(a.trips), a.datnov, a.datkon, fmt.Sprint(a.pain_manager), fmt.Sprint(a.pain), fmt.Sprint(a.wr))
			res = append(res, fmt.Sprint(a.vyx), fmt.Sprint(a.lic), fmt.Sprint(a.prog), fmt.Sprint(a.nedor))
			for i := 0; i < 7; i++ {
				if a.vyh_[i] == 0 && a.nedor_[i] == 0 && a.lic_[i] == 0 && a.prog_[i] == 0 {
					res = append(res, fmt.Sprint(a.trip[i]))
				}
				if a.trip[i] == 0 {
					if a.vyh_[i] > 0 {
						res = append(res, "Вых")
					} else if a.nedor_[i] > 0 {
						res = append(res, "Нед")
					} else if a.lic_[i] > 0 {
						res = append(res, "Бол")
					} else if a.prog_[i] > 0 {
						res = append(res, "Прг")
					}
				} else {
					if a.vyh_[i] > 0 {
						res = append(res, fmt.Sprint(a.trip[i])+" / "+"Вых")
					} else if a.nedor_[i] > 0 {
						res = append(res, fmt.Sprint(a.trip[i])+" / "+"Нед")
					} else if a.lic_[i] > 0 {
						res = append(res, fmt.Sprint(a.trip[i])+" / "+"Бол")
					} else if a.prog_[i] > 0 {
						res = append(res, fmt.Sprint(a.trip[i])+" / "+"Прг")
					}
				}
			}

			return
		}

		rs := make([]vodit, 0)
		for _, vd := range list_vod {

			r := vodit{name: vd}
			r.sity = sit.name
			for n, d := range day {
				num := 0

				cars := make([]string, 0)

				for _, g := range graf {

					// менеджеров не обрабатываем !!!
					if strings.Contains(g.name_vod1, "(M)") || strings.Contains(g.name_vod2, "(M)") {
						continue
					}

					if comp_string(g.car_nomer, cars) {
						continue
					}

					if d.Format(TNS) != g.data.Format(TNS) {
						continue
					}

					if strings.Contains(g.name_vod1, vd) || strings.Contains(g.name_vod2, vd) {
						num++
						var v1, v2, v3, v4 vod_day
						if if_name_voditel(g.name_vod1) && !if_name_voditel(g.name_vod2) {
							r.wr++
							if strings.Contains(g.name_vod1, "(UK)") {
								v1 = search_vod_day(g.name_vod1, day[n].Format(TNS))
								v2 = search_vod_day(strings.Replace(g.name_vod1, "(UK) ", "", 1), day[n].Format(TNS))
							} else {
								v1 = search_vod_day(g.name_vod1, day[n].Format(TNS))
								v2 = search_vod_day("(UK) "+g.name_vod1, day[n].Format(TNS))
							}
						} else if if_name_voditel(g.name_vod2) && !if_name_voditel(g.name_vod1) {
							r.wr++
							if strings.Contains(g.name_vod1, "(UK)") {
								v1 = search_vod_day(g.name_vod2, day[n].Format(TNS))
								v2 = search_vod_day(strings.Replace(g.name_vod2, "(UK) ", "", 1), day[n].Format(TNS))
							} else {
								v3 = search_vod_day(g.name_vod2, day[n].Format(TNS))
								v4 = search_vod_day("(UK) "+g.name_vod2, day[n].Format(TNS))
							}
						} else if if_name_voditel(g.name_vod2) && if_name_voditel(g.name_vod1) {
							r.wr++
							if strings.Contains(g.name_vod1, g.name_vod2) || strings.Contains(g.name_vod2, g.name_vod1) {
								v1 = search_vod_day(g.name_vod1, day[n].Format(TNS))
								v2 = search_vod_day(g.name_vod2, day[n].Format(TNS))
							} else {
								if strings.Contains(g.name_vod1, "(UK)") {
									v1 = search_vod_day(g.name_vod1, day[n].Format(TNS))
									v2 = search_vod_day(strings.Replace(g.name_vod1, "(UK) ", "", 1), day[n].Format(TNS))
								} else {
									v1 = search_vod_day(g.name_vod1, day[n].Format(TNS))
									v2 = search_vod_day("(UK) "+g.name_vod1, day[n].Format(TNS))
								}
								if strings.Contains(g.name_vod2, "(UK)") {
									v3 = search_vod_day(g.name_vod2, day[n].Format(TNS))
									v4 = search_vod_day(strings.Replace(g.name_vod2, "(UK) ", "", 1), day[n].Format(TNS))
								} else {
									v3 = search_vod_day(g.name_vod2, day[n].Format(TNS))
									v4 = search_vod_day("(UK) "+g.name_vod2, day[n].Format(TNS))
								}
							}
						}

						r.trips = r.trips + v1.trip_sum + v2.trip_sum + v3.trip_sum + v4.trip_sum
						r.trip[n] = v1.trip_sum + v2.trip_sum + v3.trip_sum + v4.trip_sum

						if novvod(g.name_vod1, g.data.Format(TNS)) || novvod(g.name_vod2, g.data.Format(TNS)) {
							r.nov++
							r.datnov = g.data.Format(TNS)
						}

						r.datkon = oldvod(vd)

						if status[g.name_vod1] == "Вихідний" || status[g.name_vod2] == "Вихідний" {
							r.vyx++
							r.vyh_[n]++
						} else if status[g.name_vod1] == "Лікарняний" || status[g.name_vod2] == "Лікарняний" {
							r.lic++
							r.lic_[n]++
						} else {
							if r.trip[n] == 0 {
								r.prog++
								r.prog_[n]++
							} else if r.trip[n] < 8 {
								r.nedor++
								r.nedor_[n]++
							} else if r.trip[n] > 19 {
								r.norma++
							}
						}
					}
				}
				if num > 1 {
					r.trips = r.trips / num
					r.trip[n] = r.trip[n] / num
					r.wr = r.wr / num
				}
				if r.wr > 7 {
					r.wr = 7
				}

				// модуль просчета штрафов
				if r.vyx > 2 {
					r.pain_manager = 100
				}
				if r.lic > 3 {
					r.pain_manager = 100
				}
				if r.lic > 3 && r.vyx > 1 {
					r.pain_manager = 100
				}
				if r.prog > 1 {
					r.pain_manager = 100
					r.pain = 350 * r.prog
				}
				if r.vyx == 2 && r.trips < 120 {
					r.pain_manager = 99
				}

			}
			rs = append(rs, r)
		}

		sort.Slice(rs, func(i, j int) (less bool) {
			return rs[i].trips > rs[j].trips
		})

		for _, i := range rs {
			res = append(res, v_str(i))
			if i.pain > 0 || i.pain_manager > 0 {
				res_p = append(res_p, v_str(i))
			}
		}
	}

	trs.Rec("1Y-E8coo3Td3BwysXQzIECO0FrXcsOHikhN26OpMk8tg", "'LAST_7_DAY'!A1", res, true, true)
	trs.Rec("1Y-E8coo3Td3BwysXQzIECO0FrXcsOHikhN26OpMk8tg", "'LAST_7_DAY_PAIN'!A1", res_p, true, true)

}

// отчет по перепробегам
func REPORT_OVER_RUN(days []time.Time, ignor_Rec bool) {

	OS(NS("Run ---  REPORT_OVER_RUN(dn, false) ", 40))

	list_top, err := trs.Read_sheets_CASHBOX_err("1Y-E8coo3Td3BwysXQzIECO0FrXcsOHikhN26OpMk8tg", "'TOP_DRIVER'!A3:G", 0)
	if err != nil {
		errorLog.Println(err)
	}
	TOP_DRIVER := make(map[string]int, 700)
	for _, i := range list_top {
		TOP_DRIVER[i[0]+"|"+trs.Clear_name_driver(i[1])] = string_to_int(i[2])
	}
	// for m, n := range TOP_DRIVER {
	// 	fmt.Println(n, m)
	// }

	// тариф за перепробег
	str_tar := float64(4)

	sitys := create_sitys() //[1]
	//day := last_days(1)

	day_pr := last_days(37)

	res := make([][]string, 0)
	res_Best := make([][]string, 0)
	res_Pokaz := make([][]string, 0)

	var report list_Trips
	report.s = get_Trips(sitys, days)

	for _, sit := range sitys {
		//fmt.Println(sit.name)

		for _, day := range days {

			if sit.name == "Варшава" {
				continue
			}

			// if sit.name != "Вінниця" {
			// 	continue
			// }

			// пробег по часам за указанные дни
			trip_h := trip_hour([]time.Time{day})

			// график за день
			graf := get_grafik(sit, []time.Time{day})
			graf_pr := get_grafik(sit, day_pr)

			novvod := func(s, d string) bool {
				s1 := strings.Replace(s, "(UK) ", "", 1)
				if status[s] != "" {
					return false
				}
				t, _ := time.Parse(TNS, d)
				for _, i := range graf_pr {
					if t.Before(i.data) {
						continue
					}
					if i.name_vod1 == s || i.name_vod2 == s || strings.Replace(i.name_vod1, "(UK) ", "", 1) == s || strings.Replace(i.name_vod2, "(UK) ", "", 1) == s {
						return false
					}
					if i.name_vod1 == s1 || i.name_vod2 == s1 { // || strings.Replace(i.name_vod1, "(UK) ", "", 1) == s || strings.Replace(i.name_vod2, "(UK) ", "", 1) == s
						return false
					}
				}
				return true
			}

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

			ck := make([]float64, 0)
			kt := make([]float64, 0)

			graf_sit := make([]grafik_over_run, 0)

			cars := make([]string, 0)

			for _, g := range graf {

				if comp_string(g.car_nomer, cars) {
					continue
				}

				// менеджеров не обрабатываем !!!
				if strings.Contains(g.name_vod1, "(M)") || strings.Contains(g.name_vod2, "(M)") {
					continue
				}

				cars = append(cars, g.car_nomer)

				r := grafik_over_run{}
				r.gr = g

				vd1 := trs.Clear_name_driver(g.name_vod1)
				vd2 := trs.Clear_name_driver(g.name_vod2)

				car_D := report.Calc(vd1, vd2, day.Format(TNS), sit.id_BGQ)

				_, trv := searh_car_ch(g.car_nomer)
				trip := car_D.Trip
				cash := car_D.Cash
				r.trav = trv
				r.trip = float64(trip)
				r.cash = cash
				r.travUgerDobavka = float64(car_D.Uber) * 5

				if trv < 5 {
					continue
				}

				if trip > 0 {
					cash_km := cash / trv
					km_trip := trv / float64(trip)

					r.cash_km = cash_km
					r.km_trip = km_trip

					if cash_km > 3 && cash_km < 15 {
						ck = append(ck, cash_km)
					}

					if km_trip < 20 && km_trip > 4 {
						kt = append(kt, km_trip)
					}
				}

				// находим ответственного
				v1, _, _ := g.Vidpovid()

				r.vidpovidalnyi = v1

				graf_sit = append(graf_sit, r)
			}

			sort.Float64s(ck)
			sort.Float64s(kt)

			// fmt.Println(len(ck), ck)
			// fmt.Println(len(kt), kt)

			best_ck := average_best(ck, true)
			best_kt := average_best(kt, false)

			for r := 0; r < len(graf_sit); r++ {
				graf_sit[r].koeff = modul(graf_sit[r].cash_km-best_ck, graf_sit[r].km_trip-best_kt)
			}

			kr := grafik_over_run{koeff: 1000}
			for _, i := range graf_sit {

				komm := ""
				rs := make([]string, 0)
				rs = append(rs, sit.name, day.Format(TNS), i.gr.car_nomer, i.vidpovidalnyi)
				rs = append(rs, float_to_string(best_ck), float_to_string(i.cash_km), float_to_string(best_kt), float_to_string(i.km_trip))
				rs = append(rs, float_to_string(i.trip), float_to_string(i.trav), float_to_string(i.cash), komm)
				res_Pokaz = append(res_Pokaz, rs)

				if i.vidpovidalnyi == "do not process" {
					continue
				}
				if i.trip < 20 {
					continue
				}
				if i.koeff < kr.koeff {
					kr = i
				}
			}

			// fmt.Println("*******", kr)
			// fmt.Println(best_ck, best_kt)

			if kr.trip > 19 {
				komm := "Типовий взірець гарного водія для Вашого міста на сьогодні"
				rs := make([]string, 0)
				rs = append(rs, sit.name, day.Format(TNS), kr.gr.car_nomer, kr.vidpovidalnyi)
				rs = append(rs, float_to_string(best_ck), float_to_string(kr.cash_km), float_to_string(best_kt), float_to_string(kr.km_trip))
				rs = append(rs, float_to_string(kr.trip), float_to_string(kr.trav), float_to_string(kr.cash), komm)
				res_Best = append(res_Best, rs)
			}

			for i := 0; i < len(graf_sit); i++ {

				trv := graf_sit[i].trav
				trip := graf_sit[i].trip
				cash := graf_sit[i].cash

				graf_sit[i].cash_km = best_ck
				graf_sit[i].km_trip = best_kt
				var overrun float64
				if trv < 5 {
					continue
				}

				if trip == 0 {
					graf_sit[i].overrun = trv
					continue
				}

				skidka := trip * (float64(30) / float64(20))

				skidkaUber := graf_sit[i].travUgerDobavka
				skidka = skidka + skidkaUber

				ideal1 := cash / best_ck
				ideal2 := trip * best_kt

				overrun1 := trv - ideal1
				overrun2 := trv - ideal2

				if overrun1 > overrun2 {
					overrun = overrun2 - skidka
				} else {
					overrun = overrun1 - skidka
				}

				if overrun < 12 {
					continue
				}

				graf_sit[i].overrun = overrun

			}

			sort.Slice(graf_sit, func(i, j int) (less bool) {
				return graf_sit[i].overrun > graf_sit[j].overrun
			})

			for _, i := range graf_sit {

				if i.overrun < 12 {
					continue
				}

				// находим ответственного
				v1, v2, ek := i.gr.Vidpovid()
				if v1 == "do not process" {
					continue
				}

				// исключаем тестовых водителей
				if test_driver[(sit.id_BGQ + "|" + v1)] {
					continue
				}

				typ_vid := "Водій"

				if v1 == "DM" {
					v1 = sit.DM
					typ_vid = "Manager"
				}
				if v1 == "Kvoliti" {
					v1 = sit.Kvoliti
					typ_vid = "Manager"
				}

				if novvod(v1, day.Format(TNS)) {
					continue
				}
				if novvod(v2, day.Format(TNS)) {
					continue
				}

				komment := ""

				_, top := TOP_DRIVER[sit.name+"|"+v1]

				if top {
					komment = "Топ водій" + ": "
					if i.trip < 21 {
						i.overrun = i.overrun - 30
					} else {
						i.overrun = i.overrun - (i.trip-20)*1.5
					}
				}

				if i.overrun < 12 {
					continue
				}

				vlpr := float64(0)

				vlpr = i.trip * 1.5

				if top {
					if i.trip < 21 {
						vlpr = 30
					} else {
						vlpr = i.trip * 1.5
					}
				}

				if !ek {
					if i.travUgerDobavka > 0 {
						komment = komment + fmt.Sprintf("Оплата за наднормовий пробіг у сумі %sгрн.,  на %s %s, за %s, авто %s, поіздок за добу %s, пробіг за добу %s, кеш за добу %s, власний пробіг %skm. враховано. На Uber додано %dкм.", fmt.Sprint(int(i.overrun*str_tar)), typ_vid, v1, day.Format(TNS), i.gr.car_nomer, float_to_string(i.trip), float_to_string(i.trav), float_to_string(i.cash), float_to_string(vlpr), int(i.travUgerDobavka))
					} else {
						komment = komment + fmt.Sprintf("Оплата за наднормовий пробіг у сумі %sгрн.,  на %s %s, за %s, авто %s, поіздок за добу %s, пробіг за добу %s, кеш за добу %s, власний пробіг %skm. враховано.", fmt.Sprint(int(i.overrun*str_tar)), typ_vid, v1, day.Format(TNS), i.gr.car_nomer, float_to_string(i.trip), float_to_string(i.trav), float_to_string(i.cash), float_to_string(vlpr))
					}
					rs := make([]string, 0)
					rs = append(rs, sit.name, day.Format(TNS), "", "", "FALSE", "", "", float_to_string(str_tar), i.gr.car_nomer, typ_vid, v1, i.gr.name_vod1, i.gr.name_vod2, float_to_string(i.overrun))
					rs = append(rs, fmt.Sprint(int(i.overrun*str_tar)), float_to_string(i.cash_km), float_to_string(i.km_trip))
					rs = append(rs, float_to_string(i.trip), float_to_string(i.trav), float_to_string(i.cash), komment)
					res = append(res, rs)
					//fmt.Println(rs)
				} else {
					komment = fmt.Sprintf("Eкіпаж: Оплата за наднормовий пробіг 50 відсотків у сумі %dгрн.,  на %s %s, за %s, авто %s, поіздок за добу %s, пробіг за добу %s, кеш за добу %s, власний пробіг враховано.", int(i.overrun*str_tar/2), typ_vid, v1, day.Format(TNS), i.gr.car_nomer, float_to_string(i.trip), float_to_string(i.trav), float_to_string(i.cash))
					rs := make([]string, 0)
					rs = append(rs, sit.name, day.Format(TNS), "", "", "FALSE", "", "", float_to_string(str_tar), i.gr.car_nomer, typ_vid, v1, i.gr.name_vod1, i.gr.name_vod2, float_to_string(i.overrun/2))
					rs = append(rs, fmt.Sprint(int(i.overrun*str_tar/2)), float_to_string(i.cash_km), float_to_string(i.km_trip))
					rs = append(rs, float_to_string(i.trip), float_to_string(i.trav), float_to_string(i.cash), komment)
					res = append(res, rs)
					//fmt.Println(rs)
					komment = fmt.Sprintf("Eкіпаж: Оплата за наднормовий пробіг 50 відсотків у сумі %dгрн.,  на %s %s, за %s, авто %s, поіздок за добу %s, пробіг за добу %s, кеш за добу %s, власний пробіг враховано.", int(i.overrun*str_tar/2), typ_vid, v2, day.Format(TNS), i.gr.car_nomer, float_to_string(i.trip), float_to_string(i.trav), float_to_string(i.cash))
					rs = make([]string, 0)
					rs = append(rs, sit.name, day.Format(TNS), "", "", "FALSE", "", "", float_to_string(str_tar), i.gr.car_nomer, typ_vid, v2, i.gr.name_vod1, i.gr.name_vod2, float_to_string(i.overrun/2))
					rs = append(rs, fmt.Sprint(int(i.overrun*str_tar/2)), float_to_string(i.cash_km), float_to_string(i.km_trip))
					rs = append(rs, float_to_string(i.trip), float_to_string(i.trav), float_to_string(i.cash), komment)
					res = append(res, rs)
					//fmt.Println(rs)
				}
			}
		}
	}

	// for n, i := range res {
	// 	fmt.Println(n, i[:20])
	// 	fmt.Println(i[20])
	// 	fmt.Println()
	// }

	if !ignor_Rec {
		list := "'OVER_RUN'!A2"
		trs.Rec("1Y-E8coo3Td3BwysXQzIECO0FrXcsOHikhN26OpMk8tg", list, res, true, true)
	}

	if !ignor_Rec {
		list := "'Typical Driver'!A2"
		trs.Rec("1Y-E8coo3Td3BwysXQzIECO0FrXcsOHikhN26OpMk8tg", list, res_Best, true, true)
	}

	if !ignor_Rec {
		list := "'Drv'!A2"
		trs.Rec("1boU2to5mT77THBlyIL1ZsCJLsqUKE393DcgJga3Z0_w", list, res_Pokaz, true, true)
	}

}

// отчет за вчера
func NEW_REPORT_DM_DAY(days []time.Time) {

	OS(NS("Run ---  NEW_REPORT_DM_DAY(dn) ", 40))

	res := make([][]string, 0)

	var report list_Trips
	report.s = get_Trips(create_sitys(), days)

	day_pr := last_days(37)
	// res = append(res, []string{fmt.Sprintf("Отчет обновлен %s", time.Now().Format("02.01.2006  15:04")), "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""})
	res = append(res, []string{"Місто", "Дата", "Тип відповідального", "Прізвище відповідального", "Корегуванняштрафу", "Сумма  штрафа до внесення", "Комментар ДМ", "Водитель", "Штраф", "Комментар DriverCSS", "AVTO"})

	for _, sit := range create_sitys() {

		// if sit.id_BGQ != "Lutsk" {
		// 	continue
		// }

		for _, day := range days {

			// данные графика за указанный период
			graf := get_grafik(sit, []time.Time{day})
			graf_pr := get_grafik(sit, day_pr)

			list_vod := make([]string, 0)

			for _, vd := range report.s {
				s1 := vd.Name
				if !comp_string(s1, list_vod) {
					list_vod = append(list_vod, s1)
				}
			}

			cars := make([]string, 0)
			for _, vd := range graf {

				if comp_string(vd.car_nomer, cars) {
					continue
				}

				vd.CorrecthionName()

				cars = append(cars, vd.car_nomer)
				s1 := trs.Clear_name_driver(vd.name_vod1)
				s2 := trs.Clear_name_driver(vd.name_vod2)

				if !comp_string(s1, list_vod) && status[s1] == "" && s1 != "" {
					list_vod = append(list_vod, s1)
				}
				if !comp_string(s2, list_vod) && status[s2] == "" && s2 != "" {
					list_vod = append(list_vod, s2)
				}
			}

			sort.Strings(list_vod)

			novvod := func(s, d string) bool {
				s1 := strings.Replace(s, "(UK) ", "", 1)
				if status[s] != "" {
					return false
				}
				t, _ := time.Parse(TNS, d)
				for _, i := range graf_pr {
					if t.Before(i.data) {
						continue
					}
					if i.name_vod1 == s || i.name_vod2 == s || strings.Replace(i.name_vod1, "(UK) ", "", 1) == s || strings.Replace(i.name_vod2, "(UK) ", "", 1) == s {
						return false
					}
					if i.name_vod1 == s1 || i.name_vod2 == s1 { // || strings.Replace(i.name_vod1, "(UK) ", "", 1) == s || strings.Replace(i.name_vod2, "(UK) ", "", 1) == s
						return false
					}
				}
				return true
			}

			type vodit struct {
				sity   string
				name   string
				car    string
				trip   int
				status string
				pain   int
			}

			v_str := func(a vodit) (res []string) {
				res = append(res, a.sity, day.Format(TNS), "", "", "FALSE", "", "")
				res = append(res, a.name)
				res = append(res, fmt.Sprint(a.pain), a.status, a.car)
				return
			}

			rs := make([]vodit, 0)

			for _, vd := range list_vod {

				if vd == "" {
					continue
				}

				r := vodit{name: vd}
				r.sity = sit.name

				for _, g := range graf {

					if strings.Contains(g.name_vod1, vd) || strings.Contains(g.name_vod2, vd) {
						r.car = g.car_nomer
						g.CorrecthionName()
						// менеджеров не обрабатываем !!!
						if strings.Contains(g.name_vod1, "(M)") || strings.Contains(g.name_vod2, "(M)") {
							continue
						}

						if status[g.name_vod1] == "ДТП" || status[g.name_vod2] == "ДТП" {
							continue
						}

						if status[g.name_vod1] == "Ремонт" || status[g.name_vod2] == "Ремонт" {
							continue
						}

						if status[g.name_vod1] == "Штрафмайданчик" || status[g.name_vod2] == "Штрафмайданчик" {
							continue
						}

						if status[g.name_vod1] == "Документи" || status[g.name_vod2] == "Документи" {
							continue
						}

						if status[g.name_vod1] == "Угон" || status[g.name_vod2] == "Угон" {
							continue
						}

						if status[g.name_vod1] == "Ключі" || status[g.name_vod2] == "Ключі" {
							continue
						}

						if status[g.name_vod1] == "Волонтер" || status[g.name_vod2] == "Волонтер" {
							continue
						}

						if status[g.name_vod1] == "Списання" || status[g.name_vod2] == "Списання" {
							continue
						}

						if status[g.name_vod1] == "Відновлення-авто" || status[g.name_vod2] == "Відновлення-авто" {
							continue
						}

						if g.name_vod1 == "" && g.name_vod2 == "" {
							continue
						}

						if novvod(g.name_vod1, g.data.Format(TNS)) || novvod(g.name_vod2, g.data.Format(TNS)) {
							continue
						}
						if status[g.name_vod1] == "Вихідний" || status[g.name_vod2] == "Вихідний" {
							continue
						} else if status[g.name_vod1] == "Лікарняний" || status[g.name_vod2] == "Лікарняний" {
							continue
						} else if status[g.name_vod1] == "Сервіс" || status[g.name_vod2] == "Сервіс" {
							continue
						} else if status[g.name_vod1] == "Оренда" || status[g.name_vod2] == "Оренда" {
							continue
						} else {

							vd1 := trs.Clear_name_driver(g.name_vod1)
							vd2 := trs.Clear_name_driver(g.name_vod2)
							trip := report.Calc(vd1, vd2, day.Format(TNS), sit.id_BGQ)
							r.trip = trip.Trip
							if strings.Contains(g.name_vod1, vd) {
								r.name = g.name_vod1
							} else {
								r.name = g.name_vod2
							}

							//fmt.Println(r.name, r.trip)
							// g.Print()
							// fmt.Println(r.trip)

							if r.trip == 0 {
								r.status = fmt.Sprintf("Штраф за прогул -  %s  -  %s  - поездок %d выходного, больничного, ТО-сервис в графике не проставлено", vd, day.Format(TNS), r.trip)
								r.pain = 350
								rs = append(rs, r)
							} else if r.trip < 8 {
								r.status = fmt.Sprintf("Штраф за недорабртку -  %s  -  %s  - поездок %d выходного, больничного, ТО-сервис в графике не проставлено", vd, day.Format(TNS), r.trip)
								r.pain = 350
								rs = append(rs, r)
							}
						}
					}
				}
			}
			for _, i := range rs {
				res = append(res, v_str(i))
			}
		}
	}

	// for _, i := range res {
	// 	fmt.Println(i[7], i)
	// }

	trs.Rec("1Y-E8coo3Td3BwysXQzIECO0FrXcsOHikhN26OpMk8tg", "'DM_DAY'!A1", res, true, true)

}

func corian() {
	// s := "SELECT * FROM Kharkiv.Kharkiv_shedule" // WHERE Date = '2022-12-22'
	// gr, err := BGQ.Read_BGQ(s)
	// fmt.Println(err)
	// for _, i := range gr {
	// 	fmt.Println(i)
	// }

	list, _ := trs.Read_sheets_CASHBOX_err("1r7-bcbpZmzWWTeHLm4WGJeYRyZxoDTSQVyFpDu3OYiM", "Заполнено", 1)
	type work struct {
		name     string
		city     string
		id       string
		accaunts []string
	}
	list_vod := make([]work, 0)
	for _, i := range list {
		if i[2] == "777" {
			continue
		}
		c := i[0]
		id := i[2]
		s := trs.Clear_name_driver(i[1])
		r := work{name: s, city: c, id: id}
		for _, j := range list {
			if j[0] == i[0] && j[2] == i[2] {
				if !trs.Compaire_String(trs.Clear_name_driver(j[1]), r.accaunts) {
					r.accaunts = append(r.accaunts, trs.Clear_name_driver(j[1]))
				}

			}
		}
		list_vod = append(list_vod, r)

	}

	trips := get_Trips(create_sitys(), period_days("13.12.2022", "13.01.2023"))
	res := make([][]string, 0)

	for _, cit := range create_sitys() {

		gr := grafiks{}
		gr.Get(cit, period_days("13.12.2022", "13.01.2023"))
		driver := gr.ListDriver()

		for _, v := range driver {
			r := Trips{Name: v, City: cit.name, Servise: "Total"}
			for _, tr := range trips {
				if v == tr.Name && tr.City == cit.id_BGQ {
					r.Add(tr)
				}
			}
			res = append(res, r.Report())
		}

	}

	trs.Rec("1QC3qRzY11G-LiufJdiVHNNyxlEg1TzZ83pu0z23UsU0", "'Report'!A2", res, true, true)
}
