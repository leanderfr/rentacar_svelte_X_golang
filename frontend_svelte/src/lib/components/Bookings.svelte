
<script>
import moment from 'moment';

import { Terms, backendUrl, imagesUrl , BookingCalendar_CurrentDate,  imagesStillLoading, clientIp } from '$lib/stores/stores.js'


import {dateToIsoStringConsideringLocalUTC, range, showFormErrors, leadingZeros, logAPI} from '$js/utils.js'

import { slidingMessage, getAndShowCarDetails, selectedCountry, showOrHideTooltip, hourFormat,
       makeWindowDraggable  } from '$js/utils.js'

import RecordMetadataWindow from '$lib/components/RecordMetadataWindow.svelte'

import Calendar from '$lib/components/Calendar.svelte'
import BookingForm from '$edit_forms/BookingForm.svelte';

import { setContext } from 'svelte';

// recordset com os dados de cards de carros
let carOptions

// carro atualmente selecionado no browser de carros
let bookingSelectedCarId 

// exibe/oculta o calendario
let showCalendar = false

// exibe/oculta form de reserva veiculo
let showBookingForm = false

// formHttpMethodApply pode ser: POST, PATCH ou DELETE - o mesmo form de edicao (BookingForm.svelte) é usado para insert/edit/delete
let formHttpMethodApply = ''

// contem o ID do registro que sera manipulado
let currentBookingFormRecordId

// a <div> de reserva pode ser movida ou clicada, qdo esta sendo arrastada, a variavel 'draggingBookingDivYet', desativa temporarimente o evento 'click' 
let draggingBookingDivYet = false

// se showBookingRecordMetadata= exibe tela com metadados 
let showBookingRecordMetadata = false

 // cada chamada API envia nome/id do grupo para que o usuario só manipule dados do seu grupo
let _currentWorkgroupName = localStorage.getItem("rentacar_workgroup_name");


/************************************************************************************************************************************************************
 monta o titulo do calendario
 $BookingCalendar_CurrentDate ==>  inicia sendo= hoje  (em Stores.js) 
 $BookingCalendar_CurrentDate sera atualizada a medida que usuario avançar/retroceder semanas
************************************************************************************************************************************************************/

