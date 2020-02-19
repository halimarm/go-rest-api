package main

type Users struct {
  Username  string `form:"username" json:"username"`
  Password  string `form:"password" json:"password"`
}

type JwtToken struct {
  Token     string `json:"token"`
}

type Exception struct {
  Message   string `json:"message"`
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
