package public

import (
  "log"
  "net/url"
  "net/http"
  "crypto/rand"
  "database/sql"
  "encoding/hex"
  "html/template"
  _ "github.com/go-sql-driver/mysql"
)

var DBaddr, Domain string

func HomePage(w http.ResponseWriter, r *http.Request) {
  index, err := template.ParseFiles("templates/index.html")
  CheckErr(err)
  if r.Method == "GET" {
    index.Execute(w, nil)
    return
  }else if r.Method == "POST" {
    r.ParseForm()
    inputurl := r.PostForm.Get("url")
    if URLValidate(inputurl) == false {
      index.Execute(w, "Invalid URL.")
      return
    }
    if len(inputurl) > 300 {
      index.Execute(w, "Error, URL too long!")
      return
    }
    path := RandomPath()
    status := Insert(inputurl, path)
    if status != 1 {
      index.Execute(w, "Internal Server Error")
      return
    }
    index.Execute(w, "http://" + Domain + "/u/" + path)
  }
}

func GetURL(w http.ResponseWriter, r *http.Request) {
  if r.Method != "GET" {
    http.Error(w, "Method Not Allowed", 405)
    return
  }
  path := URLParse(r.URL.Path)
  if path == "" {
    http.NotFound(w, r)
    return
  }
  ur := Query(path)
  if ur == "" {
    http.NotFound(w, r)
    return
  }
  http.Redirect(w, r, ur, 302)
}

func Insert(ur, path string) (row int64) {
  DB, err := sql.Open("mysql", DBaddr)
  CheckErr(err)
  ins, err := DB.Prepare(`INSERT INTO shorturl(url, path) VALUES(?, ?)`)
  CheckErr(err)
  res, err := ins.Exec(ur, path)
  CheckErr(err)
  row, err = res.RowsAffected()
  CheckErr(err)
  return row
}

func Query(path string) (uri string) {
  DB, err := sql.Open("mysql", DBaddr)
  CheckErr(err)
  q, err := DB.Query(`SELECT url FROM shorturl WHERE path=?`, path)
  CheckErr(err)
  if q.Next() == true {
    q.Scan(&uri)
  }else if q.Next() == false {
    uri = ""
  }
  return uri
}

func RandomPath() (encoded string) {
  a := make([]byte, 5)
  rand.Read(a)
  encoded = hex.EncodeToString(a)
  return encoded
}

func URLParse(path string) string {
  if len(path) <= 3 {
    return ""
  }
	if string(path[len(path) - 1]) == "/" {
		return string(path[3:len(path) - 1])
	}else {
		return string(path[3:])
	}
}

func URLValidate(uri string) bool {
    u, err := url.Parse(uri)
    return err == nil && u.Scheme != "" && u.Host != ""
}


func CheckErr(err error) {
	if err != nil {
    log.Fatal(err.Error())
    return
	}
}
