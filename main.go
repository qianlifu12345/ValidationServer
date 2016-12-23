package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	simplejson "github.com/bitly/go-simplejson"
)

type authorizedID struct {
	Name map[string][]string
}

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
	if r.Method == "GET" {
		if r.URL.Path == "/v1-auth-filter/validateAuthToken" {

			cookie, err := r.Cookie("token")
			if err == nil {
				// fmt.Fprintln(w, "Name:", cookie.Name)
				// fmt.Fprintln(w, "Value:", cookie.Value)
				if cookie.Value != "" {
					accountID := getValue("accounts", cookie.Value)
					projectID := getValue("projects", cookie.Value)
					// fmt.Fprintf(w, "X-API-Account-Id:%q X-API-Project-Id:%q\n", accountID, projectID)
					if accountID != "" && projectID != "" {
						var result authorizedID
						result = map[string][]string{
							"header": []string{"", ""},
                            "X-API-Account-Id": []string{accountID},
                            "X-API-Project-Id": []string{projectID},
                    }}
						body, err := json.Marshal(result)
						if err != nil {
							panic(err.Error())
						}
						body2, _ := body.String()
						fmt.Fprintln(w, "Value:", body2)
						w.WriteHeader(200)
					}

					// w.Header().Add("X-API-Account-Id", accountID)
					// w.Header().Add("X-API-Project-Id", projectID)

				} else if accountID != "Unauthorized" && projectID != "Unauthorized" {
					w.WriteHeader(401)
				}
			}

		}
	}
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8011", nil))

}
