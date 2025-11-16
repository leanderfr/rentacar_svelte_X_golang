
const express = require('express')
const morgan = require('morgan')
const mongoose = require('mongoose')

const Usuario = require('./models/usuarios')

const app = express()

const dbURI = 'mongodb://localhost:27017/nodejs'
mongoose.connect(dbURI) 
.then((result) => {console.log("\n\nconnected!!\n\n");app.listen(3000)})
.catch((err) => {console.log('erro ', err)})


app.set('view engine', 'ejs')
//app.set('views', 'views_folder')

app.use(express.static('public'))
app.use(morgan('dev'))

app.use(express.urlencoded())


app.use((req, res, next) =>  {
//  console.log('hostname= ', req.hostname)
//  console.log('path= ', req.path)
//  console.log('method= ', req.method)
  next()
})

//********************************************************************************************************************** */
//********************************************************************************************************************** */
app.get('/add-usuario', (req, res) =>  {
  const usuario = new Usuario( {
    nome: 'Sera que vai funcionar',
    idade: 444,
  })
  usuario.save()
    .then((result) => {res.send(result)})
    .catch((err) => {res.send(err)})
})


//********************************************************************************************************************** */
//********************************************************************************************************************** */

app.get('/add-usuario', (req, res) =>  {
  const usuario = new Usuario( {
    nome: 'Sera que vai funcionar',
    idade: 444,
  })
  usuario.save()
    .then((result) => {res.send(result)})
    .catch((err) => {res.send(err)})
})

//********************************************************************************************************************** */
//********************************************************************************************************************** */


app.get('/adiciona', (req, res) =>  {
  const usuario = new Usuario( req.body )
  usuario.save()
    .then((result) => {res.send(result)})
    .catch((err) => {res.send(err)})
})


//********************************************************************************************************************** */
//********************************************************************************************************************** */

app.get('/usuarios', (req, res) =>  {
  Usuario.find().sort( { createdAt: -1 })
  .then((result) => {res.send(result)})
  .catch((err) => {res.send(err)})
})

//********************************************************************************************************************** */
//********************************************************************************************************************** */

app.get('/usuario', (req, res) =>  {
  Usuario.findById( '691909e64c471b25e6fb4260' )
  .then((result) => {res.send(result)})
  .catch((err) => {res.send(err)})

})



//********************************************************************************************************************** */
//********************************************************************************************************************** */

app.get('/', (req, res) =>  {
//  res.send("<div style='color:navy;font-size:60px'>Fabian e Tamara se amam</div>")
  res.render('index', { vaiocorrer: "Fabian casa com Tamara"})
})

//********************************************************************************************************************** */
//********************************************************************************************************************** */

app.get('/', (req, res) =>  {
//  res.send("<div style='color:navy;font-size:60px'>Fabian e Tamara se amam</div>")
  res.sendFile('./teste.html', { root: __dirname })
})

//********************************************************************************************************************** */
//********************************************************************************************************************** */

app.use('/', (req, res) =>  {
  res.status(404).send("rota nao existe")
})
