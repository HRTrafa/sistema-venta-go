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

// RegisterCashDelivery maneja la lógica para registrar una entrega de dinero.
func RegisterCashDelivery(cashRepo *repository.CashDeliveryRepo) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\n--- Registrar Entrega de Dinero ---")

	fmt.Print("Fecha (DD/MM/YYYY): ")
	dateStr, _ := reader.ReadString('\n')
	date, err := time.Parse("02/01/2006", strings.TrimSpace(dateStr))
	if err != nil {
		fmt.Println("Formato de fecha inválido. Usando la fecha actual.")
		date = time.Now()
	}

	fmt.Print("Nombre: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	fmt.Print("Descripción: ")
	description, _ := reader.ReadString('\n')
	description = strings.TrimSpace(description)

	fmt.Print("Monto: ")
	amountStr, _ := reader.ReadString('\n')
	amount, err := strconv.ParseFloat(strings.TrimSpace(amountStr), 64)
	if err != nil {
		fmt.Println("Monto inválido. Operación cancelada.")
		return
	}

	cashDelivery := models.CashDelivery{
		Date:        date,
		Name:        name,
		Description: description,
		Amount:      amount,
	}

	_, err = cashRepo.CreateCashDelivery(cashDelivery)
	if err != nil {
		fmt.Println("Error al registrar la entrega de dinero:", err)
		return
	}

	fmt.Println("Entrega de dinero registrada con éxito.")
}