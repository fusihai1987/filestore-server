package handler

import (
	"filestore-server/common"
	"filestore-server/db"
	"filestore-server/utils"
	"fmt"
	"net/http"
	"time"
)

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet{
		http.Redirect(w, r, "/static/views/signup.html", http.StatusFound)
		return
	}
	r.ParseForm()

	userName := r.Form.Get("username")
	password := r.Form.Get("password")

	if len(userName) < 3 || len(password) < 6 {
		w.Write([]byte("Invalid parameter!"))
		return
	}

	encryptPassword := utils.MD5([]byte(password))

	suc := db.SignUp(userName, encryptPassword)

	if suc {
		w.Write([]byte("SUCCESS"))
	}else{
		w.Write([]byte("FAILED"))
	}

}

func SignInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet{
		http.Redirect(w, r, "/static/views/signin.html", http.StatusFound)
		return
	}
	r.ParseForm()

	userName := r.Form.Get("username")
	password := r.Form.Get("password")

	if len(userName) < 3 || len(password) < 6 {
		w.Write([]byte("Invalid parameter!"))
		return
	}

	suc := db.Signin(userName, password)

	if suc {
		resp := common.Resp{
			Code: 0,
			Msg: "OK",
			Data:struct{
				Username string
				Token string
				Location string
			}{
				Username:userName,
				Token:GenToken(userName),
				Location:"http://" + r.Host + "/static/views/home.html",
			},
		}
		w.Write(resp.JsonBytes())
	}else{
		w.Write([]byte("FAILED"))
	}

}

func QueryUserInfo(w http.ResponseWriter, r *http.Request){
	r.ParseForm()

	userName := r.Form.Get("username")

	userInfo, err := db.GetUserInfo(userName)

	if err != nil {
		w.Write([]byte("FAILED"))
		return
	}


	resp := common.NewResp(
			0,
			"SUCCESS",
				 userInfo,
		)
	w.Write(resp.JsonBytes())
}

// GenToken : 生成token
func GenToken(username string) string {
	// 40位字符:md5(username+timestamp+token_salt)+timestamp[:8]
	ts := fmt.Sprintf("%x", time.Now().Unix())
	tokenPrefix := utils.MD5([]byte(username + ts + "_tokensalt"))
	return tokenPrefix + ts[:8]
}

func IsValidToken(token string) bool{
	return len(token) == 40
}