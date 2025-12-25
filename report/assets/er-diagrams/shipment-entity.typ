#import "@preview/cetz:0.4.2"

#cetz.canvas({
  import cetz.draw: *

  // 寄件实体矩形
  rect((0, 0), (4, 2), stroke: 2pt)
  content((2, 1), text(size: 14pt, weight: "bold")[寄件 Shipment])

  // 属性椭圆
  let attrs = (
    ("寄件ID", (2, 3)),
    ("寄件单号", (-2, 2)),
    ("寄件人信息", (6, 2)),
    ("收件人信息", (-2, 0)),
    ("快递公司", (6, 0)),
    ("重量", (-2, -1.5)),
    ("运费", (2, -1.5)),
    ("状态", (6, -1.5)),
  )

  for (name, pos) in attrs {
    circle(pos, radius: 0.8, stroke: 1pt)
    content(pos, text(size: 9pt)[#name])
    line((2, 1), pos, stroke: 0.5pt)
  }

  // 主键标记（下划线）
  line((1.5, 3), (2.5, 3), stroke: 1pt)
})
