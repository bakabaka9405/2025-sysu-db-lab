#set page(
  paper: "a4",
  margin: (x: 3cm, y: 3cm),
  numbering: "1",
)

#set text(
  font: "Source Han Serif SC",
  size: 12pt,
  lang: "zh",
)

#show math.equation: set text(font: ("New Computer Modern Math", "New Computer Modern", "Source Han Serif SC"))
#show raw: set text(font: (
  (name: "DejaVu Sans Mono", covers: "latin-in-cjk"),
  "Source Han Serif SC",
))

#import "@preview/codly:1.3.0": *
#import "@preview/codly-languages:0.1.1": *
#show: codly-init.with()

// 标题样式设置
#set heading(numbering: "1.1")
#show heading.where(level: 1): it => {
  pagebreak(weak: true)
  set text(size: 18pt, weight: "bold")
  v(1em)
  it
  v(0.8em)
}
#show heading.where(level: 2): it => {
  set text(size: 15pt, weight: "bold")
  v(0.7em)
  it
  v(0.5em)
}
#show heading.where(level: 3): it => {
  set text(size: 13pt, weight: "bold")
  v(0.5em)
  it
  v(0.3em)
}

// 段落样式
#set par(
  justify: true,
  leading: 1em,
)

#show figure: set block(breakable: true)

// 枚举样式
#set enum(numbering: "1.a.i")

// ===== 封面 =====
#set page(numbering: none)
#align(center)[
  #image("assets/sysu-logo.svg", width: 30%)

  #v(1em)

  #text(size: 34pt, weight: "bold")[数据库系统实验\ 课程设计报告]

  #v(7em)

  #block(width: 85%)[
    #set text(size: 16pt)
    #let field(name, content) = {
      grid(
        columns: (50pt, 1fr),
        gutter: 0pt,
        align(left + horizon)[#name:],
        box(width: 100%, stroke: (bottom: 0.5pt), outset: (bottom: 5pt), align(center, content)),
      )
    }

    #field("题目", "校园“最后一公里”物流驿站管理系统")
    #v(25pt)
    #field("姓名", "")
    #v(25pt)
    #field("学号", "")
    #v(25pt)
    #field("日期", "2026 年 1 月 1 日")
  ]
]

// ===== 摘要 =====
#pagebreak()
#set page(numbering: "I")
#counter(page).update(1)
#include "chapters/00-abstract.typ"

// ===== 目录 =====
#pagebreak()
#outline(
  title: [目录],
  indent: auto,
  depth: 3,
)

// ===== 正文 =====
#pagebreak()
#set page(numbering: "1")
#counter(page).update(1)

#include "chapters/01-introduction.typ"
#include "chapters/02-requirements.typ"
#include "chapters/03-database-design.typ"
#include "chapters/04-implementation.typ"
#include "chapters/05-testing.typ"
#include "chapters/06-conclusion.typ"

// ===== 参考文献 =====
#pagebreak()
#bibliography("./references.bib", title: "参考文献", style: "gb-7714-2015-numeric", full: true)