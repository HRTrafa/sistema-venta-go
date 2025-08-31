package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sales-system/internal/database"
	"sales-system/internal/handlers"
	"sales-system/internal/repository"
	"sales-system/internal/utils"
)

func main() {
	// Inicializar la base de datos y crear las tablas si no existen.
	// La base de datos se guardará en un archivo llamado "sales.db".
	database.InitDB("sales.db")

	// Crear las instancias de los repositorios que interactuarán con la base de datos.
	// Se pasa la conexión a la base de datos (database.DB) a cada repositorio.
	productRepo := repository.NewProductRepo(database.DB)
	saleRepo := repository.NewSaleRepo(database.DB)
	cashRepo := repository.NewCashDeliveryRepo(database.DB)

	var choice int
	reader := bufio.NewReader(os.Stdin)

	for {
		// Limpiar la pantalla para una mejor experiencia de usuario.
		utils.ClearScreen()
		showMainMenu()
		
		// Leer la opción del usuario del menú principal.
		choiceStr, _ := reader.ReadString('\n')
		choice, _ = strconv.Atoi(strings.TrimSpace(choiceStr))

		// Usar un switch para dirigir el flujo del programa según la elección del usuario.
		switch choice {
		case 1:
			handleSalesMenu(saleRepo, productRepo)
		case 2:
			handleProductsMenu(productRepo)
		case 3:
			handlers.RegisterCashDelivery(cashRepo)
		case 4:
			handlers.GenerateReport(saleRepo, cashRepo, productRepo)
		case 5:
			fmt.Println("Saliendo del sistema...")
			return
		default:
			fmt.Println("Opción no válida. Por favor, intente de nuevo.")
			fmt.Print("Presione Enter para continuar...")
			reader.ReadString('\n')
		}
	}
}

// showMainMenu muestra el menú principal en la consola.
func showMainMenu() {
	fmt.Println("\n--- Menú Principal ---")
	fmt.Println("1. VENTAS")
	fmt.Println("2. PRODUCTOS")
	fmt.Println("3. Entregas de dinero")
	fmt.Println("4. Reporte de Ventas")
	fmt.Println("5. Salir")
	fmt.Print("Seleccione una opción: ")
}

// handleSalesMenu maneja el submenú de ventas.
func handleSalesMenu(saleRepo *repository.SaleRepo, productRepo *repository.ProductRepo) {
	reader := bufio.NewReader(os.Stdin)
	for {
		utils.ClearScreen()
		fmt.Println("\n--- Menú de Ventas ---")
		fmt.Println("1. Registrar Venta")
		fmt.Println("2. Mostrar Ventas")
		fmt.Println("3. Editar Venta")
		fmt.Println("4. Eliminar Venta")
		fmt.Println("5. Volver al Menú Principal")
		fmt.Print("Seleccione una opción: ")

		choiceStr, _ := reader.ReadString('\n')
		choice, _ := strconv.Atoi(strings.TrimSpace(choiceStr))

		switch choice {
		case 1:
			handlers.RegisterSale(saleRepo, productRepo)
		case 2:
			handlers.ShowSales(saleRepo, productRepo)
		case 3:
			handlers.EditSale(saleRepo, productRepo)
		case 4:
			handlers.DeleteSale(saleRepo)
		case 5:
			return
		default:
			fmt.Println("Opción no válida.")
		}
		fmt.Print("Presione Enter para continuar...")
		reader.ReadString('\n')
	}
}

// handleProductsMenu maneja el submenú de productos.
func handleProductsMenu(productRepo *repository.ProductRepo) {
	reader := bufio.NewReader(os.Stdin)
	for {
		utils.ClearScreen()
		fmt.Println("\n--- Menú de Productos ---")
		fmt.Println("1. Registrar Producto")
		fmt.Println("2. Mostrar productos")
		fmt.Println("3. Editar Producto")
		fmt.Println("4. Eliminar Producto")
		fmt.Println("5. Volver al Menú Principal")
		fmt.Print("Seleccione una opción: ")

		choiceStr, _ := reader.ReadString('\n')
		choice, _ := strconv.Atoi(strings.TrimSpace(choiceStr))

		switch choice {
		case 1:
			handlers.RegisterProduct(productRepo)
		case 2:
			handlers.ShowProducts(productRepo)
		case 3:
			handlers.EditProduct(productRepo)
		case 4:
			handlers.DeleteProduct(productRepo)
		case 5:
			return
		default:
			fmt.Println("Opción no válida.")
		}
		fmt.Print("Presione Enter para continuar...")
		reader.ReadString('\n')
	}
}