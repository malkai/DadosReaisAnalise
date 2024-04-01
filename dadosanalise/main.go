package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type datahelp struct {
	name     string
	trajetos []string
	path_arq []string
}

type Str struct {
	data string
}

func readdata() {

	vehiclestt := []datahelp{}

	Arq1, err := os.ReadDir("dados")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	//Pega as datas
	for _, path1 := range Arq1 {

		//fmt.Println(ass)
		dir, err := os.ReadDir("dados/" + path1.Name())
		if err != nil {
			fmt.Println("Error:", err)

		}
		//Pega os nomes dos veiculos
		for _, path2 := range dir {
			//fmt.Println(path2)

			dir2, err := os.ReadDir("dados/" + path1.Name() + "/" + path2.Name())
			if err != nil {
				fmt.Println("Error:", err)

			}
			if !contains(vehiclestt, path2.Name()) {
				vehiclestt = append(vehiclestt, datahelp{name: path2.Name()})
			}
			//Pega os nomes dos trajetos
			for _, data2 := range dir2 {
				//fmt.Println(data2)

				indexx := getindex(vehiclestt, path2.Name())
				vehiclestt[indexx].trajetos = append(vehiclestt[indexx].trajetos, data2.Name())

				//fmt.Println(len(vehiclestt))

				dir3, err := os.ReadDir("dados/" + path1.Name() + "/" + path2.Name() + "/" + data2.Name())
				if err != nil {
					fmt.Println("Error:", err)

				}
				//Pega arquivos com os dados
				for _, data3 := range dir3 {
					//fmt.Println(data3)

					vehiclestt[indexx].path_arq = append(vehiclestt[indexx].path_arq, "dados/"+path1.Name()+"/"+path2.Name()+"/"+data2.Name()+"/"+data3.Name())

					//vehiclesfiles[path2.Name()] = append(vehiclespath[path2.Name()], Str{data: "dados/" + path1.Name() + "/" + path2.Name() + "/" + data2.Name() + "/" + data3.Name()})

				}
			}
		}

	}

	file, err := os.Create("dados.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	writer := csv.NewWriter(file)

	info := []string{}

	info = append(info, "Nome Veículo", "Nomes dos Trajeto", "Número de Tuplas", "Tempo do Trajeto (s)", "Frequência de dados (s)", "Desvio da Frequência de dados (s)", "Pontos utilizados Web ", "Distancia Web Percorrida (KM)", "Distancia Euclediano Percorrida (KM)", "Número de dados processados", "Media Combustivel (%)", "Desvio Combustivel (%)", "Filtro de Kalman (%)", "Sigma 1 (%)", "Sigma 2 (%)", "Sigma 3 (%)")
	writer.Write(info)

	writer.Flush()

	for _, n := range vehiclestt {
		//fmt.Println(i, n.name)

		info3 := [][]string{}

		for ii, m := range n.path_arq {
			if m != "dados/16_01_2024/Spin/Gara_Transporte_Restaurante/OBDLink.csv" && m != "dados/16_01_2024/Spin/Gara_Angra3/OBDLink.csv" && m != "dados/17_01_2024/Spin/Alojamento1-Transporte/OBDLink.csv" && m != "dados/17_01_2024/Spin/Transporte-Alojamento1/OBDLink.csv" && m != "dados/18_01_2024/Van/Alojamento3-Alojamento1/OBDLink_1.csv" {

				file2, _ := os.Create("comb/" + n.name + fmt.Sprintf("%d", ii+1) + "dadoscomb.csv")

				info = nil

				info = append(info, "Tempo0", "Comb0", "Kalman0", "Tempo1", "Comb1", "Kalman1", "Tempo2", "Comb2", "Kalman2", "Tempo3", "Comb3", "Kalman3")

				writer2 := csv.NewWriter(file2)
				writer2.Write(info)
				info = nil

				listtime := []float64{}
				listtime2 := []float64{}
				listdistance := []string{}

				//listdistanweb := []string{}
				listfuel := []float64{}
				listfuel2 := []float64{}
				info := []string{}

				lat := []string{}
				long := []string{}
				//latlong := [][]string{}

				// '\n'
				boolread := false
				listnames := []int{}

				info = append(info, n.name)

				info = append(info, m)

				file, err := os.Open(m)
				if err != nil {
					log.Fatal(err)
				}
				defer file.Close()

				//analise1
				//meta = (rand.Intn(4000-1000) + 1000)

				scanner := bufio.NewScanner(file)

				removepetido3 := ""

				for scanner.Scan() {

					if !boolread {
						res1 := strings.Split(scanner.Text(), ",")
						if res1[0] == "Time (sec)" {
							//Time (sec), Latitude (deg), Longitude (deg), Fuel level input (%)  Instant fuel economy (l/100 km)
							for oo, nn := range res1 {

								//fmt.Println(len(nn), nn)

								if nn == "Time (sec)" || nn == " Latitude (deg)" || nn == " Longitude (deg)" || nn == " Fuel level input (%)" {

									//fmt.Println(nn)
									listnames = append(listnames, oo)
								}
							}

							boolread = true
						}
					} else {
						res1 := strings.Split(scanner.Text(), ",")

						//listdistanweb

						stringlat := res1[listnames[1]]
						stringlong := res1[listnames[2]]

						latweb, err := strconv.ParseFloat(stringlat, 64)
						if err != nil {
							log.Fatal(err)
						}

						longweb, err := strconv.ParseFloat(stringlong, 64)
						if err != nil {
							log.Fatal(err)
						}

						if len(listnames) == 4 {

							help, err := strconv.ParseFloat(res1[listnames[0]], 64)
							if err != nil {
								log.Fatal(err)
							}

							listtime = append(listtime, help)

							s1 := fmt.Sprintf("%f", latweb)

							s2 := fmt.Sprintf("%f", longweb)

							listdistance = append(listdistance, res1[listnames[2]]+","+res1[listnames[1]])
							//listdistanweb = append(listdistanweb, s2+","+s1)

							lat = append(lat, s1)
							long = append(long, s2)

							help, err = strconv.ParseFloat(res1[listnames[3]], 64)
							if err != nil {
								log.Fatal(err)
							}

							listfuel = append(listfuel, help)

						}
						if len(listnames) == 4 && res1[listnames[3]] != "0" && removepetido3 != res1[listnames[3]] {
							removepetido3 = res1[listnames[3]]

							//fmt.Println("oi")

							help, err := strconv.ParseFloat(res1[listnames[3]], 64)
							if err != nil {
								log.Fatal(err)
							}

							listfuel2 = append(listfuel2, help)

							help, err = strconv.ParseFloat(res1[listnames[0]], 64)
							if err != nil {
								log.Fatal(err)
							}

							listtime2 = append(listtime2, help)

						}

						//res1 := strings.Split(scanner.Text(), ",")
						//fmt.Println(res1)
					}

					//fmt.Println(scanner.Text())
				}

				//acctime := 0.0

				desviolist := []float64{}

				for ii := range listtime {

					if ii < len(listtime)-1 {

						//acctime = acctime + listtime[ii+1] - listtime[ii]
						desviolist = append(desviolist, listtime[ii+1]-listtime[ii])

					}

				}

				s := "0"

				s = fmt.Sprintf("%d", len(listtime))
				info = append(info, s)
				if len(listtime)-1 > 0 {
					s = fmt.Sprintf("%f", listtime[len(listtime)-1])
				} else {
					s = fmt.Sprintf("%f", 0)
				}
				//fmt.Println(listtime[len(listtime)-1])
				info = append(info, s)

				acctimemedia := mediavector(desviolist)
				s = fmt.Sprintf("%f", acctimemedia)
				info = append(info, s)
				acctdesvio := DesvioPadrão(desviolist, acctimemedia)
				s = fmt.Sprintf("%f", acctdesvio)
				info = append(info, s)

				dist := 0.0
				dist_web := 0.0

				for index := range listdistance {

					if index != len(listdistance)-1 {

						p1 := listdistance[index]
						p2 := listdistance[index+1]

						p3, err := Distanceeucle(p1, p2)
						if err != nil {
							log.Fatal(err)
						}
						if p3 < 0 {

							p3 = p3 * -1
						}
						dist = dist + p3

					}

				}

				dist_web = 0

				tampontos := 0

				if dist_web < dist {

					if len(listdistance) == 881 {

						/*
							file3, _ := os.Create("gpscoordinterno/" + n.name + "Interno" + fmt.Sprintf("%d", ii+1) + "dadostuplas.csv")
							file4, _ := os.Create("gpscoord/" + n.name + "openstree" + fmt.Sprintf("%d", ii+1) + "dadostuplas.csv")
							file5, _ := os.Create("dadosdistance/" + n.name + fmt.Sprintf("%d", ii+1) + "dadostuplas.csv")

							writer3 := csv.NewWriter(file3)
							writer4 := csv.NewWriter(file4)
							write5 := csv.NewWriter(file5)


							tampontos, dist, dist_web = escreve(dist, dist_web, ii, tampontos, n.name, s, *writer3, *writer4, *write5, *write6, listdistance, lat, long, latlong, len(listdistance))
						*/
						file6, _ := os.Create("analise/" + n.name + fmt.Sprintf("%d", ii+1) + "dadostuplamapa.csv")
						write6 := csv.NewWriter(file6)
						escreveTudo(*write6, listdistance, lat, long, len(listdistance))
					}
					/*
						else {
							file3, _ := os.Create("gpscoordinterno/" + n.name + "Interno" + fmt.Sprintf("%d", ii+1) + "dadostuplas.csv")
							file4, _ := os.Create("gpscoord/" + n.name + "openstree" + fmt.Sprintf("%d", ii+1) + "dadostuplas.csv")
							file5, _ := os.Create("dadosdistance/" + n.name + fmt.Sprintf("%d", ii+1) + "dadostuplas.csv")

							writer3 := csv.NewWriter(file3)
							writer4 := csv.NewWriter(file4)
							write5 := csv.NewWriter(file5)

							file6, _ := os.Create("dadosgrande/" + n.name + fmt.Sprintf("%d", ii+1) + "dadostuplamapa.csv")
							write6 := csv.NewWriter(file6)
							tampontos, dist, dist_web = escreve(dist, dist_web, ii, tampontos, n.name, s, *writer3, *writer4, *write5, *write6, listdistance, lat, long, latlong, 65)

						}*/
				}

				s = fmt.Sprintf("%d", tampontos)
				info = append(info, s)

				s = fmt.Sprintf("%f", dist_web)
				info = append(info, s)

				s = fmt.Sprintf("%f", dist)
				info = append(info, s)
				/*
					s = fmt.Sprintf("%f", mediavector(listdistancem))
					info = append(info, s)
					s = fmt.Sprintf("%f", DesvioPadrão(listdistancem, mediavector(listdistancem)))
					info = append(info, s)
				*/

				s = fmt.Sprintf("%d", len(listfuel2))
				info = append(info, s)

				s = fmt.Sprintf("%f", mediavector(listfuel2))
				info = append(info, s)
				s = fmt.Sprintf("%f", DesvioPadrão(listfuel2, mediavector(listfuel2)))
				info = append(info, s)

				info, info3 = processadados(info, info3, listfuel2, listtime2, 0, n.name)
				info, info3 = processadados(info, info3, listfuel2, listtime2, 1, n.name)
				info, info3 = processadados(info, info3, listfuel2, listtime2, 2, n.name)
				info, info3 = processadados(info, info3, listfuel2, listtime2, 3, n.name)

				info3 = transpose(info3)
				writer2.WriteAll(info3)
				writer2.Flush()
				info3 = [][]string{}

				writer.Write(info)
				writer.Flush()

			}
		}

	}

}

func escreve(dist float64, dist_web float64, ii, tampontos int, n, s string, writer3, writer4, writer5, writer6 csv.Writer, listdistance []string, lat []string, long []string, latlong [][]string, valueexp int) (a int, b float64, c float64) {

	posf := []string{"0", fmt.Sprint(len(listdistance) - 1)}

	res2 := strings.Split(listdistance[0], ",")
	res3 := strings.Split(listdistance[len(listdistance)-1], ",")

	lato := []string{res2[0], res3[0]}
	longo := []string{res2[1], res3[1]}

	latlong = append(latlong, posf)
	latlong = append(latlong, lato)
	latlong = append(latlong, longo)
	latlong = append(latlong, lat)
	latlong = append(latlong, long)
	latlong = transpose2(latlong)

	writer3.WriteAll(latlong)

	latlong = nil
	lat = nil
	long = nil
	posf = nil
	lato = nil
	longo = nil

	listindex := []int{}

	//ghelp := 2

	//div := int(((len(listdistanweb) - 1) * ghelp) / 100)
	loc := 0
	ghelp := 1

	divstatus := (len(listdistance) - 1) / ghelp

	//fmt.Println(ghelp, divstatus, len(listdistanweb))

	lataux := []string{}
	longaux := []string{}

	guh := []string{}

	guh = append(guh, n+fmt.Sprintf("%d", ii))
	writer5.Write(guh)
	writer5.Flush()
	guh = []string{}

	guh = append(guh, "Numero de pontos", "Expoencial", "Distancia Web KM", "Distancia euclidiana KM")
	writer5.Write(guh)
	writer5.Flush()

	listindex = append(listindex, 0)
	// dist_web < dist &&
	for ghelp <= valueexp {

		posf = nil
		lato = nil
		longo = nil
		dist_web2 := 0.0

		for i := 1; i <= ghelp; i++ {

			if divstatus*i < len(listdistance)-1 {

				listindex = append(listindex, divstatus*i)

			}

		}

		listindex = append(listindex, len(listdistance)-1)

		for _, jk := range listindex {
			posf = append(posf, strconv.Itoa(jk))
		}

		//fmt.Println(listindex)

		//fmt.Println(listindex[len(listindex)-1], len(listdistance)-1)

		//	fmt.Println(listindex)

		tupleList := []string{}

		for i, m := range listindex {

			//&& m+1 < len(listdistanweb)

			//if i < len(listindex)-1 {

			//fmt.Println(m)

			tupleList = append(tupleList, listdistance[m])

			// len(tupleList) == len(listindex)-1 ||

			if len(tupleList) >= 100 || i == len(listindex)-1 {

				//fmt.Println(len(tupleList) - 1)

				latlong := ""
				mk := "https://map.project-osrm.org/?z=10&center=-23.005241%2C-44.130020"

				mklist := []string{}

				for i, nk := range tupleList {
					res1 := strings.Split(nk, ",")

					lato = append(lato, res1[0])

					longo = append(longo, res1[1])

					if i == 0 {
						latlong = latlong + nk
					} else {
						latlong = latlong + ";" + nk
					}

					//fmt.Println(latlong)

				}

				mk = mk + "&hl=en&alt=0&srv=0"
				mklist = append(mklist, mk)
				writer6.Write(mklist)
				writer6.Flush()
				//-22.93069,-43.97846;-22.98799,-44.23648;-22.91494,-44.34997;-23.00783,-44.48557
				//fmt.Println("curl", "http://router.project-osrm.org/route/v1/driving/"+latlong+"?overview=false")

				curl := exec.Command("curl", "http://router.project-osrm.org/route/v1/driving/"+latlong+"?overview=false")

				out, err := curl.Output()
				if err != nil {

					fmt.Println("erorr", err)

				}

				auxstring := string(out)

				//fmt.Println(auxstring)

				//fmt.Println(auxstring)

				distcalculate := 0.0
				distcalculate, lataux1, longaux1 := getdistance(auxstring, tupleList)

				for i, _ := range lataux1 {

					lataux = append(lataux, lataux1[i])

					longaux = append(longaux, longaux1[i])

				}
				tupleList = []string{}

				dist_web2 = dist_web2 + distcalculate

			}

			//auxstring := string(out)

			//}

		}

		lat = lataux
		long = longaux

		fmt.Println("tam net - ", dist_web2, " tam euller - ", dist)

		//ghelp = ghelp + 5

		//div = int(((len(listdistanweb) - 1) * ghelp) / 100)
		guh = []string{}

		guh = append(guh, fmt.Sprintf("%d", len(listindex)-2))

		s = fmt.Sprintf("%d", loc)
		guh = append(guh, s)

		loc = loc + 1

		ghelp = int(math.Pow(2, float64(loc))) + 1

		divstatus = (len(listdistance) - 1) / ghelp

		fmt.Println(ghelp, divstatus)

		s = fmt.Sprintf("%f", dist_web2)
		s = strings.Replace(s, ".", ",", -1)
		guh = append(guh, s)

		s = fmt.Sprintf("%f", dist)
		s = strings.Replace(s, ".", ",", -1)
		guh = append(guh, s)

		guh2 := []string{"Novo ciclo"}

		latlong = append(latlong, posf)
		latlong = append(latlong, lato)
		latlong = append(latlong, longo)
		latlong = append(latlong, longaux)
		latlong = append(latlong, lataux)

		latlong = transpose2(latlong)

		writer4.Write(guh2)
		writer4.Flush()
		writer4.WriteAll(latlong)
		latlong = [][]string{}

		lataux = nil
		longaux = nil
		posf = nil
		lato = nil
		longo = nil

		writer5.Write(guh)
		writer5.Flush()

		tampontos = len(listindex)
		dist_web = dist_web2

		listindex = []int{}
		listindex = append(listindex, 0)

	}
	latlong = append(latlong, posf)
	latlong = append(latlong, lato)
	latlong = append(latlong, longo)
	latlong = append(latlong, long)
	latlong = append(latlong, lat)

	latlong = transpose2(latlong)

	writer4.WriteAll(latlong)

	return tampontos, dist, dist_web

}

func escreveTudo(writer6 csv.Writer, listdistance []string, lat, long []string, valueexp int) {

	listindex := []int{}

	//ghelp := 2

	//div := int(((len(listdistanweb) - 1) * ghelp) / 100)

	ghelp := valueexp - 1

	divstatus := (len(listdistance) - 1) / ghelp

	fmt.Println(divstatus)

	//fmt.Println(ghelp, divstatus, len(listdistanweb))

	lataux := []string{}
	longaux := []string{}

	listindex = append(listindex, 0)
	// dist_web < dist &&

	dist_web2 := 0.0

	for i := 1; i <= ghelp; i++ {

		if divstatus*i < len(listdistance)-1 {

			listindex = append(listindex, divstatus*i)

		}

	}

	listindex = append(listindex, len(listdistance)-1)

	//fmt.Println(listindex)

	//fmt.Println(listindex[len(listindex)-1], len(listdistance)-1)

	//	fmt.Println(listindex)

	tupleList := []string{}

	for i, m := range listindex {

		//&& m+1 < len(listdistanweb)

		//if i < len(listindex)-1 {

		//fmt.Println(m)

		tupleList = append(tupleList, listdistance[m])

		// len(tupleList) == len(listindex)-1 ||

		if len(tupleList) == 65 || i == len(listindex)-1 {

			//fmt.Println(len(tupleList) - 1)

			latlong := ""
			mk := "https://map.project-osrm.org/?z=10&center=-23.005241%2C-44.130020"

			mklist := []string{}
			r := "&radiuses="
			for i, nk := range tupleList {

				res1 := strings.Split(nk, ",")

				mk2 := "&loc=" + res1[1]
				mk3 := "%2C" + res1[0]
				mk = mk + mk2 + mk3

				if i == 0 {
					latlong = latlong + nk
					r = r + "11"
				} else {
					latlong = latlong + ";" + nk
					r = r + ";11"
				}

				//fmt.Println(latlong)

			}

			fmt.Println(r)

			mk = mk + "&hl=en&alt=0&srv=0"
			mklist = append(mklist, mk)
			writer6.Write(mklist)
			writer6.Flush()
			//-22.93069,-43.97846;-22.98799,-44.23648;-22.91494,-44.34997;-23.00783,-44.48557

			//fmt.Println("curl", "http://router.project-osrm.org/route/v1/driving/"+latlong+"?overview=false")

			//curl 'http://router.project-osrm.org/route/v1/driving/13.388860,52.517037;13.397634,52.529407;13.428555,52.523219?overview=simplified&generate_hints=false&radiuses=49;49;49

			curl := exec.Command("curl", "http://router.project-osrm.org/route/v1/driving/"+latlong+"?overview=simplified"+r)

			out, err := curl.Output()
			if err != nil {

				fmt.Println("erorr", err)

			}

			auxstring := string(out)

			fmt.Println(auxstring)

			//fmt.Println(auxstring)

			distcalculate := 0.0
			distcalculate, lataux1, longaux1 := getdistance(auxstring, tupleList)

			for i, _ := range lataux1 {

				lataux = append(lataux, lataux1[i])

				longaux = append(longaux, longaux1[i])

			}
			tupleList = []string{}

			dist_web2 = dist_web2 + distcalculate

		}

		//auxstring := string(out)

		//}

	}
	fmt.Println(dist_web2)

	//ghelp = ghelp + 5

	//div = int(((len(listdistanweb) - 1) * ghelp) / 100)

}

func getdistance(a string, tupleList []string) (float64, []string, []string) {
	lataux, longaux := []string{}, []string{}
	r, _ := regexp.Compile("\\w+")

	auxstringl := r.FindAllString(a, -1)

	valuefloat, valuefloatdot, valuesum := 0.0, 0.0, 0.0

	o2 := len(tupleList)

	for io, m := range auxstringl {

		if m == "distance" {

			o2 = o2 - 1
			//[o2 == 0 && ]
		}
		if o2 == 0 && m == "distance" {

			//fmt.Println(m, auxstringl[io+1], auxstringl[io+2])
			valuefloat, _ = strconv.ParseFloat(auxstringl[io+1], 64)

			if auxstringl[io+2] != "steps" {
				valuefloatdot, _ = strconv.ParseFloat(auxstringl[io+2][0:1], 64)
			} else {
				valuefloatdot = 0.0
			}

			valuefloatdot = valuefloatdot / math.Pow(10, float64(len(auxstringl[io+2][0:1])))

			valuefloat = valuefloat + valuefloatdot

			valuesum = valuesum + (valuefloat)/1000

			//fmt.Println(auxstringl[io+1], auxstringl[io+2], valuesum)

		}

		if m == "location" {
			a = tupleList[0][0:1]
			if a == "-" {

				lataux = append(lataux, "-"+auxstringl[io+1]+"."+auxstringl[io+2])
				longaux = append(longaux, "-"+auxstringl[io+3]+"."+auxstringl[io+4])
			} else {
				lataux = append(lataux, auxstringl[io+1]+"."+auxstringl[io+2])
				longaux = append(longaux, auxstringl[io+3]+"."+auxstringl[io+4])

			}
		}

	}

	return valuesum, lataux, longaux
}

func processadados(info []string, info3 [][]string, listfuel []float64, listtime []float64, t float64, n string) ([]string, [][]string) {

	infocomb, infotime, sigma1 := removevalores(listfuel, listtime, t)
	info3 = append(info3, infotime)
	info3 = append(info3, infocomb)

	fuelest := 0.0
	inforkalman := []float64{}
	if n == "Van" {
		fuelest, inforkalman, _ = KalmanFilter(90.00, sigma1)
	} else {
		fuelest, inforkalman, _ = KalmanFilter(51.00, sigma1)
	}

	inforkalmanS := convertToString(inforkalman)

	info3 = append(info3, inforkalmanS)
	s := fmt.Sprintf("%f", fuelest)
	info = append(info, s)

	return info, info3

}

func removevalores(listfuel []float64, infotime []float64, t float64) ([]string, []string, []float64) {
	sigmacomb := []float64{}
	sigmatime := []float64{}

	max := mediavector(listfuel) + (DesvioPadrão(listfuel, mediavector(listfuel)) * t)
	min := mediavector(listfuel) - (DesvioPadrão(listfuel, mediavector(listfuel)) * t)

	for i := range listfuel {
		if t == 0 {
			sigmacomb = append(sigmacomb, listfuel[i])
			sigmatime = append(sigmatime, infotime[i])
		} else if listfuel[i] >= min && listfuel[i] <= max {
			sigmacomb = append(sigmacomb, listfuel[i])
			sigmatime = append(sigmatime, infotime[i])
		}
	}

	sigmacombs := convertToString(sigmacomb)
	sigmatimes := convertToString(sigmatime)

	return sigmacombs, sigmatimes, sigmacomb
}

func convertToString(a []float64) []string {
	b := []string{}
	for i := range a {
		s := fmt.Sprintf("%f", a[i])

		b = append(b, s)
	}

	return b

}

func transpose2(slice [][]string) [][]string {

	newMatrix := [][]string{}

	newValue := []string{}
	//3798 3798[3] 3679[2]  2497[1]
	//fmt.Println(slice[0])

	for i := range slice[0] {

		//fmt.Println(slice[0][i], slice[1][i], slice[2][i], slice[3][i], slice[4][i])

		newValue = append(newValue, slice[0][i], slice[2][i], slice[1][i], slice[3][i], slice[4][i])

		newMatrix = append(newMatrix, newValue)
		newValue = nil
		newValue = []string{}

	}

	return newMatrix
}

func transpose(slice [][]string) [][]string {

	newMatrix := [][]string{}

	newValue := []string{}
	//3798 3798[3] 3679[2]  2497[1]

	for i := range slice[0] {

		if i < len(slice[4])-1 {

			newValue = append(newValue, slice[0][i], slice[1][i], slice[2][i], slice[9][i], slice[10][i], slice[11][i], slice[6][i], slice[7][i], slice[8][i], slice[3][i], slice[4][i], slice[5][i])

		} else if i < len(slice[7])-1 {

			newValue = append(newValue, slice[0][i], slice[1][i], slice[2][i], slice[9][i], slice[10][i], slice[11][i], slice[6][i], slice[7][i], slice[8][i], "0", "0", "0")

		} else if i < len(slice[10])-1 {

			newValue = append(newValue, slice[0][i], slice[1][i], slice[2][i], slice[9][i], slice[10][i], slice[11][i], "0", "0", "0", "0", "0", "0")

		} else if i <= len(slice[1])-1 {

			newValue = append(newValue, slice[0][i], slice[1][i], slice[2][i], "0", "0", "0", "0", "0", "0", "0", "0", "0")

		}

		newMatrix = append(newMatrix, newValue)
		newValue = nil
		newValue = []string{}

	}

	return newMatrix
}

func getindex(s []datahelp, str string) int {
	aux := 0
	for i, v := range s {
		if v.name == str {
			aux = i
		}
	}

	return aux

}

func contains(s []datahelp, str string) bool {
	for _, v := range s {
		if v.name == str {
			return true
		}
	}

	return false
}

func main() {

	readdata()

}
