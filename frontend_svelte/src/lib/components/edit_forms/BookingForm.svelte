
<script>
import { getContext, onMount } from 'svelte';

// variaveis genericas, writeables
import { Terms , elementCSSChangedTemporarily, elementIdCSSChangedTemporarily, imagesStillLoading, imagesUrl } from '$lib/stores/stores.js'

import { createEventDispatcher } from 'svelte'
import { makeWindowDraggable, selectedCountry , showOrHideTooltip  } from '$js/utils.js'

// prepara as funcoes de CRUD que serao usadas aqui, mas estao contidas em Datatable.svelte
const { performSaveBookingRecord, performDeleteBookingRecord, getBookingFormPopulatedAndReady, invokeBookingMetadataWindow  } = getContext('offerCrudOperationsToChildComponent');

import Calendar from '$lib/components/Calendar.svelte'

const dispatch = createEventDispatcher()


// formHttpMethodApply recebe qual operacao realizar com o form de reserva: POST (novo), PATCH (editar) ou DELETE
export let formHttpMethodApply

// exibe/oculta pequeno calendario que auxilia no preenchimento das datas
let showDatePicker

// contem qual input type=text chamou calendario para escolha da data
let whichInputTextCalled



/********************************************************************************************************************************************************
 se usuario clicar no [ X ] do formulario, dispara funcao (em Bookings.svelte) que o fecha
********************************************************************************************************************************************************/
const closeBookingModalWindow = () => {
    dispatch('dispatchCloseBookingForm')
}


/********************************************************************************************************************************************************
 prepara form
********************************************************************************************************************************************************/
onMount(() => {
  getBookingFormPopulatedAndReady()
})



/*****************************************************************************************************************************************************************
// funcao acionada qdo usuario clica na div '#backDrop' ou clica em um dos botoes de fechar o date picker
*****************************************************************************************************************************************************************/
export const closeDatePicker = () => {
  showDatePicker = false    // fecha janela modal 

  // bg-color do input type=text que recebe a data foi alterado para ficar igual ao bg-color do date picker (Calendar.svelte)
  // devolvendo todo CSS anterior à mudança do bg-color do input type=text
  jq( `${$elementIdCSSChangedTemporarily}`).prop('style', $elementCSSChangedTemporarily)
  jq( `${$elementIdCSSChangedTemporarily}`).focus()
}


/************************************************************************************************************************************************************
 joga data escolhida no 'date picker' (<Calendar />)  para o input type=text de data que invocou o 'Date Picker'
************************************************************************************************************************************************************/

const putChosenDateInInputText = (custom_event) => { 

  let clickedDate = new Date(custom_event.detail.clickedDate)   // converte data recebida em milisegundos para data formato 

  let _clickedDate_ = ''

  // converte data para formato do país atualmente escolhido
  let options = {month: '2-digit', day: '2-digit', year:'2-digit'} 
  if ( selectedCountry() == 'usa') 
    _clickedDate_ = clickedDate.toLocaleDateString('en-us', options )    // mm/dd/yy
  else 
    _clickedDate_ = clickedDate.toLocaleDateString('pt-br', options )    // dd/mm/yy

  // id do input type=text foi memorizado para devolver CSS temporariamente alterado
  // usando o mesmo ID para popular a data escolhida
  jq( `${$elementIdCSSChangedTemporarily}`).val(_clickedDate_);   // converte milisegundos para Date


}

let whichComponentSvelteCalledCalendar  = 'BookingForm'

</script>


