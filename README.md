# gade-srv-boileplate-go

## Description
This is an example of implementation of Clean Architecture in Go (Golang) projects.

Rule of Clean Architecture by Uncle Bob
 * Independent of Frameworks. The architecture does not depend on the existence of some library of feature laden software. This allows you to use such frameworks as tools, rather than having to cram your system into their limited constraints.
 * Testable. The business rules can be tested without the UI, Database, Web Server, or any other external element.
 * Independent of UI. The UI can change easily, without changing the rest of the system. A Web UI could be replaced with a console UI, for example, without changing the business rules.
 * Independent of Database. You can swap out Oracle or SQL Server, for Mongo, BigTable, CouchDB, or something else. Your business rules are not bound to the database.
 * Independent of any external agency. In fact your business rules simply don’t know anything at all about the outside world.

More at https://8thlight.com/blog/uncle-bob/2012/08/13/the-clean-architecture.html

This project has  4 Domain layer :
 * Models Layer
 * Repository Layer
 * UseCase Layer  
 * Delivery Layer

#### The diagram:

![golang clean architecture](https://gade/srv-gade-point/raw/master/clean-arch.png)

The explanation about this project's structure  can read from this medium's post : https://medium.com/@imantumorang/golang-clean-archithecture-efd6d7c43047


## Exposed Endpoint
### Campaign
End Point For Content Management System :

> **POST** /campaigns

    Purposes :
    Create New Campaign

    Http Header :
    Content-Type: application/json

    Sample Payload :
    {
        "name": "Open Tabungan Emas",                           // string
        "description": "Open Tabungan Emas",                    // string
        "startDate": "2019-03-11T11:13:52.958376536+07:00",     // Timestamp format RFC3339Nano
        "endDate": "2019-04-15T11:13:52.958376536+07:00",       // Timestamp format RFC3339Nano
        "status": 1,                                            // integer
        "type": 0,                                              // integer
        "validators": {
            "channel": "01",                                    // string
            "product": "01",                                    // string
            "transactionType": "01",                            // string
            "unit": "gram",                                     // string
            "multiplier": 0.01,                                 // float
            "value": 6,                                         // integer
            "formula": "(transactionAmount/multiplier)*value"   // string
        }
    }

    Success Response :
    {
        "status": "Success",
        "message": "Successfully Saved",
        "data": {
            "id": 10,
            "name": "Open Tabungan Emas",
            "description": "Open Tabungan Emas",
            "startDate": "2019-03-11T11:13:52.958376536+07:00",
            "endDate": "2019-04-15T11:13:52.958376536+07:00",
            "status": 1,
            "type": 0,
            "validators": {
                "channel": "01",
                "product": "01",
                "transactionType": "01",
                "unit": "gram",
                "multiplier": 0.01,
                "value": 6,
                "formula": "(transactionAmount/multiplier)*value"
            },
            "updatedAt": "0001-01-01T00:00:00Z",
            "createdAt": "2019-02-22T15:54:03.739922726+07:00"
        }
    }

> **PUT** /campaigns/status/:id

> **GET** /campaigns?name=nameCampaign&status=0&startDate=timestamp&endDate=timestamp

For External End POint :

> **POST** /campaigns/value

> **GET** /campaigns/point?userId=NoUserId

### Voucher
Create directory for public images :

    public/images/vouchers

add .env :

    VOUCHER_UPLOAD_PATH=./public/images/vouchers/

    VOUCHER_ROUTE_PATH=public/images/vouchers/

    VOUCHER_PATH=/images/vouchers
  
> **POST** /vouchers

> **PUT** /statusVoucher/:id

> **GET** /vouchers/?name=nameCampaign&status=0&startDate=timestamp&endDate=timestamp

> **POST** /uploadVoucherImages
