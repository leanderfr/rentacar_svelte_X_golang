<?php
header("Content-Type: text/html; charset=utf-8"); 

// headers que permitem que o frontend esteja hospedado em outro servidor
header('Access-Control-Allow-Origin: *');
header("Access-Control-Allow-Headers: X-API-KEY, Origin, X-Requested-With, Content-Type, Accept, Access-Control-Request-Method,Access-Control-Request-Headers, Authorization");
header('Access-Control-Allow-Methods: POST, GET, DELETE, PUT, PATCH, OPTIONS');       

$method = $_SERVER['REQUEST_METHOD'];
if ($method == "OPTIONS") {
    header('Access-Control-Allow-Origin: *');
    header("Access-Control-Allow-Headers: X-API-KEY, Origin, X-Requested-With, Content-Type, Accept, Access-Control-Request-Method,Access-Control-Request-Headers, Authorization");
    header("HTTP/1.1 200 OK");
    die();
} 
header("HTTP/1.0 200 OK"); 

require('setup.php');
require('functions.php');
require('aws/aws-autoloader.php');



/*

  ***** observacao sobre o metodo PATCH:  *****

    o PHP 7.4.8, usado no servidor atual, nao da suporte para a funcao 'request_parse_body()'
    essa funcao lê o body quando o metodo= PATCH , sem essa funcao é impossivel ler BODY quando o metodo usado é 
    diferente de 'POST'.

    infelizmente o front end vai ter que usar o metodo POST, ao inves de PATCH quando o backend escolhido= PHP , 
    nos demais backends continua sendo usado PATCH quando um registro é alterado

    '$_POST' nao funciona qdo metodo= PATCH ou PUT

    no codigo abaixo, para efeito de registro é usada a palavra PATCH, mas efetivamente sera esperado o metodo 'POST'

*/

// le a rota recebida

// a url deve ser similar a:
// https://.../terms/english
// https://.../car/3,   etc etc
$route =  strip_tags(filter_input(INPUT_GET, 'url', FILTER_DEFAULT)  ?? '');

$params = explode('/', strtolower($route)) ;   
$paramsOriginalCase = explode('/', $route);   

// se a query de busca de dados for montada, a rota esta ok, e a query esta pronta para execucao
$fetchSql = '';

// se a query de CRUD (insercao, edicao, exclusao) for montada, a rota esta ok, e a query esta pronta para execucao
$crudSql = '';

$method = $_SERVER['REQUEST_METHOD'];

// o front end esta preparado para receber, em alguns momentos, um objeto json sozinho, sem 
// estar contido dentro de um array, nesses caso, $noArray = true
$noArray = false;

// controla se deve ser feito registro de notificacao de alteracao de dados 
// em caso de exclusao de registro, mudanca de status (active/inactive) nao é feito log, somente em caso de update/insert
$anyNotificationToMake = '';

/*********************************************************************************
 rota: GET   /terms/:language
 rota: GET   /:workgroup/terms/:language

 equivalente backend golang= GetTermsForFillingUpFrontEnd

 lê expressoes (termos) do idioma informado (e/ou do grupo se foi informado) 
*********************************************************************************/
if ( $method=='GET' && 
    (isset($params[0]) && $params[0] == 'terms' && isset($params[1]) && ! isset($params[3])) ||
    (isset($params[1]) && $params[1] == 'terms' && isset($params[2]) && ! isset($params[3])) )     {

  $language = ''; 
  $workgroup = 'Admin';    // se nao foi passado workgroup, usa dados do grupo 'Admin'

  //  /terms/:language
  if ( ($params[0]=='terms') && ($params[1] == 'portuguese' || $params[1] == 'english') ) { 
    $language = $params[1];
    $routeOk = true;
  }    

  // :workgroup/terms/:language
  if ( ($params[1]=='terms') && ($params[2] == 'portuguese' || $params[2] == 'english') ) { 
    $language = $params[2]; 
    $workgroup = $params[0]; 
    $routeOk = true;
  }    


  if (!  $routeOk) routeError();

  $fetchSql =  "select $language as expression, item ".
                "from terms  ".
                "where workgroup = '$workgroup' and ifnull(active, false)= true and deleted_at is null ";
  
}

/*********************************************************************************
 rota: GET     /:tablename/record_metadata/:id/:country

 equivalente backend golang= GetRecordMetadata

 lê metadados de determinado registro 
*********************************************************************************/
if ( $method=='GET' && 
    isset($params[1]) && $params[1] == 'record_metadata' && isset($params[2]) && isset($params[3]) ) {

  $tablename = $params[0];   
  $country = $params[3];   
  $recordId =  $params[2]; 

  $routeOk = false;
	if ($country == "usa" && is_numeric($recordId)) {
		$dateFormat = "%M / %d / %Y - %h:%i %p";   // mm/dd/yyyy hh:mm am/pm
		$locale = "en_US";
    $routeOk = true;
	} 

  if ($country == "brazil" && is_numeric($recordId)) {
		$dateFormat = "%d / %M / %Y - %H:%i";   // dd/mm/yyyy HH:mm
		$locale = "pt_BR";
    $routeOk = true;
	}

  if (!  $routeOk) routeError();

	$fetchSql = "select date_format($tablename.created_at, '$dateFormat') as created_at, ".
              " date_format($tablename.updated_at, '$dateFormat') as updated_at, workgroup  ".
              "from $tablename ".
              "where id=$recordId ";

  $noArray = true;   // avisa para enviar JSON sozinho (nao array), o frontend aguarda como no exemplo abaixo:
    /*
      {
          "created_at": "26 / janeiro / 2025 - 21:41",
          "updated_at": "26 / janeiro / 2025 - 21:59",
          "workgroup": "Accomplished"
      }
    */
}



