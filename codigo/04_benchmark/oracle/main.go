// =====================================================================
//  PROYECTO: ACCESO A DATOS DESDE GO  --  Programación de Bases de Datos
//  Comparativa JDBC / Python / Go: pruebas en Go sobre Oracle XE
//
//  Ejecución (desde la raíz del módulo):  go run ./codigo/04_benchmark/oracle
// =====================================================================

package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"runtime"
	"sort"
	"sync"
	"time"

	go_ora "github.com/sijms/go-ora/v2" // registra el driver "oracle"
)

// Parámetros de las pruebas: IGUALES en los seis programas
const (
	REPS        = 5     // repeticiones medidas (además de 1 de calentamiento)
	N_CONEXION  = 10    // conexiones abiertas y cerradas por repetición
	N_PK        = 2000  // consultas por clave primaria
	N_TX        = 2000  // inserciones dentro de una transacción
	N_AUTO      = 500   // inserciones con autocommit (un commit cada una)
	N_LOTE      = 20000 // filas de la carga masiva
	CONC_HILOS  = 4     // trabajadores concurrentes
	CONC_CONSUL = 500   // consultas por trabajador
)

// La instancia de la asignatura se identifica por su SID: xe
var url = go_ora.BuildUrl("localhost", 1521, "", "system", "12345",
	map[string]string{"SID": "xe"})

const INSERT = "INSERT INTO BenchPruebas VALUES (:1, :2, :3, :4)"
const CONSULTA = "SELECT nombre, megabits FROM Usuarios WHERE cod = :1"

var fecha = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

func comprobar(err error) {
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}

func ms(t0 time.Time) float64 { return float64(time.Since(t0).Microseconds()) / 1000 }

func vaciar(db *sql.DB) {
	_, err := db.Exec("TRUNCATE TABLE BenchPruebas")
	comprobar(err)
}

// 1. Abrir y cerrar una conexión física (media por conexión)
func pruebaConexion(db *sql.DB, cods []string) float64 {
	t0 := time.Now()
	for i := 0; i < N_CONEXION; i++ {
		c, err := sql.Open("oracle", url)
		comprobar(err)
		comprobar(c.Ping())
		c.Close()
	}
	return ms(t0) / N_CONEXION
}

// 2. N_PK consultas por clave primaria con sentencia preparada
func pruebaConsultaPK(db *sql.DB, cods []string) float64 {
	ps, err := db.Prepare(CONSULTA)
	comprobar(err)
	defer ps.Close()
	var nombre string
	var mb float64
	t0 := time.Now()
	for i := 0; i < N_PK; i++ {
		comprobar(ps.QueryRow(cods[i%len(cods)]).Scan(&nombre, &mb))
	}
	return ms(t0)
}

// 3. N_TX inserciones dentro de UNA transacción
func pruebaInsertTx(db *sql.DB, cods []string) float64 {
	vaciar(db)
	t0 := time.Now()
	tx, err := db.Begin()
	comprobar(err)
	ps, err := tx.Prepare(INSERT)
	comprobar(err)
	for i := 0; i < N_TX; i++ {
		_, err := ps.Exec(i, fmt.Sprint("fila", i), float64(i)*0.5, fecha)
		comprobar(err)
	}
	ps.Close()
	comprobar(tx.Commit())
	return ms(t0)
}

// 4. N_AUTO inserciones con autocommit (cada una se confirma sola)
func pruebaInsertAuto(db *sql.DB, cods []string) float64 {
	vaciar(db)
	ps, err := db.Prepare(INSERT)
	comprobar(err)
	defer ps.Close()
	t0 := time.Now()
	for i := 0; i < N_AUTO; i++ {
		_, err := ps.Exec(i, fmt.Sprint("fila", i), float64(i)*0.5, fecha)
		comprobar(err)
	}
	return ms(t0)
}

// 5. Carga masiva de N_LOTE filas. database/sql no tiene API de lotes:
// la forma estándar es una transacción con una sentencia preparada.
func pruebaInsertLote(db *sql.DB, cods []string) float64 {
	vaciar(db)
	t0 := time.Now()
	tx, err := db.Begin()
	comprobar(err)
	ps, err := tx.Prepare(INSERT)
	comprobar(err)
	for i := 0; i < N_LOTE; i++ {
		_, err := ps.Exec(i, fmt.Sprint("fila", i), float64(i)*0.5, fecha)
		comprobar(err)
	}
	ps.Close()
	comprobar(tx.Commit())
	return ms(t0)
}

