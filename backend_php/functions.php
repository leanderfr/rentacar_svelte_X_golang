<?php


session_start();

//*********************************************************************************************************
// executa query SELECT e retorna registros inicialmente em formato array, que posteriormente
// pode ser convertido para json
//*********************************************************************************************************

function executeFetchQueryAndReturnJsonResult($sql): array {

  global $dbConnection;

  try {
    $result = mysqli_query($dbConnection, $sql) or internalError('[1] Database error / Erro na base de dados');

  } catch(Exception $e)  {
    internalError( mysqli_error($dbConnection) );
  }

  $anyData = mysqli_num_rows($result) > 0;

  if ( mysqli_num_rows($result) > 1) {
      /*********************************************************************************************************************
      exemplo json com mais de 1 item (array)
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

      // converte resultset em json 
      while($row =mysqli_fetch_assoc($result))    {
        $json[] = $row;
      }
  }

  if ( mysqli_num_rows($result) == 1) {

    /*********************************************************************************************************************
    exemplo json com 1 item (json puro, nao array)
    {
        "name": "nome",
        "manufacturer_id": "173",
        "country": "usa",
        "manufacturer_name": "Wolkswagen",
        "rental_price": "550,00",
    }
    *********************************************************************************************************************/
    $json = mysqli_fetch_assoc($result);
  }

  return($json);

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
      //logNotification($dbOperation, $tablename, $recordId, $workgroup, $clientIp);
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

?>