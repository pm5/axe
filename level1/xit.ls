require! <[cheerio request fs]>
{lists-to-obj} = require 'prelude-ls'

data = []

error, response, body <- request.get 'http://axe-level-1.herokuapp.com/'
$ = cheerio.load body

subjects = $ 'table.table tr' .first!.children!.map (i, e) -> $(e).text!
subjects .= to-array!
subjects.shift!

row = $ 'table.table tr' .first!.next!
while row.text!
  columns = row.children!.map (i, e) -> $(e).text!
  columns .= to-array!
  name = columns.shift!
  data.push {name, grades: lists-to-obj subjects, columns.map -> +it}
  row = row.next!

data |> JSON.stringify |> fs.write-file-sync 'output.json', _
