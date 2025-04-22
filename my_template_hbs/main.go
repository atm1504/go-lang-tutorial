package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/aymerick/raymond"
)

type Traits struct {
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	IsEmailRefuse bool   `json:"is_email_refuse"`
}

type Banner struct {
	Description string `json:"description"`
	ImageURL    string `json:"imageUrl"`
	LinkURL     string `json:"linkUrl"`
}

type Product struct {
	Brand              string  `json:"brand"`
	HasMultiSalesPrice int     `json:"hasMultiSalesPrice"`
	HasSpecialPrice    int     `json:"hasSpecialPrice"`
	ProductCode        int     `json:"productCode"`
	ImageURL           string  `json:"imageUrl"`
	ItemAddable        *bool   `json:"itemAddable"`
	LinkURL            string  `json:"linkUrl"`
	Name               string  `json:"name"`
	NumReviews         int     `json:"numReviews"`
	Price              int     `json:"price"`
	Quantity           *string `json:"quantity"`
	Rating             float64 `json:"rating"`
	TaxPrice           int     `json:"taxPrice"`
	TaxRate            float64 `json:"taxRate"`
	TaxRateType        int     `json:"taxRateType"`
}

type RecommendButtons struct {
	Text     string `json:"text"`
	ImageURL string `json:"imageUrl"`
	LinkURL  string `json:"linkUrl"`
}

type TriggerPlacement struct {
	BannerLogicID     string             `json:"bannerLogicId"`
	Message           string             `json:"message"`
	Placement         *string            `json:"placement"`
	ProductLogicID    string             `json:"productLogicId"`
	RecommendBanners  []Banner           `json:"recommendBanners"`
	RecommendProducts []Product          `json:"recommendProducts"`
	RecommendButtons  []RecommendButtons `json:"recommendButtons"`
}

type Text struct {
	SelectionTitle       string `json:"selection_title"`
	SelectionDescription string `json:"selection_description"`
}

type Events struct {
	TriggerPlacement TriggerPlacement `json:"triggerPlacement"`
	Text             Text             `json:"text"`
}

type Recsys struct {
	DMSelection TriggerPlacement `json:"dm-selection"`
	Text        Text             `json:"text"`
}

type Data struct {
	UserID string `json:"userId"`
	Traits Traits `json:"traits"`
	Events Events `json:"events"`
	Recsys Recsys `json:"recsys"`
}

func renderTemplate(data Data, source string, results chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	result, err := raymond.Render(source, data)
	if err != nil {
		log.Fatalf("Error rendering template: %v", err)
	}
	results <- result
}

func worker(id int, jobs <-chan int, results chan<- string, wg *sync.WaitGroup, source string, data Data) {
	for range jobs {
		renderTemplate(data, source, results, wg)
	}
}

func main() {
	source, err := ioutil.ReadFile("./full_html.hbs")
	if err != nil {
		log.Fatalf("Error reading template file: %v", err)
	}

	jsonData, err := ioutil.ReadFile("./data.json")
	if err != nil {
		log.Fatalf("Error reading JSON file: %v", err)
	}

	var data Data
	if err := json.Unmarshal(jsonData, &data); err != nil {
		log.Fatalf("Error parsing JSON data: %v", err)
	}

	// Register helpers
	raymond.RegisterHelper("isMultipleOf", func(index int, multiple int) bool {
		return index%multiple == 0
	})

	raymond.RegisterHelper("isLastRowStart", func(index int, length int) bool {
		return index == length-1
	})

	raymond.RegisterHelper("isLastInRow", func(index int, length int) bool {
		return (index+1)%3 == 0
	})

	raymond.RegisterHelper("isLastProduct", func(index int, length int) bool {
		return index == length-1
	})

	raymond.RegisterHelper("eq", func(a interface{}, b interface{}) bool {
		return a == b
	})

	raymond.RegisterHelper("getImageUrl", func(context []string, index int) string {
		return context[index]
	})

	t1 := time.Now()
	fmt.Println(t1)

	numJobs := 100000
	numWorkers := runtime.NumCPU()

	jobs := make(chan int, numJobs)
	results := make(chan string, numJobs)

	var wg sync.WaitGroup

	for w := 1; w <= numWorkers; w++ {
		go worker(w, jobs, results, &wg, string(source), data)
	}

	for j := 1; j <= numJobs; j++ {
		wg.Add(1)
		jobs <- j
	}

	close(jobs)
	wg.Wait()
	close(results)

	t2 := time.Now()
	fmt.Println(t2)
	fmt.Println(t2.Sub(t1).Milliseconds())
}
