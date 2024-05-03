package manufacturer

import (
    "encoding/json"
    "fmt"

    "github.com/hyperledger/fabric/core/chaincode/shim"
    pb "github.com/hyperledger/fabric/protos/peer"
)

type ManufacturerChaincode struct {
}

func (manufacturer *ManufacturerChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
    fmt.Println("Manufacturer chaincode initialized")
    return shim.Success(nil)
}

func (manufacturer *ManufacturerChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
    function, args := stub.GetFunctionAndParameters()
    fmt.Println("Invoked function:", function)

    switch function {
    case "RegisterMedicine":
        return manufacturer.registerMedicine(stub, args)
    case "UpdateMedicine":
        return manufacturer.updateMedicine(stub, args)
    case "GetMedicineDetail":
        return manufacturer.getMedicineDetail(stub, args)
    case "TransferMedicine":
        return manufacturer.transferMedicine(stub, args)
    default:
        return shim.Error("Invalid function name")
    }
}

// Implement specific functions for registering, updating, and transferring medicines

func (manufacturer *ManufacturerChaincode) registerMedicine(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    // Validate number of arguments
    if len(args) != 6 {
        return shim.Error("Incorrect number of arguments for RegisterMedicine")
    }

    // Parse and store medicine details
    medicineID := args[0]
    name := args[1]
    manufacturerName := args[2]
    location := args[3]
    status := args[4]
    currentOwner := args[5]

    medicineData, err := json.Marshal(Medicine{
        ID:          medicineID,
        Name:        name,
        Manufacturer: manufacturerName,
        Location:    location,
        Status:      status,
        CurrentOwner: currentOwner,
    })
    if err != nil {
        return shim.Error("Error marshalling medicine data")
    }

    // Put medicine data into blockchain state
    err = stub.PutState(medicineID, medicineData)
    if err != nil {
        return shim.Error("Error storing medicine data")
    }

    return shim.Success(nil)
}

func (manufacturer *ManufacturerChaincode) getMedicineDetail(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    // Validate arguments
    if len(args) != 1 {
        return shim.Error("Incorrect number of arguments for GetMedicineDetail")
    }

    medicineID := args[0]

    // Retrieve medicine data
    medicineData, err := stub.GetState(medicineID)
    if err != nil {
        return shim.Error("Error retrieving medicine data")
    }

    // Check if medicine exists
    if medicineData == nil {
        return shim.Error("Medicine not found")
    }

    // Unmarshal medicine data
    var medicine Order
    err = json.Unmarshal(medicineData, &medicine)
    if err != nil {
        return shim.Error("Error unmarshalling medicine data")
    }

    // Create and return detailed representation
    medicineDetail := &DetailResponse{
        ID:          medicine.ID,
        Name:        medicine.Name,
        Manufacturer: medicine.Manufacturer,
        Location:    medicine.Location,
        Status:      medicine.Status,
        Owner:       medicine.Owner,
        Transactions: medicine.Transactions,
    }

    medicineDetailBytes, err := json.Marshal(medicineDetail)
    if err != nil {
        return shim.Error("Error marshalling detailed representation")
    }

    return shim.Success(medicineDetailBytes)
}


func (manufacturer *ManufacturerChaincode) updateMedicine(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    // Validate arguments
    if len(args) != 3 {
        return shim.Error("Incorrect number of arguments for UpdateMedicine")
    }

    medicineID := args[0]
    fieldName := args[1]
    fieldValue := args[2]

    // Retrieve medicine data
    medicineData, err := stub.GetState(medicineID)
    if err != nil {
        return shim.Error("Error retrieving medicine data")
    }

    // Check if medicine exists
    if medicineData == nil {
        return shim.Error("Medicine not found")
    }

    // Unmarshal medicine data
    var medicine Order
    err = json.Unmarshal(medicineData, &medicine)
    if err != nil {
        return shim.Error("Error unmarshalling medicine data")
    }

    // Validate update based on field and value
    switch fieldName {
    case "Location":
        // Perform location update logic
    case "Status":
        // Perform status update logic
    default:
        return shim.Error("Invalid field for update")
    }

    // Update medicine data and store
    medicine.Location = fieldValue // Replace with actual update logic
    updatedMedicineData, err := json.Marshal(medicine)
    if err != nil {
        return shim.Error("Error marshalling updated medicine data")
    }
    err = stub.PutState(medicineID, updatedMedicineData)
    if err != nil {
        return shim.Error("Error storing updated medicine data")
    }

    return shim.Success([]byte("Medicine updated successfully"))
}


func (manufacturer *ManufacturerChaincode) transferMedicine(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    // Validate arguments
    if len(args) != 3 {
        return shim.Error("Incorrect number of arguments for TransferMedicine")
    }

    medicineID := args[0]
    newOwner := args[1]
    quantity := args[2]

    // Retrieve medicine data
    medicineData, err := stub.GetState(medicineID)
    if err != nil {
        return shim.Error("Error retrieving medicine data")
    }

    // Check if medicine exists and has sufficient quantity
    if medicineData == nil || medicine.Quantity < quantity {
        return shim.Error("Medicine not found or insufficient quantity")
    }

    // Unmarshal medicine data
    var medicine Medicine // Replace with your actual Medicine struct within the manufacturer package
    err = json.Unmarshal(medicineData, &medicine)
    if err != nil {
        return shim.Error("Error unmarshalling medicine data")
    }

    // Check ownership authorization (e.g., manufacturer can only transfer owned medicines)
    // Implement your logic for verifying transfer authorization here

    // Update medicine data with new owner and quantity
    medicine.CurrentOwner = newOwner
    medicine.Quantity -= quantity // Update logic based on your actual quantity management

    // Update and store medicine data
    updatedMedicineData, err := json.Marshal(medicine)
    if err != nil {
        return shim.Error("Error marshalling updated medicine data")
    }
    err = stub.PutState(medicineID, updatedMedicineData)
    if err != nil {
        return shim.Error("Error storing updated medicine data")
    }

    // Create and store transaction record
    transactionData, err := json.Marshal(Transaction{
        MedicineID: medicineID,
        Quantity:   quantity,
        Source:    medicine.CurrentOwner, // Update based on your source tracking logic
        Destination: newOwner,
        Timestamp:  time.Now().Unix(), // Replace with your timestamp logic
    })
    if err != nil {
        return shim.Error("Error marshalling transaction data")
    }
    transactionID := fmt.Sprintf("%s_%d", medicineID, time.Now().UnixNano()) // Generate unique transaction ID
    err = stub.PutState(transactionID, transactionData)
    if err != nil {
        return shim.Error("Error storing transaction data")
    }

    return shim.Success([]byte("Medicine transferred successfully"))
}
