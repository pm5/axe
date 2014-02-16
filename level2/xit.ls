require! <[fs cheerio request async]>
{lists-to-obj} = require 'prelude-ls'

columns = <[town village name]>

data = []

tasks = for let i from 1 to 12
  (cb) ->
    error, res, body <- request.get "http://axe-level-1.herokuapp.com/lv2/?page=#{i}"
    $ = cheerio.load body
    row = $ 'table.table tr' .first!.next!
    d = []
    while row.text!
      d.push lists-to-obj columns, (row.children!.map (i, d) -> $(d).text!).to-array!
      row .= next!
    cb null, d

error, result <- async.parallel tasks
for d in result
  data ++= d
data |> JSON.stringify |> fs.write-file-sync 'output.json', _
