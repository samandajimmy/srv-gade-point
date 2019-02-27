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

> **POST**   ${apiUrl}/campaigns

    Purposes :
    Create new campaign

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

> **PUT**   ${apiUrl}/campaigns/status/:id

    Purposes :
    Update status campaign
    0 --> INACTIVE
    1 --> ACTIVE

    Http Header :
    Content-Type: application/json

    Sample Payload :
    ${apiUrl}/campaigns/status/${campaign.id}

    {
	    "status" : 1        //integer
    }

    Success Response :
    {
        "status": "Success",
        "message": "Successfully Updated",
        "data": ""
    }

> **GET**   ${apiUrl}/campaigns

    Purposes :
    Get all campaign

    Http Header :
    Content-Type: application/json

    Sample:

    with params :
    name = Open                                 // name campaign
    status = 1                                  // status campaign
    startDate = 2019-03-11T11:13:52.958377Z     // start date campaign
    endDate = 2019-12-11T11:13:52.958377Z       // end date campaign

    ${apiUrl}/campaigns?name=Tabungan&startDate=2019-03-11T11:13:52.958377Z&endDate=2019-12-11T11:13:52.958377Z&status=1

    Success Response :
    {
        "status": "Success",
        "message": "Success",
        "data": [
            {
                "id": 10,
                "name": "Open Tabungan Emas",
                "description": "Open Tabungan Emas",
                "startDate": "2019-03-11T11:13:52.958377Z",
                "endDate": "2019-04-15T11:13:52.958377Z",
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
                "updatedAt": "2019-02-22T16:08:19.298535Z",
                "createdAt": "2019-02-22T15:54:03.721114Z"
            },
            {
                "id": 9,
                "name": "Open Tabungan Emas",
                "description": "Open Tabungan Emas",
                "startDate": "2019-03-11T11:13:52.958377Z",
                "endDate": "2019-04-15T11:13:52.958377Z",
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
                "createdAt": "2019-02-22T15:53:41.900267Z"
            }
        ]
    }

For External End Point :

> **POST** ${apiUrl}/campaigns/value

    Purposes :
    Get value point campaign result user point

    Http Header :
    Content-Type: application/json

    Sample Payload :
    {
        "userId": "001",            // string
        "channel": "01",            // string
        "product": "01",            // string
        "transactionType": "01",    // string
        "unit": "gram",             // string
        "transactionAmount": 1.80   // float
    }

    Success Response :
    {
        "status": "Success",
        "message": "Data Successfully Sent",
        "data": {
            "userPoint": 720
        }
    }

> **GET**   ${apiUrl}/campaigns/point?userId=NoUserId

    Purposes :
    Get value point amount user point

    Http Header :
    Content-Type: application/json

    Sample:

    with params :
    userId = 001                                 // No User Id

    Success Response :
    {
        "status": "Success",
        "message": "Success",
        "data": {
            "userPoint": 720
        }
    }

### Voucher
Create directory for public images :

    public/images/vouchers

add .env :

    VOUCHER_UPLOAD_PATH=./public/images/vouchers/

    VOUCHER_ROUTE_PATH=public/images/vouchers/

    VOUCHER_PATH=/images/vouchers
  
> **POST**   ${apiUrl}/vouchers

    Purposes :
    Create new voucher

    Http Header :
    Content-Type: application/json

    Sample Payload :
    {
        "name": "voucher emas",                             // string
        "description": "voucher emas potongan harga",       // string
        "startDate": "2019-02-10T22:08:41Z",                // Timestamp format RFC3339Nano
        "endDate": "2019-04-30T22:08:41Z",                  // Timestamp format RFC3339Nano
        "point": 100,                                       // integer
        "journalAccount": "000025130101360",                // string
        "value": 20000,                                     // float
        "imageUrl": "public/images/test.png",               // string
        "status": 1,                                        // integer
        "stock": 20,                                        // integer
        "prefixPromoCode": "EM",                            // string
        "validators": {
            "channel": "001",                               // string
            "product": "002",                               // string
            "transactionType": "003",                       // string
            "unit": "gram"                                  // string
        }
    }

    Success Response :
    {
        "status": "Success",
        "message": "Successfully Saved",
        "data": {
            "id": 1,
            "name": "voucher emas",
            "description": "voucher emas potongan harga",
            "startDate": "2019-02-10T22:08:41Z",
            "endDate": "2019-04-30T22:08:41Z",
            "point": 100,
            "journalAccount": "000025130101360",
            "value": 20000,
            "imageUrl": "public/images/test.png",
            "status": 1,
            "stock": 20,
            "prefixPromoCode": "EM",
            "validators": {
                "channel": "001",
                "product": "002",
                "transactionType": "003",
                "unit": "gram"
            },
            "updatedAt": "0001-01-01T00:00:00Z",
            "createdAt": "2019-02-26T08:18:28.092667717+07:00"
        }
    }

> **PUT**   ${apiUrl}/vouchers/status/:id

    Purposes :
    Update status voucher
    0 --> INACTIVE
    1 --> ACTIVE

    Http Header :
    Content-Type: application/json

    Sample Payload :
    ${apiUrl}/vouchers/status/${voucher.id}

    {
	    "status" : 1        //integer
    }

    Success Response :
    {
        "status": "Success",
        "message": "Successfully Updated",
        "data": ""
    }

> **GET**   ${apiUrl}/vouchers

    Purposes :
    Get all vouchers

    Http Header :
    Content-Type: application/json

    Sample:

    with params :
    name = Voucher                              // name vouchers
    status = 1                                  // status vouchers
    startDate = 2019-02-11T11:13:52.958377Z     // start date vouchers
    endDate = 2019-12-11T11:13:52.958377Z       // end date vouchers

    ${apiUrl}/vouchers?name=Voucher&startDate=2019-02-11T11:13:52.958377Z&endDate=2019-12-11T11:13:52.958377Z&status=1

    Success Response :
    {
        "status": "Success",
        "message": "Success",
        "data": [
            {
                "id": 2,
                "name": "voucher name",
                "description": "voucher description",
                "startDate": "2019-02-10T22:08:41Z",
                "endDate": "2019-04-30T22:08:41Z",
                "point": 100,
                "journalAccount": "000025130101360     ",
                "value": 20000,
                "imageUrl": "public/images/test.png",
                "status": 1,
                "validators": {
                    "channel": "001",
                    "product": "002",
                    "transactionType": "003",
                    "unit": "gram"
                },
                "updatedAt": "0001-01-01T00:00:00Z",
                "createdAt": "2019-02-25T11:12:48.366314Z"
            },
            {
                "id": 1,
                "name": "voucher name",
                "description": "voucher description",
                "startDate": "2019-02-10T22:08:41Z",
                "endDate": "2019-04-30T22:08:41Z",
                "point": 100,
                "journalAccount": "000025130101360     ",
                "value": 20000,
                "imageUrl": "public/images/test.png",
                "status": 1,
                "validators": {
                    "channel": "001",
                    "product": "002",
                    "transactionType": "003",
                    "unit": "gram"
                },
                "updatedAt": "2019-02-25T11:00:41.986737Z",
                "createdAt": "2019-02-25T10:56:35.534745Z"
            }
        ]
    }

> **POST**   ${apiUrl}/vouchers/uploadImage

    Purposes :
    Upload image voucher

    Http Header :
    Content-Type: multipart/form-data

    Sample :
    body : 
    form-data
    - key = file value = namaFile

    Success Response :
    {
        "status": "Success",
        "message": "Successfully Upload",
        "data": {
            "imageUrl": "/images/vouchers/1551068264320753609.png"
        }
    }
