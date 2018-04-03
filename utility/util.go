package utility

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

//Semester - calcula el semestre actual
func Semester() (round int) {
	NowMonth := time.Now().UTC().Month()

	if NowMonth == 1 || NowMonth == 2 || NowMonth == 3 || NowMonth == 4 || NowMonth == 5 || NowMonth == 6 {
		round = 1
	} else {
		round = 3
	}
	return round
}

//GetServiceXML - Obtiene datos desde una URL de un servicio XML
func GetServiceXML(T interface{}, url string, wg *sync.WaitGroup) error {

	response, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("GET error: %v", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("Status error: %v", response.StatusCode)
	}
	data, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return fmt.Errorf("Read body: %v", err)
	}
	xml.Unmarshal(data, &T)
	if wg != nil {
		wg.Done()
	}
	return err
}

//GetInitEnd - Retorna el comienxo y fin de un semestre en formato Time
func GetInitEnd() (fromDate time.Time, toDate time.Time) {
	Semester := Semester()
	var Inicial time.Month
	var Final time.Month
	if Semester == 1 {
		Inicial = time.January
		Final = time.June
	} else {
		Inicial = time.July
		Final = time.December
	}
	fromDate = time.Date(time.Now().UTC().Year(), Inicial, 1, 0, 0, 0, 0, time.UTC)
	toDate = time.Date(time.Now().UTC().Year(), Final, 30, 0, 0, 0, 0, time.UTC)
	return fromDate, toDate
}
