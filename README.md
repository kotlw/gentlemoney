# Gentlemoney

[![License](https://img.shields.io/github/license/kotlw/gentlemoney)](https://github.com/kotlw/gentlemoney/blob/main/LICENSE)
[![codecov](https://codecov.io/gh/kotlw/gentlemoney/branch/main/graph/badge.svg?token=1TMPI2NDBQ)](https://codecov.io/gh/kotlw/gentlemoney)
[![Go Report Card](https://goreportcard.com/badge/github.com/kotlw/gentlemoney)](https://goreportcard.com/report/github.com/kotlw/gentlemoney)

Terminal money manager for personal use.
<p align="center">
  <img src="./assets/demonstration.gif" />
</p>

## Instalation
There are no trix to run this app. Just clone it and run.
```
git clone https://github.com/kotlw/gentlemoney.git
cd gentlemoney
go run ./cmd/gmon
```

## Navigation
Since navigation hints in app is missing here are some description of how to use it. [Tview](https://github.com/rivo/tview) has a bit of predefined bindings which are good such as table navigation using vim bindings ```'h'```, ```'j'```, ```'k'```, ```'l'```. As for others they are more or less intuitive. Here are the list:
 - ```Enter``` - submit
 - ```Esc``` - cancel
 - ```Tab``` - focus next item
 - ```Shift+Tab``` - focus previous item
 - ```c``` - create (transaction/account/currency/category)
 - ```u``` - update
 - ```d``` - delete
