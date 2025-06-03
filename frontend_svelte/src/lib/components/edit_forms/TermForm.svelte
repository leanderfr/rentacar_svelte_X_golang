

<script>
import { getContext, onMount } from 'svelte';

// prepara as funcoes de CRUD que serao usadas aqui, mas estao contidas em Datatable.svelte
const { performDeleteRecord, getFormPopulatedAndReady, invokeMetadataWindow, performSaveRecord  } = getContext('offerCrudOperationsToChildComponent');

// variaveis genericas, writeables
import { Terms, staticImagesUrl  } from '$lib/stores/stores.js'

import { createEventDispatcher } from 'svelte'

const dispatch = createEventDispatcher()

// formHttpMethodApply recebe qual operacao realizar com o form: POST (novo), PATCH (editar) ou DELETE
export let formHttpMethodApply


/********************************************************************************************************************************************************
 se usuario clicar no [ X ] do formulario, dispara funcao (em Datatable.svelte) que o fecha
********************************************************************************************************************************************************/
const btnCloseFormClicked = () => {
    dispatch('dispatchCloseForm')
}

/********************************************************************************************************************************************************
 prepara form
********************************************************************************************************************************************************/
onMount(() => {
  getFormPopulatedAndReady()    //  em Datatable.svelte
})



</script>

