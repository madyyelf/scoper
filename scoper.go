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
		inscopeFile  string
		outscopeFile string
	)

	// Definir los flags personalizados
	flag.StringVar(&inscopeFile, "is", "scope.cfg", "Archivo con dominios en inscope (por defecto: scope.cfg)")
	flag.StringVar(&outscopeFile, "os", "outscope.cfg", "Archivo con dominios en outscope (por defecto: outscope.cfg)")

	// Parsear los flags
	flag.Parse()

	// Leer los dominios de inscope y outscope
	inscopeDomains := readDomainsFromFile(inscopeFile)
	outscopeDomains := readDomainsFromFile(outscopeFile)

	// Crear un scanner para leer la entrada estándar (pipe)
	scanner := bufio.NewScanner(os.Stdin)

	// Recorrer las líneas de la entrada estándar
	for scanner.Scan() {
		line := scanner.Text()
		// Comprobar si el subdominio está en inscope y no en outscope
		if isInScope(line, inscopeDomains) && !isInScope(line, outscopeDomains) {
			// Escribir el subdominio en la salida estándar
			fmt.Println(line)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error al leer la entrada estándar:", err)
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