/*********************************************************************************
 rota: PATCH    /:tablename/status/:id

 equivalente backend golang= ChangeRecordStatus

 altera o status (active/inactive) de determinado registro
*********************************************************************************/
if ( $method=='PATCH' && 
  isset($params[1]) && $params[1] == 'status' && isset($params[2]) ) {

  $tablename = $params[0];   
  $recordId =  $params[2]; 

  $routeOk = false;
	if (is_numeric($recordId)) {
    $routeOk = true;
	} 

  if (!  $routeOk) routeError();

  $dbConnection -> autocommit(true);    // gravacao sem necessidade de confirmacao 
	$crudSql = "update $tablename set active = CASE WHEN ifnull(active, false)=false THEN true else false END, updated_at=now() where id=$recordId";
}


/*********************************************************************************
 rota: DELETE     /:tablename/delete/:id
	                /:tablename/batch_delete/:ids

 equivalente backend golang=  DeleteRecord
                              DeleteRecords

 exclui registro(s) de determinada tabela
*********************************************************************************/
if ( $method=='DELETE' && 
    isset($params[1]) && ($params[1] == 'delete' || $params[1] == 'batch_delete') && isset($params[2]) ) {

  $which = $params[1];

  $routeOk = false;

  $tablename = $params[0];   
  if ($which=='delete')   {
    $recordId =  $params[2]; 

    if (is_numeric($recordId)) {
      $routeOk = true;
    } 
  }

  // nao é possivel verificar se a lista de registros é toda numerica, se der erro, o mysql vai avisar e o front end vai saber
  if ($which=='batch_delete')   { 
    $recordsIds =  $params[2]; 
    $routeOk = true;
  }

  if (!  $routeOk) routeError();

  $dbConnection -> autocommit(true);    // gravacao sem necessidade de confirmacao 

  if ($which=='delete')   
  	$crudSql = "update $tablename set deleted_at = now() where id=$recordId";
  if ($which=='batch_delete')   
  	$crudSql = "update $tablename set deleted_at = now() where id in($recordsIds) ";


}

/*********************************************************************************
 rota: GET   /:workgroup/report

 equivalente backend golang= WorkgroupReport

 obtem a qtde de alteracoes (inc/exc/edicao) de dados feitas pelo grupo atual
*********************************************************************************/
if ( $method=='GET' && 
    isset($params[1]) && $params[1] == 'report' ) {

  $workgroup = $params[0];   

	$fetchSql = "select ifnull(database_changes_amount, 0) as database_changes_amount ".
              "from workgroups ".
              "where trim(name)=trim('$workgroup') and ifnull(active, false)=true and deleted_at is null ";
}

/*********************************************************************************
 rota: GET   :workgroup/car_cards/:country

 equivalente backend golang= GetCarsForCards

 obtem dados para preenchimento de cards de carros no front end
*********************************************************************************/
if ( $method=='GET' && 
    isset($params[2]) && $params[1] == 'car_cards' ) {

  $country = $params[2];   
  $workgroup = $params[0];   

  $routeOk = false;
	if ($country == "usa" || $country == "brazil")
    $routeOk = true;

  if (!  $routeOk) routeError();

  // gera string randomica para concatenar com link da imagem e tentar evitar cache do navegador
  $tempLink = rand(10000,99999);

	$fetchSql = "select cars.id as id, cars.name as name, rental_price as rental_price, concat(cars.workgroup, '_car_', ".
              "     LPAD(cars.id, 6, '0'), '.png?$tempLink') as car_image ".
              "from cars ".
              "where country = '$country' and ifnull(active, false) = true and workgroup='$workgroup' and deleted_at is null  ".
              " order by rand() ";
}

/*********************************************************************************
 rota: GET  /:workgroup/:table_name/itens_for_autocomplete

 equivalente backend golang= ItensForAutocomplete

 obtem items para oferecer como sugestoes em campo <input text>, usuario começa a
 digitar e a aplicacao oferece as sugestoes baseado nas letras que usuario ja digitou 
*********************************************************************************/
if ( $method=='GET' && 
    isset($params[2]) && $params[2] == 'itens_for_autocomplete' ) {

  $workgroup = $params[0];     
  $tablename = $params[1];   
    
	$fetchSql = "SELECT id as item_id, name as item_name ".
              "FROM $tablename ".
              " where ifnull(active, false) = true and deleted_at is null and workgroup = '$workgroup' order by name";

  try {
    $result = mysqli_query($dbConnection, $fetchSql); 

  } catch(Exception $e)  {
    internalError( mysqli_error($dbConnection) );
  }

					
	// monta string com registros e colunas para autocomplete
	// ao inves de retornar como JSON, vai retornar uma string, para facilitar o processamento pelo jscript
	$itens = '';
	while ( $row = mysqli_fetch_object($result) ) {
		if ($itens != "") 	$itens .= "|";     // separador de registros
		$itens .= $row->item_id . ";" . $row->item_name;
	}

  header("Content-Type: text/plain");
  http_response_code(200);   // 200= requisicao bem sucedida

  die($itens);
}



