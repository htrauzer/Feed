# 🏛️ Go Feed Project
A full-stack web forum built with **Go (Golang)**, featuring real-time data persistence with **SQLite** and containerized deployment via **Docker**. This project allows users to share posts, categorize discussions, and interact through likes and comments.

## 🚀 Features
* User Authentication:  Secure registration and login using bcrypt password hashing.
* Persistent Storage:   Data is stored in an SQLite database that survives container restarts.
* Interactive Content:  Users can create posts, leave comments, and react (Like/Dislike).
* Clean Architecture:   Separation of concerns between Handlers, Database logic, and Utils.
* Containerized:        Fully Dockerized for easy setup and deployment.

## 🛠️ Tech Stack
* Backend: Go (Golang)
* Database: SQLite3
* Frontend: HTML5, CSS3
* Security: Bcrypt (password encryption)
* Conteiner: Docker

## 📦 Getting Started
### Prerequisites
$$
Docker installed on your machine.
$$

### Installation & Running
    
1. Clone the repository:
2. Build the Docker image: ```bash build.sh```
3. Access the Forum: open your browser and navigate to http://localhost:8080.

## 📂 Project Structure
   
* main.go            - Entry point and server initialization.
* handlers.go        - Route handlers and HTTP request logic.
* database/          - SQL queries and database connection management.
* models/            - Struct definitions for Users, Posts, and Comments.
* utils.go           - Helper functions (Session checks, validations).
* static/            - Frontend assets:
*    ├── css/        - Stylesheets for layout and design.
*    ├── js/         - Client-side logic and interactivity.
*    └── templates/  - HTML files   
     
