// 总体关系 ER 图 - 使用 fletcher 包
#import "@preview/fletcher:0.5.8" as fletcher: diagram, edge, node
#import fletcher.shapes: diamond, rect

#diagram(
  spacing: (25pt, 20pt),
  node-stroke: 1pt,
  edge-stroke: 1.2pt,

  // ===== 实体定义 =====
  // 用户实体（左上）
  node(
    (0, 0),
    text(size: 11pt, weight: "bold")[用户\nUser],
    stroke: 2pt,
    shape: rect,
    width: 2.2cm,
    height: 1.2cm,
    name: <user>,
  ),

  // 包裹实体（中上）
  node(
    (4, 0),
    text(size: 11pt, weight: "bold")[包裹\nParcel],
    stroke: 2pt,
    shape: rect,
    width: 2.2cm,
    height: 1.2cm,
    name: <parcel>,
  ),

  // 货架实体（右上）
  node(
    (8, 0),
    text(size: 11pt, weight: "bold")[货架\nShelf],
    stroke: 2pt,
    shape: rect,
    width: 2.2cm,
    height: 1.2cm,
    name: <shelf>,
  ),

  // 寄件实体（左下）
  node(
    (0, 4),
    text(size: 11pt, weight: "bold")[寄件\nShipment],
    stroke: 2pt,
    shape: rect,
    width: 2.2cm,
    height: 1.2cm,
    name: <shipment>,
  ),

  // 滞留记录实体（中下）
  node(
    (4, 4),
    text(size: 11pt, weight: "bold")[滞留记录\nOverdue],
    stroke: 2pt,
    shape: rect,
    width: 2.5cm,
    height: 1.2cm,
    name: <overdue>,
  ),

  // 操作日志实体（右下）
  node(
    (8, 4),
    text(size: 11pt, weight: "bold")[操作日志\nLog],
    stroke: 2pt,
    shape: rect,
    width: 2.2cm,
    height: 1.2cm,
    name: <log>,
  ),

  // ===== 关系（菱形）定义 =====

  // 关系1：用户-包裹（收件）
  node((2, 0), text(size: 9pt)[收件], stroke: 1.5pt, shape: diamond, inset: 8pt, name: <rel-receive>),

  // 关系2：货架-包裹（存放）
  node((6, 0), text(size: 9pt)[存放], stroke: 1.5pt, shape: diamond, inset: 8pt, name: <rel-store>),

  // 关系3：用户-寄件（发起）
  node((0, 2), text(size: 9pt)[发起], stroke: 1.5pt, shape: diamond, inset: 8pt, name: <rel-send>),

  // 关系4：包裹-滞留记录（产生）
  node((4, 2), text(size: 9pt)[产生], stroke: 1.5pt, shape: diamond, inset: 8pt, name: <rel-overdue>),

  // 关系5：用户-操作日志（执行）
  node((4, 5.5), text(size: 9pt)[执行], stroke: 1.5pt, shape: diamond, inset: 8pt, name: <rel-log>),

  // ===== 连接线与基数标注 =====

  // 用户 - 收件 - 包裹
  edge(<user>, <rel-receive>, label: text(size: 9pt)[1], label-pos: 0.3, label-side: left),
  edge(<rel-receive>, <parcel>, label: text(size: 9pt)[N], label-pos: 0.7, label-side: left),

  // 包裹 - 存放 - 货架
  edge(<parcel>, <rel-store>, label: text(size: 9pt)[N], label-pos: 0.3, label-side: left),
  edge(<rel-store>, <shelf>, label: text(size: 9pt)[1], label-pos: 0.7, label-side: left),

  // 用户 - 发起 - 寄件
  edge(<user>, <rel-send>, label: text(size: 9pt)[1], label-pos: 0.3, label-side: left),
  edge(<rel-send>, <shipment>, label: text(size: 9pt)[N], label-pos: 0.7, label-side: left),

  // 包裹 - 产生 - 滞留记录
  edge(<parcel>, <rel-overdue>, label: text(size: 9pt)[1], label-pos: 0.3, label-side: left),
  edge(<rel-overdue>, <overdue>, label: text(size: 9pt)[1], label-pos: 0.7, label-side: left),

  // 用户 - 执行 - 操作日志（折线连接）
  edge(<user>, (0, 5.5), stroke: 1.2pt),
  edge((0, 5.5), <rel-log>, label: text(size: 9pt)[1], label-pos: 0.5, label-side: left),
  edge(<rel-log>, (8, 5.5), label: text(size: 9pt)[N], label-pos: 0.5, label-side: left),
  edge((8, 5.5), <log>, stroke: 1.2pt),
)
