#import "@preview/cetz:0.4.2"

#cetz.canvas({
  import cetz.draw: *

  // 包裹实体矩形
  rect((0, 0), (4, 2), stroke: 2pt)
  content((2, 1), text(size: 14pt, weight: "bold")[包裹 Parcel])

  // 属性椭圆
  let attrs = (
    ("包裹ID", (2, 3.5)),
    ("快递单号", (-2, 2.5)),
    ("取件码", (6, 2.5)),
    ("收件人姓名", (-2.5, 1)),
    ("收件人手机", (6.5, 1)),
    ("快递公司", (-2.5, -0.5)),
    ("包裹大小", (6.5, -0.5)),
    ("重量", (-2, -2)),
    ("状态", (2, -2)),
    ("入库时间", (6, -2)),
  )

  for (name, pos) in attrs {
    circle(pos, radius: 0.8, stroke: 1pt)
    content(pos, text(size: 9pt)[#name])
    line((2, 1), pos, stroke: 0.5pt)
  }

  // 主键标记（下划线）
  line((1.5, 3.5), (2.5, 3.5), stroke: 1pt)
})
