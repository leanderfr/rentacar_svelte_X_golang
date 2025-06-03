
<script>

// variaveis genericas, writeables
import { Terms, sourceCodeViewer  } from '$lib/stores/stores.js'
import { createEventDispatcher } from 'svelte'
import { setBackendUrl } from '$js/utils.js'


let currentChoice = ''


const dispatch = createEventDispatcher()


/************************************************************************************************************************************************************
// memoriza qual backend deve ser usado a partir de agora
************************************************************************************************************************************************************/
const proceedBackendSet = () => {

// memoriza (localstorage) e muda variavel global que informa o backend
setBackendUrl(currentChoice)  

dispatch('dispatchProceedBackendChange')   // retorna a Page.svelte com o backend escolhido
}


/************************************************************************************************************************************************************
// abre visualizador de cod fonte
************************************************************************************************************************************************************/
const viewSourceCode = () => {

// url exemplos:
// http://leanderdeveloper.store/php
// http://leanderdeveloper.store/golang

let url = $sourceCodeViewer + currentChoice
window.open(url, '_blank').focus();

}
</script>


<div class="flex flex-col w-[1000px] bg-white relative rounded-lg  h-48 border-gray-700 border-[1px] justify-around items-center px-7 py-5 "  > 

    <div class="flex flex-row w-full bg-white gap-4"  > 

            <div  class='divAsButton' aria-hidden="true" on:click={() => { currentChoice = 'frontend'}} 
                              class:divAsButtonClicked={ currentChoice == 'frontend' }         >     
              <div class='text-lg'>Front End</div>
              <div style='display:flex;flex-direction: row;'>
                <div><img src="nodejs_symbol.svg" style="height:70px;width:70px;" alt=''></div>           
                <div><img src="html_symbol.svg" style="height:60px;width:60px;" alt=''></div>           
                <div><img src="css_symbol.svg" style="height:60px;width:60px;" alt=''></div>           
                <div><img src="javascript_symbol.svg" style="height:60px;width:60px;" alt=''></div>           
                <div><img src="jquery_symbol.svg" style="height:70px;width:70px;padding-left:5px" alt=''></div>           
                <div><img src="tailwind_symbol.svg" style="height:60px;width:60px;" alt=''></div>           
                <div><img src="svelte_symbol.svg" style="height:60px;width:60px;" alt=''></div>           

              </div>
            </div>


            <div  class='divAsButton' aria-hidden="true" on:click={() => { currentChoice = 'golang'}} 
                  class:divAsButtonClicked={ currentChoice == 'golang' }         >     
              <div class='text-lg'>Backend Golang</div>
              <div ><img src="golang_symbol.svg" style='width:60px;height:60px;' alt=''> </div>
            </div>

            <div  class='divAsButton' aria-hidden="true" on:click={() => { currentChoice = 'php'}} 
                  class:divAsButtonClicked={ currentChoice == 'php' }         >     
              <div class='text-lg'>Backend PHP</div>
              <div ><img src="php_symbol.svg" style='width:60px;height:60px;' alt=''> </div>
            </div>


      </div>

      <div class='w-full flex justify-between pr-3 pt-7 '  >     
        <button class:visible={currentChoice != ''} class:invisible={currentChoice == ''}  class="btnMETADATA " aria-hidden="true" 
          style='min-width:300px;'  on:click={()=>{viewSourceCode()}} >{$Terms.view_source_code}</button>

        <button class:visible={currentChoice != '' && currentChoice != 'frontend'} 
          class:invisible={currentChoice == '' || currentChoice == 'frontend'}  class="btnFLEXIBLE_OK " aria-hidden="true" 
          style='min-width:300px;padding-right:10px'  on:click={()=>{proceedBackendSet()}} >{$Terms.confirm_backend_change}</button>
      </div>

</div>


