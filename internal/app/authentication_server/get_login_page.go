package authentication_server

import (
	"net/http"
)

func (i *Implementation) GetLoginPage(w http.ResponseWriter, r *http.Request) {
	//todo: мб добавить логику с проверкой, пришёл ли jwt
	i.fileServer.ServeHTTP(w, r)
}
