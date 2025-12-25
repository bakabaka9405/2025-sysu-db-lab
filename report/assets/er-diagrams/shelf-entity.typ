#import "@preview/cetz:0.4.2"

#cetz.canvas({
  import cetz.draw: *

  // 货架实体矩形
  rect((0, 0), (4, 2), stroke: 2pt)
  content((2, 1), text(size: 14pt, weight: "bold")[货架 Shelf])

  // 属性椭圆
  let attrs = (
    ("货架ID", (2, 3)),
    ("货架编号", (-1.5, 2)),
    ("区域", (5.5, 2)),
    ("楼层", (-1.5, 0)),
    ("列号", (5.5, 0)),
    ("类型", (-1.5, -1.5)),
    ("容量", (2, -1.5)),
    ("当前占用数", (5.5, -1.5)),
  )

  for (name, pos) in attrs {
    circle(pos, radius: 0.8, stroke: 1pt)
    content(pos, text(size: 9pt)[#name])
    line((2, 1), pos, stroke: 0.5pt)
  }

  // 主键标记（下划线）
  line((1.5, 3), (2.5, 3), stroke: 1pt)
})
