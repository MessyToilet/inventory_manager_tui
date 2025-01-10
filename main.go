package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/rivo/tview"
)

type Item struct {
	Name  string `json:"name"`
	Stock int    `json:"stock"`
}

// Global vars
var (
	inventory     = []Item{}
	inventoryFile = "inventory.json"
)

func loadInventory() {
	// _ throwen var
	// we keep error to kill
	if _, err := os.Stat(inventoryFile); err == nil {
		data, err := os.ReadFile(inventoryFile)
		if err != nil {
			log.Fatal("[!] Error in loadInventory()! - ", err)
		}
		json.Unmarshal(data, &inventory)
	}
}

func saveInventory() {
	// use this for formating, pprint
	data, err := json.MarshalIndent(inventory, "", " ")
	if err != nil {
		log.Fatal("[!] Error in saveInventroy()! - ", err)
	}
	// octal permission
	os.WriteFile(inventoryFile, data, 0644)
}

// pass params in ( )
func deleteItem(index int) {
	if index < 0 || index >= len(inventory) {
		fmt.Println("[!] Invalid index!")
	}
	// cpp like array delete
	// ... unpacks eles to appended 1 by 1
	inventory = append(inventory[:index], inventory[index+1:]...)
}

func main() {
	app := tview.NewApplication()
	loadInventory()
	// obj.meth() same as obj.
	//					  meth()
	inventoryList := tview.NewTextView().
		SetDynamicColors(true).
		SetWordWrap(true)
	inventoryList.SetBorder(true).SetTitle("Inventory Manager")

	// scope func
	refreshInventory := func() {
		inventoryList.Clear()
		if len(inventory) == 0 {
			fmt.Fprintln(inventoryList, "[!] No Items in inventory!")
		} else {
			for i, item := range inventory {
				fmt.Fprintf(inventoryList, "[%d] %s (Stock: %d)\n", i+1, item.Name, item.Stock)
			}
		}
	}

	itemNameInput := tview.NewInputField().SetLabel("Item Name: ")
	itemStockInput := tview.NewInputField().SetLabel("Stock: ")
	itemIDInput := tview.NewInputField().SetLabel("Item ID to del: ")

	form := tview.NewForm().
		AddFormItem(itemNameInput).
		AddFormItem(itemStockInput).
		AddFormItem(itemIDInput).
		AddButton("Add Item", func() {
			name := itemNameInput.GetText()
			stock := itemStockInput.GetText()
			if name != "" && stock != "" {
				quantity, err := strconv.Atoi(stock)
				if err != nil {
					fmt.Fprintln(inventoryList, "[!] Invalid stock value!")
					return
				}
				inventory = append(inventory,
					Item{
						Name:  name,
						Stock: quantity,
					},
				)
				saveInventory()
				refreshInventory()
				itemNameInput.SetText("")
				itemStockInput.SetText("")
			}
		},
		).
		AddButton("Delete Item", func() {
			idStr := itemIDInput.GetText()
			if idStr == "" {
				fmt.Fprintln(inventoryList, "ID to delete.")
				return
			}
			id, err := strconv.Atoi(idStr)
			if err != nil || id < 1 || id > len(inventory) {
				fmt.Fprintln(inventoryList, "[!] Invalid ID!")
				return
			}
			deleteItem(id - 1)
			fmt.Fprintf(inventoryList, "Item [%d] deleted.\n", id)
			refreshInventory()
			itemIDInput.SetText("")
		}).
		AddButton("Exit", func() {
			app.Stop()
		})
	form.SetBorder(true).SetTitle("Manage Inventory").SetTitleAlign(tview.AlignLeft)
	flex := tview.NewFlex().
		AddItem(inventoryList, 0, 1, false).
		AddItem(form, 0, 1, true)
	refreshInventory()

	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}
