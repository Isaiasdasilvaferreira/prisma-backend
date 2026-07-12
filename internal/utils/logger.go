package utils

import (
	"fmt"
	"os"
	"time"
)

func LogError(message string, err error) {
	// Cria a pasta logs se não existir
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		os.Mkdir("logs", 0755)
	}

	// Abre o arquivo de log
	file, err := os.OpenFile("logs/error.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("❌ Erro ao abrir arquivo de log: %v\n", err)
		return
	}
	defer file.Close()

	// Escreve o erro no arquivo
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logEntry := fmt.Sprintf("[%s] %s: %v\n", timestamp, message, err)
	
	if _, err := file.WriteString(logEntry); err != nil {
		fmt.Printf("❌ Erro ao escrever no arquivo de log: %v\n", err)
	}
}

func LogInfo(message string) {
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		os.Mkdir("logs", 0755)
	}

	file, err := os.OpenFile("logs/info.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("❌ Erro ao abrir arquivo de log: %v\n", err)
		return
	}
	defer file.Close()

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logEntry := fmt.Sprintf("[%s] %s\n", timestamp, message)
	
	if _, err := file.WriteString(logEntry); err != nil {
		fmt.Printf("❌ Erro ao escrever no arquivo de log: %v\n", err)
	}
}

func LogData(data interface{}) {
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		os.Mkdir("logs", 0755)
	}

	file, err := os.OpenFile("logs/data.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("❌ Erro ao abrir arquivo de log: %v\n", err)
		return
	}
	defer file.Close()

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logEntry := fmt.Sprintf("[%s] %+v\n", timestamp, data)
	
	if _, err := file.WriteString(logEntry); err != nil {
		fmt.Printf("❌ Erro ao escrever no arquivo de log: %v\n", err)
	}
}
