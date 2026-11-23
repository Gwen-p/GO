
# Acceso a Datos con `database/sql`: PostgreSQL y Oracle

Una de las mayores ventajas de Go frente a otros lenguajes es que incluye de serie el paquete estándar **`database/sql`**. Esto proporciona una interfaz común para interactuar con distintos motores de bases de datos relacionales sin cambiar la estructura del código.

```text
[ Código Go ] ──▶ [ database/sql ] ──▶ [ Driver (pgx / go-ora) ] ──▶ [ PostgreSQL / Oracle ]
```

Para registrar el controlador sin llamarlo directamente en el código, se importa precedido de un guion bajo (`_`):

```go
import (
    "database/sql"
    _ "[github.com/jackc/pgx/v5/stdlib](https://github.com/jackc/pgx/v5/stdlib)" // Driver PostgreSQL
    // o bien:
    _ "[github.com/sijms/go-ora/v2](https://github.com/sijms/go-ora/v2)"     // Driver Oracle
)
```



## 1. Conexión y desconexión segura (`defer`)

El flujo de conexión replica el esquema de `dbConectar()` y `dbDesconectar()` visto en clase, pero introduce dos particularidades clave de Go:

1. **`sql.Open(driver, dsn)`:** Crea un objeto `*sql.DB` que gestiona un **pool de conexiones** automático. No abre una conexión de red al instante; solo valida los parámetros.
2. **`db.Ping()`:** Comprueba realmente que el servidor de base de datos responde y las credenciales son válidas.
3. **`defer db.Close()`:** Programa el cierre automático de la conexión cuando termine la función actual, ocurra o no un error más adelante.

```go
db, err := sql.Open("pgx", dsn)
if err != nil {
    log.Fatal(err)
}
defer db.Close()

if err := db.Ping(); err != nil {
    log.Fatal("No se pudo conectar:", err)
}
```


## 2. Consultas (`SELECT`) y recorrido de resultados

Para ejecutar consultas que devuelven filas se utiliza `db.Query()`. El objeto `rows` resultante actúa como el cursor/iterador de Python:

```go
rows, err := db.Query("SELECT cod, nombre FROM sectores ORDER BY cod")
if err != nil {
    log.Fatal(err)
}
defer rows.Close() // Imprescindible liberar el cursor al terminar

for rows.Next() {
    var cod int
    var nombre string
    if err := rows.Scan(&cod, &nombre); err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Código: %d | Nombre: %s\n", cod, nombre)
}
```


## 3. Inserción, modificación y borrado (`INSERT`, `UPDATE`, `DELETE`)

Las sentencias que no devuelven conjuntos de filas se ejecutan con **`db.Exec()`**. Para prevenir ataques de **Inyección SQL**, nunca se deben concatenar variables dentro de la cadena SQL; siempre se utilizan **consultas parametrizadas** (vinculación de variables).



## 4. Diferencias clave entre PostgreSQL y Oracle en Go

Aunque ambos scripts utilizan exactamente las mismas funciones de `database/sql`, existen diferencias sintácticas importantes entre ambos SGBD:

| Aspecto | PostgreSQL (`pgx`) | Oracle (`go-ora/v2`) |
| --- | --- | --- |
| **Formato de DSN** | `postgres://usuario:pass@localhost:5432/empresa` | `oracle://usuario:pass@localhost:1521/XE` |
| **Marcadores de parámetros** | Posicionales con dólar: `$1, $2, $3...` | Dos puntos: `:1, :2...` o por nombre `:codigo` |
| **Ejemplo de `INSERT**` | `db.Exec("INSERT INTO sectores VALUES ($1, $2)", c, n)` | `db.Exec("INSERT INTO SECTORES VALUES (:1, :2)", c, n)` |
| **Creación de tabla segura** | Soporta `CREATE TABLE IF NOT EXISTS` | En versiones `< 23ai` no soporta `IF NOT EXISTS`; se captura el error `ORA-00955` |
| **Tipos de columna típicos** | `INT`, `VARCHAR(100)` | `NUMBER`, `VARCHAR2(100)` |
| **Nombres de atributos** | Devueltos por defecto en **minúsculas** | Devueltos por defecto en **MAYÚSCULAS** |



## 5. Ejecución y evidencias

### PostgreSQL (`codigo/02_postgresql/02_acceso_postgresql.go`)

Incluye una entrada interactiva por teclado pensada para la demostración en vivo en clase.

```bash
go run ./codigo/02_postgresql/02_acceso_postgresql.go
```

### Oracle (`codigo/03_oracle/03_acceso_oracle.go`)

Ejecuta la creación segura de la tabla `SECTORES`, inserción, actualización, consulta y borrado del registro de prueba.

```bash
go run ./codigo/03_oracle/03_acceso_oracle.go
```

