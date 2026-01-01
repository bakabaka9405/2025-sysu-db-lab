// 流程图公共组件库 - 使用 fletcher 包
#import "@preview/fletcher:0.5.8" as fletcher: diagram, edge, node
#import fletcher.shapes: diamond, ellipse, rect

// 终端节点（开始/结束）- 圆角矩形
#let terminal(pos, label) = node(
  pos,
  text(size: 9pt)[#label],
  stroke: 1pt,
  shape: rect,
  corner-radius: 12pt,
  fill: rgb("#e8f5e9"),
  inset: 8pt,
)

// 处理节点 - 矩形
#let process(pos, label) = node(
  pos,
  text(size: 9pt)[#label],
  stroke: 1pt,
  shape: rect,
  fill: rgb("#e3f2fd"),
  inset: 8pt,
)

// 判断节点 - 菱形
#let decision(pos, label) = node(
  pos,
  text(size: 8pt)[#label],
  stroke: 1pt,
  shape: diamond,
  fill: rgb("#fff3e0"),
  inset: 8pt,
)

// 箭头边
#let arrow(from, to, label: none, label-side: left) = {
  if label != none {
    edge(from, to, "-|>", stroke: 0.8pt, label: text(size: 8pt)[#label], label-side: label-side)
  } else {
    edge(from, to, "-|>", stroke: 0.8pt)
  }
}

// 普通连接线（无箭头）
#let line-edge(from, to) = edge(from, to, stroke: 0.8pt)
