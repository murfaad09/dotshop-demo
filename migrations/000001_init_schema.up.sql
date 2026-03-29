
-- Create Users table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create Sessions table
CREATE TABLE sessions (
    id SERIAL PRIMARY KEY,
    user_id INT,
    session_token VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES Users(id)
);

-- Create Roles table
CREATE TABLE roles (
    id SERIAL PRIMARY KEY,
    role_name VARCHAR(255) UNIQUE NOT NULL
);

-- Create Permissions table
CREATE TABLE permissions (
    id SERIAL PRIMARY KEY,
    role_id INT,
    permission_name VARCHAR(255) UNIQUE NOT NULL,
    description TEXT,
    FOREIGN KEY (role_id) REFERENCES Roles(id)
);

-- Create Logs table
CREATE TABLE logs (
    id SERIAL PRIMARY KEY,
    user_id INT,
    action TEXT,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    details TEXT,
    FOREIGN KEY (user_id) REFERENCES Users(id)
);

-- Create Audit_Trail table
CREATE TABLE audit_trail (
    id SERIAL PRIMARY KEY,
    user_id INT,
    action TEXT,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    details TEXT,
    table_affected VARCHAR(255),
    record_id INT,
    FOREIGN KEY (user_id) REFERENCES Users(id)
);

-- Create Vendor_Configurations table
CREATE TABLE vendor_configurations (
    id SERIAL PRIMARY KEY,
    vendor_name VARCHAR(255) UNIQUE NOT NULL,
    api_key TEXT,
    endpoint TEXT,
    status VARCHAR(50)
);

-- Create Creators table
CREATE TABLE creators (
    id SERIAL PRIMARY KEY,
    user_id INT UNIQUE,
    name VARCHAR(255),
    bio TEXT,
    social_media_links TEXT[],
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES Users(id)
);

-- Create Suppliers table
CREATE TABLE suppliers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    address TEXT,
    contact_info TEXT,
    email VARCHAR(255),
    phone VARCHAR(20),
    website VARCHAR(255)
);

-- Create Brands table
CREATE TABLE brands (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    description TEXT,
    origin VARCHAR(100),
    founder VARCHAR(255),
    established_date DATE,
    website VARCHAR(255)
);
-- Create Products table
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    brand_id INT,
    supplier_id INT,
    name VARCHAR(255),
    description TEXT,
    category VARCHAR(255),
    price DECIMAL(10, 2),
    quantity INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (brand_id) REFERENCES Brands(id),
    FOREIGN KEY (supplier_id) REFERENCES Suppliers(id)
);

-- Create Media table
CREATE TABLE media (
    id SERIAL PRIMARY KEY,
    creator_id INT,
    media_type VARCHAR(50),
    media_url TEXT,
    caption TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (creator_id) REFERENCES Creators(id)
);

-- Create Orders table
CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    user_id INT,
    total_amount DECIMAL(10, 2),
    status VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES Users(id)
);

-- Create Order_Items table
CREATE TABLE order_items (
    id SERIAL PRIMARY KEY,
    order_id INT,
    product_id INT,
    quantity INT,
    price DECIMAL(10, 2),
    FOREIGN KEY (order_id) REFERENCES Orders(id),
    FOREIGN KEY (product_id) REFERENCES Products(id)
);

-- Create Shipping_Info table
CREATE TABLE shipping_info (
    id SERIAL PRIMARY KEY,
    order_id INT,
    address TEXT,
    city VARCHAR(100),
    state VARCHAR(100),
    country VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (order_id) REFERENCES Orders(id)
);

-- Create Payments table
CREATE TABLE payments (
    id SERIAL PRIMARY KEY,
    order_id INT,
    amount DECIMAL(10, 2),
    status VARCHAR(50),
    payment_method VARCHAR(100),
    transaction_id VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (order_id) REFERENCES Orders(id)
);

-- Reviews and Ratings
CREATE TABLE reviews (
    id SERIAL PRIMARY KEY,
    user_id INT,
    product_id INT,
    rating INT,
    review_text TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES Users(id),
    FOREIGN KEY (product_id) REFERENCES Products(id)
);

-- Promotions and Discounts
CREATE TABLE promotions (
    id SERIAL PRIMARY KEY,
    promotion_id VARCHAR(50) UNIQUE,
    discount DECIMAL(5, 2),
    valid_from TIMESTAMP,
    valid_to TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Wishlists
CREATE TABLE wishlists (
    id SERIAL PRIMARY KEY,
    user_id INT,
    product_id INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES Users(id),
    FOREIGN KEY (product_id) REFERENCES Products(id)
);

-- Favorites
CREATE TABLE favorites (
    id SERIAL PRIMARY KEY,
    user_id INT,
    product_id INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES Users(id),
    FOREIGN KEY (product_id) REFERENCES Products(id)
);

-- Notifications and Alerts
CREATE TABLE notifications (
    id SERIAL PRIMARY KEY,
    user_id INT,
    message TEXT,
    is_read BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES Users(id)
);

-- Refunds and Returns
CREATE TABLE refunds (
    id SERIAL PRIMARY KEY,
    order_id INT,
    reason TEXT,
    status VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (order_id) REFERENCES Orders(id)
);

-- Support Tickets
CREATE TABLE support_tickets (
    id SERIAL PRIMARY KEY,
    user_id INT,
    subject VARCHAR(255),
    message TEXT,
    status VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES Users(id)
);

-- Gift Cards
CREATE TABLE gift_cards (
    id SERIAL PRIMARY KEY,
    gift_card_id VARCHAR(50) UNIQUE,
    balance DECIMAL(10, 2),
    is_redeemed BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);