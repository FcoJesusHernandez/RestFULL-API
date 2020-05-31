package main

import (
	"container/list"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Materia struct {
	Id     uint64
	Nombre string
}

type Alumno struct {
	Id     uint64
	Nombre string
}

type Calificacion struct {
	Id           uint64
	Alumno       Alumno
	Materia      Materia
	Calificacion float64
}

type Calificaciones struct {
	Calificaciones list.List
}

var lista_calificaciones = Calificaciones{}
var lista_alumnos list.List
var lista_materias list.List

var cont_alumnos uint64 = 0
var cont_materias uint64 = 0
var cont_calificaciones uint64 = 0

func (this *Calificaciones) Evaluar(datos []string, respuesta *string, danger *bool) error {
	alumno_ := Alumno{
		Id:     cont_alumnos,
		Nombre: datos[0],
	}

	materia_ := Materia{
		Id:     cont_materias,
		Nombre: datos[1],
	}

	num, err := strconv.ParseFloat(datos[2], 64)
	if err == nil {
		evaluacion := Calificacion{
			Id:           cont_calificaciones,
			Alumno:       alumno_,
			Materia:      materia_,
			Calificacion: num,
		}

		var bandera = false
		for e := lista_alumnos.Front(); e != nil; e = e.Next() {
			if e.Value.(Alumno).Nombre == alumno_.Nombre {
				bandera = true
			}
		}

		if !bandera {
			lista_alumnos.PushBack(alumno_)
			cont_alumnos++
		}

		bandera = false
		for e := lista_materias.Front(); e != nil; e = e.Next() {
			if e.Value.(Materia).Nombre == materia_.Nombre {
				bandera = true
			}
		}

		if !bandera {
			lista_materias.PushBack(materia_)
			cont_materias++
		}

		bandera = false
		for e := lista_calificaciones.Calificaciones.Front(); e != nil; e = e.Next() {
			if e.Value.(Calificacion).Alumno.Nombre == evaluacion.Alumno.Nombre && e.Value.(Calificacion).Materia.Nombre == evaluacion.Materia.Nombre {
				bandera = true
			}
		}

		if !bandera {
			lista_calificaciones.Calificaciones.PushBack(evaluacion)
			cont_calificaciones++
			*respuesta = "Evaluación anexada con éxito"
			*danger = false
		} else {
			*respuesta = "Error, Evaluación ya existente"
			*danger = true
		}
	} else {
		*respuesta = "Error, Evaluación debe ser numerica"
		*danger = true
	}

	return nil
}

func (this *Calificaciones) Promedio(datos []string, respuesta *float64, danger *bool) error {
	var total float64
	var promedio float64

	tipo := datos[0]
	auxiliar := datos[1]

	if tipo == "1" { // promedio de alumno
		alumno_ := Alumno{
			Id:     12,
			Nombre: auxiliar,
		}

		for e := lista_calificaciones.Calificaciones.Front(); e != nil; e = e.Next() {
			if e.Value.(Calificacion).Alumno.Nombre == alumno_.Nombre {
				total += 1
				promedio += e.Value.(Calificacion).Calificacion
			}
		}
		if total == 0 {
			*respuesta = 0
			*danger = true
		} else {
			*respuesta = promedio / total
			*danger = false
		}
	} else if tipo == "2" { // promedio general / todos
		for e := lista_calificaciones.Calificaciones.Front(); e != nil; e = e.Next() {
			total += 1
			promedio += e.Value.(Calificacion).Calificacion
		}
		if total == 0 {
			*respuesta = 0
			*danger = true
		} else {
			*respuesta = promedio / total
			*danger = false
		}
	} else if tipo == "3" { // promedio de materia
		materia_ := Materia{
			Id:     12,
			Nombre: auxiliar,
		}

		for e := lista_calificaciones.Calificaciones.Front(); e != nil; e = e.Next() {
			if e.Value.(Calificacion).Materia.Nombre == materia_.Nombre {
				total += 1
				promedio += e.Value.(Calificacion).Calificacion
			}
		}
		if total == 0 {
			*respuesta = 0
			*danger = true
		} else {
			*respuesta = promedio / total
			*danger = false
		}
	} else {
		*respuesta = 0.0
		*danger = true
	}
	return nil
}

func root(res http.ResponseWriter, req *http.Request) {
	res.Header().Set(
		"Content-Type",
		"text/html",
	)
	fmt.Fprintf(res,
		cargarHtml("index.html", "inicio", "Hola Bienvenido", "", false, false),
		cargaAlumnosHTML(),
		cargaMateriasHTML())
}

var clf = new(Calificaciones)

func calificacion(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST":
		if err := req.ParseForm(); err != nil {
			fmt.Fprintf(res, "ParseForm() error %v", err)
			return
		}
		fmt.Println(req.PostForm)
		datos := []string{req.FormValue("alumno"), req.FormValue("materia"), req.FormValue("calificacion")}
		var result string
		var danger bool

		clf.Evaluar(datos, &result, &danger)

		res.Header().Set(
			"Content-Type",
			"text/html",
		)

		fmt.Fprintf(
			res,
			cargarHtml("index.html", "inicio", "Hola, Bienvenido", result, danger, false),
			cargaAlumnosHTML(),
			cargaMateriasHTML(),
		)
	}
}

