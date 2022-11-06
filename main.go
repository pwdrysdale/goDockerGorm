package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func createUser (name string) int {
	db, err := gorm.Open(mysql.Open("tester:secret@tcp(db:3306)/test"), &gorm.Config{})
	if err != nil {
		log.Print(err.Error())
	}

	var user User = User{Name: name}
	
	db.Create(&user)

	return user.ID
}

func createRandomUser() int {

	var firstNames []string = []string{"John", "Paul", "George", "Ringo", "Pete", "Stuart", "Mick", "Keith", "Brian", "Roger"}
	var lastNames []string = []string{"Lennon", "McCartney", "Harrison", "Starr", "Townshend", "Sutcliffe", "Jagger", "Richards", "Jones", "Daltrey"}

	var name string = ""
	name += firstNames[rand.Intn(len(firstNames))]
	name += " " + lastNames[rand.Intn(len(lastNames))]

	return createUser(name)
}

func getUsers () []*User {
	db, err := gorm.Open(mysql.Open("tester:secret@tcp(db:3306)/test"), &gorm.Config{})
	if err != nil {
		log.Print(err.Error())
	}

	var users []*User
	db.Find(&users)

	return users
}

func editUser (id int, name string) {
	db, err := gorm.Open(mysql.Open("tester:secret@tcp(db:3306)/test"), &gorm.Config{})
	if err != nil {
		log.Print(err.Error())
	}

	var user User
	db.First(&user, id)
	user.Name = name
	db.Save(&user)
}

func deleteUser (id int) {
	db, err := gorm.Open(mysql.Open("tester:secret@tcp(db:3306)/test"), &gorm.Config{})
	if err != nil {
		log.Print(err.Error())
	}

	var user User
	db.First(&user, id)
	db.Delete(&user)
}

func createRandomUserPage(w http.ResponseWriter, r *http.Request) {
	id := createRandomUser()

	fmt.Println("Endpoint Hit: createUserPage")
	json.NewEncoder(w).Encode(id)
}

func usersPage(w http.ResponseWriter, r *http.Request) {
	users := getUsers()

	json.NewEncoder(w).Encode(users)
}

func createUserPage(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var user User
	json.Unmarshal(reqBody, &user)
	id := createUser(user.Name)

	json.NewEncoder(w).Encode(id)
}

func editUserPage (w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var user User
	json.Unmarshal(reqBody, &user)
	editUser(user.ID, user.Name)
}

func deleteUserPage (w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var user User
	json.Unmarshal(reqBody, &user)
	deleteUser(user.ID)
}



func main() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/createrandomuser", createRandomUserPage)
	http.HandleFunc("/users", usersPage)
	http.HandleFunc("/createuser", createUserPage)
	http.HandleFunc("/edituser", editUserPage)
	http.HandleFunc("/deleteuser", deleteUserPage)
	
	log.Fatal(http.ListenAndServe(":8080", nil))
}
