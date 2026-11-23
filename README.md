# GO básico y su acceso a bases de datos (PostgreSQL y Oracle)

![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-pgx-336791?style=for-the-badge&logo=postgresql&logoColor=white)
![Oracle](https://img.shields.io/badge/Oracle-go--ora-F80000?style=for-the-badge&logo=oracle&logoColor=white)

Se introduce el lenguaje de programación **Go (Golang)** replicando los fundamentos de Python e implementando operaciones CRUD completas sobre **PostgreSQL** y **Oracle** mediante el paquete estándar `database/sql`.


## Estructura del repositorio

```text
GO/
├── codigo/
│   ├── 01_introduccion/      # Fundamentos del lenguaje (variables, slices, maps, bucles, errores)
│   ├── 02_postgresql/        # Conexión y CRUD parametrizado ($1, $2) con pgx
│   └── 03_oracle/            # Conexión y CRUD parametrizado (:1, :2) con go-ora
└── docs/
    ├── 01_instalacion.md     # Guía de instalación de Go y drivers
    ├── 02_introduccion_go.md # Teoría y comparativa Python vs. Go
    └── 03_acceso_datos_sql.md# Uso de database/sql, defer y diferencias entre SGBD

```


## Requisitos previos

* **Go** (versión 1.21 o superior) instalado en el sistema.
* Acceso a una instancia local o remota de **PostgreSQL** y/o **Oracle Database**.
* Drivers utilizados (se descargan automáticamente con Go):
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
Edita la cadena de conexión (`dsn`) en `codigo/02_postgresql/02_acceso_postgresql.go` y `codigo/03_oracle/03_acceso_oracle.go` con tu usuario, contraseña, host y puerto.
4. **Ejecuta los scripts:**
```bash
# 1. Introducción a Go
go run ./codigo/01_introduccion/01_introduccion.go

# 2. Acceso a PostgreSQL (incluye inserción interactiva para demo)
go run ./codigo/02_postgresql/02_acceso_postgresql.go

# 3. Acceso a Oracle
go run ./codigo/03_oracle/03_acceso_oracle.go

```




## Documentación del proyecto

* [Guía de instalación paso a paso](https://www.google.com/search?q=docs/01_instalacion.md)
* [Introducción a Go y comparativa con Python](https://www.google.com/search?q=docs/02_introduccion_go.md)
* [Acceso a datos con `database/sql`: PostgreSQL vs. Oracle](https://www.google.com/search?q=docs/03_acceso_datos_sql.md)



## ¿Quieres probar Go sin instalar nada?

Si te ha picado la curiosidad y quieres escribir tu primer programa en Go directamente desde el navegador, te recomendamos completar el tutorial oficial interactivo en español:

* **[Un Tour por Go (en español)](https://go-tour-lat.appspot.com/list)** — Índice completo de lecciones interactivas para practicar sintaxis, estructuras de control y concurrencia paso a paso.