async function refreshBookingDatesAndContent() { 

  // necessario abrir evento assincrono para exibir div ajax loading, caso contrario navegador nao atualiza a tela
  setTimeout(() => {showLoadingGif(); }, 1);

  let currentDate = new Date($BookingCalendar_CurrentDate.getFullYear(), $BookingCalendar_CurrentDate.getMonth(), $BookingCalendar_CurrentDate.getDate());

  // retrocede ate achar o ultimo domingo antes da data atual (BookingCalendar_CurrentDate)
  while (currentDate.getDay()!=0)  {
    currentDate.setDate(currentDate.getDate() - 1);
  }

  // exibe o dia semana (3 letras)  e a data dd/mm 
  let weekday, weekday_str

  let weekdays =  {0: $Terms.sunday_short, 1: $Terms.monday_short, 2: $Terms.tuesday_short, 3: $Terms.wednesday_short, 
                    4: $Terms.thursday_short, 5: $Terms.friday_short, 6: $Terms.saturday_short} 


  let country = selectedCountry()

  let today = new Date();  // hoje
  let _today_ = new Date(today.getFullYear(), today.getMonth(), today.getDate());

  let options = {month: '2-digit', day: '2-digit'} 

  let displayedYears = []

  let firstDayWeek = new Date(currentDate.getFullYear(), currentDate.getMonth(), currentDate.getDate());
  let lastDayWeek

  for (weekday=0; weekday<7; weekday++) {

    lastDayWeek = new Date(currentDate.getFullYear(), currentDate.getMonth(), currentDate.getDate());

    weekday_str = weekdays[weekday]

    // exibe a data na respectiva <div>
    if ( country == 'usa') 
      jq(`#datecolumn${weekday}`).html( weekday_str + ' ' +currentDate.toLocaleDateString('en-us', options ))    // mm/dd
    else 
      jq(`#datecolumn${weekday}`).html( weekday_str + ' ' + currentDate.toLocaleDateString( 'pt-br', options ))  // dd/mm

    // usa a propriedade (inventada) 'realDate' para memorizar a data real da coluna
    jq(`#datecolumn${weekday}`).attr('real_date', currentDate.toLocaleDateString("fr-CA", {year:"numeric", month: "2-digit", day:"2-digit"})  )    // yyyy-mm-dd

    // marca o dia de hoje em vermelho
    if (currentDate.getTime() == _today_.getTime()) {
      jq(`#datecolumn${weekday}`).css('border-color', 'blue')
      jq(`#datecolumn${weekday}`).attr('today', 'true')
    } else {
      jq(`#datecolumn${weekday}`).css('border-color', 'transparent')
      jq(`#datecolumn${weekday}`).attr('today', 'false')
    }

    // usuario passou/retirou mouse sobre a data no titulo do calendario
    jq(`#datecolumn${weekday}`).on('mouseleave', function()   {       
        if ( jq(this).attr('today')!='true' )   jq(this).css("border-color","transparent")
    });      
    jq(`#datecolumn${weekday}`).on('mouseenter', function()   {       
        if ( jq(this).attr('today')!='true' )   jq(this).css("border-color","black")
        else    jq(this).css("border-color","blue")
    });      

    // concatena ano atual para exibir no titulo do calendario de reservas
    if (displayedYears.indexOf(currentDate.getFullYear())==-1)  displayedYears.push(currentDate.getFullYear()  )

    // avanca 1 dia
    currentDate.setDate(currentDate.getDate() + 1);
  }


  // exibe no titulo do calendario , ano  ou anos atuais sendo visualizado(s), se estiver no inicio/final de um ano, aparecerão ambos os anos envolvidos  
  let _displayedYears_ = ''
  let y
  for (y=0; y<displayedYears.length; y++) {
    if (_displayedYears_!='')  _displayedYears_ += ' / '
    _displayedYears_ += displayedYears[y]
  }

  jq('#currentYear').html(_displayedYears_)

  // carrega as reservas da semana que esta sendo visualizada

  let __firstDayWeek = firstDayWeek.toLocaleDateString("fr-CA", {year:"numeric", month: "2-digit", day:"2-digit"})    
  let __lastDayWeek = lastDayWeek.toLocaleDateString("fr-CA", {year:"numeric", month: "2-digit", day:"2-digit"})    

  try {
      let _route_ = `${$backendUrl}/${_currentWorkgroupName}/bookings/${country}/${bookingSelectedCarId}/${__firstDayWeek}/${__lastDayWeek}`
      logAPI('GET', _route_)
      await fetch(_route_, {method: 'GET'})

      .then( (response) => {

        if (!response.ok) {
          throw new Error(`Bookings Prepare Err Fatal= ${response.status}`);
        }
        return response.json();
      })

      .then( (bookings) => {

        // exemplo json retornado
        /*
        [
            {
                "booking_id": 19,
                "car_id": 7,
                "pickup_formatted": "04/12 09:30",
                "pickup_reference": "2024-12-04|09:30",
                "dropoff_formatted": "05/12 15:00",
                "dropoff_reference": "2024-12-05|15:00",
                "driver_name": "teste 1",
                "car_image": "car_000007.png"
            },
            {
                "booking_id": 20,
                "car_id": 7,
                "pickup_formatted": "04/12 08:55",
                "pickup_reference": "2024-12-04|08:55",
                "dropoff_formatted": "06/12 18:30",
                "dropoff_reference": "2024-12-06|18:30",
                "driver_name": "agora vai",
                "car_image": "car_000007.png"
            }
        ]
        */


        // coloca na propriedade (inventada) 'bookings_this_day', contida no header de cada coluna de dia da semana, 
        // as reservas feitas para aquele dia, dessa forma o jscript vai poder mais adiante exibir <div>'s com as reservas naquela coluna

        // as reservas do dia serao concatenadas em um string, cada reserva separada por '^', e cada campo da reserva, separado por '|'
        // e colocadas na propriedade (inventada) 'bookings_this_day'

        let bookingsThisDay 
        let weekday

        // remove informacao usada antes, sobre quais reservas devem ser exibidas em cada coluna de dia da semana
        for (weekday=0; weekday<7; weekday++) {
          jq(`#datecolumn${weekday}`).attr('bookings_this_day','') 
        }

        let availableColors = [ ['rgb(204, 224, 255)', 'rgb(77, 148, 255)'],
                                    ['rgb(214, 245, 214)', 'rgb(0, 128, 0)'],
                                    ['rgb(255, 235, 204)', 'rgb(255, 165, 0)'],
                                    ['rgb(255, 235, 230)', 'rgb(255, 0, 0)'],
                                    ['rgb(242, 230, 217)', 'rgb(172, 115, 57)'] ]
        let whichColor = 0

        for (let bks = 0; bks < bookings.length; bks++)   {
        
            let pickup = bookings[bks]['pickup_reference'].split('|')[0].split('-')   // pega a data yyyy-mm-dd   (pickup_reference= yyyy-mm-dd|HH:mm)
            let dropoff = bookings[bks]['dropoff_reference'].split('|')[0].split('-')

            let _pickup = new Date(pickup[0], parseInt(pickup[1], 10)-1, pickup[2]);
            let _dropoff = new Date(dropoff[0], parseInt(dropoff[1], 10)-1, dropoff[2]); 

            let pickupHour = bookings[bks]['pickup_reference'].split('|')[1]   // pega a hora HH:mm   (pickup_reference= yyyy-mm-dd|HH:mm)
            let dropoffHour = bookings[bks]['dropoff_reference'].split('|')[1]

            // percorre cada dia da semana e verifica se a reserva atual deve ser exibida no dia
            let currentDay = new Date(firstDayWeek.getFullYear(), firstDayWeek.getMonth(), firstDayWeek.getDate());


            for (weekday=0; weekday<7; weekday++) {

              // se a data da coluna atual esta dentro do intervalo da reserva sendo lida, registra (bookings_this_day) que a <div> da reserva deve aparecer aqui
              // o tamanho da <div> vai depender da qtde de horas reservadas no dia atual
              if (currentDay >= _pickup && currentDay <= _dropoff)   {

                bookingsThisDay = jq(`#datecolumn${weekday}`).attr('bookings_this_day') 

                // se o inicio de reserva nao é no dia atual, considera reservado desde o começo do dia (05:00), pois ela comecou em um dia anterior ao atual
                let startingHour = 5, startingMinute = 0                  

                // se o inicio da reserva é o dia atual, obtem a hora do inicio para informar ao jscript onde iniciar a exibicao da <div> com a reserva
                if (currentDay.getTime() == _pickup.getTime())  {
                  startingHour = parseInt( pickupHour.substring(0,2), 10)  // obtem o numero da hora inicial (HH)
                  startingMinute = parseInt( pickupHour.substring(3), 10)  // obtem o minuto inicial (mm)  
                }

                // se o fim de reserva nao é no dia atual, considera reservado ate o fim do dia (23:59), pois ela terminará somente em um dia posterior ao atual
                // usa a hora invalida 24 para avisar o algoritmo que essa reserva nao terminara no dia atual, só em um proximo dia
                let endingHour = 24, endingMinute = 0
                // se o fim da reserva é o dia atual, obtem a hora do fim para informar ao jscript onde finalizar a exibicao da <div> com a reserva
                if (currentDay.getTime() == _dropoff.getTime())  {
                  endingHour = parseInt( dropoffHour.substring(0,2), 10 )  // obtem o numero da hora final  (HH)
                  endingMinute = parseInt( dropoffHour.substring(3), 10 )  // obtem o minuto da hora final  (mm)
                }

                bookingsThisDay += bookingsThisDay=='' ? '' : '^'
                bookingsThisDay +=  bookings[bks]['booking_id'] + '|' + 
                                    bookings[bks]['car_id'] + '|' +  
                                    bookings[bks]['pickup_formatted'] + '|' + 
                                    bookings[bks]['dropoff_formatted'] + '|' + 
                                    bookings[bks]['driver_name'] + '|' +
                                    startingHour + '|' + startingMinute + '|' +
                                    endingHour + '|' + endingMinute + '|' +
                                    bookings[bks]['car_image'] + '|' +
                                    availableColors[whichColor][0] + '|' + availableColors[whichColor][1]    // bgcolor|border color

                // memoriza que na coluna de data atual, esta reserva deve ser exibida
                jq(`#datecolumn${weekday}`).attr('bookings_this_day', bookingsThisDay) 
              }

              // avanca 1 dia
              currentDay.setDate(currentDay.getDate() + 1);
            }

            whichColor = (whichColor < availableColors.length-1) ? (whichColor + 1) : 0   // pega proxima cor didsponivel para <div>
        }

        // exibe as <div>'s de reserva 
        postItBookingDivs()

    })

  } 
  catch(err) {
    throw new Error(`Bookings Prepare Err Fatal= ${err.message}`);
  }

  // necessario abrir evento assincrono para exibir div ajax loading, caso contrario navegador nao atualiza a tela
  setTimeout(() => {hideLoadingGif()  }, 300);
  
}

/****************************************************************************************************
usuario clicou no botao 'metadados'
****************************************************************************************************/
const invokeBookingMetadataWindow = () =>  {
  showBookingRecordMetadata=true
}

/*****************************************************************************************************************************************************************
// fehca janela modal com metadados do registro
*****************************************************************************************************************************************************************/
export const closeBookingMetadataWindow = () => {
  showBookingRecordMetadata = false 
}


 
/************************************************************************************************************************************************************
usuario escolheu um carro para ver agenda
************************************************************************************************************************************************************/

