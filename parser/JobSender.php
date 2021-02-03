<?php


require_once  '../vendor/autoload.php'; //__DIR__ .
use PhpAmqpLib\Connection\AMQPStreamConnection;
use PhpAmqpLib\Message\AMQPMessage;

$connection = new AMQPStreamConnection('192.168.4.201', 5672, 'dmx_test', 'dmx_test');
$channel = $connection->channel();

$channel->queue_declare('anker_1', false, true, false, false);


$job = array("Url" => "https://shopee.vn/api/v2/search_items/?by=pop&limit=100&match_id=16461019&newest=0&order=desc&page_type=shop&version=2", "Interval" => 600);

echo "message is ", $job["Url"], "\n";
$msg = new AMQPMessage(json_encode($job));
$channel->basic_publish($msg, 'dmx_test_exchange', 'anker_1_key');


echo " [x] Sent message to anker 1\n";

?>