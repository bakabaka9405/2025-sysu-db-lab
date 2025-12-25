#import "@preview/cetz:0.4.2"

#cetz.canvas(length: 1cm, {
  import cetz.draw: *

  // 用户实体
  rect((0, 8), (3, 10), stroke: 2pt)
  content((1.5, 9), text(weight: "bold")[用户\nUser])

  // 包裹实体
  rect((8, 8), (11, 10), stroke: 2pt)
  content((9.5, 9), text(weight: "bold")[包裹\nParcel])

  // 货架实体
  rect((14, 8), (17, 10), stroke: 2pt)
  content((15.5, 9), text(weight: "bold")[货架\nShelf])

  // 寄件实体
  rect((0, 2), (3, 4), stroke: 2pt)
  content((1.5, 3), text(weight: "bold")[寄件\nShipment])

  // 滞留记录实体
  rect((8, 2), (11, 4), stroke: 2pt)
  content((9.5, 3), text(weight: "bold")[滞留记录\nOverdue])

  // 操作日志实体
  rect((14, 2), (17, 4), stroke: 2pt)
  content((15.5, 3), text(weight: "bold")[操作日志\nLog])

  // 用户-包裹关系（收件）
  line((3, 9), (8, 9), stroke: 1.5pt)
  content((5.5, 9.5), [收件])
  content((3.5, 9.3), text(size: 10pt)[1])
  content((7.5, 9.3), text(size: 10pt)[N])

  // 货架-包裹关系（存放）
  line((11, 9), (14, 9), stroke: 1.5pt)
  content((12.5, 9.5), [存放])
  content((11.5, 9.3), text(size: 10pt)[N])
  content((13.5, 9.3), text(size: 10pt)[1])

  // 用户-寄件关系
  line((1.5, 8), (1.5, 4), stroke: 1.5pt)
  content((0.5, 6), [寄件])
  content((1.8, 7.5), text(size: 10pt)[1])
  content((1.8, 4.5), text(size: 10pt)[N])

  // 包裹-滞留记录关系
  line((9.5, 8), (9.5, 4), stroke: 1.5pt)
  content((10.5, 6), [滞留])
  content((9.8, 7.5), text(size: 10pt)[1])
  content((9.8, 4.5), text(size: 10pt)[1])

  // 用户-操作日志关系（折线）
  line((1.5, 8), (1.5, 6), stroke: 1.5pt)
  line((1.5, 6), (15.5, 6), stroke: 1.5pt)
  line((15.5, 6), (15.5, 4), stroke: 1.5pt)
  content((8.5, 6.5), [操作])
  content((2, 7), text(size: 10pt)[1])
  content((15.8, 4.5), text(size: 10pt)[N])
})
