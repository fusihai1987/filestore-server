package handler

import (
	"net/http"
)

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet{
		//signUpHtml, err := ioutil.ReadFile("./static/views/sign.html")

		//if err != nil {
		//	io.WriteString(w, "File store server internal error!")
		//	return
		//}

		//io.WriteString(w, string(signUpHtml))
		//return
		http.Redirect(w, r, "/static/view/signup.html", http.StatusFound)
		return
	}
	r.ParseForm()

	userName := r.Form.Get("user_name");



}