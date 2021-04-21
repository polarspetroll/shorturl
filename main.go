package main

import (
	"fmt"
	"net/http"
	"os"
	"public"
)

func main() {
	public.Domain = os.Getenv("DOMAIN")
	public.DBaddr = fmt.Sprintf("%v:%v@tcp(%v:3306)/%v", os.Getenv("DBUSR"), os.Getenv("DBPWD"), os.Getenv("DBADDR"), os.Getenv("DBNAME"))
	port := fmt.Sprintf(":" + os.Getenv("LPORT"))

	file := http.FileServer(http.Dir("./static"))
	http.Handle("/statics/", http.StripPrefix("/statics/", file))
	http.HandleFunc("/", RootRedirect)
	http.HandleFunc("/index/", public.HomePage)
	http.HandleFunc("/u/", public.GetURL)
	err := http.ListenAndServe(port, nil)
	public.CheckErr(err)
}

func RootRedirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/index/", 302)
}