const showCarBooking = () => {

jq('.carCardBooking').removeClass('carCardBookingClicked')  

// bookingSelectedCarId= -1   == exibir agenda de todos os carros, nao ha card de carro selecionado nesse caso
if (bookingSelectedCarId!=-1)  jq(`#carCard${bookingSelectedCarId}`).addClass('carCardBookingClicked')

refreshBookingDatesAndContent()

}




/************************************************************************************************************************************************************
exibe os post it's (<div>'s) com as reservas em seus respectivos dias
************************************************************************************************************************************************************/

const postItBookingDivs = () => {

  let $bookingsTable = jq('#bookingsTable')
  $bookingsTable.scrollTop(0)


  let divCount = 0

  // remove eventuais <div>'s com reservas criadas anteriormente
  jq('.bookTemporaryDiv').off('click')  
  jq('.bookTemporaryDiv').remove()


  for (let weekday=0; weekday<7; weekday++) {

    /* exemplo de 'bookings_this_day'

    coluna de segunda feira:
      bookings_this_day="Car 1 name|11|03/12/2024 05:00|03/12/2024 11:35|myself the driver|5|0|11|35^Car 2 name|9|03/12/2024 19:00|03/12/2024 19:30|me again driving|19|0|19|30^Titano|7.......

    */

    let bookingsThisDay = jq(`#datecolumn${weekday}`).attr('bookings_this_day') 
    if (bookingsThisDay=='') continue

    let bookings = bookingsThisDay.split('^')   // separador de reservas

    // memoriza IDs de reservas para auxiliar na montagem do evento 'draggable' mais abaixo
    let bookingsIDs = []

    // lê as reservas para a coluna do dia de semana atual
    for (let bks=0; bks < bookings.length; bks++)  {

      let booking = bookings[bks].split('|')   // separador de campos da reserva

      let startingHour = booking[5]
      let startingMinute = booking[6]

      let endingHour = booking[7]
      let endingMinute = booking[8]

      // converte as horas/minutos para suas respectivas ROWS <div>'s da agenda
      let tableRowBookingTop = startingHour - 5    // 5= hora inicial do dia, 05:00
      let tableRowBookingBottom = endingHour - 5

      // a <div> que exibe a agenda (bookingsTable), possui <div>'s filhas que sao as horas do dia (05:00 - 23:00)
      // e cada hora do dia possui <div>'s filhas que sao os dias da semana (segunda - sexta)

      let $divBookingTop = $bookingsTable.children().eq(tableRowBookingTop).children().eq(weekday+1)      // weekday+1 por causa da 1a coluna (hora)
      let $divBookingBottom = $bookingsTable.children().eq(tableRowBookingBottom).children().eq(weekday+1)


      let bookDivWidth = $divBookingTop.width() 

      // endingHour= 24 significa que a reserva nao termina no dia de hoje.. continuar para proximo(s) dia(s)

      let pickupMoment = booking[2]
      let dropoffMoment = booking[3]
                        

      let carImage = booking[9]
      let driverName = booking[4]

      let bookingHtml = `<div style='display:flex;flex-direction:column;'>`+
                        `   <div style='display:flex;flex-direction:row;justify-content: space-between;margin-bottom:20px;align-items: center;;'>`+
                        `       <div style="background-repeat: no-repeat;background-size: contain;width:80%;height:50px;background-image: url('${$imagesUrl}${carImage}');margin-bottom:20px "></div>`+
                        `       <div class='bookingDivDrag'>&nbsp;</div>`+
                        `   </div>`+
                        `   <div style='display:flex;flex-direction:row;;margin-bottom:20px;align-items: center;;'>`+
                        `       <div class='bookingStartingHourInfo'>&nbsp;</div>`+
                        `       <div style='padding-top:8px;padding-left:10px' >${pickupMoment}</div>`+
                        `   </div>`+
                        `   <div style='display:flex;flex-direction:row;;margin-bottom:20px;align-items: center ;;'>`+
                        `       <div class='bookingEndingHourInfo'>&nbsp;</div>`+
                        `       <div style='padding-top:8px;padding-left:10px' >${dropoffMoment}</div>`+
                        `   </div>`+
                        `   <div style='display:flex;flex-direction:row;;margin-bottom:20px;align-items: center ;;'>`+
                        `       <div class='bookingDriverInfo'>&nbsp;</div>`+
                        `       <div style='padding-top:8px;padding-left:10px' >${driverName}</div>`+
                        `   </div>`+
                        `</div>`


      // <div> de reserva colocada na ultima coluna (ult dia semana), necessario diminuir 2 pixels, se nao browser esconde a borda
      if (weekday==6) bookDivWidth -= 2;

      // calcula altura da <div> da reserva, considerando a distancia entre a hora final/inicial da reserva
      // 60 pixels, o tamanho da <div> que cada hora tem
      let bookDivHeight = (tableRowBookingBottom - tableRowBookingTop) * 62
      bookDivHeight -= parseInt(startingMinute, 10)
      bookDivHeight += parseInt(endingMinute, 10)

      /* cria div  */
      let $divBOOKING = jq("<div />").css({
          position: "absolute",
          overflow:'hidden',
          padding: '5px',
          'border-radius': '5px',
          cursor: 'pointer',          
          backgroundColor: booking[10],                  // alterna cor da reserva para diferenciar
          border: `solid 2px ${booking[11]}`,
          height: bookDivHeight,
          width: bookDivWidth ,
      }).appendTo( $bookingsTable );

      $divBOOKING.attr('id', `bookTemporaryDiv${divCount}`);
      $divBOOKING.addClass('bookTemporaryDiv')    // bookTemporaryDiv nao possui CSS, serve somente para identificar que é <div> de reserva

      let bookingId = booking[0]
      $divBOOKING.attr('booking_id', bookingId)  // une as <div>'s que dizem respeito à mesma reserva, para que qdo usuario passar mouse sobre, coloque destauqe sobre todas ao mesma tempo

      $divBOOKING.html( bookingHtml )

      // se clicar na <div> de reserva, abre edicao
      $divBOOKING.on('click', function(e)   { 
        if (! draggingBookingDivYet)   editBookingRecord(e, jq(this).attr('booking_id') )
      })
        
     
      // qdo usuario passar/retirar mouse de <div>'s que se referem à mesma reserva, coloca/retira destaque de todas ao mesmo tempo, dando a impressao de serem a mesma
      $divBOOKING.mouseenter(function()  {
        jq( `[booking_id=${bookingId}]`).css('border', 'solid 2px black')
        jq( `[booking_id=${bookingId}]`).css('background-color', '#ffffcc')
      });
      $divBOOKING.mouseleave(function()  {
        jq( `[booking_id=${bookingId}]`).css('border', `solid 2px ${booking[11]}`)
        jq( `[booking_id=${bookingId}]`).css('background-color', `${booking[10]}`  )
      })

      // memoriza os IDs das reservas para futuro agrupamento de <div>'s  mais abaixo
      // exemplo:  uma reserva de 2, 3 dias... ela sera exibida em 2, 3 <div>'s - estas divs serao movidas (draggable) e seram destacadas (hover) ao mesmo tempo  

      if ( bookingsIDs.indexOf(bookingId) == -1 ) bookingsIDs.push( bookingId )


      let bookTopHourDivPosition = document.getElementById( $divBookingTop.attr('id') ).getBoundingClientRect()   // captura posicao da <div> hora inicio da reserva
      let bookingsTablePosition = document.getElementById( 'bookingsTable' ).getBoundingClientRect()   // captura posicao da <div> container das reservas

      $divBOOKING.css("left", bookTopHourDivPosition.left - bookingsTablePosition.left + 2);   // posiciona a <div> da reserva no seu dia/horario 

      // se a hora inicial da reserva nao é hora cheia, possui minutos, considerar estes minutos no posicionamento da <div> da reserva
      // e sendo que a <div> da hora possui exatos 60 pixels, basta somar a qtde de minutos (startingMinute) à posicao top da <div> reserva

      let divTopPosition = parseInt(bookTopHourDivPosition.top, 10) - parseInt(bookingsTablePosition.top, 10)
      divTopPosition += parseInt(startingMinute , 10)    // minutos = pixels

      $divBOOKING.css("top", divTopPosition);

      divCount++

    }

    // todas as <div>'s que se referem à mesma reserva, serao movidas (draggable) juntas
    // exemplo:  uma reserva de 2, 3 dias... ela sera exibida em 2, 3 <div>'s - estas divs serao movidas (draggable) e seram destacadas (hover) ao mesmo tempo  
    setTimeout(() => {
        for (let ids=0; ids < bookingsIDs.length; ids++)  {
          jq( `[booking_id=${bookingsIDs[ids]}]`).multiDraggable({ 
            group: jq(`[booking_id=${bookingsIDs[ids]}]`),

            // impede disparo de 'click' na div reserva enqto estiver arrastando a div
            startNative: function (event,ui) {draggingBookingDivYet = true},
            stopNative : function (event,ui) { setTimeout(() => {draggingBookingDivYet = false}, 100); },
          });
        }
    }, 100);


  }
}


