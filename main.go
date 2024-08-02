package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Product struct {
	ID       int
	Name     string
	Quantity int
	Price    float64
}

type Inventory interface {
	AddProduct(p Product) error
	GetProduct(id int) (Product, error)
	UpdateProduct(id int, newId *int, newName *string, newQuantity *int, newPrice *float64) error
	DeleteProduct(id int) (string, error)
	ListProducts() ([]Product, error)
}

type items struct {
	products map[int]Product
}

func newInventory() *items {
	return &items{
		products: make(map[int]Product),
	}
}

func (inv *items) AddProduct(p Product) error {
	if _, exists := inv.products[p.ID]; exists {
		return fmt.Errorf("Product ID #%d already exists", p.ID)
	}
	inv.products[p.ID] = p
	return nil
}

func (inv *items) GetProduct(id int) (Product, error) {
	if p, exists := inv.products[id]; exists {
		return p, nil
	}
	return Product{}, fmt.Errorf("Product with ID#:%d not found in our items", id)
}

func (inv *items) validateNewID(newID int) error {
	if _, exists := inv.products[newID]; exists {
		return fmt.Errorf("Product ID %d already exists", newID)
	}
	return nil
}

func (inv *items) UpdateProduct(id int, newId *int, newName *string, newQuantity *int, newPrice *float64) error {
	product, exists := inv.products[id]
	if !exists {
		return fmt.Errorf("Product ID #%d not found", id)
	}

	if newId != nil {
		if err := inv.validateNewID(*newId); err != nil {
			return err
		}
		delete(inv.products, id)
		product.ID = *newId
	}

	if newName != nil {
		product.Name = *newName
		fmt.Print("New name is", product.Name)
	}

	if newQuantity != nil {
		product.Quantity = *newQuantity
	}

	if newPrice != nil {
		product.Price = *newPrice
	}

	inv.products[product.ID] = product
	return nil
}

func (inv *items) DeleteProduct(id int) (string, error) {
	if _, exists := inv.products[id]; exists {
		fmt.Printf("Are you sure you want to delete Product #%d from the items? (Type 'yes' to confirm): ", id)
		var confirmation string
		fmt.Scanln(&confirmation)
		if strings.ToLower(confirmation) == "yes" {
			delete(inv.products, id)
			return fmt.Sprintf("Product #%d removed successfully", id), nil
		} else {
			return fmt.Sprintf("Deletion of Product #%d canceled", id), nil
		}
	}
	return "", fmt.Errorf("Product ID #%d not found", id)
}

func (inv *items) ListProducts() ([]Product, error) {
	if len(inv.products) == 0 {
		return nil, fmt.Errorf("No products found in the items")
	}

	var productList []Product
	for _, p := range inv.products {
		productList = append(productList, p)
	}
	return productList, nil
}

func main() {
	inv := newInventory()
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Welcome to our Inventory Control!")
	for {
		time.Sleep(time.Second * 1)
		fmt.Println("What would you like to do?")
		fmt.Println("1. Add Product")
		fmt.Println("2. Update Product")
		fmt.Println("3. Delete Product")
		fmt.Println("4. List Products")
		fmt.Println("5. Exit")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		choice, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("Invalid input.\nPlease select a number from 1 to 5")
			continue
		}
		switch choice {
		case 1:
			var id int
			var name string
			var quantity int
			var price float64
			fmt.Println("Enter Product ID:")
			input, _ = reader.ReadString('\n')
			input = strings.TrimSpace(input)
			id, err = strconv.Atoi(input)
			if err != nil {
				fmt.Println("Invalid input for ID")
				continue
			}
			fmt.Println("Enter Product Name:")
			name, _ = reader.ReadString('\n')
			name = strings.TrimSpace(name)

			fmt.Println("Enter Product Quantity:")
			input, _ = reader.ReadString('\n')
			input = strings.TrimSpace(input)
			quantity, err = strconv.Atoi(input)
			if err != nil {
				fmt.Println("Invalid input for Quantity")
				continue
			}
			fmt.Println("Enter Product Price in US $:")
			input, _ = reader.ReadString('\n')
			input = strings.TrimSpace(input)
			price, err = strconv.ParseFloat(input, 64)
			if err != nil {
				fmt.Println("Invalid input for Price")
				continue
			}
			p := Product{ID: id, Name: name, Quantity: quantity, Price: price}
			err = inv.AddProduct(p)
			if err != nil {
				fmt.Println("Error adding Product", err)
			} else {
				fmt.Println("Product added successfully.")
			}
		case 2:
			var id int
			fmt.Println("Enter Product ID to update it")
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)
			id, err = strconv.Atoi(input)
			if err != nil {
				fmt.Println("Invalid input for ID. Please select a number.")
				continue
			}

			_, err = inv.GetProduct(id)
			if err != nil {
				fmt.Println("Error getting the Product:", err)
				continue
			} else {
				var newId *int
				var newName *string
				var newQuantity *int
				var newPrice *float64

				fmt.Println("Enter Product NEW ID: (leave it empty to keep current values)")
				input, _ = reader.ReadString('\n')
				input = strings.TrimSpace(input)
				if input != "" {
					newIdValue, err := strconv.Atoi(input)
					if err != nil {
						fmt.Println("Invalid input for ID")
						continue
					}
					newId = &newIdValue
				}

				fmt.Println("Enter Product NEW Name: (leave it empty to keep current values)")
				input, _ = reader.ReadString('\n')
				input = strings.TrimSpace(input)
				if input != "" {
					newName = &input
				}

				fmt.Println("Enter Product NEW Quantity: (leave it empty to keep current values)")
				input, _ = reader.ReadString('\n')
				input = strings.TrimSpace(input)
				if input != "" {
					newQuantityValue, err := strconv.Atoi(input)
					if err != nil {
						fmt.Println("Invalid input for Quantity")
						continue
					}
					newQuantity = &newQuantityValue
				}

				fmt.Println("Enter Product NEW Price in US $: (leave it empty to keep current values)")
				input, _ = reader.ReadString('\n')
				input = strings.TrimSpace(input)
				if input != "" {
					newPriceValue, err := strconv.ParseFloat(input, 64)
					if err != nil {
						fmt.Println("Invalid input for Price")
						continue
					}
					newPrice = &newPriceValue
				}

				err = inv.UpdateProduct(id, newId, newName, newQuantity, newPrice)
				if err != nil {
					fmt.Println("Error updating product:", err)
				} else {
					fmt.Println("Product updated successfully.")
				}
			}
		case 3:
			var id int
			fmt.Println("Enter product ID to REMOVE it\n(WARNING: THIS ACTION CANNOT BE UNDONE!)")
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)
			id, err = strconv.Atoi(input)
			if err != nil {
				fmt.Println("Invalid input for ID")
				continue
			}

			msg, err := inv.DeleteProduct(id)
			if err != nil {
				fmt.Println("Error deleting product:", err)
			} else {
				fmt.Println(msg)
			}
		case 4:
			products, err := inv.ListProducts()
			if err != nil {
				fmt.Println("Error listing products:", err)
			} else {
				for _, product := range products {
					fmt.Printf("ID: %d, Name: %s, Quantity: %d, Price: %.2f\n", product.ID, product.Name, product.Quantity, product.Price)
				}
			}
		case 5:
			fmt.Println("Exiting Inventory Control...")
			time.Sleep(time.Second * 1)
			os.Exit(0)
		default:
			fmt.Println("Invalid option. Please select a number from 1 to 5.")
		}
	}
}
