module Main where

import Network.HTTP
import Text.HTML.TagSoup
import Text.JSON
import Data.List
import Control.Monad (liftM)
import qualified Codec.Binary.UTF8.String as U

main :: IO()
main = do
  tags <- liftM parseTags $ openURL "http://axe-level-1.herokuapp.com/"
  let rows = map parseRow $ sections (~== "<tr>") tags
  writeFile "output.json" $
    encode $
    compileJSON $ tail rows

  where
    openURL :: String -> IO String
    openURL x = getResponseBody =<< simpleHTTP (getRequest x)

    parseRow :: [Tag String] -> [String]
    parseRow = map parseColumn . sections (~== "<td>") . takeWhile (~/= "</tr>")

    parseColumn :: [Tag String] -> String
    parseColumn = fromTagText . head . filter isTagText

    compileText :: [[String]] -> String
    compileText = unlines . map unwords

    compileJSON :: [[String]] -> [JSObject JSValue]
    compileJSON = map compileRow

    compileRow :: [String] -> JSObject JSValue
    compileRow row = toJSObject $
      [ ("name", JSString $ toJSString $ U.decodeString $ head row)
      , ("grades", JSObject $ compileGrades $ tail row) ]

    compileGrades :: [String] -> JSObject JSValue
    compileGrades row = toJSObject $ zip
      [ "國語"
      , "數學"
      , "自然"
      , "社會"
      , "健康教育" ]
      (map (JSRational True) (map toRational (map read row :: [Int])))
