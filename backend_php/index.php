<?php
// exige que toda variavel seja tipada
declare(strict_types=1);

header("Content-Type: text/html; charset=utf-8"); 

// headers que permitem front end nao hospedado no mesmo servidor chamar o backend atual
//  CORS = Cross-Origin Resource Sharing
header('Access-Control-Allow-Origin: *');
header("Access-Control-Allow-Headers: X-API-KEY, Origin, X-Requested-With, Content-Type, Accept, Access-Control-Request-Method,Access-Control-Request-Headers, Authorization");
header('Access-Control-Allow-Methods: POST, GET, DELETE, PUT, PATCH, OPTIONS');       

$method = $_SERVER['REQUEST_METHOD'];

// OPTIONS= browser enviou só sinal de teste, cai fora
if ($method == "OPTIONS") {
    header('Access-Control-Allow-Origin: *');
    header("Access-Control-Allow-Headers: X-API-KEY, Origin, X-Requested-With, Content-Type, Accept, Access-Control-Request-Method,Access-Control-Request-Headers, Authorization");
    header("HTTP/1.1 200 OK");
    die();
} 
header("HTTP/1.0 200 OK"); 

require 'setup.php';
require 'functions.php';
require "Router.php";

require "handlers/Workgroups.php";  

require "models/Workgroups.php";  

//********************************************************************
// analisa rota recebida e decide o que fazer
//********************************************************************
$path = parse_url($_SERVER["REQUEST_URI"], PHP_URL_PATH);

// prepara handlers (controllers)
$WorkgroupsHandler = new Workgroups;

$router = new Router;

$getRequest = $_SERVER['REQUEST_METHOD']==='GET';
$postRequest = $_SERVER['REQUEST_METHOD']==='POST';

// PHP < 8.4 cannot deal with PATCH or PUT methods unfortunately, the test below will never be true
$patchRequest = $_SERVER['REQUEST_METHOD']==='PATCH';

$deleteRequest = $_SERVER['REQUEST_METHOD']==='DELETE';

// all returns will be json format
header('Content-Type: application/json');


//*********************************************************************************************************************************************************
// get method
//*********************************************************************************************************************************************************
if ($getRequest) {




}

//*********************************************************************************************************************************************************
// post method 
//*********************************************************************************************************************************************************
if ($postRequest) {

  // http://localhost:8082/workgrop/generate/smart (POST)
  $router->Post("/workgroup/{action}/{workgroup}", function($action, $workgroup) use($WorkgroupsHandler)  {  
    $WorkgroupsHandler->GetNewWorkgroupOrResetCurrentOrChooseRandomlyAnother($action, $workgroup);
  });






}

//*********************************************************************************************************************************************************
// delete method
//*********************************************************************************************************************************************************
if ($deleteRequest) {




}




$router->dispatch($path);

//********************************************************************
//********************************************************************


?>