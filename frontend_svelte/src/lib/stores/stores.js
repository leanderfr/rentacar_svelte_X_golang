import { writable } from 'svelte/store';


// qual idioma/país esta sendo usado no front end
export const Terms = writable([]);
export let UnitedStates_selected = writable(false);

// HomeCalendar_CurrentFirstDate  vai iniciar como sendo o 1o dia do mês vigente
// auxilia o calendario exibido em Home.svelte
let today = new Date();  // hoje
let firstDayMonth = new Date(today.getFullYear(), today.getMonth(), 1);
export let HomeCalendar_CurrentFirstDate = writable(firstDayMonth);

// auxilia a tela de reservas de carros no componente Bookings.svelte
let _today_ = new Date(today.getFullYear(), today.getMonth(), today.getDate());
export let BookingCalendar_CurrentDate = writable(_today_);


// inicia a aplicacao com menu lateral 'principal' clicado
export const clickedSidebarMenu = writable('main')

// backendUrl= default, mas sera substituito pelo backend escolhido pelo cliente (golang, php ou c#)
export const backendUrl = writable('')

// backends possiveis
//export const backendUrlPHP = writable('http://127.0.0.1:80')

export const backendUrlPHP = writable('http://ec2-54-233-183-5.sa-east-1.compute.amazonaws.com:8072')
//export const backendUrlPHP = writable('http://127.0.0.1:8072')

// URL usada durante o desenvolvimento da app
//export let backendUrlGolang = writable('http://127.0.0.1:8070')  
export let backendUrlGolang = writable('http://ec2-54-233-183-5.sa-east-1.compute.amazonaws.com:8070')  

// repositorio de imagens de carros, fabricantes, etc
export const imagesUrl = writable('https://devs-app.s3.sa-east-1.amazonaws.com/rentacar_images/')

export const staticImagesUrl = writable('https://devs-app.s3.sa-east-1.amazonaws.com/rentacar_images/static')

// visualizador de cod fonte
export const sourceCodeViewer = writable('http://leanderdeveloper.store/rentacar_source_code/')

// controla se a aplicacao emite sons de alerta
export let soundEnabled = writable(true);

// controla se o ID dos registros sera exibido
export let recordIdsVisible = writable(false)

// controla se as chamadas API serao visualizadas (parte inferior da tela)
export let apiCallsVisible = writable(false)

// contem array com opcoes de autocompletar para input type= fabricante
export let manufacturersAutocomplete = writable([]);

export let clientIp = writable('');

// qdo um elemento input type=text, é usado para registrar uma data, um calendario auxiliar <Calendar.svelte />, pode ser chamado para preencher a data
// ao chamar Calendar.svelte, 'background-color' do input type=text é alterado para ficar parecido com o CSS do calendario, e dar a ideia de 'uniao'
// isso estraga nao so a bg-color, mas todo CSS do elemento, por isso as 2 variaveis abaixo memorizam qual o elemento alterado e qual era o CSS anterior,
// para que, ao fechar o calendario auxiliar, o CSS volte ao que era
export let elementCSSChangedTemporarily = writable([]);
export let elementIdCSSChangedTemporarily = writable('');

// avisa se imagem de algum formulario sendo carregada ainda
export let imagesStillLoading = writable(false);

// se reportedError for preenchida a qq momento, a div que reporta erro (dentro de +page.svelte) sera exibida com a mensagem do erro contida em 'reportedError'
export let reportedErrorMessage = writable('');

// verifica se ja ha um grupo definido para o cliente (navegador) atual
export let workgroupReady  = writable(true);

// verifica se o idioma ja esta definido
export let isLanguageChosen  = writable(true);

// controla qual é o carro atualmente selecionado em Home.svelte
export let currentHomeSelectedCarId = writable(0); 

// controla qual backend esta fazendo a clonagem dos dados no momento
export let cloningCurrentBackend = writable('')

