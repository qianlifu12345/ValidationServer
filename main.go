package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	simplejson "github.com/bitly/go-simplejson"
)

// type authorizedID struct {
// 	account string `json:"X-API-Account-Id"`
// 	project string `json:"X-API-Project-Id"`
// }

func getValue(path string, token string) string {
	client := &http.Client{}
	requestURL := "http://54.255.182.226:8080/v2-beta/" + path
	req, err := http.NewRequest("GET", requestURL, nil)
	cookie := http.Cookie{Name: "token", Value: token}
	req.AddCookie(&cookie)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	js, _ := simplejson.NewJson(bodyText)
	authorized, _ := js.Get("message").String()
	if authorized == "Unauthorized" {
		return "Unauthorized"
	} else {
		jsq, _ := simplejson.NewJson(bodyText)
		id, _ := jsq.Get("data").GetIndex(0).Get("id").String()
		return id
	}

}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		if r.URL.Path == "/v1-auth-filter/validateAuthToken" {
			cookie, err := r.Cookie("token")
			if err == nil {
				fmt.Fprintln(w, "Name:", cookie.Name)
				fmt.Fprintln(w, "Value:", cookie.Value)
				if cookie.Value != "" {
					accountID := getValue("accounts", cookie.Value)
					projectID := getValue("projects", cookie.Value)

					if accountID != "" && projectID != "" {
						var m map[string][]string = make(map[string][]string)
						m["headers"] = []string{"header"}
						m["X-API-Project-Id"] = []string{projectID}
						m["X-API-Account-Id"] = []string{accountID}
						if bs, err := json.Marshal(m); err != nil {
							panic(err)
						} else {
							//result --> {"C":"No.3","Go":"No.1","Java":"No.2"}
							fmt.Println(string(bs))
							fmt.Fprintln(w, string(bs))
						}
					} else if accountID != "Unauthorized" && projectID != "Unauthorized" {
						w.WriteHeader(401)
					}

					// w.Header().Add("X-API-Account-Id", accountID)
					// w.Header().Add("X-API-Project-Id", projectID)

				}

			}
		}
	}
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8011", nil))

}
