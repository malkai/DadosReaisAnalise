package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

const R = 6371 //raio da Terra em km

func mediavector(a []float64) float64 {
	media := 0.0
	for _, leitura := range a {
		media = media + leitura
	}

	media = media / float64(len(a)-1)

	return media

}

func DesvioPadrão(a []float64, media float64) float64 {
	var aux float64 = 0

	for _, leitura := range a {
		aux = aux + (math.Pow(leitura-media, 2))

	}
	aux = math.Sqrt((aux) / float64(len(a)))
	//fmt.Println(aux)
	return aux
}

func KalmanFilter(capacidade float64, medições []float64) (float64, []float64, error) {
	if len(medições) > 0 {
		media := mediavector(medições)

		if DesvioPadrão(medições, media) != 0 {

			Gerro := errovector(medições, math.Round(media*10000)/10000) // desvio global

			auxe := []float64{}

			auxe = append(auxe, medições[1])
			Lerro := errovector(auxe, math.Round(media*10000)/10000) // desvio local
			//fmt.Println(Lerro, auxe, math.Ceil(media*100)/100, math.Ceil(math.Pow(Lerro, 2)*100)/100, math.Ceil(math.Pow(Gerro, 2)*10000)/10000)
			k := (math.Ceil(math.Pow(Lerro, 2)*100000) / 100000) / (math.Ceil(math.Pow(Lerro, 2)*100000)/100000 + math.Ceil(math.Pow(Gerro, 2)*100000)/100000)

			var estima float64
			leiturasPercentuaispos := []float64{}

			for _, leitura := range medições {
				//k := (math.Round(math.Pow(Lerro, 2)*1000) / 1000) / ((math.Round(math.Pow(Lerro, 2)*1000) / 1000) + (math.Round(math.Pow(Gerro, 2)*1000) / 1000))

				estima = media + k*((math.Ceil(leitura*100000)/100000)-(math.Ceil(media*100000)/100000))

				if estima > 100 {
					estima = 100
				} else if estima < 0 {
					estima = 0
				}
				//fmt.Println(leitura, estima)

				Lerro = (1.0 - k) * (math.Ceil(Lerro*100000) / 100000)

				media = estima

				k = (math.Ceil(Lerro*100000) / 100000) / (math.Ceil(Lerro*100000)/100000 + math.Ceil(math.Pow(Gerro, 2)*100000)/100000)

				//fmt.Println("Kalman gain", k, "estima", estima, "erro local", Lerro, "erro glbal", Gerro)
				leiturasPercentuaispos = append(leiturasPercentuaispos, math.Round(estima*1000)/1000)

			}

			//fmt.Println(leiturasPercentuaispos[0], leiturasPercentuaispos[len(leiturasPercentuaispos)-1], len(leiturasPercentuaispos))

			//resultadotanque := ((medições[0] - estima) * capacidade) / 100
			//fmt.Println(media, leiturasPercentuaispos[0], medições[0], leiturasPercentuaispos[len(leiturasPercentuaispos)-3])
			//media2 := mediavector(leiturasPercentuaispos)
			return ((leiturasPercentuaispos[0] - leiturasPercentuaispos[len(leiturasPercentuaispos)-1]) * 100) / capacidade, leiturasPercentuaispos, nil

		} else {
			return medições[0] - medições[len(medições)-1], medições, nil

		}
	}
	return 0, medições, nil

}

// https://physics.stackexchange.com/questions/704367/how-to-quantify-the-uncertainty-of-the-time-series-average
func errovector(a []float64, media float64) float64 {
	var aux float64 = 0
	errovec := []float64{}
	for _, leitura := range a {
		errovec = append(errovec, math.Pow(leitura-media, 2))

	}
	for _, leitura := range errovec {
		aux = +leitura

	}
	if len(a) == 1 {
		aux = aux / float64(len(a))
	} else {
		aux = aux / float64(len(a)-1)
	}

	return aux
}

func CoordenadasCartesianas(latitude, longitude float64) []float64 {
	theta := latitude * (math.Pi / 180.0) //teta
	phi := longitude * (math.Pi / 180.0)  //fi

	x := R * math.Cos(theta) * math.Cos(phi)
	y := R * math.Cos(theta) * math.Sin(phi)
	z := R * math.Sin(theta)
	var a []float64
	a = append(a, x)
	a = append(a, y)
	a = append(a, z)
	return a
}

func Distanceeucle(latlongA, latlongB string) (float64, error) {

	distancia := 0.0
	res1 := strings.Split(latlongA, ",")

	latitudeA, err := strconv.ParseFloat(res1[1], 64)
	if err != nil {
		return 0.0, fmt.Errorf("\n Erro checar LatA. %v", err)
	}
	longitudeA, _ := strconv.ParseFloat(res1[0], 64)
	if err != nil {
		return 0.0, fmt.Errorf("\n Erro checar LonA. %v", err)
	}

	//return 0.0, fmt.Errorf("\n %f %f ", latitudeA, longitudeA)

	res1 = strings.Split(latlongB, ",")
	latitudeB, err := strconv.ParseFloat(res1[1], 64)
	if err != nil {
		return 0.0, fmt.Errorf("\n Erro checar LatB. %v", err)
	}
	longitudeB, err := strconv.ParseFloat(res1[0], 64)
	if err != nil {
		return 0.0, fmt.Errorf("\n Erro checar LonB. %v", err)
	}
	//distancia = latitudeA + longitudeA + latitudeB + longitudeB
	//return 0.0, fmt.Errorf("\n %f %f ", latitudeB, longitudeB)

	a := CoordenadasCartesianas(latitudeA, longitudeA)
	b := CoordenadasCartesianas(latitudeB, longitudeB)
	//return 0.0, fmt.Errorf("\n %f %f ", latitudeB, longitudeB)

	x1, y1, z1 := a[0], a[1], a[2]
	x2, y2, z2 := b[0], b[1], b[2]
	distancia = math.Sqrt(math.Pow(x2-x1, 2) + math.Pow(y2-y1, 2) + math.Pow(z2-z1, 2)) //calculando a distancia euclidiana entre os pontos A e B
	//return 0.0, fmt.Errorf("\n %f %f %f", distancia, a, b)
	return distancia, nil

}
