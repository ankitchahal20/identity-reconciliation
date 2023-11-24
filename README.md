# Identity Reconciliation Service

Problem Statement:
1. Develop a web service capable of processing JSON payloads with "email" and
"phoneNumber" fields, creating a shadowy infrastructure that consolidates contact
information across multiple purchases.
2. Craft a response mechanism that returns an HTTP 200 status code, along with a
JSON payload containing consolidated contact details. The payload should be
cunningly structured with a "primaryContactId," "emails," "phoneNumbers," and
"secondaryContactIds."
3. Be prepared for a scenario where no existing contacts match the incoming request. In
such cases, your service should craft a new entry in the database, treating it as a
discreet individual with no affiliations.

4. Implement a strategy where incoming requests matching existing contacts trigger the
creation of "secondary" contact entries, providing a covert mechanism for storing new
information.
5. Exercise caution as primary contacts can seamlessly transform into secondary contacts if subsequent requests reveal overlapping contact information. This dual-
purpose functionality adds an extra layer of complexity to your mission.

Requirements:
1. Meticulously implement the /identify endpoint, ensuring that it operates with utmost
discretion and precision.
2. Execute the creation of a new "Contact" entry with linkPrecedence="primary" when no
existing contacts match the incoming request.
3. Employ a covert strategy for creating "secondary" contact entries when incoming
requests match existing contacts and introduce new information.
4. Maintain the integrity of the database state, executing updates seamlessly with each
incoming request.

## Prerequisites

Before running the Identity Reconciliation Service, make sure you have the following prerequisites installed on your system:

- Go programming language (go1.21.3)
- PostgreSQL(14.8)

## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/ankitchahal20/identity-reconciliation.git
   ```

2. Navigate to the project directory:

   ```bash
   cd identity-reconciliation
   ```

3. Install the required dependencies:

   ```bash
   go mod tidy
   ```

4. Defaults.toml
Add the values to defaults.toml and execute `go run main.go` from the cmd directory.

5. DB setup
    ```
    Use the scripts inside sql-scripts directory to create the tables in your db.
    ```

## APIs
There is only 1 API which this repo currently supports.

Generate Identity Reconciliation
```
curl -i -k -X POST \
   http://127.0.0.1:8080/v1/ \
  -H "transaction-id: 288a59c1-b826-42f7-a3cd-bf2911a5c351" \
  -H "content-type: application/json" \
  -d '{
  "email":"hello2@gmail.com",
  "phone_number": "9410051646"
}'
```
Response
```
{
  "primary_contact_id": 113,
  "emails": [
    "hello2@gmail.com"
  ],
  "phone_numbers": [
    "+919410051646",
    "9410051646"
  ],
  "secondary_contact_ids": [
    114
  ]
}
```

## Project Structure

The project follows a standard Go project structure:

- `config/`: Configuration file for the application.
- `internal/`: Contains the internal packages and modules of the application.
  - `config/`: Global configuration which can be used anywhere in the application.
  - `constants/`: Contains constant values used throughout the application.
  - `db/`: Contains the database package for interacting with PostgreSQL.
  - `models/`: Contains the data models used in the application.
  - `middleware/`: Contains code for input and token validation
  - `identityReconciliationError`: Defines the errors in the application
  - `service/`: Contains the business logic and services of the application.
  - `server/`: Contains the server logic of the application.
  - `utils/`: Contains utility functions and helpers.
- `cmd/`:  Contains command you want to build.
    - `main.go`: Main entry point of the application.
- `README.md`: README.md contains the description for the Identity Reconciliation Service

## Contributing

Contributions to the Identity Reconciliation Service are welcome. If you find any issues or have suggestions for improvement, feel free to open an issue or submit a pull request.

## License

The Identity Reconciliation Service is open-source and released under the [MIT License](LICENSE).

## Contact

For any inquiries or questions, please contact:

- Ankit Chahal
- ankitchahal20@gmail.com

Feel free to reach out with any feedback or concerns.
