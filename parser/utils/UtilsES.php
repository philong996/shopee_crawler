<?php

require_once __DIR__ . '/../../vendor/autoload.php';

use Elasticsearch\ClientBuilder;

$client = ClientBuilder::create()->build();

$params = [
    'index' => 'anker-shopee-1',
    'body'  => ['name' => 'cde',
                'price' => 400000,
                'stock' => 50]
];

$response = $client->index($params);
print_r($response);


?>