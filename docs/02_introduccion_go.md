# Introducción al lenguaje Go (Comparativa con Python)

**Go (Golang)** es un lenguaje de programación de código abierto creado por Google en 2009. Nació con el objetivo de combinar la **legibilidad y rapidez de desarrollo de Python** con el **rendimiento y seguridad de compilación de C/C++**.

### Características principales

* **Compilado a código máquina:** Genera un único archivo ejecutable nativo sin necesidad de instalar intérpretes ni máquinas virtuales.
* **Tipado estático con inferencia:** Detecta errores de tipo antes de ejecutar el programa, pero permite escribir código conciso mediante el operador `:=`.
* **Minimalista:** Solo cuenta con **25 palabras reservadas**, lo que facilita aprender su sintaxis en muy poco tiempo.
* **Concurrencia nativa:** Permite ejecutar miles de tareas en paralelo de forma ligera mediante *goroutines* y canales (*channels*).
* **Gestión de errores explícita:** No utiliza excepciones (`try/except`); los errores se devuelven como valores y se comprueban con `if err != nil`.
* **Punteros sin aritmética:** Permite pasar datos por referencia de forma segura, sin los riesgos de C/C++.


## Comparativa directa: Python vs. Go

Tomando como referencia los conceptos vistos en clase en Python, esta tabla resume sus equivalentes en Go:

| Concepto | Python (visto en clase) | Go (equivalente) |
| :--- | :--- | :--- |
| **Comentarios** | `# una línea` / `""" varias """` | `// una línea` / `/* varias */` |
| **Variables** | Dinámicas (`edad = 20`) | Explícitas (`var edad int = 20`) o inferidas (`edad := 20`) |
| **Constantes** | No existe un tipo específico | `const` y `iota` para enumerados |
| **Entrada / Salida** | `print()` y `input()` | `fmt.Println()` / `fmt.Printf()` y `fmt.Scanln()` |
| **Conversión de tipos** | Implícita (ej. `int + float`) | Siempre **explícita** (ej. `float64(entero)`) |
| **Listas / Vectores** | Listas (`["Python", "Go"]`) | *Slices* (`[]string{"Python", "Go"}`) y `append()` |
| **Diccionarios** | `{"clave": "valor"}` | *Maps* (`map[string]int{"PostgreSQL": 5432}`) |
| **Bloques de código** | Indentación obligatoria y `:` | Llaves `{ }` obligatorias |
| **Condicionales** | `if` / `elif` / `else` | `if` / `else if` / `else` (y `switch` sin `break`) |
| **Bucles** | `while` y `for` | **Solo existe `for`** (sirve como `for`, `while` y `for range`) |
| **Funciones** | `def dividir(a, b):` | `func dividir(a, b float64) (float64, error)` |
| **Retornos múltiples** | `return a, b` (tupla) | `return a, b` (nativo, sin tuplas) |
| **Funciones anónimas** | `lambda x: x+1` | `func(x int) int { return x+1 }` (closures completas) |
| **Funciones variádicas** | `def f(*args):` | `func f(valores ...float64)` |
| **Control de errores** | `try` / `except` | Retorno explícito de tipo `error` y `if err != nil` |
| **Errores centinela** | Clases de excepción (`ValueError`) | `var ErrX = errors.New(...)` y `errors.Is(err, ErrX)` |
| **Clases** | `class Servidor:` | `type Servidor struct { ... }` con métodos asociados |
| **Herencia / Polimorfismo** | Herencia de clases | **Interfaces implícitas** (no se escribe `implements`) |
| **Punteros** | Todo es referencia | Explícitos con `*` y `&`, pero sin aritmética |
| **`defer`** | `with` / `finally` | `defer recurso.Close()` — se ejecuta al salir de la función (LIFO) |
| **Cierre de recursos** | `with open(...) as f:` | `defer db.Close()` / `defer rows.Close()` |


## Ejemplo práctico del repositorio

El archivo `codigo/01_introduccion/main.go` implementa todos los conceptos de la tabla anterior en 7 partes:

| Parte | Contenido | Conceptos demostrados |
|-------|-----------|----------------------|
| **1** | Estructura del programa | `package`, `import`, `const`, `iota`, funciones auxiliares |
| **2** | Variables y tipos | `var`, `:=`, valor cero, conversiones explícitas, `fmt.Printf` con `%T` |
| **3** | Control de flujo | `if` con inicialización, `for` clásico, `for` como `while`, `for range`, `switch` |
| **4** | Slices y mapas | `append`, subslices, `map`, `delete`, *comma ok*, recorrido ordenado |
| **5** | Funciones | Retornos múltiples, funciones variádicas, retornos con nombre, closures |
| **6** | Errores y `defer` | `strconv.Atoi`, errores centinela, `defer` en orden LIFO, error envuelto con `%w` |
| **7** | Structs, métodos e interfaces | Receptores por valor y puntero, `Describible` como interfaz implícita, `errors.Is` |


### Ejecución

```bash
go run ./codigo/01_introduccion
```

> **Nota:** Este script no requiere ninguna base de datos; se ejecuta de forma autónoma para demostrar las características del lenguaje.


## Recurso interactivo recomendado

Para practicar estos fundamentos directamente desde el navegador sin necesidad de instalar nada en el equipo, recomendamos seguir el tutorial oficial interactivo en español:

* **[Un Tour por Go (en español)](https://go-tour-lat.appspot.com/list)** — Incluye ejercicios guiados sobre paquetes, variables, control de flujo, *slices*, *maps*, métodos y concurrencia.
