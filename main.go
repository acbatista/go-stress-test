package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"sort"
	"sync"
	"time"
)

type Result struct {
	StatusCode int
	Duration   time.Duration
	Error      error
}

type Report struct {
	TotalTime     time.Duration
	TotalRequests int
	StatusCodes   map[int]int
	Errors        int
}

// Mapa de descrições dos códigos HTTP mais comuns
var httpStatusDescriptions = map[int]string{
	200: "OK - Requisição bem sucedida",
	201: "Created - Recurso criado com sucesso",
	204: "No Content - Requisição processada, sem conteúdo para retornar",
	301: "Moved Permanently - Recurso movido permanentemente",
	302: "Found - Redirecionamento temporário",
	400: "Bad Request - Requisição mal formatada",
	401: "Unauthorized - Não autorizado",
	403: "Forbidden - Acesso proibido",
	404: "Not Found - Recurso não encontrado",
	408: "Request Timeout - Tempo limite da requisição excedido",
	429: "Too Many Requests - Muitas requisições",
	500: "Internal Server Error - Erro interno do servidor",
	502: "Bad Gateway - Gateway inválido",
	503: "Service Unavailable - Serviço indisponível",
	504: "Gateway Timeout - Tempo limite do gateway excedido",
}

// Função para categorizar códigos HTTP
func getStatusCategory(code int) string {
	switch {
	case code >= 100 && code < 200:
		return "Informacional (1xx)"
	case code >= 200 && code < 300:
		return "Sucesso (2xx)"
	case code >= 300 && code < 400:
		return "Redirecionamento (3xx)"
	case code >= 400 && code < 500:
		return "Erro do Cliente (4xx)"
	case code >= 500 && code < 600:
		return "Erro do Servidor (5xx)"
	default:
		return "Código Desconhecido"
	}
}

func makeRequest(url string) Result {
	start := time.Now()
	resp, err := http.Get(url)
	duration := time.Since(start)

	if err != nil {
		return Result{
			StatusCode: 0,
			Duration:   duration,
			Error:      err,
		}
	}
	defer resp.Body.Close()

	return Result{
		StatusCode: resp.StatusCode,
		Duration:   duration,
		Error:      nil,
	}
}

func runStressTest(url string, totalRequests int, concurrency int) Report {
	results := make(chan Result, totalRequests)
	var wg sync.WaitGroup
	report := Report{
		StatusCodes: make(map[int]int),
	}

	start := time.Now()

	// Create worker pool
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < totalRequests/concurrency; j++ {
				results <- makeRequest(url)
			}
		}()
	}

	// Handle remaining requests if total is not evenly divisible by concurrency
	remaining := totalRequests % concurrency
	if remaining > 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < remaining; i++ {
				results <- makeRequest(url)
			}
		}()
	}

	// Close results channel when all workers are done
	go func() {
		wg.Wait()
		close(results)
	}()

	// Process results
	for result := range results {
		if result.Error != nil {
			report.Errors++
		} else {
			report.StatusCodes[result.StatusCode]++
		}
	}

	report.TotalTime = time.Since(start)
	report.TotalRequests = totalRequests

	return report
}

func printStatusCodeDistribution(statusCodes map[int]int) {
	// Agrupar códigos por categoria
	categories := make(map[string]map[int]int)
	for code, count := range statusCodes {
		category := getStatusCategory(code)
		if categories[category] == nil {
			categories[category] = make(map[int]int)
		}
		categories[category][code] = count
	}

	// Ordenar categorias
	categoryNames := make([]string, 0, len(categories))
	for category := range categories {
		categoryNames = append(categoryNames, category)
	}
	sort.Strings(categoryNames)

	// Imprimir resultados por categoria
	for _, category := range categoryNames {
		fmt.Printf("\n%s:\n", category)

		// Ordenar códigos dentro da categoria
		codes := make([]int, 0, len(categories[category]))
		for code := range categories[category] {
			codes = append(codes, code)
		}
		sort.Ints(codes)

		for _, code := range codes {
			count := categories[category][code]
			description := httpStatusDescriptions[code]
			if description == "" {
				description = "Status HTTP"
			}
			fmt.Printf("  HTTP %d (%s): %d requisições\n", code, description, count)
		}
	}
}

func main() {
	url := flag.String("url", "", "URL do serviço a ser testado")
	requests := flag.Int("requests", 0, "Número total de requisições a serem feitas")
	concurrency := flag.Int("concurrency", 1, "Número de requisições simultâneas")

	flag.Parse()

	if *url == "" {
		log.Fatal("URL é obrigatória")
	}
	if *requests <= 0 {
		log.Fatal("Número de requisições deve ser maior que 0")
	}
	if *concurrency <= 0 {
		log.Fatal("Concorrência deve ser maior que 0")
	}
	if *concurrency > *requests {
		log.Fatal("Concorrência não pode ser maior que o número total de requisições")
	}

	fmt.Printf("Iniciando teste de carga para %s\n", *url)
	fmt.Printf("Total de requisições: %d\n", *requests)
	fmt.Printf("Nível de concorrência: %d\n\n", *concurrency)

	report := runStressTest(*url, *requests, *concurrency)

	fmt.Println("\nResultados do Teste:")
	fmt.Printf("Tempo total: %v\n", report.TotalTime)
	fmt.Printf("Total de requisições: %d\n", report.TotalRequests)

	fmt.Println("\nDistribuição de Códigos de Status:")
	printStatusCodeDistribution(report.StatusCodes)

	if report.Errors > 0 {
		fmt.Printf("\nRequisições com falha: %d\n", report.Errors)
	}
}
