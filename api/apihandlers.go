package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"RESTGvkGitLab/globals"

	"google.golang.org/api/iterator"

	"RESTGvkGitLab/firedb"
)

//SetupHandlers for easier test
func SetupHandlers() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", NilHandler)
	mux.HandleFunc("/repocheck/v1/commits/", CommitsHandler)
	mux.HandleFunc("/repocheck/v1/languages/", LangHandler)
	mux.HandleFunc("/repocheck/v1/issues/", IssueHandler)
	mux.HandleFunc("/repocheck/v1/status/", StatusHandler)
	mux.HandleFunc("/repocheck/v1/webhooks/", WebhookHandler)
	return mux
}

//NilHandler throws a Bad Request
func NilHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Invalid request", http.StatusBadRequest)
}

//CommitsHandler handler to find the amount of commits
func CommitsHandler(w http.ResponseWriter, r *http.Request) {
	if !isGetRequest(w, r) {
		return
	}
	var repo Repos
	repo.Auth = false
	//Default public
	commitFileName := globals.PUBLIC
	auth := r.FormValue(globals.AUTHP)
	if auth != "" {
		//Make a personal json file for authorized users
		//Should be deleted after XX hours/Days
		repo.Auth = true
		commitFileName = auth
	} else {
		auth = globals.PUBLIC
	}
	commitFileName = commitFileName + globals.COMMITFILE
	_, limit, offset, ok := genericGetHandler(w, r, commitFileName, globals.COMMITDIR, &repo, auth)
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
	http.Header.Add(w.Header(), "content-type", "application/json")
	json.NewEncoder(w).Encode(repo)

	param := []string{strconv.FormatInt(limit, 10), strconv.FormatInt(offset, 10), strconv.FormatBool(repo.Auth)}
	err := activateWebhook(globals.CommitE, param)
	if err != nil {
		//No need to throw a webhook error to user, so just print it for sys admin
		fmt.Println("Some error involving activating webhook:", err)
	}
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
	auth := r.FormValue(globals.AUTHP)
	if auth != "" {
		//Makes a personal json file for authorized users
		lang.Auth = true
		langFileName = auth
	} else {
		auth = globals.PUBLIC
	}
	langFileName = langFileName + globals.LANGFILE
	projects, limit, offset, ok := genericGetHandler(w, r, langFileName, globals.LANGDIR, &lang, auth)
	//general rule: if takes ResponseWriter, and it returns a bool and not an err,
	//the error handeling is done in the function
	if ok == false {
		return
	}
	var payload PayloadLang
	ok, err := GetPayload(r, &payload)
	if err != nil {
		fmt.Println("Failed to read body", err.Error())
		http.Error(w, "Invalid Body", http.StatusBadRequest)
	}
	//Does not care if body/payload is empty
	if ok && len(payload.ProjectName) > 0 {
		var temp []Project
		for i := range payload.ProjectName {
			for j := range projects {
				if payload.ProjectName[i] == projects[j].NamePath {
					temp = append(temp, projects[j])
					break
				}
			}
		}
		//Find languages based of the projects
		lang.Language = subAPICallsForLang(temp, auth, w)
	}
	//Make sure limit is in range
	if limit+offset > int64(len(lang.Language)) {
		if offset >= int64(len(lang.Language)) {
			offset = 0
		}
		limit = int64(len(lang.Language)) - offset

	}
	lang.Language = lang.Language[offset : limit+offset]

	http.Header.Add(w.Header(), "content-type", "application/json")
	json.NewEncoder(w).Encode(lang)
	//I decide to send the actual limit the user get and not what it asks for
	param := []string{strconv.FormatInt(limit, 10), strconv.FormatInt(offset, 10), strconv.FormatBool(lang.Auth)}
	err = activateWebhook(globals.LanguagesE, param)
	if err != nil {
		//No need to throw a webhook error to user, so just print it for sys admin
		fmt.Println("Some error involving activating webhook:", err)
	}
}

