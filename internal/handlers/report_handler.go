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

	"github.com/jung-kurt/gofpdf"
)

// GenerateReport maneja la generación de reportes diarios, semanales o mensuales.
func GenerateReport(saleRepo *repository.SaleRepo, cashRepo *repository.CashDeliveryRepo, productRepo *repository.ProductRepo) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("\n--- Reportes de Ventas ---")
	fmt.Println("1. Diario")
	fmt.Println("2. Semanal")
	fmt.Println("3. Mensual")
	fmt.Print("Seleccione un tipo de reporte: ")
	choiceStr, _ := reader.ReadString('\n')
	choice, _ := strconv.Atoi(strings.TrimSpace(choiceStr))

	var sales []models.Sale
	var deliveries []models.CashDelivery
	var reportTitle string

	now := time.Now()
	var start, end time.Time

	switch choice {
	case 1: // Diario
		reportTitle = "Reporte de Ventas Diario"
		start = now.Truncate(24 * time.Hour)
		end = start.Add(24*time.Hour - time.Second)
	case 2: // Semanal
		reportTitle = "Reporte de Ventas Semanal"
		weekday := int(now.Weekday())
		start = now.AddDate(0, 0, -weekday).Truncate(24 * time.Hour)
		end = start.AddDate(0, 0, 7).Add(-time.Second)
	case 3: // Mensual
		reportTitle = "Reporte de Ventas Mensual"
		start = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		end = start.AddDate(0, 1, 0).Add(-time.Second)
	default:
		fmt.Println("Opción no válida.")
		return
	}

	sales, err := saleRepo.GetSalesByDateRange(start, end)
	if err != nil {
		fmt.Println("Error al obtener ventas para el reporte:", err)
		return
	}

	deliveries, err = cashRepo.GetCashDeliveriesByDateRange(start, end)
	if err != nil {
		fmt.Println("Error al obtener entregas de dinero para el reporte:", err)
		return
	}

	// Sumatoria de todos los totales
	var totalSalesAmount float64
	var totalProductsSold int
	for _, s := range sales {
		totalSalesAmount += s.Total
		totalProductsSold += s.Quantity
	}
	var totalCashDelivered float64
	for _, d := range deliveries {
		totalCashDelivered += d.Amount
	}

	// Mostrar reporte en consola
	fmt.Printf("\n--- %s ---\n", reportTitle)
	fmt.Printf("Período: %s a %s\n", start.Format("02/01/2006"), end.Format("02/01/2006"))
	
	// Tabla de ventas
	fmt.Println("\nDetalles de Ventas:")
	fmt.Printf("%-5s | %-12s | %-20s | %-10s | %-8s | %-8s\n", "ID", "Fecha", "Cliente", "Producto", "Cantidad", "Total")
	fmt.Println("-------------------------------------------------------------------------------------")
	for _, s := range sales {
		product, _ := productRepo.GetProductByID(s.ProductID)
		productName := "N/A"
		if product != nil {
			productName = product.Name
		}
		fmt.Printf("%-5d | %-12s | %-20s | %-10s | %-8d | %-8.2f\n", s.ID, s.Date.Format("02/01/2006"), s.Client, productName, s.Quantity, s.Total)
	}

	// Resumen del reporte
	fmt.Println("\n--- Resumen del Reporte ---")
	fmt.Printf("Total de Ventas: %.2f\n", totalSalesAmount)
	fmt.Printf("Total de Productos Vendidos: %d\n", totalProductsSold)
	fmt.Printf("Total de Dinero Entregado: %.2f\n", totalCashDelivered)
	
	neto := totalSalesAmount - totalCashDelivered
	fmt.Printf("Neto (Ventas - Entregas): %.2f\n", neto)

	fmt.Print("\n¿Desea exportar este reporte a PDF? (s/n): ")
	exportChoice, _ := reader.ReadString('\n')
	if strings.ToLower(strings.TrimSpace(exportChoice)) == "s" {
		ExportReportToPDF(reportTitle, start, end, sales, deliveries, totalSalesAmount, totalProductsSold, totalCashDelivered, productRepo)
		fmt.Println("Reporte exportado a PDF con éxito.")
	}
}

// ExportReportToPDF genera y guarda un archivo PDF del reporte.
func ExportReportToPDF(title string, start, end time.Time, sales []models.Sale, deliveries []models.CashDelivery, totalSalesAmount float64, totalProductsSold int, totalCashDelivered float64, productRepo *repository.ProductRepo) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	
	// Título del PDF
	pdf.Cell(40, 10, title)
	pdf.Ln(8)
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, fmt.Sprintf("Periodo: %s a %s", start.Format("02/01/2006"), end.Format("02/01/2006")))
	pdf.Ln(12)

	// Encabezados de la tabla
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(15, 7, "ID")
	pdf.Cell(25, 7, "Fecha")
	pdf.Cell(40, 7, "Cliente")
	pdf.Cell(30, 7, "Producto")
	pdf.Cell(20, 7, "Cantidad")
	pdf.Cell(20, 7, "Total")
	pdf.Cell(20, 7, "Estado")
	pdf.Ln(-1)

	// Líneas de la tabla
	pdf.SetFont("Arial", "", 10)
	for _, s := range sales {
		product, _ := productRepo.GetProductByID(s.ProductID)
		productName := "N/A"
		if product != nil {
			productName = product.Name
		}
		pdf.Cell(15, 7, strconv.Itoa(s.ID))
		pdf.Cell(25, 7, s.Date.Format("02/01/2006"))
		pdf.Cell(40, 7, s.Client)
		pdf.Cell(30, 7, productName)
		pdf.Cell(20, 7, strconv.Itoa(s.Quantity))
		pdf.Cell(20, 7, fmt.Sprintf("%.2f", s.Total))
		pdf.Cell(20, 7, s.Status)
		pdf.Ln(-1)
	}

	pdf.Ln(10) // Espacio entre la tabla y el resumen

	// Resumen de la tabla
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(50, 7, "Resumen:")
	pdf.Ln(-1)
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(50, 7, fmt.Sprintf("Total de Ventas: %.2f", totalSalesAmount))
	pdf.Ln(-1)
	pdf.Cell(50, 7, fmt.Sprintf("Total de Productos Vendidos: %d", totalProductsSold))
	pdf.Ln(-1)
	pdf.Cell(50, 7, fmt.Sprintf("Total de Dinero Entregado: %.2f", totalCashDelivered))
	pdf.Ln(-1)
	pdf.Cell(50, 7, fmt.Sprintf("Neto (Ventas - Entregas): %.2f", totalSalesAmount-totalCashDelivered))

	// Guardar el PDF
	fileName := strings.ReplaceAll(title, " ", "_") + "_" + time.Now().Format("2006-01-02") + ".pdf"
	err := pdf.OutputFileAndClose(fileName)
	if err != nil {
		fmt.Printf("Error al crear el archivo PDF: %v\n", err)
	}
}