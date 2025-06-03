
<script>

// carrega 'stores' (variaveis writable) que serao usadas aqui
import {  currentHomeSelectedCarId, Terms, elementIdCSSChangedTemporarily, elementCSSChangedTemporarily, HomeCalendar_CurrentFirstDate, BookingCalendar_CurrentDate } from '$lib/stores/stores.js'

import { getAndShowCarDetails } from '$js/utils.js'

$: month_names = [$Terms.january, $Terms.february, $Terms.march, $Terms.april, $Terms.may, $Terms.june,
                  $Terms.july, $Terms.august, $Terms.september, $Terms.october, $Terms.november, $Terms.december]

// prepara disparo de eventos no componente pai
import { createEventDispatcher } from 'svelte'
const dispatch = createEventDispatcher()



// recebe a informacao de qual componente (.svelte) chamou <Calendar />
export let whichComponentSvelteCalledCalendar 




// recebe ID do input type=text que invocou o calendario para escolha da data
// pode vir em branco, qdo <Calendar /> foi invocado por um componente sem input type= text
export let whichInputTextCalled 


// altera algumas caracteristicas visuais dependendo do componente que vai usar <Calendar />
let svgArrowHeight, svgArrowWidth, dayLineHeight, monthFontSize,  showCloseButton, dayFontSize
let buttonsContainerHeight, buttonsDivHeight, dayWeekFontSize, titleLineHeight, ptWeekDays, currentMonthDaysColor

// calendario usado no canto direito da pagina Home
if (whichComponentSvelteCalledCalendar =='Home')     {             //  Home.svelte invocou Calendar.svelte
  svgArrowHeight = '40px'
  svgArrowWidth = '70px'
  dayLineHeight = 'h-11'   // classe tailwind 
  monthFontSize = '14px'
  showCloseButton = false
  dayFontSize = ''
  buttonsDivHeight = 'h-16'
  buttonsContainerHeight = 'h-10'
  titleLineHeight = "h-12"
  ptWeekDays = "pt-5"
  currentMonthDaysColor = "text-gray-500"
}

// calendario invocado para escolha de data, a partir da agenda de veiculos (Bookings.svelte)
else if (whichComponentSvelteCalledCalendar =='Booking') {           //  Bookings.svelte invocou Calendar.svelte
  svgArrowHeight = '44px'
  svgArrowWidth = '120px'
  dayLineHeight = 'h-14'         // classe tailwind 
  monthFontSize = '24px'
  showCloseButton = true
  dayFontSize = ''
  buttonsDivHeight = 'h-16'
  buttonsContainerHeight = 'h-10'
  titleLineHeight = "h-12"
  ptWeekDays = "pt-7"
  currentMonthDaysColor = "text-gray-500"

}

// calendario invocado quando digitando alguma data (input type=text) no formulario de reserva de veiculo
else if (whichComponentSvelteCalledCalendar =='BookingForm')   {   // qdpo chamado por BookingForm.svelte, servira como date picker de campos input type=text que recebem datas
  svgArrowHeight = '24px'
  svgArrowWidth = '50px'
  dayLineHeight = 'h-7'         // classe tailwind 
  monthFontSize = '14px'
  showCloseButton = true
  dayFontSize = 'text-[13px]'
  buttonsDivHeight = 'h-11'
  buttonsContainerHeight = 'h-6'
  dayWeekFontSize = 'text-sm'
  titleLineHeight = "h-9"
  ptWeekDays = "pt-3"
  currentMonthDaysColor = "text-black"  // date picker tem fundo cinza, necessario colocar cor mais forte nos numeros
}




/************************************************************************************************************************************************************
 monta o calendario em formato de array
 $CalendarCurrentDateHome ==>  1o dia do mês vigente (inicia em Stores.js) 
 $CalendarCurrentDateHome sera atualizada a medida que usuario avançar/retroceder meses
************************************************************************************************************************************************************/
let calendar 

