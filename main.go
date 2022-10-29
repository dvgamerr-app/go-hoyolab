package main

import (
	"fmt"
	"os"

	"github.com/zellyn/kooky"
	_ "github.com/zellyn/kooky/browser/chrome" // register cookie store finders!
)

// var json jsoniter.API = jsoniter.ConfigCompatibleWithStandardLibrary

func main() {
	daily := &DailyHoyolab{
		ActID:    "e202102251931481",
		Endpoint: "https://sg-hk4e-api.hoyolab.com",
		Sign:     "/event/sol/sign",
		Info:     "/event/sol/info",
		Lang:     "en-us",
	}

	for _, store := range kooky.FindAllCookieStores() {
		_, err := os.Stat(store.FilePath())

		if store.Browser() != "chrome" || os.IsNotExist(err) {
			continue
		}

		fmt.Printf("%+v\n", store.FilePath())
	}

	if err := daily.CheckServer(); err != nil {
		fmt.Printf("API: %v", err)
	}

	// cookies := kooky.ReadCookies(kooky.DomainContains("hoyolab.com"))
	// for _, cookie := range cookies {
	// 	fmt.Printf("%+v\n", cookie)
	// }

	// dir, _ := os.UserCacheDir() // "/<USER>/Library/Application Support/"
	// cookiesFile := path.Join(dir, "Google/Chrome/User Data/Default/Network/Cookies")
	// cookies, err := chrome.ReadCookies(cookiesFile)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// r, _ := regexp.Compile(".+?hoyolab.com$")
	// for _, cookie := range cookies {
	// 	if r.MatchString(cookie.Domain) {
	// 		fmt.Println(cookie.Cookie)
	// 	}
	// }
}

// fetch("https://sg-hk4e-api.hoyolab.com/event/sol/sign?lang=en-us", {
//   "headers": {
//     "accept": "*/*",
//     "accept-language": "en-US,en;q=0.9,th;q=0.8",
//     "cache-control": "no-cache",
//     "pragma": "no-cache",
//     "sec-fetch-dest": "empty",
//     "sec-fetch-mode": "cors",
//     "sec-fetch-site": "same-site"
//   },
//   "referrer": "https://act.hoyolab.com/",
//   "referrerPolicy": "strict-origin-when-cross-origin",
//   "body": null,
//   "method": "OPTIONS",
//   "mode": "cors",
//   "credentials": "omit"
// });

// fetch("https://sg-hk4e-api.hoyolab.com/event/sol/sign?lang=en-us", {
//   "headers": {
//     "accept": "application/json, text/plain, */*",
//     "accept-language": "en-US,en;q=0.9,th;q=0.8",
//     "cache-control": "no-cache",
//     "content-type": "application/json;charset=UTF-8",
//     "pragma": "no-cache",
//     "sec-ch-ua": "\"Chromium\";v=\"106\", \"Google Chrome\";v=\"106\", \"Not;A=Brand\";v=\"99\"",
//     "sec-ch-ua-mobile": "?0",
//     "sec-ch-ua-platform": "\"Windows\"",
//     "sec-fetch-dest": "empty",
//     "sec-fetch-mode": "cors",
//     "sec-fetch-site": "same-site"
//   },
//   "referrer": "https://act.hoyolab.com/",
//   "referrerPolicy": "strict-origin-when-cross-origin",
//   "body": "{\"act_id\":\"e202102251931481\"}",
//   "method": "POST",
//   "mode": "cors",
//   "credentials": "include"
// });
// cookie: mi18nLang=en-us; _MHYUUID=5ca37bbe-de66-4c87-8509-e3e51eb94e58; DEVICEFP_SEED_ID=291b04db126bbf2f; DEVICEFP_SEED_TIME=1666923188358; DEVICEFP=38d7eafefdad9; _ga=GA1.2.1160952229.1666923189; _gid=GA1.2.2134928089.1666923189; _gat_gtag_UA_201411121_1=1; ltoken=jcoOfAJNBFLHj4KK4nT87IX1a97Wj0weMvr9OBxE; ltuid=1104305; cookie_token=2xHLthsjoMlWkdWTDqoJcXrS7NosQsOWBxN0P0rj; account_id=1104305; _ga_54PBK3QDF4=GS1.1.1666923189.1.1.1666923238.0.0.0

// fetch("https://sg-hk4e-api.hoyolab.com/event/sol/info?lang=en-us&act_id=e202102251931481", {
//   "headers": {
//     "accept": "application/json, text/plain, */*",
//     "accept-language": "en-US,en;q=0.9,th;q=0.8",
//     "cache-control": "no-cache",
//     "pragma": "no-cache",
//     "sec-ch-ua": "\"Chromium\";v=\"106\", \"Google Chrome\";v=\"106\", \"Not;A=Brand\";v=\"99\"",
//     "sec-ch-ua-mobile": "?0",
//     "sec-ch-ua-platform": "\"Windows\"",
//     "sec-fetch-dest": "empty",
//     "sec-fetch-mode": "cors",
//     "sec-fetch-site": "same-site"
//   },
//   "referrer": "https://act.hoyolab.com/",
//   "referrerPolicy": "strict-origin-when-cross-origin",
//   "body": null,
//   "method": "GET",
//   "mode": "cors",
//   "credentials": "include"
// });
