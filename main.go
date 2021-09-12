package main

import (
	"fmt"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

var (
	steamUrls		 = [3]string{"http://steamcommunity.com/games/1029630/members?p=", "http://steamcommunity.com/games/1379930/members?p=", "http://steamcommunity.com/games/1683320/members?p="}
	eversimUrls      = [3]string{"http://eversim.com/_geol/_geolreg_p19_sr.php", "http://eversim.com/_geol/_geolreg_p20_sr.php", "http://eversim.com/_geol/_geolreg_p21_sr.php"}
	eversimHandlers  = [3]string{"/_geol/_geolreg_p19_sr.php", "/_geol/_geolreg_p20_sr.php", "/_geol/_geolreg_p21_sr.php"}
	re               = regexp.MustCompile("steamcommunity.com/profiles/(.*)'")
	keyFound         = false
	year       		 = 2
	bytes            []byte
)

func main() {
	yearPtr := flag.Int("year", 2021, "Edition of the game")
	flag.Parse()
	year = *yearPtr - 2019
	if year != 2 {
		fmt.Println("pnrGetID is running in legacy mode. Please be aware that legacy mode activates only 2019 and 2020 Edition of Power and Revolution.")
	}

	if year < 0 || year >= len(steamUrls) {
		fmt.Println("Wrong year has been selected. 2021 will be used instead.")
		year = 2
	}

	fmt.Println("Starting server...")
	http.HandleFunc(eversimHandlers[year], geolreg)
	
	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Press the any Key to close window. Have a good day :)")
	fmt.Scanln() // wait for Enter Key
}

func geolreg(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\nGet request from _start.exe")
	if err := r.ParseForm(); err != nil {
		log.Fatal(err)
	}

	if keyFound {
		if _, err := w.Write(bytes); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Saved key has been sent! Registration done! Enjoy :)")
		return
	}
	fmt.Println("Get steam ids and find free id:")
	for i := 65; i > 0; i-- {
		if steamIds, ok := grabSteamIds(i); ok {
			for _, id := range steamIds[1] {
				fmt.Print(".")
				if bytes, ok = geolregPost(id, r.Form); ok {
					if _, err := w.Write(bytes); err != nil {
						log.Fatal(err)
					}
					return
				}
			}
		}
	}
}

func geolregPost(id string, form url.Values) (body []byte, ok bool) {
	form["v4"] = []string{"#SK" + id + "ZZZZZZZZZZ"}
	form["v6"] = []string{id + "@steam."}
	resp, err := http.PostForm(eversimUrls[year], form)
	
	if err != nil {
		fmt.Println(err)
		return
	}
	bytes, _ := ioutil.ReadAll(resp.Body)
	geolregExitCode, _ := strconv.Atoi(strings.Split(strings.Split(string(bytes), "\n")[0], "=")[1])
	if geolregExitCode > 0 {
		keyFound = true
		fmt.Printf("\n---ID found! Don't close this window! Use Steam ID: %s in 'Emulator Setting' tab of SmartSteamEmu and restart _start!\n", id)
		fmt.Printf("---ID найден! Не закрывайте это окно! Пропишите Steam ID: %s в файле \"%AppData%/Goldberg SteamEmu Saves\\settings\\user_steam_id.txt\" и перезапустите _start!\n", id)
		return bytes, true
	}
	return
}

func grabSteamIds(page int) (matches [][]string, ok bool) {
	resp, err := http.Get(steamUrls[year] + strconv.Itoa(page))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	bytes, _ := ioutil.ReadAll(resp.Body)
	matches = re.FindAllStringSubmatch(string(bytes), -1)
	if len(matches) == 0 {
		return
	}
	return matches, true
}
