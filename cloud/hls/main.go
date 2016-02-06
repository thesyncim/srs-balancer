/*
package main

import("net/http")



func main()  {
	http.HandleFunc("/crossdomain.xml",func(w http.ResponseWriter, r *http.Request){
		w.Write([]byte(`<cross-domain-policy>
  <allow-access-from domain="*"/>
</cross-domain-policy>`))



	})
	http.HandleFunc("/livestream.m3u8",func(w http.ResponseWriter, r *http.Request){
		http.Redirect(w, r, "http://149.202.163.159/live/livestream.m3u8", 302)

	})
	panic(http.ListenAndServe(":9999",nil))


}*/
package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"log"

)

func cross(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, `<cross-domain-policy>
  <allow-access-from domain="*"/>
</cross-domain-policy>`)
}

func Hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if ps.ByName("name") == "crossdomain.xml" {
		fmt.Fprint(w, `<cross-domain-policy>
        <allow-access-from domain="*"/>
        </cross-domain-policy>`)
	}else {



		s := "http://149.202.163.159/live/%s"
		http.Redirect(w, r, fmt.Sprintf(s, ps.ByName("name")), 302)
	}
}

func main() {
	router := httprouter.New()

	router.GET("/:name", Hello)

	log.Fatal(http.ListenAndServe(":9999", router))
}