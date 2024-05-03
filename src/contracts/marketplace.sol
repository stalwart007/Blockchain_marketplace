pragma solidity ^0.5.0;

contract Marketplace {
    string public name;

    uint public productCount = 0;
    mapping(uint => Product) public products;


    struct Product {
        uint id;
        string name;
        uint price;
        address payable owner;
        bool purchased;
    }

    event ProductCreated(
        uint id,
        string name,
        uint price,
        address payable owner,
        bool purchased
    );

    event ProductPurchased(
        uint id,
        string name,
        uint price,
        address payable owner,
        bool purchased
    );

    constructor() public {
        name = "Dapp";
    }

    function createProduct(string memory _name, uint _price) public{
        // require a valid name
        require(bytes(_name).length > 0);
        // Require a valid price
        require(_price > 0);
        // Make sure parameters are correct
        // Increment product count
        productCount++;
        // Create the product
        products[productCount] = Product(productCount, _name, _price, msg.sender, false);
        // Trigger an event
        emit ProductCreated(productCount , _name, _price, msg.sender, false);
    }

    function purchaseProduct(uint _id) public payable{
        // fetch the product
        Product memory _product = products[_id];
        // fetch the owner
        address payable _seller = _product.owner;
        // make sure the product is valid
        require(_product.id > 0 && _product.id <= productCount);
        // check if there is enough ether in the transaction
        require(msg.value >= _product.price);
        // check that the product is not already purchased
        require(!_product.purchased);
        // check that the buyer is not the seller
        require(_seller != msg.sender);
        // Purchase it (Transfer ownership to the buyer)
        _product.owner = msg.sender;
        // Mark as purchased
        _product.purchased = true;
        // update the product
        products[_id] = _product;
        // Pay the seller by sendeing them ether
        address(_seller).transfer(msg.value);
        // trigger an event
        emit ProductPurchased(productCount, _product.name, _product.price, msg.sender, true);

        

    }
}