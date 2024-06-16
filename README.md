
# Mini-Aspire API

## Overview
This project is a mini-aspire API designed to guide authenticated users through a comprehensive loan application process. The API provides the following functionalities:

1. **Creating Loan Applications:**
   - Users can create loan applications by providing the minimum required information to calculate their loan eligibility.

2. **Allocating Loans for Approval:**
   - The system allocates loans to employees for approval. If the approved amount exceeds the eligible amount, an override request must be made.

3. **Fetching Application List:**
   - Both users and employees can fetch a list of their respective loan applications.

4. **Viewing Loan Details:**
   - Users can view detailed loan information, including interest rates, principal outstanding, and EMI (Equated Monthly Installment) outstanding.

5. **Loan Approval by Employees:**
   - Employees can approve loan applications. This endpoint is restricted for user type "customer." During the approval process, if the approved amount is greater than the eligible amount, the employee must approve the amount with an override. The system will then recalculate the repayment schedule for the application.

6. **Making Repayments:**
   - Users can make repayments, which will adjust the outstanding amount and settle installments based on the payment amount.

These features are designed to ensure a streamlined and efficient loan management process, from application creation to approval and repayment.
## Features
- Loan application creation and participant management
- Loan approval and override request handling
- Repayment schedule generation
- Loan details viewing
- Repayment processing

## Project Structure
```
boilerplate
└── app
    └── /common
    └── /configs
    └── /controllers
        └── /loan
            └── controller.go       # It include loan create, loan details and approve loan api's controller
            └── controller_test.go  # Unit test case with apis
        └── /repayment
            └── controller.go       # It include repayment api
            └── controller_test.go  # Unit test case for repayment api
        └── /user
            └── controller.go       # It include user auth api which is common for both employee and user
            └── controller_test.go  # Unit test case for auth api
    └── /database
    └── /dto                        # In this directory we maintain Data transfer object (DTO) of all our apis
    └── /logger                     # Include customer logger for each api request with response
    └── /middleware                 # Middleware for auth and access restriction
    └── /models                     # In the directory we can find all the database models used in the project
        └── constant.go
        └── country.go
        └── country_currency.go
        └── currency.go
        └── loan_application.go
        └── loan_application_participant.go
        └── loan_eligibility_config.go
        └── payment.go
        └── repayment.go
        └── repayment_payment_log.go
        └── user.go
        └── user_kyc.go
    └── /routes                     # This directory include routes
        └── routers.go              # It initalise the route provider and setup route version
        └── v1.go                   # all the v1 routing and services initialise happens here
    └── /services                   # This is our service directory which contain all the business logic
        └── /loan
            └── mock_loan_service.go        # mockgen generated file for handing loan service
            └── service.go                  # loan service interface
            └── loan_service.go             # loan service methods
        └── /repayment
            └── mock_repayment_service.go   # mockgen generated file for handing repayment service
            └── service.go                  # repayment service interface
            └── repayment_service.go        # repayemnt service methods
        └── /user
            └── mock_user_service.go         # mockgen generated file for handing user service
            └── service.go                   # user service interface
            └── user_service.go              # user service methods
└── migration
    └── /schema                     # this directory contains all up and down sql file
    └── migration.go                # migration file execution
├── .env.example                    # Example environment variables
├── .gitignore                      # Git ignore file
├── Dockerfile                      # Docker configuration
├── docker-compose.yaml             # Docker compose file which contains env required for all to run
├── README.md                       # Project documentation
├── go.mod                          # Go module dependencies
├── go.sum                          # Go module checksum
        
```

## Getting Started

### Prerequisites
- [Go](https://golang.org/doc/install) 1.20 or later
- [Docker](https://docs.docker.com/get-docker/)

### Installation
1. Clone the repository:
    ```sh
    git clone https://github.com/nishanthrk/aspire-lms.git
    cd aspire-lms
    ```

2. Copy the example environment file and edit as needed:
    ```sh
    cp sample.env .env
    ```

3. Install dependencies:
    ```sh
    go mod download
    ```
### Running the Application
1. Run the application:
    ```sh
    go run main.go
    ```

2. Alternatively, you can use Docker:
    ```sh
    docker-compose up --build
    ```

### Environment Variables
- `ENV=development`: The application environment (e.g., development, production).
- `APP_PORT=8080`: The port on which the application will run.
- `APP_HOST=localhost`: The host address for the application.
- `ENV_LOAD_METHOD=LOCAL`: Method to load environment variables (e.g., LOCAL, REMOTE).
- `ENV_LOAD_PATH=`: Path to the environment file if `ENV_LOAD_METHOD` is set to LOCAL.
- `JWT_ACCESS_SIGN_KEY=greatest-secret-ever`: Secret key used for signing JWT access tokens.
- `JWT_REFRESH_SIGN_KEY=greatest-secret-ever`: Secret key used for signing JWT refresh tokens.
- `JWT_ISSUER=ASPIRE`: Issuer of the JWT tokens.
- `TENANT=ASPIRE`: Tenant name for the application.
- `NEW_RELIC_LICENSE=`: New Relic license key for monitoring.
- `DB_HOST=host.docker.internal`: Database host address.
- `DB_PORT=3306`: Database port number.
- `DB_DRIVER=mysql`: Database driver (e.g., mysql, postgres).
- `DB_USER=root`: Database username.
- `DB_PASSWORD=nishanth`: Database password.
- `DB_NAME=aspire_lms`: Database name.

## Postman Collection

A Postman collection is provided to test the API endpoints. You can import the collection into Postman using the following steps:

1. Open Postman.
2. Click on the Import button.
3. Select the `Aspire LMS.postman_collection.json` file from the root of the project directory.
4. Click Import to add the collection to your Postman workspace.

Now you can use the collection to test the different API endpoints. The collection includes requests for:

- Application Create
- Application List
- Approve Application
- Auth - Employee
- Auth - Customer
- Application Detail
- Repayment

You can find the Postman collection file in the collection directory.

## Running Tests

### Unit Tests

The following command are used to generate mock services:
- mockgen -source=app/services/loan/service.go -destination=app/services/loan/mock_loan_service.go -package=loan_service
- mockgen -source=app/services/user/service.go -destination=app/services/user/mock_user_service.go -package=user_service
- mockgen -source=app/services/repayment/service.go -destination=app/services/repayment/mock_repayment_service.go -package=repayment_service

To run unit tests for the controllers, you can use the following command:
```bash
go test ./app/controllers/...