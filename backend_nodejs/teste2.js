
const express = require('express')
const morgan = require('morgan')
const mongoose = require('mongoose')



const app = express()

const usuariosRouter = require('./routes/usuariosRouter')

const dbURI = 'mongodb://localhost:27017/nodejs'
mongoose.connect(dbURI) 
.then((result) => {console.log("\n\nconnected!!\n\n");app.listen(3000)})
.catch((err) => {console.log('erro ', err)})


app.set('view engine', 'ejs')
//app.set('views', 'views_folder')

app.use(express.static('public'))
app.use(morgan('dev'))

app.use(express.urlencoded())
app.use('/usuarios', usuariosRouter)


