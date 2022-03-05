package routes

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "embed"

	"github.com/thanhpk/randstr"
)

type AdminLoginParams struct {
	LoginError  bool
	PrefillName string
}

func StartSession() {

}

func AdminLogin(w http.ResponseWriter, r *http.Request) {
	// Serve the form and end
	if r.Method == http.MethodGet {
		err := t.ExecuteTemplate(w, "admin-login.html", &AdminLoginParams{})
		if err != nil {
			w.WriteHeader(500)
			log.Println("error with admin page: " + err.Error())
		}
		return
	}

	r.ParseForm()

	var row struct {
		Id         string
		Nombre     string
		Contraseña string
	}
	err := db.QueryRowx("SELECT id, nombre, contraseña FROM autores WHERE id=?;", r.PostFormValue("userid")).StructScan(&row)
	if err == sql.ErrNoRows {
		//TODO log failed login
		t.ExecuteTemplate(w, "admin-login.html", &AdminLoginParams{
			PrefillName: r.PostFormValue("userid"),
			LoginError:  true,
		})
	}

	pass := ValidatePassword(r.PostFormValue("password"), row.Contraseña)

	if !pass {
		t.ExecuteTemplate(w, "admin-login.html", &AdminLoginParams{
			PrefillName: r.PostFormValue("userid"),
			LoginError:  true,
		})
		return
	}

	token := randstr.String(20)

	http.SetCookie(w, &http.Cookie{
		Name:     "sess",
		Value:    token,
		Path:     "/admin",
		Domain:   r.URL.Host,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
	})

	fmt.Fprintf(w, "id => %s\nnombre => %s\ncontraseña => %t\n", row.Id, row.Nombre, pass)
}
