func (t *IPFSChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

// Invoke - Our entry point for Invocations
func (t *IPFSChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	switch function {
	case "storeOnIPFS":
		return t.storeOnIPFS(stub, args)
	case "retrieveFromIPFS":
		return t.retrieveFromIPFS(stub, args)
	default:
		return shim.Error("Invalid function name. Expecting 'storeOnIPFS' or 'retrieveFromIPFS'")
	}
}

// Store data on IPFS
func (t *IPFSChaincode) storeOnIPFS(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	ipfsShell := shell.NewShell("localhost:5001")
	data := []byte(args[0])

	// Add data to IPFS
	hash, err := ipfsShell.Add(bytes.NewReader(data))
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to add data to IPFS: %s", err))
	}

	return shim.Success([]byte(hash))
}

// Retrieve data from IPFS
func (t *IPFSChaincode) retrieveFromIPFS(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	ipfsShell := shell.NewShell("localhost:5001")
	hash := args[0]

	// Retrieve data from IPFS
	data, err := ipfsShell.Cat(hash)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to retrieve data from IPFS: %s", err))
	}

	return shim.Success(data)
}

func main() {
	err := shim.Start(new(IPFSChaincode))
	if err != nil {
		fmt.Printf("Error starting IPFSChaincode: %s", err)
	}
}
