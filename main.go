package main

import (
  "encoding/json"
	"fmt"
	"log"
	"net/http"
  "strings"

  _ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
  "github.com/dgrijalva/jwt-go"
  "github.com/gorilla/context"
  "github.com/mitchellh/mapstructure"
)


func CreateTokenEndpoint(w http.ResponseWriter, req *http.Request) {
  var users Users
  _ = json.NewDecoder(req.Body).Decode(&users)
  token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims {
    "username": users.Username,
    "password": users.Password,
  })
  tokenString, error := token.SignedString([]byte("secret"))
  if error != nil {
    fmt.Print(error)
  }
  json.NewEncoder(w).Encode(JwtToken{Token: tokenString})
}

func ProtectedEndpoint(w http.ResponseWriter, req *http.Request) {
  params := req.URL.Query()
  token, _ := jwt.Parse(params["token"][0], func(token *jwt.Token) (interface{}, error) {
    if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
      return nil, fmt.Errorf("There was an error")
    }
    return []byte("secret"), nil
  })
  if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
    var users Users
    mapstructure.Decode(claims, &users)
    json.NewEncoder(w).Encode(users)
  } else {
    json.NewEncoder(w).Encode(Exception{Message: "Invalid authorization token"})
  }
}

func ValidateMiddleware(next http.HandlerFunc) http.HandlerFunc {
  return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
    authorizationHeader := req.Header.Get("authorization")
    if authorizationHeader != "" {
      bearerToken := strings.Split(authorizationHeader, " ")
      if len(bearerToken) == 2 {
        token, error := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface {}, error) {
          if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("There was an error")
          }
          return []byte("secret"), nil
        })
        if error != nil {
          json.NewEncoder(w).Encode(Exception{Message: error.Error()})
          return
        }
        if token.Valid {
          context.Set(req, "decoded", token.Claims)
          next(w, req)
        } else {
          json.NewEncoder(w).Encode(Exception{Message: "Invalid authorization token"})
        }
      }
    } else {
      json.NewEncoder(w).Encode(Exception{Message: "An authorization header is required"})
    }
  })
}

func TestEndpoint(w http.ResponseWriter, req *http.Request) {
  decoded := context.Get(req, "decoded")
  var users Users
  mapstructure.Decode(decoded.(jwt.MapClaims), &users)
  json.NewEncoder(w).Encode(users)
}

func getAllUsers(w http.ResponseWriter, r *http.Request) {
  var users Users
  var arr_user []Users
  var response ResUsers

  db := connect()
  defer db.Close()

  rows, err := db.Query("Select id, username, location from users")
  if err != nil {
    log.Print(err)
  }

  for rows.Next() {
    if err := rows.Scan(&users.Id, &users.Username, &users.Location); err != nil {
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
  router.HandleFunc("/authenticate", CreateTokenEndpoint).Methods("POST")
  router.HandleFunc("/protected", ProtectedEndpoint).Methods("GET")
  router.HandleFunc("/blogs", ValidateMiddleware(getAllBlogs)).Methods("GET")
  router.HandleFunc("/users", ValidateMiddleware(getAllUsers)).Methods("GET")
  router.HandleFunc("/test", ValidateMiddleware(TestEndpoint)).Methods("GET")
  http.Handle("/", router)
  fmt.Println("Connected to port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
