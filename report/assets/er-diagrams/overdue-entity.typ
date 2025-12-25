#import "@preview/cetz:0.4.2"

#cetz.canvas({
  import cetz.draw: *

  // 滞留记录实体矩形
  rect((0, 0), (4, 2), stroke: 2pt)
  content((2, 1), text(size: 14pt, weight: "bold")[滞留记录])

  // 属性椭圆
  let attrs = (
    ("记录ID", (2, 3)),
    ("滞留开始时间", (-2, 1.5)),
    ("提醒次数", (6, 1.5)),
    ("最后提醒时间", (-2, 0)),
    ("滞留费用", (6, 0)),
    ("处理状态", (2, -1.5)),
  )

  for (name, pos) in attrs {
    circle(pos, radius: 0.9, stroke: 1pt)
    content(pos, text(size: 9pt)[#name])
    line((2, 1), pos, stroke: 0.5pt)
  }

  // 主键标记（下划线）
  line((1.5, 3), (2.5, 3), stroke: 1pt)
})
