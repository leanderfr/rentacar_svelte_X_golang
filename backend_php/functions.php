<?php

require 'setup.php';

use Aws\S3\S3Client;

session_start();

//*********************************************************************************************************
// converte registros mysql em json

// $noArray= true ==> o front end esta preparado para receber, em alguns momentos, um objeto json sozinho, sem 
//                    estar contido dentro de um array, nesses caso, $noArray = true
// 
// $noArray= false ==> devolve um array de JSON's
//*********************************************************************************************************

function executeFetchQueryAndReturnJsonResult($sql, $noArray=false) {

  global $dbConnection;
   
  try {
    $result = mysqli_query($dbConnection, $sql) or internalError('[1] Database error / Erro na base de dados');    

  } catch(Exception $e)  {
    internalError( mysqli_error($dbConnection) );
  }

  $anyData = mysqli_num_rows($result) > 0;

  // ! $noArray = monta array de objetos JSON
  if (! $noArray)        {
      /*********************************************************************************************************************
      retorna para o frontend em formato array de JSON, exemplo:  
        [
            {
                "expression": "Seja bem vindo, Visitante!",
                "item": "welcome"
            },
            {
                "expression": "Carros disponiveis",
                "item": "available_cars"
            },
        ]
      *********************************************************************************************************************/
      $json = array();

      // converte resultset em json para o front end
      while($row =mysqli_fetch_assoc($result))    {
        $json[] = $row;
      }
  }

  // $noArray = monta 1 unico objeto JSON
  if ($noArray)  
    /*********************************************************************************************************************
    retorna para o frontend em formato 1 JSON , exemplo:  
    {
        "name": "nome",
        "manufacturer_id": "173",
        "country": "usa",
        "manufacturer_name": "Wolkswagen",
        "rental_price": "550,00",
    }
    *********************************************************************************************************************/
    $json = mysqli_fetch_array($result, MYSQLI_ASSOC);


  //********************************************************************************************
  // envia dados para front end
  //********************************************************************************************
  header('Content-type: application/json'); 
  http_response_code(200);   // 200= requisicao bem sucedida

  if ($anyData)     {
      // converte texto para html entities, caso contrario, dará erro no front end
      //array_walk_recursive($json, function(&$item, $key) {
        //  $item = htmlentities($item);
      //});

      die( json_encode($json) );     
  }
  else   
    die( json_encode([]) );
}



//*********************************************************************************************************
// execute query de edicao de dados (crud)
// $needToReturnId= true, necessario retornar para que a funcao anterior execute algo a mais, como por
// exemplo fazer o upload da imagem de um carro, ao retornar, informa ID do registro recem inserido,
// se for o caso

// quando a operacao = insert/update os dados para gravacao de notificacao (workgroup, clientip, etc)
// sao enviados 
// quando a operacao = mudanca status, exclusao, nao é feita notificacao
//*********************************************************************************************************

function executeCrudQueryAndReturnResult($sql, $needToReturnId = false, $anyNotificationToMake='', $workgroup='', $clientIp='' ) {
  global $dbConnection;

  // executa a operacao (insert/update/delete)
  try {
    mysqli_query($dbConnection, $sql) or internalError('[2] Database error / Erro na base de dados');  

    // obtem ID do registro que foi inserido
    // em caso de 'update', $dbConnection->insert_id = 0
    //$newRecordId = $dbConnection->insert_id;    

    $lastId = mysqli_query($dbConnection, "select LAST_INSERT_ID() as record_id" ) or internalError('[3] Database error / Erro na base de dados');    
    if ($___lastID = mysqli_fetch_object($lastId))   $newRecordId = $___lastID->record_id;
    else internalError('[4] Database error / Erro na base de dados');    


    // cria notificacao para que os demais usuarios do grupo saibam que houve alteracao na base
    if ( $anyNotificationToMake != '' ) {
      $tmp = explode('|', $anyNotificationToMake) ;   

      $dbOperation = $tmp[0];    // insert ou update
      $tablename = $tmp[1];    // nome da tabela que foi manipulada

      if ($dbOperation=='update') $recordId = $tmp[2];
      else  $recordId = $newRecordId;
      logNotification($dbOperation, $tablename, $recordId, $workgroup, $clientIp);
    }

    http_response_code(200);   // 200= requisicao bem sucedida
    if ($needToReturnId) return "__success__|$newRecordId";
    else die( '__success__' );

  } catch(Exception $e)  {
    internalError( mysqli_error($dbConnection) );
  }

}




