

<script>

// variaveis genericas, writeables
import { Terms, backendUrlGolang } from '$lib/stores/stores.js'
import { onMount, createEventDispatcher } from 'svelte'

import { logAPI } from '$js/utils.js'

const dispatch = createEventDispatcher()

// contabiliza qtde de mudanças feitas na base ja
let databaseChangesMadeAmount = 'reading'

/********************************************************************************************************************************************************
 se usuario clicar no [ X ] do formulario, despacha comando para fechar janela
********************************************************************************************************************************************************/
const btnCloseFormClicked = () => {
    dispatch('dispatchCloseWorkgroupReset')
}

/********************************************************************************************************************************************************
 usuario confirmou reset dos dados do grupo
********************************************************************************************************************************************************/
const performWorkgroupDataReset = () => {
    dispatch('dispatchPerformGroupDataReset')
}


/********************************************************************************************************************************************************
********************************************************************************************************************************************************/
onMount(async () => {

// necessario abrir evento assincrono para exibir div ajax loading, caso contrario navegador nao atualiza a tela
setTimeout(() => {
  showLoadingGif();   // mostra animacao 'processando..' mesmo sem tela ter sido renderizada
}, 1);

// obtem qtde de modificacoes que o grupo ja fez na base de dados
let _currentWorkgroupName = localStorage.getItem("rentacar_workgroup_name");

// manipulacao de grupos e notificacoes sao feitas somente pelo backend golang
let _route_ = `${$backendUrlGolang}/${_currentWorkgroupName}/report/`
logAPI('GET', _route_)

await fetch(_route_, { method: 'GET' })

.then((response) => {
  if (!response.ok) {
    throw new Error(`Access workgroup report Err fatal= ${response.status}`);
  }
  return response.text();
})

.then((_databaseChangesMadeAmount) => {
  setTimeout( () => {hideLoadingGif();}, 1)

  databaseChangesMadeAmount = _databaseChangesMadeAmount.split('|')[1]      //  _databaseChangesMadeAmount=  __success__|qtde alteracoes feitas
  databaseChangesMadeAmount = databaseChangesMadeAmount=='0' ? 'none' : databaseChangesMadeAmount

})

})

/********************************************************************************************************************************************************
 mensagem com a qtde de alteracoes feitas na base
********************************************************************************************************************************************************/
const databaseChangesAmount = () => {

let msg = $Terms.note_about_data_reset_confirm
msg = msg.replaceAll('@amount', databaseChangesMadeAmount)

return msg
}


</script>

<div  class="flex flex-col w-auto  overflow-hidden bg-white p-[4px] rounded-xl font-Roboto "  id='workgroupDataReset'  >

    <div  class="flex flex-col w-[900px]  border-2  relative rounded-lg " style='background-color: #ffe6e6;border-color:red' >

        <div class="flex flex-col w-full h-auto pt-6 pb-10 pr-5 pl-3 " >

          <div class="text-xl pb-8" >
            {$Terms.to_reset_workgroup_data}
          </div>

          <div class="flex flex-row justify-between pb-14 gap-4" >
            <div style='text-align: justify;text-justify: inter-word;height:50px;' >
                {#if databaseChangesMadeAmount=='reading'} 
                  &nbsp;&nbsp;&nbsp;

                {:else if databaseChangesMadeAmount=='none'} 
                  {$Terms.note_about_data_reset_no_changes}

                {:else} 
                  {@html databaseChangesAmount() }


                {/if}
            </div>
          </div>


        </div>

        <!-- botoes confirmar ou cancelar  -->
        <div class="flex flex-row w-auto justify-between px-6 border-t-[1px] border-t-gray-300 py-3 gap-20">
          <button  id="btnCLOSE_RESET_WINDOW" class="btnCANCEL" on:click={btnCloseFormClicked} >{$Terms.button_cancel}</button>

          {#if databaseChangesMadeAmount!='none' && databaseChangesMadeAmount!='reading' } 
            <button  id="btnCONTINUE" class="btnCONTINUE" aria-hidden="true" on:click={()=>{performWorkgroupDataReset()}} >{$Terms.button_yes_continue}</button>
          {/if}
        </div>


    </div> 

</div>