/************************************************************************************************************************************************************
usuario pediu para ver metadata (data criacao, id criador, etc) do registro de reserva
************************************************************************************************************************************************************/
export async function fillRecordMetadataWindow (record_id) {

  // necessario abrir evento assincrono para exibir div ajax loading, caso contrario navegador nao atualiza a tela
  setTimeout(() => {
    showLoadingGif();   // mostra animacao 'processando..' mesmo sem tela ter sido renderizada
  }, 1);

  let country = selectedCountry()

  // necessario passar o país para que as datas (created_at, etc) retornem no formato usado naquele país
  const url = new URL(`${$backendUrl}/bookings/record_metadata/${record_id}/${country}/`);
  
  logAPI('GET', url)
  await fetch(url, {method: 'GET'})

  .then((response) => {
    if (!response.ok) {
      throw new Error(`Read Record Metadata Err fatal= ${response.status}`);
    }
    return response.json();
  })

  .then((result) => {
    setTimeout( () => {hideLoadingGif();}, 500)

    let fields = result

    // joga campo json recebido para seu respectivo <span>
    // exemplo:  json: created_at  ==>   <span id='info_created_at'>
    Object.keys(fields).forEach(function(key, idx, arr)   {
      jq(`#info_${key}`).html( fields[key] )      
    })
  })
}



/************************************************************************************************************************************************************
recebe JSON que foi obtido via FETCH API e popula elementos do form de reserva de veiculo
************************************************************************************************************************************************************/
export async function populateBookingFormDestinationFields(pairs) {

  jq('#txtPickUpDate').val( pairs['pickup_date'] )
  jq('#txtPickUpHour').val( pairs['pickup_hour'] )
  jq('#txtDropOffDate').val( pairs['dropoff_date'] )
  jq('#txtDropOffHour').val( pairs['dropoff_hour'] )
  jq('#txtDriverName').val( pairs['driver_name'] )

  getAndShowCarDetails( pairs['car_id'] )

}


/************************************************************************************************************************************************************
avanca/retorna semana no calendario
************************************************************************************************************************************************************/
const browseBookingCalendar = (days) => {
  // obtem o mês/ano da data exibida atualmente para avançar/retroceder
  let tmpDate = new Date($BookingCalendar_CurrentDate.getFullYear(), $BookingCalendar_CurrentDate.getMonth(), $BookingCalendar_CurrentDate.getDate());
  tmpDate.setDate(tmpDate.getDate() + days);

  $BookingCalendar_CurrentDate = tmpDate

  refreshBookingDatesAndContent()
}


/************************************************************************************************************************************************************
monta lista de carros que ficara no canto direito da tela
************************************************************************************************************************************************************/
async function mountCarList () {

  // necessario abrir evento assincrono para exibir div ajax loading, caso contrario navegador nao atualiza a tela
  setTimeout(() => {showLoadingGif(); }, 1);

  bookingSelectedCarId = 0

  let country = selectedCountry()

  try {
      let _route_ = `${$backendUrl}/${_currentWorkgroupName}/car_cards/${country}`
      logAPI('GET', _route_)
      await fetch(_route_, {method: 'GET'})

      .then( (response) => {

        if (!response.ok) {
          throw new Error(`Car Cards Prepare Err Fatal= ${response.status}`);
        }
        return response.json();
      })

      .then( (carOptionsData) => {
        // obtem o ID do primeiro carro, para colocar destaque no mesmo e ja começar mostrando suas propriedades
        // o usuario pode mudar o carro atualmente selecioado (bookingSelectedCarId) clicando sobre o respectivo 
        bookingSelectedCarId = carOptionsData[0].id

        carOptions = carOptionsData    

        setTimeout( () => {

          // define tooltip dos botoes parte superior
          // putWhiteTooltip sao <div>'s com detalhes do veiculo, mas cuja tooltip deve aparecer em branco, mais discreto (em Bookings.svelte)
          if (typeof jq('.putWhiteTooltip').tooltip !== "undefined") {
            // define parte visual e o conteudo (title) do elemento
            jq('.putWhiteTooltip').tooltip({
              tooltipClass: 'prettierTitle_white',
              show: false,   // sem animacao ao exibir 
              hide: false,   // sem animacao ao ocultar
              position: { my: "left top", at: "left top-60", collision: "flipfit" }
            })

            // tooltip para cards de carro (lado direito da agenda)
            jq('.carCardBooking').tooltip({ 
              tooltipClass: 'prettierTitle_blue',
              show: false,   // sem animacao ao exibir 
              hide: false,   // sem animacao ao ocultar
              position: { my: "left top", at: "left-150 top", collision: "flipfit" }
            })
          }


          setTimeout(() => {hideLoadingGif(); }, 1);          
          setTimeout(() => {showCarBooking() ; }, 100);
        }, 500)    
      })
  } 
  catch(err) {
    throw new Error(`Car Cards Prepare Err Fatal= ${err.message}`);
  }
}


