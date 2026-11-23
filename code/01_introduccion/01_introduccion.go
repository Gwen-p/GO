package main

import (
	"errors"
	"fmt"
)

// Función con retorno múltiple y manejo de errores (sustituye a try/except)
func dividir(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("no se puede dividir entre cero")
	}
	return a / b, nil
}

func main() {
	// 1. Variables (declaración explícita e inferida) y E/S
	var nombre string
	edad := 20 // Inferencia de tipo

	fmt.Print("Introduce tu nombre: ")
	fmt.Scanln(&nombre)
	fmt.Print("Introduce tu edad (si pulsas ENTER se supondrá que tienes 20): ")
	fmt.Scanln(&edad)
	fmt.Printf("Hola %s, tienes %d años.\n\n", nombre, edad)

	// 2. Slices (equivalente a listas de Python)
	lenguajes := []string{"Python", "Go", "SQL"}
	lenguajes = append(lenguajes, "Rust")
	fmt.Println("Lenguajes:", lenguajes)

	// 3. Maps (equivalente a diccionarios de Python)
	puertos := map[string]int{
		"PostgreSQL": 5432,
		"Oracle":     1521,
	}
	fmt.Println("Puerto de PostgreSQL:", puertos["PostgreSQL"])

	// 4. Condicionales
	if edad >= 18 {
		fmt.Println("Estado: Mayor de edad")
	} else {
		fmt.Println("Estado: Menor de edad")
	}

	// 5. Bucle for (en Go no existe 'while', todo se hace con 'for')
	fmt.Println("\nRecorriendo el slice:")
	for indice, valor := range lenguajes {
		fmt.Printf("  [%d] %s\n", indice, valor)
	}

	// 6. Llamada a función y comprobación de error
	resultado, err := dividir(10, 0)
	if err != nil {
		fmt.Println("\nError controlado:", err)
	} else {
		fmt.Println("\nResultado:", resultado)
	}
}