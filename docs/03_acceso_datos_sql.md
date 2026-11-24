# Acceso a Datos con `database/sql`: PostgreSQL y Oracle

Una de las mayores ventajas de Go frente a otros lenguajes es que incluye de serie el paquete estándar **`database/sql`**. Esto proporciona una interfaz común para interactuar con distintos motores de bases de datos relacionales sin cambiar la estructura del código.

```mermaid
flowchart LR
    A["Código Go"] --> B["database/sql"]
    B --> C["Driver<br/>(pgx / go-ora)"]
    C --> D["PostgreSQL / Oracle"]

    style A fill:#cfe2ff,color:#000,stroke:#333,stroke-width:2px
    style B fill:#d9ead3,color:#000,stroke:#333,stroke-width:2px
    style C fill:#fff2cc,color:#000,stroke:#333,stroke-width:2px
    style D fill:#f4cccc,color:#000,stroke:#333,stroke-width:2px
```

Para registrar el controlador sin llamarlo directamente en el código, se importa precedido de un guion bajo (`_`):

```go
import (
    "database/sql"
    _ "github.com/jackc/pgx/v5/stdlib" // Driver PostgreSQL
    // o bien:
    _ "github.com/sijms/go-ora/v2"     // Driver Oracle
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

### Configuración del pool de conexiones

El objeto `*sql.DB` no es una única conexión, sino un **pool** que gestiona automáticamente múltiples conexiones. Se puede configurar:

```go
db.SetMaxOpenConns(5)                    // máximo de conexiones abiertas simultáneamente
db.SetConnMaxLifetime(30 * time.Minute)  // renovar conexiones antiguas
```


## 2. Consultas (`SELECT`) y recorrido de resultados

### Consultar múltiples filas: `db.Query()`

Para ejecutar consultas que devuelven varias filas se utiliza `db.Query()`. El objeto `rows` resultante actúa como el cursor/iterador de Python:

```go
rows, err := db.Query("SELECT cod, nombre FROM Servidores ORDER BY cod")
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
// Comprobar si hubo errores durante el recorrido
if err := rows.Err(); err != nil {
    log.Fatal(err)
}
```

### Consultar una única fila: `db.QueryRow()`

Cuando se espera una sola fila (por ejemplo, una búsqueda por clave primaria), se usa `QueryRow`. Si no encuentra resultados, devuelve `sql.ErrNoRows`:

```go
var nombre string
var megabits float64
err := db.QueryRow("SELECT nombre, megabits FROM Usuarios WHERE cod = $1", cod).
    Scan(&nombre, &megabits)

switch {
case errors.Is(err, sql.ErrNoRows):
    fmt.Println("No existe ningún usuario con ese código")
case err != nil:
    fmt.Println("Error en la consulta:", err)
default:
    fmt.Printf("Usuario: %s, %.2f megabits\n", nombre, megabits)
}
```

### Manejo de valores `NULL` con tipos `sql.Null*`

Cuando una columna puede contener `NULL`, se usa un tipo especial como `sql.NullInt64`, que tiene dos campos: `Int64` (el valor) y `Valid` (si no es NULL):

```go
var servidor sql.NullInt64
rows.Scan(&cod, &nombre, &servidor)
if servidor.Valid {
    fmt.Println("Servidor:", servidor.Int64)
} else {
    fmt.Println("Servidor: NULL")
}
```


## 3. Inserción, modificación y borrado (`INSERT`, `UPDATE`, `DELETE`)

Las sentencias que no devuelven conjuntos de filas se ejecutan con **`db.Exec()`**. Para prevenir ataques de **inyección SQL**, nunca se deben concatenar variables dentro de la cadena SQL; siempre se utilizan **consultas parametrizadas** (vinculación de variables).

```go
// INSERT
res, err := db.Exec("INSERT INTO Servidores VALUES ($1, $2, 0, 0)", cod, nombre)
n, _ := res.RowsAffected()
fmt.Println("Filas insertadas:", n)

// UPDATE
res, err = db.Exec("UPDATE Servidores SET nombreS = $1 WHERE codS = $2", nuevoNombre, cod)

