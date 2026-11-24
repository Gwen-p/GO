// =====================================================================
//  PROYECTO: ACCESO A DATOS DESDE GO  --  Programación de Bases de Datos
//  Bloque B. Acceso a datos desde Go: Oracle Database XE (SID xe)
//
//  Funciones: dbConectar, dbDesconectar, verUsuarios, verServidores,
//  dbObtenerUsuario, dbInsertarServidor, dbModificarServidor,
//  dbEliminarServidor, dbInsertarVariosServidores y
//  dbBuscarUsuarioConcatenado (solo para demostrar la inyección SQL).
//
//  Ejecución (desde la raíz del módulo):  go run ./codigo/03_oracle
// =====================================================================

package main

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	go_ora "github.com/sijms/go-ora/v2" // registra el driver "oracle"
)

// conexion es la variable global con el acceso a la base de datos, igual
// que en JDBC y en Python. En Go, *sql.DB no es una única conexión sino
// un POOL de conexiones que el paquete abre y reutiliza según necesite.
var conexion *sql.DB

// Servidor representa una tupla de la tabla Servidores.
type Servidor struct {
	Cod      int
	Nombre   string
	Uso      float64
	Megabits float64
}

// =====================================================================
// dbConectar: abre la conexión con Oracle XE y la comprueba.
// Devuelve true si la conexión se ha establecido.
// =====================================================================
func dbConectar() bool {
	fmt.Println("---dbConectar---")
	// BuildUrl compone la URL oracle://... y escapa la contraseña.
	// La instancia de la asignatura se identifica por su SID: xe.
	url := go_ora.BuildUrl("localhost", 1521, "", "system", "12345",
		map[string]string{"SID": "xe"})

	var err error
	conexion, err = sql.Open("oracle", url) // prepara el pool, NO conecta
	if err != nil {
		fmt.Println("Error en sql.Open:", err)
		return false
	}
	if err = conexion.Ping(); err != nil { // aquí se conecta de verdad
		fmt.Println("Error de conexión:", err)
		return false
	}
	conexion.SetMaxOpenConns(5)                   // máximo de conexiones abiertas
	conexion.SetConnMaxLifetime(30 * time.Minute) // renovar conexiones antiguas

	var version string
	conexion.QueryRow("SELECT banner FROM v$version WHERE ROWNUM = 1").Scan(&version)
	fmt.Println("Conectado a:", version)
	return true
}

// =====================================================================
// dbDesconectar: cierra el pool y todas sus conexiones.
// =====================================================================
func dbDesconectar() bool {
	fmt.Println("---dbDesconectar---")
	if err := conexion.Close(); err != nil {
		fmt.Println("Error al desconectar:", err)
		return false
	}
	fmt.Println("Desconexión realizada")
	return true
}

// =====================================================================
// verUsuarios: consulta y muestra TODAS las tuplas de Usuarios.
// =====================================================================
func verUsuarios() {
	fmt.Println("---verUsuarios---")
	rows, err := conexion.Query(
		"SELECT cod, nombre, fechaalta, megabits, servidor FROM Usuarios ORDER BY cod")
	if err != nil {
		fmt.Println("Error en la consulta:", err)
		return
	}
	defer rows.Close() // libera el cursor al salir, pase lo que pase

	n := 0
	for rows.Next() { // avanza a la siguiente fila; false al terminar
		var cod, nombre string
		var alta time.Time
		var megabits float64
		var servidor sql.NullInt64 // la columna admite NULL
		if err := rows.Scan(&cod, &nombre, &alta, &megabits, &servidor); err != nil {
			fmt.Println("Error en Scan:", err)
			return
		}
		srv := "NULL"
		if servidor.Valid {
			srv = fmt.Sprint(servidor.Int64)
		}
		fmt.Printf("  %-3s %-10s %s %7.2f  servidor %s\n",
			cod, nombre, alta.Format("2006-01-02"), megabits, srv)
		n++
	}
	if err := rows.Err(); err != nil { // error ocurrido durante el recorrido
		fmt.Println("Error al recorrer:", err)
	}
	fmt.Println("Usuarios recuperados:", n)
}

