# BasicAuthProxyTransformer

![Logo](https://raw.githubusercontent.com/nickhartjes/basic-auth-proxy-transformer/main/docs/bapt-logo.png)

BasicAuthProxyTransformer is a specialized proxy application written in Go. It is designed to handle requests with Basic Authentication, transform them into JWT (JSON Web Tokens), and forward them to their original destination.

The application works by intercepting incoming requests and checking for the presence of a Basic Auth header. If such a header is found, the application removes it and makes a request to an OAuth2 server to obtain a JWT token. This token is then added to the original request's headers, and the request is forwarded to its original destination.

To prevent high load on the OAuth2 server and to increase response times, BasicAuthProxyTransformer utilizes caching. It supports both in-memory caching and Redis (via Valkey) caching. This allows the application to store and quickly retrieve JWT tokens, reducing the need for frequent requests to the OAuth2 server.

This process ensures that the downstream services receive requests with JWT tokens, which are more secure and versatile than Basic Auth credentials. This transformation process is transparent to the client making the request, making BasicAuthProxyTransformer a valuable tool for enhancing the security of your services without requiring changes to the client applications.

## Getting Started

These instructions will guide you through the process of setting up BasicAuthProxyTransformer for development and testing purposes.


### Prerequisites

- Go
- Docker
- Docker Compose

### Installing

Clone the repository:
```bash
git clone https://github.com/nickhartjes/BasicAuthProxyTransformer.git
cd BasicAuthProxyTransformer
```
Build the application:
```bash
make build
```
Build the Docker image:
```bash
make docker-build
```
Start the Docker Compose services:
```bash
make docker-up
```
## Running the tests

Run the tests:
```bash
make test
```
Generate the test coverage report:
```bash
make cover
```
## Built With

- [Go](https://golang.org/)
- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)

## Authors

- Nick Hartjes

## License
