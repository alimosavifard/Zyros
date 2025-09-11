# Zyros CMS

A content management system with a Go backend (Gin, GORM, PostgreSQL, Redis) and a React + Vite frontend.

## Repository Structure

- `backend/`: Go-based backend
- `frontend/`: React + Vite frontend

## Prerequisites

- Docker and Docker Compose
- Go 1.21+ (optional for non-Docker setup)
- Node.js 18+ (optional for non-Docker setup)
- Nginx (for production deployment)

## Installation (Single Instance)

1. Clone the repository:

   <code>
   git clone https://github.com/alimosavifard/Zyros.git
   cd Zyros
   </code>


Set up environment files:
cp backend/.env.example backend/.env
cp frontend/.env.example frontend/.env

Edit backend/.env and frontend/.env with your settings (e.g., VITE_API_URL, DB_NAME, MODULE_NAME).

Run with Docker Compose:
docker-compose up


Access:

Backend health check: http://localhost:8080/api/v1/health
Frontend: http://localhost:3000
Swagger docs: http://localhost:8080/swagger/index.html



Running Multiple Instances
To run multiple instances (e.g., for different CMS sites):

Clone into separate directories:
</code>
git clone https://github.com/alimosavifard/Zyros.git zyros-instance1
git clone https://github.com/alimosavifard/Zyros.git zyros-instance2
</code>

For each instance:

Edit backend/.env with unique settings (e.g., MODULE_NAME=zyros-backend-instance1, DB_NAME=zyros_db_instance1, SERVER_PORT=:8081, REDIS_DB=1).
Edit frontend/.env with the corresponding backend URL (e.g., VITE_API_URL=http://localhost:8081/api/v1, VITE_FRONTEND_PORT=3001).
Update docker-compose.yml to use unique ports (e.g., 8081:8080, 3001:3000).


Run each instance:
<code>
cd zyros-instance1
docker-compose up
</code>


Nginx Configuration
For production, configure Nginx to serve the frontend static files and proxy API requests:
<code>
server {
    listen 80;
    server_name domain.com;

    root /var/www/domain.com/frontend/public;
    index index.html;

    location / {
        try_files $uri $uri/ /index.html;
    }

    location /api/ {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
</code>
Environment Variables
See backend/.env.example and frontend/.env.example for required variables.
Notes

Use unique MODULE_NAME, DB_NAME, and REDIS_DB for each instance to avoid conflicts.
For HTTPS, configure SSL certificates in Nginx.
For better performance, ensure build.rollupOptions.manualChunks is configured in vite.config.js.