// 6. Leer las N_LOTE filas de BenchPruebas, las cuatro columnas
func pruebaLectura(db *sql.DB, cods []string) float64 {
	t0 := time.Now()
	rows, err := db.Query("SELECT id, nombre, valor, fecha FROM BenchPruebas")
	comprobar(err)
	n := 0
	for rows.Next() {
		var id int
		var nombre string
		var valor float64
		var f time.Time
		comprobar(rows.Scan(&id, &nombre, &valor, &f))
		n++
	}
	comprobar(rows.Err())
	rows.Close()
	t := ms(t0)
	if n != N_LOTE {
		comprobar(fmt.Errorf("lectura: %d filas, se esperaban %d", n, N_LOTE))
	}
	return t
}

// 7. CONC_HILOS goroutines, cada una con su conexión, CONC_CONSUL consultas
func pruebaConcurrente(db *sql.DB, cods []string) float64 {
	ctx := context.Background()
	conexiones := make([]*sql.Conn, CONC_HILOS) // se abren fuera de la medida
	for w := range conexiones {
		c, err := db.Conn(ctx)
		comprobar(err)
		conexiones[w] = c
	}
	var wg sync.WaitGroup
	t0 := time.Now()
	for w := 0; w < CONC_HILOS; w++ {
		wg.Add(1)
		go func(w int) {
			defer wg.Done()
			var nombre string
			var mb float64
			for i := 0; i < CONC_CONSUL; i++ {
				comprobar(conexiones[w].QueryRowContext(ctx, CONSULTA,
					cods[(w+i)%len(cods)]).Scan(&nombre, &mb))
			}
		}(w)
	}
	wg.Wait()
	t := ms(t0)
	for _, c := range conexiones {
		c.Close()
	}
	return t
}

func main() {
	db, err := sql.Open("oracle", url)
	comprobar(err)
	comprobar(db.Ping())
	defer db.Close()
	db.SetMaxOpenConns(CONC_HILOS + 1)
	db.SetMaxIdleConns(CONC_HILOS + 1)

	var servidor string
	db.QueryRow("SELECT banner FROM v$version WHERE ROWNUM = 1").Scan(&servidor)
	fmt.Println("=== COMPARATIVA: Go + Oracle ===")
	fmt.Println("Lenguaje:", runtime.Version(), "| Driver: go-ora v2")
	fmt.Println("Servidor:", servidor)
	fmt.Printf("Repeticiones: %d (+1 de calentamiento)\n\n", REPS)

	// Tabla auxiliar para las pruebas (se borra al terminar).
	// Oracle XE 21c no admite DROP TABLE IF EXISTS: si no existe, se ignora el error.
	db.Exec("DROP TABLE BenchPruebas")
	_, err = db.Exec("CREATE TABLE BenchPruebas (id NUMBER(10) PRIMARY KEY, " +
		"nombre VARCHAR2(30), valor NUMBER(10,2), fecha DATE)")
	comprobar(err)

	var cods []string
	rows, err := db.Query("SELECT cod FROM Usuarios ORDER BY cod")
	comprobar(err)
	for rows.Next() {
		var c string
		rows.Scan(&c)
		cods = append(cods, c)
	}
	rows.Close()

	nombres := []string{"conexion (ms/conexión)", "consulta_pk", "insert_tx",
		"insert_autocommit", "insert_lote", "lectura_masiva", "concurrente"}
	pruebas := []func(*sql.DB, []string) float64{pruebaConexion, pruebaConsultaPK,
		pruebaInsertTx, pruebaInsertAuto, pruebaInsertLote, pruebaLectura, pruebaConcurrente}
	valores := make([][]float64, len(pruebas))

	for r := 0; r <= REPS; r++ { // r = 0 es la repetición de calentamiento
		fmt.Printf("Repetición %d/%d...\n", r, REPS)
		for i, f := range pruebas {
			v := f(db, cods)
			if r > 0 {
				valores[i] = append(valores[i], v)
			}
		}
	}
	db.Exec("DROP TABLE BenchPruebas")

	fmt.Printf("\n%-24s %10s %10s %10s\n", "PRUEBA (ms)", "MEDIANA", "MIN", "MAX")
	for i, v := range valores {
		sort.Float64s(v)
		fmt.Printf("%-24s %10.1f %10.1f %10.1f\n", nombres[i], v[len(v)/2], v[0], v[len(v)-1])
	}
}
