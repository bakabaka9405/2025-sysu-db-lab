#import "@preview/cetz:0.4.2"

#cetz.canvas({
  import cetz.draw: *

  // 用户实体矩形
  rect((0, 0), (4, 2), stroke: 2pt)
  content((2, 1), text(size: 14pt, weight: "bold")[用户 User])

  // 属性椭圆
  let attrs = (
    ("用户ID", (2, 3.5)),
    ("用户名", (-1, 2.5)),
    ("密码", (5, 2.5)),
    ("真实姓名", (-1.5, 1)),
    ("手机号", (5.5, 1)),
    ("身份证号", (-1.5, -0.5)),
    ("邮箱", (5.5, -0.5)),
    ("角色", (2, -1.5)),
  )

  for (name, pos) in attrs {
    circle(pos, radius: 0.8, stroke: 1pt)
    content(pos, text(size: 9pt)[#name])
    line((2, 1), pos, stroke: 0.5pt)
  }

  // 主键标记（下划线）
  line((1.5, 3.5), (2.5, 3.5), stroke: 1pt)
})
