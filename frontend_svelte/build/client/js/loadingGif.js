

/*******************************************************************************/
/*******************************************************************************/
// exibe a animacao ajax
function showLoadingGif() {

// devido a muitos eventos asincronos que requisitam exibicao da animacao 'processando...'
// verifica se ela ja nao esta visivel , se sim, cai fora
if (typeof jq('#blockAppInteraction').attr('id')!='undefined') return

// .spinner é a animacao em si, localizada dentro da div 'divLoading'
jq('.spinner').show()

// centra a animacao pq via CSS nao funciona de jeito nenhum, e a posicao central muda dinamicamente 
$divAjax = jq('#divLoading');  

$divAjax.css("position","absolute");
$divAjax.css("top", Math.max(0, ((jq(window).height() - $divAjax.outerHeight()) / 2) + 
                                            jq(window).scrollTop()) + "px");
$divAjax.css("left", Math.max(0, ((jq(window).width() - $divAjax.outerWidth()) / 2) + 
                                            jq(window).scrollLeft()) + "px");

jq("#divLoading").removeClass("cssLOADING_SHOW").removeClass("cssLOADING_HIDE");
jq("#divLoading").addClass("cssLOADING_SHOW");

jq('body').append('<div id="blockAppInteraction" style="position: absolute;top:0;left:0;width: 100%;height:100%;z-index:2;background:rgba(0,0,0,0);"></div>');

}


/*******************************************************************************/
/*******************************************************************************/
// oculta a animacao ajax
function hideLoadingGif() {

// .spinner é a animacao em si, localizada dentro da div 'divLoading'
// se ela nao é escondida antes, ocorre um 'buraco branco' ao esconder 'divLoading' 
jq('.spinner').hide()

setTimeout(() => {
jq('#divLoading').removeClass("cssLOADING_SHOW").removeClass("cssLOADING_HIDE");
jq('#divLoading').addClass("cssLOADING_HIDE");

  jq('#blockAppInteraction').remove();  
}, 50);
}
 
