package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

func main() {
	// Configuração do servidor HTTP com endpoints para o front-end
	http.HandleFunc("/start-monitoring", startMonitoringHandler)
	http.HandleFunc("/logs", latestLogHandler)
	http.HandleFunc("/list-logs", listLogsHandler)
	http.HandleFunc("/view-log", viewLogHandler)
	http.HandleFunc("/exit", exitHandler)
	http.Handle("/", http.FileServer(http.Dir("../web"))) // Serve a pasta de front-end

	fmt.Println("Server is running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Error starting server:", err)
	}
}

func iniciarMonitoramento() {
	fmt.Println("Starting monitoring...")

	sites := lerSitesDoArquivo()
	if len(sites) == 0 {
		fmt.Println("No sites found to monitor.")
		return
	}

	for _, site := range sites {
		testaSite(site)
		fmt.Println("-------------------------------------------------------------")
	}

	fmt.Println("Monitoring completed.")
}

func startMonitoringHandler(w http.ResponseWriter, r *http.Request) {
	iniciarMonitoramento()
	fmt.Fprintln(w, "Monitoring completed!")
}

func latestLogHandler(w http.ResponseWriter, r *http.Request) {
	logFile, err := getLatestLogFile()
	if err != nil {
		http.Error(w, "Could not read logs", http.StatusInternalServerError)
		return
	}

	data, err := os.ReadFile(logFile)
	if err != nil {
		http.Error(w, "Could not read log file", http.StatusInternalServerError)
		return
	}

	w.Write(data)
}

func getLatestLogFile() (string, error) {
	files, err := os.ReadDir(".")
	if err != nil {
		return "", err
	}

	var logFiles []string
	for _, file := range files {
		if strings.HasPrefix(file.Name(), "Log_") && strings.HasSuffix(file.Name(), ".txt") {
			logFiles = append(logFiles, file.Name())
		}
	}

	if len(logFiles) == 0 {
		return "", fmt.Errorf("no log files found")
	}

	sort.Strings(logFiles)
	return logFiles[len(logFiles)-1], nil
}

func listLogsHandler(w http.ResponseWriter, r *http.Request) {
	files, err := os.ReadDir(".") // Lê o diretório atual
	if err != nil {
		http.Error(w, "Could not read log directory", http.StatusInternalServerError)
		return
	}

	var logFiles []string
	for _, file := range files {
		if strings.HasPrefix(file.Name(), "Log_") && strings.HasSuffix(file.Name(), ".txt") {
			logFiles = append(logFiles, file.Name())
		}
	}

	for _, logFile := range logFiles {
		fmt.Fprintln(w, logFile)
	}
}

func viewLogHandler(w http.ResponseWriter, r *http.Request) {
	logFile := r.URL.Query().Get("file")
	if logFile == "" || !strings.HasPrefix(logFile, "Log_") || !strings.HasSuffix(logFile, ".txt") {
		http.Error(w, "Invalid log file", http.StatusBadRequest)
		return
	}

	data, err := os.ReadFile(logFile)
	if err != nil {
		http.Error(w, "Could not read log file", http.StatusInternalServerError)
		return
	}

	w.Write(data)
}

func exitHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Exiting program...")
	os.Exit(0)
}

func testaSite(site string) {
	logFileName := fmt.Sprintf("Log_%s.txt", time.Now().Format("2006-01-02_15-04-05"))
	arquivo, err := os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("Error opening log file: %v", err)
		return
	}
	defer arquivo.Close()

	resp, err := http.Get(site)
	if err != nil {
		fmt.Printf("Error accessing site %s: %v\n", site, err)
		registraLog(arquivo, site, false)
		return
	}
	defer resp.Body.Close()

	status := resp.StatusCode == 200
	fmt.Printf("Site: %s - Status: %s\n", site, statusText(status))
	registraLog(arquivo, site, status)
}

func registraLog(arquivo *os.File, site string, status bool) {
	linha := fmt.Sprintf("%s - %s - Online: %t\n", time.Now().Format("02/01/2006 15:04:05"), site, status)
	if _, err := arquivo.WriteString(linha); err != nil {
		log.Printf("Error writing to log file: %v", err)
	}
}

func lerSitesDoArquivo() []string {
	var sites []string
	arquivo, err := os.Open("sites.txt")
	if err != nil {
		log.Printf("Error opening sites file: %v", err)
		return sites
	}
	defer arquivo.Close()

	leitor := bufio.NewScanner(arquivo)
	for leitor.Scan() {
		linha := strings.TrimSpace(leitor.Text())
		if linha != "" {
			sites = append(sites, linha)
		}
	}

	if err := leitor.Err(); err != nil {
		log.Printf("Error reading sites file: %v", err)
	}
	return sites
}

// Função auxiliar para exibir o texto do status
func statusText(status bool) string {
	if status {
		return "Online"
	}
	return "Offline"
}