<!-- container do form edicao -->
<div  class="flex flex-col w-[90%] max-w-[1200px] overflow-hidden pt-8 font-Roboto"  id='recordForm'>

    <!-- imagem oferecendo enter/setas, qdo form usado para visualizacao do registro (exemplo: delete), nao oferece as teclas -->
    {#if formHttpMethodApply!='DELETE'} 
      <div  class="flex flex-row w-full justify-end pr-[20%] z-20">
        <div class="absolute bg-center -mt-6 w-44 p-3 h-14 bg-contain bg-no-repeat cursor-pointer containsTooltip" 
            style="background-image: url('{$staticImagesUrl}/enter_arrows.png')  "
            _originalTooltip_={$Terms.use_arrow_keys} >
        </div>
      </div>
    {/if}

    <!-- container do form edicao 
    necessario ser 'relative' para que a slidingMessage ao final da div que é 'absolute' fique dentro da relative -->

    <div  class="flex flex-col w-full bg-white relative rounded-lg"  >

        <!-- titulo e botao fechar -->
        <div id='divWINDOW_TOP'>
          
          <div id='divWINDOW_TITLE'>
            {#if (formHttpMethodApply=='POST')  }
              {$Terms.new_language_form_title}  
            {:else if (formHttpMethodApply=='PATCH')  }
              {$Terms.edit_language_form_title}  
            {:else if (formHttpMethodApply=='DELETE')  }
              {$Terms.delete_language_form_title}  
            {/if}
          </div>

          <div class='flex flex-row '>
              <div class='bg-icon-div-draggable w-7 bg-transparent bg-no-repeat bg-contain bg-center mr-6 containsTooltip' _originalTooltip_={$Terms.msg_click_to_drag}  >
                &nbsp;
              </div>


              <div class='divWINDOW_BUTTON mr-2'  aria-hidden="true" >
                &nbsp;&nbsp;[ ? ]&nbsp;&nbsp;
              </div>

              <div class='divWINDOW_BUTTON mr-6'  on:click={btnCloseFormClicked} aria-hidden="true" >
                &nbsp;&nbsp;[ X ]&nbsp;&nbsp;
              </div>
          </div>
        
        </div>

        <!-- campos  do formulario -->
        <div class="flex flex-col w-full h-auto  px-4 mt-6" >

          <div class="flex flex-col w-full">

            <div class="flex flex-row w-full">
              <div class='w-[15%] text-right pr-5 pt-2'>{$Terms.fieldname_language_item}:</div>
              <div class="flex flex-col w-[90%]">  
                      {#if formHttpMethodApply=='DELETE'}
                        <span class='span_formFieldValue' id='txtItem'>&nbsp;</span>
                      {:else}
                        <input type="text" autocomplete="off" sequence="1"  id="txtItem" maxlength='200' minlength='5' class='text_formFieldValue' on:focus={(e)=>{e.target.select()}} >
                      {/if}
                    <div class='noerrorTextbox' id="errorItem">{$Terms.errormsg_fill_language_item}</div>


              </div>
            </div>

            <div class="flex flex-row w-full">
              <div class="w-[15%] text-right pr-5 pt-2">{@html $Terms.fieldname_language_portuguese}:</div>
              <div class="flex flex-col w-[90%]">
                    {#if formHttpMethodApply=='DELETE'}
                      <div class='span_formFieldValue' id='txtPortuguese' style='height:80px;overflow-y:auto;border:solid 1px lightgray;padding:5px;'>&nbsp;</div>
                    {:else}            
                      <textarea class='h-[80px] w-full' autocomplete="off" sequence="2" id='txtPortuguese' maxlength='1200' minlength='2' on:focus={(e)=>{e.target.select()}} ></textarea>
                    {/if}
                    <div class='noerrorTextbox' id="errorPortuguese">{$Terms.errormsg_fill_language_portuguese}</div>
              </div>  
            </div>

            <div class="flex flex-row w-full">
              <div class='w-[15%] text-right pr-5 pt-2'>{@html $Terms.fieldname_language_english}:</div>
              <div class="flex flex-col w-[90%]">
                    {#if formHttpMethodApply=='DELETE'}
                      <div class='span_formFieldValue' id='txtEnglish' style='height:80px;overflow-y:auto;border:solid 1px lightgray;padding:5px;'>&nbsp;</div>
                    {:else}            
                      <textarea class='h-[80px] w-full' autocomplete="off" sequence="3" id='txtEnglish' maxlength='1200' minlength='2' on:focus={(e)=>{e.target.select()}} ></textarea>
                    {/if}
                    <div class='noerrorTextbox' id="errorEnglish">{$Terms.errormsg_fill_language_english}</div>
              </div>
            </div>

            <div class="flex flex-row w-full">
              <div class='w-[15%] text-right pr-5 pt-2'>{@html $Terms.fieldname_language_description}:</div>
              <div class="flex flex-col w-[90%]">
                    {#if formHttpMethodApply=='DELETE'}
                      <div class='span_formFieldValue' id='txtDescription' style='height:80px;overflow-y:auto;border:solid 1px lightgray;padding:5px;'>&nbsp;</div>
                      <div class='deleteWarning ' >{$Terms.warning_delete_record_form}</div>
                    {:else}            
                      <textarea class='h-[80px] w-full' sequence="4" id='txtDescription' maxlength='200' minlength='10' on:focus={(e)=>{e.target.select()}} ></textarea>
                      <div class='noerrorTextbox' id="errorDescription">{$Terms.errormsg_fill_language_description}</div>
                    {/if}
                    <div class='noerrorTextbox' id="errorEnglish">{$Terms.errormsg_fill_language_english}</div>
              </div>
            </div>


          </div>

        </div>

        <!-- botoes salvar/sair -->
        <div class="flex flex-row w-full justify-between px-6 border-t-[1px] border-t-gray-300 py-2">
          <button  id="btnCLOSE" class="btnCANCEL" on:click={btnCloseFormClicked} >{$Terms.button_cancel}</button>

          {#if formHttpMethodApply!='POST'}
            <button  id="btnMETADATA" class="btnMETADATA" on:click={ () => {invokeMetadataWindow()}} aria-hidden="true">{$Terms.button_metadata}</button>
          {/if}


          {#if formHttpMethodApply=='DELETE'} 
            <button  id="btnDELETE" class="btnDELETE" on:click={()=>{performDeleteRecord()}} aria-hidden="true">{$Terms.button_delete}</button>
          {:else}
            <button  id="btnSAVE" class="btnSAVE" on:click={performSaveRecord} aria-hidden="true">{$Terms.button_save}</button>
          {/if}
        </div>

        <!-- mensagem rolante usada para avisos -->
        <div  id="slidingFormMessage" >
          &nbsp;
        </div>
 

    </div> 

</div> 


