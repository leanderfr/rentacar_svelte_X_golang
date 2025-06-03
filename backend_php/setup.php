<?php


require 'definicoes.php';

// se estiver rodando em desenvolvimento, mostra erro e warning
if ( strpos($_SERVER['HTTP_HOST'], 'localhost')!==false ) {
  ini_set('display_errors', '1');
  ini_set('display_startup_errors', '1');
  error_reporting(E_ALL);
}

// em producao nao mostra warning, principalmente
else {  
  error_reporting(0);
  ini_set('display_errors', 0);
}

 
// compartilha a conexao por toda aplicacao
$dbConnection = mysqli_connect($server, $login, $password);
if (mysqli_connect_errno())    {
    echo "Erro conectar base: " . mysqli_connect_error();
}

mysqli_select_db($dbConnection, $database) or die(mysqli_error($dbConnection));

mysqli_query($dbConnection, "SET NAMES 'utf8'");
mysqli_query($dbConnection,'SET character_set_connection=utf8');
mysqli_query($dbConnection,'SET character_set_client=utf8');
mysqli_query($dbConnection,'SET character_set_results=utf8');







?>