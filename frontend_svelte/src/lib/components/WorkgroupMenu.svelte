

<script>

// variaveis genericas, writeables
import { Terms  } from '$lib/stores/stores.js'

import { createEventDispatcher } from 'svelte'

const dispatch = createEventDispatcher()


/********************************************************************************************************************************************************
 usuario pediu para acessar dados de outro grupo , entrar em outro grupo
********************************************************************************************************************************************************/
const accessAnotherWorkgroup = () => {
    dispatch('dispatchCallMenuLogin')
}

/******************************************************************************************************************************************************** 
texto com explicacoes de como funcionam os workgroups (topo do menu)
********************************************************************************************************************************************************/

const prepareGroupsExplanation = () => {

let explanation = $Terms.workgroups_explanation

let _currentWorkgroupName = localStorage.getItem("rentacar_workgroup_name");
explanation = explanation.replaceAll('@name', _currentWorkgroupName)

return explanation
}


/******************************************************************************************************************************************************** 
usuario pediu para resetar dados do grupo
********************************************************************************************************************************************************/
const resetWorkgroupData = () => {
  dispatch('dispatchResetWorkgroupData')
}

/******************************************************************************************************************************************************** 
usuario pediu para escolher (aplicacao escolher aleatoriamente) outro nome de grupo 
********************************************************************************************************************************************************/
const randomChooseAnotherWorkgroup = () => {
  dispatch('dispatchRandomChooseAnother')
}




</script>



<div  class="flex flex-col font-Roboto absolute overflow-hidden w-[750px]  border-2 border-gray-400 rounded-lg bg-gray-100 p-2"  id='workgroupMenu'>

    <div  class="flex flex-col  bg-white relative  overflow-x-hidden overflow-y-visible h-auto rounded-lg" >

        <div class="flex flex-col w-full h-auto   pr-5 pl-3 " >

          <div class="flex flex-col gap-2 px-2 pb-4">

              <div class="WORKGROUP_EXPLANATION" style="text-align: justify;text-justify: inter-word;">
                {@html prepareGroupsExplanation()}
              </div>


              <button  id="btnWorkgroupReset" class="btnWORKGROUP_RESET_DATA" on:click={ () => {resetWorkgroupData()}} >{$Terms.reset_workgroup_data}</button>

              <button  id="btnChoseRandomWorkgroup" class="btnWORKGROUP_RANDOM_CHOOSE" on:click={ () => {randomChooseAnotherWorkgroup()}} >{$Terms.choose_random_workgroup}</button>

              <button  id="btnLoginExistingWorkgroup" class="btnLOG_IN_EXISTING_GROUP" on:click={ () => {accessAnotherWorkgroup()}} >{$Terms.login_existing_workgroup}</button>

          </div>

        </div>


    </div> 

</div> 