func promedio(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		res.Header().Set(
			"Content-Type",
			"text/html",
		)

		fmt.Fprintf(
			res,
			cargarHtml("index.html", "agregar calificación", "Hola, Bienvenido", "Captura una calificación", false, false),
			cargaAlumnosHTML(),
			cargaMateriasHTML(),
		)
	case "POST":
		if err := req.ParseForm(); err != nil {
			fmt.Fprintf(res, "ParseForm() error %v", err)
			return
		}
		fmt.Println(req.PostForm)
		datos := []string{}
		var salida string
		var result float64
		var danger bool

		if req.FormValue("alumno") != "" {
			salida = "El promedio de " + req.FormValue("alumno") + " es : "
			datos = []string{"1", req.FormValue("alumno")}
		} else if req.FormValue("materia") != "" {
			salida = "El promedio de " + req.FormValue("materia") + " es : "
			datos = []string{"3", req.FormValue("materia")}
		} else {
			salida = "Petición desconocida"
			datos = []string{"4", "error"}
			result = 0.0
			danger = true
		}

		clf.Promedio(datos, &result, &danger)

		res.Header().Set(
			"Content-Type",
			"text/html",
		)

		fmt.Fprintf(
			res,
			cargarHtml("index.html", "promedio", "Hola, Bienvenido", salida+fmt.Sprintf("%f", result), danger, false),
			cargaAlumnosHTML(),
			cargaMateriasHTML(),
		)
	}
}

func promedio_gen(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		if err := req.ParseForm(); err != nil {
			fmt.Fprintf(res, "ParseForm() error %v", err)
			return
		}
		fmt.Println(req.PostForm)
		datos := []string{}
		var salida string
		var result float64
		var danger bool

		salida = "El promedio general es : "
		datos = []string{"2", ""}

		clf.Promedio(datos, &result, &danger)

		res.Header().Set(
			"Content-Type",
			"text/html",
		)

		fmt.Fprintf(
			res,
			cargarHtml("index.html", "promedio general", "Hola, Bienvenido", salida+fmt.Sprintf("%f", result), danger, false),
			cargaAlumnosHTML(),
			cargaMateriasHTML(),
		)
	}
}

func cargaAlumnosHTML() string {
	var html string
	html += "<option value='null'>Selecciona una opción</option>"
	for e := lista_alumnos.Front(); e != nil; e = e.Next() {
		html += "<option value='" + e.Value.(Alumno).Nombre + "'>" + e.Value.(Alumno).Nombre + "</option>"
	}

	return html
}

func cargaMateriasHTML() string {
	var html string
	html += "<option value='null'>Selecciona una opción</option>"
	for e := lista_materias.Front(); e != nil; e = e.Next() {
		html += "<option value='" + e.Value.(Materia).Nombre + "'>" + e.Value.(Materia).Nombre + "</option>"
	}

	return html
}

func cargarHtml(a string, titulo string, mensaje string, auxiliar string, danger_aux bool, danger_msj bool) string {
	html, _ := ioutil.ReadFile(a)
	salida := strings.Replace(string(html), "$__TITULO__$", titulo, -1)
	salida = strings.Replace(salida, "$__MENSAJE__$", mensaje, -1)
	salida = strings.Replace(salida, "$__AUXILIAR__$", auxiliar, -1)

	if danger_aux {
		salida = strings.Replace(salida, "$__CLASS_AUX__$", "alert-danger", -1)
	} else {
		salida = strings.Replace(salida, "$__CLASS_AUX__$", "alert-success", -1)
	}

	if danger_msj {
		salida = strings.Replace(salida, "$__CLASS_MSJ__$", "alert-danger", -1)
	} else {
		salida = strings.Replace(salida, "$__CLASS_MSJ__$", "alert-secondary", -1)
	}
	return salida
}

