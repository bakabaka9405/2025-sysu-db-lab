// 包裹实体 ER 图 - 使用 fletcher 包
#import "@preview/fletcher:0.5.8" as fletcher: diagram, edge, node
#import fletcher.shapes: ellipse, rect

#diagram(
  spacing: (15pt, 12pt),
  node-stroke: 1pt,
  edge-stroke: 0.8pt,

  // 中央实体
  node(
    (0, 0),
    text(size: 12pt, weight: "bold")[包裹 Parcel],
    stroke: 2pt,
    shape: rect,
    width: 3cm,
    height: 1.2cm,
    name: <entity>,
  ),

  // 属性节点 - 主键
  node((0, -2.5), text(size: 9pt)[#underline[包裹ID]], stroke: 1pt, shape: ellipse, inset: 6pt, name: <id>),

  // 属性节点 - 普通属性（上方）
  node((-2, -1.5), text(size: 9pt)[快递单号], stroke: 1pt, shape: ellipse, inset: 6pt, name: <tracking>),
  node((2, -1.5), text(size: 9pt)[取件码], stroke: 1pt, shape: ellipse, inset: 6pt, name: <pickup-code>),

  // 属性节点 - 左右两侧
  node((-2.8, -0.5), text(size: 9pt)[收件人姓名], stroke: 1pt, shape: ellipse, inset: 6pt, name: <recipient-name>),
  node((2.8, -0.5), text(size: 9pt)[收件人手机], stroke: 1pt, shape: ellipse, inset: 6pt, name: <recipient-phone>),
  node((-2.8, 0.5), text(size: 9pt)[快递公司], stroke: 1pt, shape: ellipse, inset: 6pt, name: <company>),
  node((2.8, 0.5), text(size: 9pt)[包裹大小], stroke: 1pt, shape: ellipse, inset: 6pt, name: <size>),

  // 属性节点 - 下方
  node((-2, 1.5), text(size: 9pt)[重量], stroke: 1pt, shape: ellipse, inset: 6pt, name: <weight>),
  node((2, 1.5), text(size: 9pt)[入库时间], stroke: 1pt, shape: ellipse, inset: 6pt, name: <checkin-time>),
  node((0, 2.5), text(size: 9pt)[状态], stroke: 1pt, shape: ellipse, inset: 6pt, name: <status>),

  // 连接线
  edge(<entity>, <id>),
  edge(<entity>, <tracking>),
  edge(<entity>, <pickup-code>),
  edge(<entity>, <recipient-name>),
  edge(<entity>, <recipient-phone>),
  edge(<entity>, <company>),
  edge(<entity>, <size>),
  edge(<entity>, <weight>),
  edge(<entity>, <checkin-time>),
  edge(<entity>, <status>),
)
