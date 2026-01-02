// 滞留记录实体 ER 图 - 使用 fletcher 包
#import "@preview/fletcher:0.5.8" as fletcher: diagram, edge, node
#import fletcher.shapes: ellipse, rect

#diagram(
  spacing: (15pt, 12pt),
  node-stroke: 1pt,
  edge-stroke: 0.8pt,

  // 中央实体
  node(
    (0, 0),
    text(size: 12pt, weight: "bold")[滞留记录 Overdue],
    stroke: 2pt,
    shape: rect,
    width: 4.2cm,
    height: 1.2cm,
    name: <entity>,
  ),

  // 属性节点 - 主键
  node(
    (0, -1.5),
    text(size: 9pt)[#underline[记录ID]],
    stroke: 1pt,
    shape: ellipse,
    inset: 6pt,
    name: <id>,
  ),

  // 属性节点 - 普通属性
  node(
    (-2, -1),
    text(size: 9pt)[滞留开始时间],
    stroke: 1pt,
    shape: ellipse,
    inset: 6pt,
    name: <start-time>,
  ),
  node(
    (2, -1),
    text(size: 9pt)[提醒次数],
    stroke: 1pt,
    shape: ellipse,
    inset: 6pt,
    name: <remind-count>,
  ),
  node(
    (-2, 0.5),
    text(size: 9pt)[最后提醒时间],
    stroke: 1pt,
    shape: ellipse,
    inset: 6pt,
    name: <last-remind>,
  ),
  node((2, 0.5), text(size: 9pt)[滞留费用], stroke: 1pt, shape: ellipse, inset: 6pt, name: <fee>),
  node((0, 1.5), text(size: 9pt)[处理状态], stroke: 1pt, shape: ellipse, inset: 6pt, name: <status>),

  // 连接线
  edge(<entity>, <id>),
  edge(<entity>, <start-time>),
  edge(<entity>, <remind-count>),
  edge(<entity>, <last-remind>),
  edge(<entity>, <fee>),
  edge(<entity>, <status>),
)
