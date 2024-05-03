//const { assert } = require("chai")

require('chai').use(require('chai-as-promised')).should()

const Marketplace = artifacts.require('./Marketplace.sol')

contract ('Marketplace', ([deployer, seller, buyer]) => {
    let marketplace

    before(async () => {
        marketplace = await Marketplace.deployed()
    })

    describe('deployment', async () => {
        it('deplys successfully', async () => {
            const address = await marketplace.address
            assert.notEqual(address,0x0)
            assert.notEqual(address,'')
            assert.notEqual(address, null)
            assert.notEqual(address, undefined)
        })

        it('has a name', async () => {
            const name = await marketplace.name()
            assert.equal(name,'Dapp')
        })
    })

    describe('Products', async () => {

        let result, productCount
        before(async () => {
            result = await marketplace.createProduct("iPhone X", web3.utils.toWei('1', 'Ether'), { from:seller})
            productCount = await marketplace.productCount()
        })

        it('Creates Products', async () => {
            //SUCCESS
            assert.equal(productCount, 1)
            const event = result.logs[0].args
            assert.equal(event.id.toNumber(), productCount.toNumber(), 'ID is correct')
            assert.equal(event.name, 'nice', 'Name is correct')
            assert.equal(event.price, '1000000000000000000', 'Price is correct')
            assert.equal(event.owner, seller, 'Owner is correct')
            assert.equal(event.purchased, false, 'purchased is correct')
            // FAILURE
            await marketplace.createProduct("", web3.utils.toWei('1', 'Ether'), { from:seller}).should.be.rejected;
            await marketplace.createProduct("iPhone X", 0, { from:seller}).should.be.rejected;
        }) 

        it('Lists products', async () => {
            const product = await marketplace.products(productCount)
            assert.equal(product.id.toNumber(), productCount.toNumber(), 'ID is correct')
            assert.equal(product.name, 'iPhone X', 'Name is correct')
            assert.equal(product.price, '1000000000000000000', 'Price is correct')
            assert.equal(product.owner, seller, 'Owner is correct')
            assert.equal(product.purchased, false, 'purchased is correct')
        }) 

        it('Sells products', async () => {

            // Track the seller balance before purchase
            let oldSellerBalance
            oldSellerBalance = await web3.eth.getBalance(seller)
            oldSellerBalance = new web3.utils.BN(oldSellerBalance)
            
            // SUCCESS
            result = await marketplace.purchaseProduct(productCount, { from: buyer, value: web3.utils.toWei('1', 'Ether')})
            
            // check logs
            const event = result.logs[0].args
            assert.equal(event.id.toNumber(), productCount.toNumber(), 'ID is correct')
            assert.equal(event.name, 'iPhone X', 'Name is correct')
            assert.equal(event.price, '1000000000000000000', 'Price is correct')
            assert.equal(event.owner, buyer, 'Owner is correct')
            assert.equal(event.purchased, true, 'purchased is correct')

            //check that seller received funds
            let newSellerBalance
            newSellerBalance = await web3.eth.getBalance(seller)
            newSellerBalance = new web3.utils.BN(newSellerBalance)

            let price
            price = web3.utils.toWei('1','Ether')
            price = new web3.utils.BN(price)

            const expectedBalance = oldSellerBalance.add(price)

            assert.equal(newSellerBalance.toString(), expectedBalance.toString())

            // FAILURE

            // tries to buy a product that does not exist
            await marketplace.purchaseProduct(99, {from: buyer, value: web3.utils.toWei('1','Ether')}).should.be.rejected;
            // tries to buy without enough ether
            await marketplace.purchaseProduct(productCount, {from: buyer, value: web3.utils.toWei('0.5','Ether')}).should.be.rejected;
            // deplyer tries to buy the product
            await marketplace.purchaseProduct(productCount, {from: deployer, value: web3.utils.toWei('1','Ether')}).should.be.rejected;
            // Buyer tries to buy again
            await marketplace.purchaseProduct(productCount, {from: buyer, value: web3.utils.toWei('1','Ether')}).should.be.rejected;
        }) 

    })
})