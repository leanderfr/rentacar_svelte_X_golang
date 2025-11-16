
const mongoose = require('mongoose')


const Schema = mongoose.Schema 

const UsuarioSchema = new Schema({
  nome: {
    type: String,
    required: true
  },
  idade: {
    type: Number,
    min: 0, // Optional: define minimum value
    max: 120, // Optional: define maximum value
    required: true
  },

}, { timestamps: true }) 


const Usuario = mongoose.model('Usuarios', UsuarioSchema)

module.exports = Usuario

