package main

import (
	"encoding/json"
	"fmt"
	sqlserver "forms/database"
	models "forms/models"
	utils "forms/utils"
	"html/template"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/tawesoft/golib/v2/dialog"
)

func handleRequests() {
	//client side (http://localhost:9070)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", routeIndexGet)
	http.HandleFunc("/process", routeSubmitPost)
	http.HandleFunc("/about", routeAbout)
	http.HandleFunc("/policy", routePolicy)
	http.HandleFunc("/login", routeLoginGet)
	http.HandleFunc("/home", routeHomePost)
	http.HandleFunc("/profile", routeProfile)

	fmt.Println("server started at localhost:9070")
	go func() {
		http.ListenAndServe(":9070", nil)
	}()
	//server side (http://localhost:9080/users)
	myrouter := mux.NewRouter().StrictSlash(false)
	myrouter.HandleFunc("/users", allUsers).Methods("GET")
	myrouter.HandleFunc("/user/{id}", getUser).Methods("GET")
	http.ListenAndServe(":9080", myrouter)
}

func getUsersFromDb() []models.User {
	var Users []models.User
	sqlserver.Database.Find(&Users)
	return Users
}

func allUsers(w http.ResponseWriter, r *http.Request) {
	utils.JsonResponse(w, getUsersFromDb())
}

func getUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var user models.User
	result := sqlserver.Database.Where("Id = ?", id).First(&user)
	if result.Error == nil {
		utils.JsonResponse(w, user)
	}
}

func routeIndexGet(w http.ResponseWriter, r *http.Request) {
	var tmpl = template.Must(template.New("form").ParseFiles("views/regis.html"))
	var err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func routeLoginGet(w http.ResponseWriter, r *http.Request) {
	var tmpl = template.Must(template.New("form").ParseFiles("views/login.html"))
	var err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func routeAbout(w http.ResponseWriter, r *http.Request) {
	var tmpl = template.Must(template.New("aboutUs").ParseFiles("views/about.html"))
	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func routePolicy(w http.ResponseWriter, r *http.Request) {
	var tmpl = template.Must(template.New("policy").ParseFiles("views/about.html"))
	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func routeSubmitPost(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var tmpl = template.Must(template.New("hasil").ParseFiles("views/regis.html"))
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var name = r.FormValue("name")
		var phone, _ = strconv.Atoi(r.FormValue("phone"))
		var pass = r.FormValue("pass")
		var email = r.FormValue("email")
		var data = map[string]any{"name": name, "phone": phone, "password": pass, "email": email}
		//Print input data
		fmt.Println()
		fmt.Println(data)
		reqBody, _ := json.Marshal(data)
		var newUser models.User
		utils.JsonDeserialize(reqBody, &newUser)
		//Insert into database
		fmt.Println("New user added!")
		sqlserver.Database.Create(&newUser)

		if err := tmpl.Execute(w, data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	http.Error(w, "", http.StatusBadRequest)
}

var email string

func routeHomePost(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var pass = r.FormValue("pass")
		email = r.FormValue("email")
		var data = map[string]any{"password": pass, "email": email}
		fmt.Println()
		fmt.Println(data) //Show credentials
		var person models.User
		sqlserver.Database.Where("Email = ?", email).Where(sqlserver.Database.Where("Password = ?", pass)).Find(&person)
		name := person.Name
		if name == "" {
			fmt.Println("Someone tried to logged in")
			dialog.Alert("Email atau password tidak terdaftar!")
			//Reload form
			var tmpl = template.Must(template.New("form").ParseFiles("views/login.html"))
			if err := r.ParseForm(); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			if err := tmpl.Execute(w, data); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			fmt.Println(name + " has logged in")
			//Open homepage
			var tmpl = template.Must(template.New("berhasil").ParseFiles("views/home.html"))
			if err := r.ParseForm(); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			} //Pass variable person
			if err := tmpl.Execute(w, person); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	}
}

func routeProfile(w http.ResponseWriter, r *http.Request) {
	var tmpl = template.Must(template.New("profil").ParseFiles("views/profile.html"))
	var person models.User
	sqlserver.Database.Where("Email = ?", email).Find(&person)
	if err := tmpl.Execute(w, person); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	sqlserver.Init()
	handleRequests()
}
