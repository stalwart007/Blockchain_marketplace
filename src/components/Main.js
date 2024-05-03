import React, { Component } from 'react';


class Main extends Component {

  render() {
    return (
        <div id="content">
            <h1>Add Product</h1>
            <form onSubmit={(event) => {
                event.preventDefault()
                const name = this.productName.value
                const price = window.web3.utils.toWei(this.productPrice.value.toString(), 'Ether')
                this.props.createProduct(name, price)
            }}>
            <div className="form-group mr-sm-2">
                <input
                id="productName"
                type="text"
                ref={ (input) => { this.productName = input } }
                className="form-control"
                placeholder="Product Name"
                required />
            </div>
            <div className="form-group mr-sm-2">
                <input
                id="productPrice"
                type="text"
                ref={ (input) => { this.productPrice = input } }
                className="form-control"
                placeholder="Product Price"
                required />
            </div>
            <button type="submit" className="btn btn-primary">Add Product</button>
            </form>


            <p> </p>


            <h1>Buy Product</h1>
            <table className="table">
            <thead>
                <tr>
                <th scope="col">#</th>
                <th scope="col">Name</th>
                <th scope="col">Price</th>
                <th scope="col">Owner</th>
                <th scope="col"></th>
                </tr>
            </thead>
            <tbody id="productList"> {
                this.props.products.map((product, key) => {
                    return(
                        <tr key = {key}>
                            <th scope="row">{product.id.toString()}</th>
                            <td>{product.name}</td>
                            <td>{window.web3.utils.fromWei(product.price.toString(), 'Ether')} Eth</td>
                            <td>{product.owner}</td>
                            <td><button className="buyButton"  name = {product.id} value = {product.price} onClick={(event) => {
                                this.props.purchaseProduct(event.target.name, event.target.price)
                            }}>
                            Buy
                            </button>
                            </td>
                        </tr>
                    )
                })
            }
                
            </tbody>
            </table>
            
            <h1>Track Details</h1>
            <form onSubmit={(event) => {
                event.preventDefault()
                
            }}>
            <div className="form-group mr-sm-2">
                <input
                id="productId"
                type="number"
                ref={ (input) => { this.productName = input } }
                className="form-control"
                placeholder="Product ID"
                required />
            </div>
            <button type="submit" className="btn btn-primary">Get Details</button>
            </form>

        </div>
    );
  }
}

export default Main;
