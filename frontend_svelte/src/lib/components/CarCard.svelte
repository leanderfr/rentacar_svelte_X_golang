

<script>

// variaveis genericas, writeables

export let id, name, car_image_url

import {  getAndShowCarDetails } from '$js/utils.js'
import { currentHomeSelectedCarId } from '$lib/stores/stores.js'

let carCardHome = 'carCardHome'

/************************************************************************************************************************************************************
usuario escolheu um carro para ver detalhes (Home.svelte)
************************************************************************************************************************************************************/

const showCarDetailsHome = (carId) => {

jq('.carCardHomeClicked').removeClass('carCardHomeClicked')  
jq(`#carCard${carId}`).addClass('carCardHomeClicked')

// memoriza na variavel writable qual id do carro selecionada em Home.svelte
// faz isso para que o calendario saiba qual carro verificar se possui dias reservados
$currentHomeSelectedCarId = carId

getAndShowCarDetails(carId)
}


</script>

<!--
*************************************************************************************************************************************
*************************************************************************************************************************************
*************************************************************************************************************************************

HTML 

*************************************************************************************************************************************
*************************************************************************************************************************************
*************************************************************************************************************************************
-->

<div class="w-44  h-full flex flex-col shrink-0 "  id="car_card_{id}"> 

    <!-- necessario colocar uma div em branco como um espaço, pois a imagem do carro subirá um pouco (-top-2) -->
    <div class="h-6 w-full ">
        &nbsp;
    </div>

    <div class="car_card w-full rounded-2xl flex flex-col items-center cursor-pointer  bg-gray-200 h-[100px] "  aria-hidden="true" 
       class:carCardHomeClicked={ $currentHomeSelectedCarId == id } class:carCardHome on:click={ () => showCarDetailsHome(id)}   id='carCard{id}'  >

      <div class="w-full h-[80px]  justify-center flex  ">
          <div class="w-[90%] bg-no-repeat bg-center bg-contain h-full  relative -top-8"    style="background-image: url('{car_image_url}')" ></div>    
      </div>

      <div class="h-[20px] relative -top-3">      
        <div class="h-3/5 text-base font-bold">{name}</div>
      </div>

    </div>

</div>


