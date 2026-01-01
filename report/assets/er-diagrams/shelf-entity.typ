// 货架实体 ER 图 - 使用 fletcher 包
#import "@preview/fletcher:0.5.8" as fletcher: diagram, edge, node
#import fletcher.shapes: ellipse, rect

#diagram(
  spacing: (15pt, 12pt),
  node-stroke: 1pt,
  edge-stroke: 0.8pt,

  // 中央实体
  node(
    (0, 0),
    text(size: 12pt, weight: "bold")[货架 Shelf],
    stroke: 2pt,
    shape: rect,
    width: 3cm,
    height: 1.2cm,
    name: <entity>,
  ),

  // 属性节点 - 主键
  node((0, -2), text(size: 9pt)[#underline[货架ID]], stroke: 1pt, shape: ellipse, inset: 6pt, name: <id>),

  // 属性节点 - 普通属性
  node((-2, -1), text(size: 9pt)[货架编号], stroke: 1pt, shape: ellipse, inset: 6pt, name: <code>),
  node((2, -1), text(size: 9pt)[区域], stroke: 1pt, shape: ellipse, inset: 6pt, name: <area>),
  node((-2.5, 0), text(size: 9pt)[楼层], stroke: 1pt, shape: ellipse, inset: 6pt, name: <floor>),
  node((2.5, 0), text(size: 9pt)[列号], stroke: 1pt, shape: ellipse, inset: 6pt, name: <column>),
  node((-2, 1), text(size: 9pt)[类型], stroke: 1pt, shape: ellipse, inset: 6pt, name: <type>),
  node((2, 1), text(size: 9pt)[当前占用数], stroke: 1pt, shape: ellipse, inset: 6pt, name: <current>),
  node((0, 2), text(size: 9pt)[容量], stroke: 1pt, shape: ellipse, inset: 6pt, name: <capacity>),

  // 连接线
  edge(<entity>, <id>),
  edge(<entity>, <code>),
  edge(<entity>, <area>),
  edge(<entity>, <floor>),
  edge(<entity>, <column>),
  edge(<entity>, <type>),
  edge(<entity>, <current>),
  edge(<entity>, <capacity>),
)
