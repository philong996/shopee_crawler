<?php

function CreateConnection($servername, $username, $password, $dbname) {
    // Create connection
    $conn = new mysqli($servername, $username, $password, $dbname);

    // Check connection
    if ($conn->connect_error) {
    die("Connection failed: " . $conn->connect_error);
    }

    return $conn;
}

function ExecuteQuery($query, $conn) {

    if ($conn->query($query) === TRUE) {
        echo "Query executed successfully", "\n";
        } else {
        echo "Error executing query: " . $conn->error;
        }
}


?>