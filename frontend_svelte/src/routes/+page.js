
import { selectedLanguage } from '$js/utils.js'
import { backendUrlGolang  } from '$lib/stores/stores.js'


//*****************************************************************************************************************
// carrega as expressoes do idioma atual, para preenchimento inicial de todos os textos na aplicacao
// inicia usando os termos default, do workgroup 'Admin'
// à medida que usuario iniciar seu proprio workgroup, buscara os termos do mesmo
//*****************************************************************************************************************
export const load = async ({fetch}) => {

// le conteudo atual da variavel writable (backendUrl) ==> url do back end
// o carregamento inicial dos termos sera feito pelo backend golang
// nesse momento nao é possivel saber qual backend o usuario escolheu 
// localstorage nao funciona aqui em '+page.js'
let _backendUrlGolang_
backendUrlGolang.subscribe((url) => {
  _backendUrlGolang_ = url
})

let language = selectedLanguage()

const languageRes = await fetch(`${_backendUrlGolang_}/terms/${language}`, { method: 'GET' })
const languageData = await languageRes.json();

return {
  languageData
}

}
