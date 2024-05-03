package buyer

import (
    "encoding/json"
    "fmt"

    "github.com/hyperledger/fabric/core/chaincode/shim"
    pb "github.com/hyperledger/fabric/protos/peer"
)

type BuyerChaincode struct {
}

func (buyer *BuyerChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
    fmt.Println("Buyer chaincode initialized")
    return shim.Success(nil)
}

func (buyer *BuyerChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
    function, args := stub.GetFunctionAndParameters()
    fmt.Println("Invoked function:", function)

    switch function {
    case "PlaceOrder":
        return buyer.placeOrder(stub, args)
    case "ConfirmOrder":
        return buyer.confirmOrder(stub, args)
    case "ReceiveMedicine":
        return buyer.receiveMedicine(stub, args)
    case "SellMedicine":
        return buyer.sellMedicine(stub, args)
    default:
        return shim.Error("Invalid function name")
    }
}


func (buyer *BuyerChaincode) placeOrder(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    // Validate number of arguments
    if len(args) != 4 {
        return shim.Error("Incorrect number of arguments for PlaceOrder")
    }

    // Parse and store order details
    medicineID := args[0]
    manufacturerID := args[1]
    quantity := args[2]

    orderData, err := json.Marshal(Order{
        MedicineID: medicineID,
        SellerID:   manufacturerID,
        BuyerID:    stub.GetCreator(), // Get caller's ID as
		quantity,
        Status:  "Placed",
        Transactions:  make([]Transaction, 0),
    })
    if err != nil {
        return shim.Error("Error marshalling order data")
    }

    // Put order data into blockchain state
    err = stub.PutState(stub.GetTxID(), orderData)
    if err != nil {
        return shim.Error("Error storing order data")
    }

    return shim.Success(nil)
}

func (buyer *BuyerChaincode) confirmOrder(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    // Validate arguments
    if len(args) != 1 {
        return shim.Error("Incorrect number of arguments for ConfirmOrder")
    }

    orderID := args[0]

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
    order.Status = "Confirmed"

    // Store updated order data
    updatedOrderData, err := json.Marshal(order)
    if err != nil {
        return shim.Error("Error marshalling order data")
    }
    err = stub.PutState(orderID, updatedOrderData)
    if err != nil {
        return shim.Error("Error storing order data")
    }

    return shim.Success([]byte("Order confirmed successfully"))
}

func (buyer *BuyerChaincode) receiveMedicine(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    // Validate arguments
    if len(args) != 2 {
        return shim.Error("Incorrect number of arguments for ReceiveMedicine")
    }

    medicineID := args[0]
    orderID := args[1]

    // Retrieve medicine data
    medicineData, err := stub.GetState(medicineID)
    if err != nil {
        return shim.Error("Error retrieving medicine data")
    }

    // Unmarshal medicine data
    var medicine Medicine
    err = json.Unmarshal(medicineData, &medicine)
    if err != nil {
        return shim.Error("Error unmarshalling medicine data")
    }

    // Update medicine's current owner
    medicine.CurrentOwner = stub.GetCreator() // Set to caller's ID

    // Update order status to "Delivered"
    orderData, err := stub.GetState(orderID)
    if err != nil {
        return shim.Error("Error retrieving order data")
    }
    var order Order
    err = json.Unmarshal(orderData, &order)
    if err != nil {
        return shim.Error("Error unmarshalling order data")
    }
    order.Status = "Delivered"

    // Store updated medicine and order data
    updatedMedicineData, err := json.Marshal(medicine)
    if err != nil {
        return shim.Error("Error marshalling medicine data")
    }
    err = stub.PutState(medicineID, updatedMedicineData)
    if err != nil {
        return shim.Error("Error storing medicine data")
    }
    updatedOrderData, err := json.Marshal(order)
    if err != nil {
        return shim.Error("Error marshalling order data")
    }
    err = stub.PutState(orderID, updatedOrderData)
    if err != nil {
        return shim.Error("Error storing order data")
    }

    // Trigger payment transaction (implementation depends on payment logic)
    // ...

    return shim.Success([]byte("Medicine received successfully"))
}

func (buyer *BuyerChaincode) sellMedicine(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    // Validate arguments
    if len(args) != 3 {
        return shim.Error("Incorrect number of arguments for SellMedicine")
    }

    medicineID := args[0]
    buyerID := args[1]
    quantity := args[2]

    // Retrieve medicine data to check ownership
    medicineData, err := stub.GetState(medicineID)
    if err != nil {
        return shim.Error("Error retrieving medicine data")
    }
    var medicine Medicine
    err = json.Unmarshal(medicineData, &medicine)
    if err != nil {
        return shim.Error("Error unmarshalling medicine data")
    }
    if medicine.CurrentOwner != stub.GetCreator() {
        return shim.Error("Caller does not own the medicine")
    }

    // Create a new order for the sale
    orderData, err := json.Marshal(Order{
        MedicineID: medicineID,
        SellerID:   stub.GetCreator(),
        BuyerID:    buyerID,
        Quantity:   quantity,
        Status:     "Placed",
        Transactions: make([]Transaction, 0),
    })
    if err != nil {
        return shim.Error("Error marshalling order data")
    }
    err = stub.PutState(stub.GetTxID(), orderData)
    if err != nil {
        return shim.Error("Error storing order data")
    }

    return shim.Success([]byte("New order for sale created successfully"))
}