func respaldo(res http.ResponseWriter, req *http.Request) {
	outFile, err := os.Create("calificaciones.json")
	if err != nil {
		fmt.Println("Error al convertir a JSON", err.Error())
		return
	}
	var temp_cal []Calificacion

	for e := lista_calificaciones.Calificaciones.Front(); e != nil; e = e.Next() {
		temp_cal = append(temp_cal, e.Value.(Calificacion))
	}

	err = json.NewEncoder(outFile).Encode(temp_cal)
	if err != nil {
		fmt.Println("Error al convertir a JSON", err.Error())
		return
	}
	outFile.Close()

	res.Header().Set(
		"Content-Type",
		"text/html",
	)

	fmt.Fprintf(
		res,
		cargarHtml("index.html", "Respaldo exitoso ", "Información respaldada con éxito ", "Respaldamos toda la información almacenada hasta el momento", true, true),
		cargaAlumnosHTML(),
		cargaMateriasHTML(),
	)
}

func recuperacion(res http.ResponseWriter, req *http.Request) {
	inFile, err := os.Open("calificaciones.json")
	if err != nil {
		fmt.Println("Error al abrir el archivo", err.Error())
		return
	}

	var temp_cal []Calificacion

	err = json.NewDecoder(inFile).Decode(&temp_cal)

	danger := false
	if err != nil {
		danger = true
		fmt.Println("Error de conversión", err.Error())
		return
	}

	for i := 0; i < len(temp_cal); i++ {
		var bandera = false
		for e := lista_alumnos.Front(); e != nil; e = e.Next() {
			if e.Value.(Alumno).Nombre == temp_cal[i].Alumno.Nombre {
				bandera = true
			}
		}

		if !bandera {
			lista_alumnos.PushBack(temp_cal[i].Alumno)
		}

		bandera = false
		for e := lista_materias.Front(); e != nil; e = e.Next() {
			if e.Value.(Materia).Nombre == temp_cal[i].Materia.Nombre {
				bandera = true
			}
		}

		if !bandera {
			lista_materias.PushBack(temp_cal[i].Materia)
		}

		bandera = false
		for e := lista_calificaciones.Calificaciones.Front(); e != nil; e = e.Next() {
			if e.Value.(Calificacion).Alumno.Nombre == temp_cal[i].Alumno.Nombre && e.Value.(Calificacion).Materia.Nombre == temp_cal[i].Materia.Nombre {
				bandera = true
			}
		}

		if !bandera {
			lista_calificaciones.Calificaciones.PushBack(temp_cal[i])
		}
	}

	fmt.Println(lista_calificaciones.Calificaciones)

	inFile.Close()

	res.Header().Set(
		"Content-Type",
		"text/html",
	)

	fmt.Fprintf(
		res,
		cargarHtml("index.html", "Recuperación exitosa ", "Información Restaurada con éxito ", "Restauramos toda la información almacenada hasta el momento", danger, danger),
		cargaAlumnosHTML(),
		cargaMateriasHTML(),
	)
}

func restauracion(res http.ResponseWriter, req *http.Request) {
	lista_calificaciones.Calificaciones.Init()
	lista_alumnos.Init()
	lista_materias.Init()

	res.Header().Set(
		"Content-Type",
		"text/html",
	)

	fmt.Fprintf(
		res,
		cargarHtml("index.html", "Vaciado exitoso ", "Información eliminada con éxito ", "Eliminamos toda la información almacenada hasta el momento", true, true),
		cargaAlumnosHTML(),
		cargaMateriasHTML(),
	)
}

func GetID(id uint64) ([]byte, error) {
	jsonData := []byte(`{}`)

	var temp_cal []Calificacion

	for e := lista_calificaciones.Calificaciones.Front(); e != nil; e = e.Next() {
		if e.Value.(Calificacion).Id == id {
			temp_cal = append(temp_cal, e.Value.(Calificacion))
		}
	}

	jsonData, err := json.MarshalIndent(temp_cal, "", "    ")
	if err != nil {
		return jsonData, err
	}
	return jsonData, nil
}

