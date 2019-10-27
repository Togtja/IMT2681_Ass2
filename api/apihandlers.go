package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"../globals"
	"google.golang.org/api/iterator"

	"../firedb"
)

//NilHandler throws a Bad Request
func NilHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Invalid request", http.StatusBadRequest)
}

//CommitsHandler handler to find the amout of commits
func CommitsHandler(w http.ResponseWriter, r *http.Request) {
	if !isGetRequest(w, r) {
		return
	}
	var repo Repos
	repo.Auth = false
	//Default public
	commitFileName := globals.PUBLIC
	auth := r.FormValue("auth")
	if auth != "" {
		//Make a personal json file for authorized users
		//Should be deleted after XX hours/Days
		repo.Auth = true
		commitFileName = auth
	} else {
		auth = globals.PUBLIC
	}
	commitFileName = commitFileName + globals.COMMITFILE
	limit, offset, ok := genericGetHandler(w, r, commitFileName, globals.COMMITDIR, &repo, auth)
	if ok == false {
		return
	}
	if limit+offset > int64(len(repo.Repos)) {
		if offset >= int64(len(repo.Repos)) {
			offset = 0
		}
		limit = int64(len(repo.Repos)) - offset

	}
	repo.Repos = repo.Repos[offset : limit+offset]
	json.NewEncoder(w).Encode(repo)

	return

}

//LangHandler handles the Programming Language requests
func LangHandler(w http.ResponseWriter, r *http.Request) {
	if !isGetRequest(w, r) {
		return
	}
	var lang Lang
	lang.Auth = false
	langFileName := globals.PUBLIC
	auth := r.FormValue("auth")
	if auth != "" {
		//Make a personal json file for authorized users
		//TODO: Should be deleted after XX hours/Days
		lang.Auth = true
		langFileName = auth
	} else {
		auth = globals.PUBLIC
	}
	langFileName = langFileName + globals.LANGFILE
	limit, offset, ok := genericGetHandler(w, r, langFileName, globals.LANGDIR, &lang, auth)
	if ok == false {
		return
	}
	if limit+offset > int64(len(lang.Language)) {
		if offset >= int64(len(lang.Language)) {
			offset = 0
		}
		limit = int64(len(lang.Language)) - offset

	}
	//TODO: Get payload
	lang.Language = lang.Language[offset : limit+offset]
	json.NewEncoder(w).Encode(lang)

}

//IssueHandler handles issue request
func IssueHandler(w http.ResponseWriter, r *http.Request) {
	if !isGetRequest(w, r) {
		return
	}
	_type := r.FormValue("type")
	if _type == "users" {

	} else if _type == "labels" {

	} else {

	}
}

//StatusHandler get status code from db/ external api and uptime and version for thid API
func StatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "only get method allowed", http.StatusNotImplemented)
		return
	}
	http.Header.Add(w.Header(), "content-type", "application/json")
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 5 {
		http.Error(w, "Expecting format .../", http.StatusBadRequest)
		return
	}
	gitlab, err := http.Get("https://git.gvk.idi.ntnu.no/api/v4/projects")
	if err != nil {
		errmsg := "The HTTP request failed with error" + err.Error()
		http.Error(w, errmsg, http.StatusInternalServerError)
		return
	}
	//TODO: FIX THIS (I.E call database)
	db, err := http.Get("https://restcountries.eu/rest/v2/")
	if err != nil {
		errmsg := "The HTTP request failed with error" + err.Error()
		http.Error(w, errmsg, http.StatusInternalServerError)
		return
	}
	uptimeString := fmt.Sprintf("%.0f seconds", uptime().Seconds())
	diag := Status{gitlab.StatusCode, db.StatusCode, uptimeString, globals.Version}
	json.NewEncoder(w).Encode(diag)
	return

}

//WebhookHandler the handler for webhooks
func WebhookHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		//Get payload
		var webhook Webhook
		err := json.NewDecoder(r.Body).Decode(&webhook)
		if err != nil {
			http.Error(w, "Something went wrong: "+err.Error(), http.StatusBadRequest)
			return
		}

		//Finds all the current ids
		iter := firedb.Client.Collection("webhooks").Documents(firedb.Ctx)
		var ids []int
		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return
			}
			//Should be no errors, as the ID is a string that we insert
			id, _ := strconv.Atoi(doc.Ref.ID)
			ids = append(ids, id)
		}

		//Finds a new ID (ID's starts from 1)
		var newid int
		sort.Ints(ids)
		newid = 1
		for _, id := range ids {
			if id == newid {
				newid++
			} else {
				//We found an unused id
				break
			}
		}
		if webhook.Event == "" || webhook.URL == "" {
			http.Error(w, "Please provide both event and url", http.StatusBadRequest)
			return
		}
		//Some form of type cheking
		webhook.Event = strings.ToLower(webhook.Event)
		if !EventOK(webhook.Event) {
			http.Error(w, "Please provide an event of type (commits|languages|issues|status)", http.StatusBadRequest)
			return
		}
		webhook.ID = strconv.Itoa(newid)
		webhook.Time = time.Now().String()
		_, err = firedb.Client.Collection("webhooks").Doc(strconv.Itoa(newid)).Set(firedb.Ctx, webhook)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(webhook)

	case http.MethodGet:
		http.Header.Add(w.Header(), "content-type", "application/json")
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) != 5 {
			http.Error(w, "Expecting format webhooks/<id>", http.StatusBadRequest)
			return
		}
		if parts[4] != "" {
			id, err := strconv.Atoi(parts[4])
			if err != nil {
				http.Error(w, "Could not turn id into a integer", http.StatusBadRequest)
				return
			}
			doc, err := firedb.Client.Collection("webhooks").Doc(strconv.Itoa(id)).Get(firedb.Ctx)
			if err != nil {
				fmt.Println(err)
				return
			}
			m := doc.Data()
			//Create and Encode the struct
			var event, time string = fmt.Sprint(m["Event"]), fmt.Sprint(m["Time"])
			json.NewEncoder(w).Encode(WebhookGet{id, event, time})
			return
		}
		// For now just return all webhooks, don't respond to specific resource requests
		iter := firedb.Client.Collection("webhooks").Documents(firedb.Ctx)
		var webhooks []WebhookGet
		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				fmt.Println(err)
				return
			}
			m := doc.Data()
			var id, event, time string = fmt.Sprint(m["ID"]), fmt.Sprint(m["Event"]), fmt.Sprint(m["Time"])
			wid, _ := strconv.Atoi(id)
			webhooks = append(webhooks, WebhookGet{wid, event, time})

		}
		json.NewEncoder(w).Encode(webhooks)
	case http.MethodDelete:
		//TODO: Do deleting stuff
	default:
		http.Error(w, "Invalid method "+r.Method, http.StatusBadRequest)
	}
	return
}
