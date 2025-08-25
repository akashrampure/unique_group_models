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

// to convert Model data from csv to map
func csvToMap(filePath string) (map[string]model.UniqueModel, error) {
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

// to insert unique models to database
func insertDataToDB(models map[string]model.UniqueModel) error {
	uniqueModels := make([]model.UniqueModel, 0, len(models))
	for _, model := range models {
		uniqueModels = append(uniqueModels, model)
	}

	err := service.CreateUniqueModels(uniqueModels)
	if err != nil {
		log.Println("Error creating models", err)
		return err
	}
	fmt.Println("Models created successfully")
	return nil
}

// to insert model count to database
func insertModelCountToDB(models map[string]model.UniqueModelCount) error {
	uniqueModels := make([]model.UniqueModelCount, 0, len(models))
	for _, model := range models {
		uniqueModels = append(uniqueModels, model)
	}

	err := service.CreateUniqueModelsCount(uniqueModels)
	if err != nil {
		log.Println("Error inserting models count to database", err)
		return err
	}
	fmt.Println("Models count inserted to database successfully")
	return nil
}

// to output unique models
func outputDataToCSV(models map[string]model.UniqueModel) error {
	file, err := os.Create("output.csv")
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write([]string{"GroupID", "GroupName", "ModelId", "ModelName"}); err != nil {
		return err
	}

	for _, model := range models {
		if err := writer.Write([]string{model.GroupId, model.GroupName, model.ModelId, model.ModelName}); err != nil {
			return err
		}
	}

	return nil
}

// to get model count from csv file
func modelCount(filePath string) (map[string]model.UniqueModelCount, error) {
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

func main() {
	modelsCount, err := modelCount("/home/akash/Documents/getmyvds/data/No_of_devices_under_each_group.csv")
	if err != nil {
		log.Fatalf("Error reading csv file: %v", err)
	}

	err = insertModelCountToDB(modelsCount)
	if err != nil {
		log.Fatalf("Error inserting models count to database: %v", err)
	}

	server.StartServer()

	// models, err := csvToMap("/home/akash/Documents/getmyvds/data/No_of_devices_under_each_group.csv")
	// if err != nil {
	// 	log.Fatalf("Error reading csv file: %v", err)
	// }

	// err = outputDataToCSV(models)
	// if err != nil {
	// 	log.Fatalf("Error writing csv file: %v", err)
	// }

	// err = insertDataToDB(models)
	// if err != nil {
	// 	log.Fatalf("Error inserting data to database: %v", err)
	// }
}
