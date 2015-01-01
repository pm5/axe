'use strict'
require! <[cheerio superagent async fs]>

level-url = "http://axe-level-1.herokuapp.com/lv3/"

getTableData = (bodyText) ->
  $ = cheerio.load bodyText
  rows = $ 'tr' .to-array!
  rows.shift!
  rows.map ->
    [ it.children[i].children[0].data for i in [1, 3, 5] ]

res <- superagent.get level-url .end
throw res.text unless res.ok

var session-id
for s in res.header["set-cookie"]
  [, session-id] = (s == /PHPSESSID=(.*?);/)
  break if session-id
throw "No PHPSESSID found." unless session-id

first-page-data = getTableData res.text

tasks = [
  (done) ->
    res <- superagent.get "#{level-url}?page=next"
      .set 'cookie', "PHPSESSID=#{session-id}"
      .end (res) ->
        | res.ok    => done null, getTableData res.text
        | otherwise => done res.text, null
] * 75

err, data <- async.series tasks
throw err if err
output = first-page-data ++ data.reduce (prev, cur) -> prev ++ cur

[{town: x[0], village: x[1], name: x[2]} for x in output] |> JSON.stringify |> fs.writeFileSync 'output.json', _
