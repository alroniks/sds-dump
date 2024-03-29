package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

const PRICES = "resources/prices.csv"
const IMPORT = "resources/import.csv"
const OUTPUT = "resources/output.json"
const VIDEOS = "resources/videos.json"

type ProductValue struct {

	ID int `json:"id"`
	Article string `json:"article"`
	Title string `json:"title"`
	Image string `json:"image"`
	Brand string `json:"brand"`
	Price string `json:"price"`
	Units string `json:"units"`
	InPack string `json:"inpack"`
	Description string `json:"description"`
	Availability int `json:"availability"`
	Category int `json:"category"`
	Video string `json:"video"`
	Link string `json:"link"`

}

type Video struct {
	Product int `json:"id"`
	Video string `json:"video"`
}

func main() {

	pricesFile, err := os.Open(PRICES)
	defer pricesFile.Close()
	if err != nil {
		panic(err)
	}
	lines, err := csv.NewReader(pricesFile).ReadAll();
	if err != nil {
		panic(err)
	}

	prices := make(map[string]string)

	for _, line := range lines {
		prices[line[0]] = line[1]
	}

	videosFile, err := os.Open(VIDEOS)
	defer videosFile.Close()
	if err != nil {
		panic(err)
	}
	data, _ := ioutil.ReadAll(videosFile)
	var videos []Video
	json.Unmarshal(data, &videos)

	videosMap := make(map[int]string)

	for _, line := range videos {
		videosMap[line.Product] = line.Video
	}

	// updating search file
	outputFile, err := os.Open(OUTPUT)
	defer outputFile.Close()
	if err != nil {
		panic(err)
	}
	output, _ := ioutil.ReadAll(outputFile)
	var products []ProductValue
	json.Unmarshal(output, &products)

	importFile, _ := os.Create(IMPORT)
	defer importFile.Close()
	writer := csv.NewWriter(importFile)
	defer writer.Flush()

	for i, product := range products {
		var actualPrice string
		if val, ok := prices[product.Article]; ok {
			products[i].Price = val
			actualPrice = val
		} else {
			products[i].Price = ""
			actualPrice = ""
		}

		var actualVideo string
		if val, ok := videosMap[product.ID]; ok {
			products[i].Video = val
			actualVideo = val
		} else {
			products[i].Video = ""
			actualVideo = ""
		}

		csvError := writer.Write([]string {
			strconv.Itoa(product.ID),
			strconv.Itoa(product.Category),
			product.Article,
			product.Brand,
			actualPrice,
			product.Title,
			product.Units,
			product.InPack,
			strconv.Itoa(product.Availability),
			actualVideo,
		})

		if csvError != nil {
			fmt.Println("Error of writing record to csv: ", csvError)
		}
	}

	jsonProducts, _ := json.MarshalIndent(products, "", "  ")

	ioutil.WriteFile(OUTPUT, jsonProducts, 0644)
}