/************************************************************************************************************************************
 rotas: GET  /:workgroup/cars   
        GET  /:workgroup/manufacturers
        GET  /:workgroup/terms
        GET  /:workgroup/workgroups

 equivalente backend golang= GetCarsForDatatable
                             GetManufacturersForDatatable
                             GetTermsForDatatable
                             GetWorkgroupsForDatatable

 obtem dados de carros/fabricantes/exrpessoes ingles, portugues (termos)/grupos de trabalho para preenchimento da datatable no front
************************************************************************************************************************************/
if ( ( ( $method=='GET' && isset($params[1]) && $params[1] == 'cars' ) ||
      ( $method=='GET' && isset($params[1]) && $params[1] == 'manufacturers' ) ||   
      ( $method=='GET' && isset($params[1]) && $params[1] == 'terms' )   ||
      ( $method=='GET' && isset($params[1]) && $params[1] == 'workgroups' ) ) && 
      ! isset($params[2]) )  {  

  $workgroup = $params[0];   
  $tablename = $params[1];   

  // separa parametros que vieram na URL 
  // URLs serao parecidas com estas:
  // http://localhost/admin/cars?country=brazil&order_by=name&order_direction=asc&search_txt=&only_active_or_inactive_records=
  // http://localhost/admin/manufacturers?country=brazil&order_by=name&order_direction=asc&search_txt=&only_active_or_inactive_records=
  // http://localhost/admin/terms?country=brazil&order_by=name&order_direction=asc&search_txt=&only_active_or_inactive_records=
  // http://localhost/admin/workgroups?country=brazil&order_by=name&order_direction=asc&search_txt=&only_active_or_inactive_records=

  $url= parse_url($_SERVER['REQUEST_URI']);
  parse_str($url['query'], $details);

  // exige todas os parametros de montagem da datatable no front end
  if ( ! isset($details['order_by']) || ! isset($details['country']) || ! isset($details['order_direction']) || ! isset($details['search_txt']) ||
       !  isset($details['only_active_or_inactive_records'])  ) 
    routeError();


  $orderBy = $details['order_by'];
  $country = $details['country'];
  $orderDirection = $details['order_direction'];
  $searchText = "'%". $details['search_txt'] . "%'";
  $onlyActive = $details['only_active_or_inactive_records'];


  // gera string randomica para concatenar com link da imagem e tentar evitar cache do navegador
  $tempLink = rand(10000,99999);

  if ($tablename=='cars') 
	  // arq imagem do carro:   workgroup_car_999999.png
    $fetchSql =
      " select  cars.id , cars.name,  concat(cars.workgroup, '_car_', LPAD(cars.id, 6, '0'), '.png?$tempLink') as car_image,  " .
      "         ifnull(cars.active, 0) as active, concat(cars.workgroup, '_manufacturer_', LPAD(manufacturers.id, 6, '0'), '.png?$tempLink') as manufacturer_logo  ".
      "from cars  ".
      "left join manufacturers on manufacturers.id = cars.manufacturer_id ".
      "where country = '$country' and cars.workgroup= '$workgroup' ";

  if ($tablename=='manufacturers') 
	  // arq logotipo do fabricante:   workgroup_manufacturer_999999.png
    // nao separa fabricantes por país (country), fabricantes sao multinacionais usados por ambos os países
    $fetchSql =
      "select id , name, ifnull(active, 0) as active, concat(workgroup, '_manufacturer_', LPAD(id, 6, '0'), '.png?$tempLink') as manufacturer_logo  ".
      "from manufacturers ".
      "where workgroup= '$workgroup' ";
  
  if ($tablename=='terms') 
    $fetchSql =
      "select id, portuguese, english, item, description, ifnull(active, 0) as active  ".
      "from terms ".
      "where workgroup= '$workgroup' ";

  if ($tablename=='workgroups') {
    // formato de datas muda conforme país usado no front 
    if ($country == "usa" )
      $fields =
        " id, name, ifnull(in_use, false) as in_use, ifnull(client_ip, '') as client_ip, ifnull(client_country, '') as client_country, ".
        " ifnull(client_city, '') as client_city, ifnull(client_loc, '') as client_loc, " .
          " date_format(updated_at, '%m/%d/%Y %l:%i - %p') as updated_at, " .
          " date_format(deleted_at, '%m/%d/%Y %l:%i - %p') as deleted_at, ifnull(active, false) as active ";
    

    if ($country == "brazil" )
      $fields =
        " id, name, ifnull(in_use, false) as in_use, ifnull(client_ip, '') as client_ip, ifnull(client_country, '') as client_country, ".
        " ifnull(client_city, '') as client_city, ifnull(client_loc, '') as client_loc , " .
          " date_format(updated_at, '%d/%m/%Y %H:%i') as updated_at, " .
          " date_format(deleted_at, '%d/%m/%Y %H:%i') as deleted_at, ifnull(active, false) as active ";
    
    $fetchSql =
      "select $fields  ".
      "from workgroups ".
      "where 1=1 ";

    // se o grupo atual no front end nao é 'admin', exibe somente grupos disponiveis
    // nao deixa usuario comum ver grupos abertos por outras pessoas para que ele nao possa ingressar
    // usuario comum pode compartilhar o proprio grupo com outras pessoas
    if ( strtolower(trim($workgroup)) != "admin" ) {
      $fetchSql .= " and ifnull(in_use, false) = false ";
    }
  }


  $where = '';
	if ($searchText != "")  {  
    if ($tablename=='cars') 		
      $where .= " and (cars.name like $searchText or manufacturers.name like $searchText) ";

    if ( ($tablename=='manufacturers') || 
         ($tablename=='workgroups') )
      $where .= " and (name like $searchText) ";

    if ($tablename=='terms') 		
      $where .= " and (item like $searchText or portuguese like $searchText or english like $searchText or description like $searchText) ";


  }
	
  // so ativos
	if ($onlyActive == "active") 
		$where .= " and ifnull($tablename.active, false) = true ";
	
	// so inativos		
  if ($onlyActive == "inactive") 
		$where .= " and ifnull($tablename.active, false) = false ";

  // nao considera exluidos logicamente
  $where .= " and $tablename.deleted_at is null ";
	
  $fetchSql .= $where;

	// se nao foi passado campo order by inicialmente
	if ($orderBy == "") {
		$orderBy = "$tablename.name";
		$oderDirection = "asc";
	}

  $fetchSql .= "order by $orderBy $orderDirection";

}



