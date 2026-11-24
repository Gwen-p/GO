// =====================================================================
//  PROYECTO: ACCESO A DATOS DESDE GO  --  Programación de Bases de Datos
//  Bloque A. Introducción a Go: script de pruebas del lenguaje
//  Ejecución (desde la raíz del módulo):  go run ./codigo/01_introduccion
// =====================================================================

package main

// ---------------------------------------------------------------------
// PARTE 1. Estructura de un programa: paquete, imports y constantes
// ---------------------------------------------------------------------
import (
	"errors"  // crear y comparar errores
	"fmt"     // imprimir por pantalla con formato
	"sort"    // ordenar slices
	"strconv" // convertir texto a número
	"strings" // utilidades de cadenas
)

const Asignatura = "Programación de Bases de Datos" // constante sin tipo explícito

// iota numera constantes consecutivas: sirve para definir enumerados.
const (
	Oracle   = iota // 0
	Postgres        // 1
)

// seccion imprime un título para separar las partes del script.
func seccion(titulo string) {
	fmt.Println("\n==== " + strings.ToUpper(titulo) + " ====")
}

// ---------------------------------------------------------------------
// PARTE 5. Funciones (se declaran fuera de main y se usan más abajo)
// ---------------------------------------------------------------------

// dividir devuelve DOS valores: el resultado y un error.
func dividir(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("división por cero")
	}
	return a / b, nil // nil significa "sin error"
}

// estadisticas es VARIÁDICA (recibe cualquier número de valores) y usa
// retornos CON NOMBRE: total y media se declaran en la firma.
func estadisticas(valores ...float64) (total, media float64) {
	for _, v := range valores {
		total += v
	}
	if len(valores) > 0 {
		media = total / float64(len(valores))
	}
	return // devuelve total y media
}

// ---------------------------------------------------------------------
// PARTE 6. Errores y defer
// ---------------------------------------------------------------------

// ErrNoEncontrado es un error "centinela": un valor fijo con el que
// comparar, igual que sql.ErrNoRows en el acceso a datos.
var ErrNoEncontrado = errors.New("servidor no encontrado")

// Conexion simula un recurso que hay que cerrar siempre (como una
// conexión o un cursor de base de datos).
type Conexion struct{ nombre string }

func abrir(nombre string) *Conexion {
	fmt.Println("  abrir", nombre)
	return &Conexion{nombre: nombre}
}

func (c *Conexion) Close() { fmt.Println("  cerrar", c.nombre) }

// usarRecursos muestra que defer ejecuta las llamadas al salir de la
// función, en orden inverso (LIFO), aunque se salga por un error.
func usarRecursos(fallar bool) error {
	con := abrir("conexión")
	defer con.Close() // se ejecutará la última
	cur := abrir("cursor")
	defer cur.Close() // se ejecutará la primera
	if fallar {
		return errors.New("fallo a mitad de la consulta")
	}
	fmt.Println("  trabajando con los datos...")
	return nil
}

// ---------------------------------------------------------------------
// PARTE 7. Structs, métodos e interfaces
// ---------------------------------------------------------------------

// Servidor agrupa los datos de una tupla de la tabla Servidores.
type Servidor struct {
	Cod      int
	Nombre   string
	Uso      float64
	Megabits float64
}

// Usuario agrupa algunos datos de una tupla de la tabla Usuarios.
type Usuario struct {
	Cod      string
	Nombre   string
	Servidor int
}

// Método con receptor por VALOR: trabaja sobre una copia.
func (s Servidor) Describir() string {
	return fmt.Sprintf("Servidor %d (%s): %.2f megabits", s.Cod, s.Nombre, s.Megabits)
}

// Método con receptor PUNTERO: modifica el struct original.
func (s *Servidor) RegistrarTrafico(mb float64) {
	s.Megabits += mb
}

func (u Usuario) Describir() string {
	return fmt.Sprintf("Usuario %s (%s) en el servidor %d", u.Cod, u.Nombre, u.Servidor)
}

// Describible es una INTERFAZ: cualquier tipo con un método Describir()
// la cumple automáticamente, sin escribir "implements".
type Describible interface {
	Describir() string
}

// buscar recorre un slice y devuelve el servidor o un error envuelto.
func buscar(lista []Servidor, cod int) (Servidor, error) {
	for _, s := range lista {
		if s.Cod == cod {
			return s, nil
		}
	}
	// %w "envuelve" el error centinela añadiendo contexto
	return Servidor{}, fmt.Errorf("buscar(%d): %w", cod, ErrNoEncontrado)
}

