<?php
include("./utils/UtilsDB.php");

function parseData($msg, $conn) {
    
    $output = utf8_decode($msg);
    $data = json_decode($output, true);
    
    $items = $data["items"];
    
    $result = array();
    foreach ($items as $item) {
        
        $array_item = array();

        $array_item["name"] = $item["name"];
        $array_item["sale_price"] = $item["price"] / 100000;
        $array_item["rrp_price"] =  $item["price_before_discount"] / 100000;
        $array_item["stock"] = $item["stock"];
        $array_item["url"] = "https://shopee.vn/" . "P" . "-i." . $item["shopid"] . "." . $item["itemid"];
        
        // create query
        $sql = "INSERT INTO shopee_product (name, url, rrp_price, sale_price, stock) VALUES ('" . $array_item["name"]. "', '" . $array_item["url"] . "', " . $array_item["rrp_price"] . ", " . $array_item["sale_price"] . ", " . $array_item["stock"] . ");" ; 

        // insert data to database
        ExecuteQuery($sql, $conn);

        array_push($result, $array_item);
    };

    return $result;
}


?>