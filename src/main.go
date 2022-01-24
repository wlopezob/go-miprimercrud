package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

type Empleado struct {
	Id    int
	Name  string
	Email string
}

func conexionBD() *sql.DB {
	driver := "postgres"
	user := "miusuario"
	password := "mipass"
	bd := "mibd"
	server := "miserver"

	conectionString := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=require",
		server, user, password, bd)
	fmt.Println(conectionString)
	conexion, err := sql.Open(driver, conectionString)
	if err != nil {
		panic(err.Error())
	}
	return conexion
}

var templates = template.Must(template.ParseGlob("plantillas/*"))

func doNothing(w http.ResponseWriter, r *http.Request) {}
func main() {
	http.HandleFunc("/favicon.ico", doNothing)
	http.HandleFunc("/", Init)
	http.HandleFunc("/add", Add)
	http.HandleFunc("/save", Save)
	http.HandleFunc("/delete", Delete)
	http.HandleFunc("/edit", Edit)
	http.HandleFunc("/update", Update)
	log.Println("Servidor corriendo...")
	http.ListenAndServe(":8080", nil)
}
func Update(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		id := r.FormValue("id")
		name := r.FormValue("name")
		email := r.FormValue("email")
		conexionPq := conexionBD()
		defer conexionPq.Close()
		_, err := conexionPq.Exec("UPDATE empleados SET nombre=$1, correo=$2 WHERE id=$3", name, email, id)
		if err != nil {
			panic(err.Error())
		}
		http.Redirect(w, r, "/", 301)
	}
}
func Edit(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	conexionPq := conexionBD()
	defer conexionPq.Close()
	empleado := Empleado{}
	selectRecord, err := conexionPq.Query("SELECT id,nombre,correo FROM empleados WHERE id=$1", id)
	if err != nil {
		panic(err.Error())
	}
	for selectRecord.Next() {
		var id int
		var name, email string
		selectRecord.Scan(&id, &name, &email)
		empleado.Id = id
		empleado.Name = name
		empleado.Email = email
	}
	templates.ExecuteTemplate(w, "update", empleado)
}
func Delete(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	fmt.Println(id)
	conexionPq := conexionBD()
	defer conexionPq.Close()
	deleteRecord, err := conexionPq.Prepare("DELETE FROM empleados WHERE id=$1")
	if err != nil {
		panic(err.Error())
	}
	deleteRecord.Exec(id)
	http.Redirect(w, r, "/", 301)
}
func Save(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		name := r.FormValue("name")
		email := r.FormValue("email")
		conexionPq := conexionBD()
		defer conexionPq.Close()
		insertRecord, err := conexionPq.Prepare("INSERT INTO empleados(nombre,correo) VALUES($1,$2)")
		if err != nil {
			panic(err.Error())
		}
		//insertamos los registros
		insertRecord.Exec(name, email)
		http.Redirect(w, r, "/", 301)
	}
}
func Init(w http.ResponseWriter, r *http.Request) {
	conexionPq := conexionBD()
	defer conexionPq.Close()
	selectRecord, err := conexionPq.Query("SELECT id,nombre,correo FROM empleados")
	if err != nil {
		panic(err.Error())
	}
	empleado := Empleado{}
	listadoEmpleado := []Empleado{}

	for selectRecord.Next() {
		var id int
		var name, email string
		err := selectRecord.Scan(&id, &name, &email)
		if err != nil {
			panic(err.Error())
		}
		empleado.Id = id
		empleado.Name = name
		empleado.Email = email

		listadoEmpleado = append(listadoEmpleado, empleado)
	}
	templates.ExecuteTemplate(w, "inicio", listadoEmpleado)
}

func Add(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "add", nil)
}
