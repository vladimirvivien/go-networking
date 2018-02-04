# Currency Service
This directory contains the implementations of 
a simple currency service which returns ISO currency
information ([see currency website](http://www.currency-iso.org/en/home/tables/table-a1.html))
stored in file [./data.csv](./data.csv).

Code in directories servertxtXX implement a simple structured text protocol.  
Directories with serverjsonXX implements the currency server using JSON-encoded
data.  The client sends request:

```JSON
{
    "get":<currency-name or code>
}
```
The server returns currencies information that
matches the request:
```JSON
[
    {
        "currency_code":<string>,
        "currency_name":<string>,
        "currency_number":<string>,
        "currency_number":<string>
    }
]
```

The supporting data types and functions are declared
in package [lib](https://github.com/vladimirvivien/go-networking/blog/master/currency/lib/curlib.go).