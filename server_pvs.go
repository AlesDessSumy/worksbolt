package main

import (
	"fmt"
	"time"
	"trs"
)

type Select_Data struct {
	Name   string
	book   string
	sheet  string
	City   string
	dn     string
	dk     string
	sort   string
	z_sort string
	z1     string
	z2     string
	z3     string
	z4     string
	z5     string
	z6     string
	z7     string
	z8     string
	z9     string
	z10    string
	z11    string
	z12    string
	z13    string
}

// Печать элементаSelect_Data
func (a *Select_Data) PrintSm() {
	fmt.Println()
	fmt.Println("Select_Data         :")
	fmt.Println("Name                :", a.Name)
	fmt.Println("book                :", a.book)
	fmt.Println("sheet               :", a.sheet)
	fmt.Println("City                :", a.City)
	fmt.Println("dn                  :", a.dn)
	fmt.Println("dk                  :", a.dk)
	fmt.Println("sort                :", a.sort)
	fmt.Println("z_sort              :", a.z_sort)
}

// Печать элементаSelect_Data
func (a *Select_Data) Print() {
	fmt.Println()
	fmt.Println("Select_Data         :")
	fmt.Println("Name                :", a.Name)
	fmt.Println("book                :", a.book)
	fmt.Println("sheet               :", a.sheet)
	fmt.Println("City                :", a.City)
	fmt.Println("dn                  :", a.dn)
	fmt.Println("dk                  :", a.dk)
	fmt.Println("sort                :", a.sort)
	fmt.Println("z_sort              :", a.z_sort)
	fmt.Println("z1                  :", a.z1)
	fmt.Println("z2                  :", a.z2)
	fmt.Println("z3                  :", a.z3)
	fmt.Println("z4                  :", a.z4)
	fmt.Println("z5                  :", a.z5)
	fmt.Println("z6                  :", a.z6)
	fmt.Println("z7                  :", a.z7)
	fmt.Println("z8                  :", a.z8)
	fmt.Println("z9                  :", a.z9)
	fmt.Println("z10                 :", a.z10)
	fmt.Println("z11                 :", a.z11)
	fmt.Println("z12                 :", a.z12)
	fmt.Println("z13                 :", a.z13)
}

func (a *Select_Data) Create(b []string) {
	a.Name = b[0]
	a.book = b[1]
	a.sheet = b[2]
	a.City = b[3]
	a.dn = b[4]
	a.dk = b[5]
	a.sort = b[6]
	a.z_sort = b[7]
	a.z1 = b[8]
	a.z2 = b[9]
	a.z3 = b[10]
	a.z4 = b[11]
	a.z5 = b[12]
	a.z6 = b[13]
	a.z7 = b[14]
	a.z8 = b[15]
	a.z9 = b[16]
	a.z10 = b[17]
	a.z11 = b[18]
	a.z12 = b[19]
	a.z13 = b[20]
}

func server_pvs() {

	OS("Run ---  go server_pvs()")

	sl := make(map[string]Select_Data, 0)
	sl["Штрафы история"] = Select_Data{}
	sl["КПД Отчет Autoreport V4"] = Select_Data{}
	sl["История ДМ работа с водителями"] = Select_Data{}
	sl["Прибыль парка"] = Select_Data{}
	sl["Просмотр кассы"] = Select_Data{}
	sl["Просмотр кассы Фесенко"] = Select_Data{}
	sl["Контроль 5.1"] = Select_Data{}
	sl["Прибыль от бренда"] = Select_Data{}
	sl["ZvitHRDM"] = Select_Data{}
	sl["ZvitHRDM_small"] = Select_Data{}
	sl["Driver_Trips"] = Select_Data{}

	for {
		list, err := trs.Read_sheets_CASHBOX_err("1tNOveo3bMmdavTBQgK83wW6Xg0_Won50XXRdD6AtbCc", "'l1'!A1:U", 0)
		if err != nil {
			fmt.Println(err)
		}
		for n, i := range list {
			if n == 0 {
				continue
			}
			r := Select_Data{}
			r.Create(i)
			e := sl[r.Name]
			// if e.Name == "Просмотр кассы" {

			// }
			if r != e {
				sl[r.Name] = r
				// e.PrintSm()
				// r.PrintSm()
				Select_Run(r)
				OS(NS("Select_Run(r) ------  "+r.Name, 40))
			}
		}
		<-time.After(time.Second)
	}
}

func Select_Run(a Select_Data) {

	switch a.Name {
	case "Штрафы история":
		fmt.Println("RUN  *****      ", a)
	case "КПД Отчет Autoreport V4":
		fmt.Println("RUN  *****      ", a)
	case "История ДМ работа с водителями":
		fmt.Println("RUN  *****      ", a)
	case "Прибыль парка":
		fmt.Println("RUN  *****      ", a)
		go Report_Summa_Kassa_Profit_Period(a.dn, a.dk, a.sort)
	case "Просмотр кассы":
		fmt.Println("RUN  *****      ", a)
		go Kassa_Rashod_Period("1z6Tp6wvwIP3rWc8SK-Q9ZVWA_Xy_Bkm7w7d9N5q3Q1g", a.dn, a.dk, a.sort, a.City)
	case "Просмотр кассы Фесенко":
		fmt.Println("RUN  *****      ", a)
		go Kassa_Rashod_Period("1E_vX93lm76U7h6gUlovX6TuEZXZz5s4RsWgoQyDY5pI", a.dn, a.dk, a.sort, a.City)
	case "Контроль 5.1":
		fmt.Println("RUN  *****      ", a)
		go Check_kassa_51(a.dn, a.dk)
	case "Прибыль от бренда":
		fmt.Println("RUN  *****      ", a)
		go Report_Brend_Period(a.dn, a.dk, a.City)
	case "ZvitHRDM":
		fmt.Println("RUN  *****      ", a)
		go Zvit_HRDM_History(a.City, a.dn, a.dk)
	case "ZvitHRDM_small":
		fmt.Println("RUN  *****      ", a)
		go Zvit_HRDM_History_Small(a.City, a.dn, a.dk)
	case "Driver_Trips":
		fmt.Println("RUN  *****      ", a)
		if a.z13 == "Обрахувати" {
			go Search_Driver_Trips_Period(a.City, a.dn, a.dk)
		}
	}
}

// Прибыль от бренда
