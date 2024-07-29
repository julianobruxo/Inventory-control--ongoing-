package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type product struct { // tipos que serão usados
	ID       int
	Name     string
	Quantity int
	Price    float64
}

type Inventory interface { //interface que lista os metodos que serão utilizados
	AddProduct(p product) error
	GetProduct(id int) (product, error)
	UpdateProduct(p product) error
	DeleteProduct(p product) error
	ListProducts() ([]product, error)
}

type inventory struct { // mapa que armazena os productos
	products map[int]product //utiliza um int da ID como valor e o product como valor
}

// Métodos
// // mètodo que lista os produtos no inventário e retorna um ponteiro do inventário
func newInventory() *inventory {
	return &inventory{
		products: make(map[int]product),
	}
}

func (inv *inventory) AddProduct(p product) error {
	if _, exists := inv.products[p.ID]; exists {
		return fmt.Errorf("Product ID #%d already exists", p.ID) //Erro caso o produto com aquela ID já exista em estoque
	}
	inv.products[p.ID] = p // adiciona o produto p ao mapa inv.products com a chave p.ID
	return nil
}

// Buscar Produto: método aponta para inventory, recebe um int da ID e retorn um produto e um erro
func (inv *inventory) GetProduct(id int) (product, error) {
	if p, exists := inv.products[id]; exists {
		//se o produto existir, retorna o produto
		return p, nil
	}
	return product{}, fmt.Errorf("Product with ID#:%d not found in our inventory", id)
	// do contrário, retorna um erro
}

// Atualiza os produtos no inventário
func (inv *inventory) UpdateProduct(p product) error {
	if _, exists := inv.products[p.ID]; exists {
		inv.products[p.ID] = p //se o produto existir, retorna o produto
		return nil
	}
	return fmt.Errorf("Product with ID #:%d not found in our inventory", p.ID)
	// do contrário, retorna um erro
}

func (inv *inventory) DeleteProduct(id int) (string, error) {
	if _, exists := inv.products[id]; exists {
		fmt.Printf("Are you sure you want to delete product #%d from the inventory?", id)
		var confirmation string
		fmt.Scanln(&confirmation)
		if confirmation == "Sim" || confirmation == "sim" {
			delete(inv.products, id) //se o produto existir,deleta o produto
			return fmt.Sprintf("Product #%d removed successfully", id), nil
		} else {
			return fmt.Sprintf("Deletion of Product #%d canceled", id), nil // mensagem de erro
		}
	}
	return fmt.Sprintf("Product ID #%d not found", id), nil

}

func (inv *inventory) ListProducts() ([]product, error) {
	if len(inv.products) == 0 {
		return nil, fmt.Errorf("No products found in the inventory")
	}

	var productList []product
	for _, p := range inv.products {
		productList = append(productList, p)
	}
	return productList, nil
}

func main() {
	inv := newInventory()
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("Welcome to our Inventory Control!\nWhat would you like to do?")
		fmt.Println("1. Add Product")
		fmt.Println("2. Update Product")
		fmt.Println("3. Delete Product")
		fmt.Println("4. List Products")
		fmt.Println("5. Exit")
		var choice int
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		choice, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("Invalid input.\nPlease select a number from 1 to 5")

			var discard string
			fmt.Scanln(&discard)
			continue
		}
		switch choice {
		case 1:
			var id int
			var name string
			var quantity int
			var price float64
			fmt.Println("Enter product ID:")
			input, _ = reader.ReadString('\n')
			input = strings.TrimSpace(input)
			id, err = strconv.Atoi(input)
			if err != nil {
				fmt.Println("Invalid input for ID")
				continue
			}
			fmt.Println("Enter product Name:")
			name, _ = reader.ReadString('\n')
			name = strings.TrimSpace(name)

			fmt.Println("Enter product Quantity:")
			input, _ = reader.ReadString('\n')
			input = strings.TrimSpace(input)
			quantity, err = strconv.Atoi(input)
			if err != nil {
				fmt.Println("Invalid input for Quantity")
				continue
			}
			fmt.Println("Enter product Price in US $:")
			input, _ = reader.ReadString('\n')
			input = strings.TrimSpace(input)
			price, err = strconv.ParseFloat(input, 64)
			if err != nil {
				fmt.Println("Invalid input for Price")
				continue
			}
			//add product
			p := product{ID: id, Name: name, Quantity: quantity, Price: price}
			err := inv.AddProduct(p)
			if err != nil {
				fmt.Println("Error addidng product", err)
			} else {
				fmt.Println("Product added sucessfully.")
			}
		case 2:
			//update product
			var id int
			fmt.Println("Enter product ID to update it")
			fmt.Scanln(&id)

			p, err := inv.GetProduct(id)
			if err != nil {
				fmt.Println("Error getting the product")
			} else {
				// solicita novos dados para atualizar o produto
				var id int
				var name string
				var quantity int
				var price float64

				fmt.Println("Enter product NEW ID: (leave it empty to keep current values)")
				fmt.Scanln(&id)
				fmt.Println("Enter product NEW Name:(leave it empty to keep current values)")
				fmt.Scanln(&name)
				fmt.Println("Enter product NEW Quantity:(leave it empty to keep current values)")
				fmt.Scanln(&quantity)
				fmt.Println("Enter product NEW Price in US $:(leave it empty to keep current values)")
				fmt.Scanln(&price)

				if name != "" {
					p.Name = name
				}
				if quantity != 0 {
					p.Quantity = quantity
				}
				if id != 0 {
					p.ID = id
				}
				if price != 0 {
					p.Price = price
				}

				//TRATAMENTO DE ERRO
				err := inv.UpdateProduct(p)
				if err != nil {
					fmt.Println("Error updating product", err)
				} else {
					fmt.Println("Product updated successfully")
				}
			}

		case 3:
			//delete product
			var id int
			fmt.Println("Enter product ID to REMOVE it\n(WARNING:\ntHIS ACTION CANNOT BE UNDONE!)")
			fmt.Scanln(&id)

			msg, err := inv.DeleteProduct(id)
			if err != nil {
				fmt.Println("Error removing product. Try again", err)
			} else {
				fmt.Println(msg)
			}

		case 4:
			//list products
			products, err := inv.ListProducts()
			if err != nil {
				fmt.Println("Error fetching the products", err)
			} else {
				fmt.Println("List of available produts in inventory:")
				for _, p := range products {
					fmt.Printf("ID: %d, Name: %s, Quantity: %d, Price: %.2f\n", p.ID, p.Name, p.Quantity, p.Price)
				}
			}

		case 5:
			fmt.Println("Thank you for using our  Inventory Control")
			return
		default:
			fmt.Println("Invalid choice. Please try again.")
		}

	}

}
