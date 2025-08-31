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

// PreviewProducts muestra una lista simple de productos (ID y Nombre) para referencia.
func PreviewProducts(productRepo *repository.ProductRepo) {
	products, err := productRepo.GetAllProducts()
	if err != nil {
		fmt.Println("Error al obtener productos:", err)
		return
	}

	fmt.Println("\n--- Productos Registrados (ID y Nombre) ---")
	fmt.Printf("%-5s | %-20s\n", "ID", "Producto")
	fmt.Println("-------------------------------")
	for _, p := range products {
		fmt.Printf("%-5d | %-20s\n", p.ID, p.Name)
	}
	fmt.Println("-------------------------------")
}

// RegisterSale maneja la lógica para registrar una nueva venta.
func RegisterSale(saleRepo *repository.SaleRepo, productRepo *repository.ProductRepo) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\n--- Registrar Venta ---")

	// Previsualización de productos antes de pedir el ID.
	PreviewProducts(productRepo)

	fmt.Print("Fecha (DD/MM/YYYY): ")
	dateStr, _ := reader.ReadString('\n')
	date, err := time.Parse("02/01/2006", strings.TrimSpace(dateStr))
	if err != nil {
		fmt.Println("Formato de fecha inválido. Usando la fecha actual.")
		date = time.Now()
	}

	fmt.Print("Nombre del Cliente: ")
	client, _ := reader.ReadString('\n')
	client = strings.TrimSpace(client)

	fmt.Print("ID del Producto: ")
	productIDStr, _ := reader.ReadString('\n')
	productID, err := strconv.Atoi(strings.TrimSpace(productIDStr))
	if err != nil {
		fmt.Println("ID de producto inválido.")
		return
	}

	product, err := productRepo.GetProductByID(productID)
	if err != nil {
		fmt.Println("Error: Producto no encontrado. Verifique el ID.")
		return
	}

	fmt.Print("Cantidad: ")
	quantityStr, _ := reader.ReadString('\n')
	quantity, err := strconv.Atoi(strings.TrimSpace(quantityStr))
	if err != nil {
		fmt.Println("Cantidad inválida. Operación cancelada.")
		return
	}

	// El sistema multiplica la cantidad por el precio para obtener el total.
	total := float64(quantity) * product.Price

	fmt.Print("Estado (1. Pagado, 2. Pendiente): ")
	statusChoiceStr, _ := reader.ReadString('\n')
	statusChoice, err := strconv.Atoi(strings.TrimSpace(statusChoiceStr))
	if err != nil || (statusChoice != 1 && statusChoice != 2) {
		fmt.Println("Opción de estatus inválida. Usando 'Pendiente'.")
		statusChoice = 2
	}
	status := models.StatusPending
	if statusChoice == 1 {
		status = models.StatusPaid
	}

	newSale := models.Sale{
		Date:      date,
		Client:    client,
		ProductID: productID,
		Quantity:  quantity,
		Price:     product.Price,
		Total:     total,
		Status:    status,
	}

	id, err := saleRepo.CreateSale(newSale)
	if err != nil {
		fmt.Println("Error al registrar la venta:", err)
		return
	}

	fmt.Printf("Venta registrada con éxito. ID: %d\n", id)
}

// ShowSales visualiza todas las ventas registradas o los detalles de una venta específica.
func ShowSales(saleRepo *repository.SaleRepo, productRepo *repository.ProductRepo) {
	reader := bufio.NewReader(os.Stdin)
	
	sales, err := saleRepo.GetAllSales()
	if err != nil {
		fmt.Println("Error al obtener las ventas:", err)
		return
	}

	fmt.Println("\n--- Listado de Ventas ---")
	fmt.Printf("%-5s | %-12s | %-20s | %-10s | %-8s\n", "ID", "Fecha", "Cliente", "Cantidad", "Estado")
	fmt.Println("------------------------------------------------------------------")
	for _, sale := range sales {
		fmt.Printf("%-5d | %-12s | %-20s | %-10d | %-8s\n", sale.ID, sale.Date.Format("02/01/2006"), sale.Client, sale.Quantity, sale.Status)
	}

	fmt.Print("\nIngrese el ID de la venta para ver detalles completos (o presione Enter para volver): ")
	idStr, _ := reader.ReadString('\n')
	idStr = strings.TrimSpace(idStr)

	if idStr != "" {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			fmt.Println("ID inválido.")
			return
		}
		
		sale, err := saleRepo.GetSaleByID(id)
		if err != nil {
			fmt.Println("Venta no encontrada.")
			return
		}

		product, _ := productRepo.GetProductByID(sale.ProductID)
		productName := "N/A"
		if product != nil {
			productName = product.Name
		}
		
		fmt.Println("\n--- Detalles de la Venta ---")
		fmt.Println("ID:", sale.ID)
		fmt.Println("Fecha:", sale.Date.Format("02/01/2006"))
		fmt.Println("Cliente:", sale.Client)
		fmt.Println("Producto ID:", sale.ProductID)
		fmt.Println("Nombre del Producto:", productName)
		fmt.Println("Cantidad:", sale.Quantity)
		fmt.Println("Precio Unitario:", sale.Price)
		fmt.Println("Total:", sale.Total)
		fmt.Println("Estatus:", sale.Status)
	}
}

