package main

import (
	"fmt"
	"os"
	"time"
)

var apiKey = os.Getenv("OWM_API_KEY")

func main() {
	//if err := godotenv.Load(); err != nil {
	//	logrus.Fatalf("error loading env variables: %s", err.Error())
	//}
	//
	//w, err := owm.NewCurrent("C", "ru", "60ccce665c20e05f092484ed03193b52") // fahrenheit (imperial) with Russian output
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//
	//w.CurrentByName("Moscow")
	//fmt.Printf("%#v", w)
	//w.CurrentByName("Oludeniz")
	//w.Main.Temp
	//fmt.Printf("%#v", w)

	fmt.Println(timeDurations(10))

}

func timeDurations(minutes int) time.Duration {
	return time.Duration(minutes * 1e9 * 60)
}