// DELETE
res, err = db.Exec("DELETE FROM Servidores WHERE codS = $1", cod)
```


## 4. Transacciones (`BEGIN`, `COMMIT`, `ROLLBACK`)

Una transacción permite ejecutar varias operaciones como una unidad atómica: o se confirman todas, o no se guarda ninguna.

```go
tx, err := db.Begin()         // inicia la transacción
if err != nil {
    log.Fatal(err)
}
defer tx.Rollback()           // deshace todo si salimos antes del Commit

ps, err := tx.Prepare("INSERT INTO Servidores VALUES ($1, $2, 0, 0)")
if err != nil {
    log.Fatal(err)
}
defer ps.Close()

for _, s := range listaServidores {
    if _, err := ps.Exec(s.Cod, s.Nombre); err != nil {
        fmt.Println("Error, se hace ROLLBACK:", err)
        return  // el defer hace el Rollback automáticamente
    }
}

if err := tx.Commit(); err != nil {  // confirma todos los cambios
    log.Fatal(err)
}
```

> **Patrón clave:** Se usa `defer tx.Rollback()` inmediatamente después de `Begin()`. Si se llega al `Commit()`, el `Rollback` posterior no tiene efecto. Si se sale por un error antes del `Commit`, el `Rollback` se ejecuta automáticamente.


## 5. Inyección SQL: consultas parametrizadas vs. concatenación

Los scripts incluyen una función `dbBuscarUsuarioConcatenado` que demuestra el problema de la inyección SQL, contraponiéndola con la versión segura parametrizada:

```go
//  MAL: concatenación (vulnerable a inyección SQL)
consulta := "SELECT cod FROM Usuarios WHERE cod = '" + cod + "'"
rows, err := db.Query(consulta)
// Con la entrada: ' OR '1'='1    →  devuelve TODAS las filas

//  BIEN: parámetro vinculado (seguro contra inyección SQL)
err := db.QueryRow("SELECT nombre FROM Usuarios WHERE cod = $1", cod).Scan(&nombre)
// Con la entrada: ' OR '1'='1    →  no encuentra nada (busca ese literal)
```


## 6. Diferencias clave entre PostgreSQL y Oracle en Go

Aunque ambos scripts utilizan exactamente las mismas funciones de `database/sql`, existen diferencias sintácticas importantes entre ambos SGBD:

| Aspecto | PostgreSQL (`pgx`) | Oracle (`go-ora/v2`) |
| --- | --- | --- |
| **Formato de DSN** | `postgres://usuario:pass@localhost:5432/empresa` | `oracle://usuario:pass@localhost:1521/XE` u `oracle.BuildUrl(...)` |
| **Marcadores de parámetros** | Posicionales con dólar: `$1, $2, $3...` | Dos puntos: `:1, :2...` o por nombre `:codigo` |
| **Ejemplo de `INSERT`** | `db.Exec("INSERT INTO s VALUES ($1, $2)", c, n)` | `db.Exec("INSERT INTO s VALUES (:1, :2)", c, n)` |
| **Creación de tabla segura** | Soporta `CREATE TABLE IF NOT EXISTS` | En versiones `< 23ai` no soporta `IF NOT EXISTS`; se captura el error `ORA-00955` |
| **Tipos de columna típicos** | `INT`, `VARCHAR(100)`, `DATE` | `NUMBER`, `VARCHAR2(100)`, `DATE` |
| **Nombres de atributos** | Devueltos por defecto en **minúsculas** | Devueltos por defecto en **MAYÚSCULAS** |
| **Versión del servidor** | `SELECT version()` | `SELECT banner FROM v$version WHERE ROWNUM = 1` |
| **Autocommit** | Habilitado por defecto fuera de transacciones | Habilitado por defecto fuera de transacciones |


## 7. Ejecución de los scripts

### PostgreSQL (`codigo/02_postgresql/main.go`)

Realiza el ciclo completo: conexión, consultas de todas las tablas, búsqueda por PK, inserción, modificación, borrado, transacciones y demostración de inyección SQL.

```bash
go run ./codigo/02_postgresql
```

### Oracle (`codigo/03_oracle/main.go`)

Mismo flujo que PostgreSQL pero adaptado a la sintaxis de Oracle (marcadores `:1`, `:2` y `BuildUrl` para la conexión).

```bash
go run ./codigo/03_oracle
```
