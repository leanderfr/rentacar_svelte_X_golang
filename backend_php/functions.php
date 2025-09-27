<?php

use Aws\S3\S3Client;
use Aws\S3\Exception\S3Exception as S3;

$THEMES = [
    'inovacao' => 'Inovação',
    'tecnologia' => 'Tecnologia',
    'gestao' => 'Gestão',
    'lideranca' => 'Liderança',
    'marketing' => 'Marketing',
    'vendas' => 'Vendas',
    'financas' => 'Finanças',
    'comunicacao' => 'Comunicação',
    'negocios' => 'Negócios',
    'empreendedorismo' => 'Empreendedorismo',
];


//*********************************************************************************
// informs the received route is incorrect
//*********************************************************************************
function routeError($detail='') {
  http_response_code(500);     
  die('Error with the route= '.$detail);
}


/*********************************************************************************************************
 internal server error
*********************************************************************************************************/

function internalError($message = 'Internal Error') {
  //http_response_code(500);   
  header("HTTP/1.1 500");
  die( json_encode(array('errors'=>$message)) );
}


/*********************************************************************************************************
 fetch data and convert records into json 
  $noArrayNeeded = true, returning one record of course,   false, returning multiple records inside an array
*********************************************************************************************************/

function executeFetchQueryAndReturnJsonResult($sql, $noArrayNeeded=false): array {

  global $dbConnection;
   
  try {
    $stmt = $dbConnection->query($sql) or internalError('[1] Database error / Erro na base de dados');    

  } catch(Exception $e)  {
    internalError( $dbConnection->errorInfo() );
  }

  if ($noArrayNeeded) {
    $row = $stmt->fetch(PDO::FETCH_ASSOC);
    return($row);
  }

  if (! $noArrayNeeded) {
    $json = array();

    while ($row = $stmt->fetch(PDO::FETCH_ASSOC)) {
      $json[] = $row;
    }
  }

  //********************************************************************************************
  // send data to requesting function
  //********************************************************************************************

  return( $json );     
}

/*********************************************************************************************************
 converts 'themes' field (courses table) to an array

 themes string example:  ["tecnologia", "educacao"]  will be converted to an array like this:

  "themes": [
      {
          "id": "comunicacao",
          "name": "Comunica\u00e7\u00e3o"
      },
      {
          "id": "tecnologia",
          "name": "Tecnologia"
      }
  ],

*********************************************************************************************************/
function themesstr_to_themesarray($_val): array  {

global $THEMES;

$theme_string = $_val;
$theme_string = str_replace('[','',$theme_string);
$theme_string = str_replace(']','',$theme_string);
$theme_string = str_replace('"','',$theme_string);

$themes_array = explode(',', $theme_string);
$themes_result = array();
for ($e=0; $e<count($themes_array); $e++)  {
  $theme_id = trim($themes_array[$e]);
  array_push($themes_result, 
      array('id'=> $theme_id, 'name'=>$THEMES[$theme_id]) );
}

return($themes_result);
}

//*********************************************************************************************************
// execute query to post or patch  (crud)
//*********************************************************************************************************

function executeCrudQueryAndReturnResult($sql) {
  global $dbConnection;

  try {
    $stmt = $dbConnection->query($sql) or internalError('[1] Database error / Erro na base de dados');    

  } catch(Exception $e)  {
    internalError( $dbConnection->errorInfo() );
  }

}


/*********************************************************************************************************
 converts 'themes' array  to an string to be recorded in the 'themes' field (courses table)
*********************************************************************************************************/
function themesarray_to_themesstr($themes): string  {

$_themesstr='';
for ($x=0; $x<count($themes); $x++)   {
  $_themesstr .= ($_themesstr == '') ? '' : ',';
  $_themesstr .= '"'.$themes[$x].'"' ;
}
$_themesstr = "[$_themesstr]";


return($_themesstr);

}

/*********************************************************************************************************
 tests date 
*********************************************************************************************************/
function validateDate($date, $format = 'Y-m-d') {
    $d = DateTime::createFromFormat($format, $date);

    return $d && $d->format($format) === $date;
}





?>
