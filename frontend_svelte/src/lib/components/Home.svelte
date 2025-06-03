<style>
</style>

<script>

// carrega 'stores' (variaveis writable) que serao usadas aqui
import { Terms, backendUrl, imagesUrl, currentHomeSelectedCarId } from '$lib/stores/stores.js'
import { selectedCountry, getAndShowCarDetails, logAPI } from '$js/utils.js'

// Card dos carros que sera exibido no navegador de carros (homeCarsBrowser)
import CarCard from '$lib/components/CarCard.svelte'

import Calendar from '$lib/components/Calendar.svelte'


// recordset com os dados de cards de carros
let carCards



/************************************************************************************************************************************************************
funcao que sera carregada ao renderizar o componente Home.svelte
************************************************************************************************************************************************************/
async function mountCarBrowser () {

  // necessario abrir evento assincrono para exibir div ajax loading, caso contrario navegador nao atualiza a tela
  setTimeout(() => {
    showLoadingGif();
  }, 1);

  let country = selectedCountry()

  $currentHomeSelectedCarId = 0

  let _currentWorkgroupName = localStorage.getItem("rentacar_workgroup_name");

  if (_currentWorkgroupName == null )   return;
  
  let _route_ =  `${$backendUrl}/${_currentWorkgroupName}/car_cards/${country}`
  logAPI('GET', _route_)
  await fetch(_route_, {method: "GET"} )

  .then((response) => {
  
    if (!response.ok) {
      throw new Error(`Car Browser Mounting Err Fatal= ${response.status}`);
    }
    return response.json();
  }) 

  .then( (carCardsData) => {
    // obtem o ID do primeiro carro obtido, para colocar destaque no mesmo e ja começar mostrando suas propriedades
    // o usuario pode mudar o carro atualmente selecioado (currentHomeSelectedCarId) clicando sobre seu respectivo card
    if (carCardsData.length>0)
      $currentHomeSelectedCarId = carCardsData[0].id   

    setTimeout(() => {
      hideLoadingGif();
    }, 1);

    carCards = carCardsData    
  })
}

// se mudar o idioma/país, recarrega browser de carros
// o svelte dispara qq funcao vinculada à alteracao do valor de $Terms (por ser uma declaracao reativa), 
// e qdo Home.svelte é carregado, svelte considera mudança de valor da variavel, e dispara automatico, sem necessidade de criar onMount()
$: $Terms, mountCarBrowser()

// ligam/desligam o slide automatico do navegador de carros
let start_slide_car_to_right
let start_slide_car_to_left

/************************************************************************************************************************************************************
faz slide automatico do navegador de carros qdo usuario passa mouse sobre setas direita/esquerda
************************************************************************************************************************************************************/
const start_slidingCarCards = (direction) => {
  if (direction==1) 
    start_slide_car_to_right = setInterval(function () {jq('#homeCarsBrowser').scrollLeft( jq('#homeCarsBrowser').scrollLeft() + 40 )}, 100);
    
  else 
    start_slide_car_to_left = setInterval(function () {jq('#homeCarsBrowser').scrollLeft( jq('#homeCarsBrowser').scrollLeft() - 40 )}, 100);

}

/************************************************************************************************************************************************************
desliga slide automatico do navegador de carros qdo usuario passa retira o mouse das setas direita/esquerda
************************************************************************************************************************************************************/
const stop_slidingCarCards = (e) => {
  clearInterval(start_slide_car_to_left)
  clearInterval(start_slide_car_to_right)
}



/************************************************************************************************************************************************************
avança/retrocede scroll horizontal do navegador de carros se o usuario rolar o mouse wheel (vertical)
************************************************************************************************************************************************************/
const handleWheel = e => {
  if (e.type == "wheel") {
    var getDelta = e.deltaY;

    if (getDelta>0) 
      jq('#homeCarsBrowser').scrollLeft( jq('#homeCarsBrowser').scrollLeft() + 100 )
    else 
      jq('#homeCarsBrowser').scrollLeft( jq('#homeCarsBrowser').scrollLeft() - 100 )
  };
}




