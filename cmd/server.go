package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"bitbucket.org/voranskripto/nameconverters/hotwalletconverter"
	"bitbucket.org/voranskripto/nameconverters/tickerconverter"

	"github.com/gorilla/mux"
	tb "github.com/labusk/tokenbalance"
)

func router() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/balance/{contract}/{wallet}", getBalanceHandler).Methods("GET")
	r.HandleFunc("/token/{contract}/{wallet}", getTokenHandler).Methods("GET")
	r.HandleFunc("/health", getHealthHandler)
	r.HandleFunc("/labustoken/{contract}/{wallet}", getLabusTokenHandler)
	return r
}

func startServer(ip string, port int) error {
	log.Printf("TokenBalance Server Running: http://%v:%v", ip, port)
	srv := &http.Server{
		Addr:         fmt.Sprintf("%v:%v", ip, port),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router(),
	}
	return srv.ListenAndServe()
}

func collectVars(r *http.Request) (string, string) {
	vars := mux.Vars(r)
	return vars["contract"], vars["wallet"]
}

func getHealthHandler(w http.ResponseWriter, r *http.Request) {
	var health map[string]interface{}
	chainId, err := tb.Geth.NetworkID(context.TODO())
	if err != nil {
		health = map[string]interface{}{
			"online": false,
			"chain":  0,
		}
	} else {
		health = map[string]interface{}{
			"online": true,
			"chain":  chainId.Int64(),
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(health)
}

func getTokenHandler(w http.ResponseWriter, r *http.Request) {
	contract, wallet := collectVars(r)
	log.Println("Fetching /token for Wallet:", wallet, "at Contract:", contract)
	query, err := tb.New(contract, wallet)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		m := errorResponse{
			Error:   true,
			Message: err.Error(),
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(m)
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(query.ToJSON()))
	}
}

func getBalanceHandler(w http.ResponseWriter, r *http.Request) {
	contract, wallet := collectVars(r)
	log.Println("Fetching /balance for Wallet:", wallet, "at Contract:", contract)
	query, err := tb.New(contract, wallet)
	w.Header().Set("Content-Type", "text/plain")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("0.0"))
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(query.BalanceString()))
	}
}

func convertVars(symbol, exchange string) (contract, wallet string, err error) {
	contract, err = tickerconverter.ConvertSymbol(symbol)
	if err != nil {
		return "", "", err
	}
	wallet, err = hotwalletconverter.ConvertExchangeName(exchange)
	if err != nil {
		return "", "", err
	}
	return
}

func getLabusTokenHandler(w http.ResponseWriter, r *http.Request) {
	symbol, exchange := collectVars(r)
	contract, wallet, err := convertVars(symbol, exchange)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err != nil {
		m := errorResponse{
			Error:   true,
			Message: err.Error(),
		}
		json.NewEncoder(w).Encode(m)
		return
	}
	query, err := tb.New(contract, wallet)
	if err != nil {
		m := errorResponse{
			Error:   true,
			Message: err.Error(),
		}
		json.NewEncoder(w).Encode(m)
		return
	}
	json.NewEncoder(w).Encode(query)

	log.Println("Fetching /balance for Wallet:", wallet, "at Contract:", contract)
}
