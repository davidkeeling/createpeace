package createpeace

import (
    "net/http"
    "html/template"
    "appengine"
    "appengine/user"
    "appengine/datastore"
)

type ActOfPeace struct {
    ProjectName string
    ProjectDescription string
    ContactInfo string
    AreaOfFocus string
}

func init() {
 http.HandleFunc("/", index)
 http.HandleFunc("/createact", createact)
 http.HandleFunc("/projectlist", projectlist)
}


func index (w http.ResponseWriter, r *http.Request) {
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
    
    creationForm, err := template.ParseFiles("createpeacetemplate.html")
    if err != nil {
        c.Warningf("Error parsing template file: %s", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    creationForm.Execute(w, u.String())
}


func createact (w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    newProject := ActOfPeace{
        ProjectName: r.FormValue("projectname"),
        ProjectDescription: r.FormValue("projectdescription"),
        ContactInfo: r.FormValue("contactinfo"),
        AreaOfFocus: r.FormValue("focusarea"),
    }
    
    _, err := datastore.Put(c, datastore.NewIncompleteKey(c, "ActOfPeace", ancestorKey(c)), &newProject)
    if err != nil {
        c.Warningf("Error saving the project: %s", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    http.Redirect(w, r, "/projectlist", http.StatusSeeOther)
}


func projectlist (w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    listTemplate, err := template.ParseFiles("projectlisttemplate.html")
    if err != nil {
        c.Warningf("Error parsing template file: %s", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    q := datastore.NewQuery("ActOfPeace").Ancestor(ancestorKey(c)).Order("ProjectName").Limit(50)
    theActs := make([]ActOfPeace, 0, 50)
    if _, err := q.GetAll(c, &theActs); err != nil {
        c.Warningf("Error retrieving projects: %s", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    if len(theActs) == 0 {
        blankproject := ActOfPeace{
            ProjectName: "Your project...",
            ProjectDescription: "...will be the first!",
            ContactInfo: "What are you going to do?",
            AreaOfFocus: "",
        }
        var noprojects []ActOfPeace
        noprojects = append(noprojects, blankproject)
        listTemplate.Execute(w, noprojects)
        return
    }

    listTemplate.Execute(w, theActs)
}


func ancestorKey (c appengine.Context) *datastore.Key {
    return datastore.NewKey(c, "OngoingActs", "default_actofpeace", 0, nil)
}