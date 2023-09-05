package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	var (
		inputFile   string
		outputFile  string
		inscopeFile string
		outscopeFile string
	)

	// Definir los flags personalizados
	flag.StringVar(&inputFile, "if", "", "Archivo de entrada (inputFile)")
	flag.StringVar(&outputFile, "of", "", "Archivo de salida (outputFile)")
	flag.StringVar(&inscopeFile, "is", "", "Archivo con dominios en inscope")
	flag.StringVar(&outscopeFile, "os", "", "Archivo con dominios en outscope")

	// Parsear los flags
	flag.Parse()

	// Verificar que se han proporcionado los cuatro parámetros
	if inputFile == "" || outputFile == "" || inscopeFile == "" || outscopeFile == "" {
		fmt.Println("Faltan parámetros. Uso: programa -if inputFile -of outputFile -is inscope -os outscope")
		os.Exit(1)
	}

	// Leer los dominios de inscope y outscope
	inscopeDomains := readDomainsFromFile(inscopeFile)
	outscopeDomains := readDomainsFromFile(outscopeFile)

	// Abrir el archivo de entrada para lectura
	input, err := os.Open(inputFile)
	if err != nil {
		fmt.Println("No se puede abrir el archivo de entrada:", err)
		os.Exit(1)
	}
	defer input.Close()

	// Abrir el archivo de salida para escritura
	output, err := os.Create(outputFile)
	if err != nil {
		fmt.Println("No se puede crear el archivo de salida:", err)
		os.Exit(1)
	}
	defer output.Close()

	scanner := bufio.NewScanner(input)

	// Recorrer las líneas del archivo de entrada
	for scanner.Scan() {
		line := scanner.Text()
		// Comprobar si el subdominio está en inscope y no en outscope
		if isInScope(line, inscopeDomains) && !isInScope(line, outscopeDomains) {
			// Escribir el subdominio en el archivo de salida
			fmt.Println(line)
			output.WriteString(line + "\n")
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error al leer el archivo de entrada:", err)
		os.Exit(1)
	}
}

func readDomainsFromFile(filename string) map[string]bool {
	domains := make(map[string]bool)
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("No se puede abrir el archivo %s: %v\n", filename, err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		domains[line] = true
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error al leer el archivo %s: %v\n", filename, err)
		os.Exit(1)
	}

	return domains
}

func isInScope(domain string, scope map[string]bool) bool {
	// Verificar si el dominio o subdominio está en el scope
	for key := range scope {
		if strings.HasSuffix(domain, key) {
			return true
		}
	}
	return false
}