/************************************************************************************************************************************************************
formulario para nova reserva
************************************************************************************************************************************************************/
const newBookingRecord = () => {
setTimeout(() => {showLoadingGif(); }, 1);
formHttpMethodApply  = 'POST'
showBookingForm = true
}

/****************************************************************************************************
 abre form para edicao de registro
****************************************************************************************************/
const editBookingRecord = (event, _id_) =>  {
  event.stopPropagation();
  formHttpMethodApply = 'PATCH'
  currentBookingFormRecordId = _id_

  setTimeout(() => {showLoadingGif() }, 1);
  showBookingForm=true
}

// se mudar o idioma/país, recarrega lista de carros
// a linha abaixo faz disparar a funcao 'mountCarList()',  
// o svelte dispara qq funcao vinculada à alteracao do valor de $Terms (por ser uma declaracao reativa), 
// e qdo Bookings.svelte é carregado, svelte considera mudança de valor da variavel, e dispara automatico, sem necessidade de criar onMount()
$: $Terms, mountCarList()
/*****************************************************************************************************************************************************************
// funcao acionada qdo usuario clica na div '#backDrop' ou clica em um dos botoes de fechar calendario
*****************************************************************************************************************************************************************/
export const closeAnyBookingModalWindow = () => {
  showCalendar = false    // fecha janela modal
  showBookingForm = false   // fecha form de reserva veiculo
}

/********************************************************************************************************************************************************
 exclui o registro atual sendo exibido no form edicao
********************************************************************************************************************************************************/
const performDeleteBookingRecord = () =>  {
  performBookingCrudOperation(null, 'DELETE', `bookings/delete/${currentBookingFormRecordId}`)  
}

/********************************************************************************************************************************************************
 valida dados do formulario e se tudo ok, tenta gravar
********************************************************************************************************************************************************/
const performSaveBookingRecord = () =>  {

  // remove qq msg erro que foi exibida previamente
  jq('.errorTextbox').removeClass('errorTextbox').addClass('noerrorTextbox')      // cada div que exibe erro, e possui a classe 'errorTextbox'

  let errors = []
  let country = selectedCountry()

  // valida datas / horas de retirada/devolucao do carro 

  //******************************************************************************************/
  // pick up date
  //******************************************************************************************/
  let _txtPickUpDate_ = jq('#txtPickUpDate').val()
  let txtPickUpDate = _txtPickUpDate_.split('/')  
  let pickUpIso8601Format

  if (country == 'usa')  
    pickUpIso8601Format = '20'+txtPickUpDate[2] +'-' + txtPickUpDate[0] + '-'+txtPickUpDate[1]    // yyyy-mm-dd

  if (country == 'brazil')   
    pickUpIso8601Format = '20'+txtPickUpDate[2] +'-' + txtPickUpDate[1] + '-'+txtPickUpDate[0]    // yyyy-mm-dd


  if (! moment(pickUpIso8601Format).isValid())     errors.push('txtPickUpDate')

  //******************************************************************************************/
  // pick up hour
  //******************************************************************************************/

  let regex_usa = /^([0]\d|[1][0-2]):([0-5]\d)\s?(?:AM|PM)$/i;          
  let regex_brazil = /^(?:[01][0-9]|2[0-3]):[0-5][0-9](?::[0-5][0-9])?$/;

  let txtPickUpHour = jq('#txtPickUpHour').val()

  if   (   (country == 'usa'  &&  ! regex_usa.test(txtPickUpHour))     ||     (country == 'brazil' &&  ! regex_brazil.test(txtPickUpHour))    )
    errors.push('txtPickUpHour')

  

  
  //******************************************************************************************/
  // drop off date
  //******************************************************************************************/

  let _txtDropOffDate_ = jq('#txtDropOffDate').val()
  let txtDropOffDate = _txtDropOffDate_.split('/')  
  let dropOffIso8601Format

  if (country == 'usa')  
    dropOffIso8601Format = '20'+txtDropOffDate[2] +'-' + txtDropOffDate[0] + '-'+txtDropOffDate[1]    // yyyy-mm-dd

  if (country == 'brazil')    
    dropOffIso8601Format = '20'+txtDropOffDate[2] +'-' + txtDropOffDate[1] + '-'+txtDropOffDate[0]    // yyyy-mm-dd

  if (! moment(dropOffIso8601Format).isValid())     errors.push('txtDropOffDate')


  //******************************************************************************************/
  // drop off hour
  //******************************************************************************************/

  let txtDropOffHour = jq('#txtDropOffHour').val()

  if   (   (country == 'usa'  &&  ! regex_usa.test(txtDropOffHour))     ||     (country == 'brazil' &&  ! regex_brazil.test(txtDropOffHour))    )
    errors.push('txtDropOffHour')

  // nome do motorista
  if ( jq('#txtDriverName').val().trim().length < 3 )  errors.push('txtDriverName')


  // se encontrou erros, exibe
  if (Object.keys(errors).length>0) {showFormErrors(errors) ;  return;}


  var formData = new FormData(); 
  formData.append('country', country)    // alguns tabelas precisam do país selecionado atualmente - exemplo: tabela carros, grava/separa os carros por país 


  let pickUpHour, pickUpMinute, dropOffHour, dropOffMinute, pickupAlmostReady, dropoffAlmostReady

  
  // converte hora formato USA 12 AM/PM para formato sql HH:MM
  if (country == 'usa')  {
    pickUpHour = moment(txtPickUpHour, ["h:mm A"]).format("HH");
    pickUpMinute = moment(txtPickUpHour, ["h:mm A"]).format("mm");

    dropOffHour = moment(txtDropOffHour, ["h:mm A"]).format("HH");
    dropOffMinute = moment(txtDropOffHour, ["h:mm A"]).format("mm");

    pickupAlmostReady = new Date('20'+txtPickUpDate[2], parseInt(txtPickUpDate[0], 10)-1, txtPickUpDate[1], pickUpHour, pickUpMinute)   // mm/dd/yy
    dropoffAlmostReady = new Date('20'+txtDropOffDate[2], parseInt(txtDropOffDate[0], 10)-1, txtDropOffDate[1], dropOffHour, dropOffMinute)
  }

  if (country == 'brazil')  { 
    pickUpHour = moment(txtPickUpHour, ["HH:mm"]).format("HH");
    pickUpMinute = moment(txtPickUpHour, ["HH:mm"]).format("mm");

    dropOffHour = moment(txtDropOffHour, ["HH:mm"]).format("HH");
    dropOffMinute = moment(txtDropOffHour, ["HH:mm"]).format("mm");

    pickupAlmostReady = new Date('20'+txtPickUpDate[2], parseInt(txtPickUpDate[1], 10)-1, txtPickUpDate[0], pickUpHour, pickUpMinute)  // dd/mm/yy
    dropoffAlmostReady = new Date('20'+txtDropOffDate[2], parseInt(txtDropOffDate[1], 10)-1, txtDropOffDate[0], dropOffHour, dropOffMinute)
  }

  // verifica se data devolucao é maior do que data retirada
  let datesDifference = ( dropoffAlmostReady - pickupAlmostReady ) / 36e5;
  if (datesDifference < 0) {
    slidingMessage('slidingFormMessage', $Terms.booking_dates_difference_error, 4000)
    jq('#txtDropOffDate').focus() 
    return
  }

  // reserva deve ser feita com pelo menos 1 hora de antecedencia
  let _minimumDate = new Date();  
  let minimumDate = _minimumDate.getTime() + (1*60*60*1000);   // adiciona 1 hora

  let datesDifference1 = ( dropoffAlmostReady - minimumDate ) / 36e5;
  let datesDifference2 = ( pickupAlmostReady - minimumDate ) / 36e5;

  if (datesDifference1 < 0 || datesDifference2 < 0) {
    slidingMessage('slidingFormMessage', $Terms.booking_dates_not_in_advance, 4000)
    jq('#txtDropOffDate').focus() 
    return
  }


  formData.append('dropoff_datetime', dateToIsoStringConsideringLocalUTC(dropoffAlmostReady))    // formato iso8601 datetime
  formData.append('pickup_datetime', dateToIsoStringConsideringLocalUTC(pickupAlmostReady))    // formato iso8601 datetime
  formData.append('driver_name', jq('#txtDriverName').val())
  formData.append('car_id', bookingSelectedCarId)


  let route = `${_currentWorkgroupName}/`

  if (formHttpMethodApply=='POST') 
    route += 'booking'        
  if (formHttpMethodApply=='PATCH') 
    route += `booking/${currentBookingFormRecordId}`  

  // informa IP do cliente para que nao exiba notificacoes para quem fez a alteracao, 
  // somente para quem pertence ao grupo mas nao fez a alteracao
  formData.append('client_ip', $clientIp)    

  // formHttpMethodApply= POST, PATCH ou DELETE
  performBookingCrudOperation(formData, formHttpMethodApply, route)  
}