func Delete(id uint64) []byte {
	bandera := false

	for e := lista_calificaciones.Calificaciones.Front(); e != nil; e = e.Next() {
		if e.Value.(Calificacion).Id == id {
			lista_calificaciones.Calificaciones.Remove(e)
			bandera = true
			break
		}
	}

	if bandera == false {
		return []byte(`{"code": "noexiste"}`)
	}

	return []byte(`{"code": "ok"}`)
}

func Update(id uint64, calificacion Calificacion) []byte {
	bandera := false
	for e := lista_calificaciones.Calificaciones.Front(); e != nil; e = e.Next() {
		if e.Value.(Calificacion).Id == id {
			lista_calificaciones.Calificaciones.Remove(e)
			lista_calificaciones.Calificaciones.PushFront(calificacion)
			bandera = true
		}
	}

	if bandera == false {
		return []byte(`{"code": "noexiste"}`)
	} else {
		return []byte(`{"code": "ok"}`)
	}
}

func api_id(res http.ResponseWriter, req *http.Request) {
	id, err := strconv.ParseUint(strings.TrimPrefix(req.URL.Path, "/alumno/"), 10, 64)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println(req.Method)
	switch req.Method {
	case "GET":
		res_json, err := GetID(id)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
		res.Header().Set(
			"Content-Type",
			"application/json",
		)
		res.Write(res_json)
	case "DELETE":
		res_json := Delete(id)
		res.Header().Set(
			"Content-Type",
			"application/json",
		)
		res.Write(res_json)
	case "PUT":
		var calificacion Calificacion
		err := json.NewDecoder(req.Body).Decode(&calificacion)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		res_json := Update(id, calificacion)
		res.Header().Set(
			"Content-Type",
			"application/json",
		)
		res.Write(res_json)
	}
}

func Get() ([]byte, error) {
	var temp_cal []Calificacion

	for e := lista_calificaciones.Calificaciones.Front(); e != nil; e = e.Next() {
		temp_cal = append(temp_cal, e.Value.(Calificacion))
	}

	jsonData, err := json.MarshalIndent(temp_cal, "", "    ")
	if err != nil {
		return jsonData, nil
	}
	return jsonData, err
}

func Add(temp_cal Calificacion) []byte {
	var bandera = false
	for e := lista_alumnos.Front(); e != nil; e = e.Next() {
		if e.Value.(Alumno).Nombre == temp_cal.Alumno.Nombre {
			bandera = true
		}
	}

	if !bandera {
		lista_alumnos.PushBack(temp_cal.Alumno)
	}

	bandera = false
	for e := lista_materias.Front(); e != nil; e = e.Next() {
		if e.Value.(Materia).Nombre == temp_cal.Materia.Nombre {
			bandera = true
		}
	}

	if !bandera {
		lista_materias.PushBack(temp_cal.Materia)
	}

	bandera = false
	for e := lista_calificaciones.Calificaciones.Front(); e != nil; e = e.Next() {
		if e.Value.(Calificacion).Alumno.Nombre == temp_cal.Alumno.Nombre && e.Value.(Calificacion).Materia.Nombre == temp_cal.Materia.Nombre {
			bandera = true
		}
	}

	if !bandera {
		lista_calificaciones.Calificaciones.PushBack(temp_cal)
		return []byte(`{"code": "ok"}`)
	} else {
		return []byte(`{"code": "error, calificacion antes capturada"}`)
	}
}

func api(res http.ResponseWriter, req *http.Request) {
	fmt.Println(req.Method)
	switch req.Method {
	case "POST":
		var calificacion Calificacion
		err := json.NewDecoder(req.Body).Decode(&calificacion)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Println(calificacion)
		res_json := Add(calificacion)
		res.Header().Set(
			"Content-Type",
			"application/json",
		)
		res.Write(res_json)
	case "GET":
		res_json, err := Get()
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		res.Header().Set(
			"Content-Type",
			"application/json",
		)
		res.Write(res_json)
	}
}

func main() {
	http.HandleFunc("/calificacion", calificacion)
	http.HandleFunc("/promedio", promedio)
	http.HandleFunc("/general", promedio_gen)
	http.HandleFunc("/inicio", root)
	http.HandleFunc("/respaldo", respaldo)
	http.HandleFunc("/recuperacion", recuperacion)
	http.HandleFunc("/restauracion", restauracion)
	http.HandleFunc("/api", api)
	http.HandleFunc("/alumno/", api_id)
	fmt.Println("Arrancando el servidor...")
	http.ListenAndServe(":9000", nil)
}
