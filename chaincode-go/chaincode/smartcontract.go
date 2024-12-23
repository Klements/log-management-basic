package chaincode

import (
    "encoding/json"
    "fmt"

    "github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

// SmartContract provides functions for managing an Log
type SmartContract struct {
    contractapi.Contract
}

// Log describes basic details of what makes up a simple asset
// Insert struct field in alphabetic order => to achieve determinism across languages
// golang keeps the order when marshal to json but doesn't order automatically
type Log struct {
    ID        string `json:"ID"`
    Timestamp string `json:"timestamp"`
    Server    string `json:"server"`
    User      string `json:"user"`
    Event     string `json:"event"`
    Port      int    `json:"port"`
    IPAddress string `json:"ipAddress"`
}

// InitLedger adds a base set of assets to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
    assets := []Log{
        {ID: "log1", Timestamp: "2020-06-25 01:53:37", Server: "mailsv1", User: "guest", Event: "access denied", Port: 20, IPAddress: "168.156.11.24"},
        {ID: "log2", Timestamp: "2020-08-19 16:51:15", Server: "mailsv1", User: "guest", Event: "access denied", Port: 22, IPAddress: "168.156.11.24"},
        {ID: "log3", Timestamp: "2020-06-06 18:57:31", Server: "websv2", User: "guest", Event: "access denied", Port: 80, IPAddress: "168.156.11.24"},
        {ID: "log4", Timestamp: "2021-04-04 21:47:12", Server: "mailsv1", User: "guest", Event: "access denied", Port: 20, IPAddress: "168.156.11.24"},
        {ID: "log5", Timestamp: "2022-04-08 05:12:14", Server: "mailsv1", User: "guest", Event: "access denied", Port: 20, IPAddress: "168.156.11.24"},
        {ID: "log6", Timestamp: "2024-05-06 10:28:58", Server: "mailsv1", User: "guest", Event: "access denied", Port: 20, IPAddress: "168.156.11.24"},
    }

    for _, asset := range assets {
        assetJSON, err := json.Marshal(asset)
        if err != nil {
            return err
        }

        err = ctx.GetStub().PutState(asset.ID, assetJSON)
        if err != nil {
            return fmt.Errorf("failed to put to world state. %v", err)
        }
    }

    return nil
}

// CreateAsset issues a new asset to the world state with given details.
func (s *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, id string, timestamp string, server string, user string, event string, port int, ipAddress string) error {
    exists, err := s.AssetExists(ctx, id)
    if err != nil {
        return err
    }
    if exists {
        return fmt.Errorf("the log %s already exists", id)
    }

    asset := Log{
        ID:        id,
        Timestamp: timestamp,
        Server:    server,
        User:      user,
        Event:     event,
        Port:      port,
        IPAddress: ipAddress,
    }
    assetJSON, err := json.Marshal(asset)
    if err != nil {
        return err
    }

    return ctx.GetStub().PutState(id, assetJSON)
}

// ReadAsset returns the asset stored in the world state with given id.
func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, id string) (*Log, error) {
    assetJSON, err := ctx.GetStub().GetState(id)
    if err != nil {
        return nil, fmt.Errorf("failed to read from world state: %v", err)
    }
    if assetJSON == nil {
        return nil, fmt.Errorf("the asset %s does not exist", id)
    }

    var asset Log
    err = json.Unmarshal(assetJSON, &asset)
    if err != nil {
        return nil, err
    }

    return &asset, nil
}

// AssetExists returns true when asset with given ID exists in world state
func (s *SmartContract) AssetExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
    assetJSON, err := ctx.GetStub().GetState(id)
    if err != nil {
        return false, fmt.Errorf("failed to read from world state: %v", err)
    }

    return assetJSON != nil, nil
}

// GetAllAssets returns all assets found in world state
func (s *SmartContract) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]*Log, error) {
    // range query with empty string for startKey and endKey does an
    // open-ended query of all assets in the chaincode namespace.
    resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
    if err != nil {
        return nil, err
    }
    defer resultsIterator.Close()

    var assets []*Log
    for resultsIterator.HasNext() {
        queryResponse, err := resultsIterator.Next()
        if err != nil {
            return nil, err
        }

        var asset Log
        err = json.Unmarshal(queryResponse.Value, &asset)
        if err != nil {
            return nil, err
        }
        assets = append(assets, &asset)
    }

    return assets, nil
}
