package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"vds/config"
	"vds/model"
	"vds/server"
	"vds/service"
	"vds/utils"
)

func init() {
	utils.LoadEnv()

	username := utils.GetEnv("DB_USERNAME", "")
	password := utils.GetEnv("DB_PASSWORD", "")
	host := utils.GetEnv("DB_HOST", "")
	port := utils.GetEnv("DB_PORT", "")
	schema := utils.GetEnv("DB_SCHEMA", "")

	err := config.ConnectDB(username, password, host, port, schema)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Database connected successfully")
}

func getUniqueModelsFromCSV(filePath string) (map[string]model.UniqueModel, error) {
	csvFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer csvFile.Close()

	csvReader := csv.NewReader(csvFile)

	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	uniqueModels := make(map[string]model.UniqueModel, len(records)-1)
	for _, record := range records[1:] {
		groupid := record[0]
		modelid := record[1]
		groupname := record[2]
		modelname := record[3]

		key := fmt.Sprintf("%s_%s", groupid, modelid)

		uniqueModels[key] = model.UniqueModel{
			GroupId:   groupid,
			GroupName: groupname,
			ModelId:   modelid,
			ModelName: modelname,
		}
	}
	return uniqueModels, nil
}

func insertUniqueModelsToDB(filePath string) error {
	models, err := getUniqueModelsFromCSV(filePath)
	if err != nil {
		log.Fatalf("Error reading csv file: %v", err)
	}

	uniqueModels := make([]model.UniqueModel, 0, len(models))
	for _, model := range models {
		uniqueModels = append(uniqueModels, model)
	}

	err = service.CreateUniqueModels(uniqueModels)
	if err != nil {
		log.Println("Error inserting unique models to database", err)
		return err
	}
	fmt.Println("Unique models inserted to database successfully")
	return nil
}

func getModelCountFromCSV(filePath string) (map[string]model.UniqueModelCount, error) {
	csvFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer csvFile.Close()

	csvReader := csv.NewReader(csvFile)

	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	uniqueModels := make(map[string]model.UniqueModelCount)
	for _, record := range records[1:] {
		groupid := record[0]
		modelid := record[1]
		deviceid := record[4]

		key := fmt.Sprintf("%s_%s", groupid, modelid)

		if _, exists := uniqueModels[key]; !exists {
			uniqueModels[key] = model.UniqueModelCount{
				GroupID: groupid,
				ModelId: modelid,
				Count:   0,
			}
		}

		if deviceid != "" {
			modelCount := uniqueModels[key]
			modelCount.Count++
			uniqueModels[key] = modelCount
		}
	}

	return uniqueModels, nil
}

func insertModelCountToDB(filePath string) error {
	models, err := getModelCountFromCSV(filePath)
	if err != nil {
		log.Fatalf("Error reading csv file: %v", err)
	}

	uniqueModels := make([]model.UniqueModelCount, 0, len(models))
	for _, model := range models {
		uniqueModels = append(uniqueModels, model)
	}

	err = service.CreateUniqueModelsCount(uniqueModels)
	if err != nil {
		log.Println("Error inserting model count to database", err)
		return err
	}
	fmt.Println("Model count inserted to database successfully")
	return nil
}

func main() {
	var err error

	err = server.StartServer()
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

	if len(os.Args) < 2 {
		log.Fatalf("Error: No file path provided")
	}
	filePath := os.Args[1]

	err = insertModelCountToDB(filePath)
	if err != nil {
		log.Fatalf("Error inserting models count to database: %v", err)
	}

	err = insertUniqueModelsToDB(filePath)
	if err != nil {
		log.Fatalf("Error inserting data to database: %v", err)
	}
}
