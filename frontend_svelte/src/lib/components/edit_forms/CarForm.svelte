
<script>
import { getContext, onMount } from 'svelte';

// prepara as funcoes de CRUD que serao usadas aqui, mas estao contidas em Datatable.svelte
const { performSaveRecord, performDeleteRecord, getFormPopulatedAndReady, invokeMetadataWindow  } = getContext('offerCrudOperationsToChildComponent');
  
// variaveis genericas, writeables
import { imagesUrl, staticImagesUrl, Terms, backendUrl , manufacturersAutocomplete, imagesStillLoading  } from '$lib/stores/stores.js'

import { createEventDispatcher } from 'svelte'
import { forceImgRefresh, slidingMessage , showFormErrors, encodeHTMLEntities, selectedCountry, getAutocompleteItemId   } from '$js/utils.js'

const dispatch = createEventDispatcher()

// formHttpMethodApply recebe qual operacao realizar com o form: POST (novo), PATCH (editar) ou DELETE
export let formHttpMethodApply

// fileCarImage => representa o input type=file da foto  do veiculo
let fileCarImage

// responde (bind) pelo checkbox tipo de transmissao (automatica ou manual)
// transmissiona manual= checkbox true,  automatica= checkbox false
let transmissionManual


/********************************************************************************************************************************************************
 se usuario clicar no [ X ] do formulario, dispara funcao (no datatable.svelte) que o fecha
********************************************************************************************************************************************************/
const btnCloseFormClicked = () => {
    dispatch('dispatchCloseForm')
}

/********************************************************************************************************************************************************
obtem situacao do checkbox do tipo de transmissao do carro
********************************************************************************************************************************************************/
const transmissionChanged = () => {
// true= manual
// false= cambio automatico

// coloca o tipo de cambio escolhido atualmente em cor mais forte
setTimeout(() => {
  transmissionManual = jq('#chkTransmissionManual').prop('checked')  
}, 1000);

}


/********************************************************************************************************************************************************
 prepara form
********************************************************************************************************************************************************/
onMount(() => {


  getFormPopulatedAndReady()    //  em Datatable.svelte

})

// se mudar tipo de transmissao (variavel bind 'transmissionManual') , atualiza visualmente o tipo atual
$: transmissionManual, transmissionChanged()

/********************************************************************************************************************************************************
 usuario mudou foto do carro usando botao de upload, atualiza <img> com o arquivo escolhido
*******************************************************************************************************************************************************/
const carImageChanged = async () =>  { 

jq('#imgCarImage').attr('src', window.URL.createObjectURL( document.getElementById('fileCarImage').files[0] )) 

}

</script>

