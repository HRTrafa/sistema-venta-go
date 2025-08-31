package utils

import (
	"os"
	"os/exec"
	"runtime"
)

// ClearScreen limpia la pantalla de la consola.
// Funciona en sistemas Unix-like (Linux, macOS) y en Windows.
func ClearScreen() {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}