// controla qual componente cv
let whichComponentSvelteCalledCalendar = 'Booking'




/**********************************************************************************************
a funcao abaixo insere, edita (update), apaga (delete) reservas de veiculos
**********************************************************************************************/

const performBookingCrudOperation = async (body, formHttpMethodApply, route) => {

  // mostra animacao indicando pausa para processamento
  setTimeout(() => {showLoadingGif();}, 1);

  let successmsg
  let errormsg

  // prepara as msgs de erro ou sucesso, para cada tipo de operacao 
  switch (formHttpMethodApply) {

    case 'GET':
      // se leitura ocorrer ok, nao precisa mostrar aviso, nao existe 'Terms.successmsg_record_read'
      errormsg = $Terms.errormsg_reading_record
      break

    case 'POST':
      successmsg= $Terms.successmsg_record_inserted
      errormsg = $Terms.errormsg_saving_record
      break

    case 'PATCH':
      successmsg= $Terms.successmsg_record_updated
      errormsg = $Terms.errormsg_saving_record
      break

    case 'DELETE':
      successmsg= $Terms.successmsg_record_deleted
      errormsg = $Terms.errormsg_deleting_record
      break
  }

  try {
    const url = new URL(`${$backendUrl}/${route}`);

    let options

    // metodo 'GET' nao possui body, parametros ja vieram na URL (route)
    if (formHttpMethodApply=='GET')      
      options = {method: 'GET'}

    // metodo DELETE nao possui body, os parametros sao passados via URL  
    // metodos POST e PATCH possuem body (form-data com dados para gravacao)
    else {
      if ( body != null )  {
        // infelizmente o PHP, ate a versao 8.2, nao trabalha bem com metodo PATCH, necessario mudar para 'POST'
        // a funcao php, request_parse_body(), que promete ler body qdo metodo= PATCH() nao funciona na versao PHP do backend (7.4)
        // request_parse_body() so funciona na versao 8.4

        // os demais backends manipulam PATCH sem problemas
        let backend = localStorage.getItem("rentacar_current_backend")
        if (backend=='php' && formHttpMethodApply=='PATCH')   
          options = {method: 'POST', body}
        else
          options = {method: formHttpMethodApply, body}
      }
      else
        // metodo PATCH com parametros passados na URL mesmo
        options = {method: formHttpMethodApply}     
    }
    
    // chama backend para GET ou POST
    logAPI(formHttpMethodApply, url)
    await fetch(url, options)

    .then((response) => {
      if (formHttpMethodApply=='GET') return response.json()   // read= retorno é json contento campos do registro
      else return response.text()     // delete, patch, post = retorno é um texto informando se operacao teve sucesso ou nao
    })                      

    // se executou backend com sucesso
    .then((result) => {
      setTimeout(() => {
        hideLoadingGif(); 
      },  1);

      // read= sucesso ao ler dados do registro
      if (formHttpMethodApply=='GET') { 

          let pairs = result

          // preenche os campos da reserva com o conteudo JSON recebido
          populateBookingFormDestinationFields(pairs)
          putFocusInFirstInputText_AndOthersParticularitiesOfTheBookingForm();          
      }

      // != read= sucesso ao gravar dados do registro (update, insert, delete, delete_selected)
      else {
        let text = result

        // se gravou ok
        if (text.indexOf('__success__')!=-1)  {

          slidingMessage('slidingWindowMessage', successmsg, 2000)   

          // coloca como data inicial da agenda o dia da retirada do veiculo (pickupDate)
          let _txtPickUpDate_ = jq('#txtPickUpDate').val()
          let txtPickUpDate = _txtPickUpDate_.split('/')  

          if (selectedCountry()=='usa')  
            $BookingCalendar_CurrentDate = new Date('20'+txtPickUpDate[2], parseInt(txtPickUpDate[0], 10)-1, txtPickUpDate[1]);   // mm/dd/yy
          else 
            $BookingCalendar_CurrentDate = new Date('20'+txtPickUpDate[2], parseInt(txtPickUpDate[1], 10)-1, txtPickUpDate[0]);   // dd/mm/yy

          // fecha formulario e atualiza agenda
          closeAnyBookingModalWindow()
          refreshBookingDatesAndContent()  

          // se usuario alterou algum dado que necessita uma funcao a mais apos sua alteracao, executa esta funcao
          // exemplo: usuario alterou alguma expressao ingles/portugues, necessario recarregar todas as expressoes para que svelte atualize a expressao alterada na tela
          //if (specificFunctionCallWhenSuccess != null)  eval(specificFunctionCallWhenSuccess)

        // backend retornou problema
        } else  {
          // se formulario nem foi aberto, usa a div maior (window) para informar erro
          // se formulario aberto, usa a div interna do formulario

          slidingMessage('slidingFormMessage', errormsg + '&nbsp;&nbsp;=> '+text, 4000)   
          return('error')
        }
      }
    })

    // se retornou erro, erro execucao da API
    .catch(function (error)   { 

      // se GET, fecha form pois nao foi possivel ler campos

      setTimeout(() => {hideLoadingGif();}, 1); 

      if (formHttpMethodApply=='GET') {  
        setTimeout(() => {closeAnyBookingModalWindow()}, 5500); 
        slidingMessage('slidingWindowMessage', $Terms.errormsg_general_database + '&nbsp;&nbsp;Error position 2= '+error.message, 5000)   
      }
      // erro gravacao/exclusao, avisa dentro do proprio form o erro e mantem form aberto
      // usa div rolante no form aberto para informar erro 
      else {
        slidingMessage('slidingFormMessage', errormsg + '&nbsp;&nbsp;Error: situation 2= '+error.message, 5000)   
        return('error')
      }

    });


  // se erro conexao, erro jscript
  } catch (error) {

    setTimeout(() => {hideLoadingGif();}, 1);

    // se GET, fecha form pois nao foi possivel ler campos
    // GET, changingStatus e batchDeleting usa div rolante na janela principal

    if (formHttpMethodApply=='GET' || changingStatus || batchDeleting) { 
      setTimeout(() => {closeCurrentModalWindow()}, 5500); 
      slidingMessage('slidingWindowMessage', $Terms.errormsg_general_database + '&nbsp;&nbsp;Error position 3= '+error.message, 5000)
    }
    // erro gravacao/exclusao, avisa dentro do proprio form o erro e mantem form aberto
    // usa div rolante no form aberto para informar erro
    else {
      slidingMessage('slidingFormMessage', errormsg + '&nbsp;&nbsp;Error situation 3= '+error.message+'-'+formHttpMethodApply, 5000)    
      return('error')
    }
  }

}

