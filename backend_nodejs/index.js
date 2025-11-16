

//const tmp = require('./teste')

//const os = require('os')

const http = require('http')
const fs = require('fs')
const _ = require('lodash')

const server = http.createServer((req, res) => {
  console.log(req.url, req.method)
//  res.setHeader('Content-Type', 'text/plain')
  res.setHeader('Content-Type', 'text/html')
//  res.write("<div style='font-size:40px;color:red'>teste legal</div>")

  fs.readFile('./teste.html', (err, data) => {
    if (err) {
      console.log("erro= ", err) 
    } else {
      //res.write(data)
      res.statusCode=400
      res.end(data);
    }
  })
  
})

server.listen(3000, 'localhost') 

//console.log(tmp.people)
//console.log(os.platform() , os.homedir())