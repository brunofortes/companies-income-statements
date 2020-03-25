# Demonstração de Resultados de empresas listadas na B3 / *B3 companies income statements*

Uma fonte para consulta de dados históricos das demonstrações de resultados de empresas listadas em bolsas de valores. Os dados iniciais são obtidos a partir de *scraper* na API do Yahoo! Finance. São coletados os indicadores considerados mais importantes para analises fundamentalistas das empresas. Embora desenvolvido inicialmente para empresas listadas B3, é possível coletar dados de empresas listadas em outras bolsas.

Os dados coletados são armazenados em um banco MongoDB e serão providos através de Rest API.

Uma empresa contém a seguinte estrutura:

```
type Company struct {
	Name       string     
	Labels     []string    // or tickets, ex: PETR4, VALE3, AAPL, AMZN...
	Sectors    []string   
	Industries []string   
	Financials []Financial
}
```

Onde *Financials* é uma lista de *Financial* com a seguinte estrutura:

```
type Financial struct {
	Date                     time.Time
	AvarageSharesOutstanding float64   // Quantidade de ações em circulação
	StockholdersEquity       float64   // Patrimônio Líquido
	TotalRevenue             float64   // Receita Líquida
	NetIncome                float64   // Lucro Líquido
	EBIT                     float64  
}
```

## Pré requisitos

* [Go v1.13 or higher](https://golang.org/)
* [MongoDB](https://www.mongodb.com/)


## 1ª Etapa

Scraper do Yahoo! Finance para obter os dados históricos das demonstrações de resultados. Após clonar o repositório e garantir que Go e MongoDB estejam rodando, basta rodar o seguinte comando a partir da raiz do repositório:

```
go run scrap.go [quarterly|annual] periodBegin periodEnd
```

onde ```[quarterly|annual]``` determina se os dados devem ser trimestrais ou acumulado do ano

sendo  que ```periodBegin``` e ```periodEnd``` tem o seguinte formato ```yyyy-MM-dd```

Exemplo:
```
go run scrap.go quarterly 2017-12-20 2020-06-30
```

O ponto de partida para scraper são as empresas listadas na seguinte url: https://br.financas.yahoo.com/noticias/acoes-mais-negociadas , porém isso pode ser alterado no arquivo scraper.go, onde logo no início do arquivo é definida a URL para consulta:

```
const (
	URL = "https://br.financas.yahoo.com/noticias/acoes-mais-negociadas" // B3
	//URL = "https://finance.yahoo.com/sector/ms_technology" // Technology
	//URL = https://finance.yahoo.com/sector/ms_energy // Energy
)
```

## 2ª Etapa

  Rest API para acesso as informações - Em desenvolvimento...


## Obs

* Este projeto não tem como objetivo fornecer o preço das ações ao longo do tempo, existem diversas outras ferramentas abertas que podem ser utilizadas para esse fim;

* O preço das ações combinados com os dados aqui coletados são capazes de fornecer indicadores como: LPA(Lucro por Ação), VPA(Valor Patrimonial por Ação), Margem EBIT, Margem Liquida, P/L (preço dividido por LPA), P/VP (preço dividido por VPA), Valor de Mercado, Valor da Firma...


## Licença

[![License](http://img.shields.io/:license-mit-blue.svg?style=flat-square)](http://badges.mit-license.org)

## Referências
* [Learn Go With Tests](https://github.com/larien/learn-go-with-tests) - Ponto de partida do meu aprendizado na linguagem Go;
* [How to use Go with MongoDB using official mongodb-go-driver!](https://medium.com/dev-howto/how-to-use-go-with-mongodb-using-official-mongodb-go-driver-76c1628dae1e)
* [Repository, Service Patern Go](https://www.reddit.com/r/golang/comments/9h7dnn/repository_service_patern_go/)
* [Fundamentus](https://www.fundamentus.com.br/) e [InvestSite](https://www.investsite.com.br/) - Sites com indicadores fundamentalistas de empresas listadas na B3, me ajudou a identificar quais são os indicadores mais utilizados.
* [Colly](http://go-colly.org/) - Scraping Framework for Go
