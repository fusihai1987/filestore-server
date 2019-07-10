package handler

import (
	"FILESTORE-SERVER/common"
	"net/http"
)

func HttpMiddle(h http.HandlerFunc) http.HandlerFunc{
	return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
		r.ParseForm()
		username := r.Form.Get("username")
		token := r.Form.Get("token")

		if len(username)<3 || !IsValidToken(token) {
			resp := common.NewResp(
				int(common.StatusInvalidToken),
				"Invalid Token",
				nil,
				)
			w.Write(resp.JsonBytes())
			return
		}
		h(w, r)
	})
}