// 上下文数据流图 - 使用 fletcher 包
#import "@preview/fletcher:0.5.8" as fletcher: diagram, edge, node
#import fletcher.shapes: diamond, ellipse, rect
#diagram(
  spacing: (30pt, 25pt),
  node-stroke: 1pt,
  edge-stroke: 0.8pt,

  // 外部实体 - 四角布局
  node(
    (0, 0),
    text(size: 9pt, weight: "bold")[普通用户],
    stroke: 1.2pt,
    shape: rect,
    fill: rgb("#fff9c4"),
    width: 2.2cm,
    height: 0.8cm,
    name: <user>,
  ),
  node(
    (4, 0),
    text(size: 9pt, weight: "bold")[驿站员工],
    stroke: 1.2pt,
    shape: rect,
    fill: rgb("#fff9c4"),
    width: 2.2cm,
    height: 0.8cm,
    name: <staff>,
  ),
  node(
    (0, 4),
    text(size: 9pt, weight: "bold")[快递员],
    stroke: 1.2pt,
    shape: rect,
    fill: rgb("#fff9c4"),
    width: 2.2cm,
    height: 0.8cm,
    name: <courier>,
  ),
  node(
    (4, 4),
    text(size: 9pt, weight: "bold")[管理员],
    stroke: 1.2pt,
    shape: rect,
    fill: rgb("#fff9c4"),
    width: 2.2cm,
    height: 0.8cm,
    name: <admin>,
  ),

  // 中央处理系统
  node(
    (2, 2),
    text(size: 9pt, weight: "bold")[驿站管理系统],
    stroke: 1.2pt,
    shape: circle,
    fill: rgb("#e3f2fd"),
    width: 2.8cm,
    height: 2.8cm,
    name: <system>,
  ),

  // 数据流 - 用户
  edge(<user>, <system>, "-|>", label: text(size: 7pt)[取件请求/寄件信息], label-side: left, bend: 30deg),
  edge(<system>, <user>, "-|>", label: text(size: 7pt)[包裹状态/通知], label-side: right, bend: 30deg),

  // 数据流 - 员工
  edge(<staff>, <system>, "-|>", label: text(size: 7pt)[入库操作/状态更新], label-side: right, bend: 30deg),
  edge(<system>, <staff>, "-|>", label: text(size: 7pt)[报表/日志], label-side: left, bend: 30deg),

  // 数据流 - 快递员
  edge(<courier>, <system>, "-|>", label: text(size: 7pt)[包裹投递], label-side: left, bend: 30deg),
  edge(<system>, <courier>, "-|>", label: text(size: 7pt)[取件码/通知], label-side: right, bend: 30deg),

  // 数据流 - 管理员
  edge(<admin>, <system>, "-|>", label: text(size: 7pt)[系统配置], label-side: right, bend: 30deg),
  edge(<system>, <admin>, "-|>", label: text(size: 7pt)[包裹信息/统计], label-side: left, bend: 30deg),
)
