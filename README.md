Product Management System
This is a Product Management System built using Go, Gin framework, MySQL, and Redis. It provides CRUD (Create, Read, Update, Delete) functionality for managing products and includes caching for improved performance.

Features
1. CRUD Operations: Create, retrieve, update, and delete product details.
2. Caching: Utilizes Redis for caching product data to improve read performance.
3. Database Integration: MySQL database for persistent data storage.
4. Logging: Logs request details including method, path, status, and response time.
5. Unit Testing: Includes basic unit and benchmark tests for CRUD endpoints.

Prerequisites
Before running the application, ensure the following are installed on your system:
1. Go 
2. MySQL 
3. Redis

API Endpoints
Method
1. POST	-> Create a new product
2. GET -> Retrieve a product by ID or all products
3. PUT	-> Update a product
4. DELETE	->	Delete a product by ID

Create Product
Sample:
{
  "user_id": 1,
  "product_name": "Test Product",
  "product_description": "Sample Description",
  "product_price": "Decimal"
}

Architectural Choices
1. Gin Framework
  Chosen for its lightweight nature and fast performance.
  Ideal for building RESTful APIs.
2. MySQL
  Used for persistent data storage.
  Relational model fits well with structured data.
3. Redis
  Used for caching product data to reduce database calls and improve response times.
4. Logrus
  Used for structured logging, providing insights into API usage and debugging.
5. Testing
  Mocking is used to simulate database and Redis calls.
  Unit tests and benchmarks validate and measure performance of key functionalities.

Assumptions
1. Products have a fixed schema with fields like ProductName, ProductDescription, ProductPrice, and UserID.
2. Redis is used primarily as a cache layer and not as the main datastore.
3. MySQL is running locally with the provided credentials.