/************************************************************************************************************************
 rotas: POST   /:workgroup/car   , inserir carro
        PATCH  /:workgroup/car/:id  , editar carro   (sera usado POST, conforme explicado no topo do arquivo)

 equivalente backend golang= SaveCar

 insere registro de carro ou edita registro existente
************************************************************************************************************************/
if ( ($method=='POST' &&  isset($params[1]) && $params[1] == 'car') ||
     ($method=='POST' &&  isset($params[2]) && $params[1] == 'car' && is_numeric($params[2])) )  {   

  $workgroup = $paramsOriginalCase[0];    

  // verifica se todos os campos vieram no 'form data'
  // formato [fieldName, minSize, maxSize]
  $fields = [ ['string', 'country', 3, 10],
              ['string', 'year', 4, 4 ],
              ['string', 'name', 2, 50 ],
              ['int', 'manufacturer_id', 1, 6],
              ['string', 'rental_price', 5, 9 ],
              ['string', 'odometer', 2, 7 ], 
              ['string', 'mpg', 1, 5 ],
              ['string', 'cylinders', 1, 1 ],
              ['hp', 'hp', 2, 3 ],
              ['bool', 'transmission_manual', 0, 0 ],
              ['string', 'doors', 1, 1 ],
              ['string', 'cc', 2, 5 ],
              ['string', 'client_ip', 5, 50 ] ];


  $dataError = '';
  for ($i=0; $i < count($fields); $i++)  {

    $fieldType = $fields[$i][0];
    $fieldName = $fields[$i][1];
    $minSize = $fields[$i][2];
    $maxSize = $fields[$i][3];

    $fieldValue = $_POST[$fieldName];

    // verifica se campo numerico veio como numero mesmo  
    if ($fields[$i][0] == 'int') {
      if (! is_numeric($fieldValue)) {
        $dataError = 'Not numeric / não numérico';
        break;
      }
    }

    // verifica tamanho minimo/maximo dos campos 
    if ($fieldType=='string') {
        if ( strlen($fieldValue) < $minSize || strlen($fieldValue) > $maxSize )  {
          $dataError = $fieldName . ' - String size error / erro no tamanho da string';
          break;
        }
    }
  }

  if ($dataError!='') internalError( $dataError );

  $country = $_POST['country'];
  $name = $_POST['name'];
  $year = $_POST['year'];
  $manufacturerId = $_POST['manufacturer_id'];
  $rentalPrice = $_POST['rental_price'];
  $odometer = $_POST['odometer'];
  $mpg = $_POST['mpg'];
  $cylinders = $_POST['cylinders'];
  $hp = $_POST['hp'];
  $transmissionManual = $_POST['transmission_manual'];
  $doors = $_POST['doors'];
  $cc = $_POST['cc'];
  $clientIp = $_POST['client_ip'];
  $bypassImageUpload = $_POST['bypass_image_upload'];

  // verifica a imagem do carro, caso tenha sido escolhida
  if ( $bypassImageUpload == 'false' )  {
    // arquivo veio?
    $imgInfo = getimagesize($_FILES['chosen_image_file']['tmp_name']);
    if ($imgInfo === FALSE) 
      internalError( 'File erro / Erro no arquivo');

    // so aceita PNG
    if ($imgInfo[2] !== IMAGETYPE_PNG) 
      internalError( 'Image must be PNG / Imagem deve ser PNG');

    // tamanho maximo
    if ($_FILES[$fileName]['size'] > 1500000) 
      internalError( 'Max size 1.5 MB / Tamanho máximo 1,5 MB');
  }

  $carId = '';   

  // se nao foi passado ID, é uma insercao de registro de carro
  if (! isset($params[2]))    {
    $crudSql = "insert into cars(name, country, year, manufacturer_id, rental_price, odometer, mpg, cylinders, hp, transmission_manual, doors, cc, workgroup, active, created_at, updated_at) ". 
               "select '$name', '$country', '$year', $manufacturerId, '$rentalPrice', '$odometer', '$mpg', ".
               "       '$cylinders', '$hp', $transmissionManual, '$doors', '$cc', '$workgroup', true, now(), now() ";
    $dbOperation = 'insert';
  }

  // se foi passado ID, é uma edicao de registro
  if ( isset($params[2]))   {
    $carId = $params[2];   

    $crudSql = "update cars set name='$name', year='$year', manufacturer_id=$manufacturerId, rental_price='$rentalPrice', odometer='$odometer', mpg='$mpg', ".
              "        cylinders='$cylinders', hp='$hp', transmission_manual=$transmissionManual, doors='$doors', cc='$cc', updated_at=now() ". 
               "where id = $carId ";
    $dbOperation = 'update';
  }


  $anyNotificationToMake = "$dbOperation|cars|$carId";   // avisa que deve ser feito log da operacao 'operation', em relacao à tabela 'cars'

  if ( $bypassImageUpload == 'false' )  
    $dbConnection -> autocommit(false);    // so vai efetivar o sql abaixo apos o upload da imagem do carro ter sido feito com sucesso
  else 
    $dbConnection -> autocommit(true);   

  // efetua SQL e pede o ID do registro recem manipulado (3o parametro, true)
  $result = executeCrudQueryAndReturnResult($crudSql, true, $anyNotificationToMake, $workgroup, $clientIp);  

  
  // qdo usuario nao alterou imagem do carro, nao eh necessario subi la novamente (front end avisa: bypass_image_upload= true)
  if ( $bypassImageUpload == 'false' )  {
  
    /*
     obtem ID do registro que sera usado no nome do arquivo de imagem
    */

    // se foi insercao, obtem ID do registro recem adicionado
    if ($dbOperation == 'insert')  {
      if ( strpos($result, "__success__")  === false ) {}      
      else {
        $tmp = explode('|', $result);     // result= "__success__|$newRecordId";
        $carId = $tmp[1];
      } 
    }
    // nao edicao, ID foi passado na url
    if ($dbOperation == 'update')  
      $carId = $params[2];   


    // sobe para AWS S3 a imagem escolhida
    uploadImageToAWS_S3( 'car', 'chosen_image_file', $carId, $workgroup );
  }

  // alteracoes e upload imagem bem sucedidos, confirma alteracoes
  if ( $bypassImageUpload == 'false' )  
    $dbConnection -> commit();

  // devolve ao front end o ID do registro manipulado (inserido/editado) 
  // para que o front end saiba em qual registro colcoar destaque na datatable
  if ($dbOperation == 'update')   die( '__success__' );
  else die( $result );    // __success__|id registro
}


/*********************************************************************************
 rotas: GET   /car/:id                       , busca alguns detalhes do carro
        GET   /car/:id/:first_day/:last_day  , busca alguns detalhes do carro e os dias do mes/ano (dentro de um periodo) em que ele ja esta reservado
       
 equivalente backend golang= GetCar

 ********************************************************************************/