<!-- container do form edicao -->
<div  class="flex flex-col w-[95%] max-w-[1400px] overflow-hidden pt-8 font-roboto"  id='recordForm'>

    <!-- imagem oferecendo enter/setas, qdo form usado para visualizacao do registro (exemplo: DELETE), nao oferece as teclas -->
    {#if formHttpMethodApply!='DELETE'}   
      <div  class="flex flex-row w-full justify-end pr-[20%] z-20">
        <div class="absolute bg-center -mt-6 w-44 p-3 h-14 bg-contain bg-no-repeat  cursor-pointer containsTooltip" 
          style="background-image: url('{$staticImagesUrl}/enter_arrows.png')"
          _originalTooltip_={$Terms.use_arrow_keys} >
        </div>
      </div>
    {/if}

    <!-- container do form edicao 
    necessario ser 'relative' para que a slidingMessage ao final da div que é 'absolute' fique dentro da relative -->

    <div  class="flex flex-col w-full bg-white relative rounded-lg"  >

        <!-- titulo e botao fechar -->
        <div id="divWINDOW_TOP" >
          
          <div id='divWINDOW_TITLE'>
            {#if (formHttpMethodApply=='POST')  }
              {$Terms.new_car_form_title}  

            {:else if (formHttpMethodApply=='PATCH')  }
              {$Terms.edit_car_form_title}  

            {:else if (formHttpMethodApply=='DELETE')  }
              {$Terms.delete_car_form_title}  
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

          <div class="flex flex-row w-full gap-9">

            <!-- **************************************************************************************
              imagem do carro
            ************************************************************************************** -->


            <div class="flex  w-[20%]   h-[200px] justify-center  items-center align-middle ">
                {#if $imagesStillLoading} 
                  <img id='imgCarImage'  alt=''  class='gifImageLoading' src='loading.gif'>  
                {:else}
                  <img id='imgCarImage'  alt=''  class='imageLoadedInForm' src=''>  
                {/if}
            </div>

            <!-- **************************************************************************************
              campos texto 1a linha 
            ************************************************************************************** -->
            <div class="flex flex-col w-[80%] h-[230px]">

              <div class="flex flex-row w-full gap-5 ">

                  <div class="flex flex-col w-[25%]">
                        <div>{$Terms.fieldname_country }:</div>
                        <div id='divCountryFlag'></div>
                  </div>

                  <div class="flex flex-col w-[25%]">
                        <div>{$Terms.fieldname_generic_name }:</div>
                        <div>
                          {#if formHttpMethodApply=='DELETE'}
                            <span class='span_formFieldValue' id='txtName'>&nbsp;</span>
                          {:else}
                            <input type="text" autocomplete="off" sequence="1" maxlength='50' minlength='2'  id="txtName" class='text_formFieldValue' on:focus={(e)=>{e.target.select()}} >

                          {/if}
                        </div>
                    <div class='noerrorTextbox' id="errorName">{$Terms.errormsg_fill_car_model}</div>
                  </div>


                  <div class="flex flex-col w-[25%]">
                        <div>{$Terms.fieldname_year_manufacture }:</div>
                        <div>
                          {#if formHttpMethodApply=='DELETE'}
                            <span class='span_formFieldValue'   id='txtYear'>&nbsp;</span>
                          {:else}
                            <input type="text" autocomplete="off" sequence="2" maxlength='4' minlength='4' id="txtYear" class='text_formFieldValue' on:focus={(e)=>{e.target.select()}} >
                          {/if}
                        </div>
                    <div class='noerrorTextbox' id="errorYear">{$Terms.errormsg_fill_car_year}</div>
                  </div>

                  <div class="flex flex-col w-[25%]">
                        <div>{$Terms.fieldname_car_manufacturer }:</div>
                        <div>
                          {#if formHttpMethodApply=='DELETE'}
                            <span class='span_formFieldValue'  id='txtManufacturerName'>&nbsp;</span>
                          {:else}
                            <input type="text" autocomplete="off" sequence="3" maxlength='50' minlength='2'  id="txtManufacturerName" class='text_formFieldValue' on:focus={(e)=>{e.target.select()}} >
                          {/if}
                        </div>
                    <div class='noerrorTextbox' id="errorManufacturerName">{$Terms.errormsg_fill_manufacturer_name}</div>
                  </div>
              </div>

              <!-- **************************************************************************************
                campos texto 2a linha 
              ************************************************************************************** -->

              <div class="flex flex-row w-full gap-5 ">

                  <div class="flex flex-col w-[25%]">
                        <div>{$Terms.odometer }:</div>
                        <div>
                          {#if formHttpMethodApply=='DELETE'}
                            <span class='span_formFieldValue'   id='txtOdometer'>&nbsp;</span>
                          {:else}
                            <input type="text" autocomplete="off"  sequence="4"  maxlength='7' minlength='2'  id="txtOdometer" class='text_formFieldValue' on:focus={(e)=>{e.target.select()}} >
                          {/if}
                        </div>
                    <div class='noerrorTextbox' id="errorOdometer">{$Terms.errormsg_fill_odometer}</div>
                  </div>

                  <div class="flex flex-col w-[25%]">
                        <div>{$Terms.car_consume }:</div>
                        <div>
                          {#if formHttpMethodApply=='DELETE'}
                            <span class='span_formFieldValue' id='txtMpg'>&nbsp;</span>
                          {:else} 
                            <input type="text" autocomplete="off"  maxlength='5' minlength='1'  sequence="5"  id="txtMpg" class='text_formFieldValue' on:focus={(e)=>{e.target.select()}} >
                          {/if}
                        </div>
                    <div class='noerrorTextbox' id="errorMpg">{$Terms.errormsg_fill_car_consume}</div>
                  </div>


                  <div class="flex flex-col w-[25%]">
                        <div>{$Terms.cylinders }:</div>
                        <div>
                          {#if formHttpMethodApply=='DELETE'}
                            <span class='span_formFieldValue' id='txtCylinders'>&nbsp;</span>
                          {:else}
                            <input type="text" autocomplete="off"  maxlength='1' minlength='1'  sequence="6"  id="txtCylinders" class='text_formFieldValue' on:focus={(e)=>{e.target.select()}} >
                          {/if}
                        </div>
                    <div class='noerrorTextbox' id="errorCylinders">{$Terms.errormsg_fill_cylinders}</div>
                  </div>

                  <div class="flex flex-col w-[25%]">
                        <div>{$Terms.fieldname_transmision }:</div>
                        {#if formHttpMethodApply=='DELETE'}
                          <span class='span_formFieldValue' id='txtTransmission'>&nbsp;</span>

                        {:else}

                          <div class="flex flex-row gap-4 ">
                              <div class='p-1 cursor-pointer' on:click={ () => {transmissionManual = false}} aria-hidden="true" 
                                class:text-gray-300={ transmissionManual }   class:text-black={ ! transmissionManual } class:font-bold={ ! transmissionManual }   >     
                                {$Terms.automatic_transmission_short}
                              </div>

                              <div class='p-[2px]'>     
                                <label for="chkTransmissionManual" class="switch_field_value"   >
                                  <input id="chkTransmissionManual" type="checkbox"  bind:checked={transmissionManual}  >
                                  <span class="slider_field_value round"></span>  
                                </label>
                              </div>

                              <div class='p-1 cursor-pointer' on:click={ () => {transmissionManual = true} } aria-hidden="true" 
                                class:text-gray-300={ ! transmissionManual }   class:text-black={ transmissionManual }  class:font-bold={ transmissionManual } >     
                                {$Terms.manual_transmission_short}
                              </div>
                          </div>
                        {/if}
                  </div>
              </div>

              <!-- **************************************************************************************
                campos texto 3a linha 
              ************************************************************************************** -->

              <div class="flex flex-row w-full gap-5 ">

                  <div class="flex flex-col w-[25%]">
                        <div>{$Terms.car_power }:</div>
                        <div>
                          {#if formHttpMethodApply=='DELETE'}
                            <span class='span_formFieldValue' id='txtHp'>&nbsp;</span>
                          {:else}
                            <input type="text" autocomplete="off"   maxlength='3' minlength='2'  sequence="7" id="txtHp" class='text_formFieldValue' on:focus={(e)=>{e.target.select()}} >
                          {/if}
                        </div>
                    <div class='noerrorTextbox' id="errorHp">{$Terms.errormsg_fill_car_power}</div>
                  </div>

                  <div class="flex flex-col w-[25%]">
                        <div>{$Terms.doors }:</div>
                        <div>
                          {#if formHttpMethodApply=='DELETE'}
                            <span class='span_formFieldValue' id='txtDoors'>&nbsp;</span>
                          {:else}
                            <input type="text" autocomplete="off"  maxlength='1' minlength='1'  sequence="8"  id="txtDoors" class='text_formFieldValue' on:focus={(e)=>{e.target.select()}} >
                          {/if}
                        </div>
                    <div class='noerrorTextbox' id="errorDoors">{$Terms.errormsg_fill_doors}</div>
                  </div>


                  <div class="flex flex-col w-[25%]">
                        <div>{$Terms.cylinder_capacity }:</div>
                        <div>
                          {#if formHttpMethodApply=='DELETE'}
                            <span class='span_formFieldValue' id='txtCc'>&nbsp;</span>
                          {:else}
                            <input type="text" autocomplete="off"  maxlength='5' minlength='2'  sequence="9"  id="txtCc" class='text_formFieldValue' on:focus={(e)=>{e.target.select()}} >
                          {/if}
                        </div>
                    <div class='noerrorTextbox' id="errorCc">{$Terms.errormsg_fill_cc}</div>
                  </div>

                  <div class="flex flex-col w-[25%]">
                        <div>{$Terms.fieldname_rental_price }:</div>
                        <div>
                          {#if formHttpMethodApply=='DELETE'}
                            <span class='span_formFieldValue' id='txtRentalPrice'>&nbsp;</span>
                          {:else}
                            <input type="text" autocomplete="off" sequence="10" maxlength="8" minlength='4'  id="txtRentalPrice" class='text_formFieldValue' on:focus={(e)=>{e.target.select()}} >
                          {/if}  
                        </div>
                    <div class='noerrorTextbox' id="errorRentalPrice">{$Terms.errormsg_fill_rental_price}</div>
                  </div>
              </div>


            </div>



          </div>

          <!-- alerta sobre exclusao do registro -->
          {#if formHttpMethodApply=='DELETE'}
            <div class="flex flex-row w-full gap-5 ">
              <div class='deleteWarning' >{$Terms.warning_delete_record_form}</div>
            </div>
          {/if}



        </div>

        <!-- botoes salvar/sair -->
        <div class="flex flex-row w-full justify-between px-6 border-t-[1px] border-t-gray-300 py-2">
          <button  id="btnCLOSE" class="btnCANCEL" on:click={btnCloseFormClicked} >{$Terms.button_cancel}</button>

          {#if formHttpMethodApply!='POST'}
            <button  id="btnMETADATA" class="btnMETADATA" on:click={ () => {invokeMetadataWindow()}} aria-hidden="true">{$Terms.button_metadata}</button>
          {/if}


          {#if formHttpMethodApply!='DELETE'}
            <button  id="btnUPLOAD" class="btnUPLOAD" on:click={ () => {fileCarImage.click()}} >{$Terms.button_upload_image}</button>
          {/if}

          {#if formHttpMethodApply=='DELETE'} 
            <button  id="btnDELETE" class="btnDELETE" on:click={ () => {performDeleteRecord()}} aria-hidden="true">{$Terms.button_delete}</button>
          {:else}
            <button  id="btnSAVE" class="btnSAVE" aria-hidden="true"
                    on:click={()=>{performSaveRecord()}} >{$Terms.button_save}</button>
          {/if}
        </div>

          <!-- botao upload, ficara oculto e sera acionado por btnUPLOAD -->
          <input type="file" accept="image/png" style="width: 0px; height: 0px; overflow: hidden;"  id='fileCarImage' 
                        on:change={ () => {carImageChanged()}} bind:this={fileCarImage}  >

        <!-- mensagem rolante usada para avisos -->
        <div  id="slidingFormMessage" >  
          &nbsp;
        </div>
 


    </div> 

</div> 


