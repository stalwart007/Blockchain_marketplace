package order

import (
    "encoding/json"
    "fmt"

    "github.com/hyperledger/fabric/core/chaincode/shim"
    pb "github.com/hyperledger/fabric/protos/peer"
)

type OrderChaincode struct {
}

func (order *OrderChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
    fmt.Println("Order chaincode initialized")
    return shim.Success(nil)
}

func (order *OrderChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
    function, args := stub.GetFunctionAndParameters()
    fmt.Println("Invoked function:", function)

    switch function {
    case "GetOrderDetail":
        return order.getOrderDetail(stub, args)
    case "UpdateOrderStatus":
        return order.updateOrderStatus(stub, args)
    case "AddTransaction":
        return order.addTransaction(stub, args)
    default:
        return shim.Error("Invalid function name")
    }
}

// Implement specific functions for getting order details, updating order status, and adding transactions

func (order *OrderChaincode) getOrderDetail(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    // Validate arguments
    if len(args) != 1 {
        return shim.Error("Incorrect number of arguments for GetOrderDetail")
    }

    orderID := args[0]

    // Get order data from state
    orderData, err := stub.GetState(orderID)
    if err != nil {
        return shim.Error("Error retrieving order data")
    }

    if orderData == nil {
        return shim.Error("Order not found")
    }

    // Unmarshal and return order details
    var order Order
    err = json.Unmarshal(orderData, &order)
    if err != nil {
        return shim.Error("Error unmarshalling order data")
    }

    return shim.Success(orderData)
}

func (order *OrderChaincode) updateOrderStatus(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    // Validate arguments
    if len(args) != 2 {
        return shim.Error("Incorrect number of arguments for UpdateOrderStatus")
    }

    orderID := args[0]
    newStatus := args[1]

    // Retrieve order data
    orderData, err := stub.GetState(orderID)
    if err != nil {
        return shim.Error("Error retrieving order data")
    }

    // Unmarshal order data
    var order Order
    err = json.Unmarshal(orderData, &order)
    if err != nil {
        return shim.Error("Error unmarshalling order data")
    }

    // Update order status
    order.Status = newStatus

    // Store updated order data
    updatedOrderData, err := json.Marshal(order)
    if err != nil {
        return shim.Error("Error marshalling order data")
    }
    err = stub.PutState(orderID, updatedOrderData)
    if err != nil {
        return shim.Error("Error storing order data")
    }

    return shim.Success([]byte("Order status updated successfully"))
}

func (order *OrderChaincode) addTransaction(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    // Validate arguments
    if len(args) != 4 {
        return shim.Error("Incorrect number of arguments for AddTransaction")
    }

    orderID := args[0]
    amount := args[1]
    paymentMethod := args[2]
    timestamp := args[3]

    // Retrieve order data
    orderData, err := stub.GetState(orderID)
    if err != nil {
        return shim.Error("Error retrieving order data")
    }

    // Unmarshal order data
    var order Order
    err = json.Unmarshal(orderData, &order)
    if err != nil {
        return shim.Error("Error unmarshalling order data")
    }

    // Create and add a new transaction
    transaction := Transaction{
        OrderID:    orderID,
        Amount:     amount,
        PaymentMethod: paymentMethod,
        Timestamp:   timestamp,
    }
    order.Transactions = append(order.Transactions, transaction)

    // Store updated order data
    updatedOrderData, err := json.Marshal(order)
    if err != nil {
        return shim.Error("Error marshalling order data")
    }
    err = stub.PutState(orderID, updatedOrderData)
    if err != nil {
        return shim.Error("Error storing order data")
    }

    return shim.Success([]byte("Transaction added successfully"))
}

func (order *OrderChaincode) changeOwnership(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    // Validate arguments
    if len(args) != 3 {
        return shim.Error("Incorrect number of arguments for ChangeOwnership")
    }

    orderID := args[0]
    newHash := args[1]
    rfidTag := args[2]

    // Retrieve order data
    orderData, err := stub.GetState(orderID)
    if err != nil {
        return shim.Error("Error retrieving order data")
    }

    // Unmarshal order data
    var order Order
    err = json.Unmarshal(orderData, &order)
    if err != nil {
        return shim.Error("Error unmarshalling order data")
    }

    // Verify RFID tag signature with existing hash
    validSignature, err := verifySignature(order.CurrentHash, rfidTag)
    if err != nil {
        return shim.Error("Error verifying RFID tag signature")
    }
    if !validSignature {
        return shim.Error("Invalid RFID tag signature")
    }

    // Update order ownership information
    order.CurrentHash = newHash
    order.Owner = args[3] 

}


