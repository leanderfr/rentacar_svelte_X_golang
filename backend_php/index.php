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

//********************************************************************
// analisa rota recebida e decide o que fazer
//********************************************************************
$path = parse_url($_SERVER["REQUEST_URI"], PHP_URL_PATH);

// prepara handlers (controllers)
$handlerWorkgroups = new Workgroups;

$router = new Router;

$getRequest = $_SERVER['REQUEST_METHOD']==='GET';
$postRequest = $_SERVER['REQUEST_METHOD']==='POST';

// PHP < 8.4 cannot deal with PATCH or PUT methods unfortunately, the test below will never be true
$patchRequest = $_SERVER['REQUEST_METHOD']==='PATCH';

$deleteRequest = $_SERVER['REQUEST_METHOD']==='DELETE';

// all returns will be json format
header('Content-Type: application/json');

die('jklljk');


//*********************************************************************************************************************************************************
// get method
//*********************************************************************************************************************************************************
if ($getRequest) {


  // http://localhost:8082/course  (GET)
  // http://localhost:8082/course?themes[]=inovacao&themes[]=vendas
  $router->Get("/course", function() use($handlerCourses) {  
    $searchbox = $_GET['search'] ?? '';
    $active = $_GET['active'] ?? '';
    $themes = $_GET['themes'] ?? [];

    $data = $handlerCourses->getCourses($active, $searchbox, $themes);

    // if there's an error, the getCourses function already returned the 500 error code along with the error message
    die( json_encode($data) );  
  });

  // http://localhost:8082/course/6853faed73c0f  (GET)
  $router->Get("/course/{id}", function($id) use($handlerCourses) {  
    $data = $handlerCourses->getCourseById($id);

    // if there's an error, the getCourseById function already returned the 500 error code along with the error message
    die( json_encode($data) );  
  });

  // http://localhost:8082/session/783933dd73c0f  (GET)
  $router->Get("/session/{id}", function($id) use($handlerSessions) {  
    $data = $handlerSessions->getSessionById($id);

    die( json_encode($data) );  
  });

  // http://localhost:8082/session/course/7638943h73c0f  (GET)
  $router->Get("/session/course/{id}", function($id) use($handlerSessions) {  
    $data = $handlerSessions->getSessionsByCourseId($id);

    die( json_encode($data) );  
  });

  // http://localhost:8082/enrollment/session/4309843ggbc0f  (GET)
  $router->Get("/enrollment/session/{id}", function($id) use($handlerEnrollments) {  
    $data = $handlerEnrollments->getEnrollmentsBySessionId($id);

    die( json_encode($data) );  
  });

  // http://localhost:8082/user/email  (GET)
  // the route /user/email is needed to keep the default created by Dot Digital backend
  $router->Get("/user/email", function() use($handlerUsers) {  
    $searchbox = $_GET['search'] ?? '';
    $active = $_GET['active'] ?? '';

    $data = $handlerUsers->getUsers($active, $searchbox);

    // if there's an error, the getUsers function already returned the 500 error code along with the error message
    die( json_encode($data) );  
  });


  // http://localhost:8082/user  (GET)
  $router->Get("/user", function() use($handlerUsers) {  
    $searchbox = $_GET['search'] ?? '';
    $active = $_GET['active'] ?? '';

    $data = $handlerUsers->getUsers($active, $searchbox);

    // if there's an error, the getUsers function already returned the 500 error code along with the error message
    die( json_encode($data) );  
  });


  // http://localhost:8082/user/7638943h73c0f  (GET)
  $router->Get("/user/{id}", function($id) use($handlerUsers) {  
    $data = $handlerUsers->getUserById($id);

    die( json_encode($data) );  
  });



}

//*********************************************************************************************************************************************************
// post method 
//*********************************************************************************************************************************************************
if ($postRequest) {

  // http://localhost:8082/workgrop/generate/smart (POST)
  $router->Post("/workgroup/{action}/{workgroup}", function($action, $workgroup) use($handlerWorkgroups)  {  
    $handlerWorkgroups->GetNewWorkgroupOrResetCurrentOrChooseRandomlyAnother($action, $workgroup);
  });

  // http://localhost:8082/course (POST)
  $router->Post("/course", function() use($handlerCourses)  {  
    $handlerCourses->updateOrInsertCourse();
  });

  //http://localhost:8080/course/6853fb18d3ac3/activity  (POST)
  $router->Post("/course/{id}/activity", function($id) use($handlerCourses)  {  
    $handlerCourses->ChangeStatus($id);
  });

  //http://localhost:8080/session/6853789443433/update  (POST)
  $router->Post("/session/{id}/update", function($id) use($handlerSessions)  {  
    $handlerSessions->updateOrInsertSession($id);
  });

  //http://localhost:8080/session  (POST)
  $router->Post("/session", function() use($handlerSessions)  {  
    $handlerSessions->updateOrInsertSession();
  });

  //http://localhost:8080/session/6853789443433/status/expired  (POST)
  //http://localhost:8080/session/6853789443433/status/available  (POST)
  $router->Post("/session/{id}/status/{status}", function($id, $status) use($handlerSessions)  {  
    $handlerSessions->ChangeAvailability($id, $status);
  });

  //http://localhost:8080/enrollment  (POST)
  $router->Post("/enrollment", function() use($handlerEnrollments)  {  
    $handlerEnrollments->Enroll();
  });

  //http://localhost:8080/user/6853fb18d3ac3/activity  (POST)
  $router->Post("/user/{id}/activity", function($id) use($handlerUsers)  {  
    $handlerUsers->ChangeStatus($id);
  });


  // http://localhost:8082/user/6853faed73c0f/update (POST)
  $router->Post("/user/{id}/update", function($id) use($handlerUsers)  {  
    $handlerUsers->updateOrInsertUser($id);
  });

  // http://localhost:8082/course (POST)
  $router->Post("/user", function() use($handlerUsers)  {  
    $handlerUsers->updateOrInsertUser();
  });




}

//*********************************************************************************************************************************************************
// delete method
//*********************************************************************************************************************************************************
if ($deleteRequest) {


  // http://localhost:8082/course/6853faed73c0f (DELETE)
  $router->Delete("/course/{id}", function($id) use($handlerCourses)  {  
    $handlerCourses->deleteCourse($id);
  });

  // http://localhost:8082/session/6853faed73c0f (DELETE)
  $router->Delete("/session/{id}", function($id) use($handlerSessions)  {  
    $handlerSessions->deleteSession($id);
  });

  // http://localhost:8082/enrollment/6852378923c0f (DELETE)
  $router->Delete("/enrollment/{id}", function($id) use($handlerEnrollments)  {  
    $handlerEnrollments->deleteEnrollment($id);
  });


  // http://localhost:8082/user/6852378923c0f (DELETE)
  $router->Delete("/user/{id}", function($id) use($handlerUsers)  {  
    $handlerUsers->deleteUser($id);
  });



}




$router->dispatch($path);

//********************************************************************
//********************************************************************


?>