function mountCalendar() {
  calendar = [] 

  let firstDayCurrentMonth = new Date($HomeCalendar_CurrentFirstDate.getFullYear(), $HomeCalendar_CurrentFirstDate.getMonth(), $HomeCalendar_CurrentFirstDate.getDate());
  let currentDate = new Date($HomeCalendar_CurrentFirstDate.getFullYear(), $HomeCalendar_CurrentFirstDate.getMonth(), $HomeCalendar_CurrentFirstDate.getDate());

  // retrocede ate achar o ultimo domingo antes do 1o dia do mes sendo visualizado
  while (currentDate.getDay()!=0)  {
    currentDate.setDate(currentDate.getDate() - 1);
  }

  let weekday, week, weekdays

  let dayColor    // dias do mês atual terão cor normal, dias de mês anterior ou posterior, cor mais fraca
  let today = new Date();  // hoje
  let _today_ = new Date(today.getFullYear(), today.getMonth(), today.getDate());

  for (week=0; week<6; week++ )  {
    weekdays=[]
    for (weekday=0; weekday<7; weekday++) {

      // dias em meses anteriores, mas exibidos no calendario ficarão em cor mais fraca
      if (currentDate.getMonth() != firstDayCurrentMonth.getMonth())    dayColor = 'weaker'
      else    dayColor = 'normal'

      if (currentDate.getTime() == _today_.getTime()) dayColor = 'today_color'

      weekdays.push( [currentDate.getDate(), dayColor, currentDate.getTime()] )   // memoriza o dia, respectiva cor, e a data em formato jscript caso o dia seja clicado, escolhido
      currentDate.setDate(currentDate.getDate() + 1);    // avanca 1 dia
    }
    calendar.push(weekdays)  
  }

  let year = $HomeCalendar_CurrentFirstDate.getFullYear()

  // atualiza o texto do botoes para avançar/retroceder ano, dependendo da data atual  ($HomeCalendar_CurrentFirstDate)
  setTimeout(() => {
    jq('#btnPreviousYear').html( year-1 )
    jq('#btnNextYear').html( year+1 )  

    // whichComponentSvelteCalledCalendar  == BookingForm, usuario esta editando uma reserva de veiculo e clicou no icone para escolha de data

    // posiciona o calendario logo abaixo do input type=text que o invocou
    if (whichComponentSvelteCalledCalendar =='BookingForm')  {
      
      // vai alterar bg-color do input type=text para ficar igual ao do date picker, necessario memorizar todo CSS para devolver ao fechar o date picker (<Calendar.svelte />)
      $elementCSSChangedTemporarily = jq( `#${whichInputTextCalled}` ).prop('style') 
      $elementIdCSSChangedTemporarily = `#${whichInputTextCalled}` 

      //jq( `#${whichInputTextCalled}` ).css('border-color','black') 
      jq( `#${whichInputTextCalled}` ).css('background-color','#F3F4F6')   // mesma cor do date picker (Calendar.svelte)

      // BookingFormCalendarContainer =>  div container do calendario (em BookingForm.svelte)
      // posiciona o calendario logo abaixo do input type=text que o invocou
      jq( "#BookingFormCalendarContainer" ).position({
          of: jq( `#${whichInputTextCalled}` ),
          my: 'left top',
          at: 'left bottom',
      });

      // torna a div container visivel novamente, ela tinha sido escondida pq nao estava posicionada embaixo do input type=text
      jq( "#BookingFormCalendarContainer" ).css('visibility','visible')

    }

  }, 100);

  return calendar
}


/************************************************************************************************************************************************************
avanca/retorna mes no calendario
************************************************************************************************************************************************************/
const browseCalendar = (unit, number) => {

  // obtem o mês/ano da data exibida atualmente para avançar/retroceder
  let tmpFirstDate = new Date($HomeCalendar_CurrentFirstDate.getFullYear(), $HomeCalendar_CurrentFirstDate.getMonth(), $HomeCalendar_CurrentFirstDate.getDate());

  if (unit  == 'month')  tmpFirstDate.setMonth(tmpFirstDate.getMonth() + number);
  if (unit  == 'year')  tmpFirstDate.setFullYear(tmpFirstDate.getFullYear() + number);

  // dia de hoje, na verdade, para estar de acordo com a logica do calendario, calcula o 1o dia do mes vigente
  if (unit  == 'today')  {

    let today = new Date();  // hoje 


    // se esta escolhendo data para o formulario de reserva, e clicou em HOJE, joga valor da data para o textbox destino (atraves de dispatch)
    if (whichComponentSvelteCalledCalendar =='BookingForm')        {
      calendarDateClicked(today)
      return  
    }
    else {
      tmpFirstDate = new Date(today.getFullYear(), today.getMonth(), 1);
    }
  }

  $HomeCalendar_CurrentFirstDate = tmpFirstDate

  mountCalendar()

  // se a tela atualmente visivel é Home.svelte, apos montar o calendario, verifica no mes/ano visualizado se o veiculo possui reservas
  // faz isso para colocar um X vermelho sobre os dias reservados
  let isHomeVisible = typeof jq('#homeCarsBrowser').attr('id')!='undefined'
  if (isHomeVisible)   getAndShowCarDetails($currentHomeSelectedCarId)  
}


