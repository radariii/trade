package main

import (
	"bytes"
	"errors"
	"fmt"

	"encoding/json"
	"strconv"

	"log"

	"encoding/binary"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

// Order is a structure that defines an order
type Order struct {
	ID               string  `json:"id,omitempty"`
	OrderTimestamp   string  `json:"orderTimestamp,omitempty"`
	ShippedTimestamp string  `json:"shippedTimestamp,omitempty"`
	ArrivedTimestamp string  `json:"arrivedTimestamp,omitempty"`
	Quantity         int     `json:"quantity,omitempty"`
	TotalPrice       float32 `json:"totalPrice,omitempty"`
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

	// Create the table that will hold the orders
	var c1 = shim.ColumnDefinition{Name: "ID", Type: shim.ColumnDefinition_STRING, Key: true}
	var c2 = shim.ColumnDefinition{Name: "BuyerName", Type: shim.ColumnDefinition_STRING, Key: false}
	var c3 = shim.ColumnDefinition{Name: "SellerName", Type: shim.ColumnDefinition_STRING, Key: false}
	var c4 = shim.ColumnDefinition{Name: "Quantity", Type: shim.ColumnDefinition_UINT32, Key: false}
	var c5 = shim.ColumnDefinition{Name: "TotalPrice", Type: shim.ColumnDefinition_BYTES, Key: false}
	var columnDefs []*shim.ColumnDefinition

	columnDefs = append(columnDefs, &c1, &c2, &c3, &c4, &c5)
	err := stub.CreateTable("Orders", columnDefs)
	if err != nil {
		return nil, errors.New("Failed to initialize order table")
	}

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

func (t *SimpleChaincode) producerFactory(producerName string) Producer {
	return Producer{Name: producerName, CurrentInventory: 100}
}

// Invoke is where new things happen
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	var producerName string
	var err error

	if function == "harvestCoffee" {
		var coffeeAmtHarvested int
		if len(args) != 2 {
			return nil, errors.New("Incorrect number of arguments. Expecting 2")
		}
		producerName = args[0]
		coffeeAmtHarvested, _ = strconv.Atoi(args[1])

		producer, err := t.getProducer(stub, producerName)
		if err != nil {
			return nil, errors.New("harvestCoffee: Error retrieving producer: " + err.Error())
		}

		producer.CurrentInventory = producer.CurrentInventory + coffeeAmtHarvested

		log.Printf("Producer %s just harvested %d pounds of coffee beans. Current Inventory = %d", producerName, coffeeAmtHarvested, producer.CurrentInventory)

		producerOut, _ := json.Marshal(producer)
		stub.PutState(producerName, producerOut)
		return producerOut, nil

	} else if function == "buyCoffee" {
		if len(args) != 3 {
			return nil, errors.New("Incorrect number of arguments. Expecting 3")
		}
		producerName = args[0]
		buyerName := args[1]
		var order Order
		log.Printf("buyCoffee: order JSON input: '%s'", args[2])
		err = json.Unmarshal([]byte(args[2]), &order)
		if err != nil {
			return nil, errors.New("buyCoffee: Error parsing order: " + err.Error())
		}

		producer, err := t.getProducer(stub, producerName)
		if err != nil {
			return nil, errors.New("buyCoffee: Error retrieving producer: " + err.Error())
		}

		if producer.CurrentInventory < order.Quantity {
			return nil, errors.New("Could not complete purchase. Producer has insufficient inventory.")
		}

		//var producer = t.get(stub, producerName, t.producerFactory).(Producer) // Type assertion to Producer
		producer.CurrentInventory -= order.Quantity

		buf := new(bytes.Buffer)
		binary.Write(buf, binary.BigEndian, order.TotalPrice)
		_, err = stub.InsertRow("Orders", shim.Row{
			Columns: []*shim.Column{
				&shim.Column{Value: &shim.Column_String_{String_: order.ID}},
				&shim.Column{Value: &shim.Column_String_{String_: buyerName}},
				&shim.Column{Value: &shim.Column_String_{String_: producerName}},
				&shim.Column{Value: &shim.Column_Uint32{Uint32: uint32(order.Quantity)}},
				&shim.Column{Value: &shim.Column_Bytes{Bytes: buf.Bytes()}},
			}})
		if err != nil {
			return nil, errors.New("An error occurred adding the new order to the Orders table: " + err.Error())
		}

		producer.Orders = append(producer.Orders, order)

		log.Printf("Buyer '%s' just purchased %d units from Producer '%s' for %f, leaving it with %d units in inventory of coffee beans. ", buyerName, order.Quantity, producerName, order.TotalPrice, producer.CurrentInventory)
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

func (t *SimpleChaincode) getProducer(stub shim.ChaincodeStubInterface, producerName string) (Producer, error) {
	var producer Producer
	producerBytes, err := stub.GetState(producerName)
	if err != nil {
		return producer, err // Empty struct, but err indicates a problem
	}
	if producerBytes == nil {
		producer = t.producerFactory(producerName)
	} else {
		err = json.Unmarshal(producerBytes, &producer)
		if err != nil {
			return producer, err
		}
	}
	return producer, nil
}

func (t *SimpleChaincode) get(stub shim.ChaincodeStubInterface, key string, factory func() storedObject) storedObject {
	// Initialize this producer if not already there
	var storedObjectInst = factory()
	storedObjectBytes, _ := stub.GetState(key)
	if storedObjectBytes != nil {
		json.Unmarshal(storedObjectBytes, &storedObjectInst)
		//var producer Producer
		// json.Unmarshal(storedObjectBytes, &producer)
		// storedObjectInst = producer
		// var loadedObj = reflect.Zero(objType)
		// json.Unmarshal(storedObjectBytes, &loadedObj)
		// storedObjectInst = loadedObj
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
