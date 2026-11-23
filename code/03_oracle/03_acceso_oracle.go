package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/sijms/go-ora/v2" // Driver puro de Oracle (sin Instant Client)
)

func main() {
	// 1. Cadena de conexión: oracle://usuario:contraseña@host:puerto/servicio
	dsn := "oracle://system:12345@localhost:1521/XE"

	db, err := sql.Open("oracle", dsn)
	if err != nil {
		log.Fatal("Error en la configuración:", err)
	}
	defer db.Close()

	// 2. Verificar conexión real
	if err := db.Ping(); err != nil {
		log.Fatal("No se pudo conectar a Oracle:", err)
	}
	fmt.Println("✔ Conectado a Oracle con éxito.")

	// 2.5. Crear la tabla si no existe (compatible con cualquier versión de Oracle)
	// Si en clase estáis usando Oracle 23ai, entonces sí puedes escribir directamente CREATE TABLE IF NOT EXISTS SECTORES (...)
	_, err = db.Exec(`CREATE TABLE SECTORES (
		COD NUMBER PRIMARY KEY,
		NOMBRE VARCHAR2(100)
	)`)
	if err != nil {
		// Si ya existe (ORA-00955), simplemente continuamos; si es otro error, lo mostramos
		fmt.Println("Nota sobre la tabla SECTORES:", err)
	} else {
		fmt.Println("✔ Tabla SECTORES creada correctamente.")
	}

	// 3. INSERT parametrizado (Oracle usa :1, :2)
	cod := 101
	nombre := "Tecnología"
	_, err = db.Exec("INSERT INTO SECTORES (COD, NOMBRE) VALUES (:1, :2)", cod, nombre)
	if err != nil {
		log.Println("Aviso al insertar (puede que el código ya exista):", err)
	} else {
		fmt.Println("✔ Registro insertado en Oracle.")
	}

	// 4. UPDATE parametrizado
	_, err = db.Exec("UPDATE SECTORES SET NOMBRE = :1 WHERE COD = :2", "Tecnología e IA", cod)
	if err != nil {
		log.Println("Error al actualizar:", err)
	}

	// 5. SELECT (lectura con rows.Next)
	fmt.Println("\n--- Listado de Sectores en Oracle ---")
	rows, err := db.Query("SELECT COD, NOMBRE FROM SECTORES ORDER BY COD")
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

	// 6. DELETE parametrizado
	_, err = db.Exec("DELETE FROM SECTORES WHERE COD = :1", cod)
	if err != nil {
		log.Println("Error al borrar:", err)
	} else {
		fmt.Println("\n✔ Registro de prueba eliminado.")
	}
}