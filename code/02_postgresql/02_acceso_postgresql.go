package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib" // Driver de PostgreSQL para database/sql
)

func main() {
	// 1. Cadena de conexión (ajusta usuario, contraseña, puerto y nombre de BD)
	dsn := "postgres://postgres:12345@localhost:5432/empresa"

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal("Error en la configuración:", err)
	}
	defer db.Close() // Cierre automático al terminar main()

	// 2. Verificar conexión real
	if err := db.Ping(); err != nil {
		log.Fatal("No se pudo conectar a PostgreSQL:", err)
	}
	fmt.Println("✔ Conectado a PostgreSQL con éxito.")

	// 3. Crear tabla de ejemplo si no existe (útil para probar rápido)
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS sectores (
		cod INT PRIMARY KEY,
		nombre VARCHAR(100)
	)`)
	if err != nil {
		log.Fatal("Error creando tabla:", err)
	}

	// 4. INSERT parametrizado (PostgreSQL usa $1, $2) — Ideal para la DEMO EN VIVO
	var cod int
	var nombre string
	fmt.Print("\n[DEMO] Introduce un código numérico para el nuevo sector: ")
	fmt.Scanln(&cod)
	fmt.Print("[DEMO] Introduce el nombre del sector: ")
	fmt.Scanln(&nombre)

	_, err = db.Exec("INSERT INTO sectores (cod, nombre) VALUES ($1, $2)", cod, nombre)
	if err != nil {
		log.Println("Aviso al insertar:", err)
	} else {
		fmt.Println("✔ Registro insertado correctamente.")
	}

	// 5. UPDATE parametrizado
	_, err = db.Exec("UPDATE sectores SET nombre = $1 WHERE cod = $2", nombre+" (Verificado)", cod)
	if err != nil {
		log.Println("Error al actualizar:", err)
	}

	// 6. SELECT (lectura con cursor/iterador)
	fmt.Println("\n--- Listado de Sectores en PostgreSQL ---")
	rows, err := db.Query("SELECT cod, nombre FROM sectores ORDER BY cod")
	if err != nil {
		log.Fatal("Error en SELECT:", err)
	}
	defer rows.Close()

	for rows.Next() {
		var c int
		var n string
		if err := rows.Scan(&c, &n); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Código: %d | Nombre: %s\n", c, n)
	}

	// 7. DELETE (descomenta si quieres probar el borrado)
	// _, err = db.Exec("DELETE FROM sectores WHERE cod = $1", cod)
}