/************************************************************************************************************************************************************
um dia foi clicado, escolhido no calendario
atualiza um elemento DOM (tag html) do componente pai baseado em sua ID 
clickedDate = milisegundos (getTime())
************************************************************************************************************************************************************/
const calendarDateClicked = (clickedDate) => {

if (whichComponentSvelteCalledCalendar =='Home') {}   // componente Home.svelte exibe calendario apenas para informar, nao para escolher data

if (whichComponentSvelteCalledCalendar =='Booking') {
  $BookingCalendar_CurrentDate = new Date(clickedDate);   // converte milisegundos para Date
  dispatch('dispatchCloseCalendar')    // fecha calendario
  dispatch('dispatchrefreshBookingDatesAndContent')    // dispara atualizacao do calendario de reservas de carros (Bookings.svelte)
}

// form de reserva de veiculo invocou Calendar.svelte, joga a data escolhida para o input type=text que invocou Calendar
if (whichComponentSvelteCalledCalendar =='BookingForm') {       
  dispatch('dispatchPutChosenDateInInputText', {clickedDate: clickedDate})    // dispara funcao que joga a data para o input type=text destino
  dispatch('dispatchCloseCalendar')    // fecha calendario
}

}


/********************************************************************************************************************************************************
 se usuario clicar no botao 'fechar' , dispara evento (no componete pai) para fechar calendario
********************************************************************************************************************************************************/
const btnCloseCalendarClicked = () => {
    dispatch('dispatchCloseCalendar')
}


mountCalendar()

</script>

