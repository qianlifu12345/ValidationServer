package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	simplejson "github.com/bitly/go-simplejson"
)

func basicAuth(user string, pwd string) string {
	// var username string = "9E7544ADF1AA1E188FBC"
	// var passwd string = "nfBa3HwdruPxwGtHrpdb3sSkmn68bKCuquSLZKTv"
	var username string = user
	var passwd string = pwd
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://54.255.182.226:8080/v2-beta/accounts", nil)
	req.SetBasicAuth(username, passwd)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	// s := string(bodyText)
	// fmt.Printf(string(bodyText))
	js, _ := simplejson.NewJson(bodyText)
	authorized, _ := js.Get("message").String()
	if authorized == "Unauthorized" {
		return "Error 401 unauthorized"
	} else {
		jsq, _ := simplejson.NewJson(bodyText)
		// fmt.Print(bodyText)
		accountID, _ := jsq.Get("data").GetIndex(0).Get("id").String()
		return accountID
	}

	// js := gojson.Json(s).Get("data")
	// js, err := simplejson.NewJson(bodyText)
	// temp, _ := js.Get("data").Get("id").String()
	// temp, _ := js.String()
	// fmt.Printf("value")

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
				fmt.Fprintln(w, "Name:", cookie.Name)
				fmt.Fprintln(w, "Value:", cookie.Value)
				if cookie.Value != "" {
					accountID := getValue("accounts", cookie.Value)
					projectID := getValue("projects", cookie.Value)
					fmt.Fprintf(w, "X-API-Account-Id:%q X-API-Project-Id:%q\n", accountID, projectID)
					if accountID != "" && projectID != "" {
						w.Header().Add("X-API-Account-Id", accountID)
						w.Header().Add("X-API-Project-Id", projectID)
						w.WriteHeader(http.StatusOK)
					} else if accountID != "Unauthorized" && projectID != "Unauthorized" {
						w.WriteHeader(http.StatusUnauthorized)
					}
				}

			} else {
				w.WriteHeader(http.StatusBadRequest)
			}

		}

	} else if r.Method == "POST" {
		if r.URL.Path == "/post" {
			result, _ := ioutil.ReadAll(r.Body)
			r.Body.Close()
			fmt.Printf("POST %s\n", result)
			token, _ := simplejson.NewJson(result)
			usr, _ := token.Get("usr").String()
			pwd, _ := token.Get("pwd").String()
			// fmt.Printf(usr + "" + pwd)
			accountID := basicAuth(usr, pwd)
			// fmt.Printf("result" + accountID)
			// jsonObject, _ := json.Marshal(temp)
			fmt.Fprintf(w, "%q\n", accountID)
		}
	}
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8000", nil))

}