if ( ($method=='GET' &&  isset($params[1]) && $params[0] == 'car') ||
     ($method=='GET' &&  isset($params[3]) && $params[0] == 'car') ) {

  $carId = $params[1];   

	if (! is_numeric($carId))   routeError();

  $reservedDays;     // em principio, nao vai pesquisar periodo

  // se informadas datas para pesquisa, busca reservas daquele carro feitas no periodo
  // verifica validade das datas
  if (isset($params[3])) {
      $iniDate = $params[2];
      $endDate = $params[3];

      for ($dates=1; $dates<=2; $dates++)  {
        if ($dates==1) $tmp = explode('-', $iniDate);
        if ($dates==2) $tmp = explode('-', $endDate);

        if ( count($tmp)<3 )  routeError();

        // checkdate(mes, dia, ano), $iniDate= yyyy-mm-dd
        if (! checkdate($tmp[1], $tmp[2], $tmp[0]))  routeError();
      }



      $tmpSql .=  " SELECT DAY(pickup_datetime) as pickup, DAY(dropoff_datetime) as dropoff " .
                  " from bookings " .
                  " where " .
                  "  ( date_format(pickup_datetime, '%Y-%m-%d') between '$iniDate' and '$endDate' or " .
                  " date_format(dropoff_datetime, '%Y-%m-%d') between '$iniDate' and '$endDate' )  and car_id = $carId " .
                  " AND deleted_at IS null ";

      

      try {
        $bookings = mysqli_query($dbConnection, $tmpSql) or die (mysqli_error($conexao).'-'.$sql);
      } catch(Exception $e)  {
        internalError( mysqli_error($dbConnection) );
      }

      // monta  string com a lista de dias em que o carro esta ocupado, separados por ','
      while ($row = mysqli_fetch_object($bookings))  {

        for ($day = $row->pickup; $day <= $row->dropoff; $day++)  {
          if ($reservedDays != "") {
            $reservedDays .= ",";
          }
          $reservedDays .= $day;
        }
      }
  }

 
	// logo_filename ficara parecido com: 	manufacturer_000002.png, manufacturer_000041.png, etc
	// car_image ficara parecido com: 	car_000002.png, car_000041.png, etc

  // gera string randomica para concatenar com link da imagem e tentar evitar cache do navegador
  $tempLink = rand(10000,99999);

	$fetchSql = " select cars.name, manufacturer_id, country, manufacturers.name as manufacturer_name, rental_price, cars.id, " .
              " concat(cars.workgroup, '_car_', LPAD(cars.id, 6, '0'), '.png?$tempLink') as car_image, year, doors, odometer, mpg, cylinders, hp, cc,  " .
              " concat(cars.workgroup, '_manufacturer_', LPAD(cars.manufacturer_id, 6, '0'), '.png?$tempLink') as manufacturer_logo, " .
              " ifnull(transmission_manual, false) as transmission_manual, '$reservedDays' as days_reserved ".
              " from cars ".
              " left join ".
              "   manufacturers on manufacturers.id = cars.manufacturer_id ".
              " where cars.id = $carId ";
  $noArray = true;   // avisa para enviar JSON sozinho (nao array), o frontend aguarda como no exemplo abaixo:
    /*
    {
        "country": "usa",
        "year": "2025",
        "name": "nome",
        "manufacturer_name": "Wolkswagen",
        "manufacturer_logo": "Accomplished_manufacturer_000173.png",
        "rental_price": "550,00",
        "odometer": "27000",
        "car_image": "Accomplished_car_001461.png",
        "mpg": "15",
        "cylinders": "4",
        "hp": "150",
        "transmission_manual": "1",
        "doors": "5",
        "cc": "1600",
        "days_reserved": ""
    }
    */
}




/*********************************************************************************
 rotas: GET   /manufacturer/:id          , busca registro de fabricante
       
 equivalente backend golang= GetManufacturer

 ********************************************************************************/
if ($method=='GET' &&  isset($params[1]) && $params[0] == 'manufacturer') {

  $manufacturerId = $params[1];   

	if (! is_numeric($manufacturerId))   routeError();
 
	// logo_filename ficara parecido com: 	manufacturer_000002.png, manufacturer_000041.png, etc

  // gera string randomica para concatenar com link da imagem e tentar evitar cache do navegador
  $tempLink = rand(10000,99999);
	$fetchSql = " select name, concat(workgroup, '_manufacturer_', LPAD(id, 6, '0'), '.png?$tempLink') as manufacturer_logo  ".
              " from manufacturers ".
              " where id = $manufacturerId ";

  $noArray = true;   // avisa para enviar JSON sozinho (nao array), o frontend aguarda como no exemplo abaixo:
    /*
      {
          "id": "",
          "name": "Wolkswagen",
          "manufacturer_logo": "Accomplished_manufacturer_000173.png",
          "active": "1"
      }
    */
}