//*********************************************************************************************************
// retorna erro interno do servidor
//*********************************************************************************************************

function internalError($message = 'Internal Error') {
  http_response_code(500);   // 500= erro interno
  die( $message );
}

//*********************************************************************************
// informa que a rota passada nao existe
//*********************************************************************************
function routeError() {
  http_response_code(500);     // 500= erro interno
  die('Missing route');
}

//*********************************************************************************
// registra na base de dados o status atual da clonagem 
//*********************************************************************************
function logCloningStatus() {
global $dbConnection, $cloningPercentReady, $cloningStatus;

mysqli_query($dbConnection, "update config set cloning_status='$cloningStatus', cloning_percent_ready='$cloningPercentReady'") or internalError('[5] Database error / Erro na base de dados');  
}


//***********************************************************************
// grava arquivo de imagem (carro ou logotipo de fabricante) no 
// respositorio da AWS S3
// $recordType pode ser:  car ou manufacturer
//***********************************************************************
function uploadImageToAWS_S3( $recordType, $fileName, $recordId, $workgroup ) {

global $AWS_S3_APIKEY, $AWS_S3_SECRETKEY, $AWS_S3_BUCKET, $AWS_S3_IMAGES_FOLDER;

// cria nome baseado no ID do registro 
$suffixWithID =  str_pad($recordId, 6, '0', STR_PAD_LEFT); 

// nome do arquivo para gravacao local e para gravacao na AWS S3
$uniqueFileName = $workgroup . '_' . $recordType . '_' .  $suffixWithID . '.png'; // workgroup_car_000005.png, workgroup_manufacturer_000053.png, etc
$localFile = "tmp/$uniqueFileName"; // gravacao local

if (! move_uploaded_file( $_FILES[$fileName]['tmp_name'], $localFile))  
  internalError( $localFile);
//  internalError( 'Image upload failed / Upload da imaddgem falhou'.$localFile);

// conecta AWS S3
$s3Client = new S3Client([
    'region' => 'sa-east-1',
    'version' => 'latest',
    'credentials' => [
        'key' => $AWS_S3_APIKEY,
        'secret' => $AWS_S3_SECRETKEY
    ]
]);

$result = $s3Client->putObject([
  'Bucket' => $AWS_S3_BUCKET,
  'Key'    => "$AWS_S3_IMAGES_FOLDER/".basename($localFile),
  'Body'   => fopen($localFile, 'r'),
]);
 
    

}


