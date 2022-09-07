# parsePDF
京都大学が公開している「各授業の平均学生在籍数」のPDFをパースしてjson形式で保存するスクリプト

## 使い方
1. 京都大学の各授業の平均学生在籍数のページからPDFをダウンロードする
2. ファイルを`parsePDF/pdf`ディレクトリに移動する（ファイル名は`{年度}.pdf`）
3. `main.go`の`YEAR`をいい感じに設定する
4. `go run main.go`を実行する


## データの出力先
`parsePDF/export/{年度}.json`に出力される

## いろいろ
- `malformed PDF: reading at offset 0: stream not present`が起きる時
```bash
mutool clean -s 2019.pdf 2019_out.pdf
```
を実行すれば直る