/************************************************************************************************************************************************************
coloca foco no 1o input textbox e prepara alguns campos para digitacao
************************************************************************************************************************************************************/
export const putFocusInFirstInputText_AndOthersParticularitiesOfTheBookingForm = () => { 

  // faz o form ser deslocavel
  makeWindowDraggable('divWINDOW_TOP', 'bookingRecordForm') 

  // prepara mascara de digitacao das datas/horas
  if ( selectedCountry()=='usa') {
    jq.mask.definitions['h'] = "[A^Pa^p]"
    jq('#txtPickUpHour, #txtDropOffHour').mask('99:99 hm', {placeholder:"hh:mm _m "})    // data formato USA
  } else {
    jq('#txtPickUpHour, #txtDropOffHour').mask('99:99', {placeholder:"hh:mm"})           // data formato Brasil  
  }

  jq('#txtPickUpDate, #txtDropOffDate').mask('99/99/99');   

  $imagesStillLoading = false

  setTimeout(() => {
    let firstInputText = jq('#bookingRecordForm').find('input[type=text]').filter(':visible:first');
    firstInputText.focus();  
    firstInputText.select();          
    showOrHideTooltip();
    hideLoadingGif()   // esconde div 'carregando...'  caso nao tenha sido escondida ainda
  }, 300);

}


/********************************************************************************************************************************************************
prepara funcionalidades iniciais do formulario de reserva (BookingForm.svelte) invocado por Bookings.svelte e 
dispara funcao que carrega os dados do registro sendo editado
********************************************************************************************************************************************************/

const getBookingFormPopulatedAndReady = async () => { 

  /* ajustes iniciais no form */  

  $imagesStillLoading = true    // mostra animacao gif 'processando'.. durante o carregamento de imagens

  /* popula o form  */  

  // aguarda 300 milisegundos para que os ajustes jscript acima fiquem prontos
  setTimeout(() => { 

      // se vai manipular um registro ja existente (!=POST), carrega seus dados primeiro
      // POST = novo registro
      // != POST= PATCH, DELETE
      if (formHttpMethodApply!='POST')    {
        performBookingCrudOperation(null, 'GET', `booking/${currentBookingFormRecordId}`)     // putFocusInFirstInputText_AndOthersParticularitiesOfTheForm() sera executada dentro de performCrudOperation
      }

      // vai adicionar registro , campos ficam em branco , coloca foco no 1o campo 
      if (formHttpMethodApply=='POST')    {
          // exibe detalhes do carro que sera reservado
          getAndShowCarDetails(bookingSelectedCarId)   // foto carro, tipo cambio, qtde portas, etc  
          //$imagesStillLoading = false
          putFocusInFirstInputText_AndOthersParticularitiesOfTheBookingForm()
      }
  }, 300);
   
}


// disponibiliza ao componente BookingForm.svelte funcoes de CRUD
setContext('offerCrudOperationsToChildComponent', { 
  performSaveBookingRecord,  
  putFocusInFirstInputText_AndOthersParticularitiesOfTheBookingForm,
  performBookingCrudOperation,
  getBookingFormPopulatedAndReady,
  performDeleteBookingRecord,
  invokeBookingMetadataWindow,
  fillRecordMetadataWindow
});



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


