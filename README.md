# Real-Time Forum

A single-page application forum with real-time features using WebSockets. This project implements a complete forum system with user registration, post creation, commenting, and a real-time private messaging system with typing indicators.

## Technologies
- **Backend**: Go (Golang) with Gorilla WebSockets
- **Database**: SQLite
- **Frontend**: Vanilla JavaScript, HTML, CSS (Single Page Application)

## Core Features
- User authentication (registration and login)
- Post creation with categories
- Post commenting
- Real-time private messaging system
- "Typing in progress" indicators
- Online/offline user status

## Development Guidelines
- No frontend frameworks or libraries (pure JavaScript)
- Limited Go packages (standard packages, Gorilla WebSocket, SQLite3, bcrypt, and UUID packages)
- Single HTML file with JavaScript-based routing

## Dids
- Login and register pages front, back (can register and log)
- Home Page (list of users and posts displayed)
- All CRUD functions working

## DOs
- Separate connected and disconnected users (and display them separatly in the list of users)
- Implement cookie session system
- Add Comments to post
- Enhance to header (Name of the user, notif, disconnect button)