func main() {
	fmt.Println("--- Introducción a Go:", Asignatura, "---")

	// -----------------------------------------------------------------
	seccion("Parte 2. Variables y tipos")
	// -----------------------------------------------------------------
	var usuario string = "system" // declaración completa
	var puerto = 1521             // el tipo (int) se deduce del valor
	sid := "xe"                   // declaración corta, solo dentro de funciones
	var sinValor float64          // sin inicializar: toma su "valor cero" (0)
	activo := true
	host, bd := "localhost", "Empresa" // asignación múltiple

	fmt.Println("usuario:", usuario, "| puerto:", puerto, "| sid:", sid)
	fmt.Println("host:", host, "| base de datos:", bd, "| activo:", activo)
	fmt.Println("valor cero de un float64:", sinValor)
	fmt.Println("motores:", Oracle, Postgres)

	megabits := 300              // int
	uso := float64(megabits) / 7 // conversión SIEMPRE explícita
	fmt.Printf("%d megabits / 7 = %.3f (tipo %T)\n", megabits, uso, uso)
	fmt.Printf("tipos: %T, %T, %T, %T\n", usuario, puerto, uso, activo)

	// -----------------------------------------------------------------
	seccion("Parte 3. Control de flujo")
	// -----------------------------------------------------------------
	// if con sentencia de inicialización: n solo existe dentro del if
	if n := len("12345"); n < 8 {
		fmt.Println("Contraseña corta: tiene", n, "caracteres")
	}

	// for es el ÚNICO bucle de Go. Forma clásica:
	for i := 1; i <= 3; i++ {
		fmt.Print("intento ", i, "  ")
	}
	fmt.Println()

	// Forma "while": solo condición
	pendientes := 3
	for pendientes > 0 {
		pendientes--
	}
	fmt.Println("tuplas pendientes:", pendientes)

	// for range recorre colecciones devolviendo índice y valor
	for i, nombre := range []string{"Jupiter", "Mercurio", "Saturno"} {
		fmt.Printf("  [%d] %s\n", i, nombre)
	}

	// switch: no necesita break y un case admite varios valores
	for _, motor := range []int{Oracle, Postgres} {
		switch motor {
		case Oracle:
			fmt.Println("Oracle usa los marcadores :1, :2 ...")
		case Postgres:
			fmt.Println("PostgreSQL usa los marcadores $1, $2 ...")
		}
	}

	// -----------------------------------------------------------------
	seccion("Parte 4. Slices y mapas")
	// -----------------------------------------------------------------
	// Slice: lista dinámica (equivale a la lista de Python o al ArrayList)
	nombres := []string{"Jupiter", "Mercurio"}
	nombres = append(nombres, "Saturno", "Urano") // añadir elementos
	fmt.Println("servidores:", nombres, "| len =", len(nombres))
	fmt.Println("subslice [1:3]:", nombres[1:3]) // del índice 1 al 2

	// Mapa: tabla clave -> valor (equivale al diccionario de Python)
	consumo := map[string]float64{"123": 300, "234": 5, "AAA": 600}
	consumo["BBB"] = 20    // insertar o modificar
	delete(consumo, "234") // borrar
	fmt.Println("consumo del usuario AAA:", consumo["AAA"])

	// "comma ok": saber si una clave existe
	if mb, ok := consumo["XYZ"]; !ok {
		fmt.Println("XYZ no existe; se obtiene el valor cero:", mb)
	}

	// El orden de recorrido de un mapa NO está garantizado: se ordenan
	// las claves para obtener siempre la misma salida.
	claves := make([]string, 0, len(consumo))
	for k := range consumo {
		claves = append(claves, k)
	}
	sort.Strings(claves)
	for _, k := range claves {
		fmt.Printf("  %s -> %.0f megabits\n", k, consumo[k])
	}

	// -----------------------------------------------------------------
	seccion("Parte 5. Funciones con varios valores de retorno")
	// -----------------------------------------------------------------
	if r, err := dividir(353, 5); err == nil {
		fmt.Printf("353 / 5 = %.1f\n", r)
	}
	if _, err := dividir(1, 0); err != nil { // _ descarta el resultado
		fmt.Println("Error:", err)
	}

	total, media := estadisticas(300, 5, 20, 3, 25)
	fmt.Printf("megabits del servidor 1: total = %.0f, media = %.1f\n", total, media)

	// Función anónima (closure) que recuerda la variable contador
	contador := 0
	siguienteID := func() int { contador++; return contador }
	fmt.Println("IDs generados:", siguienteID(), siguienteID(), siguienteID())

	// -----------------------------------------------------------------
	seccion("Parte 6. Errores y defer")
	// -----------------------------------------------------------------
	// strconv.Atoi devuelve un error si el texto no es un número
	for _, texto := range []string{"7", "siete"} {
		cod, err := strconv.Atoi(texto)
		if err != nil {
			fmt.Println("No es un código válido:", err)
			continue
		}
		fmt.Println("Código leído:", cod)
	}

	fmt.Println("Ejecución correcta:")
	_ = usarRecursos(false)
	fmt.Println("Ejecución con error:")
	if err := usarRecursos(true); err != nil {
		fmt.Println("  error devuelto:", err)
	}

	// -----------------------------------------------------------------
	seccion("Parte 7. Structs, métodos e interfaces")
	// -----------------------------------------------------------------
	servidores := []Servidor{
		{Cod: 1, Nombre: "Jupiter"},
		{Cod: 2, Nombre: "Mercurio"},
		{3, "Saturno", 0, 0}, // también por posición
	}

	s := &servidores[0]       // puntero al primer servidor
	s.RegistrarTrafico(125.5) // el método modifica el original
	fmt.Println(servidores[0].Describir())

	if encontrado, err := buscar(servidores, 2); err == nil {
		fmt.Println("Encontrado:", encontrado.Describir())
	}
	_, err := buscar(servidores, 9)
	if errors.Is(err, ErrNoEncontrado) { // compara con el error envuelto
		fmt.Println("Error controlado:", err)
	}

	// Dos tipos distintos guardados en un slice de la misma interfaz
	u := Usuario{Cod: "123", Nombre: "Mendez", Servidor: 1}
	elementos := []Describible{servidores[1], u}
	for _, e := range elementos {
		fmt.Printf("  %-13T -> %s\n", e, e.Describir())
	}

	fmt.Println("\n--- Fin del script ---")
}