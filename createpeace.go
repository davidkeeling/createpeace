package createpeace

import (
	"appengine"
	"appengine/datastore"
	"appengine/user"
	"encoding/json"
	"html/template"
	"net/http"
	"time"
)

type ActOfPeace struct {
	ProjectName        string
	ProjectDescription string
	ContactInfo        string
	AreaOfFocus        string
	WhenCreated        time.Time
}

var funcMap = template.FuncMap{
	"getFocusAreas":    getFocusAreas,
	"getFocusAreaURL":  getFocusAreaURL,
	"getFocusAreaJSON": getFocusAreaJSON,
}

var focusAreaList = map[string]string{
	"Access to Water and Natural Resources": "jonata-Water-bottle.png",
	"Advancing Women and Children":          "mother-son-line-art-ArtFavor-Mom-holding-childs-hand.png",
	"Education and Community Development":   "eleve-posant-une-question.png",
	"Global Health and Wellness":            "papapishu-Doctor-examining-a-patient.png",
	"Environmental Sustainability":          "sheikh-tuhin-Save-Environment.png",
	"Conflict Resolution":                   "cyberscooty-shaking-hands3.png",
	"Inclusivity and Cooperation":           "tetheredtogether.png",
	"Human Rights":                          "Human-rights-icon.png",
	"Alleviation of Extreme Poverty":        "chovynz-Money-Bag-Icon.png",
	"Weapons Access and Proliferation":      "no-guns.png",
}

func getFocusAreas() map[string]string {
	return focusAreaList
}

func getFocusAreaURL(theKey string) (theValue string) {
	return focusAreaList[theKey]
}

func getFocusAreaJSON() (focusAreaJSON string) {
	theJSON, _ := json.Marshal(focusAreaList)
	return string(theJSON)
}

func index(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	if u == nil {
		url, err := user.LoginURL(c, r.URL.String())
		if err != nil {
			c.Warningf("Error logging in: %s", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusFound)
		return
	}

	creationForm := template.Must(template.New("creationForm").Funcs(funcMap).ParseFiles("createpeacetemplate.html"))

	creationForm.ExecuteTemplate(w, "createpeacetemplate.html", u)
}

func createact(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	newProject := ActOfPeace{
		ProjectName:        r.FormValue("projectname"),
		ProjectDescription: r.FormValue("projectdescription"),
		ContactInfo:        r.FormValue("contactinfo"),
		AreaOfFocus:        r.FormValue("focusarea"),
		WhenCreated:        time.Now(),
	}

	_, err := datastore.Put(c, datastore.NewIncompleteKey(c, "ActOfPeace", ancestorKey(c)), &newProject)
	if err != nil {
		c.Warningf("Error saving the project: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/projectlist", http.StatusSeeOther)
}

func projectlist(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	listTemplate := template.Must(template.New("listTemplate").Funcs(funcMap).ParseFiles("projectlisttemplate.html"))

	q := datastore.NewQuery("ActOfPeace").Ancestor(ancestorKey(c)).Order("-WhenCreated").Limit(50)
	theActs := make([]ActOfPeace, 0, 50)
	if _, err := q.GetAll(c, &theActs); err != nil {
		c.Warningf("Error retrieving projects: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(theActs) == 0 {
		blankproject := ActOfPeace{
			ProjectName:        "Your project...",
			ProjectDescription: "...will be the first!",
			ContactInfo:        "What are you going to do?",
			AreaOfFocus:        "Human Rights",
		}
		var noprojects []ActOfPeace
		noprojects = append(noprojects, blankproject)
		listTemplate.ExecuteTemplate(w, "projectlisttemplate.html", noprojects)
		return
	}

	listTemplate.ExecuteTemplate(w, "projectlisttemplate.html", theActs)
}

func ancestorKey(c appengine.Context) *datastore.Key {
	return datastore.NewKey(c, "OngoingActs", "default_actofpeace", 0, nil)
}

func init() {
	http.HandleFunc("/", index)
	http.HandleFunc("/createact", createact)
	http.HandleFunc("/projectlist", projectlist)

}
