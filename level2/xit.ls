require! <[fs cheerio request async]>

cols = <[town village name]>

tasks = [1 to 12].map (page) ->
  (cb) ->
    error, res, body <- request.get "http://axe-level-1.herokuapp.com/lv2/?page=#{page}"
    $ = cheerio.load body
    row = $ 'table.table tr' .first!.next!
    data = []
    while row.text!
      data.push {[cols[i], d] for d, i in (row.children!.map (, it) -> $(it).text!).to-array!}
      row .= next!
    cb null, data

error, result <- async.series tasks
data = result.reduce (prev, cur) -> prev ++ cur
data |> JSON.stringify |> fs.write-file-sync 'output.json', _
