# Pruebas de rendimiento (Benchmark): PostgreSQL y Oracle

Los scripts de benchmark permiten medir y comparar el rendimiento de las operaciones más comunes de acceso a datos desde Go, tanto sobre **PostgreSQL** como sobre **Oracle**. Están diseñados para ser comparables con implementaciones equivalentes en Java (JDBC) y Python.


## 1. Parámetros de las pruebas

Todos los parámetros son **idénticos** en los programas de los distintos lenguajes para garantizar una comparación justa:

| Constante | Valor | Descripción |
|-----------|-------|-------------|
| `REPS` | 5 | Repeticiones medidas (+ 1 de calentamiento) |
| `N_CONEXION` | 10 | Conexiones abiertas y cerradas por repetición |
| `N_PK` | 2.000 | Consultas por clave primaria con sentencia preparada |
| `N_TX` | 2.000 | Inserciones dentro de una transacción |
| `N_AUTO` | 500 | Inserciones con autocommit (un commit por cada una) |
| `N_LOTE` | 20.000 | Filas de la carga masiva |
| `CONC_HILOS` | 4 | Goroutines (trabajadores concurrentes) |
| `CONC_CONSUL` | 500 | Consultas por trabajador concurrente |


## 2. Las 7 pruebas

| # | Prueba | Qué mide |
|---|--------|----------|
| 1 | **Conexión** | Tiempo medio de abrir y cerrar una conexión física (no del pool). |
| 2 | **Consulta por PK** | Tiempo total de `N_PK` consultas `SELECT` por clave primaria con sentencia preparada. |
| 3 | **Insert en transacción** | Tiempo de insertar `N_TX` filas dentro de una sola transacción (`BEGIN`...`COMMIT`). |
| 4 | **Insert con autocommit** | Tiempo de insertar `N_AUTO` filas con un commit implícito por cada inserción. |
| 5 | **Carga masiva (lote)** | Tiempo de insertar `N_LOTE` filas en una transacción con sentencia preparada. |
| 6 | **Lectura masiva** | Tiempo de leer las `N_LOTE` filas insertadas, recorriendo las cuatro columnas. |
| 7 | **Concurrencia** | Tiempo total de `CONC_HILOS` goroutines ejecutando `CONC_CONSUL` consultas cada una, en paralelo. |


## 3. Requisitos previos

Además de los requisitos generales del proyecto (ver [01_instalacion.md](01_instalacion.md)):

* La base de datos debe contener la tabla **`Usuarios`** con al menos una fila (se usa para las consultas por PK).
* La tabla **`BenchPruebas`** se crea y destruye automáticamente durante la ejecución; no es necesario crearla manualmente.


## 4. Ejecución

```bash
# Benchmark sobre PostgreSQL
go run ./codigo/04_benchmark/postgresql

# Benchmark sobre Oracle
go run ./codigo/04_benchmark/oracle
```

> **Nota:** La primera repetición (`0/5`) es de calentamiento y no se incluye en los resultados. Esto asegura que la JIT del servidor, las cachés y el pool de conexiones estén estabilizados.


## 5. Interpretación de los resultados

La salida muestra una tabla con tres estadísticas para cada prueba:

```
PRUEBA (ms)                MEDIANA        MIN        MAX
conexion (ms/conexión)        12.3       11.8       13.1
consulta_pk                  145.2      140.0      155.3
insert_tx                     85.0       82.1       90.4
insert_autocommit            520.3      510.2      535.1
insert_lote                  850.0      830.5      880.2
lectura_masiva               120.4      115.3      125.0
concurrente                  200.1      195.0      210.3
```

| Columna | Significado |
|---------|-------------|
| **MEDIANA** | El valor central de las `REPS` mediciones. Es el indicador más fiable, ya que no se ve afectado por picos puntuales. |
| **MIN** | El mejor tiempo obtenido en las `REPS` repeticiones. |
| **MAX** | El peor tiempo obtenido. Una diferencia grande entre MIN y MAX sugiere variabilidad en la carga del sistema. |

### Observaciones típicas

* **Insert autocommit vs. transacción:** El autocommit es significativamente más lento porque cada inserción requiere un `fsync` al disco. Usar transacciones agrupa los cambios en un solo `fsync`.
* **Concurrencia en Go:** Las goroutines son mucho más ligeras que los hilos de Java o los procesos de Python, lo que se refleja en tiempos competitivos con pocos recursos.
* **Conexión:** Go con `go-ora` (puro) suele ser más rápido que un driver que dependa de un cliente nativo (como Oracle Instant Client).