<!-- usuario pediu para exibir metadados do registro -->
{#if showBookingRecordMetadata}      

  <!-- backDropMetadata é a div fumê que ficara por trás da janela modal de  metadados
  backDropMetadata impede que o usuario clique em elementos que estão atrás da janela modal
  qdo backDrop for clicada, fecha o form edicao -->

  <div id='backDropMetadata' class='w-full h-full  absolute flex items-center justify-center left-0 top-0 z-20 bg-[rgba(0,0,0,0.5)]' on:click|self={closeBookingMetadataWindow} aria-hidden="true"  >  

      <RecordMetadataWindow  on:dispatchCloseForm={closeBookingMetadataWindow}  currentFormRecordId={currentBookingFormRecordId}  metadataTableName={$Terms.table_bookings}  />     

  </div>

{/if}      


{#if showCalendar}      
  <!-- backDrop é uma div fumê que ficara por trás do calendario
  backDrop impede que o usuario clique em elementos que estão atrás da janela modal
  qdo backDrop for clicada, fecha a janela modal <Calendar /> -->

  <div id='backDrop' class='w-full h-full  absolute flex items-center justify-center left-0 top-0 z-10 bg-[rgba(0,0,0,0.5)]' on:click|self={closeAnyBookingModalWindow} aria-hidden="true"  >  

     <div class="w-[40%] h-auto flex align-top text-lg items-start justify-start  flex-col  bg-white  border-2 border-white"  > 
      <Calendar  {whichComponentSvelteCalledCalendar} on:dispatchCloseCalendar={closeAnyBookingModalWindow}  on:dispatchrefreshBookingDatesAndContent={refreshBookingDatesAndContent}   /> 
     </div>

  </div>

{/if}


{#if showBookingForm}       
  <!-- backDrop é uma div fumê que ficara por trás do form
  backDrop impede que o usuario clique em elementos que estão atrás da janela modal
  qdo backDrop for clicada, fecha a janela modal <BookingForm /> -->

  <div id='backDrop' class='w-full h-full  absolute flex items-center justify-center left-0 top-0 z-10 bg-[rgba(0,0,0,0.5)]' on:click|self={closeAnyBookingModalWindow} aria-hidden="true"  >  
      <BookingForm  {formHttpMethodApply} on:dispatchCloseBookingForm={closeAnyBookingModalWindow}         /> 
  </div>

{/if}




<div class="flex flex-col w-full h-[95%] border-l-red-200  ">

  <div class="flex flex-row h-[60px] w-[85%] justify-between border-b-2">

    <!-- ano sendo exibido -->
    <div class="flex flex-row text-2xl font-bold pt-4" id='currentYear'></div>

    <!-- botoes de acoes -->
    <div class="flex flex-row">
        <!-- botao novo registro -->
        <div  class='btnBOOKING_ADD_CAR_RESERVATION putWhiteTooltip'  title={$Terms.datatable_new_booking}   on:click={ () => {newBookingRecord()}} aria-hidden="true"></div>   

        <!-- botao exbir calendario -->
        <div  class='btnBOOKING_CALENDAR putWhiteTooltip' title={$Terms.choose_date} on:click={ () => {showCalendar=true}} aria-hidden="true"></div>    

        <!-- botao exbir agenda de todos os carros -->
        <div  class='btnBOOKING_ALL_CARS putWhiteTooltip' title={$Terms.all_cars_booking} on:click={() => {bookingSelectedCarId=-1; showCarBooking()}} aria-hidden="true"></div>    

        <!-- botao retroceder semana (seta esquerda) -->
        <div  class='btnBOOKING_LEFT_ARROW putWhiteTooltip'  title={$Terms.previous_week} on:click={() => {browseBookingCalendar(-7)}} aria-hidden="true"></div>   

        <!-- botao avançar semana (seta direita) -->
        <div  class='btnBOOKING_RIGHT_ARROW putWhiteTooltip' title={$Terms.next_week} on:click={() => {browseBookingCalendar(+7)}} aria-hidden="true"></div>   
    </div> 

  </div>

  <div class="flex flex-row h-[calc(100%-60px)]  w-full  gap-x-[1px] ">

      <!-- lado esquerdo, agenda e detalhes do veiculo escolhido -->
      <div class="flex flex-col w-[85%]  border-b-2  h-full ">

            <!-- dias da semana  -->
            <div class="w-full border-b-2 border-b-gray-300 text-lg py-1 grow-0" >  
                <div class="w-[calc(100%-22px)] flex flex-row text-gray-500  text-lg font-bold text-center h-12 justify-center cursor-pointer items-center" >
                  <div class='w-[9%] tdBookingHeader'>&nbsp;</div>
                  <div class='w-[13%] tdBookingHeader rounded-2xl' id='datecolumn0' bookings_this_day='' real_date=''></div> 
                  <div class='w-[13%] tdBookingHeader rounded-2xl' id='datecolumn1' bookings_this_day='' real_date=''></div>
                  <div class='w-[13%] tdBookingHeader rounded-2xl' id='datecolumn2' bookings_this_day='' real_date=''></div>
                  <div class='w-[13%] tdBookingHeader rounded-2xl' id='datecolumn3' bookings_this_day='' real_date=''></div>
                  <div class='w-[13%] tdBookingHeader rounded-2xl' id='datecolumn4' bookings_this_day='' real_date=''></div> 
                  <div class='w-[13%] tdBookingHeader rounded-2xl' id='datecolumn5' bookings_this_day='' real_date=''></div>
                  <div class='w-[13%] tdBookingHeader rounded-2xl' id='datecolumn6' bookings_this_day='' real_date=''></div>
                </div>
            </div>

            <!-- looping para listar os horarios de 05:00 ate 23:00 -->
            <div class="w-full flex flex-col  overflow-y-auto  flex-1 border-l-0 border-gray-200 border-r-0 relative" id='bookingsTable' >  
                  {#each range(5, 24, 1) as hour, index}

                    <div class="w-full flex flex-row  leading-[60px]  justify-center cursor-pointer border-b-2 border-gray-300 hover:bg-gray-100"  >
                      <div class='w-[9%] tdBookingCell flex justify-center'>{hourFormat(hour)}</div>
                      <div class='w-[13%] tdBookingCell' id='bookingHourDay0{hour}' ></div>
                      <div class='w-[13%] tdBookingCell' id='bookingHourDay1{hour}'></div>
                      <div class='w-[13%] tdBookingCell' id='bookingHourDay2{hour}'></div>
                      <div class='w-[13%] tdBookingCell' id='bookingHourDay3{hour}'></div>
                      <div class='w-[13%] tdBookingCell' id='bookingHourDay4{hour}'></div>
                      <div class='w-[13%] tdBookingCell' id='bookingHourDay5{hour}'></div>
                      <div class='w-[13%] tdBookingCell' id='bookingHourDay6{hour}'></div>
                    </div>
        
                  {/each}
            </div>

      </div>

      <!-- lado direito, lista com veiculos disponiveis -->
      <div class="w-[15%]  mt-[40px]  overflow-y-scroll pl-3 pr-2" id='bookingCarsBrowser'>

            {#if carOptions}

              <div class="flex flex-col overflow-hidden " >

                  <!-- cards de carros disponiveis -->
                  {#each carOptions as {name, id, car_image}} 
        
                    <div class='carCardBooking putWhiteTooltip'    style="background-image: url('{$imagesUrl + car_image}') ;"  aria-hidden="true"  
                              on:click={ () => {bookingSelectedCarId=id; showCarBooking()} }  title={name} id='carCard{id}'  >
                    </div>    
                  {/each}

              </div>    


            {:else}
                <div class="flex flex-col overflow-hidden h-full" >
                      <div class='flex justify-center items-center h-full '>
                        <img src='loading.gif' style='width:54px' alt='...'  >
                      </div>
                </div>
            {/if}
      </div>
  </div>

</div>
