// 用户实体 ER 图 - 使用 fletcher 包
#import "@preview/fletcher:0.5.8" as fletcher: diagram, edge, node
#import fletcher.shapes: ellipse, rect

#diagram(
  spacing: (15pt, 12pt),
  node-stroke: 1pt,
  edge-stroke: 0.8pt,

  // 中央实体
  node(
    (0, 0),
    text(size: 12pt, weight: "bold")[用户 User],
    stroke: 2pt,
    shape: rect,
    width: 3cm,
    height: 1.2cm,
    name: <entity>,
  ),

  // 属性节点 - 主键
  node((0, -2), text(size: 9pt)[用户ID], stroke: 1pt, shape: ellipse, inset: 6pt, name: <id>),

  // 属性节点 - 普通属性
  node((-2, -1), text(size: 9pt)[用户名], stroke: 1pt, shape: ellipse, inset: 6pt, name: <username>),
  node((2, -1), text(size: 9pt)[密码], stroke: 1pt, shape: ellipse, inset: 6pt, name: <password>),
  node((-2.5, 0), text(size: 9pt)[真实姓名], stroke: 1pt, shape: ellipse, inset: 6pt, name: <realname>),
  node((2.5, 0), text(size: 9pt)[手机号], stroke: 1pt, shape: ellipse, inset: 6pt, name: <phone>),
  node((-2, 1), text(size: 9pt)[身份证号], stroke: 1pt, shape: ellipse, inset: 6pt, name: <idcard>),
  node((2, 1), text(size: 9pt)[邮箱], stroke: 1pt, shape: ellipse, inset: 6pt, name: <email>),
  node((0, 2), text(size: 9pt)[角色], stroke: 1pt, shape: ellipse, inset: 6pt, name: <role>),

  // 连接线
  edge(<entity>, <id>),
  edge(<entity>, <username>),
  edge(<entity>, <password>),
  edge(<entity>, <realname>),
  edge(<entity>, <phone>),
  edge(<entity>, <idcard>),
  edge(<entity>, <email>),
  edge(<entity>, <role>),
)
