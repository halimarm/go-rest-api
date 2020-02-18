package main

import (
  "encoding/json"
	"fmt"
	"log"
	"net/http"

  _ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func getAllUsers(w http.ResponseWriter, r *http.Request) {
  var users Users
  var arr_user []Users
  var response ResUsers

  db := connect()
  defer db.Close()

  rows, err := db.Query("Select id, username, email, first_name, last_name, location from users")
  if err != nil {
    log.Print(err)
  }

  for rows.Next() {
    if err := rows.Scan(&users.Id, &users.Username, &users.Email, &users.FirstName, &users.LastName, &users.Location); err != nil {
      log.Fatal(err.Error())
    } else {
      arr_user = append(arr_user, users)
    }
  }

  response.Status = true
  response.Message = "Success"
  response.Data = arr_user

  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(response)

}

func getAllBlogs(w http.ResponseWriter, r *http.Request) {
  var blogs Blogs
  var arr_blog []Blogs
  var response ResBlogs

  db := connect()
  defer db.Close()

  rows, err := db.Query("Select id, title, slug, content from blogs")
  if err != nil {
    log.Print(err)
  }

  for rows.Next() {
    if err := rows.Scan(&blogs.Id, &blogs.Title, &blogs.Slug, &blogs.Content); err != nil {
      log.Fatal(err.Error())
    } else {
      arr_blog = append(arr_blog, blogs)
    }
  }

  response.Status = true
  response.Message = "Success"
  response.Data = arr_blog

  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(response)

}

func index(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "null")
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
  router.HandleFunc("/", index)
  router.HandleFunc("/users", getAllUsers).Methods("GET")
  router.HandleFunc("/blogs", getAllBlogs).Methods("GET")
  http.Handle("/", router)
  fmt.Println("Connected to port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
