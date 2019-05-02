# API Interview Test Project 
This is an example Golang API service using static data in a json file located in the mock data directory.

Instructions
1.	Create a document/text file with a list of possible assumptions that you would make based on the Feature Criteria Requested
2.	Note any questions you would surface to the ScrumMaster/PM or Business Stakeholders/Product Owners regarding the feature request(s)
3.	Create a API Solution in your prefererd language/technology that will use the provided json file (auto.leads.json) as the Minimally Viable Product iteration’s  dataset
4.	Provide detailed instructions as to: 
a.	How to install/configure/run your solution on any Windows machine.  
b.	Please explicitly define the expected URL to access your API via postman once the API solution is being served.
5.	Please keep the level of effort to complete the exercise reasonable between one and four hours.
Features
1.	Create a method/endpoint to return a collection of all list items
2.	Create a method/endpoint to retrieve item detail by ID
3.	Create a method/endpoint to return a collection of all list items filtering on consumer’s state value
4.	Create a method/endpoint to return a collection of all list items filtering on vehicle’s make value
5.	Create a method/endpoint to return a collection of all list items filtering on former insurer


Delivery
The preferred delivery method of this assignment would be to provide a GitHub or BitBucket URL that allows us to download & review your work.   

Decided to render the test both in Go lang and Node.  See the repo tetrainfo/APIDemoNode for the Node approach.

Go Notes

For the Go approach I decided to keep it really simple.  No dependancies other than what comes standard.  One file. Ideal for a microservice.

That meant that the query schema had to change.  The delivered http router only supports static paths. So a query where
the id was in the url like this /:id  had to change to ?id=.  While I was at it, I decided to a versioning schema to the url.
That adds a /v1 in front the quotes path.

Unlike JavaScript and Node, Go requires that json data be unmarshalled and marshalled into maps, or arrays of maps. Could have explicitly defined the data structure into various types.  Decided instead to use the Go empty interface{}. That means that the solution is much more flexible, requiring only the key query items to be in the json blob.

Go compiles to an exe.  In theory that's all that's required for a Windows machine.  The exec can be compiled from the source.

To make the executable from the source code download the repo and install Go at https://golang.org/dl/. Set the GOPATH.


```bash
$ git clone https://github.com/tetrainfo/APIDemoGo.git
```

then build the exe

```bash
$ go build
```

To run the project, navigate to the APIDemo folder and enter

```bash
$ ./APIDemoGo.exe

Other shells try .\APIDemoGo.exe
```

To test, use these sample urls with an tool like Postman at https://www.getpostman.com 
```
http://localhost:8080/v1/quotes?id=998 should return a single response for a customer with id = 998

http://localhost:8080/v1/quotes?state=IL should return a list for consumers with state=IL or state="IL"

http://localhost:8080/v1/quotes?make=ford should return a list for consumers have one or more Ford vehicles

http://localhost:8080/v1/quotes?former_insurer="Monolith Casualty" should return a list for consumers with the specified insurer

http://localhost:8080/v1/quotes?list=all should return the entire list

Recently added: the ability to send the parameters via form data


make:ford


or  json post (only one query will be satisfied)

{
    "state": "IL",
    "make" : "Ford",
    "former_insurer": "Monolith Casualty",
    "list": "all",
    "id": "998"
}


```To Do:

Make the search criteria a bit fuzzier using a wildcard
Implement real logging
Write tests

```

