package main

import (
	"flag"
	"fmt"
	"googleSuggest/gsuggest"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

var (
	keyword, lang, country string
	err                    error
)

func main() {

	flag.StringVar(&keyword, "keyword", "", "поисковый запрос")
	flag.StringVar(&lang, "lang", "", "язык, ISO 639-1")
	flag.StringVar(&country, "country", "", "страна, ISO 3166-1 alpha-2")
	flag.Parse()

	if keyword == "" || lang == "" || country == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	file, err := os.OpenFile(keyword+"."+lang+"."+country+".txt", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0755)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	keys := make(map[string]bool)
	keys[keyword] = false

	for {

		var parsedKey string
		for k, processed := range keys {
			if _, ok := keys[""]; ok {
				delete(keys, "")
			}
			if !processed {
				parsedKey = k
				break
			}
		}
		if parsedKey == "" {
			break
		}
		log.Printf("- %s", parsedKey)

		suggests := make([]string, 0)
		i := 0
		for {
			if i == 10 {
				panic(err)
			}

			suggests, err = gsuggest.Get(parsedKey, lang, country)
			if err != nil {
				log.Printf("\t[%d]: %s", i, err)
				sleep()
				i++
				continue
			}

			break
		}

		keys[parsedKey] = true

		for _, k := range suggests {
			if !strings.Contains(k, keyword) {
				continue
			}
			trimKey := strings.TrimSpace(k)
			if _, ok := keys[trimKey]; !ok {
				keys[trimKey] = false
				log.Printf("+ %s\n", trimKey)
				fmt.Fprintln(file, trimKey)
			}
		}
	}
}

func sleep() {
	sleepTime := rand.Intn(5)
	log.Printf("\twait for %d s", sleepTime)
	time.Sleep(time.Duration(sleepTime) * time.Second)
}
