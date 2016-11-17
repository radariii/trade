package main

import (
	"errors"
	"fmt"

	"encoding/json"
	"strconv"

	"os"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

// Order is a structure that defines an order
type Order struct {
	OrderTimestamp   int     `json:"orderTimestamp"`
	ShippedTimestamp int     `json:"shippedTimestamp"`
	ArrivedTimestamp int     `json:"arrivedTimestamp"`
	Quantity         int     `json:"quantity"`
	TotalPrice       float32 `json:"totalPrice"`
}

// Producer is a structure that defines a producer
type Producer struct {
	Name             string  `json:"name"`
	CurrentInventory int     `json:"currentInventory"`
	Orders           []Order `json:"orders"`
}

type storedObject interface {
}

// Init is where it all begins
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

func (t *SimpleChaincode) producerFactory() storedObject {
	return Producer{Name: "hello", CurrentInventory: 100}
}

// Invoke is where new things happen
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	var producerName string
	var producer Producer

	if function == "harvestCoffee" {
		var coffeeAmtHarvested int
		if len(args) != 2 {
			return nil, errors.New("Incorrect number of arguments. Expecting 2")
		}
		producerName = args[0]
		coffeeAmtHarvested, _ = strconv.Atoi(args[1])

		// Initialize this producer if not already there
		producerBytes, _ := stub.GetState(producerName)
		if producerBytes == nil {
			// producer not found
			producer = Producer{Name: producerName, CurrentInventory: 0}
		} else {
			json.Unmarshal(producerBytes, &producer)
		}

		producer.CurrentInventory = producer.CurrentInventory + coffeeAmtHarvested

		outputStr := fmt.Sprintf("Producer %s just harvested %d pounds of coffee beans. Current Inventory = %d", producerName, coffeeAmtHarvested, producer.CurrentInventory)

		producerOut, _ := json.Marshal(producer)
		stub.PutState(producerName, producerOut)
		return []byte(outputStr), nil

	} else if function == "buyCoffee" {
		if len(args) != 4 {
			return nil, errors.New("Incorrect number of arguments. Expecting 4")
		}
		producerName = args[0]
		buyerName := args[1]
		amountPurchased, _ := strconv.Atoi(args[2])
		totalPrice, _ := strconv.Atoi(args[3])

		var producer = t.get(stub, producerName, t.producerFactory).(Producer) // Type assertion to Producer
		producer.CurrentInventory -= amountPurchased

		fmt.Fprintf(os.Stderr, "Buyer '%s' just purchased %d units from Producer '%s' for %d, leaving it with %d units in inventory of coffee beans. ", buyerName, amountPurchased, producerName, totalPrice, producer.CurrentInventory)
		producerOut, _ := json.Marshal(producer)
		stub.PutState(producerName, producerOut)

		return producerOut, nil

	} else if function == "shipCoffee" {

	} else if function == "coffeeArrives" {

	} else if function == "makePayment" {

	} else if function == "coffeeArrivesAtBorder" {

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

func (t *SimpleChaincode) get(stub shim.ChaincodeStubInterface, key string, factory func() storedObject) storedObject {
	// Initialize this producer if not already there
	storedObjectBytes, _ := stub.GetState(key)
	var storedObjectInst storedObject
	if storedObjectBytes == nil {
		storedObjectInst = factory()
	} else {
		json.Unmarshal(storedObjectBytes, &storedObjectInst)
	}
	return storedObjectInst
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
