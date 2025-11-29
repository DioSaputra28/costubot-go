CREATE TABLE payments (
    payment_id INT PRIMARY KEY AUTO_INCREMENT,
    product_id INT,
    FOREIGN KEY (product_id) REFERENCES products (product_id),
    order_id INT,
    FOREIGN KEY (order_id) REFERENCES orders (order_id),
    amount INT NOT NULL,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL,
    phone VARCHAR(15),
    method VARCHAR(50) NOT NULL,
    status VARCHAR(50) NOT NULL,
    external_id VARCHAR(100) NOT NULL,
    payment_url TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME DEFAULT NULL
);