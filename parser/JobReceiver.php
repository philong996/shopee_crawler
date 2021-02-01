<?php
include("Parser.php");

require_once  '../vendor/autoload.php'; //__DIR__ .
use PhpAmqpLib\Connection\AMQPStreamConnection;


$connection = new AMQPStreamConnection('192.168.4.201', 5672, 'dmx_test', 'dmx_test');
$channel = $connection->channel();

$channel->queue_declare('anker_1_download_results', false, false, false, false);

echo " [*] Waiting for messages. To exit press CTRL+C\n";

$callback = function ($msg) {
  echo ' [V] Received ', $msg->body, "\n";

  // set a variable for api
  
  $servername = "192.168.4.200";
  $username = "new_engineer";
  $password = "New@Team";
  $dbname = "test";

  $conn = CreateConnection($servername, $username, $password, $dbname);

  $products = parseData($msg->body, $conn);

  echo "Created date is " . date("Y-m-d h:i:sa"), "\n";
  echo "number of products: ", count($products), "\n";
};

$channel->basic_consume('anker_1_download_results', 'dmx_test_exchange', false, true, false, false, $callback);

while ($channel->is_consuming()) {
    $channel->wait();
}
$conn->close();

?>