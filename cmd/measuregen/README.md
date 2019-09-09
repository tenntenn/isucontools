# measuregen

標準入力にGoのソースコードを与えると
[measure](https://github.com/najeira/measure)の計測用の関数呼び出しを追加します。

```sh
$ cat main.go | measuregen
```

なお、同じソースに何度書けても同じ関数には2度は同じ関数呼び出しを追加しないようになっています。
そのため、新しく追加した関数に対して計測用の関数呼び出しを差し込みたい場合は単にもう一度`measuregen`をかければ良いです。

```sh
# 結果は同じ
$ cat main.go | measuregen | measuregen
```
