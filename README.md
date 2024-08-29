# Go Currency Converter API

This is a simple Go-based API server that provides currency conversion using the [ExchangeRate-API](https://www.exchangerate-api.com/). The server listens on port `8080` and requires an API key to function.

## Getting Started

### Prerequisites

- Go (version 1.15 or higher)
- [ExchangeRate-API](https://www.exchangerate-api.com/) API Key

### Installation

1. Clone the repository:

    ```sh
    git clone https://github.com/a-drawn-pixel/go_currency_converter.git
    ```

2. Get your API key from [ExchangeRate-API](https://www.exchangerate-api.com/):

   - Sign up on their website.
   - You will receive an API key after registration.

3. Run the server:

    ```sh
    go run main.go <API_KEY>
    ```

   Replace `<API_KEY>` with your actual API key.

### Usage

The server will start listening on port `8080`. You can access the currency conversion endpoints by making requests to `http://localhost:8080`.

### Testing

To run the tests, use the following command:

```sh
npm run test
```
