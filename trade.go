package main

import (
	"errors"
	"fmt"

	"encoding/json"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type Order struct {
	orderTimestamp   int     `json:"orderTimestamp"`
	shippedTimestamp int     `json:"shippedTimestamp"`
	arrivedTimestamp int     `json:"arrivedTimestamp"`
	quantity         int     `json:"quantity"`
	totalPrice       float32 `json:"totalPrice"`
}

type Producer struct {
	name             string  `json:"name"`
	currentInventory int     `json:"currentInventory"`
	orders           []Order `json:"orders"`
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	// When we start, we want to initialize the global state

	// var A, B string    // Entities
	// var Aval, Bval int // Asset holdings
	// var err error

	// if len(args) != 4 {
	// 	return nil, errors.New("Incorrect number of arguments. Expecting 4")
	// }

	// // Initialize the chaincode
	// A = args[0]
	// Aval, err = strconv.Atoi(args[1])
	// if err != nil {
	// 	return nil, errors.New("Expecting integer value for asset holding")
	// }
	// B = args[2]
	// Bval, err = strconv.Atoi(args[3])
	// if err != nil {
	// 	return nil, errors.New("Expecting integer value for asset holding")
	// }
	// fmt.Printf("Aval = %d, Bval = %d\n", Aval, Bval)

	// // Write the state to the ledger
	// err = stub.PutState(A, []byte(strconv.Itoa(Aval)))
	// if err != nil {
	// 	return nil, err
	// }

	// err = stub.PutState(B, []byte(strconv.Itoa(Bval)))
	// if err != nil {
	// 	return nil, err
	// }

	return nil, nil
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	var producerName string
	var coffeeAmtHarvested int
	var producer Producer

	if function == "harvestCoffee" {
		if len(args) != 2 {
			return nil, errors.New("Incorrect number of arguments. Expecting 2")
		}
		producerName = args[0]
		coffeeAmtHarvested, _ = strconv.Atoi(args[1])

		// Initialize this producer if not already there
		producerBytes, _ := stub.GetState(producerName)
		if producerBytes == nil {
			// producer not found
			producer = Producer{name: producerName, currentInventory: 0}
		} else {
			json.Unmarshal(producerBytes, &producer)
		}

		producer.currentInventory = producer.currentInventory + coffeeAmtHarvested

		outputStr := fmt.Sprintf("Producer %s just harvested %d pounds of coffee beans. Current Inventory = %d", producerName, coffeeAmtHarvested, producer.currentInventory)

		producerOut, _ := json.Marshal(producer)
		stub.PutState(producerName, producerOut)
		return []byte(outputStr), nil
		//stub.PutState(producerName, []byte("HelloWorld"))

	}

	// if function == "delete" {
	// 	// Deletes an entity from its state
	// 	return t.delete(stub, args)
	// }

	// var A, B string    // Entities
	// var Aval, Bval int // Asset holdings
	// var X int          // Transaction value
	// var err error

	// if len(args) != 3 {
	// 	return nil, errors.New("Incorrect number of arguments. Expecting 3")
	// }

	// A = args[0]
	// B = args[1]

	// // Get the state from the ledger
	// // TODO: will be nice to have a GetAllState call to ledger
	// Avalbytes, err := stub.GetState(A)
	// if err != nil {
	// 	return nil, errors.New("Failed to get state")
	// }
	// if Avalbytes == nil {
	// 	return nil, errors.New("Entity not found")
	// }
	// Aval, _ = strconv.Atoi(string(Avalbytes))

	// Bvalbytes, err := stub.GetState(B)
	// if err != nil {
	// 	return nil, errors.New("Failed to get state")
	// }
	// if Bvalbytes == nil {
	// 	return nil, errors.New("Entity not found")
	// }
	// Bval, _ = strconv.Atoi(string(Bvalbytes))

	// // Perform the execution
	// X, err = strconv.Atoi(args[2])
	// if err != nil {
	// 	return nil, errors.New("Invalid transaction amount, expecting a integer value")
	// }
	// Aval = Aval - X
	// Bval = Bval + X
	// fmt.Printf("Aval = %d, Bval = %d\n", Aval, Bval)

	// // Write the state back to the ledger
	// err = stub.PutState(A, []byte(strconv.Itoa(Aval)))
	// if err != nil {
	// 	return nil, err
	// }

	// err = stub.PutState(B, []byte(strconv.Itoa(Bval)))
	// if err != nil {
	// 	return nil, err
	// }

	return nil, nil
}

// Deletes an entity from state
func (t *SimpleChaincode) delete(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	A := args[0]

	// Delete the key from the state in ledger
	err := stub.DelState(A)
	if err != nil {
		return nil, errors.New("Failed to delete state")
	}

	return nil, nil
}

// Query callback representing the query of a chaincode
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	if function == "getProducer" {
		if len(args) != 1 {
			return nil, errors.New("Incorrect number of arguments. Expecting 1")
		}
		producerName := args[0]
		producerBytes, _ := stub.GetState(producerName)
		if producerBytes == nil {
			return nil, errors.New("No producer found with name" + producerName)
		}

		return producerBytes, nil
	}

	return nil, errors.New("Incorrect query function name")

	// if function != "query" {
	// 	return nil, errors.New("Invalid query function name. Expecting \"query\"")
	// }
	// var A string // Entities
	// var err error

	// if len(args) != 1 {
	// 	return nil, errors.New("Incorrect number of arguments. Expecting name of the person to query")
	// }

	// A = args[0]

	// // Get the state from the ledger
	// Avalbytes, err := stub.GetState(A)
	// if err != nil {
	// 	jsonResp := "{\"Error\":\"Failed to get state for " + A + "\"}"
	// 	return nil, errors.New(jsonResp)
	// }

	// if Avalbytes == nil {
	// 	jsonResp := "{\"Error\":\"Nil amount for " + A + "\"}"
	// 	return nil, errors.New(jsonResp)
	// }

	// jsonResp := "{\"Name\":\"" + A + "\",\"Amount\":\"" + string(Avalbytes) + "\"}"
	// fmt.Printf("Query Response:%s\n", jsonResp)
	// return Avalbytes, nil
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
