
# Guía de instalación y configuración del entorno

En este documento se detalla el proceso paso a paso para instalar el compilador de **Go**, inicializar el módulo del proyecto y descargar los controladores (*drivers*) necesarios para conectar con **PostgreSQL** y **Oracle**.


## 1. Versiones recomendadas

| Software / Componente | Versión recomendada | Notas |
| :--- | :--- | :--- |
| **Go (Golang)** | `1.22` o superior | Incluye el paquete estándar `database/sql`. |
| **PostgreSQL** | `15` / `16` | Puerto por defecto: `5432`. |
| **Oracle Database** | `21c XE` o `23ai Free` | Puerto por defecto: `1521` (servicio `XE` o `FREEPDB1`). |
| **Driver PostgreSQL (`pgx`)** | `v5` (`github.com/jackc/pgx/v5/stdlib`) | Compatible directamente con `database/sql`. |
| **Driver Oracle (`go-ora`)** | `v2` (`github.com/sijms/go-ora/v2`) | Escrito 100 % en Go puro; **no requiere** instalar Oracle Instant Client ni compilador C. |
| **Editor recomendado** | Visual Studio Code | Instalar la extensión oficial **Go** (de Go Team at Google). |



## 2. Instalación de Go

### En Windows
1. Descarga el instalador `.msi` desde la web oficial: [https://go.dev/dl/](https://go.dev/dl/).
2. Ejecuta el asistente manteniendo la ruta por defecto (`C:\Program Files\Go`). El instalador añadirá Go automáticamente a la variable de entorno `PATH`.
3. Abre una nueva terminal (PowerShell o CMD) y comprueba la instalación:
```bash
   go version
```

### En Linux (Ubuntu / Debian)

1. Descarga el archivo `.tar.gz` desde [https://go.dev/dl/](https://go.dev/dl/) y extráelo en `/usr/local`:
```bash
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.22.x.linux-amd64.tar.gz
```


2. Añade Go al `PATH` editando `~/.bashrc` o `~/.zshrc`:
```bash
export PATH=$PATH:/usr/local/go/bin
```


3. Recarga la configuración y verifica la versión:
```bash
source ~/.bashrc
go version
```



### En macOS

* Instala mediante el paquete `.pkg` oficial desde [https://go.dev/dl/](https://go.dev/dl/) o usando Homebrew:
```bash
brew install go
go version
```


## 3. Inicialización del proyecto y descarga de drivers

Si partes desde cero (sin clonar el repositorio), abre una terminal dentro de la carpeta raíz del proyecto y ejecuta:

```bash
# 1. Crear el archivo go.mod que gestiona las dependencias
go mod init proyecto_go_bd

# 2. Descargar el driver de PostgreSQL
go get [github.com/jackc/pgx/v5/stdlib](https://github.com/jackc/pgx/v5/stdlib)

# 3. Descargar el driver de Oracle (Go puro)
go get [github.com/sijms/go-ora/v2](https://github.com/sijms/go-ora/v2)

# 4. Limpiar y sincronizar dependencias en go.mod y go.sum
go mod tidy
```

Si ya has clonado este repositorio, basta con ejecutar `go mod tidy` en la raíz para descargar automáticamente las versiones exactas registradas en `go.mod`.



## 4. Verificación previa de las bases de datos

Antes de ejecutar los scripts de acceso a datos, asegúrate de que los servicios de base de datos están activos:

* **PostgreSQL:** Servicio activo en `localhost:5432` con la base de datos `empresa` creada (o ajusta el nombre en la variable `dsn`).
* **Oracle:** Instancia activa en `localhost:1521` y servicio `XE` (o el nombre de servicio configurado en tu máquina) escuchando conexiones.
