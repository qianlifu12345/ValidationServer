package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	simplejson "github.com/bitly/go-simplejson"
)

//get the projectID and accountID from rancher API
func getValue(host string, path string, token string) string {
	result := ""
	client := &http.Client{}
	requestURL := host + "v2-beta/" + path
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
		result = "Unauthorized"
	} else {
		jsonBody, _ := simplejson.NewJson(bodyText)
		id, err := jsonBody.Get("data").GetIndex(0).Get("id").String()
		if err != nil {
			log.Fatal(err)
			result = "No id found"
		} else {
			result = id
		}

	}

	return result
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		if r.URL.Path == "/v1-auth-filter/validateAuthToken" {
			cookie, err := r.Cookie("token")
			if err == nil {
				//check if the token value is empty or not
				if cookie.Value != "" {
					accountID := getValue("http://54.255.182.226:8080/", "accounts", cookie.Value)
					projectID := getValue("http://54.255.182.226:8080/", "projects", cookie.Value)
					//check if the accountID or projectID is empty
					if accountID != "" && projectID != "" {
						if accountID == "Unauthorized" && projectID == "Unauthorized" {
							w.WriteHeader(401)
						} else {
							//construct the responseBody
							var responseBody map[string][]string = make(map[string][]string)
							for k, v := range r.Header {
								responseBody[k] = v
							}
							responseBody["X-API-Project-Id"] = []string{projectID}
							responseBody["X-API-Account-Id"] = []string{accountID}
							if responseBodyString, err := json.Marshal(responseBody); err != nil {
								panic(err)
							} else {
								fmt.Fprintln(w, string(responseBodyString))
							}
						}
					}

				}

			}
		}
	}
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8011", nil))

}
