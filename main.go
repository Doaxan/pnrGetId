package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

var (
	steamUrl   = "http://steamcommunity.com/games/1029630/members?p="
	eversimUrl = "http://eversim.com/_geol/_geolreg_p19_sr.php"
	re         = regexp.MustCompile("steamcommunity.com/profiles/(.*)'")
	keyFound   = false
	bytes      []byte
)

func main() {
	fmt.Println("Start server...")
	http.HandleFunc("/_geol/_geolreg_p19_sr.php", geolreg)
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
	resp, err := http.PostForm(eversimUrl, form)
	if err != nil {
		fmt.Println(err)
		return
	}
	bytes, _ := ioutil.ReadAll(resp.Body)
	geolregExitCode, _ := strconv.Atoi(strings.Split(strings.Split(string(bytes), "\n")[0], "=")[1])
	if geolregExitCode > 0 {
		keyFound = true
		fmt.Printf("\n---ID found! Don't close this window! Use Steam ID: %s in 'Emulator Setting' tab of SmartSteamEmu and restart _start!\n", id)
		fmt.Printf("---ID найден! Не закрывайте это окно! Используйте Steam ID: %s во вкладке 'Emulator Setting' в SmartSteamEmu и перезапустите _start!\n", id)
		return bytes, true
	}
	return
}

func grabSteamIds(page int) (matches [][]string, ok bool) {
	resp, err := http.Get(steamUrl + strconv.Itoa(page))
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
