package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/rivo/tview"
)

// Define a estrutura do item que será armazenado no inventário
type Item struct {
	Name  string `json:"name"`  //nome do item
	Stock int    `json:"value"` //quantidade em estoque
}

var (
	inventory    = []Item{}
	invetoryFile = "inventory.json"
)

// 1 carregar o inventário from json file
func loadInventory() {
	if _, err := os.Stat(invetoryFile); err == nil {
		data, err := os.ReadFile(invetoryFile)
		if err != nil {
			log.Fatal("Error reading inventory file! - ", err)
		}
		json.Unmarshal(data, &inventory)
	}
}

// 2 save inventory function
func saveInventory() {
	data, err := json.MarshalIndent(inventory, "", " ")
	if err != nil {
		log.Fatal("Erro saving inventory! - ", err)
	}

	os.WriteFile(invetoryFile, data, 0644)
}

// 3 delete item
func deleteItem(index int) {
	if index < 0 || index >= len(inventory) {
		fmt.Println("Invalid item index!")
		return
	}

	inventory = append(inventory[:index], inventory[index+1:]...)
	saveInventory()
}

func main() {
	app := tview.NewApplication() // TUI application
	loadInventory()
	inventoryList := tview.NewTextView().
		SetDynamicColors(true).
		SetWordWrap(true)

	inventoryList.SetBorder(true).SetTitle("Inventory Items")

	// refreshes inventory display
	refreshInventory := func() {
		inventoryList.Clear() // Clear

		if len(inventory) == 0 {
			fmt.Fprintln(inventoryList, "No items in inventory.")
		} else {
			// Iterate through inventory
			for i, item := range inventory {
				fmt.Fprintf(inventoryList, "[%d] %s (Stock: %d)\n", i+1, item.Name, item.Stock)
			}
		}
	}

	// input fields item and quantity
	itemNameInput := tview.NewInputField().SetLabel("Item Name: ")
	itemStockInput := tview.NewInputField().SetLabel("Stock: ")

	// input field deleting
	itemIDInput := tview.NewInputField().SetLabel("Item ID to delete: ")

	// Create a form that lets the user add or delete items
	form := tview.NewForm().
		AddFormItem(itemNameInput).    // Add item name input to the form
		AddFormItem(itemStockInput).   // Add item stock input to the form
		AddFormItem(itemIDInput).      // Add item ID input for deletion
		AddButton("Add Item", func() { // Button to add a new item
			// Get the text input for name and stock
			name := itemNameInput.GetText()
			stock := itemStockInput.GetText()
			// Check if both fields are filled
			if name != "" && stock != "" {
				// Convert the stock input to an integer
				quantity, err := strconv.Atoi(stock)
				if err != nil {
					fmt.Fprintln(inventoryList, "Invalid stock value.")
					return
				}
				// Add the new item to the inventory slice
				inventory = append(inventory, Item{Name: name, Stock: quantity})
				// Save the updated inventory
				saveInventory()
				// Refresh the inventory display
				refreshInventory()
				// Clear the input fields after adding the item
				itemNameInput.SetText("")
				itemStockInput.SetText("")
			}
		}).
		AddButton("Delete Item", func() { // Button to delete an item
			idStr := itemIDInput.GetText()
			// Ensure the ID field is not empty
			if idStr == "" {
				fmt.Fprintln(inventoryList, "Please enter an item ID to delete.")
				return
			}
			// Convert the ID to an integer and check if it's valid
			id, err := strconv.Atoi(idStr)
			if err != nil || id < 1 || id > len(inventory) {
				fmt.Fprintln(inventoryList, "Invalid item ID.")
				return
			}
			// Delete the item (adjust for zero-based index)
			deleteItem(id - 1)
			fmt.Fprintf(inventoryList, "Item [%d] deleted.\n", id)
			// Refresh the inventory display after deletion
			refreshInventory()
			itemIDInput.SetText("") // Clear the ID input field
		}).
		AddButton("Exit", func() { // Button to exit the application
			app.Stop()
		})

	// Set a border and title for the form
	form.SetBorder(true).SetTitle("Manage Inventory").SetTitleAlign(tview.AlignLeft)

	// Create a layout using Flex to display the inventory list and the form side by side
	flex := tview.NewFlex().
		AddItem(inventoryList, 0, 1, false). // Left side: inventory list
		AddItem(form, 0, 1, true)            // Right side: form for adding/deleting items

	// Initial inventory display
	refreshInventory()

	// Start TUI application
	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}
