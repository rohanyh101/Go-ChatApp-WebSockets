# Real-Time Chat Application in Golang

This repository contains a real-time chat application built using Golang and the Gorilla WebSocket framework. The application includes several advanced features such as OTP (One Time Password) protection, jumbo frames, client-side egress for concurrent messaging, secure WebSockets (WSS) with locally generated certificates, room change functionality, and event-based messaging for scalability.

## Features

- **Authentication**: Users must log in to access the chat service. Authentication includes OTP protection for added security.
- **Jumbo Frames**: User message size is regulated to ensure efficient communication.
- **Client-Side Egress**: Implements egress on the client side to handle concurrent messaging efficiently.
- **Secure WebSockets (WSS)**: Communication is secured using WebSocket Secure (WSS) with locally generated certificates.
- **Room Change**: Users can change chat rooms seamlessly.
- **Event-Based Messaging**: Implements event-based messaging to enhance scalability.

## Prerequisites

- Go 1.16 or higher
- Gorilla WebSocket package
- OpenSSL (for generating certificates)

## Installation

1. **Clone the repository:**

   ```bash
   git clone https://github.com/rohanyh101/BreadcrumbsGo-ChatApp-With-GorillaWebSockets.git
   cd realtime-chat-app
   ```

2. Install dependencies:

```bash
go get github.com/gorilla/websocket
```

3. Generate SSL certificates using,
```bash
bash keygen.sh
```

## Running the Application
1. Start the server:

```bash
go run main.go
```

2. Access the application:
 Open your web browser and navigate to https://localhost:8080.

## Usage
### Login:

 - Navigate to the login page and enter your credentials.
 - You will receive an OTP to verify your identity.

### Chat:

 - Once logged in, you can start sending and receiving messages.
 - Use the room change functionality to switch between different chat rooms.

## Contributing
Contributions are welcome! Please fork the repository and create a pull request with your changes. Make sure to update the documentation as necessary.

<!-- ## License
This project is licensed under the MIT License. See the LICENSE file for details.
-->

## Acknowledgements
 - Gorilla WebSocket - A fast, well-tested, and widely used WebSocket implementation for Go.
 - OpenSSL - A robust, full-featured open-source toolkit implementing the Secure Sockets Layer (SSL) and Transport Layer Security (TLS) protocols.

## Contact
For any questions or inquiries, don't hesitate to get in touch with rohanyh101@gmail.com.