//IssueHandler handles issue request
func IssueHandler(w http.ResponseWriter, r *http.Request) {
	if !isGetRequest(w, r) {
		return
	}
	var authBool bool
	auth := r.FormValue("auth")
	if auth != "" {
		//Make a personal json file for authorized users
		authBool = true
	} else {
		auth = globals.PUBLIC
		authBool = false
	}
	var project PayloadIssue
	_, err := GetPayload(r, &project)
	if err != nil {
		fmt.Println("Failed to read body", err.Error())
		http.Error(w, "Invalid Body", http.StatusBadRequest)
		return
	}
	_type := r.FormValue("type")

	http.Header.Add(w.Header(), "content-type", "application/json")
	//TODO: find out what I need to count for users
	issues, err := findIssuesForProject(project, auth, w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if len(issues) == 0 {
		http.Error(w, "Could not find any issues with that name", http.StatusBadRequest)
		return
	}
	if _type == "users" {
		users := findAuthorsInIssues(issues, authBool)
		if len(users.Users) == 0 {
			http.Error(w, "Could not find issues by users in the project", http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(users)
	} else if _type == "labels" {
		labels := findLabelsInIssues(issues, authBool)
		if len(labels.Labels) == 0 {
			http.Error(w, "Could not find issues in the project", http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(labels)
	} else {
		http.Error(w, "Invalid type", http.StatusBadRequest)
		return
	}
	param := []string{_type, strconv.FormatBool(authBool)}
	err = activateWebhook(globals.IssuesE, param)
	if err != nil {
		//No need to throw a webhook error to user, so just print it for sys admin
		fmt.Println("Some error involving activating webhook:", err)
	}
}

//StatusHandler get status code from db/ external api and uptime and version for this API
func StatusHandler(w http.ResponseWriter, r *http.Request) {
	if !isGetRequest(w, r) {
		return
	}
	gitlab, err := http.Get("https://git.gvk.idi.ntnu.no/api/v4/projects")
	if err != nil {
		errmsg := "The HTTP request failed with error" + err.Error()
		http.Error(w, errmsg, http.StatusInternalServerError)
		return
	}

	dbStatus := 200 //Assumes it to be ok
	//Doc("0") is reserved for status checks
	_, err = firedb.Client.Collection(globals.WebhookF).Doc("0").Get(firedb.Ctx)
	if err != nil {
		//Can not get to server for unknown reason
		//Gives a Service Unavailable error
		dbStatus = 503
	}
	uptimeString := fmt.Sprintf("%.0f seconds", uptime().Seconds())

	http.Header.Add(w.Header(), "content-type", "application/json")
	diag := Status{gitlab.StatusCode, dbStatus, uptimeString, globals.Version}
	json.NewEncoder(w).Encode(diag)
	var param []string //Empty parameters, as Status does not take in parameters
	err = activateWebhook(globals.StatusE, param)
	if err != nil {
		//No need to throw a webhook error to user, so just print it for sys admin
		fmt.Println("Some error involving activating webhook:", err)
	}
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
			http.Error(w, "Invalid Body: "+err.Error(), http.StatusBadRequest)
			return
		}

		//Finds all the current ids
		iter := firedb.Client.Collection(globals.WebhookF).Documents(firedb.Ctx)
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
		_, err = firedb.Client.Collection(globals.WebhookF).Doc(strconv.Itoa(newid)).Set(firedb.Ctx, webhook)
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
		//if Id provided
		if parts[4] != "" {
			id, err := strconv.Atoi(parts[4])
			if err != nil {

				return
			}
			doc, err := firedb.Client.Collection(globals.WebhookF).Doc(strconv.Itoa(id)).Get(firedb.Ctx)
			if err != nil {
				http.Error(w, "Could not find ID", http.StatusBadRequest)
				fmt.Println(err)
				return
			}
			m := doc.Data()
			//Create and Encode the struct
			var event, time string = fmt.Sprint(m[globals.EventF]), fmt.Sprint(m[globals.TimeF])
			json.NewEncoder(w).Encode(WebhookGet{id, event, time})
			return
		}
		// For now just return all webhooks, don't respond to specific resource requests
		iter := firedb.Client.Collection(globals.WebhookF).Documents(firedb.Ctx)
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
			var id, event, time string = fmt.Sprint(m[globals.IDF]), fmt.Sprint(m[globals.EventF]), fmt.Sprint(m[globals.TimeF])
			wid, _ := strconv.Atoi(id)
			webhooks = append(webhooks, WebhookGet{wid, event, time})

		}
		sort.SliceStable(webhooks, func(i, j int) bool {
			return webhooks[i].ID < webhooks[j].ID
		})
		json.NewEncoder(w).Encode(webhooks[1:])
	case http.MethodDelete:
		var delID int
		ok, err := GetPayload(r, &delID)
		if err != nil {
			http.Error(w, "Invalid Body", http.StatusBadRequest)
		}
		//We did not get ID in body
		if !ok {
			parts := strings.Split(r.URL.Path, "/")
			//We will now see if it was part of the request
			if parts[4] != "" {
				delID, err = strconv.Atoi(parts[4])
				if err != nil {
					http.Error(w, "Could not turn id into a integer", http.StatusBadRequest)
					return
				}
			} else {
				http.Error(w, "No form of resource ID provided", http.StatusBadRequest)
				return
			}
		}
		if delID == 0 {
			http.Error(w, "Failed to find resource ID", http.StatusBadRequest)
			return
		}
		_, err = firedb.Client.Collection(globals.WebhookF).Doc(strconv.Itoa(delID)).Delete(firedb.Ctx)
		if err != nil {
			http.Error(w, "Could not Delete the resource", http.StatusInternalServerError)
			return
		}
		msg := "Successfully deleted resource" + strconv.Itoa(delID)
		bmsg := []byte(msg)
		w.Write(bmsg)
	default:
		http.Error(w, "Invalid method "+r.Method, http.StatusBadRequest)
	}
	return
}
