<?php
include("UtilsDB.php");

$servername = "192.168.4.200";
$username = "new_engineer";
$password = "New@Team";
$dbname = "test";

// sql to create table
$sql = "CREATE TABLE shopee_product (
    id INT(6) UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(250) NOT NULL,
    categories JSON,
    url VARCHAR(100) NOT NULL,
    rrp_price FLOAT,
    stock INT,
    sale_price INT,	
    rating FLOAT,
    sold INT,
    view_count INT,
    liked_count INT,
    discount INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
    )";

// create data table
$conn = CreateConnection($servername, $username, $password, $dbname);
ExecuteQuery($sql, $conn);


?>