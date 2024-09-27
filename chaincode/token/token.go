package main

import (
    "encoding/json"
    "fmt"
    "github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
    contractapi.Contract
}

type Token struct {
    ID      string  `json:"id"`
    Balance float64 `json:"balance"`
}

func (s *SmartContract) InitToken(ctx contractapi.TransactionContextInterface, id string, amount float64) error {
    token := Token{ID: id, Balance: amount}
    tokenJSON, err := json.Marshal(token)
    if err != nil {
        return fmt.Errorf("failed to marshal token: %s", err.Error())
    }

    return ctx.GetStub().PutState(id, tokenJSON)
}

func (s *SmartContract) TransferToken(ctx contractapi.TransactionContextInterface, from string, to string, amount float64) error {
    fromTokenJSON, err := ctx.GetStub().GetState(from)
    if err != nil {
        return fmt.Errorf("failed to read from token: %s", err.Error())
    }
    if fromTokenJSON == nil {
        return fmt.Errorf("from token does not exist")
    }

    var fromToken Token
    err = json.Unmarshal(fromTokenJSON, &fromToken)
    if err != nil {
        return fmt.Errorf("failed to unmarshal from token: %s", err.Error())
    }

    if fromToken.Balance < amount {
        return fmt.Errorf("insufficient balance")
    }

    toTokenJSON, err := ctx.GetStub().GetState(to)
    if err != nil {
        return fmt.Errorf("failed to read to token: %s", err.Error())
    }

    var toToken Token
    if toTokenJSON == nil {
        toToken = Token{ID: to, Balance: 0}
    } else {
        err = json.Unmarshal(toTokenJSON, &toToken)
        if err != nil {
            return fmt.Errorf("failed to unmarshal to token: %s", err.Error())
        }
    }

    fromToken.Balance -= amount
    toToken.Balance += amount

    fromTokenJSON, err = json.Marshal(fromToken)
    if err != nil {
        return fmt.Errorf("failed to marshal from token: %s", err.Error())
    }

    toTokenJSON, err = json.Marshal(toToken)
    if err != nil {
        return fmt.Errorf("failed to marshal to token: %s", err.Error())
    }

    err = ctx.GetStub().PutState(from, fromTokenJSON)
    if err != nil {
        return fmt.Errorf("failed to update from token: %s", err.Error())
    }

    return ctx.GetStub().PutState(to, toTokenJSON)
}

func main() {
    chaincode, err := contractapi.NewChaincode(new(SmartContract))
    if err != nil {
        fmt.Printf("Error creating chaincode: %s", err.Error())
        return
    }

    if err := chaincode.Start(); err != nil {
        fmt.Printf("Error starting chaincode: %s", err.Error())
    }
}
