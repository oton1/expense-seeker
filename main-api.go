package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type Gastos struct {
	ID    string  `json:"id"`
	Valor float64 `json:"valor"`
	Data  string  `json:"data"`
	Desc  string  `json:"desc"`
}

const expensesFile = "expenses.json"

var expenses []Gastos

func main() {
	router := gin.Default()

	if err := loadExpenses(); err != nil {
		fmt.Println("Erro ao carregar despesas:", err)
	}

	router.GET("/expenses", getExpenses)
	router.POST("/expenses", postExpense)

	router.Run("localhost:8080")
}

func getExpenses(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, expenses)
}

func postExpense(c *gin.Context) {
	var newExpense Gastos

	if err := c.ShouldBindJSON(&newExpense); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}
	if newExpense.Valor == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "O campo 'valor' é obrigatório"})
		return
	}
	if newExpense.Data == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "O campo 'data' é obrigatório"})
		return
	}
	if newExpense.Desc == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "O campo 'desc' é obrigatório"})
		return
	}

	newExpense.ID = generateID(newExpense.Data, newExpense.Desc)
	expenses = append(expenses, newExpense)

	if err := saveExpenses(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao salvar despesa"})
		return
	}

	c.JSON(http.StatusCreated, newExpense)
}

func generateID(data string, desc string) string {
	dataToHash := data + desc
	hash := sha256.Sum256([]byte(dataToHash))
	return hex.EncodeToString(hash[:])
}

func loadExpenses() error {
	file, err := os.Open(expensesFile)
	if err != nil {
		if os.IsNotExist(err) {
			expenses = []Gastos{}
			return nil
		}
		return err
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(&expenses); err != nil {
		return err
	}
	return nil
}

func saveExpenses() error {
	file, err := os.Create(expensesFile)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(expenses)
}
