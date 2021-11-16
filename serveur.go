package main

import (
	"bufio"
	"essai/matrixMult"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

var count = 0
var taille = make([]int, 4)

func handleConnection(c net.Conn) {
	fmt.Print(".")
	var taille []int

	for {

		netData, err := bufio.NewReader(c).ReadString('\n')

		if err != nil {
			fmt.Println(err)
			return
		}

		data := strings.TrimSpace(string(netData))
		//i, err := strconv.Atoi(temp)

		for _, s := range data {
			j, _ := strconv.Atoi(string(s))
			if j != 0 {
				taille = append(taille, j)
			}
		}

		var a [][]int
		var b [][]int

		//if len(taille) == 4 {
		//	rowA := taille[0]
		//	colA := taille[1]
		//	rowB := taille[2]
		//	colB := taille[3]
		//a, b=matrixMult.GenerateRandomMatrix(rowA, colA), matrixMult.GenerateRandomMatrix(rowB, colB)

		//}
		a, b = matrixMult.GenerateRandomMatrix(1000, 1000), matrixMult.GenerateRandomMatrix(1000, 1000)

		if data == "STOP" {
			break
		}
		//messages renvoyés par le serveur:

		matrixMult.CountGoroutines(a, b)
		//result := matrixMult.MultiplyMatricesGoRoutine(a, b)
		//pour écrire les élements du produit de matrice dans un fichier
		//matrixMult.WriteMatrixElements(result)
		//écrire le temps d'exécution pour la multiplication sans goroutine puis avec
		matrixMult.WriteExecTime(a, b)

		nbGoroutines := strconv.Itoa(matrixMult.CountGoroutines(a, b)) + "\n"
		counter := strconv.Itoa(count) + "\n"

		c.Write([]byte(string(nbGoroutines))) //envoyer le message au client
		c.Write([]byte(string(counter)))

		taille = taille[:0] //vider le tableau contenant les éléments indiquant la taille des deux matrices

	}

	c.Close()
}

func main() {

	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide a port number!")
		return
	}

	PORT := ":" + arguments[1]
	l, err := net.Listen("tcp4", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go handleConnection(c)
		count++
	}
}

//
