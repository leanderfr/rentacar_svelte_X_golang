<?php


class WorkgroupsHandler

{
  // prepara models
  private $WorkgroupsModel;

  public function __construct() {
      $this->WorkgroupsModel = new Workgroups;
  }

  //***************************************************************************************************************************************
  //***************************************************************************************************************************************

  public function GetNewWorkgroupOrResetCurrentOrChooseRandomlyAnother(string $action, string $workgroup)   {

  // impede grupo admin de ser resetado ou substituido por outro (sorteio)
	if ( ($action == "reset" || $action == "another") && ($workgroup == "" || $workgroup == "admin") ) {
    http_response_code(500);   // 500= erro interno
    die("Missing workgroup");
  }

  // action invalida
	if ($action != "reset" && $action != "another" && $action != "generate" ) {
    http_response_code(500);   // 500= erro interno
    die("Missing whatToDo info");
  }

  $this->WorkgroupsModel->GetNewWorkgroupOrResetCurrentOrChooseRandomlyAnother($action, $workgroup);


    //die("chegou = $workgroup - $action");

  }




}