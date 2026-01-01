// 数据流图(DFD)公共组件库 - 使用 fletcher 包
#import "@preview/fletcher:0.5.8" as fletcher: diagram, edge, node

// 外部实体 - 矩形
#let entity(pos, label, size: "normal") = {
  let (w, h, fs) = if size == "large" {
    (2.5cm, 0.8cm, 9pt)
  } else {
    (1.8cm, 0.6cm, 8pt)
  }
  node(
    pos,
    text(size: fs, weight: "bold")[#label],
    stroke: 1.2pt,
    shape: rect,
    fill: rgb("#fff9c4"),
    width: w,
    height: h,
  )
}

// 处理过程 - 圆形（简单版，只有标签）
#let process-simple(pos, label, radius: 1.2cm) = node(
  pos,
  text(size: 9pt, weight: "bold")[#label],
  stroke: 1.2pt,
  shape: circle,
  fill: rgb("#e3f2fd"),
  width: radius * 2,
  height: radius * 2,
)

// 处理过程 - 圆形（带编号）
#let process-numbered(pos, id, label) = node(
  pos,
  [#text(size: 7pt, weight: "bold")[#id]\ #text(size: 7pt)[#label]],
  stroke: 1pt,
  shape: circle,
  fill: rgb("#e3f2fd"),
  width: 1.8cm,
  height: 1.8cm,
)

// 数据存储 - 矩形
#let store(pos, id, label) = node(
  pos,
  [#text(size: 7pt)[#id] │ #text(size: 7pt)[#label]],
  stroke: 1pt,
  shape: rect,
  fill: rgb("#f5f5f5"),
  inset: 6pt,
)

// 数据流 - 带标签的箭头
#let dataflow(from, to, label, label-side: left) = edge(
  from,
  to,
  "-|>",
  stroke: 0.8pt,
  label: text(size: 7pt)[#label],
  label-side: label-side,
)
