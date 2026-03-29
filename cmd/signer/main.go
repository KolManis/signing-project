package main

import (
	"fmt"
	"log"
	"os"

	"github.com/KolManis/signing-project/internal/crypto"
)

func main() {
	certFile := "cert.crt"
	keyFile := "private.key"
	inputFile := "document.pdf"

	// Создаем тестовый файл, если его нет
	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		os.WriteFile(inputFile, []byte("Содержимое диплома: Теория криптографии..."), 0644)
		fmt.Println("Создан тестовый файл:", inputFile)
	}

	// Генерация ключей (если их нет)
	if _, err := os.Stat(keyFile); os.IsNotExist(err) {
		fmt.Println("Генерация ключей...")
		if err := crypto.GenerateKeys(certFile, keyFile); err != nil {
			log.Fatal(err)
		}
	}

	// Читаем данные файла для подписи
	data, err := os.ReadFile(inputFile)
	if err != nil {
		log.Fatal(err)
	}

	// Создаем подпись
	fmt.Printf("Подписываем файл %s...\n", inputFile)
	signature, err := crypto.SignDocument(data, certFile, keyFile)
	if err != nil {
		log.Fatalf("Ошибка подписи: %v", err)
	}

	// СОХРАНЯЕМ В .p7s
	signatureFile := inputFile + ".p7s" // Получится document.pdf.p7s
	err = os.WriteFile(signatureFile, signature, 0644)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(" Подпись сохранена в:", signatureFile)

	err = crypto.VerifySignature(data, signature)
	if err != nil {
		fmt.Println(" ОШИБКА: Подпись не прошла проверку")
	} else {
		fmt.Println(" Проверка успешна")
	}
}
