# GO básico y su acceso a bases de datos (PostgreSQL y Oracle)

![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-pgx-336791?style=for-the-badge&logo=postgresql&logoColor=white)
![Oracle](https://img.shields.io/badge/Oracle-go--ora-F80000?style=for-the-badge&logo=oracle&logoColor=white)

Se introduce el lenguaje de programación **Go (Golang)** replicando los fundamentos de Python e implementando operaciones CRUD completas sobre **PostgreSQL** y **Oracle** mediante el paquete estándar `database/sql`. Además se incluye una batería de pruebas de rendimiento comparativas.


## Estructura del repositorio

```text
GO/
├── codigo/
│   ├── 01_introduccion/          # Fundamentos del lenguaje (variables, tipos, slices, maps,
│   │                             #   bucles, funciones, errores, defer, structs, interfaces)
│   ├── 02_postgresql/            # Conexión y CRUD parametrizado ($1, $2) con pgx
│   │                             #   + transacciones y demostración de inyección SQL
│   ├── 03_oracle/                # Conexión y CRUD parametrizado (:1, :2) con go-ora
│   │                             #   + transacciones y demostración de inyección SQL
│   └── 04_benchmark/             # Pruebas de rendimiento comparativas
│       ├── oracle/               #   Benchmark sobre Oracle
│       └── postgresql/           #   Benchmark sobre PostgreSQL
├── docs/
│   ├── 01_instalacion.md         # Guía de instalación de Go y drivers
│   ├── 02_introduccion_go.md     # Teoría y comparativa Python vs. Go
│   ├── 03_acceso_datos_sql.md    # Uso de database/sql, defer, transacciones y diferencias entre SGBD
│   └── 04_benchmark.md          # Descripción de las pruebas de rendimiento
├── go.mod                        # Definición del módulo y dependencias
├── go.sum                        # Sumas de verificación de dependencias
├── .gitignore
└── README.md
```


## Requisitos previos

* **Go** (versión 1.22 o superior) instalado en el sistema — ver [guía de instalación](docs/01_instalacion.md).
* Acceso a una instancia local o remota de **PostgreSQL** y/o **Oracle Database**.
* Drivers utilizados (se descargan automáticamente con `go mod tidy`):
  * PostgreSQL: `github.com/jackc/pgx/v5/stdlib`
  * Oracle: `github.com/sijms/go-ora/v2` (100 % Go puro, no requiere Oracle Instant Client).


## Puesta en marcha rápida

1. **Clona el repositorio:**
   ```bash
   git clone https://github.com/Gwen-p/GO.git
   cd GO
   ```

2. **Descarga las dependencias del módulo:**
   ```bash
   go mod tidy
   ```

3. **Configura tus credenciales:**
   Edita la cadena de conexión (`dsn` / `url`) en los archivos de acceso a datos:
   * `codigo/02_postgresql/main.go` — usuario, contraseña, host, puerto y base de datos.
   * `codigo/03_oracle/main.go` — usuario, contraseña, host, puerto y SID.

4. **Ejecuta los scripts:**
   ```bash
   # 1. Introducción a Go (no requiere base de datos)
   go run ./codigo/01_introduccion

   # 2. Acceso a PostgreSQL (CRUD + transacciones + inyección SQL)
   go run ./codigo/02_postgresql

   # 3. Acceso a Oracle (CRUD + transacciones + inyección SQL)
   go run ./codigo/03_oracle

   # 4. Benchmark sobre PostgreSQL
   go run ./codigo/04_benchmark/postgresql

   # 5. Benchmark sobre Oracle
   go run ./codigo/04_benchmark/oracle
   ```


## Contenido del proyecto

| Bloque | Script | Descripción |
|--------|--------|-------------|
| **A. Introducción** | `codigo/01_introduccion` | Variables, tipos, control de flujo, slices, maps, funciones con retornos múltiples, errores, `defer`, structs, métodos, punteros e interfaces. |
| **B. PostgreSQL** | `codigo/02_postgresql` | Pool de conexiones (`sql.Open`/`Ping`), consultas (`Query`/`QueryRow`), CRUD parametrizado (`$1, $2`), transacciones (`Begin`/`Commit`/`Rollback`), manejo de `NULL` y demostración de inyección SQL. |
| **C. Oracle** | `codigo/03_oracle` | Mismo flujo que PostgreSQL pero con marcadores `:1, :2`, `go-ora` como driver y `BuildUrl` para componer la URL de conexión. |
| **D. Benchmarks** | `codigo/04_benchmark` | 7 pruebas de rendimiento (conexión, consulta PK, insert TX, autocommit, carga masiva, lectura y concurrencia) para comparar Go con otros lenguajes. |


## Documentación del proyecto

* [Guía de instalación paso a paso](docs/01_instalacion.md)
* [Introducción a Go y comparativa con Python](docs/02_introduccion_go.md)
* [Acceso a datos con `database/sql`: PostgreSQL vs. Oracle](docs/03_acceso_datos_sql.md)
* [Pruebas de rendimiento (benchmark)](docs/04_benchmark.md)


## ¿Quieres probar Go sin instalar nada?

Si te ha picado la curiosidad y quieres escribir tu primer programa en Go directamente desde el navegador, te recomendamos completar el tutorial oficial interactivo en español:

* **[Un Tour por Go (en español)](https://go-tour-lat.appspot.com/list)** — Índice completo de lecciones interactivas para practicar sintaxis, estructuras de control y concurrencia paso a paso.
