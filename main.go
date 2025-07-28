package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	mcp_golang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport/stdio"

	dotenv "github.com/joho/godotenv"
)

type ExchangeReqArgs struct {
	From   string
	To     string
	Amount float32
}

type ExchangeRateResponse struct {
	Success bool    `json:"success"`
	Terms   string  `json:"terms"`
	Result  float32 `json:"result"`
}

func getCurrency(args ExchangeReqArgs, apiKey string) (*ExchangeRateResponse, error) {
	apiUrl := fmt.Sprintf("https://api.exchangerate.host/convert?from=%s&to=%s&amount=%.2f&access_key=%s", args.From, args.To, args.Amount, apiKey)

	resp, err := http.Get(apiUrl)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var exchangeRateResponse ExchangeRateResponse

	if err = json.NewDecoder(resp.Body).Decode(&exchangeRateResponse); err != nil {
		return nil, err
	}

	return &exchangeRateResponse, nil
}

func main() {

	dotenv.Load()

	var EXCHANGE_RATE_API_KEY string = os.Getenv("EXCHANGE_RATE_API_KEY")

	done := make(chan struct{})

	mcpServer := mcp_golang.NewServer(stdio.NewStdioServerTransport())

	err := mcpServer.RegisterTool("exchange-mcp", "Get Exchange Currency", func(arguments ExchangeReqArgs) (*mcp_golang.ToolResponse, error) {
		exchangeResp, err := getCurrency(arguments, EXCHANGE_RATE_API_KEY)

		if err != nil {
			return nil, err
		}

		res, err := json.Marshal(exchangeResp)

		if err != nil {
			return nil, err
		}

		msgContext := fmt.Sprintf("O Resultado da conversão do valor %s para %s é: %s", arguments.From, arguments.To, string(res))

		return mcp_golang.NewToolResponse(mcp_golang.NewTextContent(msgContext)), nil
	})

	if err != nil {
		panic(err)
	}

	err = mcpServer.Serve()

	if err != nil {
		panic(err)
	}

	<-done
}