<!-- date picker é o mesmo componente <Calendar />  mas menor e posicionado logo abaixo dos input type=text responsaveis pelas datas  -->
{#if showDatePicker}       

  <!-- backDrop é uma div fumê que ficara por trás do calendario
  backDrop impede que o usuario clique em elementos que estão atrás da janela modal
  qdo backDrop for clicada, fecha a janela modal <Calendar /> 

  backDrop nao terá bg-color escura porque nao é uma janela modal, é uma janela auxiliar de preenchimento de data  
  -->

  <div id='backDrop' class='w-full h-full  absolute flex items-center justify-center left-0 top-0 z-10 bg-[rgba(0,0,0,0.0)]' on:click|self={closeDatePicker} aria-hidden="true"  >  

      <!-- 

      a div que contem o calendario, nesse caso, nao vai ficar centrada na tela, 
        - ela vai possuir um ID (BookingFormCalendarContainer), 
        - vai ser 'position:absolute'
        - vai iniciar como visibility:invisible

      isso porque apos o componente <Calendar /> ser montado, o jquery vai usar a funcao 'jquery-ui.position()' para posicionar o calendario EMBAIXO do input type=text 
      que recebera a data escolhida no calendario 

     -->  

     <div class="w-[500px]  h-auto flex align-top text-lg items-start invisible justify-start  flex-col  bg-white  border-2 border-gray-300 absolute" id='BookingFormCalendarContainer' > 
          <Calendar  {whichInputTextCalled}  {whichComponentSvelteCalledCalendar}  on:dispatchCloseCalendar={closeDatePicker}  on:dispatchPutChosenDateInInputText={putChosenDateInInputText}   /> 
     </div>

  </div>

{/if}



<!-- container do form edicao -->
<div  class="flex flex-col w-[95%] max-w-[1300px] overflow-hidden pt-8 "  id='bookingRecordForm'>

    <!-- imagem oferecendo enter/setas, qdo form usado para visualizacao do registro (exemplo: delete), nao mostra essa imagem  -->
    <div  class="flex flex-row w-full justify-end pr-[20%] z-20">
      <div class="absolute bg-center -mt-6 w-44 p-3 h-14 bg-contain bg-no-repeat cursor-pointer containsTooltip" 
          style="background-image: url('{$imagesUrl}/enter_arrows.png')  "
          _originalTooltip_={$Terms.use_arrow_keys} > 
      </div>
    </div>

    <!-- container do form edicao 
    necessario ser 'relative' para que a slidingMessage ao final da div que é 'absolute' fique dentro da relative -->

    <div  class="flex flex-col w-full bg-white relative rounded-lg"  >

        <!-- titulo e botao fechar -->
        <div id='divWINDOW_TOP'>
          
          <div id='divWINDOW_TITLE'>
            {#if (formHttpMethodApply=='POST')  }
              {$Terms.new_booking_form_title}  
            {:else if (formHttpMethodApply=='PATCH')  }
              {$Terms.edit_booking_form_title}  
            {/if}
          </div>

          <div class='flex flex-row '>
              <div class='bg-icon-div-draggable w-7 bg-transparent bg-no-repeat bg-contain bg-center mr-6 containsTooltip' _originalTooltip_={$Terms.msg_click_to_drag}  >
                &nbsp;
              </div>


              <div class='divWINDOW_BUTTON mr-2'  aria-hidden="true" >
                &nbsp;&nbsp;[ ? ]&nbsp;&nbsp;
              </div>

              <div class='divWINDOW_BUTTON mr-6'  on:click={closeBookingModalWindow} aria-hidden="true" >
                &nbsp;&nbsp;[ X ]&nbsp;&nbsp;
              </div>
          </div>
        
        </div>

        <!-- campos  do formulario -->
        <div class="flex flex-col w-full h-auto  px-4 my-6 " >

          <!-- 1a linha, data/hora retirada e devolucao veiculo, nome do motorista -->
          <div class="flex flex-row w-full gap-[10px] border-b-2 pb-1">

             <!-- data/hora retirada do veiculo -->
            <div class="flex flex-col w-[calc(40%-10px)] ">
              <div class="flex flex-row w-full  ">  
                <div class='w-[69%]'>{$Terms.fieldname_pick_up_car}:</div>
                <div class='w-[30%]'>{$Terms.hour_format_to_use}</div>
              </div>

              <div class="flex flex-row w-full  ">  

                  <div class="flex w-[70%] pr-3">  
                    
                      <!-- data -->
                      <div class='flex flex-col w-full'  >
                          <div class='flex flex-row w-full' >
                            <input type="text" autocomplete="off" sequence="1"  id="txtPickUpDate" placeholder={$Terms.date_placeholder}  maxlength='10' class='text_formFieldValue w-full' on:focus={(e)=>{e.target.select()}} >
                            <!-- icone que aciona escolha de data pelo calendario -->  
                            <div  class='btnDATE_PICKER' on:click={ () => {whichInputTextCalled='txtPickUpDate';showDatePicker=true}} aria-hidden="true"></div>    
                          </div>
                          <!-- msg erro caso ocorra -->
                          <div class='noerrorTextbox' id="errorPickUpDate">{$Terms.errormsg_fill_invalid_date}</div>
                      </div>

                  </div>

                  <!-- hora -->
                  <div class="flex w-[30%]">  
                    <div class='flex flex-col w-full'  >
                      <input type="text" autocomplete="off" sequence="2"  id="txtPickUpHour" placeholder={$Terms.hour_placeholder} maxlength='8' class='text_formFieldValue w-full' on:focus={(e)=>{e.target.select()}} >
                      <!-- msg erro caso ocorra -->
                      <div class='noerrorTextbox' id="errorPickUpHour">{$Terms.errormsg_fill_invalid_time}</div>  
                    </div>
                  </div>

              </div>
            </div>

             <!-- data/hora devolucao do veiculo -->
            <div class="flex flex-col w-[calc(40%-10px)]">

              <div class="flex flex-row w-full  ">  
                <div class='w-[69%]'>{$Terms.fieldname_drop_off_car}:</div>
                <div class='w-[30%]'>{$Terms.hour_format_to_use}</div>
              </div>
              <div class="flex flex-row w-full ">  

                    <div class="flex w-[70%]  pr-3">  

                        <!-- data -->
                        <div class='flex flex-col w-full'  >
                            <div class='flex flex-row w-full' >
                                <input type="text" autocomplete="off" sequence="3"  id="txtDropOffDate" placeholder={$Terms.date_placeholder} maxlength='10' class='text_formFieldValue w-full' on:focus={(e)=>{e.target.select()}} >
                                <!-- icone que aciona escolha de data pelo calendario -->  
                                <div  class='btnDATE_PICKER' on:click={ () => {whichInputTextCalled='txtDropOffDate';showDatePicker=true}} aria-hidden="true"></div>    
                            </div>  
                            <!-- msg erro caso ocorra -->
                            <div class='noerrorTextbox' id="errorDropOffDate">{$Terms.errormsg_fill_invalid_date}</div>
                        </div>  

                    </div>  

                    <!-- hora -->
                    <div class="flex w-[30%]">  
                      <div class='flex flex-col w-full'  >
                          <input type="text" autocomplete="off" sequence="4"  id="txtDropOffHour" placeholder={$Terms.hour_placeholder} maxlength='5' class='text_formFieldValue w-full' on:focus={(e)=>{e.target.select()}} >
                          <!-- msg erro caso ocorra -->
                          <div class='noerrorTextbox' id="errorDropOffHour">{$Terms.errormsg_fill_invalid_time}</div>
                      </div>
                    </div>  
                    
              </div>
            </div>

            <!-- nome motorista -->
            <div class="flex flex-col w-[20%]">
              <div class='w-full'>{$Terms.fieldname_driver_name}:</div>
              <div class="flex flex-row w-full">  
                    <div class="flex flex-col w-full">  
                      <input type="text" autocomplete="off" sequence="5"  id="txtDriverName" maxlength='50' class='text_formFieldValue w-full' on:focus={(e)=>{e.target.select()}} >
                      <!-- msg erro caso ocorra -->
                      <div class='noerrorTextbox' id="errorDriverName">{$Terms.errormsg_fill_driver_name}</div>
                    </div>
              </div>
            </div>

          </div>

          <!-- 1a linha, data/hora retirada e devolucao veiculo, nome do motorista -->
          <div class="flex flex-row w-full gap-[10px] h-[350px] items-end">

                <!-- putWhiteTooltip  nao temm propriedades CSS, é usado somente para endereçar TOOLTIP no jscript!   -->

                <!-- logo fabricante, modelo carro, foto carro -->
                <div class="flex w-[60%] h-full ">

                      <!-- fabricante, foto veiculo, valor da diaria-->
                      <div class="w-full flex  items-start  flex-col pt-3 pl-3 h-full   ">

                          <!--  foto do carro --> 
                          <div class="flex flex-row w-full h-[80%] justify-center "> 

                              <div class="bg-no-repeat bg-center bg-contain w-full h-full items-center flex justify-center" id='carDetails_Picture'>
                              
                                  {#if $imagesStillLoading} 
                                      <img id='carDetails_Picture'  alt=''  class='gifImageLoading' src='loading.gif'>
                                  {/if}


                              </div>
                              
                          </div>

                          <!--  logotipo da fabricante, nome fabricante, nome modelo carro e ano fabricacao --> 
                          <div class="flex flex-row h-[20%] items-center justify-around  w-full"> 

                            <div class="flex flex-col  ">
                                <div class="font-bold text-lg" id='carDetails_manufacturer_name'>&nbsp;</div>    <!-- nome fabricante  -->
                                <!--  modelo / ano carro -->
                                <div class="text-gray-500 pt-2">
                                    <span id='carDetails_name'>&nbsp;</span>
                                </div>
                            </div>

                            <!-- logotipo fabricante  -->
                            <div class="w-24 h-16 max-w-24 max-h-16 min-w-24 min-h-16 bg-no-repeat bg-center bg-contain" id='carDetails_ManufacturerLogo'></div>


                            <!--  diária --> 
                            <div class="flex pt-4 pl-2  justify-end "> 
                              <div class="text-gray-500 pb-4 ">{$Terms.fieldname_rental_price}: <span id='carDetails_rental_price' style='padding-left:10px'>&nbsp;</span></div> 
                            </div>
                          </div>

                      </div>

                </div>

                <!-- detalhes do carro, tipo cambio, qtde portas, etc -->                
                <div class="flex w-[40%] h-full flex-col items-start">

                    <div class="flex w-full flex-row justify-around h-[25%]">
                        <!-- cilindrada-->
                        <div class="flex flex-col w-[50%] justify-center cursor-pointer  items-center putWhiteTooltip"  >
                          <div class="w-20 h-11 rounded-full p-2 bg-icon-cylinder-capacity bg-center bg-no-repeat bg-[length:24px_24px]" > 
                          </div>

                          <div class="bookingCarDetail invisible"><span id='carDetails_cc'>&nbsp;</span>&nbsp;{$Terms.cylinder_capacity}</div>
                        </div>              

                        <!-- qtde portas -->
                        <div class="flex flex-col w-[50%] justify-center cursor-pointer  items-center putWhiteTooltip"  >
                          <div class="w-20 h-11 rounded-full p-2 bg-icon-doors bg-center bg-no-repeat bg-[length:24px_24px]" > 
                          </div>

                          <div class="bookingCarDetail invisible"><span id='carDetails_doors'>&nbsp;</span>&nbsp;{$Terms.doors}</div>
                        </div>              
                    </div>

                    <div class="flex w-full flex-row justify-around  h-[25%]">
                        <!-- consumo de combustivel-->
                        <div class="flex flex-col w-[50%] justify-center cursor-pointer  items-center putWhiteTooltip"  >
                          <div class="w-20 h-11 rounded-full p-2 bg-icon-consume bg-center bg-no-repeat bg-[length:24px_24px]" > 
                          </div>

                          <div class="bookingCarDetail invisible"><span id='carDetails_mpg'>&nbsp;</span>&nbsp;{$Terms.car_consume}</div>
                        </div>              

                        <!-- tipo da transmissao -->
                        <div class="flex flex-col w-[50%] justify-center cursor-pointer  items-center putWhiteTooltip"   >
                          <div class="w-20 h-11 rounded-full p-2 bg-icon-transmission bg-center bg-no-repeat bg-[length:24px_24px]" > 
                          </div>

                          <div class="bookingCarDetail invisible"><span id='carDetails_transmission_short'>&nbsp;</span></div>
                        </div>              
                    </div>

                    <div class="flex w-full flex-row justify-around  h-[25%]">
                        <!-- força motor (hps) -->
                        <div class="flex flex-col w-[50%] justify-center  cursor-pointer  items-center putWhiteTooltip"   >
                          <div class="w-20 h-11 rounded-full p-2 bg-icon-power bg-center bg-no-repeat bg-[length:24px_24px]" > 
                          </div>

                          <div class="bookingCarDetail invisible"><span id='carDetails_hp'>&nbsp;</span>&nbsp;{$Terms.car_power}</div>
                        </div>              

                        <!-- odometro -->
                        <div class="flex flex-col w-[50%] justify-center cursor-pointer   items-center putWhiteTooltip"   >
                          <div class="w-20 h-11 rounded-full p-2 bg-icon-odometer bg-center bg-no-repeat bg-[length:24px_24px]" > 
                          </div>

                          <div class="bookingCarDetail invisible"><span id='carDetails_odometer'>&nbsp;</span>&nbsp;{$Terms.odometer_unit}</div>
                        </div>              
                    </div>

                    <div class="flex w-full flex-row justify-around  h-[25%]">
                        <!-- ano de fabricacao -->
                        <div class="flex flex-col w-[50%] justify-center cursor-pointer  items-center putWhiteTooltip"   >
                          <div class="w-20 h-11 rounded-full p-2 bg-icon-year bg-center bg-no-repeat bg-[length:24px_24px]" > 
                          </div>

                          <div class="bookingCarDetail invisible"><span id='carDetails_year'>a&nbsp;</span></div>
                        </div>              

                        <!-- qtde cilindros -->
                        <div class="flex flex-col w-[50%] justify-center cursor-pointer  items-center putWhiteTooltip"   >
                          <div class="w-20 h-11 rounded-full p-2 bg-icon-cylinders bg-center bg-no-repeat bg-[length:24px_24px]" > 
                          </div>

                          <div class="bookingCarDetail invisible"><span id='carDetails_cylinders'>&nbsp;</span>&nbsp;{$Terms.cylinders}</div>
                        </div>              

                    </div>
                </div>

            </div>  

        </div>

        <!-- botoes salvar/sair -->
        <div class="flex flex-row w-full justify-between px-6 border-t-[1px] border-t-gray-300 py-2">
          <button  id="btnCLOSE" class="btnCANCEL" on:click={closeBookingModalWindow} >{$Terms.button_cancel}</button>

          {#if formHttpMethodApply!='POST'}
            <button  id="btnMETADATA" class="btnMETADATA" on:click={ () => {invokeBookingMetadataWindow()}} aria-hidden="true">{$Terms.button_metadata}</button>
            <button  id="btnDELETE" class="btnDELETE" on:click={()=>{performDeleteBookingRecord()}} aria-hidden="true">{$Terms.button_booking_delete}</button>
          {/if}

          <button  id="btnSAVE" class="btnSAVE" on:click={() => {performSaveBookingRecord()}} aria-hidden="true">{$Terms.button_book_car}</button>
        </div>

        <!-- mensagem rolante usada para avisos -->
        <div  id="slidingFormMessage" >
          &nbsp;
        </div>


    </div> 

</div> 