// cores de fundo usadas em PropertyCards
let carPropertyColor = [ 
  { name: 'brown', normal: '#FBE0CB', stronger: '#F2D1BA' }, 
  { name: 'blue', normal: '#D6E5FC', stronger: '#B4D0F9' },  
  { name: 'orange', normal: '#FBE1CC', stronger: '#F2D1BA' },
  { name: 'purple', normal: '#F7EDFE', stronger: '#EBDDF7' },
  { name: 'green', normal: '#EEEFB3', stronger: '#E2E3AA' },
  { name: 'navy', normal: '#EDECF9', stronger: '#D0DAF7' },
  { name: 'red', normal: '#FFE9F2', stronger: '#ffcce1' },
  { name: 'blue', normal: '#EEF1FA', stronger: '#E6E9F3' },
  { name: 'lighterblue', normal: '#E9F2FF', stronger: '#E0EAF9' }
]


/************************************************************************************************************************************************************
 obtem cor de fundo para o PropertyCard
************************************************************************************************************************************************************/
const getPropertyCarColor = (name, which) => {
  const result = carPropertyColor.filter((color) => color.name === name);

  if (which=="normal")   return result[0].normal
  else return result[0].stronger
}

let whichComponentSvelteCalledCalendar = 'Home'

</script>


<div class="w-full flex  items-center flex-col font-Roboto  "> 

    <!-- 1a linha: titulo "bem vindo, visitante" -->
    <div class="w-full mb-2 h-1 text-lg " > 
        {$Terms.available_cars}
    </div>


    <!-- 2a linha: navegacao pelos veiculos disponiveis -->
    <div class="w-full flex  items-center h-44  flex-row relative "> 

        <!-- seta para esquerda-->
        <div class="  absolute z-20 bg-slate-100 rounded-full p-2 cursor-pointer hover:bg-slate-300 -ml-4"  
            on:mouseenter={() => {start_slidingCarCards(-1)}}  on:mouseleave={stop_slidingCarCards}  aria-hidden="true"  >
          <svg width="48px" height="48px" viewBox="0 0 1024 1024" class="icon" version="1.1" xmlns="http://www.w3.org/2000/svg" fill="#000000"><g id="SVGRepo_bgCarrier" stroke-width="0"></g><g id="SVGRepo_tracerCarrier" stroke-linecap="round" stroke-linejoin="round"></g><g id="SVGRepo_iconCarrier"><path d="M768 903.232l-50.432 56.768L256 512l461.568-448 50.432 56.768L364.928 512z" fill="#000000"></path></g></svg>
        </div>

        <!-- navegador de car cards -->
        <div class="w-full grow flex flex-row gap-8 overflow-x-auto z-10 overflow-y-visible " id="homeCarsBrowser" on:wheel={(e) => {handleWheel(e)}} >  
                {#if carCards}

                    {#each carCards as {name, id, car_image}} 
                      <CarCard 
                            {id} 
                            car_image_url = {$imagesUrl + car_image}
                            {name}  
                      />
                    {/each}

                    <!-- terminou listar os carros, mostra detalhes do 1o carro  -->
                    <!-- se nao colocar como hidden, svelte tenta mostrar conteudo da funcao abaixo -->

                    <div class="hidden">{  getAndShowCarDetails($currentHomeSelectedCarId) }</div> 
                    
                {:else}
                  <div class="divLOADING_MSG">{$Terms.loading}</div>
                {/if}

        </div>


        <!-- seta para direita -->
        <div class="absolute z-20 bg-slate-100 rounded-full p-2 cursor-pointer hover:bg-slate-300 -mr-4 right-0" 
                on:mouseenter={() => {start_slidingCarCards(1)}} on:mouseleave={stop_slidingCarCards}  aria-hidden="true"  >
          <svg width="48px" height="48px" viewBox="0 0 1024 1024" class="icon" version="1.1" xmlns="http://www.w3.org/2000/svg" fill="#000000"><g id="SVGRepo_bgCarrier" stroke-width="0"></g><g id="SVGRepo_tracerCarrier" stroke-linecap="round" stroke-linejoin="round"></g><g id="SVGRepo_iconCarrier"><path d="M256 120.768L306.432 64 768 512l-461.568 448L256 903.232 659.072 512z" fill="#000000"></path></g></svg>
        </div>

    </div>

    <!-- 3a linha: detalhes do veiculo clicado (lado esquerdo) e calendario agendamento (lado direito) -->
    <div class="w-full flex  flex-row  gap-5  "> 

        <!-- lado esquerdo, propriedades do carro -->
        <div class="w-[75%] flex  items-center   flex-row  bg-white rounded-3xl border-2 border-gray-300 h-full gap-4 "> 

            <!-- fabricante, foto veiculo, valor da diaria-->
            <div class="w-2/5 flex  items-start  flex-col pt-3 pl-3 h-full   "> 

                <!--  foto do carro --> 
                <div class="flex flex-row w-full h-[40%] justify-center "> 
                    <!--  foto do veiculo  --> 
                    <div class="bg-no-repeat bg-center bg-contain w-full h-full" id='carDetails_Picture'>&nbsp;</div> 
                </div>

                <!--  logotipo da fabricante, nome fabricante, nome modelo carro e ano fabricacao --> 
                <div class="flex flex-row h-[30%] items-center"> 

                  <!-- logotipo fabricante  -->
                  <div class="w-24 h-20 max-w-24 max-h-20 min-w-24 min-h-20 bg-no-repeat bg-center bg-contain " id='carDetails_ManufacturerLogo'></div>

                  <div class="pl-6 flex flex-col w-full pt-3 ">
                      <div class="font-bold text-lg" id='carDetails_manufacturer_name'>&nbsp;</div>    <!-- nome fabricante  -->
                      <!--  modelo / ano carro -->
                      <div class="text-gray-500 pt-2">
                          <span id='carDetails_name'>&nbsp;</span> - 
                          <span id='carDetails_year'>&nbsp;</span>
                      </div>
                  </div>
                </div>


                <!--  diária --> 
                <div class="flex pt-4 pl-2 flex-col h-[25%]  justify-end "> 
                  <div class="text-gray-500 pb-4 ">{$Terms.fieldname_rental_price}: <span id='carDetails_rental_price' style='padding-left:10px'>&nbsp;</span></div> 
                </div>

            </div>


            <!-- lado direito, propriedades  -->
            <!-- propriedades do veiculo => cilindrada, potencia, etc -->
            <div class="w-3/5 flex items-start flex-col pt-3 pr-3 h-full"> 

              <!-- 1a linha com 6 propriedades genéricas -->
              <div class="w-full flex gap-4 "> 

                <!-- 1a coluna, informacoes: cilindrada e km/litro -->
                <div class="flex flex-col gap-3 w-1/3"  >

                    <div class="h-CarPropertyCardType1 w-full rounded-2xl flex flex-col items-center justify-center border-transparent border-2 cursor-pointer hover:border-2 hover:border-gray-500"  
                              style="background-color: {getPropertyCarColor('brown', 'normal')}" >

                      <!-- icone cilindradas  -->
                      <div class="w-12 h-12 rounded-full p-2 bg-icon-cylinder-capacity bg-center bg-no-repeat bg-[length:24px_24px]" style="background-color: {getPropertyCarColor('brown', 'stronger')}"> 
                      </div>
                      
                      <div class="flex flex-col pt-3 items-center text-center">
                        <div class="homeCarDetail"><span id='carDetails_cc'>&nbsp;</span>&nbsp;{$Terms.cylinder_capacity}</div>
                      </div>

                    </div>

                    <div class="h-CarPropertyCardType2 w-full rounded-2xl flex flex-col items-center justify-center border-transparent border-2 cursor-pointer hover:border-2 hover:border-gray-500"
                            style="background-color: {getPropertyCarColor('green', 'normal')}" >

                      <!-- icone qtde cilindros -->
                      <div class="w-12 h-12 rounded-full p-2 bg-icon-cylinders bg-center bg-no-repeat bg-[length:24px_24px]" style="background-color: {getPropertyCarColor('green', 'stronger')}">  
                      </div>
                      
                      <div class="flex flex-col pt-3 items-center">
                        <div class="homeCarDetail"><span id='carDetails_cylinders'>&nbsp;</span>&nbsp;{$Terms.cylinders}</div>
                      </div>

                    </div>



                </div>


                <!-- 2a coluna, qtde portas e tipo cambio -->
                <div class="flex flex-col gap-3 w-1/3"  >

                  <div class="h-CarPropertyCardType2 w-full rounded-2xl flex flex-col items-center justify-center border-transparent border-2 cursor-pointer hover:border-2 hover:border-gray-500"
                          style="background-color: {getPropertyCarColor('blue', 'normal')}"> 


                    <!-- icone qtde portas -->
                    <div class="w-12 h-12 rounded-full p-2 bg-icon-doors bg-center bg-no-repeat bg-[length:24px_24px]" style="background-color: {getPropertyCarColor('blue', 'stronger')}">
                    </div>
                    
                    <div class="flex flex-col pt-3 items-center">
                      <div class="homeCarDetail"><span id='carDetails_doors'>&nbsp;</span>&nbsp;{$Terms.doors}</div>
                    </div>

                  </div>

                  <div class="h-CarPropertyCardType1 w-full rounded-2xl flex flex-col items-center justify-center border-transparent border-2 cursor-pointer hover:border-2 hover:border-gray-500"  
                            style="background-color: {getPropertyCarColor('purple', 'normal')}">

                    <!-- icone tipo de transmissao -->
                    <div class="w-12 h-12 rounded-full p-2 bg-icon-transmission bg-center bg-no-repeat bg-[length:24px_24px]" style="background-color: {getPropertyCarColor('purple', 'stronger')}">
                    </div>
                    
                    <div class="flex flex-row pt-3 items-center text-center">
                      <div class="homeCarDetail"><span id='carDetails_transmission'>&nbsp;</span></div>
                    </div>

                  </div>

                </div>

                <!-- 3a coluna: km/litros potencia motor  -->
                <div class="flex flex-col gap-3 w-1/3"  >

                    <div class="h-CarPropertyCardType1 w-full  rounded-2xl flex flex-col items-center justify-center border-transparent border-2 cursor-pointer hover:border-2 hover:border-gray-500"  
                              style="background-color: {getPropertyCarColor('red', 'normal')}" > 

                      <!-- icone consumo km/litro ou miles per galon -->
                      <div class="w-12 h-12 rounded-full p-2 bg-icon-consume bg-center bg-no-repeat bg-[length:24px_24px]" style="background-color: {getPropertyCarColor('red', 'stronger')}"> 
                      </div>
                      
                      <div class="flex flex-col pt-3 items-center">
                        <div class="homeCarDetail"><span id='carDetails_mpg'>&nbsp;</span>&nbsp;{$Terms.car_consume}</div>
                      </div>

                    </div>

                    <div class="h-CarPropertyCardType2 w-full rounded-2xl flex flex-col items-center justify-center border-transparent border-2 cursor-pointer hover:border-2 hover:border-gray-500"
                            style="background-color: {getPropertyCarColor('navy', 'normal')}" >

                      <!-- icone força (hps) -->
                      <div class="w-12 h-12 rounded-full p-2 bg-icon-power bg-center bg-no-repeat bg-[length:24px_24px]" style="background-color: {getPropertyCarColor('navy', 'stronger')}">  
                      </div>
                      
                      <div class="flex flex-col pt-3 items-center">
                        <div class="homeCarDetail"><span id='carDetails_hp'>&nbsp;</span>&nbsp;{$Terms.car_power}</div>
                      </div>
                    </div>

                </div>

              </div>

              <!-- 2a linha, odometro do carro  -->
              <div class="w-full flex mt-3 "> 

                  <div class="w-full p-4 rounded-2xl flex flex-row items-center justify-start border-transparent border-2 cursor-pointer hover:border-2 hover:border-gray-500"
                          style="background-color: {getPropertyCarColor('navy', 'normal')}" >

                    <!-- icone odometro -->
                    <div class="w-12 h-12 rounded-full p-2 bg-icon-odometer bg-center bg-no-repeat bg-[length:24px_24px]" style="background-color: {getPropertyCarColor('navy', 'stronger')}">  
                    </div>
                    
                    <div class="flex  pl-5  items-center flex-row"> 
                      <div class="homeCarDetail">{ $Terms.odometer }</div>
                      <div class="pl-4"><span id='carDetails_odometer'>&nbsp;</span>&nbsp;{$Terms.odometer_unit}</div>
                    </div>

                  </div>
              </div>

              <!-- espacamento ao final do card 'odometro' -->
              <div class="w-full flex h-3 "></div> 

            </div>


        </div>

        <!-- lado direito, calendario com agenda de uso do veiculo -->
        <div class="w-[25%] flex align-top text-lg items-start justify-start  flex-col  bg-white rounded-3xl border-2 border-gray-300 h-full"> 
          <Calendar {whichComponentSvelteCalledCalendar}  />
        </div>

    </div>

    <div class="flex h-50  flex-col">
        <br><br><br><br>
    </div>





</div>

