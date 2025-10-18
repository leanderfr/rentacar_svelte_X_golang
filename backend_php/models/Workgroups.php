<?php


class Workgroups
{

  //***************************************************************************************************************************************
  //***************************************************************************************************************************************

  public function GetNewWorkgroupOrResetCurrentOrChooseRandomlyAnother(string $action, string $workgroup)   {


    // foi pedido para gerar grupo novo (1a utilizacao) ou criar outro grupo (reset)
    if ($action == "generate" || $action == "another") {
      $sql = "select id, name from workgroups where ifnull(in_use, false)=false order by rand() limit 1";
      $workgroup = executeFetchQueryAndReturnJsonResult($sql);
      $workgroup_id = $workgroup['id'];

      // dados do cliente que acessou a aplicacao
      $ip = $_POST['ip'];      
      $city = $_POST['city'];      
      $country = $_POST['country'];      
      $hostname = $_POST['hostname'];      
      $loc = $_POST['loc'];      
      $org = $_POST['org'];      
      $postal = $_POST['postal'];      
      $region = $_POST['region'];      
      $timezone = $_POST['timezone'];      

      $sql = "update workgroups set active=true, in_use=true, client_ip = '$ip', client_city='$city', client_country='$country', client_hostname='$hostname', ".
             "client_loc='$loc', client_org='$org', client_postal='$postal', client_region='$region', client_timezone='$timezone', updated_at=now() ".
             "where id='$workgroup_id' ";

      $result = executeCrudQueryAndReturnResult($sql);
    } 

      die("chegou = $workgroup - $action - $result - $sql");

  }

}




?>