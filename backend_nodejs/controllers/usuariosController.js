
const Usuario = require('../models/usuarios')


const usuarios_index = (req, res) => {

console.log('tudo certo')

  Usuario.find().sort( { createdAt: -1 })
  .then((result) => {res.send(result)})
  .catch((err) => {res.send(err)})

}

module.exports = {
  usuarios_index
}