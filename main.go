package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
)

type Lesson struct {
	Section string
	Title   string
	Url     string
}

var wg sync.WaitGroup

func parseLessons(fileName string) []Lesson {
	var lessons []Lesson
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("\nError opening input CSV file:", err)
		return lessons
	}
	defer file.Close()

	// Create a new CSV reader
	reader := csv.NewReader(file)

	// Read all records from CSV
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("\nError reading CSV records:", err)
		return lessons
	}

	// Process each record
	for idx, record := range records {
		if idx == 0 {
			continue
		}
		lessons = append(lessons, Lesson{
			record[0],
			record[1],
			record[2],
		})
	}
	return lessons
}
func getDownloadLink(id string) (string, error) {
	var link string
	lessonData := "https://fast.wistia.com/embed/medias/" + id + ".json"
	response, err := http.Get(lessonData)
	if err != nil {
		fmt.Println("\nError making HTTP request:", err)
		return "", err
	}
	defer response.Body.Close()

	// Read the response body
	jsonData, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("\nError reading response body:", err)
		return "", err
	}
	var data map[string]interface{}
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		fmt.Println("\nError decoding JSON:", err)
		return "", err
	}

	// Extract and print URLs
	assets, ok := data["media"].(map[string]interface{})["assets"].([]interface{})
	if !ok {
		fmt.Println("\nError accessing assets")
		return "", err
	}

	for _, asset := range assets {
		url, ok := asset.(map[string]interface{})["url"].(string)
		if !ok {
			fmt.Println("\nError accessing URL")
			continue
		}
		lessonType, ok := asset.(map[string]interface{})["type"].(string)
		if !ok {
			fmt.Println("\nError accessing URL")
			continue
		}
		if lessonType == "original" {
			link = url
		}
	}
	return link, nil
}

func main() {
	fmt.Print("Parsing input..")
	lessons := parseLessons("input.csv")
	fmt.Println("\rFound", len(lessons), "videos")
	for idx, lesson := range lessons {
		idx = idx + 1
		fmt.Print("\rDownloading ", "(", idx, "/", len(lessons), ").. [", lesson.Title, "]")
		url, err := getDownloadLink(lesson.Url)
		if err != nil {
			fmt.Println("\nFailed to fetch download link for ", lesson.Title)
			continue
		}
		err = downloadFile(url, "downloads/"+lesson.Section+"/"+strings.Replace(lesson.Title, "/", "-", -1)+".mp4")
		if err != nil {
			fmt.Println("\nError downloading file:", err)
			return
		}
	}
	fmt.Println("\nFinished.")
}

func downloadFile(url string, outputPath string) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}
