CREATE TABLE Product (
    Id INT IDENTITY(1,1) PRIMARY KEY,
    Name VARCHAR(150) NOT NULL,
    Product_Category VARCHAR(50) CONSTRAINT CK_Product_Category CHECK (Product_Category IN ('Rokok', 'Obat', 'Lainnya')),
    Description VARCHAR(255)
);

CREATE TABLE Price (
    Id INT IDENTITY(1,1) PRIMARY KEY,
    Product_Id INT NOT NULL,
    Unit VARCHAR(100),
    FOREIGN KEY (Product_Id) REFERENCES Product(Id)
);

CREATE TABLE PriceDetail (
    Id INT IDENTITY(1,1) PRIMARY KEY,
    Price_Id INT NOT NULL,
    Tier VARCHAR(50) CONSTRAINT CK_Tier CHECK (Tier IN ('Non Member', 'Basic', 'Premium')),
    Price INT,
    FOREIGN KEY (Price_Id) REFERENCES Price(Id)
);