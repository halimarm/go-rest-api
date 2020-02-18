package main

type Users struct {
	Id        string `form:"id" json:"id"`
  Username  string `form:"username" json:"username"`
  Email     string `form:"email" json:"email"`
	FirstName string `form:"firstname" json:"firstname"`
	LastName  string `form:"lastname" json:"lastname"`
  Location  string `form:"location" json:"location"`
}

type Blogs struct {
  Id        string `form:"id" json:"id"`
  Title     string `form:"title" json:"title"`
  Slug      string `form:"slug" json:"slug"`
  Content   string `form:"content" json:"content"`
}

type ResUsers struct {
	Status  bool    `json:"status"`
	Message string `json:"message"`
  Data    []Users `json:"results"`
}

type ResBlogs struct {
	Status  bool    `json:"status"`
	Message string `json:"message"`
  Data    []Blogs `json:"results"`
}