/************************************************************************************************************************
 rotas: POST   /:workgroup/manufacturer        , inserir fabricante
        PATCH  /:workgroup/manufacturer/:id    , editar fabricante   (sera usado POST, conforme explicado no topo do arquivo)

 equivalente backend golang= SaveManufacturer

 insere registro de fabricante ou edita registro existente
************************************************************************************************************************/
if ( ($method=='POST' &&  isset($params[1]) && $params[1] == 'manufacturer') ||
     ($method=='POST' &&  isset($params[2]) && $params[1] == 'manufacturer' && is_numeric($params[2])) )  {   

  $workgroup = $paramsOriginalCase[0];    
  
  // verifica se todos os campos vieram no 'form data'
  // formato [fieldName, minSize, maxSize]
  $fields = [ ['string', 'name', 3, 50] ]  ;

  $dataError = '';
  for ($i=0; $i < count($fields); $i++)  {

    $fieldType = $fields[$i][0];
    $fieldName = $fields[$i][1];
    $minSize = $fields[$i][2];
    $maxSize = $fields[$i][3];

    $fieldValue = $_POST[$fieldName];

      // verifica tamanho minimo/maximo dos campos 
    if ($fieldType=='string') {
        if ( strlen($fieldValue) < $minSize || strlen($fieldValue) > $maxSize )  {
          $dataError = $fieldName . ' - String size error / erro no tamanho da string';
          break;
        }
    }
  }

  if ($dataError!='') internalError( $dataError );

  $name = $_POST['name'];
  $clientIp = $_POST['client_ip'];
  $bypassImageUpload = $_POST['bypass_image_upload'];

  // verifica logotipo do fabricante, caso tenha sido escolhido
  if ( $bypassImageUpload == 'false' )  {
    // arquivo veio?
    $imgInfo = getimagesize($_FILES['chosen_image_file']['tmp_name']);
    if ($imgInfo === FALSE) 
      internalError( 'File erro / Erro no arquivo');

    // so aceita PNG
    if ($imgInfo[2] !== IMAGETYPE_PNG) 
      internalError( 'Image must be PNG / Imagem deve ser PNG');

    // tamanho maximo
    if ($_FILES[$fileName]['size'] > 1500000) 
      internalError( 'Max size 1.5 MB / Tamanho máximo 1,5 MB');
  }


  $manufacturerId = '';

  // se nao foi passado ID, é uma insercao de registro de fabricante
  if (! isset($params[2]))    {
    $crudSql = "insert into manufacturers(name, active, workgroup, created_at, updated_at) ". 
               "select '$name', true, '$workgroup', now(), now() ";
    $dbOperation = 'insert';
  }

  // se foi passado ID, é uma edicao de registro
  if ( isset($params[2]))   {
    $manufacturerId = $params[2];   

    $crudSql = "update manufacturers set name='$name', updated_at=now() ". 
               "where id = $manufacturerId ";
    $dbOperation = 'update';
  }

  $anyNotificationToMake = "$dbOperation|manufacturers|$manufacturerId";   // avisa que deve ser feito log da operacao 'operation', em relacao à tabela 'manufacturers'

  if ( $bypassImageUpload == 'false' )  
    $dbConnection -> autocommit(false);    // so vai efetivar o sql abaixo apos o upload do logotiupo do fabricante ter sido feito com sucesso
  else 
    $dbConnection -> autocommit(true);    

  // efetua SQL e pede o ID do registro recem manipulado (3o parametro, true)
  $result = executeCrudQueryAndReturnResult($crudSql, true, $anyNotificationToMake, $workgroup, $clientIp);  

  
  // qdo usuario nao alterou logotipo do fabricante, nao eh necessario subi lo novamente (front end avisa: bypass_image_upload= true)
  if ( $bypassImageUpload == 'false' )  {
  
    /*
     obtem ID do registro que sera usado no nome do arquivo de imagem
    */

    // se foi insercao, obtem ID do registro recem adicionado
    if ($dbOperation == 'insert')  {
      if ( strpos($result, "__success__")  === false ) {}      
      else {
        $tmp = explode('|', $result);    // result= "__success__|$newRecordId";
        $manufacturerId = $tmp[1];
      } 
    }
    // nao edicao, ID foi passado na url
    if ($dbOperation == 'update')  
      $manufacturerId = $params[2];   


    // sobe para AWS S3, o logotipo do fabricante
    uploadImageToAWS_S3( 'manufacturer', 'chosen_image_file', $manufacturerId, $workgroup );

  }

  // alteracoes e upload do logotipo bem sucedidos, confirma alteracoes
  if ( $bypassImageUpload == 'false' )  
    $dbConnection -> commit();

  // devolve ao front end o ID do registro manipulado (inserido/editado) 
  // para que o front end saiba em qual registro colcoar destaque na datatable
  if ($dbOperation == 'update')   die( '__success__' );
  else die( $result );    // __success__|id registro
}





/*************************************************************************************************
 rotas: GET   /term/:id      , busca registro de expressao ingles/portugues (term)
       
 equivalente backend golang= GetManufacturer

*************************************************************************************************/
if ($method=='GET' &&  isset($params[1]) && $params[0] == 'term')  {

  $termId = $params[1];   

	if (! is_numeric($termId))   routeError();
 
	$fetchSql = " select portuguese, english, item, description  ".
              " from terms ".
              " where id = $termId ";

  $noArray = true;   // avisa para enviar JSON sozinho (nao array), o frontend aguarda como no exemplo abaixo:
    /*
      {
          "item": "api_fetch_visible",
          "portuguese": "Exibir / ocultar chamadas API ",
          "english": "Show / hide api calls",
          "description": "yes yes yes"
      }
    */
}