// EditSale maneja la edición de los datos de una venta.
func EditSale(saleRepo *repository.SaleRepo, productRepo *repository.ProductRepo) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\n--- Editar Venta ---")
	ShowSales(saleRepo, productRepo)

	fmt.Print("\nIngrese el ID de la venta a editar: ")
	idStr, _ := reader.ReadString('\n')
	id, err := strconv.Atoi(strings.TrimSpace(idStr))
	if err != nil {
		fmt.Println("ID inválido.")
		return
	}

	sale, err := saleRepo.GetSaleByID(id)
	if err != nil {
		fmt.Println("Venta no encontrada.")
		return
	}

	fmt.Println("\nDeje los campos en blanco para mantener el valor actual.")

	fmt.Printf("Fecha (actual: %s): ", sale.Date.Format("02/01/2006"))
	dateStr, _ := reader.ReadString('\n')
	if strings.TrimSpace(dateStr) != "" {
		newDate, err := time.Parse("02/01/2006", strings.TrimSpace(dateStr))
		if err == nil {
			sale.Date = newDate
		}
	}

	fmt.Printf("Nombre del Cliente (actual: %s): ", sale.Client)
	clientStr, _ := reader.ReadString('\n')
	if strings.TrimSpace(clientStr) != "" {
		sale.Client = strings.TrimSpace(clientStr)
	}
	
	fmt.Printf("Cantidad (actual: %d): ", sale.Quantity)
	quantityStr, _ := reader.ReadString('\n')
	if strings.TrimSpace(quantityStr) != "" {
		newQuantity, err := strconv.Atoi(strings.TrimSpace(quantityStr))
		if err == nil {
			sale.Quantity = newQuantity
		}
	}

	fmt.Printf("Precio (actual: %.2f): ", sale.Price)
	priceStr, _ := reader.ReadString('\n')
	if strings.TrimSpace(priceStr) != "" {
		newPrice, err := strconv.ParseFloat(strings.TrimSpace(priceStr), 64)
		if err == nil {
			sale.Price = newPrice
		}
	}

	// Recalcular el total si la cantidad o el precio cambiaron
	sale.Total = float64(sale.Quantity) * sale.Price

	fmt.Printf("Estado (actual: %s - 1. Pagado, 2. Pendiente): ", sale.Status)
	statusChoiceStr, _ := reader.ReadString('\n')
	if strings.TrimSpace(statusChoiceStr) != "" {
		newStatusChoice, err := strconv.Atoi(strings.TrimSpace(statusChoiceStr))
		if err == nil && (newStatusChoice == 1 || newStatusChoice == 2) {
			if newStatusChoice == 1 {
				sale.Status = models.StatusPaid
			} else {
				sale.Status = models.StatusPending
			}
		}
	}

	err = saleRepo.UpdateSale(*sale)
	if err != nil {
		fmt.Println("Error al actualizar la venta:", err)
		return
	}
	fmt.Println("Venta actualizada con éxito.")
}

// DeleteSale maneja la eliminación de una venta.
func DeleteSale(saleRepo *repository.SaleRepo) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\n--- Eliminar Venta ---")
	sales, err := saleRepo.GetAllSales()
	if err != nil {
		fmt.Println("Error al obtener las ventas:", err)
		return
	}
	// Muestra una lista simple para facilitar la elección del usuario
	fmt.Printf("%-5s | %-12s | %-20s\n", "ID", "Fecha", "Cliente")
	fmt.Println("--------------------------------------")
	for _, s := range sales {
		fmt.Printf("%-5d | %-12s | %-20s\n", s.ID, s.Date.Format("02/01/2006"), s.Client)
	}

	fmt.Print("\nIngrese el ID de la venta a eliminar: ")
	idStr, _ := reader.ReadString('\n')
	id, err := strconv.Atoi(strings.TrimSpace(idStr))
	if err != nil {
		fmt.Println("ID inválido.")
		return
	}

	fmt.Print("¿Está seguro de que desea eliminar esta venta? (s/n): ")
	confirmation, _ := reader.ReadString('\n')
	if strings.ToLower(strings.TrimSpace(confirmation)) != "s" {
		fmt.Println("Operación cancelada.")
		return
	}

	err = saleRepo.DeleteSale(id)
	if err != nil {
		fmt.Println("Error al eliminar la venta:", err)
		return
	}
	fmt.Println("Venta eliminada con éxito.")
}