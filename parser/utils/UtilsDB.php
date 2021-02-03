<?php

function CreateConnection($servername, $username, $password, $dbname) {
    // Create connection
    $conn = new mysqli($servername, $username, $password, $dbname);

    // Check connection
    if ($conn->connect_error) {
    die("Connection failed: " . $conn->connect_error);
    }

    /* change character set to utf8 */
    if (!$conn->set_charset("utf8")) {
        printf("Error loading character set utf8: %s\n", $conn->error);
    } else {
        printf("Current character set: %s\n", $conn->character_set_name());
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