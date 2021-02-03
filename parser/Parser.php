<?php
include("./utils/UtilsDB.php");

require_once __DIR__ . '/../vendor/autoload.php';

use Elasticsearch\ClientBuilder;

function parseData($msg, $conn) {
    
    $output = utf8_decode($msg);
    $data = json_decode($output, true);
    
    $item = $data["item"];

    $array_item = array();

    $array_item["name"] = $item["name"];
    $array_item["sale_price"] = $item["price"] / 100000;
    $array_item["rrp_price"] =  $item["price_before_discount"] / 100000;
    $array_item["stock"] = $item["stock"];
    $array_item["url"] = "https://shopee.vn/" . "P" . "-i." . $item["shopid"] . "." . $item["itemid"];
    $array_item["sold"] = $item["sold"];
    $array_item["rating"] = $item["item_rating"]["rating_star"];
    $array_item["discount"] = (int)str_replace("%", "",$item["discount"]);
    $array_item["view_count"] = $item["view_count"];
    $array_item["liked_count"] = $item["liked_count"];
    $array_item["categories"] = $item["categories"];

    // insert data to database
    $sql = "INSERT INTO shopee_product (name, url, rrp_price, sale_price, stock, sold, rating, discount, view_count, liked_count, categories) 
            VALUES ('" . $array_item["name"] . "', '"
                    .$array_item["url"]. "',"
                    .$array_item["sale_price"] . ","
                    .$array_item["rrp_price"] . ","
                    .$array_item["stock"] . ","
                    .$array_item["sold"] . ","
                    .$array_item["rating"] . ","
                    .$array_item["discount"] . ","
                    .$array_item["view_count"] . ","
                    .$array_item["liked_count"] . ",'"
                    .json_encode($array_item["categories"], JSON_UNESCAPED_UNICODE) .
            "');";
    ExecuteQuery($sql, $conn);


    // insert data to elasticsearch
    $client = ClientBuilder::create()->build();
    $params = [
        'index' => 'anker-shopee-1',
        'body'  => $array_item];
    $response = $client->index($params);
    print_r("add data to ES: " ,$response["result"]);
}


?>