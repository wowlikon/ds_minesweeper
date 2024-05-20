package main

import (
	"fmt"
	"math/rand"
)

const side = 10
const count = 6

var field [][]int

func main() {
	// Initialize the field
	field = make([][]int, side)
	for i := 0; i < side; i++ {
		field[i] = make([]int, side)
	}

	for i := 0; i < side; i++ {
		for j := 0; j < side; j++ {
			field[i][j] = 0
		}
	}

	//Generate
	for i := 0; i < count; i++ {
		x := rand.Intn(side)
		y := rand.Intn(side)
		field[x][y] = -1
	}

	//Calculate
	for i := 0; i < side; i++ {
		for j := 0; j < side; j++ {
			if field[i][j] != -1 {
				count := 0
				for k := i - 1; k <= i+1; k++ {
					for l := j - 1; l <= j+1; l++ {
						if k >= 0 && k < side && l >= 0 && l < side {
							if field[k][l] == -1 {
								count++
							}
						}
					}
				}
				field[i][j] = count
			}
		}
	}

	//Print
	fmt.Printf("Mines:%d\n", count)

	for i := 0; i < side; i++ {
		for j := 0; j < side; j++ {
			switch field[i][j] {
			case 0:
				fmt.Print("||0️⃣||")
				break
			case 1:
				fmt.Print("||1️⃣||")
				break
			case 2:
				fmt.Print("||2️⃣||")
				break
			case 3:
				fmt.Print("||3️⃣||")
				break
			case 4:
				fmt.Print("||4️⃣||")
				break
			case 5:
				fmt.Print("||5️⃣||")
				break
			case 6:
				fmt.Print("||6️⃣||")
				break
			case 7:
				fmt.Print("||7️⃣||")
				break
			case 8:
				fmt.Print("||8️⃣||")
				break
			case 9:
				fmt.Print("||9️⃣||")
				break
			default:
				fmt.Print("||💣||")
				break
			}
		}
		fmt.Println()
	}
}
