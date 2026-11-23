
# Introducción al lenguaje Go (Comparativa con Python)

**Go (Golang)** es un lenguaje de programación de código abierto creado por Google en 2009. Nació con el objetivo de combinar la **legibilidad y rapidez de desarrollo de Python** con el **rendimiento y seguridad de compilación de C/C++**.

### Características principales
* **Compilado a código máquina:** Genera un único archivo ejecutable nativo sin necesidad de instalar intérpretes ni máquinas virtuales.
* **Tipado estático con inferencia:** Detecta errores de tipo antes de ejecutar el programa, pero permite escribir código conciso mediante el operador `:=`.
* **Minimalista:** Solo cuenta con **25 palabras reservadas**, lo que facilita aprender su sintaxis en muy poco tiempo.
* **Concurrencia nativa:** Permite ejecutar miles de tareas en paralelo de forma ligera mediante *goroutines* y canales (*channels*).


## Comparativa directa: Python vs. Go

Tomando como referencia los conceptos vistos en clase en Python, esta tabla resume sus equivalentes en Go:

| Concepto | Python (visto en clase) | Go (equivalente) |
| :--- | :--- | :--- |
| **Comentarios** | `# una línea` / `""" varias """` | `// una línea` / `/* varias */` |
| **Variables** | Dinámicas (`edad = 20`) | Explícitas (`var edad int = 20`) o inferidas (`edad := 20`) |
| **Entrada / Salida** | `print()` y `input()` | `fmt.Println()` / `fmt.Printf()` y `fmt.Scanln()` |
| **Listas / Vectores** | Listas (`["Python", "Go"]`) | *Slices* (`[]string{"Python", "Go"}`) y `append()` |
| **Diccionarios** | `{"clave": "valor"}` | *Maps* (`map[string]int{"PostgreSQL": 5432}`) |
| **Bloques de código** | Indentación obligatoria y `:` | Llaves `{ }` obligatorias |
| **Condicionales** | `if` / `elif` / `else` | `if` / `else if` / `else` (y `switch`) |
| **Bucles** | `while` y `for` | **Solo existe `for`** (sirve como `for`, `while` y `for range`) |
| **Funciones** | `def dividir(a, b):` | `func dividir(a, b float64) (float64, error)` |
| **Control de errores** | Excepciones (`try` / `except`) | Retorno explícito de tipo `error` y comprobación `if err != nil` |


## Ejemplo práctico del repositorio

El archivo `codigo/01_introduccion/01_introduccion.go` implementa todos los elementos de la tabla anterior:
1. Lectura de texto por teclado y formateo de salida con `fmt.Printf`.
2. Creación y ampliación de un *slice* con `append`.
3. Uso de un *map* para almacenar los puertos de PostgreSQL y Oracle.
4. Bifurcación `if / else` y recorrido de colecciones con `for indice, valor := range`.
5. Función `dividir(a, b)` que devuelve dos valores (resultado y error), demostrando cómo Go sustituye el bloque `try / except`.

### Ejecución
```bash
go run ./codigo/01_introduccion/01_introduccion.go
```

## Recurso interactivo recomendado

Para practicar estos fundamentos directamente desde el navegador sin necesidad de instalar nada en el equipo, recomendamos seguir el tutorial oficial interactivo en español:

* **[Un Tour por Go (en español)](https://go-tour-lat.appspot.com/list)** — Incluye ejercicios guiados sobre paquetes, variables, control de flujo, *slices*, *maps*, métodos y concurrencia.


