

<script>
import { getContext, onMount } from 'svelte';

// prepara as funcoes de CRUD que serao usadas aqui, mas estao contidas em Datatable.svelte
const { performDeleteRecord, getFormPopulatedAndReady , invokeMetadataWindow, performSaveRecord  } = getContext('offerCrudOperationsToChildComponent');

// variaveis genericas, writeables
import { imagesUrl, staticImagesUrl, Terms, backendUrl, imagesStillLoading  } from '$lib/stores/stores.js'

import { createEventDispatcher } from 'svelte'
import { forceImgRefresh, slidingMessage , showFormErrors, encodeHTMLEntities   } from '$js/utils.js'

const dispatch = createEventDispatcher()

// formHttpMethodApply recebe qual operacao realizar com o form: POST (novo), PATCH (editar) ou DELETE
export let formHttpMethodApply


// representa o input type=file do logotipo
let fileManufacturerLogo  


/********************************************************************************************************************************************************
 se usuario clicar no [ X ] do formulario, dispara funcao (no datatable.svelte) que o fecha
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

/********************************************************************************************************************************************************
 usuario mudou logotipo do fabricnte usando botao de upload, atualiza <img> com o arquivo escolhido
*******************************************************************************************************************************************************/
const manufacturerLogoChanged = async () =>  { 

jq('#imgManufacturerLogo').attr('src', window.URL.createObjectURL( document.getElementById('fileManufacturerLogo').files[0] )) 

}




</script>

<!-- container do form edicao -->
<div  class="flex flex-col w-[70%] max-w-[1200px] overflow-hidden pt-8 font-Roboto"  id='recordForm'>

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
              {$Terms.new_manufacturer_form_title}  

            {:else if (formHttpMethodApply=='PATCH')  }
              {$Terms.edit_manufacturer_form_title}  

            {:else if (formHttpMethodApply=='DELETE')  }
              {$Terms.delete_manufacturer_form_title}  

            {/if}
          </div>

          <div class='flex flex-row '>
              <div id='divWINDOW_DRAG' _originalTooltip_={$Terms.msg_click_to_drag}  class='containsTooltip'>
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

          <div class="flex flex-row w-full gap-5">

            <!-- nome fabricante -->
            <div class="flex flex-col w-1/2">
              <div class="flex flex-col w-full"> 
                  <div>{$Terms.fieldname_generic_name }:</div>  
                  <div>
                    {#if formHttpMethodApply=='DELETE'}
                      <span class='span_formFieldValue' id='txtName'>&nbsp;</span>
                    {:else}
                      <input type="text" autocomplete="off" sequence="1"  maxlength='50' minlength='3' id="txtName" class='text_formFieldValue' on:focus={(e)=>{e.target.select()}} >
                    {/if}
                  </div>
              </div>
              <div class='noerrorTextbox' id="errorName">{$Terms.errormsg_fill_manufacturer_name}</div>
            </div>

            <!-- logotipo -->
            <div class="flex flex-col w-1/2 ">      

                <div class="w-[100px] h-[90px] mb-8 ml-20 ">
                  {#if $imagesStillLoading} 
                    <img id='imgManufacturerLogo'  alt=''  class='gifImageLoading' src='loading.gif'>  
                  {:else}
                    <img id='imgManufacturerLogo'  alt=''  class='imageLoadedInForm' src=''>  
                  {/if}
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


          {#if formHttpMethodApply!='DELETE'}
            <button  id="btnUPLOAD" class="btnUPLOAD" on:click={ () => {fileManufacturerLogo.click()}} >{$Terms.button_upload_logo}</button>
          {/if}

          {#if formHttpMethodApply=='DELETE'} 
            <button  id="btnDELETE" class="btnDELETE" on:click={()=>{performDeleteRecord()}} aria-hidden="true">{$Terms.button_delete}</button>
          {:else}
            <button  id="btnSAVE" class="btnSAVE" on:click={performSaveRecord} aria-hidden="true">{$Terms.button_save}</button>
          {/if}
        </div>

          <!-- botao upload, ficara oculto e sera acionado por btnUPLOAD -->
          <input type="file" accept="image/png" style="width: 0px; height: 0px; overflow: hidden;"  id='fileManufacturerLogo' 
                      on:change={ () => {manufacturerLogoChanged()}} bind:this={fileManufacturerLogo}  >

        <!-- mensagem rolante usada para avisos -->
        <div  id="slidingFormMessage" >
          &nbsp;
        </div>
 


    </div> 

</div> 