<!-- se o calendario foi aberto em forma de janela e 'showCloseButton' esta exibido, 
cria div 'calendarModalWindow' para informar ao jscript que a tecla Esc pode fechar o formulario   -->
{#if showCloseButton}
  <div id='calendarModalWindow'></div>
{/if}


<!-- mes ano (texto) e setas avança / retrocede mes/ano --> 
<div class="pt-2  w-full   font-bold flex flex-row columns-3 justify-around pb-2 {titleLineHeight} " >  
    <div class="ml-1  cursor-pointer hover:bg-[#B1DCFB] w-[{svgArrowWidth}] items-center flex justify-center {dayLineHeight}" on:click={() => {browseCalendar('month',-1)}} aria-hidden="true" >
      <svg width={svgArrowWidth} height={svgArrowHeight} viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg"><g id="SVGRepo_bgCarrier" stroke-width="0"></g><g id="SVGRepo_tracerCarrier" stroke-linecap="round" stroke-linejoin="round"></g><g id="SVGRepo_iconCarrier"> <path fill-rule="evenodd" clip-rule="evenodd" d="M11.7071 4.29289C12.0976 4.68342 12.0976 5.31658 11.7071 5.70711L6.41421 11H20C20.5523 11 21 11.4477 21 12C21 12.5523 20.5523 13 20 13H6.41421L11.7071 18.2929C12.0976 18.6834 12.0976 19.3166 11.7071 19.7071C11.3166 20.0976 10.6834 20.0976 10.2929 19.7071L3.29289 12.7071C3.10536 12.5196 3 12.2652 3 12C3 11.7348 3.10536 11.4804 3.29289 11.2929L10.2929 4.29289C10.6834 3.90237 11.3166 3.90237 11.7071 4.29289Z" fill="#000000"></path> </g></svg>
    </div> 
    <div class="  grow text-center flex {dayLineHeight} justify-center items-center" style='font-size:{monthFontSize}'>
      { month_names[$HomeCalendar_CurrentFirstDate.getMonth()] }, {$HomeCalendar_CurrentFirstDate.getFullYear()}
    </div>
    <div class="mr-1  cursor-pointer hover:bg-[#B1DCFB] w-[{svgArrowWidth}]  items-center flex justify-center {dayLineHeight}"  on:click={() => {browseCalendar('month',1)}} aria-hidden="true" >
      <svg width={svgArrowWidth} height={svgArrowHeight} viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg"><g id="SVGRepo_bgCarrier" stroke-width="0"></g><g id="SVGRepo_tracerCarrier" stroke-linecap="round" stroke-linejoin="round"></g><g id="SVGRepo_iconCarrier"> <path fill-rule="evenodd" clip-rule="evenodd" d="M12.2929 4.29289C12.6834 3.90237 13.3166 3.90237 13.7071 4.29289L20.7071 11.2929C21.0976 11.6834 21.0976 12.3166 20.7071 12.7071L13.7071 19.7071C13.3166 20.0976 12.6834 20.0976 12.2929 19.7071C11.9024 19.3166 11.9024 18.6834 12.2929 18.2929L17.5858 13H4C3.44772 13 3 12.5523 3 12C3 11.4477 3.44772 11 4 11H17.5858L12.2929 5.70711C11.9024 5.31658 11.9024 4.68342 12.2929 4.29289Z" fill="#000000"></path> </g></svg>
    </div>
</div> 



<!-- titulo dias semana abreviados -->
<div class="w-full flex flex-row text-gray-400  {dayWeekFontSize} font-bold text-center {dayLineHeight} {ptWeekDays}  columns-7 justify-around " > 
  <div>{$Terms.sunday_short}</div>
  <div>{$Terms.monday_short}</div>
  <div>{$Terms.tuesday_short}</div>
  <div>{$Terms.wednesday_short}</div>
  <div>{$Terms.thursday_short}</div>
  <div>{$Terms.friday_short}</div>
  <div>{$Terms.saturday_short}</div>
</div>


<!-- dias do mes/ano atuais - 6 linhas para 6 semanas -->
{#each {length: 6} as _, week}

  <div class="w-full flex flex-row  font-bold text-center {dayLineHeight} columns-7 justify-around pt-2 mb-1" >
    {#each {length: 7} as _, day}

      <!-- dias exibidos de mes anterior/posterior, ficarão cor mais fraca -->
      {#if calendar[week][day][1]=='weaker'}
          <div class=" cursor-pointer hover:bg-[#B1DCFB] text-gray-300 hover:text-white w-full flex justify-center items-center {dayFontSize} {dayLineHeight}"
                  on:click={() => calendarDateClicked(calendar[week][day][2])} aria-hidden="true">{ calendar[week][day][0] }</div>

      <!-- dias exibidos de mes atual, ficarão cor normal -->
      {:else if calendar[week][day][1]=='normal'}
          <div class="dayOfCalendar cursor-pointer hover:bg-[#B1DCFB] {currentMonthDaysColor}   w-full flex justify-center items-center {dayFontSize} {dayLineHeight}" 
                  on:click={() => calendarDateClicked(calendar[week][day][2])} aria-hidden="true">{ calendar[week][day][0] }</div>

      <!-- o dia de hoje recebera borda vermelha -->
      {:else if calendar[week][day][1]=='today_color'}
          <div class="dayOfCalendar cursor-pointer border-blue-600  border-2 hover:bg-slate-300 text-gray-500  w-full flex justify-center items-center {dayFontSize} {dayLineHeight}"
                  on:click={() => calendarDateClicked(calendar[week][day][2])} aria-hidden="true">{ calendar[week][day][0] }</div>

      {/if}
    {/each}
  </div>

{/each}


<div class="w-full flex flex-col {buttonsDivHeight}  border-t-2" >

    <div class="w-full flex flex-row h-3 " ></div>

    <!-- botoes semana/mês -->
    <div class="text-gray-500 w-full flex flex-row  text-lg text-center {buttonsContainerHeight} pl-3 pr-3  gap-2 " >

        <div class="btnCalendarBottomButtons {dayFontSize}" id='btnPreviousYear'  on:click={() => {browseCalendar('year',-1)}} aria-hidden="true">&nbsp;</div>
        <div class="btnCalendarBottomButtons {dayFontSize}" id='btnToday'  on:click={() => {browseCalendar('today')}} aria-hidden="true">{$Terms.calendar_today}</div>
        <div class="btnCalendarBottomButtons {dayFontSize}" id='btnNextYear' on:click={() => {browseCalendar('year',1)}} aria-hidden="true">&nbsp;</div>

        {#if showCloseButton}
          <div class="btnCalendarBottomButtons {dayFontSize}" id='btnCLOSE' on:click={btnCloseCalendarClicked} aria-hidden="true">{$Terms.button_cancel}</div>
        {/if}
    </div>

</div>