//***********************************************************************************************
// cria notificacao para que os demais usuarios do grupo saibam que houve alteracao na base
// demais usuarios= IP diferente do IP que gerou a notificacao
//***********************************************************************************************
function logNotification($dbOperation, $tablename, $recordId, $workgroup, $clientIp) {

global $dbConnection, $imagesLocation, $staticImagesLocation;

//**********************************************************************
// prepara notificacao de que  algo foi alterado na tabela de carros
//**********************************************************************

if ($tablename=='cars') {
  // prepara dados dinamicos que serao colocados na notificacao de insercao/edicao de registro 
  if ($dbOperation == 'insert') 
    [$english, $portuguese] = GetTermByItem("notification_created_car", $workgroup);

  if ($dbOperation == 'update') 
    [$english, $portuguese] = GetTermByItem("notification_updated_car", $workgroup);

  $carName = '';
  $sql =  "select name from cars where id=$recordId";
  $result = mysqli_query($dbConnection, $sql ) or internalError('[6] Database error / Erro na base de dados');    
  if ($row = mysqli_fetch_object($result))   $carName = $row->name;
  else internalError('[7] Database error / Erro na base de dados');    

  // monta nome da imagem do carro baseado em seu ID
  $suffixWithID =  str_pad($recordId, 6, '0', STR_PAD_LEFT); 
  $uniqueFileName = $workgroup . '_car_' .  $suffixWithID . '.png'; // workgroup_car_000005.png, etc

  $imageUrl = "$imagesLocation/$uniqueFileName";

  $english = str_replace('@img' , $imageUrl, $english);   // substitui a string '@img' pelo link real da imagem
  $portuguese = str_replace('@img' , $imageUrl, $portuguese);   

  $english = str_replace('@name' , $carName, $english);   // substitui a string '@name' pelo nome real do carro
  $portuguese = str_replace('@name' , $carName, $portuguese);   
}

//**********************************************************************
// // prepara  notificacao de que  algo foi alterado na tabela de fabricantes
//**********************************************************************
if ($tablename=='manufacturers') {
  // prepara dados dinamicos que serao colocados na notificacao de insercao/edicao de registro 
  if ($dbOperation == 'insert') 
    [$english, $portuguese] = GetTermByItem("notification_created_manufacturer", $workgroup);

  if ($dbOperation == 'update') 
    [$english, $portuguese] = GetTermByItem("notification_updated_manufacturer", $workgroup);

  $manufacturerName = '';
  $sql =  "select name from manufacturers where id=$recordId";
  $result = mysqli_query($dbConnection, $sql ) or internalError('[8] Database error / Erro na base de dados');    
  if ($row = mysqli_fetch_object($result))   $manufacturersName = $row->name;
  else internalError('[9] Database error / Erro na base de dados');    

  // monta nome da imagem do logotipo do fabricante
  $suffixWithID =  str_pad($recordId, 6, '0', STR_PAD_LEFT); 
  $uniqueFileName = $workgroup . '_manufacturer_' .  $suffixWithID . '.png'; // workgroup_manufacturer_000005.png, etc

  $imageUrl = "$imagesLocation/$uniqueFileName";

  $english = str_replace('@img' , $imageUrl, $english);   // substitui a string '@img' pelo link real do logotipo
  $portuguese = str_replace('@img' , $imageUrl, $portuguese);   

  $english = str_replace('@name' , $manufacturerName, $english);   // substitui a string '@name' pelo nome real do fabricante
  $portuguese = str_replace('@name' , $manufacturerName, $portuguese);   
}


//*******************************************************************************************
// // prepara notificacao de que  algo foi alterado na tabela de expressoes ing/port
//*******************************************************************************************
if ($tablename=='terms')  {
  // prepara dados dinamicos que serao colocados na notificacao de insercao/edicao de registro 
  if ($dbOperation == 'insert') 
    [$english, $portuguese] = GetTermByItem("notification_created_expression", $workgroup);

  if ($dbOperation == 'update') 
    [$english, $portuguese] = GetTermByItem("notification_updated_expression", $workgroup);

  $termItem = '';
  $sql =  "select item from terms where id=$recordId";
  $result = mysqli_query($dbConnection, $sql ) or internalError('[10] Database error / Erro na base de dados');    
  if ($row = mysqli_fetch_object($result))   $termItem = $row->item;
  else internalError('[11] Database error / Erro na base de dados');    

  $imageUrl = "$staticImagesLocation/language_notification.png";

  $english = str_replace('@img' , $imageUrl, $english);   // substitui a string '@img' pelo link real da imagem ilustrativa
  $portuguese = str_replace('@img' , $imageUrl, $portuguese);   

  $english = str_replace('@item' , $termItem, $english);   // substitui a string '@item' pelo item real que foi alterado/inserido
  $portuguese = str_replace('@item' , $termItem, $portuguese);   
}


//**********************************************************************
// prepara notificacao de que reserva de carro foi feita ou alterada
//**********************************************************************

if ($tablename=='bookings') {
  // prepara dados dinamicos que serao colocados na notificacao de insercao/edicao de registro de carro
  if ($dbOperation == 'insert') 
    [$english, $portuguese] = GetTermByItem("notification_created_booking", $workgroup);

  if ($dbOperation == 'update') 
    [$english, $portuguese] = GetTermByItem("notification_updated_booking", $workgroup);

  // obtem detalhes da reserva
  $carName = '';
  $sql =  "select country, date_format(pickup_datetime, '%m/%d/%Y') as date_usa, date_format(updated_at, '%d/%m/%Y') as date_brazil, car_id, ".
          "       date_format(pickup_datetime, '%l:%i - %p') as hour_usa, date_format(pickup_datetime, '%H:%i') as hour_brazil " .
          "from bookings ".
          "where id = $recordId";
  $result = mysqli_query($dbConnection, $sql ) or internalError('[12] Database error / Erro na base de dados');    

  // le a data/hora no formato do país em que foi gravada a reserva de carro
  if ($row = mysqli_fetch_object($result))   {
    $date_usa = $row->date_usa;
    $hour_usa = $row->hour_usa;
    $date_brazil = $row->date_brazil;
    $hour_brazil = $row->hour_brazil;

    $carId = $row->car_id;
  }
  else internalError('[13] Database error / Erro na base de dados');    

  // obtem detalhes do carro reservado
  $carName = '';
  $sql =  "select name from cars where id = $carId";
  $result = mysqli_query($dbConnection, $sql ) or internalError('[14] Database error / Erro na base de dados');    
  if ($row = mysqli_fetch_object($result))   $carName = $row->name;
  else internalError('[15] Database error / Erro na base de dados');    


  // monta nome da imagem do carro baseado em seu ID
  $suffixWithID =  str_pad($carId, 6, '0', STR_PAD_LEFT); 
  $uniqueFileName = $workgroup . '_car_' .  $suffixWithID . '.png'; // workgroup_car_000005.png, etc

  $imageUrl = "$imagesLocation/$uniqueFileName";

  $english = str_replace('@img' , $imageUrl, $english);   // substitui a string '@img' pelo link real da imagem
  $portuguese = str_replace('@img' , $imageUrl, $portuguese);   

  $english = str_replace('@date' , $date_usa, $english);   // coloca dinamicamente a data do inicio da reserva 
  $english = str_replace('@hour' , $hour_usa, $english);   // coloca dinamicamente a hora do inicio da reserva 

  $portuguese = str_replace('@date' , $date_brazil, $portuguese);   
  $portuguese = str_replace('@hour' , $hour_brazil, $portuguese);   

}


//**********************************************************************
// grava a notificacao montada acima
//**********************************************************************

$portuguese = str_replace("'" , "\'", $portuguese);   // escape  de aspas simples que pode haver na frase
$english = str_replace("'" , "\'", $english);

// grava a notification de alteracao na tabela 'cars'
$sqlNotification = "insert into notifications(workgroup, made_by_ip, description_english, description_portuguese) ".
                    "select '$workgroup', '$clientIp', '$english', '$portuguese' ";		

mysqli_query($dbConnection, $sqlNotification) or internalError('[16] Database error / Erro na base de dados');  

// contabilizar 1 alteracao a mais na base de dados, feita pelo grupo atual
// faz isso para que, caso o usuario queira resetar os dados do grupo, o front end permita
// caso nenhuma alteracao tenha sido feita pelo grupo, o front end avisa que nao é necessario resetar os dados
mysqli_query($dbConnection, "update workgroups set database_changes_amount = ifnull(database_changes_amount, 0) + 1 where name='$workgroup'") 
    or internalError('[17] Database error / Erro na base de dados');  

}


//***********************************************************************************************
// obtem termo em ingles e portugues de determinado item na tabela de expressoes ing/port
//***********************************************************************************************
function GetTermByItem($item, $workgroup) {

global $dbConnection;

$sql =  "select portuguese, english ".
        "from terms ".
        "where item = '$item' and deleted_at is null and ifnull(active, false)=true and workgroup='$workgroup' ";

$result = mysqli_query($dbConnection, $sql ) or internalError('[18] Database error / Erro na base de dados');    

if ($row = mysqli_fetch_object($result)) 
  return ([$row->english, $row->portuguese]);

else 
  return (['* error *', '* error *']);
}



?>