// =====================================================================
// verServidores: consulta y muestra TODAS las tuplas de Servidores.
// =====================================================================
func verServidores() {
	fmt.Println("---verServidores---")
	rows, err := conexion.Query("SELECT codS, nombreS, usoS, megabitsS FROM Servidores ORDER BY codS")
	if err != nil {
		fmt.Println("Error en la consulta:", err)
		return
	}
	defer rows.Close()

	var lista []Servidor // las tuplas se guardan en un slice de structs
	for rows.Next() {
		var s Servidor
		if err := rows.Scan(&s.Cod, &s.Nombre, &s.Uso, &s.Megabits); err != nil {
			fmt.Println("Error en Scan:", err)
			return
		}
		lista = append(lista, s)
	}
	for _, s := range lista {
		fmt.Printf("  %d %-10s uso=%.2f megabits=%.2f\n", s.Cod, s.Nombre, s.Uso, s.Megabits)
	}
	fmt.Println("Servidores recuperados:", len(lista))
}

// =====================================================================
// dbObtenerUsuario: consulta UNA tupla por clave primaria con un
// parámetro vinculado (:1). El valor nunca se mezcla con el texto SQL.
// =====================================================================
func dbObtenerUsuario(cod string) {
	fmt.Printf("---dbObtenerUsuario(%q)---\n", cod)
	var nombre string
	var megabits float64
	err := conexion.QueryRow(
		"SELECT nombre, megabits FROM Usuarios WHERE cod = :1", cod).Scan(&nombre, &megabits)
	switch {
	case errors.Is(err, sql.ErrNoRows): // la consulta no devolvió ninguna fila
		fmt.Println("No existe ningún usuario con ese código")
	case err != nil:
		fmt.Println("Error en la consulta:", err)
	default:
		fmt.Printf("Usuario %s: %s, %.2f megabits\n", cod, nombre, megabits)
	}
}

// =====================================================================
// dbInsertarServidor: inserta un servidor con uso y megabits a 0.
// =====================================================================
func dbInsertarServidor(cod int, nombre string) {
	fmt.Printf("---dbInsertarServidor(%d, %q)---\n", cod, nombre)
	res, err := conexion.Exec("INSERT INTO Servidores VALUES (:1, :2, 0, 0)", cod, nombre)
	if err != nil {
		fmt.Println("Error al insertar:", err)
		return
	}
	n, _ := res.RowsAffected() // número de filas afectadas
	fmt.Println("Filas insertadas:", n)
}

// =====================================================================
// dbModificarServidor: cambia el nombre de un servidor.
// =====================================================================
func dbModificarServidor(cod int, nombre string) {
	fmt.Printf("---dbModificarServidor(%d, %q)---\n", cod, nombre)
	res, err := conexion.Exec("UPDATE Servidores SET nombreS = :1 WHERE codS = :2", nombre, cod)
	if err != nil {
		fmt.Println("Error al modificar:", err)
		return
	}
	n, _ := res.RowsAffected()
	fmt.Println("Filas modificadas:", n)
}

// =====================================================================
// dbEliminarServidor: borra un servidor. Si no existe, no hay error:
// simplemente RowsAffected devuelve 0.
// =====================================================================
func dbEliminarServidor(cod int) {
	fmt.Printf("---dbEliminarServidor(%d)---\n", cod)
	res, err := conexion.Exec("DELETE FROM Servidores WHERE codS = :1", cod)
	if err != nil {
		fmt.Println("Error al eliminar:", err)
		return
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		fmt.Println("No se ha borrado nada: el servidor no existe")
		return
	}
	fmt.Println("Filas eliminadas:", n)
}

