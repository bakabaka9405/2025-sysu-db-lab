// 寄件实体 ER 图 - 使用 fletcher 包
#import "@preview/fletcher:0.5.8" as fletcher: diagram, edge, node
#import fletcher.shapes: ellipse, rect

#diagram(
  spacing: (15pt, 12pt),
  node-stroke: 1pt,
  edge-stroke: 0.8pt,

  // 中央实体
  node(
    (0, 0),
    text(size: 12pt, weight: "bold")[寄件 Shipment],
    stroke: 2pt,
    shape: rect,
    width: 3.5cm,
    height: 1.2cm,
    name: <entity>,
  ),

  // 属性节点 - 主键
  node((0, -2), text(size: 9pt)[#underline[寄件ID]], stroke: 1pt, shape: ellipse, inset: 6pt, name: <id>),

  // 属性节点 - 普通属性
  node((-2, -1), text(size: 9pt)[寄件单号], stroke: 1pt, shape: ellipse, inset: 6pt, name: <tracking>),
  node((2, -1), text(size: 9pt)[寄件人信息], stroke: 1pt, shape: ellipse, inset: 6pt, name: <sender>),
  node((-2.5, 0), text(size: 9pt)[收件人信息], stroke: 1pt, shape: ellipse, inset: 6pt, name: <recipient>),
  node((2.5, 0), text(size: 9pt)[快递公司], stroke: 1pt, shape: ellipse, inset: 6pt, name: <company>),
  node((-2, 1), text(size: 9pt)[重量], stroke: 1pt, shape: ellipse, inset: 6pt, name: <weight>),
  node((2, 1), text(size: 9pt)[状态], stroke: 1pt, shape: ellipse, inset: 6pt, name: <status>),
  node((0, 2), text(size: 9pt)[运费], stroke: 1pt, shape: ellipse, inset: 6pt, name: <fee>),

  // 连接线
  edge(<entity>, <id>),
  edge(<entity>, <tracking>),
  edge(<entity>, <sender>),
  edge(<entity>, <recipient>),
  edge(<entity>, <company>),
  edge(<entity>, <weight>),
  edge(<entity>, <status>),
  edge(<entity>, <fee>),
)
