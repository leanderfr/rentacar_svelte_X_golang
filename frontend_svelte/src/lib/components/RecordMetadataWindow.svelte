

<script>

// variaveis genericas, writeables
import { Terms  } from '$lib/stores/stores.js'

import { onMount, getContext } from 'svelte';  
import { makeWindowDraggable, showOrHideTooltip } from '$js/utils.js'

import { createEventDispatcher } from 'svelte'

// prepara a funcao que preenche a janela atual com metadados do registro
const { fillRecordMetadataWindow  } = getContext('offerCrudOperationsToChildComponent');

// recebe o ID do registro cujos metadados serao visualizados
export let currentFormRecordId = 0 

// contem o nome da tabela cujos metadados estao sendo exibidos
export let metadataTableName

const dispatch = createEventDispatcher()

/********************************************************************************************************************************************************
 se usuario clicar no [ X ] da janela, dispara funcao (datatable.svelte) que o fecha
********************************************************************************************************************************************************/
const btnCloseMetadataWindowClicked = () => {
    dispatch('dispatchCloseForm')
}


/********************************************************************************************************************************************************
 faz com que a janela seja arrastavel a partir do topo da janela (parte azul escuro)
********************************************************************************************************************************************************/

onMount(() => {
  makeWindowDraggable('divWINDOW_TOP', 'metadataWindow')
  showOrHideTooltip()

  setTimeout(() => {
    fillRecordMetadataWindow(currentFormRecordId)
  }, 100);
})


</script>




<!-- container do form edicao -->
<div  class="flex flex-col w-auto  overflow-hidden  "  id='metadataWindow'>

    <!-- container da janela  -->

    <div  class="flex flex-col w-auto bg-white relative rounded-lg min-w-[700px]" >

        <!-- titulo e botao fechar -->
        <div id="divWINDOW_TOP" >
          
          <div id='divWINDOW_TITLE'>
            {$Terms.title_metadata_window}
          </div>

          <div class='flex flex-row '>

              <div id='divWINDOW_DRAG' _originalTooltip_={$Terms.msg_click_to_drag}  class='containsTooltip'>
                &nbsp;
              </div>

              <div class='divWINDOW_BUTTON mr-2'  aria-hidden="true" >
                &nbsp;&nbsp;[ ? ]&nbsp;&nbsp;
              </div>

              <div class='divWINDOW_BUTTON mr-6'  on:click={btnCloseMetadataWindowClicked} aria-hidden="true" >
                &nbsp;&nbsp;[ X ]&nbsp;&nbsp;
              </div>
          </div>
        
        </div>

        <!-- parte principal da janela -->
        <div class="flex flex-col w-full h-auto  py-7 pr-5 pl-3 " >

          <div class="flex flex-col gap-4 w-auto ">

            <div class="flex flex-row w-full ">
              <div class="w-[30%] flex justify-end pr-4">{$Terms.table_fieldname}:</div>
              <div class="w-[70%] ">{metadataTableName}</div>
            </div>

            <div class="flex flex-row w-full border-b-2 mb-1 pb-3">
              <div class="w-[30%] flex justify-end  pr-4">{$Terms.record_id}</div>
              <div class="w-[70%] ">{currentFormRecordId}</div> 
            </div>

            <div class="flex flex-row w-full">
              <div class="w-[30%] flex justify-end  pr-4">{$Terms.created_at}</div>
              <div class="w-[70%]" id='info_created_at'></div> 
            </div>

            <div class="flex flex-row w-full">
              <div class="w-[30%] flex justify-end  pr-4">{$Terms.last_update}</div>
              <div class="w-[70%]" id='info_updated_at'></div> 
            </div>

            <div class="flex flex-row w-full">
              <div class="w-[30%] flex justify-end  pr-4">{$Terms.creator_workgroup}</div>
              <div class="w-[70%]" id='info_workgroup'></div> 
            </div>

          </div>

        </div>

        <!-- botoes confirmar ou cancelar  -->
        <div class="flex flex-row w-auto justify-between px-6 border-t-[1px] border-t-gray-300 py-3 gap-20">
          <button  id="btnCLOSE_METADATA_WINDOW" class="btnCANCEL" on:click={btnCloseMetadataWindowClicked} >{$Terms.button_cancel}</button>
        </div>


    </div> 

</div> 