// =====================================================================
// dbInsertarVariosServidores: inserta una lista de servidores en UNA
// transacción. O se guardan todos, o no se guarda ninguno.
// =====================================================================
func dbInsertarVariosServidores(lista []Servidor) {
	fmt.Printf("---dbInsertarVariosServidores(%d servidores)---\n", len(lista))
	tx, err := conexion.Begin() // inicia la transacción
	if err != nil {
		fmt.Println("Error al iniciar la transacción:", err)
		return
	}
	defer tx.Rollback() // deshace todo si salimos antes del Commit

	ps, err := tx.Prepare("INSERT INTO Servidores VALUES (:1, :2, 0, 0)") // se prepara una vez
	if err != nil {
		fmt.Println("Error al preparar:", err)
		return
	}
	defer ps.Close()

	for _, s := range lista {
		if _, err := ps.Exec(s.Cod, s.Nombre); err != nil { // se ejecuta N veces
			fmt.Printf("Error con el servidor %d: %v\n", s.Cod, err)
			fmt.Println("ROLLBACK: no se guarda ninguno")
			return // el defer hace el Rollback
		}
		fmt.Printf("  insertado %d %s (pendiente de COMMIT)\n", s.Cod, s.Nombre)
	}
	if err := tx.Commit(); err != nil { // confirma todos los cambios
		fmt.Println("Error en el COMMIT:", err)
		return
	}
	fmt.Println("COMMIT: servidores guardados:", len(lista))
}

// =====================================================================
// contarServidores: comprobación del estado de la tabla.
// =====================================================================
func contarServidores() {
	var n int
	conexion.QueryRow("SELECT COUNT(*) FROM Servidores").Scan(&n)
	fmt.Println(">> Servidores en la tabla:", n)
}

// =====================================================================
// dbBuscarUsuarioConcatenado: MAL HECHO A PROPÓSITO. Construye el SQL
// concatenando el texto recibido, lo que permite la inyección SQL.
// Solo se usa para compararlo con dbObtenerUsuario.
// =====================================================================
func dbBuscarUsuarioConcatenado(cod string) {
	fmt.Printf("---dbBuscarUsuarioConcatenado(%q)---\n", cod)
	consulta := "SELECT cod FROM Usuarios WHERE cod = '" + cod + "'"
	fmt.Println("SQL enviado:", consulta)
	rows, err := conexion.Query(consulta)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer rows.Close()
	n := 0
	for rows.Next() {
		n++
	}
	fmt.Println("Filas devueltas:", n)
}

// =====================================================================
// PROGRAMA PRINCIPAL
// =====================================================================
func main() {
	fmt.Println("---Programa principal (Oracle XE)---")
	if !dbConectar() {
		os.Exit(1) // sin conexión no se puede continuar
	}
	defer dbDesconectar() // se ejecutará al terminar main

	// 1. Consultas
	verServidores()
	verUsuarios()
	dbObtenerUsuario("123")
	dbObtenerUsuario("ZZZ")

	// 2. Insertar, modificar y borrar un servidor
	dbInsertarServidor(5, "Venus")
	dbInsertarServidor(5, "Venus") // clave duplicada: debe fallar
	dbModificarServidor(5, "Marte")
	verServidores()
	dbEliminarServidor(5)
	dbEliminarServidor(5) // ya no existe: 0 filas
	contarServidores()

	// 3. Operaciones múltiples en una transacción
	dbInsertarVariosServidores([]Servidor{{Cod: 6, Nombre: "Europa"}, {Cod: 7, Nombre: "Io"}})
	contarServidores()
	dbInsertarVariosServidores([]Servidor{{Cod: 8, Nombre: "Titan"}, {Cod: 6, Nombre: "Europa"}})
	contarServidores() // el 8 no se ha guardado
	dbEliminarServidor(6)
	dbEliminarServidor(7)

	// 4. Inyección SQL: misma entrada maliciosa en las dos funciones
	dbBuscarUsuarioConcatenado("' OR '1'='1")
	dbObtenerUsuario("' OR '1'='1")

	// 5. Estado final: deben quedar los 4 servidores originales
	verServidores()
}