/************************************************************************************************************************
 rotas: POST   /:workgroup/term        , inserir expressao ingles/portugues (term)
        PATCH  /:workgroup/term/:id    , editar expressao ing/port   (sera usado POST, conforme explicado no topo do arquivo)

 equivalente backend golang= SaveTerm

 insere registro de expressao ing/port ou edita registro existente
************************************************************************************************************************/
if ( ($method=='POST' &&  isset($params[1]) && $params[1] == 'term') ||
     ($method=='POST' &&  isset($params[2]) && $params[1] == 'term' && is_numeric($params[2])) )  {   

  $workgroup = $paramsOriginalCase[0];    
  
  // verifica se todos os campos vieram no 'form data'
  // formato [fieldName, minSize, maxSize]
  $fields = [ ['string', 'item', 5, 200]  ,
              ['string', 'portuguese', 2, 1200]  ,
              ['string', 'english', 2, 1200]  ,
              ['string', 'description', 10, 200] ];

  $dataError = '';
  for ($i=0; $i < count($fields); $i++)  {

    $fieldType = $fields[$i][0];
    $fieldName = $fields[$i][1];
    $minSize = $fields[$i][2];
    $maxSize = $fields[$i][3];

    $fieldValue = $_POST[$fieldName];

      // verifica tamanho minimo/maximo dos campos 
    if ($fieldType=='string') {
        if ( strlen($fieldValue) < $minSize || strlen($fieldValue) > $maxSize )  {
          $dataError = $fieldName . ' - String size error / erro no tamanho da string';
          break;
        }
    }
  }

  if ($dataError!='') internalError( $dataError );

  $item = $_POST['item'];
  $portuguese = $_POST['portuguese'];
  $english = $_POST['english'];
  $description = $_POST['description'];
  $clientIp = $_POST['client_ip'];

  $termId = '';   

  // se nao foi passado ID, é uma insercao de registro 
  if (! isset($params[2]))    {
    $crudSql = "insert into terms(item, portuguese, english, description, active, workgroup, created_at, updated_at) ". 
               "select '$item', '$portuguese', '$english', '$description',  true, '$workgroup', now(), now() ";
    $dbOperation = 'insert';
  }

  // se foi passado ID, é uma edicao de registro
  if ( isset($params[2]))   {
    $termId = $params[2];   

    $crudSql = "update terms set item='$item', portuguese='$portuguese', english='$english', description='$description', updated_at=now() ". 
               "where id = $termId ";
    $dbOperation = 'update';
  }

  $anyNotificationToMake = "$dbOperation|terms|$termId";   // avisa que deve ser feito log da operacao 'operation', em relacao à tabela 'terms'

  $dbConnection -> autocommit(true);    // gravacao sem necessidade de confirmacao 

  // efetua SQL e pede o ID do registro recem manipulado (3o parametro, true)
  $result = executeCrudQueryAndReturnResult($crudSql, true, $anyNotificationToMake, $workgroup, $clientIp);  


  // devolve ao front end o ID do registro manipulado (inserido/editado) 
  // para que o front end saiba em qual registro colcoar destaque na datatable
  if ($dbOperation == 'update')   die( '__success__' );
  else die( $result );    // __success__|id registro
}




/*******************************************************************************************************************************************
 rotas: GET   /:workgroup/bookings/:country/:car_id/:first_day_week/:last_day_week  , busca reservas feitas por veiculo em determinado periodo
       
 equivalente backend golang= GetBookingsToPopulateSchedule

*******************************************************************************************************************************************/
if ($method=='GET' &&  isset($params[5]) && $params[1] == 'bookings')  {  

  $workgroup = $params[0];   
  $country = $params[2];   
  $carId = $params[3];   
  $firstDayWeek = $params[4];   
  $lastDayWeek = $params[5];   

	if (! is_numeric($carId))   routeError();
 
  // gera string randomica para concatenar com link da imagem e tentar evitar cache do navegador
  $tempLink = rand(10000,99999);

	$fetchSql = " select bookings.driver_name,  bookings.car_id, concat(cars.workgroup, '_car_', LPAD(cars.id, 6, '0'), '.png?$tempLink') as car_image, " .
              " if(bookings.country='usa', date_format(pickup_datetime, '%m/%d %h:%i - %p'), date_format(pickup_datetime, '%d/%m - %H:%i')) as pickup_formatted,   " .
              " if(bookings.country='usa', date_format(dropoff_datetime, '%m/%d %h:%i - %p'), date_format(dropoff_datetime, '%d/%m - %H:%i')) as dropoff_formatted,   " .
              " date_format(pickup_datetime, '%Y-%m-%d|%H:%i') as pickup_reference,   date_format(dropoff_datetime, '%Y-%m-%d|%H:%i') as dropoff_reference, " .
              " bookings.id as booking_id ".
              " from bookings ".
              " left join cars on bookings.car_id = cars.id " .
              " where 1=1 ";

  // car_id = -1, usuario pediu para ver reservas de todos os carros, so filtra por país
	if ($carId == '-1') {
		$fetchSql .= " and bookings.country='$country' ";
	}
	// se usuario indicou qual carro listar agenda
	if ($carId != '-1') {
		$fetchSql .= " and bookings.car_id = $carId ";
  }

  $fetchSql .= " and (DATE_FORMAT(pickup_datetime,'%Y-%m-%d') between '$firstDayWeek' and '$lastDayWeek' or DATE_FORMAT(dropoff_datetime,'%Y-%m-%d') between '$firstDayWeek' and '$lasttDayWeek') ";
  $fetchSql .= " and bookings.workgroup = '$workgroup' ";
  $fetchSql .= " AND bookings.deleted_at IS null ";

   /*
    exemplo de resultado que deve ser enviado ao front:
    [
        {
            "booking_id": 56,
            "car_id": 1462,
            "pickup_formatted": "27/01 - 16:30",
            "pickup_reference": "2025-01-27|16:30",
            "dropoff_formatted": "28/01 - 11:30",
            "dropoff_reference": "2025-01-28|11:30",
            "driver_name": "Peter ",
            "car_image": "Accomplished_car_001462.png"
        },
        {
            "booking_id": 57,
            "car_id": 1462,
            "pickup_formatted": "30/01 - 11:45",
            "pickup_reference": "2025-01-30|11:45",
            "dropoff_formatted": "30/01 - 16:20",
            "dropoff_reference": "2025-01-30|16:20",
            "driver_name": "Liliana",
            "car_image": "Accomplished_car_001462.png"
        }
    ]
    */
}





