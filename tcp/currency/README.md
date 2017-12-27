# Currency Service
This directory contains the implementation of 
simple currency service which returns ISO currency
information ([see currency website](http://www.currency-iso.org/en/home/tables/table-a1.html)).  

Currency client and server use JSON to exchange data.  The client sends request:
```JSON
{
    "get":<currency-name or code>,
    "limit":<limit>
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
in package [curlib](https://github.com/vladimirvivien/go-networking/blog/master/tcp/curlib/curlib.go).  The data source for the code is stored in text file data.csv.
