

<script>
// importa a funcao 'performCrudOperation' do componente pai, Datatable.svelte
import { getContext } from 'svelte';
const { performCrudOperation } = getContext('offerCrudOperationsToChildComponent');

// variaveis genericas, writeables
import { Terms  } from '$lib/stores/stores.js'

import { onMount } from 'svelte';  
import { makeWindowDraggable, showOrHideTooltip } from '$js/utils.js'

import { createEventDispatcher } from 'svelte'

const dispatch = createEventDispatcher()

// 'html_message' recebe o texto que deve ser mostrado na janela
export let html_message
export let button_confirm_message

/********************************************************************************************************************************************************
 se usuario clicar no [ X ] da janela, dispara funcao (datatable.svelte) que o fecha
********************************************************************************************************************************************************/
const btnCloseWindowClicked = () => {
    dispatch('dispatchCloseWindow')
}

/********************************************************************************************************************************************************
 se usuario clicar no botao 'confirmar'
********************************************************************************************************************************************************/
const btnConfirmClicked = () => {
    dispatch('dispatchDeleteSelected')
}

/********************************************************************************************************************************************************
 faz com que a janela seja arrastavel a partir do topo da janela (parte azul escuro)
********************************************************************************************************************************************************/

onMount(() => {
  makeWindowDraggable('divWINDOW_TOP', 'confirmationWindow')
  showOrHideTooltip()
})


</script>

<!-- container do form edicao -->
<div  class="flex flex-col w-auto  overflow-hidden  "  id='confirmationWindow'>

    <!-- container da janela  -->

    <div  class="flex flex-col w-auto bg-white relative rounded-lg min-w-[700px]" >

        <!-- titulo e botao fechar -->
        <div id="divWINDOW_TOP" >
          
          <div id='divWINDOW_TITLE'>
            {$Terms.title_confirmation_window}
          </div>

          <div class='flex flex-row '>

              <div id='divWINDOW_DRAG' _originalTooltip_={$Terms.msg_click_to_drag}  class='containsTooltip'>
                &nbsp;
              </div>

              <div class='divWINDOW_BUTTON mr-2'  aria-hidden="true" >
                &nbsp;&nbsp;[ ? ]&nbsp;&nbsp;
              </div>

              <div class='divWINDOW_BUTTON mr-6'  on:click={btnCloseWindowClicked} aria-hidden="true" >
                &nbsp;&nbsp;[ X ]&nbsp;&nbsp;
              </div>
          </div>
        
        </div>

        <!-- parte principal da janela -->
        <div class="flex flex-col w-full h-auto  py-7 pr-5 pl-3 " >

          <div class="flex flex-row gap-4 w-auto ">

            <div class="flex flex-col w-full">
              {@html html_message}
            </div>

          </div>

        </div>

        <!-- botoes confirmar ou cancelar  -->
        <div class="flex flex-row w-auto justify-between px-6 border-t-[1px] border-t-gray-300 py-3 gap-20">
          <button  id="btnCLOSE" class="btnCANCEL" on:click={btnCloseWindowClicked} >{$Terms.button_cancel}</button>

          <button  id="btnCONFIRM" class="btnDELETE" on:click={btnConfirmClicked} aria-hidden="true">{button_confirm_message}</button>
        </div>


    </div> 

</div> 


