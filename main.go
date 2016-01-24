package main

import (
	"fmt"
	"net/http"
	"html/template"
	"gopkg.in/mgo.v2"
	"math/rand"
	"time"
)

type Url struct {
	Id string `json:"id" bson:"_id,omitempty"`
	Url string `json:"url" bson:"url"`
	Slug string `json:"slug" bson:"slug"`
	Clicks int `json:"clicks" bson:"clicks"`
}

func main() {
	http.HandleFunc("/", Homepage)
	http.HandleFunc("/add", addLinkPage)
	http.ListenAndServe(":8080", nil)
}

func Homepage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, homepageContent)
}

const homepageContent = `
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>Super Duper URL Shortener!</title>
    <link href="http://getbootstrap.com/dist/css/bootstrap.min.css" rel="stylesheet">
  </head>

  <body>
    <div class="container">
      <form action="/add" method="post" class="col-md-4 col-md-offset-4">
        <h4>Add a URL to be shortened:</h4>
		<input name="url" class="form-control" required autofocus>
		<button class="btn btn-primary btn-block" type="submit">Shorten</button>
      </form>
    </div>
  </body>
</html>
`

var addPageTemplate = template.Must(template.New("addLinkPage").Parse(addLinkPageContent))

func addLinkPage(w http.ResponseWriter, r *http.Request) {

	session, err := mgo.Dial("mongodb://url:urldb@ds049925.mongolab.com:49925/urls")
	if err != nil {
		panic(err)
	}

	defer session.Close()

	session.SetSafe(&mgo.Safe{})

	Collection := session.DB("urls").C("urls")

	url := &Url{
		Url: r.FormValue("url"),
		Slug: RandomString(12),
		Clicks: 0,
	}

	Collection.Insert(url)
	if err != nil {
		panic(err)
	}

	addPageTemplate.Execute(w, "Your link is: http://localhost/l/" + r.FormValue("url"))
}

const addLinkPageContent = `
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>Super Duper URL Shortener!</title>
    <link href="http://getbootstrap.com/dist/css/bootstrap.min.css" rel="stylesheet">
  </head>

  <body>
    <div class="container">
      <div class="col-md-4 col-md-offset-4">
        <h4>New Link</h4>
        <p>{{html .}}</p>
      </div>
    </div>
  </body>
</html>
`

func RandomString(strlen int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}