/************************************************************************************************************************
 rotas: POST   :workgroup/booking        , inserir reserva de carro
        PATCH  :workgroup/booking/:id    , editar reserva de carro

 equivalente backend golang= SaveBooking 

 insere registro de reserva de carro, separado por workgroup / country
************************************************************************************************************************/
if ( ($method=='POST' &&  isset($params[1]) && $params[1] == 'booking') ||
     ($method=='POST' &&  isset($params[2]) && $params[1] == 'booking' && is_numeric($params[2])) )  {   

  $workgroup = $paramsOriginalCase[0];    
  
  // verifica se todos os campos vieram no 'form data'
  // formato [fieldName, minSize, maxSize]
  $fields = [ ['int', 'car_id', 1, 6]  ,  
              ['datetime', 'pickup_datetime', 5, 20]  ,
              ['datetime', 'dropoff_datetime', 5, 200]  ,
              ['string', 'driver_name', 3, 50],
              ['string', 'country', 3, 10], 
              ['string', 'client_ip', 5, 50 ] ];

  $dataError = '';
  for ($i=0; $i < count($fields); $i++)  {

    $fieldType = $fields[$i][0];
    $fieldName = $fields[$i][1];
    $minSize = $fields[$i][2];
    $maxSize = $fields[$i][3];

    $fieldValue = $_POST[$fieldName];

    // verifica se campo numerico veio como numero mesmo  
    if ($fields[$i][0] == 'int') {
      if (! is_numeric($fieldValue)) {
        $dataError = 'Not numeric / não numérico';
        break;
      }
    }

    // verifica tamanho minimo/maximo dos campos 
    if ($fieldType=='string') {
        if ( strlen($fieldValue) < $minSize || strlen($fieldValue) > $maxSize )  {
          $dataError = $fieldName . ' - String size error / erro no tamanho da string';
          break;
        }
    }
  }

  if ($dataError!='') internalError( $dataError );

  $carId = $_POST['car_id'];
  $pickupDatetime = $_POST['pickup_datetime'];
  $dropoffDatetime = $_POST['dropoff_datetime'];
  $driverName = $_POST['driver_name'];
  $country = $_POST['country'];
  $clientIp = $_POST['client_ip'];

  $bookingId = '';   

  // se nao foi passado ID, é uma insercao de registro 
  if (! isset($params[2]))    {
    $crudSql = "insert into bookings(car_id, country, pickup_datetime, dropoff_datetime, driver_name, workgroup, created_at, updated_at) ". 
               "select $carId, '$country', '$pickupDatetime', '$dropoffDatetime', '$driverName', '$workgroup', now(), now() "; 
    $dbOperation = 'insert';
  }

  // se foi passado ID, é uma edicao de registro
  if ( isset($params[2]))   {
    $bookingId = $params[2];   

    $crudSql = "update bookings set pickup_datetime='$pickupDatetime', dropoff_datetime='$dropoffDatetime', driver_name='$driverName', updated_at=now() ". 
               "where id = $bookingId ";
    $dbOperation = 'update';
  }

  $anyNotificationToMake = "$dbOperation|bookings|$bookingId";   // avisa que deve ser feito log da operacao 'operation', em relacao à tabela 'bookings'

  $dbConnection -> autocommit(true);    // gravacao sem necessidade de confirmacao 

  // efetua SQL e pede o ID do registro recem manipulado (3o parametro, true)
  $result = executeCrudQueryAndReturnResult($crudSql, true, $anyNotificationToMake, $workgroup, $clientIp);    

  // devolve ao front end o ID do registro manipulado (inserido/editado) 
  // para que o front end saiba em qual registro colcoar destaque na datatable
  if ($dbOperation == 'update')   die( '__success__' );
  else die( $result );    // __success__|id registro
}




/*************************************************************************************************
 rotas: GET   /booking/:id      , busca registro de reserva de carro
       
 equivalente backend golang= GetBooking

*************************************************************************************************/
if ($method=='GET' &&  isset($params[1]) && $params[0] == 'booking')  {

  $bookingId = $params[1];   

	if (! is_numeric($bookingId))   routeError();
 
	$fetchSql = " select 	bookings.car_id, driver_name, " .
              " if(bookings.country='usa', date_format(pickup_datetime, '%m/%d/%y'), date_format(pickup_datetime, '%d/%m/%y')) as pickup_date,   " .
              " if(bookings.country='usa', date_format(pickup_datetime, '%h:%i - %p'), date_format(pickup_datetime, '%H:%i')) as pickup_hour,   " .
              " if(bookings.country='usa', date_format(dropoff_datetime, '%m/%d/%y'), date_format(dropoff_datetime, '%d/%m/%y')) as dropoff_date,   " .
              " if(bookings.country='usa', date_format(dropoff_datetime, '%h:%i - %p'), date_format(dropoff_datetime, '%H:%i')) as dropoff_hour   ".
              ' from bookings '.
              " where bookings.id = $bookingId ";

  $noArray = true;   // avisa para enviar JSON sozinho (nao array), o frontend aguarda como no exemplo abaixo:
    /*
      {
          "item": "api_fetch_visible",
          "portuguese": "Exibir / ocultar chamadas API ",
          "english": "Show / hide api calls",
          "description": "yes yes yes"
      }
    */
}


/*************************************************************************************************
 rotas: POST   /gravar_ip/      , gravar ip do cliente que acessou o portfolio
       
 ************************************************************************************************/

if ( $method=='POST' && isset($params[0]) && $params[0] == 'gravar_ip' )     {

      $_POST = json_decode(file_get_contents("php://input"), true);

      $ip = $_POST['ip'];
      $country = $_POST['country'];
      $hostname= $_POST['hostname'];
      $city = $_POST['city'];
      $region = $_POST['region'];
      $country = $_POST['country'];
      $loc = $_POST['loc'];
      $org = $_POST['org'];
      $postal = $_POST['postal'];
      $timezone = $_POST['timezone'];

      $crudSql = "insert into acessos_portfolio (ip, hostname, city, region, country, loc, org, postal, timezone) ".
                  "select '$ip', '$hostname', '$city', '$region', '$country', '$loc', '$org', '$postal', '$timezone' ";
}


		

/********************************************************************************************************************
*********************************************************************************************************************

se a query foi montada com sucesso, a rota estava ok e so falta executar a busca e retornar o json com os dados

*********************************************************************************************************************
/********************************************************************************************************************/

// fetchSql= query de busca, retorna JSON com os dados
if ($fetchSql!='')
  executeFetchQueryAndReturnJsonResult($fetchSql, $noArray)  ;  

// crudSql= query de edicao de dados, retorna '__success__' ou retorna mensagem com o erro ocorrido
else if ($crudSql!='')
  executeCrudQueryAndReturnResult($crudSql)  ;  

else 
  // a rota nao esta ok, nao foi possivel montar query de busca
  routeError();


 

?>

