<?php

class Workgroups
{

  //***************************************************************************************************************************************
  //***************************************************************************************************************************************

  public function GetNewWorkgroupOrResetCurrentOrChooseRandomlyAnother(string $workgroup, string $action): array   {

    die("chegou = $workgroup - $action");

  }



  //***************************************************************************************************************************************
  //***************************************************************************************************************************************

  public function getWorkgroups(string $active, string $searchbox, array $themes): array   {

    $sql =  "select id, name, ifnull(active, false) as active, description, image, themes, created_at, updated_at ".
            "from courses  ".
            "where deleted_at is null ";

    // the priority is to filter based on the searchbox 
    if ($searchbox!='')  {
      $sql .= "and ( trim(name) like('%$searchbox%') or trim(description) like('%$searchbox%') )   ";
    } 

    // filter by themes if received this criteria
    if (! is_null($themes) && count($themes)>0 )  {
      $sql .= sprintf(
                ' and (EXISTS (SELECT 1 FROM json_each(themes) WHERE value IN ("%s")))  ',
                implode('","', array_values( $themes )) );
    } 

    if ($active=='true') $sql .= 'and ifnull(active, false)=true';
    else if ($active=='false') $sql .= 'and ifnull(active, false)=false';

    $sql .= ' order by name';

//die($sql);

    //  if there's an error while querying, the below function will return the 500 code and the error message
    $data = executeFetchQueryAndReturnJsonResult( $sql, false );

    // loop though fields to isolate 'themes' field and put it in an individual keyed array
    foreach($data as $key => $val)   {
      foreach($val as $_key => $_val) {

        // convert string content themes to an array 
        // themes example:  ["tecnologia", "educacao"]
        if ($_key=='themes') {
          $data[$key]['themes'] = themesstr_to_themesarray($_val);
        }
      }
    }

    // according to the main backend made by Dot Digital, the result must be inside a 'courses' key
    return array('courses'=>$data);
  }

  //***************************************************************************************************************************************
  //***************************************************************************************************************************************

  public function getWorkgroupById($id): array   {
    $sql =  "SELECT id , name, description, image, themes, created_at, updated_at ".
            "FROM courses ".
            "WHERE id = '$id' ";

    $data = executeFetchQueryAndReturnJsonResult( $sql, true);
    if (empty($data)) {
        die( json_encode(array('errors'=>array('id'=>['Course not found']))) );
    }

    // convert string content themes to an array 
    // themes example:  ["tecnologia", "educacao"]

    $data['themes'] = themesstr_to_themesarray( $data['themes'] );

    return $data; 
  }



  //***************************************************************************************************************************************
  //***************************************************************************************************************************************
  public function changeStatus($id): array   {
    $sql = "UPDATE courses SET active = CASE WHEN active IS true THEN false ELSE true END, updated_at = datetime('now') where id = '$id' ";

    executeCrudQueryAndReturnResult($sql);    

    die( json_encode(array('message'=>'Course activity status changed successfully')) );
  }



  //***************************************************************************************************************************************
  // if id comes blank, will be an record insertion
  // otherwise, will be an update
  //***************************************************************************************************************************************
  public function updateOrInsertWorkgroup($id=''): array   {

    // verify request
    $fields = [ ['string', 'name', 5, 150]  ,
                ['string', 'description', 5, 150],
                ['string', 'image', 5, 150] 
              ];

    $data = json_decode(file_get_contents('php://input'), true);
    
    // analyse the received data
    $dataError = '';
    for ($i=0; $i < count($fields); $i++)  {

      $fieldType = $fields[$i][0];
      $fieldName = $fields[$i][1];
      $minSize = $fields[$i][2];
      $maxSize = $fields[$i][3];

      $fieldValue = $data[$fieldName];

      // é numerico
      if ($fields[$i][0] == 'int') {
        if (! is_numeric($fieldValue)) {
          $dataError = 'Não numérico';
          break;
        }
      }

      // larguras min / max 
      if ( strlen($fieldValue) < $minSize || strlen($fieldValue) > $maxSize )  {
        $dataError = $fieldName . ' - Invalid string length';
        break;
      }
    }

    if ($dataError!='') internalError( $dataError );

    $name =   addslashes($data['name']);
    $description =   addslashes($data['description']);
    $image =   addslashes($data['image']);
    $themes =   themesarray_to_themesstr( $data['themes'] );


    // if received ID is in blank, it's an insert
    if ($id=='')    {
      $id = uniqid();

      $sql = "INSERT INTO courses (id, name, description, image, themes, created_at, updated_at) ". 
                "select '$id', '$name', '$description', '$image', '$themes', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP "; 
    }

    // received valid ID, it's an update
    else { 
      $sql = "UPDATE courses ".
                 "SET name = '$name', description = '$description', image = '$image', themes = '$themes', updated_at = CURRENT_TIMESTAMP ".
                 "where id = '$id'";
    } 

    // if there's an error while querying, the below function itself will stop and return the error code/message
    executeCrudQueryAndReturnResult($sql);    

    // (re)reads the record to return it
    $record = self::getWorkgroupById($id);
    die( json_encode($record) );


  }




  //***************************************************************************************************************************************
  // delete record
  //***************************************************************************************************************************************
  public function deleteWorkgroup($id): array   {

    $sql = "UPDATE courses ".
          "SET deleted_at = CURRENT_TIMESTAMP ".
          "where id = '$id'";

    // if there's an error while querying, the below function itself will stop and return the error code/message
    executeCrudQueryAndReturnResult($sql);    

    die( json_encode(array('message'=>'Course deleted successfully')) );


  }


}