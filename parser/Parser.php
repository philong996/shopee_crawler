<?php
include("./utils/UtilsDB.php");

require_once  '../vendor/autoload.php'; //__DIR__ .
use PhpAmqpLib\Connection\AMQPStreamConnection;
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

// set connection to the rabbitmq for data
$connection = new AMQPStreamConnection('192.168.4.201', 5672, 'dmx_test', 'dmx_test');
$channel = $connection->channel();
$channel->queue_declare('anker_1_download_results', false, false, false, false);

echo " [*] Waiting for messages. To exit press CTRL+C\n";

// set connection to database
$servername = "192.168.4.200";
$username = "new_engineer";
$password = "New@Team";
$dbname = "test";
$conn = CreateConnection($servername, $username, $password, $dbname);

$callback = function ($msg) {
    echo ' [V] Received ', "\n"; //$msg->body,

    global $conn;

    $product = parseData($msg->body, $conn);

    echo "Created date is " . date("Y-m-d h:i:sa"), "\n";
};

$channel->basic_consume('anker_1_download_results', 'dmx_test_exchange', false, true, false, false, $callback);

while ($channel->is_consuming()) {
    $channel->wait();
}
$conn->close();

?>