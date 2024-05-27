# Krakend Config Generator

Generates config of the given file based on the comments above the functions.
Comments should start  with 'krakend' prefix to be recognized by the config generator.


## Valid Config Parameters

### Role
Defines who can access the service.

``` // krakend:Role:admin,user ```

Sets the method of the endpoint.

``` // krakend:Method:GET ```

Sets the endpoint URL. Query strings can be defined here. (eg. `/orders/{order_id}`)

``` // krakend:Endpoint:/orders ```

Sets the service ServiceName

``` // krakend:ServiceName: order-api ```

Sets which query strings krakend should accept.

``` // krakend:QueryStrings:user_id,page,limit,sort ```



## Example Usage

```go
	// krakend:Role:admin,user
	// krakend:Method:GET
	// krakend:Endpoint:/orders
	// krakend:ServiceName: order-api
	// krakend:QueryStrings:user_id,page,limit,sort
	app.Get("/orders", h.Index)

	// krakend:Role:admin,user
	// krakend:Method:POST
	// krakend:Endpoint:/orders
	// krakend:ServiceName: order-api
	app.Post("/orders", h.Create)

	// krakend:Role:admin,user
	// krakend:Method:GET
	// krakend:Endpoint:/orders:{order_id}
	// krakend:ServiceName: order-api
	app.Get("/orders/:id", h.Get)
```