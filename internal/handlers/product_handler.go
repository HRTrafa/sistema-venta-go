package handlers

import (
	"bufio"
	"fmt"
	"os"
	"sales-system/internal/models"
	"sales-system/internal/repository"
	"strconv"
	"strings"
	"time"
)

// RegistrarProducto maneja la opción para registrar un nuevo producto.
func RegisterProduct(productRepo *repository.ProductRepo) {
	reader := bufio.NewReader(os.Stdin)
	
	fmt.Println("\n--- Registrar Producto ---")

	fmt.Print("Fecha (DD/MM/YYYY): ")
	dateStr, _ := reader.ReadString('\n')
	date, err := time.Parse("02/01/2006", strings.TrimSpace(dateStr))
	if err != nil {
		fmt.Println("Formato de fecha inválido. Usando la fecha actual.")
		date = time.Now()
	}

	fmt.Print("Nombre del Producto: ")
	productName, _ := reader.ReadString('\n')
	productName = strings.TrimSpace(productName)

	fmt.Print("Cantidad Inicial: ")
	quantityStr, _ := reader.ReadString('\n')
	quantity, err := strconv.Atoi(strings.TrimSpace(quantityStr))
	if err != nil {
		fmt.Println("Cantidad inválida. Usando 0.")
		quantity = 0
	}

	fmt.Print("Precio: ")
	priceStr, _ := reader.ReadString('\n')
	price, err := strconv.ParseFloat(strings.TrimSpace(priceStr), 64)
	if err != nil {
		fmt.Println("Precio inválido. Usando 0.0.")
		price = 0.0
	}

	product := models.Product{
		Date:     date,
		Name:     productName,
		Quantity: quantity,
		Price:    price,
	}

	_, err = productRepo.CreateProduct(product)
	if err != nil {
		fmt.Println("Error al registrar el producto:", err)
		return
	}

	fmt.Println("Producto registrado con éxito.")
}

// ShowProducts visualiza todos los productos registrados en una tabla.
func ShowProducts(productRepo *repository.ProductRepo) {
	products, err := productRepo.GetAllProducts()
	if err != nil {
		fmt.Println("Error al obtener productos:", err)
		return
	}

	fmt.Println("\n--- Listado de Productos ---")
	fmt.Printf("%-5s | %-12s | %-20s | %-8s | %-8s\n", "ID", "Fecha", "Producto", "Cantidad", "Precio")
	fmt.Println("------------------------------------------------------------------")
	for _, p := range products {
		fmt.Printf("%-5d | %-12s | %-20s | %-8d | %-8.2f\n", p.ID, p.Date.Format("02/01/2006"), p.Name, p.Quantity, p.Price)
	}
}

// EditProduct maneja la edición de los datos de un producto.
func EditProduct(productRepo *repository.ProductRepo) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\n--- Editar Producto ---")
	ShowProducts(productRepo) // Muestra la lista para que el usuario elija un ID

	fmt.Print("\nIngrese el ID del producto a editar: ")
	idStr, _ := reader.ReadString('\n')
	id, err := strconv.Atoi(strings.TrimSpace(idStr))
	if err != nil {
		fmt.Println("ID inválido.")
		return
	}

	product, err := productRepo.GetProductByID(id)
	if err != nil {
		fmt.Println("Producto no encontrado.")
		return
	}
	
	fmt.Println("\nDeje los campos en blanco para mantener el valor actual.")

	fmt.Printf("Fecha (actual: %s): ", product.Date.Format("02/01/2006"))
	dateStr, _ := reader.ReadString('\n')
	if strings.TrimSpace(dateStr) != "" {
		newDate, err := time.Parse("02/01/2006", strings.TrimSpace(dateStr))
		if err == nil {
			product.Date = newDate
		}
	}

	fmt.Printf("Cantidad (actual: %d): ", product.Quantity)
	quantityStr, _ := reader.ReadString('\n')
	if strings.TrimSpace(quantityStr) != "" {
		newQuantity, err := strconv.Atoi(strings.TrimSpace(quantityStr))
		if err == nil {
			product.Quantity = newQuantity
		}
	}

	fmt.Printf("Precio (actual: %.2f): ", product.Price)
	priceStr, _ := reader.ReadString('\n')
	if strings.TrimSpace(priceStr) != "" {
		newPrice, err := strconv.ParseFloat(strings.TrimSpace(priceStr), 64)
		if err == nil {
			product.Price = newPrice
		}
	}

	err = productRepo.UpdateProduct(*product)
	if err != nil {
		fmt.Println("Error al actualizar el producto:", err)
		return
	}
	fmt.Println("Producto actualizado con éxito.")
}

// DeleteProduct maneja la eliminación de un producto.
func DeleteProduct(productRepo *repository.ProductRepo) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\n--- Eliminar Producto ---")
	ShowProducts(productRepo) // Muestra la lista para que el usuario elija un ID

	fmt.Print("\nIngrese el ID del producto a eliminar: ")
	idStr, _ := reader.ReadString('\n')
	id, err := strconv.Atoi(strings.TrimSpace(idStr))
	if err != nil {
		fmt.Println("ID inválido.")
		return
	}

	fmt.Print("¿Está seguro de que desea eliminar este producto? (s/n): ")
	confirmation, _ := reader.ReadString('\n')
	if strings.ToLower(strings.TrimSpace(confirmation)) != "s" {
		fmt.Println("Operación cancelada.")
		return
	}

	err = productRepo.DeleteProduct(id)
	if err != nil {
		fmt.Println("Error al eliminar el producto:", err)
		return
	}
	fmt.Println("Producto eliminado con éxito.")
}