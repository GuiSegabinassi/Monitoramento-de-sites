package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const monitoramento = 1

const delay = 10

func main() {
	exibeIntrodução()

	for {
		exibeMenu()

		comando := lerComando()

		if comando == 1 {
			iniciarMonitoramento()
		} else if comando == 2 {
			fmt.Println("Exibindo Logs...")
			imprimeLogs()

		} else if comando == 3 {
			fmt.Println("Finalizando a execução do programa...")
			os.Exit(3)
		} else {
			fmt.Println("O comando digitado não é um comando reconhecido no menu")
			fmt.Println("Finalizando a execução do programa...")
			os.Exit(-1)
		}
	}
}

func exibeIntrodução() {
	var nome string
	var versao float32 = 1.1

	fmt.Println("Digite o seu nome: ")
	fmt.Scan(&nome)
	fmt.Println("Olá senhor,", nome)
	fmt.Println("O senhor está na versão :", versao)
	fmt.Println("...")
}

func exibeMenu() {
	fmt.Println("Escolha uma das opções abaixo: ")

	fmt.Println("1 - Iniciar Monitoramento")
	fmt.Println("2 - Exibir Logs")
	fmt.Println("3 - Sair do Programa")
}

func lerComando() int {
	var comandoLido int
	fmt.Scan(&comandoLido)
	fmt.Println("O comando escolhido foi: ", comandoLido)

	return comandoLido
}

func iniciarMonitoramento() {
	fmt.Println("Monitorando...")

	// sites := []string{
	// 	"https://httpbin.org/status/",
	// 	"https://stackoverflow.com",
	// 	"https://www.runningland.com.br/",
	// 	"https://www.riotgames.com/pt-br",
	// 	"https://www.premierleague.com/",
	// 	"https://olympics.com/pt/paris-2024/os-jogos/jogos-olimpicos-paralimpicos/jogos-paralimpicos",
	// 	"https://www.nasdaq.com",
	// 	"https://www.b3.com.br/pt_br/para-voce",
	// 	"https://www.amazon.com.br/",
	// 	"https://www.tesla.com/",
	// 	"https://www.leagueoflegends.com/pt-br/",
	// }

	sites := lerSitesDoArquivo()

	for i := 0; i < monitoramento; i++ {
		for i, site := range sites {
			fmt.Println("Testando site :", i, site)
			testaSite(site)

			fmt.Println("****************************************************************")
		}
		fmt.Println("********************************************************************")
		fmt.Println("Testando outra vez")
		fmt.Println("********************************************************************")

		time.Sleep(delay * time.Second)
	}

}

func testaSite(site string) {
	resp, err := http.Get(site)
	if err != nil {
		fmt.Println("Erro ao tentar acessar o site:", site, "-", err)
		return
	}
	defer resp.Body.Close()

	if resp != nil && resp.StatusCode == 200 {
		fmt.Println("Site:", site, "foi carregado com sucesso!", resp.StatusCode)
		registraLog(site, true)
	} else if resp != nil {
		fmt.Println("Site:", site, "está com problemas!", resp.StatusCode)
		registraLog(site, false)
	} else {
		fmt.Println("Site:", site, "não pôde ser carregado.")
	}
}

func lerSitesDoArquivo() []string {

	var sites []string
	arquivo, err := os.Open("sites.txt")

	if err != nil {
		fmt.Println("Ocorreu um erro!", err)
	}

	leitor := bufio.NewReader(arquivo)

	for {
		linha, err := leitor.ReadString('\n')
		linha = strings.TrimSpace(linha)

		sites = append(sites, linha)

		fmt.Println(linha)

		if err == io.EOF {
			break
		}

	}

	arquivo.Close()

	return sites
}

func registraLog(site string, status bool) {
	arquivo, err := os.OpenFile("Log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		fmt.Println("Ocorreu um erro: ", err)
	}

	arquivo.WriteString(time.Now().Format("02/01/2006 15:04:05") + " - " + site + "- Online: " + strconv.FormatBool(status) + "\n")

	arquivo.Close()
}

func imprimeLogs() {
	arquivo, err := os.ReadFile("Log.txt")

	if err != nil {
		fmt.Println("Ocorreu um erro: ", err)
	}

	fmt.Println(string(arquivo))

}
